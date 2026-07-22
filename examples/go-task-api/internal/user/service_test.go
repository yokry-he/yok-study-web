package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type fakeRepository struct {
	createResult User
	createErr    error
	getResult    User
	getErr       error
	listResult   Page
	listErr      error
	updateResult User
	updateErr    error

	createCalls int
	getCalls    int
	listCalls   int
	updateCalls int

	createContext context.Context
	createParams  CreateParams
	getContext    context.Context
	getID         int64
	listContext   context.Context
	listFilter    ListFilter
	updateContext context.Context
	updateID      int64
	updateStatus  Status
	updateVersion int64
}

func (f *fakeRepository) Create(ctx context.Context, params CreateParams) (User, error) {
	f.createCalls++
	f.createContext = ctx
	f.createParams = params
	return f.createResult, f.createErr
}

func (f *fakeRepository) Get(ctx context.Context, id int64) (User, error) {
	f.getCalls++
	f.getContext = ctx
	f.getID = id
	return f.getResult, f.getErr
}

func (f *fakeRepository) List(ctx context.Context, filter ListFilter) (Page, error) {
	f.listCalls++
	f.listContext = ctx
	f.listFilter = filter
	return f.listResult, f.listErr
}

func (f *fakeRepository) UpdateStatus(ctx context.Context, id int64, status Status, expectedVersion int64) (User, error) {
	f.updateCalls++
	f.updateContext = ctx
	f.updateID = id
	f.updateStatus = status
	f.updateVersion = expectedVersion
	return f.updateResult, f.updateErr
}

func TestCreateNormalizesInputAndDefaultsStatus(t *testing.T) {
	t.Parallel()

	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("request"), "create-user")
	repo := &fakeRepository{createResult: User{ID: 7}}
	service := NewService(repo)

	got, err := service.Create(ctx, CreateInput{Name: "  张三  ", Email: "  USER@EXAMPLE.COM  "})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 7 {
		t.Fatalf("Create() ID = %d, want 7", got.ID)
	}
	if repo.createContext != ctx {
		t.Fatal("Create() did not forward the original context")
	}
	want := CreateParams{Name: "张三", Email: "user@example.com", Status: StatusActive}
	if repo.createParams != want {
		t.Fatalf("Create() params = %+v, want %+v", repo.createParams, want)
	}
}

func TestCreateValidatesNameByUnicodeRuneCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantValid bool
	}{
		{name: "two Chinese runes", input: "李雷", wantValid: true},
		{name: "64 Chinese runes", input: strings.Repeat("界", 64), wantValid: true},
		{name: "one rune", input: "李", wantValid: false},
		{name: "65 runes", input: strings.Repeat("界", 65), wantValid: false},
		{name: "empty after trim", input: " \t\n ", wantValid: false},
		{name: "invalid UTF-8", input: string([]byte{0xff, 0xfe}), wantValid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{}
			service := NewService(repo)
			_, err := service.Create(context.Background(), CreateInput{Name: tt.input, Email: "user@example.com"})
			if tt.wantValid {
				if err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				return
			}
			assertValidationField(t, err, "name")
			if repo.createCalls != 0 {
				t.Fatalf("repository Create() calls = %d, want 0", repo.createCalls)
			}
		})
	}
}

func TestCreateValidatesNameCharactersWithoutRejectingUsefulFormatRunes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantValid bool
	}{
		{name: "only zero-width format runes", input: "\u200b\u200d", wantValid: false},
		{name: "only variation selectors", input: "\ufe0e\ufe0f", wantValid: false},
		{name: "only combining marks", input: "\u0301\u0308", wantValid: false},
		{name: "embedded newline", input: "张\n三", wantValid: false},
		{name: "embedded NUL", input: "张\x00三", wantValid: false},
		{name: "Chinese name", input: "张三", wantValid: true},
		{name: "hyphenated name", input: "Anne-Marie", wantValid: true},
		{name: "combining character", input: "A\u0301B", wantValid: true},
		{name: "normal internal space", input: "李 雷", wantValid: true},
		{name: "emoji joined with ZWJ", input: "👩\u200d💻", wantValid: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{}
			_, err := NewService(repo).Create(context.Background(), CreateInput{
				Name:  tt.input,
				Email: "user@example.com",
			})
			if tt.wantValid {
				if err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				if repo.createParams.Name != tt.input {
					t.Fatalf("repository name = %q, want %q", repo.createParams.Name, tt.input)
				}
				return
			}
			assertValidationField(t, err, "name")
			if repo.createCalls != 0 {
				t.Fatalf("repository Create() calls = %d, want 0", repo.createCalls)
			}
		})
	}
}

func TestCreateRejectsNonCanonicalEmailForms(t *testing.T) {
	t.Parallel()

	overlong := strings.Repeat("a", 243) + "@example.com"
	tests := []struct {
		name  string
		email string
	}{
		{name: "empty", email: "   "},
		{name: "invalid", email: "not-an-email"},
		{name: "display name", email: "Alice <alice@example.com>"},
		{name: "angle address", email: "<alice@example.com>"},
		{name: "multiple addresses", email: "a@example.com, b@example.com"},
		{name: "over 254 bytes", email: overlong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{}
			service := NewService(repo)
			_, err := service.Create(context.Background(), CreateInput{Name: "张三", Email: tt.email})
			assertValidationField(t, err, "email")
			if repo.createCalls != 0 {
				t.Fatalf("repository Create() calls = %d, want 0", repo.createCalls)
			}
		})
	}
}

func TestCreatePropagatesRepositoryErrorsUnchanged(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
	}{
		{name: "email conflict", err: fmt.Errorf("insert user: %w", ErrEmailConflict)},
		{name: "context canceled", err: context.Canceled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{createErr: tt.err}
			_, gotErr := NewService(repo).Create(context.Background(), CreateInput{Name: "张三", Email: "user@example.com"})
			if gotErr != tt.err {
				t.Fatalf("Create() error = %v, want original error %v", gotErr, tt.err)
			}
			if !errors.Is(gotErr, tt.err) {
				t.Fatalf("errors.Is(%v, %v) = false", gotErr, tt.err)
			}
		})
	}
}

func TestGetForwardsArgumentsAndRepositoryErrors(t *testing.T) {
	t.Parallel()

	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("request"), "get-user")
	repoErr := fmt.Errorf("select user: %w", ErrNotFound)
	repo := &fakeRepository{getErr: repoErr}

	_, err := NewService(repo).Get(ctx, 42)
	if err != repoErr {
		t.Fatalf("Get() error = %v, want original error %v", err, repoErr)
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get() error %v does not preserve ErrNotFound", err)
	}
	if repo.getContext != ctx || repo.getID != 42 {
		t.Fatalf("Get() forwarded context/id = %v/%d, want original/42", repo.getContext == ctx, repo.getID)
	}
}

func TestListAppliesDefaultsAndConvertsToLimitOffset(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{listResult: Page{Items: []User{{ID: 1}}, Total: 37}}
	got, err := NewService(repo).List(context.Background(), ListInput{})
	if err != nil {
		t.Fatal(err)
	}
	if repo.listFilter.Limit != 20 || repo.listFilter.Offset != 0 || repo.listFilter.Status != nil {
		t.Fatalf("List() filter = %+v, want limit=20 offset=0 no status", repo.listFilter)
	}
	if got.Page != 1 || got.PageSize != 20 || got.Total != 37 || len(got.Items) != 1 {
		t.Fatalf("List() page = %+v, want page=1 pageSize=20 total=37 with one item", got)
	}
}

func TestListConvertsExplicitPageAndStatusFilter(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{listResult: Page{Total: 250}}
	got, err := NewService(repo).List(context.Background(), ListInput{
		Page:     3,
		PageSize: 100,
		Status:   StatusDisabled,
	})
	if err != nil {
		t.Fatal(err)
	}
	filter := repo.listFilter
	if filter.Limit != 100 || filter.Offset != 200 {
		t.Fatalf("List() limit/offset = %d/%d, want 100/200", filter.Limit, filter.Offset)
	}
	if filter.Status == nil || *filter.Status != StatusDisabled {
		t.Fatalf("List() status = %v, want DISABLED", filter.Status)
	}
	if got.Page != 3 || got.PageSize != 100 || got.Total != 250 {
		t.Fatalf("List() page metadata = %+v, want page=3 pageSize=100 total=250", got)
	}
}

func TestListValidatesPaginationAndStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input ListInput
		field string
	}{
		{name: "negative page", input: ListInput{Page: -1}, field: "page"},
		{name: "negative page size", input: ListInput{PageSize: -1}, field: "pageSize"},
		{name: "page size over maximum", input: ListInput{PageSize: 101}, field: "pageSize"},
		{name: "unknown status", input: ListInput{Status: Status("PENDING")}, field: "status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{}
			_, err := NewService(repo).List(context.Background(), tt.input)
			assertValidationField(t, err, tt.field)
			if repo.listCalls != 0 {
				t.Fatalf("repository List() calls = %d, want 0", repo.listCalls)
			}
		})
	}
}

func TestListProtectsOffsetAtNativeIntBoundary(t *testing.T) {
	t.Parallel()

	const pageSize = 100
	maxInt := int(^uint(0) >> 1)
	maxSafePage := maxInt/pageSize + 1
	firstOverflowPage := maxSafePage + 1

	t.Run("largest page with representable offset", func(t *testing.T) {
		repo := &fakeRepository{}
		got, err := NewService(repo).List(context.Background(), ListInput{
			Page:     maxSafePage,
			PageSize: pageSize,
		})
		if err != nil {
			t.Fatal(err)
		}
		wantOffset := (maxSafePage - 1) * pageSize
		if repo.listCalls != 1 || repo.listFilter.Offset != wantOffset {
			t.Fatalf("repository List() calls/offset = %d/%d, want 1/%d", repo.listCalls, repo.listFilter.Offset, wantOffset)
		}
		if got.Page != maxSafePage || got.PageSize != pageSize {
			t.Fatalf("List() page metadata = %+v, want page=%d pageSize=%d", got, maxSafePage, pageSize)
		}
	})

	t.Run("first page whose offset overflows", func(t *testing.T) {
		repo := &fakeRepository{}
		_, err := NewService(repo).List(context.Background(), ListInput{
			Page:     firstOverflowPage,
			PageSize: pageSize,
		})
		assertValidationField(t, err, "page")
		if repo.listCalls != 0 {
			t.Fatalf("repository List() calls = %d, want 0", repo.listCalls)
		}
	})
}

func TestListAcceptsBothStatusesAndNoFilter(t *testing.T) {
	t.Parallel()

	for _, status := range []Status{"", StatusActive, StatusDisabled} {
		status := status
		t.Run(string(status), func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{}
			_, err := NewService(repo).List(context.Background(), ListInput{Status: status})
			if err != nil {
				t.Fatal(err)
			}
			if status == "" && repo.listFilter.Status != nil {
				t.Fatalf("empty status produced filter %v", *repo.listFilter.Status)
			}
			if status != "" && (repo.listFilter.Status == nil || *repo.listFilter.Status != status) {
				t.Fatalf("status filter = %v, want %s", repo.listFilter.Status, status)
			}
		})
	}
}

func TestListPropagatesRepositoryErrorUnchanged(t *testing.T) {
	t.Parallel()

	repoErr := context.Canceled
	repo := &fakeRepository{listErr: repoErr}
	_, err := NewService(repo).List(context.Background(), ListInput{})
	if err != repoErr || !errors.Is(err, context.Canceled) {
		t.Fatalf("List() error = %v, want original context cancellation", err)
	}
}

func TestChangeStatusForwardsZeroExpectedVersion(t *testing.T) {
	t.Parallel()

	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("request"), "change-status")
	repo := &fakeRepository{updateResult: User{ID: 9, Status: StatusDisabled, Version: 1}}
	got, err := NewService(repo).ChangeStatus(ctx, 9, ChangeStatusInput{
		Status:          StatusDisabled,
		ExpectedVersion: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 9 || repo.updateContext != ctx || repo.updateID != 9 || repo.updateStatus != StatusDisabled || repo.updateVersion != 0 {
		t.Fatalf("ChangeStatus() result/args = %+v, context=%v id=%d status=%s version=%d", got, repo.updateContext == ctx, repo.updateID, repo.updateStatus, repo.updateVersion)
	}
}

func TestChangeStatusValidatesStatusAndVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input ChangeStatusInput
		field string
	}{
		{name: "empty status", input: ChangeStatusInput{ExpectedVersion: 0}, field: "status"},
		{name: "unknown status", input: ChangeStatusInput{Status: Status("LOCKED"), ExpectedVersion: 0}, field: "status"},
		{name: "negative expected version", input: ChangeStatusInput{Status: StatusActive, ExpectedVersion: -1}, field: "expectedVersion"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{}
			_, err := NewService(repo).ChangeStatus(context.Background(), 1, tt.input)
			assertValidationField(t, err, tt.field)
			if repo.updateCalls != 0 {
				t.Fatalf("repository UpdateStatus() calls = %d, want 0", repo.updateCalls)
			}
		})
	}
}

func TestChangeStatusPropagatesRepositoryErrorsUnchanged(t *testing.T) {
	t.Parallel()

	tests := []error{
		fmt.Errorf("update status: %w", ErrVersionConflict),
		fmt.Errorf("update status: %w", ErrNotFound),
		context.Canceled,
	}
	for _, repoErr := range tests {
		repoErr := repoErr
		t.Run(repoErr.Error(), func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{updateErr: repoErr}
			_, err := NewService(repo).ChangeStatus(context.Background(), 1, ChangeStatusInput{Status: StatusActive, ExpectedVersion: 0})
			if err != repoErr {
				t.Fatalf("ChangeStatus() error = %v, want original error %v", err, repoErr)
			}
			if !errors.Is(err, repoErr) {
				t.Fatalf("errors.Is(%v, %v) = false", err, repoErr)
			}
		})
	}
}

func TestNilRepositoryReturnsStableError(t *testing.T) {
	t.Parallel()

	service := NewService(nil)
	tests := []struct {
		name string
		call func() error
	}{
		{name: "Create", call: func() error {
			_, err := service.Create(context.Background(), CreateInput{Name: "张三", Email: "user@example.com"})
			return err
		}},
		{name: "Get", call: func() error {
			_, err := service.Get(context.Background(), 1)
			return err
		}},
		{name: "List", call: func() error {
			_, err := service.List(context.Background(), ListInput{})
			return err
		}},
		{name: "ChangeStatus", call: func() error {
			_, err := service.ChangeStatus(context.Background(), 1, ChangeStatusInput{Status: StatusActive})
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); !errors.Is(err, ErrNilRepository) {
				t.Fatalf("%s() error = %v, want ErrNilRepository", tt.name, err)
			}
		})
	}
}

func TestTypedNilRepositoryReturnsStableError(t *testing.T) {
	t.Parallel()

	var repo *fakeRepository
	service := NewService(repo)
	_, err := service.Get(context.Background(), 1)
	if !errors.Is(err, ErrNilRepository) {
		t.Fatalf("Get() error = %v, want ErrNilRepository", err)
	}
}

func TestDomainSentinelsSupportErrorsIs(t *testing.T) {
	t.Parallel()

	for _, sentinel := range []error{ErrValidation, ErrNotFound, ErrEmailConflict, ErrVersionConflict, ErrNilRepository} {
		wrapped := fmt.Errorf("repository: %w", sentinel)
		if !errors.Is(wrapped, sentinel) {
			t.Fatalf("errors.Is(%v, %v) = false", wrapped, sentinel)
		}
	}
}

func assertValidationField(t *testing.T, err error, field string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected validation error for %q, got nil", field)
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("error %v does not match ErrValidation", err)
	}
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error %T does not expose *ValidationError", err)
	}
	if _, ok := validationErr.Fields[field]; !ok {
		t.Fatalf("ValidationError fields = %v, want key %q", validationErr.Fields, field)
	}
}
