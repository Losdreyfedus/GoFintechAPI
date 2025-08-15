package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheStrategy defines caching strategies
type CacheStrategy string

const (
	StrategyWriteThrough CacheStrategy = "write_through"
	StrategyWriteBehind  CacheStrategy = "write_behind"
	StrategyWriteAround  CacheStrategy = "write_around"
	StrategyCacheAside   CacheStrategy = "cache_aside"
)

// AdvancedCache represents advanced caching with multiple strategies
type AdvancedCache struct {
	client   *redis.Client
	strategy CacheStrategy
}

// NewAdvancedCache creates a new advanced cache instance
func NewAdvancedCache(client *redis.Client, strategy CacheStrategy) *AdvancedCache {
	return &AdvancedCache{
		client:   client,
		strategy: strategy,
	}
}

// CacheItem represents a cached item with metadata
type CacheItem struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	Expiration time.Time   `json:"expiration"`
	Version    int64       `json:"version"`
	Tags       []string    `json:"tags,omitempty"`
}

// SetWithStrategy sets a value with specified caching strategy
func (c *AdvancedCache) SetWithStrategy(ctx context.Context, key string, value interface{}, ttl time.Duration, tags []string) error {
	item := &CacheItem{
		Key:        key,
		Value:      value,
		Expiration: time.Now().Add(ttl),
		Version:    time.Now().UnixNano(),
		Tags:       tags,
	}

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal cache item: %w", err)
	}

	// Set cache with TTL
	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	// Store tags for invalidation
	if len(tags) > 0 {
		for _, tag := range tags {
			tagKey := fmt.Sprintf("tag:%s", tag)
			c.client.SAdd(ctx, tagKey, key)
			c.client.Expire(ctx, tagKey, ttl)
		}
	}

	// Log cache set operation
	fmt.Printf("Cache set successfully: key=%s, strategy=%s, ttl=%v, tags=%v\n", key, c.strategy, ttl, tags)

	return nil
}

// Get retrieves a value from cache
func (c *AdvancedCache) Get(ctx context.Context, key string) (interface{}, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	var item CacheItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache item: %w", err)
	}

	// Check expiration
	if time.Now().After(item.Expiration) {
		c.client.Del(ctx, key)
		return nil, fmt.Errorf("key expired: %s", key)
	}

	return item.Value, nil
}

// InvalidateByTag invalidates all keys with a specific tag
func (c *AdvancedCache) InvalidateByTag(ctx context.Context, tag string) error {
	tagKey := fmt.Sprintf("tag:%s", tag)
	keys, err := c.client.SMembers(ctx, tagKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get tag keys: %w", err)
	}

	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete tagged keys: %w", err)
		}
	}

	// Delete tag set
	c.client.Del(ctx, tagKey)

	// Log cache invalidation
	fmt.Printf("Cache invalidated by tag: tag=%s, count=%d\n", tag, len(keys))

	return nil
}

// WarmUp preloads popular data into cache
func (c *AdvancedCache) WarmUp(ctx context.Context, warmUpFunc func() (map[string]interface{}, error), ttl time.Duration) error {
	data, err := warmUpFunc()
	if err != nil {
		return fmt.Errorf("failed to warm up cache: %w", err)
	}

	for key, value := range data {
		if err := c.SetWithStrategy(ctx, key, value, ttl, []string{"warmup"}); err != nil {
			fmt.Printf("Failed to warm up cache key: %s, error: %v\n", key, err)
		}
	}

	// Log cache warm-up completion
	fmt.Printf("Cache warm-up completed: keys_count=%d, ttl=%v\n", len(data), ttl)

	return nil
}

// GetStats returns cache statistics
func (c *AdvancedCache) GetStats(ctx context.Context) (map[string]interface{}, error) {
	info, err := c.client.Info(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache info: %w", err)
	}

	stats := map[string]interface{}{
		"strategy": c.strategy,
		"info":     info,
	}

	return stats, nil
}

// Clear clears all cache data
func (c *AdvancedCache) Clear(ctx context.Context) error {
	if err := c.client.FlushDB(ctx).Err(); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	// Log cache clear operation
	fmt.Println("Cache cleared successfully")
	return nil
}
