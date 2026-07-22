CREATE TABLE roles (
  id uuid PRIMARY KEY,
  code varchar(64) NOT NULL,
  name varchar(80) NOT NULL,
  status varchar(20) NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  CONSTRAINT uk_roles_code UNIQUE (code),
  CONSTRAINT ck_roles_status CHECK (status IN ('ACTIVE', 'DISABLED'))
);

CREATE TABLE users (
  id uuid PRIMARY KEY,
  email varchar(254) NOT NULL,
  display_name varchar(80) NOT NULL,
  status varchar(20) NOT NULL,
  version bigint NOT NULL DEFAULT 0,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  CONSTRAINT ck_users_status CHECK (status IN ('ACTIVE', 'DISABLED'))
);

CREATE UNIQUE INDEX uk_users_email_lower ON users (lower(email));
CREATE INDEX idx_users_status_created ON users (status, created_at DESC, id);

CREATE TABLE user_roles (
  user_id uuid NOT NULL,
  role_id uuid NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, role_id),
  CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id)
    REFERENCES roles(id) ON DELETE RESTRICT
);

COMMENT ON TABLE users IS '后台用户主表；本项目只管理资料和角色，登录密码由后续安全项目补充。';
COMMENT ON COLUMN users.id IS '用户 UUID 主键，不承载业务语义。';
COMMENT ON COLUMN users.email IS '用户邮箱；通过 lower(email) 表达不区分大小写的业务唯一性。';
COMMENT ON COLUMN users.display_name IS '用户展示名，可修改，不用于认证。';
COMMENT ON COLUMN users.status IS '用户状态，只允许 ACTIVE 或 DISABLED。';
COMMENT ON COLUMN users.version IS 'JPA 乐观锁版本，防止并发覆盖。';
COMMENT ON COLUMN users.created_at IS '创建时间，统一保存为带时区时间。';
COMMENT ON COLUMN users.updated_at IS '最后更新时间，统一保存为带时区时间。';
COMMENT ON TABLE roles IS '角色主表；角色 code 是程序稳定引用。';
COMMENT ON COLUMN roles.code IS '角色编码，全局唯一，创建后不应随展示名称一起修改。';
COMMENT ON COLUMN roles.status IS '角色状态；停用角色不能再分配给用户。';
COMMENT ON TABLE user_roles IS '用户和角色多对多关系；复合主键防止重复授权。';
COMMENT ON COLUMN user_roles.created_at IS '授权关系建立时间，用于排查权限来源。';
