package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"backend_path/pkg/logger"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

// ScheduledTask represents a scheduled task
type ScheduledTask struct {
	ID       string                      `json:"id"`
	Name     string                      `json:"name"`
	CronExpr string                      `json:"cron_expr"`
	Handler  func(context.Context) error `json:"-"`
	LastRun  time.Time                   `json:"last_run"`
	NextRun  time.Time                   `json:"next_run"`
	Status   TaskStatus                  `json:"status"`
	Metadata map[string]interface{}      `json:"metadata"`
}

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusActive   TaskStatus = "active"
	TaskStatusPaused   TaskStatus = "paused"
	TaskStatusFailed   TaskStatus = "failed"
	TaskStatusComplete TaskStatus = "complete"
)

// Scheduler represents a task scheduler
type Scheduler struct {
	cron   *cron.Cron
	tasks  map[string]*ScheduledTask
	mu     sync.RWMutex
	logger zerolog.Logger
	ctx    context.Context
	cancel context.CancelFunc
}

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		cron:   cron.New(cron.WithSeconds()),
		tasks:  make(map[string]*ScheduledTask),
		logger: logger.GetLogger(),
		ctx:    ctx,
		cancel: cancel,
	}
}

// AddTask adds a new scheduled task
func (s *Scheduler) AddTask(task *ScheduledTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	// Parse cron expression
	schedule, err := cron.ParseStandard(task.CronExpr)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}

	// Calculate next run time
	task.NextRun = schedule.Next(time.Now())
	task.Status = TaskStatusActive

	// Add to cron scheduler
	_, err = s.cron.AddFunc(task.CronExpr, func() {
		s.executeTask(task)
	})
	if err != nil {
		return fmt.Errorf("failed to add task to cron: %w", err)
	}

	// Store task
	s.tasks[task.ID] = task

	s.logger.Info().Str("task_id", task.ID).Str("name", task.Name).Msg("Task scheduled")
	return nil
}

// RemoveTask removes a scheduled task
func (s *Scheduler) RemoveTask(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	// Remove from cron scheduler
	s.cron.Remove(cron.EntryID(0)) // This is a simplified version

	// Remove from tasks map
	delete(s.tasks, taskID)

	s.logger.Info().Str("task_id", taskID).Msg("Task removed")
	return nil
}

// PauseTask pauses a scheduled task
func (s *Scheduler) PauseTask(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	task.Status = TaskStatusPaused
	s.logger.Info().Str("task_id", taskID).Msg("Task paused")
	return nil
}

// ResumeTask resumes a paused task
func (s *Scheduler) ResumeTask(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	task.Status = TaskStatusActive
	s.logger.Info().Str("task_id", taskID).Msg("Task resumed")
	return nil
}

// GetTask returns a task by ID
func (s *Scheduler) GetTask(taskID string) (*ScheduledTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", taskID)
	}

	return task, nil
}

// GetAllTasks returns all scheduled tasks
func (s *Scheduler) GetAllTasks() []*ScheduledTask {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*ScheduledTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// executeTask executes a scheduled task
func (s *Scheduler) executeTask(task *ScheduledTask) {
	if task.Status != TaskStatusActive {
		return
	}

	s.logger.Info().Str("task_id", task.ID).Str("name", task.Name).Msg("Executing scheduled task")

	// Update last run time
	task.LastRun = time.Now()

	// Execute the task
	if err := task.Handler(s.ctx); err != nil {
		task.Status = TaskStatusFailed
		s.logger.Error().Err(err).Str("task_id", task.ID).Msg("Task execution failed")
		return
	}

	// Update next run time
	schedule, _ := cron.ParseStandard(task.CronExpr)
	task.NextRun = schedule.Next(time.Now())

	s.logger.Info().Str("task_id", task.ID).Msg("Task executed successfully")
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
	s.logger.Info().Msg("Scheduler started")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.cancel()
	s.cron.Stop()
	s.logger.Info().Msg("Scheduler stopped")
}

// TransactionScheduler specializes in scheduling financial transactions
type TransactionScheduler struct {
	*Scheduler
	transactionService TransactionService
}

// TransactionService interface for transaction operations
type TransactionService interface {
	ProcessScheduledTransaction(ctx context.Context, transactionID string) error
}

// NewTransactionScheduler creates a new transaction scheduler
func NewTransactionScheduler(transactionService TransactionService) *TransactionScheduler {
	return &TransactionScheduler{
		Scheduler:          NewScheduler(),
		transactionService: transactionService,
	}
}

// ScheduleTransaction schedules a transaction for future execution
func (ts *TransactionScheduler) ScheduleTransaction(transactionID string, cronExpr string) error {
	task := &ScheduledTask{
		ID:       fmt.Sprintf("transaction_%s", transactionID),
		Name:     fmt.Sprintf("Scheduled Transaction %s", transactionID),
		CronExpr: cronExpr,
		Handler: func(ctx context.Context) error {
			return ts.transactionService.ProcessScheduledTransaction(ctx, transactionID)
		},
		Metadata: map[string]interface{}{
			"transaction_id": transactionID,
			"type":           "transaction",
		},
	}

	return ts.AddTask(task)
}
