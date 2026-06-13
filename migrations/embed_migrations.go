// Package migrations содержит встроенные SQL-миграции базы данных.
package migrations

import (
	"embed"
)

//go:embed *.sql
var MigrationsFS embed.FS
