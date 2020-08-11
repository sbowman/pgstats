package pgstats

import "github.com/jmoiron/sqlx"

const QueryIndexUsage = `
SELECT relname,
   CASE idx_scan
     WHEN 0 THEN -1::int
     ELSE (100 * idx_scan / (seq_scan + idx_scan))::int
   END percent_of_times_index_used,
   n_live_tup rows_in_table
 FROM
   pg_stat_user_tables
 ORDER BY
   n_live_tup DESC
`

type IndexUsageRow struct {
	RelName string `db:"relname" json:"relname"`
	Percent int    `db:"percent_of_times_index_used" json:"percentUsed"`
	Rows    int    `db:"rows_in_table" json:"rowsInTable"`
}

type IndexUsageResults struct {
	Indexes []IndexUsageRow `json:"index_usage"`
}

// IndexUsage provides information on the efficiency of indexes, represented as what percentage of
// total scans were index scans. A low percentage can indicate under indexing, or wrong data being
// indexed.
func IndexUsage(db *sqlx.DB) (IndexUsageResults, error) {
	var rows []IndexUsageRow
	if err := db.Select(&rows, QueryIndexUsage); err != nil {
		return IndexUsageResults{}, err
	}

	return IndexUsageResults{rows}, nil
}

func init() {
	Commands["index_usage"] = Command{
		Name:        "index_usage",
		Description: "Provides information on the efficiency of indexes, represented as what percentage of total scans were index scans. A low percentage can indicate under indexing, or wrong data being indexed.  Effective databases are at 99% and up.",
		Fn: func(db *sqlx.DB, _ ...string) (interface{}, error) {
			return IndexUsage(db)
		},
	}
}
