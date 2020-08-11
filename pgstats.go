package pgstats

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type CommandFn func(db *sqlx.DB, args ...string) (interface{}, error)

type Command struct {
	Name string
	Description string
	Fn CommandFn
}
// Commands caches the documentation about the available database commands.
var Commands = make(map[string]Command)
