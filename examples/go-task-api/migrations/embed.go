package migrations

import "embed"

// Files 包含随程序发布的全部数据库迁移，避免部署时依赖外部 SQL 文件路径。
//
//go:embed *.sql
var Files embed.FS
