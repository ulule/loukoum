package main

import (
	"database/sql"

	"github.com/lib/pq"
	lk "github.com/ulule/loukoum"
)

// News model.
type News struct {
	ID          int64       `db:"id"`
	Status      string      `db:"status"`
	PublishedAt pq.NullTime `db:"deleted_at"`
	DeletedAt   pq.NullTime `db:"deleted_at"`
}

// PublishNews publishes a news.
func PublishNews(db *sql.DB, news News) (News, error) {
	builder := lk.Update("news").
		Set(
			lk.Pair("published_at", lk.Raw("NOW()")),
			lk.Pair("status", "published"),
		).
		Where(lk.Condition("id").Equal(news.ID)).
		And(lk.Condition("deleted_at").IsNull(true)).
		Returning("published_at")

	// query: UPDATE news SET published_at = NOW(), status = $1 WHERE ((id = $2) AND (deleted_at IS NULL))
	//        RETURNING published_at
	//  args: []interface{}{
	//            string("published"),
	//            int64(news.ID),
	//        }
	query, args := builder.Query()

	stmt, err := db.Query(query, args...)
	if err != nil {
		return news, err
	}
	defer stmt.Close()

	if stmt.Next() {
		stmt.Scan(&news.PublishedAt)
	}

	err = stmt.Err()
	return news, err
}
