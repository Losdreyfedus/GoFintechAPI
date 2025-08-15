package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

// ReplicationMode defines the replication strategy
type ReplicationMode string

const (
	ModeMasterSlave  ReplicationMode = "master_slave"
	ModeMasterMaster ReplicationMode = "master_master"
	ModeReadReplica  ReplicationMode = "read_replica"
)

// ReplicatedDB represents a database with replication support
type ReplicatedDB struct {
	master      *sql.DB
	replicas    []*sql.DB
	mode        ReplicationMode
	mu          sync.RWMutex
	healthCheck time.Duration
}

// NewReplicatedDB creates a new replicated database connection
func NewReplicatedDB(masterURL string, replicaURLs []string, mode ReplicationMode) (*ReplicatedDB, error) {
	// Connect to master
	master, err := sql.Open("mssql", masterURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master: %w", err)
	}

	// Test master connection
	if err := master.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping master: %w", err)
	}

	// Connect to replicas
	var replicas []*sql.DB
	for i, replicaURL := range replicaURLs {
		replica, err := sql.Open("mssql", replicaURL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to replica %d: %w", i, err)
		}

		// Test replica connection
		if err := replica.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping replica %d: %w", i, err)
		}

		replicas = append(replicas, replica)
	}

	return &ReplicatedDB{
		master:      master,
		replicas:    replicas,
		mode:        mode,
		healthCheck: 30 * time.Second,
	}, nil
}

// GetMaster returns the master database connection
func (rdb *ReplicatedDB) GetMaster() *sql.DB {
	return rdb.master
}

// GetReplica returns a healthy replica database connection
func (rdb *ReplicatedDB) GetReplica() (*sql.DB, error) {
	rdb.mu.RLock()
	defer rdb.mu.RUnlock()

	if len(rdb.replicas) == 0 {
		return nil, fmt.Errorf("no replicas available")
	}

	// Find healthy replica
	for _, replica := range rdb.replicas {
		if rdb.isHealthy(replica) {
			return replica, nil
		}
	}

	return nil, fmt.Errorf("no healthy replicas available")
}

// GetReplicas returns all replica connections
func (rdb *ReplicatedDB) GetReplicas() []*sql.DB {
	rdb.mu.RLock()
	defer rdb.mu.RUnlock()
	return rdb.replicas
}

// isHealthy checks if a database connection is healthy
func (rdb *ReplicatedDB) isHealthy(db *sql.DB) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.PingContext(ctx) == nil
}

// HealthCheck performs health check on all connections
func (rdb *ReplicatedDB) HealthCheck() map[string]bool {
	rdb.mu.RLock()
	defer rdb.mu.RUnlock()

	health := make(map[string]bool)

	// Check master
	health["master"] = rdb.isHealthy(rdb.master)

	// Check replicas
	for i, replica := range rdb.replicas {
		key := fmt.Sprintf("replica_%d", i)
		health[key] = rdb.isHealthy(replica)
	}

	return health
}

// StartHealthCheck starts periodic health checking
func (rdb *ReplicatedDB) StartHealthCheck(ctx context.Context) {
	ticker := time.NewTicker(rdb.healthCheck)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			health := rdb.HealthCheck()

			// Log health status
			for db, isHealthy := range health {
				if !isHealthy {
					fmt.Printf("Database %s is unhealthy\n", db)
				}
			}
		}
	}
}

// Close closes all database connections
func (rdb *ReplicatedDB) Close() error {
	var errors []error

	// Close master
	if err := rdb.master.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close master: %w", err))
	}

	// Close replicas
	for i, replica := range rdb.replicas {
		if err := replica.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close replica %d: %w", i, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing databases: %v", errors)
	}

	return nil
}

// GetStats returns replication statistics
func (rdb *ReplicatedDB) GetStats() map[string]interface{} {
	rdb.mu.RLock()
	defer rdb.mu.RUnlock()

	health := rdb.HealthCheck()
	stats := map[string]interface{}{
		"mode":           rdb.mode,
		"master_healthy": health["master"],
		"replicas_count": len(rdb.replicas),
		"health_status":  health,
	}

	return stats
}

