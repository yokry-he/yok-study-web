package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/yokry-he/yok-study-web/examples/go-task-api/internal/user"
)

var fixedNow = time.Date(2026, time.July, 21, 10, 0, 0, 0, time.UTC)

type fakeRepository struct {
	createResult       Task
	createErr          error
	getResult          Task
	getErr             error
	listResult         Page
	listErr            error
	updateResult       Task
	updateErr          error
	updateStatusResult Task
	updateStatusErr    error
	deleteErr          error

	createCalls       int
	getCalls          int
	listCalls         int
	updateCalls       int
	updateStatusCalls int
	deleteCalls       int

	createContext context.Context
	createParams  CreateParams
	getContext    context.Context
	getID         int64
	listContext   context.Context
	listFilter    ListFilter
	updateContext context.Context
	updateID      int64
	updateParams  UpdateParams
	updateVersion int64
	statusContext context.Context
	statusID      int64
	status        Status
	statusVersion int64
	deleteContext context.Context
	deleteID      int64
	deleteVersion int64
}

func (f *fakeRepository) Create(ctx context.Context, params CreateParams) (Task, error) {
	f.createCalls++
	f.createContext = ctx
	f.createParams = params
	return f.createResult, f.createErr
}

func (f *fakeRepository) Get(ctx context.Context, id int64) (Task, error) {
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

func (f *fakeRepository) Update(ctx context.Context, id int64, params UpdateParams, expectedVersion int64) (Task, error) {
	f.updateCalls++
	f.updateContext = ctx
	f.updateID = id
	f.updateParams = params
	f.updateVersion = expectedVersion
	return f.updateResult, f.updateErr
}

func (f *fakeRepository) UpdateStatus(ctx context.Context, id int64, status Status, expectedVersion int64) (Task, error) {
	f.updateStatusCalls++
	f.statusContext = ctx
	f.statusID = id
	f.status = status
	f.statusVersion = expectedVersion
	return f.updateStatusResult, f.updateStatusErr
}

func (f *fakeRepository) Delete(ctx context.Context, id int64, expectedVersion int64) error {
	f.deleteCalls++
	f.deleteContext = ctx
	f.deleteID = id
	f.deleteVersion = expectedVersion
	return f.deleteErr
}

type fakeUserReader struct {
	result user.User
	err    error
	calls  int
	ctx    context.Context
	id     int64
}

func (f *fakeUserReader) Get(ctx context.Context, id int64) (user.User, error) {
	f.calls++
	f.ctx = ctx
	f.id = id
	return f.result, f.err
}

func TestTaskJSONUsesCamelCaseAndPreservesNullableFields(t *testing.T) {
	t.Parallel()

	task := Task{
		ID:          1,
		OwnerID:     2,
		Title:       "编写文档",
		Description: nil,
		Status:      StatusTodo,
		DueAt:       nil,
		Version:     3,
		CreatedAt:   fixedNow,
		UpdatedAt:   fixedNow,
	}
	body, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"id", "ownerId", "title", "description", "status", "dueAt", "version", "createdAt", "updatedAt"} {
		if _, ok := got[key]; !ok {
			t.Fatalf("Task JSON is missing %q: %s", key, body)
		}
	}
	if got["description"] != nil || got["dueAt"] != nil {
		t.Fatalf("nullable fields = description:%v dueAt:%v, want null", got["description"], got["dueAt"])
	}
	if _, ok := got["owner_id"]; ok {
		t.Fatalf("Task JSON contains snake_case key: %s", body)
	}
}

func TestCreateNormalizesInputChecksActiveOwnerAndDefaultsTodo(t *testing.T) {
	t.Parallel()

	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("request"), "create-task")
	description := "  先写失败测试  "
	dueAt := fixedNow.Add(time.Hour)
	repo := &fakeRepository{createResult: Task{ID: 9}}
	users := &fakeUserReader{result: user.User{ID: 7, Status: user.StatusActive}}
	service := newTestService(repo, users)

	got, err := service.Create(ctx, CreateInput{
		OwnerID:     7,
		Title:       "  完成 Task Service  ",
		Description: &description,
		DueAt:       &dueAt,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 9 {
		t.Fatalf("Create() ID = %d, want 9", got.ID)
	}
	if users.calls != 1 || users.ctx != ctx || users.id != 7 {
		t.Fatalf("UserReader.Get calls/context/id = %d/%v/%d, want 1/original/7", users.calls, users.ctx == ctx, users.id)
	}
	if repo.createCalls != 1 || repo.createContext != ctx {
		t.Fatalf("Repository.Create calls/context = %d/%v, want 1/original", repo.createCalls, repo.createContext == ctx)
	}
	params := repo.createParams
	if params.OwnerID != 7 || params.Title != "完成 Task Service" || params.Status != StatusTodo {
		t.Fatalf("Create() params = %+v", params)
	}
	if params.Description == nil || *params.Description != "先写失败测试" {
		t.Fatalf("Create() description = %v, want trimmed non-nil value", params.Description)
	}
	if params.DueAt == nil || !params.DueAt.Equal(dueAt) {
		t.Fatalf("Create() dueAt = %v, want %v", params.DueAt, dueAt)
	}
	if params.Description == &description || params.DueAt == &dueAt {
		t.Fatal("Create() forwarded caller-owned pointers instead of defensive copies")
	}
}

func TestCreatePreservesDescriptionNullAndEmptySemantics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description *string
		wantNil     bool
		want        string
	}{
		{name: "nil means database null", description: nil, wantNil: true},
		{name: "whitespace remains explicit empty string", description: stringPointer(" \n\t "), want: ""},
		{name: "value is trimmed", description: stringPointer("  详细说明  "), want: "详细说明"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{}
			service := newTestService(repo, activeUserReader())
			_, err := service.Create(context.Background(), CreateInput{OwnerID: 1, Title: "标题", Description: tt.description})
			if err != nil {
				t.Fatal(err)
			}
			got := repo.createParams.Description
			if tt.wantNil {
				if got != nil {
					t.Fatalf("Description = %q, want nil", *got)
				}
				return
			}
			if got == nil || *got != tt.want {
				t.Fatalf("Description = %v, want non-nil %q", got, tt.want)
			}
		})
	}
}

func TestCreateValidatesOwnerAndMapsOwnerState(t *testing.T) {
	t.Parallel()

	readFailure := fmt.Errorf("read owner: %w", errors.New("database unavailable"))
	tests := []struct {
		name        string
		ownerID     int64
		owner       user.User
		userErr     error
		want        error
		wantField   string
		wantGetCall bool
	}{
		{name: "zero owner", ownerID: 0, wantField: "ownerId"},
		{name: "negative owner", ownerID: -1, wantField: "ownerId"},
		{name: "owner not found", ownerID: 8, userErr: fmt.Errorf("select user: %w", user.ErrNotFound), want: ErrOwnerNotFound, wantGetCall: true},
		{name: "owner disabled", ownerID: 8, owner: user.User{ID: 8, Status: user.StatusDisabled}, want: ErrOwnerDisabled, wantGetCall: true},
		{name: "reader failure", ownerID: 8, userErr: readFailure, want: readFailure, wantGetCall: true},
		{name: "context canceled", ownerID: 8, userErr: context.Canceled, want: context.Canceled, wantGetCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{}
			users := &fakeUserReader{result: tt.owner, err: tt.userErr}
			_, err := newTestService(repo, users).Create(context.Background(), CreateInput{OwnerID: tt.ownerID, Title: "标题"})
			if tt.wantField != "" {
				assertValidationField(t, err, tt.wantField)
			} else if !errors.Is(err, tt.want) {
				t.Fatalf("Create() error = %v, want errors.Is(%v)", err, tt.want)
			}
			wantCalls := 0
			if tt.wantGetCall {
				wantCalls = 1
			}
			if users.calls != wantCalls {
				t.Fatalf("UserReader.Get calls = %d, want %d", users.calls, wantCalls)
			}
			if repo.createCalls != 0 {
				t.Fatalf("Repository.Create calls = %d, want 0", repo.createCalls)
			}
		})
	}
}

func TestCreateValidatesTitleByTrimmedUnicodeRuneCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		title     string
		wantValid bool
		wantTitle string
	}{
		{name: "blank", title: " \t\n ", wantValid: false},
		{name: "one rune", title: "界", wantValid: true, wantTitle: "界"},
		{name: "128 unicode runes", title: "  " + strings.Repeat("界", 128) + "  ", wantValid: true, wantTitle: strings.Repeat("界", 128)},
		{name: "129 unicode runes", title: strings.Repeat("界", 129), wantValid: false},
		{name: "invalid utf8", title: string([]byte{0xff}), wantValid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{}
			users := activeUserReader()
			_, err := newTestService(repo, users).Create(context.Background(), CreateInput{OwnerID: 1, Title: tt.title})
			if tt.wantValid {
				if err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				if repo.createParams.Title != tt.wantTitle {
					t.Fatalf("Create() title = %q, want %q", repo.createParams.Title, tt.wantTitle)
				}
				return
			}
			assertValidationField(t, err, "title")
			if users.calls != 0 || repo.createCalls != 0 {
				t.Fatalf("invalid title called users/repo = %d/%d", users.calls, repo.createCalls)
			}
		})
	}
}

func TestCreateValidatesTitleCharactersWithoutRejectingUsefulFormatRunes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		title     string
		wantValid bool
	}{
		{name: "only zero-width format runes", title: "\u200b\u200d", wantValid: false},
		{name: "only combining marks", title: "\u0301\u0308", wantValid: false},
		{name: "only variation selectors", title: "\ufe0e\ufe0f", wantValid: false},
		{name: "embedded newline", title: "任务\n标题", wantValid: false},
		{name: "embedded NUL", title: "任务\x00标题", wantValid: false},
		{name: "Chinese title", title: "编写任务文档", wantValid: true},
		{name: "hyphenated title", title: "release-checklist", wantValid: true},
		{name: "combining character", title: "Cafe\u0301", wantValid: true},
		{name: "emoji joined with ZWJ", title: "👩\u200d💻", wantValid: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{}
			users := activeUserReader()
			_, err := newTestService(repo, users).Create(context.Background(), CreateInput{OwnerID: 1, Title: tt.title})
			if tt.wantValid {
				if err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				if repo.createParams.Title != tt.title {
					t.Fatalf("Repository.Create title = %q, want %q", repo.createParams.Title, tt.title)
				}
				return
			}
			assertValidationField(t, err, "title")
			if users.calls != 0 || repo.createCalls != 0 {
				t.Fatalf("invalid title called users/repo = %d/%d", users.calls, repo.createCalls)
			}
		})
	}
}

func TestCreateValidatesDueAtAgainstInjectedClock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		dueAt     time.Time
		wantValid bool
	}{
		{name: "past", dueAt: fixedNow.Add(-time.Nanosecond), wantValid: false},
		{name: "equal", dueAt: fixedNow, wantValid: true},
		{name: "future", dueAt: fixedNow.Add(time.Nanosecond), wantValid: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{}
			_, err := newTestService(repo, activeUserReader()).Create(context.Background(), CreateInput{OwnerID: 1, Title: "标题", DueAt: &tt.dueAt})
			if tt.wantValid {
				if err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				return
			}
			assertValidationField(t, err, "dueAt")
			if repo.createCalls != 0 {
				t.Fatalf("Repository.Create calls = %d, want 0", repo.createCalls)
			}
		})
	}
}

func TestCreatePropagatesRepositoryErrorIdentity(t *testing.T) {
	t.Parallel()

	for _, wantErr := range []error{context.Canceled, fmt.Errorf("insert task: %w", ErrVersionConflict)} {
		repo := &fakeRepository{createErr: wantErr}
		_, err := newTestService(repo, activeUserReader()).Create(context.Background(), CreateInput{OwnerID: 1, Title: "标题"})
		if err != wantErr {
			t.Fatalf("Create() error = %v, want original %v", err, wantErr)
		}
	}
}

func TestGetForwardsContextIDAndRepositoryError(t *testing.T) {
	t.Parallel()

	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("request"), "get-task")
	wantErr := fmt.Errorf("select task: %w", ErrNotFound)
	repo := &fakeRepository{getErr: wantErr}
	_, err := newTestService(repo, activeUserReader()).Get(ctx, 42)
	if err != wantErr || !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get() error = %v, want original ErrNotFound", err)
	}
	if repo.getContext != ctx || repo.getID != 42 {
		t.Fatalf("Get() context/id = %v/%d, want original/42", repo.getContext == ctx, repo.getID)
	}
}

func TestUpdateUsesOnlyEditableFieldsAndExpectedVersion(t *testing.T) {
	t.Parallel()

	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("request"), "update-task")
	description := "  修改后的说明  "
	dueAt := fixedNow.Add(2 * time.Hour)
	repo := &fakeRepository{updateResult: Task{ID: 11, Status: StatusDoing, OwnerID: 5}}
	service := newTestService(repo, activeUserReader())

	got, err := service.Update(ctx, 11, UpdateInput{
		Title:           "  新标题  ",
		Description:     &description,
		DueAt:           &dueAt,
		ExpectedVersion: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != StatusDoing || got.OwnerID != 5 {
		t.Fatalf("Update() result = %+v", got)
	}
	if repo.updateCalls != 1 || repo.updateContext != ctx || repo.updateID != 11 || repo.updateVersion != 0 {
		t.Fatalf("Repository.Update calls/context/id/version = %d/%v/%d/%d", repo.updateCalls, repo.updateContext == ctx, repo.updateID, repo.updateVersion)
	}
	params := repo.updateParams
	if params.Title != "新标题" || params.Description == nil || *params.Description != "修改后的说明" || params.DueAt == nil || !params.DueAt.Equal(dueAt) {
		t.Fatalf("Repository.Update params = %+v", params)
	}

	paramsType := reflect.TypeOf(UpdateParams{})
	for _, forbidden := range []string{"OwnerID", "Status"} {
		if _, ok := paramsType.FieldByName(forbidden); ok {
			t.Fatalf("UpdateParams exposes forbidden PUT field %s", forbidden)
		}
	}
}

func TestUpdateValidatesFieldsBeforeCallingRepository(t *testing.T) {
	t.Parallel()

	past := fixedNow.Add(-time.Second)
	tests := []struct {
		name  string
		input UpdateInput
		field string
	}{
		{name: "blank title", input: UpdateInput{Title: "   "}, field: "title"},
		{name: "title too long", input: UpdateInput{Title: strings.Repeat("界", 129)}, field: "title"},
		{name: "past due", input: UpdateInput{Title: "标题", DueAt: &past}, field: "dueAt"},
		{name: "negative expected version", input: UpdateInput{Title: "标题", ExpectedVersion: -1}, field: "expectedVersion"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{}
			_, err := newTestService(repo, activeUserReader()).Update(context.Background(), 1, tt.input)
			assertValidationField(t, err, tt.field)
			if repo.updateCalls != 0 {
				t.Fatalf("Repository.Update calls = %d, want 0", repo.updateCalls)
			}
		})
	}
}

func TestUpdateAcceptsDueAtEqualNowAndPreservesRepositoryErrors(t *testing.T) {
	t.Parallel()

	wantErr := fmt.Errorf("update task: %w", ErrVersionConflict)
	repo := &fakeRepository{updateErr: wantErr}
	_, err := newTestService(repo, activeUserReader()).Update(context.Background(), 1, UpdateInput{Title: "标题", DueAt: &fixedNow, ExpectedVersion: 2})
	if err != wantErr || !errors.Is(err, ErrVersionConflict) {
		t.Fatalf("Update() error = %v, want original version conflict", err)
	}
	if repo.updateVersion != 2 {
		t.Fatalf("Repository.Update expectedVersion = %d, want 2", repo.updateVersion)
	}
}

func TestChangeStatusAllowsOnlyDeclaredTransitions(t *testing.T) {
	t.Parallel()

	allowed := map[Status][]Status{
		StatusTodo:      {StatusDoing, StatusCancelled},
		StatusDoing:     {StatusTodo, StatusDone, StatusCancelled},
		StatusDone:      {},
		StatusCancelled: {},
	}
	all := []Status{StatusTodo, StatusDoing, StatusDone, StatusCancelled}

	for from, allowedTargets := range allowed {
		for _, to := range all {
			from, to := from, to
			wantAllowed := containsStatus(allowedTargets, to)
			t.Run(string(from)+"_to_"+string(to), func(t *testing.T) {
				t.Parallel()
				repo := &fakeRepository{
					getResult:          Task{ID: 3, Status: from, Version: 7},
					updateStatusResult: Task{ID: 3, Status: to, Version: 8},
				}
				got, err := newTestService(repo, activeUserReader()).ChangeStatus(context.Background(), 3, ChangeStatusInput{Status: to, ExpectedVersion: 7})
				if wantAllowed {
					if err != nil {
						t.Fatalf("ChangeStatus() error = %v", err)
					}
					if got.Status != to || repo.updateStatusCalls != 1 || repo.statusID != 3 || repo.status != to || repo.statusVersion != 7 {
						t.Fatalf("ChangeStatus() result/call = %+v/%d/%d/%s/%d", got, repo.updateStatusCalls, repo.statusID, repo.status, repo.statusVersion)
					}
					return
				}
				if !errors.Is(err, ErrInvalidTransition) {
					t.Fatalf("ChangeStatus() error = %v, want ErrInvalidTransition", err)
				}
				if repo.updateStatusCalls != 0 {
					t.Fatalf("Repository.UpdateStatus calls = %d, want 0", repo.updateStatusCalls)
				}
			})
		}
	}
}

func TestChangeStatusChecksVersionBeforeTransition(t *testing.T) {
	t.Parallel()

	t.Run("stale version takes precedence over invalid transition", func(t *testing.T) {
		t.Parallel()
		repo := &fakeRepository{getResult: Task{ID: 1, Status: StatusDone, Version: 4}}
		_, err := newTestService(repo, activeUserReader()).ChangeStatus(context.Background(), 1, ChangeStatusInput{
			Status:          StatusDoing,
			ExpectedVersion: 3,
		})
		if err != ErrVersionConflict {
			t.Fatalf("ChangeStatus() error = %v, want exact ErrVersionConflict", err)
		}
		if repo.updateStatusCalls != 0 {
			t.Fatalf("Repository.UpdateStatus calls = %d, want 0", repo.updateStatusCalls)
		}
	})

	t.Run("matching version still rejects invalid transition", func(t *testing.T) {
		t.Parallel()
		repo := &fakeRepository{getResult: Task{ID: 1, Status: StatusDone, Version: 3}}
		_, err := newTestService(repo, activeUserReader()).ChangeStatus(context.Background(), 1, ChangeStatusInput{
			Status:          StatusDoing,
			ExpectedVersion: 3,
		})
		if err != ErrInvalidTransition {
			t.Fatalf("ChangeStatus() error = %v, want exact ErrInvalidTransition", err)
		}
		if repo.updateStatusCalls != 0 {
			t.Fatalf("Repository.UpdateStatus calls = %d, want 0", repo.updateStatusCalls)
		}
	})

	t.Run("concurrent update after read preserves repository conflict", func(t *testing.T) {
		t.Parallel()
		original := fmt.Errorf("update status after concurrent write: %w", ErrVersionConflict)
		repo := &fakeRepository{
			getResult:       Task{ID: 1, Status: StatusTodo, Version: 3},
			updateStatusErr: original,
		}
		_, err := newTestService(repo, activeUserReader()).ChangeStatus(context.Background(), 1, ChangeStatusInput{
			Status:          StatusDoing,
			ExpectedVersion: 3,
		})
		if err != original {
			t.Fatalf("ChangeStatus() error = %v, want original %v", err, original)
		}
		if !errors.Is(err, ErrVersionConflict) {
			t.Fatalf("ChangeStatus() error = %v, want errors.Is(ErrVersionConflict)", err)
		}
		if repo.updateStatusCalls != 1 {
			t.Fatalf("Repository.UpdateStatus calls = %d, want 1", repo.updateStatusCalls)
		}
	})
}

func TestChangeStatusValidatesUnknownStatusAndExpectedVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input ChangeStatusInput
		field string
	}{
		{name: "unknown status", input: ChangeStatusInput{Status: Status("BLOCKED")}, field: "status"},
		{name: "negative expected version", input: ChangeStatusInput{Status: StatusDoing, ExpectedVersion: -1}, field: "expectedVersion"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{getResult: Task{ID: 1, Status: StatusTodo}}
			_, err := newTestService(repo, activeUserReader()).ChangeStatus(context.Background(), 1, tt.input)
			assertValidationField(t, err, tt.field)
			if repo.getCalls != 0 || repo.updateStatusCalls != 0 {
				t.Fatalf("invalid input called Get/UpdateStatus = %d/%d", repo.getCalls, repo.updateStatusCalls)
			}
		})
	}
}

func TestChangeStatusPreservesGetAndUpdateErrors(t *testing.T) {
	t.Parallel()

	getSentinel := errors.New("task repository get failed")
	wrappedGetError := fmt.Errorf("select task: %w", getSentinel)
	updateSentinel := errors.New("task repository update status failed")
	wrappedUpdateError := fmt.Errorf("update task status: %w", updateSentinel)

	tests := []struct {
		name      string
		getErr    error
		updateErr error
		wantExact error
		wantIs    error
	}{
		{name: "wrapped get error", getErr: wrappedGetError, wantExact: wrappedGetError, wantIs: getSentinel},
		{name: "get canceled", getErr: context.Canceled, wantExact: context.Canceled, wantIs: context.Canceled},
		{name: "wrapped update error", updateErr: wrappedUpdateError, wantExact: wrappedUpdateError, wantIs: updateSentinel},
		{name: "update canceled", updateErr: context.Canceled, wantExact: context.Canceled, wantIs: context.Canceled},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{getResult: Task{ID: 1, Status: StatusTodo}, getErr: tt.getErr, updateStatusErr: tt.updateErr}
			_, err := newTestService(repo, activeUserReader()).ChangeStatus(context.Background(), 1, ChangeStatusInput{Status: StatusDoing, ExpectedVersion: 0})
			if err != tt.wantExact {
				t.Fatalf("ChangeStatus() error = %v, want original %v", err, tt.wantExact)
			}
			if !errors.Is(err, tt.wantIs) {
				t.Fatalf("ChangeStatus() error = %v, want errors.Is(%v)", err, tt.wantIs)
			}
			if tt.getErr != nil && repo.updateStatusCalls != 0 {
				t.Fatalf("Repository.UpdateStatus calls = %d after Get error", repo.updateStatusCalls)
			}
		})
	}
}

func TestListAppliesDefaultsAndReturnsPageMetadata(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{listResult: Page{Items: []Task{{ID: 1}}, Total: 41}}
	got, err := newTestService(repo, activeUserReader()).List(context.Background(), ListInput{})
	if err != nil {
		t.Fatal(err)
	}
	if repo.listFilter.OwnerID != nil || repo.listFilter.Status != nil || repo.listFilter.Limit != 20 || repo.listFilter.Offset != 0 {
		t.Fatalf("Repository.List filter = %+v, want no filters limit=20 offset=0", repo.listFilter)
	}
	if got.Page != 1 || got.PageSize != 20 || got.Total != 41 || len(got.Items) != 1 {
		t.Fatalf("List() page = %+v", got)
	}
}

func TestListForwardsOwnerStatusAndLimitOffset(t *testing.T) {
	t.Parallel()

	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("request"), "list-task")
	repo := &fakeRepository{listResult: Page{Total: 250}}
	got, err := newTestService(repo, activeUserReader()).List(ctx, ListInput{Page: 3, PageSize: 100, OwnerID: 8, Status: StatusDoing})
	if err != nil {
		t.Fatal(err)
	}
	filter := repo.listFilter
	if repo.listContext != ctx || filter.OwnerID == nil || *filter.OwnerID != 8 || filter.Status == nil || *filter.Status != StatusDoing || filter.Limit != 100 || filter.Offset != 200 {
		t.Fatalf("Repository.List context/filter = %v/%+v", repo.listContext == ctx, filter)
	}
	if got.Page != 3 || got.PageSize != 100 || got.Total != 250 {
		t.Fatalf("List() page = %+v", got)
	}
}

func TestListValidatesFiltersPaginationAndOffsetOverflow(t *testing.T) {
	t.Parallel()

	maxInt := int(^uint(0) >> 1)
	tests := []struct {
		name  string
		input ListInput
		field string
	}{
		{name: "negative page", input: ListInput{Page: -1}, field: "page"},
		{name: "negative page size", input: ListInput{PageSize: -1}, field: "pageSize"},
		{name: "page size over maximum", input: ListInput{PageSize: 101}, field: "pageSize"},
		{name: "negative owner", input: ListInput{OwnerID: -1}, field: "ownerId"},
		{name: "unknown status", input: ListInput{Status: Status("BLOCKED")}, field: "status"},
		{name: "offset overflow", input: ListInput{Page: maxInt/100 + 2, PageSize: 100}, field: "page"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeRepository{}
			_, err := newTestService(repo, activeUserReader()).List(context.Background(), tt.input)
			assertValidationField(t, err, tt.field)
			if repo.listCalls != 0 {
				t.Fatalf("Repository.List calls = %d, want 0", repo.listCalls)
			}
		})
	}
}

func TestListPreservesRepositoryErrors(t *testing.T) {
	t.Parallel()

	for _, wantErr := range []error{context.Canceled, fmt.Errorf("list tasks: %w", ErrNotFound)} {
		repo := &fakeRepository{listErr: wantErr}
		_, err := newTestService(repo, activeUserReader()).List(context.Background(), ListInput{})
		if err != wantErr {
			t.Fatalf("List() error = %v, want original %v", err, wantErr)
		}
	}
}

func TestDeleteValidatesAndForwardsExpectedVersion(t *testing.T) {
	t.Parallel()

	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("request"), "delete-task")
	repo := &fakeRepository{}
	service := newTestService(repo, activeUserReader())
	if err := service.Delete(ctx, 15, DeleteInput{ExpectedVersion: 0}); err != nil {
		t.Fatal(err)
	}
	if repo.deleteCalls != 1 || repo.deleteContext != ctx || repo.deleteID != 15 || repo.deleteVersion != 0 {
		t.Fatalf("Repository.Delete calls/context/id/version = %d/%v/%d/%d", repo.deleteCalls, repo.deleteContext == ctx, repo.deleteID, repo.deleteVersion)
	}

	repo = &fakeRepository{}
	err := newTestService(repo, activeUserReader()).Delete(context.Background(), 15, DeleteInput{ExpectedVersion: -1})
	assertValidationField(t, err, "expectedVersion")
	if repo.deleteCalls != 0 {
		t.Fatalf("Repository.Delete calls = %d for invalid version", repo.deleteCalls)
	}
}

func TestDeletePreservesVersionAndContextErrors(t *testing.T) {
	t.Parallel()

	for _, wantErr := range []error{fmt.Errorf("delete task: %w", ErrVersionConflict), context.Canceled} {
		repo := &fakeRepository{deleteErr: wantErr}
		err := newTestService(repo, activeUserReader()).Delete(context.Background(), 1, DeleteInput{ExpectedVersion: 4})
		if err != wantErr {
			t.Fatalf("Delete() error = %v, want original %v", err, wantErr)
		}
	}
}

func TestServiceRejectsNilAndTypedNilDependenciesWithoutPanic(t *testing.T) {
	t.Parallel()

	var typedNilRepo *fakeRepository
	var typedNilUsers *fakeUserReader
	tests := []struct {
		name    string
		service *Service
		call    func(*Service) error
		want    error
	}{
		{name: "nil service", service: nil, call: func(s *Service) error { _, err := s.Get(context.Background(), 1); return err }, want: ErrNilRepository},
		{name: "nil repository", service: newTestService(nil, activeUserReader()), call: func(s *Service) error { _, err := s.Get(context.Background(), 1); return err }, want: ErrNilRepository},
		{name: "typed nil repository", service: newTestService(typedNilRepo, activeUserReader()), call: func(s *Service) error { _, err := s.Get(context.Background(), 1); return err }, want: ErrNilRepository},
		{name: "nil user reader", service: newTestService(&fakeRepository{}, nil), call: func(s *Service) error {
			_, err := s.Create(context.Background(), CreateInput{OwnerID: 1, Title: "标题"})
			return err
		}, want: ErrNilUserReader},
		{name: "typed nil user reader", service: newTestService(&fakeRepository{}, typedNilUsers), call: func(s *Service) error {
			_, err := s.Create(context.Background(), CreateInput{OwnerID: 1, Title: "标题"})
			return err
		}, want: ErrNilUserReader},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Fatalf("service panicked: %v", recovered)
				}
			}()
			if err := tt.call(tt.service); !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want errors.Is(%v)", err, tt.want)
			}
		})
	}
}

func TestValidationErrorSupportsErrorsIsAs(t *testing.T) {
	t.Parallel()

	err := invalidField("title", "标题不能为空")
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("errors.Is(%v, ErrValidation) = false", err)
	}
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) || validationErr.Fields["title"] == "" {
		t.Fatalf("errors.As(%v) = %v/%+v", err, validationErr != nil, validationErr)
	}
	if (*ValidationError)(nil).Error() != ErrValidation.Error() || (*ValidationError)(nil).Unwrap() != ErrValidation {
		t.Fatal("nil ValidationError receiver is not safe")
	}
}

func newTestService(repo Repository, users UserReader) *Service {
	return NewService(repo, users, WithNow(func() time.Time { return fixedNow }))
}

func activeUserReader() *fakeUserReader {
	return &fakeUserReader{result: user.User{ID: 1, Status: user.StatusActive}}
}

func stringPointer(value string) *string {
	return &value
}

func containsStatus(statuses []Status, target Status) bool {
	for _, status := range statuses {
		if status == target {
			return true
		}
	}
	return false
}

func assertValidationField(t *testing.T, err error, field string) {
	t.Helper()
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("error = %v, want ErrValidation", err)
	}
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %v, want *ValidationError", err)
	}
	if validationErr.Fields[field] == "" {
		t.Fatalf("ValidationError fields = %+v, want %q", validationErr.Fields, field)
	}
}
