package task

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/user"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
	maxTitleRunes   = 128
)

type UserReader interface {
	Get(context.Context, int64) (user.User, error)
}

type Service struct {
	repo  Repository
	users UserReader
	now   func() time.Time
}

type Option func(*Service)

func WithNow(now func() time.Time) Option {
	return func(service *Service) {
		if now != nil {
			service.now = now
		}
	}
}

func NewService(repo Repository, users UserReader, options ...Option) *Service {
	service := &Service{
		repo:  repo,
		users: users,
		now:   time.Now,
	}
	for _, option := range options {
		if option != nil {
			option(service)
		}
	}
	return service
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Task, error) {
	repo, err := s.repository()
	if err != nil {
		return Task{}, err
	}
	users, err := s.userReader()
	if err != nil {
		return Task{}, err
	}
	if input.OwnerID <= 0 {
		return Task{}, invalidField("ownerId", "负责人 ID 必须是正整数")
	}
	title, err := normalizeTitle(input.Title)
	if err != nil {
		return Task{}, err
	}
	dueAt, err := normalizeDueAt(input.DueAt, s.currentTime())
	if err != nil {
		return Task{}, err
	}

	owner, err := users.Get(ctx, input.OwnerID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return Task{}, ErrOwnerNotFound
		}
		return Task{}, err
	}
	if owner.Status != user.StatusActive {
		return Task{}, ErrOwnerDisabled
	}

	return repo.Create(ctx, CreateParams{
		OwnerID:     input.OwnerID,
		Title:       title,
		Description: normalizeDescription(input.Description),
		Status:      StatusTodo,
		DueAt:       dueAt,
	})
}

func (s *Service) Get(ctx context.Context, id int64) (Task, error) {
	repo, err := s.repository()
	if err != nil {
		return Task{}, err
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

	var ownerID *int64
	if input.OwnerID < 0 {
		return Page{}, invalidField("ownerId", "负责人 ID 必须是正整数")
	}
	if input.OwnerID > 0 {
		value := input.OwnerID
		ownerID = &value
	}
	var status *Status
	if input.Status != "" {
		if !input.Status.valid() {
			return Page{}, invalidField("status", "任务状态必须是 TODO、DOING、DONE 或 CANCELLED")
		}
		value := input.Status
		status = &value
	}

	maxInt := int(^uint(0) >> 1)
	if page > 1 && page-1 > maxInt/pageSize {
		return Page{}, invalidField("page", "页码过大")
	}
	result, err := repo.List(ctx, ListFilter{
		OwnerID: ownerID,
		Status:  status,
		Limit:   pageSize,
		Offset:  (page - 1) * pageSize,
	})
	if err != nil {
		return Page{}, err
	}
	result.Page = page
	result.PageSize = pageSize
	return result, nil
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (Task, error) {
	repo, err := s.repository()
	if err != nil {
		return Task{}, err
	}
	if input.ExpectedVersion < 0 {
		return Task{}, invalidField("expectedVersion", "预期版本不能为负数")
	}
	title, err := normalizeTitle(input.Title)
	if err != nil {
		return Task{}, err
	}
	dueAt, err := normalizeDueAt(input.DueAt, s.currentTime())
	if err != nil {
		return Task{}, err
	}
	return repo.Update(ctx, id, UpdateParams{
		Title:       title,
		Description: normalizeDescription(input.Description),
		DueAt:       dueAt,
	}, input.ExpectedVersion)
}

func (s *Service) ChangeStatus(ctx context.Context, id int64, input ChangeStatusInput) (Task, error) {
	repo, err := s.repository()
	if err != nil {
		return Task{}, err
	}
	if !input.Status.valid() {
		return Task{}, invalidField("status", "任务状态必须是 TODO、DOING、DONE 或 CANCELLED")
	}
	if input.ExpectedVersion < 0 {
		return Task{}, invalidField("expectedVersion", "预期版本不能为负数")
	}

	current, err := repo.Get(ctx, id)
	if err != nil {
		return Task{}, err
	}
	if current.Version != input.ExpectedVersion {
		return Task{}, ErrVersionConflict
	}
	if !allowedTransitions[current.Status][input.Status] {
		return Task{}, ErrInvalidTransition
	}
	return repo.UpdateStatus(ctx, id, input.Status, input.ExpectedVersion)
}

func (s *Service) Delete(ctx context.Context, id int64, input DeleteInput) error {
	repo, err := s.repository()
	if err != nil {
		return err
	}
	if input.ExpectedVersion < 0 {
		return invalidField("expectedVersion", "预期版本不能为负数")
	}
	return repo.Delete(ctx, id, input.ExpectedVersion)
}

func normalizeTitle(raw string) (string, error) {
	title := strings.TrimSpace(raw)
	if !utf8.ValidString(title) {
		return "", invalidField("title", "标题必须是有效的 UTF-8 文本")
	}
	length := utf8.RuneCountInString(title)
	if length < 1 || length > maxTitleRunes {
		return "", invalidField("title", "标题长度必须为 1 到 128 个字符")
	}
	hasVisibleCharacter := false
	for _, char := range title {
		if unicode.IsControl(char) {
			return "", invalidField("title", "标题不能包含控制字符")
		}
		if unicode.IsLetter(char) || unicode.IsNumber(char) || unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasVisibleCharacter = true
		}
	}
	if !hasVisibleCharacter {
		return "", invalidField("title", "标题必须包含字母、数字、标点或符号")
	}
	return title, nil
}

func normalizeDescription(raw *string) *string {
	if raw == nil {
		return nil
	}
	value := strings.TrimSpace(*raw)
	return &value
}

func normalizeDueAt(raw *time.Time, now time.Time) (*time.Time, error) {
	if raw == nil {
		return nil, nil
	}
	if raw.Before(now) {
		return nil, invalidField("dueAt", "截止时间不能早于当前时间")
	}
	value := *raw
	return &value, nil
}

func (s *Service) currentTime() time.Time {
	if s == nil || s.now == nil {
		return time.Now()
	}
	return s.now()
}

func (s *Service) repository() (Repository, error) {
	if s == nil || interfaceIsNil(s.repo) {
		return nil, ErrNilRepository
	}
	return s.repo, nil
}

func (s *Service) userReader() (UserReader, error) {
	if s == nil || interfaceIsNil(s.users) {
		return nil, ErrNilUserReader
	}
	return s.users, nil
}

func interfaceIsNil(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}

func invalidField(field, message string) error {
	return &ValidationError{Fields: FieldErrors{field: message}}
}
