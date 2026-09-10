package migrations

import "embed"

// FS 导出内嵌的 SQL 数据库迁移脚本
//go:embed *.sql
var FS embed.FS
