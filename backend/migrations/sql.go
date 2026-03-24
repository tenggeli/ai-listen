package migrations

import _ "embed"

//go:embed 000001_init_schema.sql
var InitSchemaSQL string
