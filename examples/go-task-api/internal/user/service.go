package user

import (
	"context"
	"net/mail"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (User, error) {
	repo, err := s.repository()
	if err != nil {
		return User{}, err
	}

	name := strings.TrimSpace(input.Name)
	if !utf8.ValidString(name) {
		return User{}, invalidField("name", "姓名必须是有效的 UTF-8 文本")
	}
	nameLength := utf8.RuneCountInString(name)
	if nameLength < 2 || nameLength > 64 {
		return User{}, invalidField("name", "姓名长度必须为 2 到 64 个字符")
	}
	hasBaseCharacter := false
	for _, r := range name {
		if unicode.IsControl(r) {
			return User{}, invalidField("name", "姓名不能包含控制字符")
		}
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsPunct(r) || unicode.IsSymbol(r) {
			hasBaseCharacter = true
		}
	}
	if !hasBaseCharacter {
		return User{}, invalidField("name", "姓名必须包含字母、数字、标点或符号")
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))
	if len(email) == 0 || len(email) > 254 {
		return User{}, invalidField("email", "邮箱长度必须为 1 到 254 个字节")
	}
	address, parseErr := mail.ParseAddress(email)
	if parseErr != nil || address.Name != "" || address.Address != email {
		return User{}, invalidField("email", "邮箱格式不正确")
	}

	return repo.Create(ctx, CreateParams{
		Name:   name,
		Email:  email,
		Status: StatusActive,
	})
}

func (s *Service) Get(ctx context.Context, id int64) (User, error) {
	repo, err := s.repository()
	if err != nil {
		return User{}, err
	}
	return repo.Get(ctx, id)
}

func (s *Service) List(ctx context.Context, input ListInput) (Page, error) {
	repo, err := s.repository()
	if err != nil {
		return Page{}, err
	}

	page := input.Page
	if page == 0 {
		page = defaultPage
	}
	if page < 0 {
		return Page{}, invalidField("page", "页码不能为负数")
	}

	pageSize := input.PageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	if pageSize < 0 || pageSize > maxPageSize {
		return Page{}, invalidField("pageSize", "每页数量必须为 1 到 100")
	}

	var status *Status
	if input.Status != "" {
		if !input.Status.valid() {
			return Page{}, invalidField("status", "用户状态必须是 ACTIVE 或 DISABLED")
		}
		value := input.Status
		status = &value
	}

	maxInt := int(^uint(0) >> 1)
	if page > 1 && page-1 > maxInt/pageSize {
		return Page{}, invalidField("page", "页码过大")
	}
	result, err := repo.List(ctx, ListFilter{
		Status: status,
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return Page{}, err
	}
	result.Page = page
	result.PageSize = pageSize
	return result, nil
}

func (s *Service) ChangeStatus(ctx context.Context, id int64, input ChangeStatusInput) (User, error) {
	repo, err := s.repository()
	if err != nil {
		return User{}, err
	}
	if !input.Status.valid() {
		return User{}, invalidField("status", "用户状态必须是 ACTIVE 或 DISABLED")
	}
	if input.ExpectedVersion < 0 {
		return User{}, invalidField("expectedVersion", "预期版本不能为负数")
	}
	return repo.UpdateStatus(ctx, id, input.Status, input.ExpectedVersion)
}

func (s *Service) repository() (Repository, error) {
	if s == nil || repositoryIsNil(s.repo) {
		return nil, ErrNilRepository
	}
	return s.repo, nil
}

func repositoryIsNil(repo Repository) bool {
	if repo == nil {
		return true
	}
	value := reflect.ValueOf(repo)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func (s Status) valid() bool {
	return s == StatusActive || s == StatusDisabled
}

func invalidField(field, message string) error {
	return &ValidationError{Fields: FieldErrors{field: message}}
}
