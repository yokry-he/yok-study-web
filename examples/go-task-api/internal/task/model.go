package task

import (
	"errors"
	"time"
)

type Status string

const (
	StatusTodo      Status = "TODO"
	StatusDoing     Status = "DOING"
	StatusDone      Status = "DONE"
	StatusCancelled Status = "CANCELLED"
)

var allowedTransitions = map[Status]map[Status]bool{
	StatusTodo: {
		StatusDoing:     true,
		StatusCancelled: true,
	},
	StatusDoing: {
		StatusTodo:      true,
		StatusDone:      true,
		StatusCancelled: true,
	},
	StatusDone:      {},
	StatusCancelled: {},
}

var (
	ErrValidation        = errors.New("task: validation failed")
	ErrNotFound          = errors.New("task: not found")
	ErrVersionConflict   = errors.New("task: version conflict")
	ErrInvalidTransition = errors.New("task: invalid status transition")
	ErrOwnerNotFound     = errors.New("task: owner not found")
	ErrOwnerDisabled     = errors.New("task: owner is disabled")
	ErrNilRepository     = errors.New("task: repository is nil")
	ErrNilUserReader     = errors.New("task: user reader is nil")
)

type Task struct {
	ID          int64      `json:"id"`
	OwnerID     int64      `json:"ownerId"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      Status     `json:"status"`
	DueAt       *time.Time `json:"dueAt"`
	Version     int64      `json:"version"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type CreateInput struct {
	OwnerID     int64
	Title       string
	Description *string
	DueAt       *time.Time
}

type UpdateInput struct {
	Title           string
	Description     *string
	DueAt           *time.Time
	ExpectedVersion int64
}

type ChangeStatusInput struct {
	Status          Status
	ExpectedVersion int64
}

type DeleteInput struct {
	ExpectedVersion int64
}

type ListInput struct {
	Page     int
	PageSize int
	OwnerID  int64
	Status   Status
}

type Page struct {
	Items    []Task `json:"items"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Total    int64  `json:"total"`
}

type FieldErrors map[string]string

// ValidationError 保留字段级错误，并通过 Unwrap 支持稳定的 errors.Is 判断。
type ValidationError struct {
	Fields FieldErrors
}

func (e *ValidationError) Error() string {
	return ErrValidation.Error()
}

func (e *ValidationError) Unwrap() error {
	return ErrValidation
}

func (s Status) valid() bool {
	_, ok := allowedTransitions[s]
	return ok
}
