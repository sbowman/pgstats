package pgstats

import (
	"time"

	"github.com/jmoiron/sqlx"
)

const QueryLocks = `
SELECT
  pg_stat_activity.pid,
  pg_class.relname,
  pg_locks.transactionid,
  pg_locks.granted,
  pg_locks.mode,
  pg_stat_activity.query AS query_snippet,
  age(now(),pg_stat_activity.query_start) AS "age"
FROM pg_stat_activity,pg_locks left
OUTER JOIN pg_class
  ON (pg_locks.relation = pg_class.oid)
WHERE pg_stat_activity.query <> '<insufficient privilege>'
  AND pg_locks.pid = pg_stat_activity.pid
  AND pg_locks.mode IN ('ExclusiveLock', 'AccessExclusiveLock', 'RowExclusiveLock')
  AND pg_stat_activity.pid <> pg_backend_pid() order by query_start
`

type LocksRow struct {
	PID           string        `db:"pid" json:"PID"`
	RelName       string        `db:"relname" json:"relname"`
	TransactionID int           `db:"transactionid" json:"transactionID"`
	Granted       bool          `db:"granted" json:"granted"`
	Mode          string        `db:"mode" json:"mode"`
	Query         string        `db:"query_snippet" json:"query"`
	Age           time.Duration `db:"age" json:"age"`
}

type LocksResults struct {
	Locks []LocksRow `json:"locks"`
}

// Locks displays queries that have taken out an exlusive lock on a relation. Exclusive locks
// typically prevent other operations on that relation from taking place, and can be a cause of
// "hung" queries that are waiting for a lock to be granted.
func Locks(db *sqlx.DB) (LocksResults, error) {
	var locks []LocksRow
	if err := db.Select(&locks, QueryLocks); err != nil {
		return LocksResults{}, err
	}

	return LocksResults{locks}, nil
}

func init() {
	Commands["locks"] = Command{
		Name:        "locks",
		Description: "Displays queries that have taken out an exlusive lock on a relation. Exclusive locks typically prevent other operations on that relation from taking place, and can be a cause of \"hung\" queries that are waiting for a lock to be granted.",
		Fn: func(db *sqlx.DB, _ ...string) (interface{}, error) {
			return Locks(db)
		},
	}
}
