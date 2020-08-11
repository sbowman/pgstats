# Pg Stats

A simple command-line binary to run analysis queries against a PostgreSQL 
database.  This functionality is based off the 
[Heroku PG Extras](https://github.com/heroku/heroku-pg-extras) and 
[Ecto PSQL Extras](https://github.com/pawurb/ecto_psql_extras).  Some of the
queries have been tweaked a little.

## Build

Use the `Makefile` to build `pg_stats`:

    make 
    
Run the `pg_stats` binary for additional helpful tips and a list of implemented 
commands. 

## Running `pg_stats`

Set the `DB_URI` environment variable for your database connection:

    export DB_URI='postgres://postgres@localhost/mydb?sslmode=disable'
    
Then run `./pg_stats` to see a list of available queries.  Some queries may
take parameters, which can be supplied in the command line.  In the help text
they are listed next to the command in brackets, with their defaults.

Implemented queries:

* all_locks
* cache_hit
* calls (requires `pg_stat_statements` extension)
* index_cache_hit
* index_usage
* locks
* outliers (requires `pg_stat_statements` extension)  
* table_cache_hit

In progress:

* blocking
* total_index_size
* index_size
* table_size
* table_indexes_size
* total_table_size
* unused_indexes
* seq_scans
* long_running_queries
* records_rank
* bloat
* vacuum_stats
* kill_all

## TODO

* Add a docker image, to support deploying into a cloud environment.
