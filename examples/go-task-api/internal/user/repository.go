package user

import "context"

type CreateParams struct {
	Name   string
	Email  string
	Status Status
}

// ListFilter 使用数据库可直接消费的 limit/offset 语义；Status 为 nil 时不筛选状态。
type ListFilter struct {
	Status *Status
	Limit  int
	Offset int
}

type Repository interface {
	Create(context.Context, CreateParams) (User, error)
	Get(context.Context, int64) (User, error)
	List(context.Context, ListFilter) (Page, error)
	UpdateStatus(context.Context, int64, Status, int64) (User, error)
}
