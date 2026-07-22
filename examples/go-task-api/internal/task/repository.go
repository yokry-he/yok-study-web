package task

import (
	"context"
	"time"
)

type CreateParams struct {
	OwnerID     int64
	Title       string
	Description *string
	Status      Status
	DueAt       *time.Time
}

// ListFilter 使用数据库可直接消费的筛选和分页语义；nil 表示不应用对应筛选。
type ListFilter struct {
	OwnerID *int64
	Status  *Status
	Limit   int
	Offset  int
}

// UpdateParams 故意不包含 OwnerID 和 Status，确保 PUT 无法绕过负责人和状态机规则。
type UpdateParams struct {
	Title       string
	Description *string
	DueAt       *time.Time
}

type Repository interface {
	Create(context.Context, CreateParams) (Task, error)
	Get(context.Context, int64) (Task, error)
	List(context.Context, ListFilter) (Page, error)
	Update(context.Context, int64, UpdateParams, int64) (Task, error)
	UpdateStatus(context.Context, int64, Status, int64) (Task, error)
	Delete(context.Context, int64, int64) error
}
