package dbx

import "embed"

type txKeyType struct{}

var TxKey = txKeyType{}

//go:embed migrations/*.sql
var Migrations embed.FS
