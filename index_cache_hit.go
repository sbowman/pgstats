package pgstats

import "github.com/jmoiron/sqlx"

const QueryIndexCacheHit = `
SELECT
  relname AS name,
  idx_blks_hit AS buffer_hits,
  idx_blks_read AS block_reads,
  idx_blks_hit + idx_blks_read AS total_read,
  CASE (idx_blks_hit + idx_blks_read)::float
    WHEN 0 THEN -1::float
    ELSE idx_blks_hit / (idx_blks_hit + idx_blks_read)::float
  END ratio
FROM
  pg_statio_user_tables
ORDER BY
  idx_blks_hit / (idx_blks_hit + idx_blks_read + 1)::float DESC`

type IndexCacheHitResults struct {
	Hits []CacheHitRow `json:"index_cache_hits"`
}

// IndexCacheHit is the same as cache_hit with each table's indexes cache hit info displayed
// seperately.  Returns -1 for the ratio if there is insufficient data.
func IndexCacheHit(db *sqlx.DB) (IndexCacheHitResults, error) {
	var rates []CacheHitRow
	if err := db.Select(&rates, QueryIndexCacheHit); err != nil {
		return IndexCacheHitResults{}, err
	}

	return IndexCacheHitResults{rates}, nil
}

func init() {
	Commands["index_cache_hit"] = Command{
		Name:        "index_cache_hit",
		Description: "The same as cache_hit with each table's indexes cache hit info displayed seperately.  Returns -1 for the ratio if there is insufficient data.",
		Fn: func(db *sqlx.DB, _ ...string) (interface{}, error) {
			return IndexCacheHit(db)
		},
	}
}
