module github.com/ulule/loukoum/examples

go 1.13

replace github.com/ulule/loukoum/v3 => ../

require (
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.2.0
	github.com/ulule/loukoum/v3 v3.0.0-00010101000000-000000000000
	google.golang.org/appengine v1.6.4 // indirect
)
