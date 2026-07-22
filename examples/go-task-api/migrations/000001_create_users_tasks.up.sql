create table users (
  id bigint generated always as identity,
  name varchar(64) not null,
  email varchar(254) not null,
  status varchar(16) not null default 'ACTIVE',
  version bigint not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  constraint pk_users primary key (id),
  constraint ck_users_status check (status in ('ACTIVE', 'DISABLED')),
  constraint ck_users_version check (version >= 0)
);

create unique index uk_users_email_lower on users (lower(email));
create index idx_users_created_id on users (created_at desc, id desc);
create index idx_users_status_created_id on users (status, created_at desc, id desc);

create table tasks (
  id bigint generated always as identity,
  owner_id bigint not null,
  title varchar(128) not null,
  description text,
  status varchar(16) not null default 'TODO',
  due_at timestamptz,
  version bigint not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  constraint pk_tasks primary key (id),
  constraint fk_tasks_owner foreign key (owner_id) references users (id) on delete restrict,
  constraint ck_tasks_status check (status in ('TODO', 'DOING', 'DONE', 'CANCELLED')),
  constraint ck_tasks_version check (version >= 0)
);

create index idx_tasks_owner_status on tasks (owner_id, status, id);
create index idx_tasks_status_due_at on tasks (status, due_at, id);
create index idx_tasks_created_id on tasks (created_at desc, id desc);
create index idx_tasks_owner_created_id on tasks (owner_id, created_at desc, id desc);
create index idx_tasks_status_created_id on tasks (status, created_at desc, id desc);
create index idx_tasks_owner_status_created_id on tasks (owner_id, status, created_at desc, id desc);

comment on table users is '系统用户信息表';
comment on column users.id is '用户唯一标识';
comment on column users.name is '用户显示名称';
comment on column users.email is '用户电子邮箱';
comment on column users.status is '用户状态：启用或禁用';
comment on column users.version is '用户记录乐观锁版本号';
comment on column users.created_at is '用户创建时间';
comment on column users.updated_at is '用户最后更新时间';
comment on constraint pk_users on users is '用户表主键约束';
comment on constraint ck_users_status on users is '用户状态取值约束';
comment on constraint ck_users_version on users is '用户版本号非负约束';
comment on index pk_users is '用户表主键索引';
comment on index uk_users_email_lower is '用户邮箱忽略大小写唯一索引';
comment on index idx_users_created_id is '按创建时间和用户标识稳定倒序分页的索引';
comment on index idx_users_status_created_id is '按状态筛选并按创建时间和用户标识稳定倒序分页的索引';

comment on table tasks is '任务信息表';
comment on column tasks.id is '任务唯一标识';
comment on column tasks.owner_id is '任务负责人用户标识';
comment on column tasks.title is '任务标题';
comment on column tasks.description is '任务详细描述';
comment on column tasks.status is '任务状态：待办、进行中、已完成或已取消';
comment on column tasks.due_at is '任务截止时间';
comment on column tasks.version is '任务记录乐观锁版本号';
comment on column tasks.created_at is '任务创建时间';
comment on column tasks.updated_at is '任务最后更新时间';
comment on constraint pk_tasks on tasks is '任务表主键约束';
comment on constraint fk_tasks_owner on tasks is '任务负责人外键约束';
comment on constraint ck_tasks_status on tasks is '任务状态取值约束';
comment on constraint ck_tasks_version on tasks is '任务版本号非负约束';
comment on index pk_tasks is '任务表主键索引';
comment on index idx_tasks_owner_status is '按负责人和状态查询任务的复合索引';
comment on index idx_tasks_status_due_at is '按状态和截止时间查询任务的复合索引';
comment on index idx_tasks_created_id is '按创建时间和任务标识稳定倒序分页的索引';
comment on index idx_tasks_owner_created_id is '按负责人筛选并按创建时间和任务标识稳定倒序分页的索引';
comment on index idx_tasks_status_created_id is '按状态筛选并按创建时间和任务标识稳定倒序分页的索引';
comment on index idx_tasks_owner_status_created_id is '按负责人和状态筛选并稳定倒序分页的索引';
