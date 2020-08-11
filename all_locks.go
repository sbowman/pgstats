package pgstats

import (
	"github.com/jmoiron/sqlx"
)

const QueryAllLocks = `
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
AND pg_stat_activity.pid <> pg_backend_pid() order by query_start
`

type AllLocksResults struct {
	Locks []LocksRow `json:"all_locks"`
}

// AllLocks displays all the current locks, regardless of their type.
func AllLocks(db *sqlx.DB) (AllLocksResults, error) {
	var locks []LocksRow
	if err := db.Select(&locks, QueryAllLocks); err != nil {
		return AllLocksResults{}, err
	}

	return AllLocksResults{locks}, nil
}

func init() {
	Commands["all_locks"] = Command{
		Name:        "all_locks",
		Description: "Displays all the current locks, regardless of their type.",
		Fn: func(db *sqlx.DB, _ ...string) (interface{}, error) {
			return AllLocks(db)
		},
	}
}
