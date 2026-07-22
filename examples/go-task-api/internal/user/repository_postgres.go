package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var errNilDatabase = errors.New("user repository: database is nil")

// PostgresRepository 使用调用方持有的连接池，不负责关闭连接池。
type PostgresRepository struct {
	db *sql.DB
}

var _ Repository = (*PostgresRepository)(nil)

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, params CreateParams) (User, error) {
	db, err := r.database()
	if err != nil {
		return User{}, err
	}

	created, err := scanUser(db.QueryRowContext(ctx, `
		insert into users (name, email, status)
		values ($1, $2, $3)
		returning id, name, email, status, version, created_at, updated_at`,
		params.Name,
		params.Email,
		params.Status,
	))
	if err != nil {
		return User{}, mapPostgresError("create user", err)
	}
	return created, nil
}

func (r *PostgresRepository) Get(ctx context.Context, id int64) (User, error) {
	db, err := r.database()
	if err != nil {
		return User{}, err
	}

	found, err := scanUser(db.QueryRowContext(ctx, `
		select id, name, email, status, version, created_at, updated_at
		from users
		where id = $1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, mapPostgresError("get user", err)
	}
	return found, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter ListFilter) (Page, error) {
	db, err := r.database()
	if err != nil {
		return Page{}, err
	}

	status := optionalStatus(filter.Status)
	rows, err := db.QueryContext(ctx, `
		with filtered as not materialized (
			select id, name, email, status, version, created_at, updated_at
			from users
			where ($1::varchar is null or status = $1)
		),
		paged as (
			select id, name, email, status, version, created_at, updated_at,
			       count(*) over () as total
			from filtered
			order by created_at desc, id desc
			limit $2 offset $3
		),
		fallback as (
			select null::bigint as id,
			       null::varchar as name,
			       null::varchar as email,
			       null::varchar as status,
			       null::bigint as version,
			       null::timestamptz as created_at,
			       null::timestamptz as updated_at,
			       (select count(*) from filtered) as total
			where not exists (select 1 from paged)
		)
		select id, name, email, status, version, created_at, updated_at, total
		from paged
		union all
		select id, name, email, status, version, created_at, updated_at, total
		from fallback
		order by created_at desc nulls last, id desc nulls last`, status, filter.Limit, filter.Offset)
	if err != nil {
		return Page{}, mapPostgresError("list users", err)
	}
	defer rows.Close()

	page := Page{Items: make([]User, 0)}
	for rows.Next() {
		item, total, present, err := scanOptionalUserWithTotal(rows)
		if err != nil {
			return Page{}, mapPostgresError("scan users", err)
		}
		page.Total = total
		if present {
			page.Items = append(page.Items, item)
		}
	}
	if err := rows.Err(); err != nil {
		return Page{}, mapPostgresError("iterate users", err)
	}

	return page, nil
}

func (r *PostgresRepository) UpdateStatus(
	ctx context.Context,
	id int64,
	status Status,
	expectedVersion int64,
) (User, error) {
	db, err := r.database()
	if err != nil {
		return User{}, err
	}

	updated, targetExists, wasUpdated, err := scanUserWriteResult(db.QueryRowContext(ctx, `
		with target as materialized (
			select id, true as present
			from users
			where id = $2
		),
		updated as (
			update users as current
			set status = $1, version = version + 1, updated_at = now()
			from target
			where current.id = target.id and current.version = $3
			returning current.id, current.name, current.email, current.status,
			          current.version, current.created_at, current.updated_at
		)
		select coalesce((select present from target limit 1), false),
		       updated.id, updated.name, updated.email, updated.status,
		       updated.version, updated.created_at, updated.updated_at
		from (values (1)) as anchor(value)
		left join updated on true`,
		status,
		id,
		expectedVersion,
	))
	if err != nil {
		return User{}, mapPostgresError("update user status", err)
	}
	if !wasUpdated {
		if targetExists {
			return User{}, ErrVersionConflict
		}
		return User{}, ErrNotFound
	}
	return updated, nil
}

func (r *PostgresRepository) database() (*sql.DB, error) {
	if r == nil || r.db == nil {
		return nil, errNilDatabase
	}
	return r.db, nil
}

type userScanner interface {
	Scan(...any) error
}

func scanUser(scanner userScanner) (User, error) {
	var found User
	err := scanner.Scan(
		&found.ID,
		&found.Name,
		&found.Email,
		&found.Status,
		&found.Version,
		&found.CreatedAt,
		&found.UpdatedAt,
	)
	return found, err
}

func scanOptionalUserWithTotal(scanner userScanner) (User, int64, bool, error) {
	var id, version sql.NullInt64
	var name, email, status sql.NullString
	var createdAt, updatedAt sql.NullTime
	var total int64
	if err := scanner.Scan(
		&id,
		&name,
		&email,
		&status,
		&version,
		&createdAt,
		&updatedAt,
		&total,
	); err != nil {
		return User{}, 0, false, err
	}
	if !id.Valid {
		return User{}, total, false, nil
	}
	if !name.Valid || !email.Valid || !status.Valid || !version.Valid || !createdAt.Valid || !updatedAt.Valid {
		return User{}, 0, false, errors.New("user repository: incomplete list row")
	}
	return User{
		ID:        id.Int64,
		Name:      name.String,
		Email:     email.String,
		Status:    Status(status.String),
		Version:   version.Int64,
		CreatedAt: createdAt.Time,
		UpdatedAt: updatedAt.Time,
	}, total, true, nil
}

func scanUserWriteResult(scanner userScanner) (User, bool, bool, error) {
	var targetExists bool
	var id, version sql.NullInt64
	var name, email, status sql.NullString
	var createdAt, updatedAt sql.NullTime
	if err := scanner.Scan(
		&targetExists,
		&id,
		&name,
		&email,
		&status,
		&version,
		&createdAt,
		&updatedAt,
	); err != nil {
		return User{}, false, false, err
	}
	if !id.Valid {
		return User{}, targetExists, false, nil
	}
	if !name.Valid || !email.Valid || !status.Valid || !version.Valid || !createdAt.Valid || !updatedAt.Valid {
		return User{}, false, false, errors.New("user repository: incomplete update row")
	}
	return User{
		ID:        id.Int64,
		Name:      name.String,
		Email:     email.String,
		Status:    Status(status.String),
		Version:   version.Int64,
		CreatedAt: createdAt.Time,
		UpdatedAt: updatedAt.Time,
	}, targetExists, true, nil
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
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrEmailConflict
	}
	return &userRepositoryError{operation: operation, cause: err}
}

type userRepositoryError struct {
	operation string
	cause     error
}

func (e *userRepositoryError) Error() string {
	return "user repository: " + e.operation + " failed"
}

func (e *userRepositoryError) Unwrap() error {
	return e.cause
}
