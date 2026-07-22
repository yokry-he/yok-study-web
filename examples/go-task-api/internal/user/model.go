package user

import (
	"errors"
	"time"
)

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusDisabled Status = "DISABLED"
)

var (
	ErrValidation      = errors.New("user: validation failed")
	ErrNotFound        = errors.New("user: not found")
	ErrEmailConflict   = errors.New("user: email conflict")
	ErrVersionConflict = errors.New("user: version conflict")
	ErrNilRepository   = errors.New("user: repository is nil")
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    Status    `json:"status"`
	Version   int64     `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateInput struct {
	Name  string
	Email string
}

type ListInput struct {
	Page     int
	PageSize int
	Status   Status
}

type ChangeStatusInput struct {
	Status          Status
	ExpectedVersion int64
}

type Page struct {
	Items    []User `json:"items"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Total    int64  `json:"total"`
}

type FieldErrors map[string]string

// ValidationError 为 HTTP 层提供稳定的字段级错误信息，同时通过 Unwrap 支持 errors.Is。
type ValidationError struct {
	Fields FieldErrors
}

func (e *ValidationError) Error() string {
	return ErrValidation.Error()
}

func (e *ValidationError) Unwrap() error {
	return ErrValidation
}
