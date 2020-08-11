package pgstats

import "github.com/jmoiron/sqlx"

const QueryTableCacheHit = `
SELECT
  relname AS name,
  heap_blks_hit AS buffer_hits,
  heap_blks_read AS block_reads,
  heap_blks_hit + heap_blks_read AS total_read,
  CASE (heap_blks_hit + heap_blks_read)::float
    WHEN 0 THEN -1::float
    ELSE heap_blks_hit / (heap_blks_hit + heap_blks_read)::float
  END ratio
FROM
  pg_statio_user_tables
ORDER BY
  heap_blks_hit / (heap_blks_hit + heap_blks_read + 1)::float DESC
`

type TableCacheHitResults struct {
	Hits []CacheHitRow `json:"table_cache_hits"`
}

// IndexCacheHit is the same as cache_hit with each table's cache hit info displayed seperately.
// Returns -1 for the ratio if there is insufficient data.
func TableCacheHit(db *sqlx.DB) (TableCacheHitResults, error) {
	var rates []CacheHitRow
	if err := db.Select(&rates, QueryTableCacheHit); err != nil {
		return TableCacheHitResults{}, err
	}

	return TableCacheHitResults{rates}, nil
}

func init() {
	Commands["table_cache_hit"] = Command{
		Name:        "table_cache_hit",
		Description: "The same as cache_hit with each table's cache hit info displayed seperately.  Returns -1 for the ratio if there is insufficient data.",
		Fn: func(db *sqlx.DB, _ ...string) (interface{}, error) {
			return TableCacheHit(db)
		},
	}
}

