package migrations

import (
	"embed"
)

//go:embed *.sql
var MigratePostgres embed.FS
