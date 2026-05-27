package db

import "embed"

//go:embed migrations/*.sql
var Migration embed.FS
