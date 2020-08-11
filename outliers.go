package pgstats

import (
	"strconv"

	"github.com/jmoiron/sqlx"
)

const QueryOutliers = `
SELECT 
  query,
  interval '1 millisecond' * total_time AS exec_time,
  (total_time/sum(total_time) OVER()) * 100 AS percent_exec_time,
  calls,
  interval '1 millisecond' * (blk_read_time + blk_write_time) AS sync_io_time
FROM pg_stat_statements WHERE userid = (SELECT usesysid FROM pg_user WHERE usename = current_user LIMIT 1)
ORDER BY total_time DESC
LIMIT $1
`

type OutliersRow struct {
	Query       string  `db:"query" json:"query"`
	TotalTime   string  `db:"exec_time" json:"totalTime"`
	PercentTime float64 `db:"percent_exec_time" json:"percent"`
	Calls       int     `db:"calls" json:"calls"`
	SyncIOTime  string  `db:"sync_io_time"`
}

type OutliersResults struct {
	Outliers []OutliersRow `json:"outliers"`
}

// Outliers displays statements, obtained from pg_stat_statements, ordered by the amount of time to
// execute in aggregate. This includes the statement itself, the total execution time for that
// statement, the proportion of total execution time for all statements that statement has taken up,
// the number of times that statement has been called, and the amount of time that statement spent
// on synchronous I/O (reading/writing from the filesystem).
//
// Typically, an efficient query will have an appropriate ratio of calls to total execution time,
// with as little time spent on I/O as possible. Queries that have a high total execution time but
// low call count should be investigated to improve their performance. Queries that have a high
// proportion of execution time being spent on synchronous I/O should also be investigated.
//
// Requires: CREATE EXTENSION pg_stat_statements;
func Outliers(db *sqlx.DB, limit int) (OutliersResults, error) {
	var outliers []OutliersRow
	if err := db.Select(&outliers, QueryOutliers, limit); err != nil {
		return OutliersResults{}, err
	}

	return OutliersResults{outliers}, nil
}

func init() {
	Commands["outliers"] = Command{
		Name:        "outliers <limit/10>",
		Description: "Displays statements, obtained from pg_stat_statements, ordered by the amount of time to execute in aggregate. This includes the statement itself, the total execution time for that statement, the proportion of total execution time for all statements that statement has taken up, the number of times that statement has been called, and the amount of time that statement spent on synchronous I/O (reading/writing from the filesystem).\n\n Typically, an efficient query will have an appropriate ratio of calls to total execution time, with as little time spent on I/O as possible. Queries that have a high total execution time but low call count should be investigated to improve their performance. Queries that have a high proportion of execution time being spent on synchronous I/O should also be investigated.\n\n Defaults to the top ten queries unless otherwise indicated.\n\n Requires pg_stat_statements module.",
		Fn: func(db *sqlx.DB, args ...string) (interface{}, error) {
			var err error

			limit := 10
			if len(args) > 0 {
				limit, err = strconv.Atoi(args[0])
				if err != nil {
					return nil, err
				}
			}

			return Outliers(db, limit)
		},
	}
}
