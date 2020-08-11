package pgstats

import (
	"strconv"

	"github.com/jmoiron/sqlx"
)

const QueryCalls = `
SELECT 
  query,
  interval '1 millisecond' * total_time AS exec_time,
  (total_time/sum(total_time) OVER()) * 100 AS percent_exec_time,
  calls,
  interval '1 millisecond' * (blk_read_time + blk_write_time) AS sync_io_time
FROM pg_stat_statements WHERE userid = (SELECT usesysid FROM pg_user WHERE usename = current_user LIMIT 1)
ORDER BY calls DESC LIMIT $1
`

// Calls is much like outliers, but ordered by the number of times a statement has been called.
//
// Requires: CREATE EXTENSION pg_stat_statements;
func Calls(db *sqlx.DB, limit int) (OutliersResults, error) {
	var outliers []OutliersRow
	if err := db.Select(&outliers, QueryCalls, limit); err != nil {
		return OutliersResults{}, err
	}

	return OutliersResults{outliers}, nil
}

func init() {
	Commands["calls"] = Command{
		Name:        "calls <limit/10>",
		Description: "Much like outliers, but ordered by the number of times a statement has been called.\n\n Requires pg_stat_statements module.",
		Fn: func(db *sqlx.DB, args ...string) (interface{}, error) {
			var err error

			limit := 10
			if len(args) > 0 {
				limit, err = strconv.Atoi(args[0])
				if err != nil {
					return nil, err
				}
			}

			return Calls(db, limit)
		},
	}
}

