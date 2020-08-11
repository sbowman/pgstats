package pgstats

import "github.com/jmoiron/sqlx"

const QueryCacheHit = `
SELECT
  'index hit rate' AS name,
  (sum(idx_blks_hit)) / nullif(sum(idx_blks_hit + idx_blks_read),0) AS ratio
FROM pg_statio_user_indexes
UNION ALL
SELECT
 'table hit rate' AS name,
  sum(heap_blks_hit) / nullif(sum(heap_blks_hit) + sum(heap_blks_read),0) AS ratio
FROM pg_statio_user_tables
`

type CacheHitRow struct {
	Name       string  `db:"name" json:"name"`
	BufferHits int     `db:"buffer_hits,omitempty" json:"bufferHits"`
	BlockReads int     `db:"block_reads,omitempty" json:"blockReads"`
	TotalRead  int     `db:"total_read,omitempty" json:"totalRead"`
	Ratio      float64 `db:"ratio" json:"ratio"`
}

type CacheHitResults struct {
	Hits []CacheHitRow `json:"cache_hits"`
}

// CacheHit provides information on the efficiency of the buffer cache, for both index reads (index
// hit rate) as well as table reads (table hit rate). A low buffer cache hit ratio can be a sign
// that the Postgres instance is too small for the workload.
func CacheHit(db *sqlx.DB) (CacheHitResults, error) {
	var rates []CacheHitRow
	if err := db.Select(&rates, QueryCacheHit); err != nil {
		return CacheHitResults{}, err
	}

	return CacheHitResults{rates}, nil
}

func init() {
	Commands["cache_hit"] = Command{
		Name:        "cache_hit",
		Description: "Provides information on the efficiency of the buffer cache, for both index reads (index hit rate) as well as table reads (table hit rate). A low buffer cache hit ratio can be a sign that the Postgres instance is too small for the workload.",
		Fn: func(db *sqlx.DB, _ ...string) (interface{}, error) {
			return CacheHit(db)
		},
	}
}
