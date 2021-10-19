package db

import "embed"

//go:generate sqlc generate

//go:embed schema
var Migrations embed.FS
