package task

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

var errNilDatabase = errors.New("task repository: database is nil")

// PostgresRepository 使用调用方持有的连接池，不负责关闭连接池。
type PostgresRepository struct {
	db *sql.DB
}

var _ Repository = (*PostgresRepository)(nil)

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, params CreateParams) (Task, error) {
	db, err := r.database()
	if err != nil {
		return Task{}, err
	}

	created, err := scanTask(db.QueryRowContext(ctx, `
		insert into tasks (owner_id, title, description, status, due_at)
		values ($1, $2, $3, $4, $5)
		returning id, owner_id, title, description, status, due_at, version, created_at, updated_at`,
		params.OwnerID,
		params.Title,
		nullableString(params.Description),
		params.Status,
		nullableTime(params.DueAt),
	))
	if err != nil {
		return Task{}, mapPostgresError("create task", err)
	}
	return created, nil
}

func (r *PostgresRepository) Get(ctx context.Context, id int64) (Task, error) {
	db, err := r.database()
	if err != nil {
		return Task{}, err
	}

	found, err := scanTask(db.QueryRowContext(ctx, `
		select id, owner_id, title, description, status, due_at, version, created_at, updated_at
		from tasks
		where id = $1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	if err != nil {
		return Task{}, mapPostgresError("get task", err)
	}
	return found, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter ListFilter) (Page, error) {
	db, err := r.database()
	if err != nil {
		return Page{}, err
	}

	ownerID := optionalOwnerID(filter.OwnerID)
	status := optionalStatus(filter.Status)
	rows, err := db.QueryContext(ctx, `
		with filtered as not materialized (
			select id, owner_id, title, description, status, due_at, version, created_at, updated_at
			from tasks
			where ($1::bigint is null or owner_id = $1)
			  and ($2::varchar is null or status = $2)
		),
		paged as (
			select id, owner_id, title, description, status, due_at, version, created_at, updated_at,
			       count(*) over () as total
			from filtered
			order by created_at desc, id desc
			limit $3 offset $4
		),
		fallback as (
			select null::bigint as id,
			       null::bigint as owner_id,
			       null::varchar as title,
			       null::text as description,
			       null::varchar as status,
			       null::timestamptz as due_at,
			       null::bigint as version,
			       null::timestamptz as created_at,
			       null::timestamptz as updated_at,
			       (select count(*) from filtered) as total
			where not exists (select 1 from paged)
		)
		select id, owner_id, title, description, status, due_at, version, created_at, updated_at, total
		from paged
		union all
		select id, owner_id, title, description, status, due_at, version, created_at, updated_at, total
		from fallback
		order by created_at desc nulls last, id desc nulls last`, ownerID, status, filter.Limit, filter.Offset)
	if err != nil {
		return Page{}, mapPostgresError("list tasks", err)
	}
	defer rows.Close()

	page := Page{Items: make([]Task, 0)}
	for rows.Next() {
		item, total, present, err := scanOptionalTaskWithTotal(rows)
		if err != nil {
			return Page{}, mapPostgresError("scan tasks", err)
		}
		page.Total = total
		if present {
			page.Items = append(page.Items, item)
		}
	}
	if err := rows.Err(); err != nil {
		return Page{}, mapPostgresError("iterate tasks", err)
	}

	return page, nil
}

func (r *PostgresRepository) Update(
	ctx context.Context,
	id int64,
	params UpdateParams,
	expectedVersion int64,
) (Task, error) {
	db, err := r.database()
	if err != nil {
		return Task{}, err
	}

	updated, targetExists, wasUpdated, err := scanTaskWriteResult(db.QueryRowContext(ctx, `
		with target as materialized (
			select id, true as present
			from tasks
			where id = $4
		),
		updated as (
			update tasks as current
			set title = $1,
			    description = $2,
			    due_at = $3,
			    version = version + 1,
			    updated_at = now()
			from target
			where current.id = target.id and current.version = $5
			returning current.id, current.owner_id, current.title, current.description,
			          current.status, current.due_at, current.version,
			          current.created_at, current.updated_at
		)
		select coalesce((select present from target limit 1), false),
		       updated.id, updated.owner_id, updated.title, updated.description,
		       updated.status, updated.due_at, updated.version,
		       updated.created_at, updated.updated_at
		from (values (1)) as anchor(value)
		left join updated on true`,
		params.Title,
		nullableString(params.Description),
		nullableTime(params.DueAt),
		id,
		expectedVersion,
	))
	if err != nil {
		return Task{}, mapPostgresError("update task", err)
	}
	if !wasUpdated {
		if targetExists {
			return Task{}, ErrVersionConflict
		}
		return Task{}, ErrNotFound
	}
	return updated, nil
}

func (r *PostgresRepository) UpdateStatus(
	ctx context.Context,
	id int64,
	status Status,
	expectedVersion int64,
) (Task, error) {
	db, err := r.database()
	if err != nil {
		return Task{}, err
	}

	updated, targetExists, wasUpdated, err := scanTaskWriteResult(db.QueryRowContext(ctx, `
		with target as materialized (
			select id, true as present
			from tasks
			where id = $2
		),
		updated as (
			update tasks as current
			set status = $1, version = version + 1, updated_at = now()
			from target
			where current.id = target.id and current.version = $3
			returning current.id, current.owner_id, current.title, current.description,
			          current.status, current.due_at, current.version,
			          current.created_at, current.updated_at
		)
		select coalesce((select present from target limit 1), false),
		       updated.id, updated.owner_id, updated.title, updated.description,
		       updated.status, updated.due_at, updated.version,
		       updated.created_at, updated.updated_at
		from (values (1)) as anchor(value)
		left join updated on true`,
		status,
		id,
		expectedVersion,
	))
	if err != nil {
		return Task{}, mapPostgresError("update task status", err)
	}
	if !wasUpdated {
		if targetExists {
			return Task{}, ErrVersionConflict
		}
		return Task{}, ErrNotFound
	}
	return updated, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64, expectedVersion int64) error {
	db, err := r.database()
	if err != nil {
		return err
	}

	var targetExists, wasDeleted bool
	err = db.QueryRowContext(ctx, `
		with target as materialized (
			select id, true as present
			from tasks
			where id = $1
		),
		deleted as (
			delete from tasks as current
			using target
			where current.id = target.id and current.version = $2
			returning true as deleted
		)
		select coalesce((select present from target limit 1), false),
		       coalesce((select deleted from deleted limit 1), false)`,
		id,
		expectedVersion,
	).Scan(&targetExists, &wasDeleted)
	if err != nil {
		return mapPostgresError("delete task", err)
	}
	if wasDeleted {
		return nil
	}
	if targetExists {
		return ErrVersionConflict
	}
	return ErrNotFound
}

func (r *PostgresRepository) database() (*sql.DB, error) {
	if r == nil || r.db == nil {
		return nil, errNilDatabase
	}
	return r.db, nil
}

type taskScanner interface {
	Scan(...any) error
}

func scanTask(scanner taskScanner) (Task, error) {
	var found Task
	var description sql.NullString
	var dueAt sql.NullTime
	err := scanner.Scan(
		&found.ID,
		&found.OwnerID,
		&found.Title,
		&description,
		&found.Status,
		&dueAt,
		&found.Version,
		&found.CreatedAt,
		&found.UpdatedAt,
	)
	if err != nil {
		return Task{}, err
	}
	setNullableFields(&found, description, dueAt)
	return found, nil
}

func scanOptionalTaskWithTotal(scanner taskScanner) (Task, int64, bool, error) {
	var id, ownerID, version sql.NullInt64
	var title, status sql.NullString
	var description sql.NullString
	var dueAt sql.NullTime
	var createdAt, updatedAt sql.NullTime
	var total int64
	err := scanner.Scan(
		&id,
		&ownerID,
		&title,
		&description,
		&status,
		&dueAt,
		&version,
		&createdAt,
		&updatedAt,
		&total,
	)
	if err != nil {
		return Task{}, 0, false, err
	}
	if !id.Valid {
		return Task{}, total, false, nil
	}
	if !ownerID.Valid || !title.Valid || !status.Valid || !version.Valid || !createdAt.Valid || !updatedAt.Valid {
		return Task{}, 0, false, errors.New("task repository: incomplete list row")
	}
	found := Task{
		ID:        id.Int64,
		OwnerID:   ownerID.Int64,
		Title:     title.String,
		Status:    Status(status.String),
		Version:   version.Int64,
		CreatedAt: createdAt.Time,
		UpdatedAt: updatedAt.Time,
	}
	setNullableFields(&found, description, dueAt)
	return found, total, true, nil
}

func scanTaskWriteResult(scanner taskScanner) (Task, bool, bool, error) {
	var targetExists bool
	var id, ownerID, version sql.NullInt64
	var title, status sql.NullString
	var description sql.NullString
	var dueAt, createdAt, updatedAt sql.NullTime
	if err := scanner.Scan(
		&targetExists,
		&id,
		&ownerID,
		&title,
		&description,
		&status,
		&dueAt,
		&version,
		&createdAt,
		&updatedAt,
	); err != nil {
		return Task{}, false, false, err
	}
	if !id.Valid {
		return Task{}, targetExists, false, nil
	}
	if !ownerID.Valid || !title.Valid || !status.Valid || !version.Valid || !createdAt.Valid || !updatedAt.Valid {
		return Task{}, false, false, errors.New("task repository: incomplete update row")
	}
	found := Task{
		ID:        id.Int64,
		OwnerID:   ownerID.Int64,
		Title:     title.String,
		Status:    Status(status.String),
		Version:   version.Int64,
		CreatedAt: createdAt.Time,
		UpdatedAt: updatedAt.Time,
	}
	setNullableFields(&found, description, dueAt)
	return found, targetExists, true, nil
}

func setNullableFields(found *Task, description sql.NullString, dueAt sql.NullTime) {
	if description.Valid {
		value := description.String
		found.Description = &value
	}
	if dueAt.Valid {
		value := dueAt.Time
		found.DueAt = &value
	}
}

func nullableString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func nullableTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return *value
}

func optionalOwnerID(ownerID *int64) any {
	if ownerID == nil {
		return nil
	}
	return *ownerID
}

func optionalStatus(status *Status) any {
	if status == nil {
		return nil
	}
	return string(*status)
}

func mapPostgresError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return ErrOwnerNotFound
	}
	return &taskRepositoryError{operation: operation, cause: err}
}

type taskRepositoryError struct {
	operation string
	cause     error
}

func (e *taskRepositoryError) Error() string {
	return "task repository: " + e.operation + " failed"
}

func (e *taskRepositoryError) Unwrap() error {
	return e.cause
}
