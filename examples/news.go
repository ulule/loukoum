package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	lk "github.com/ulule/loukoum"
)

// News model
type News struct {
	ID          int64
	Status      string      `db:"status"`
	PublishedAt pq.NullTime `db:"deleted_at"`
	DeletedAt   pq.NullTime `db:"deleted_at"`
}

// PublishNews publishes a news.
func PublishNews(db *sqlx.DB, news News) (News, error) {
	builder := lk.Update("news").
		Set(
			lk.Pair("published_at", lk.Raw("NOW()")),
			lk.Pair("status", "published"),
		).
		Where(lk.Condition("id").Equal(news.ID)).
		And(lk.Condition("deleted_at").IsNull(true)).
		Returning("published_at")

	query, args := builder.Prepare()
	// query: UPDATE news SET published_at = NOW(), status = :arg_1 WHERE ((id = :arg_2) AND (deleted_at IS NULL)) RETURNING published_at
	// args: (map[string]interface {}) (len=2) {
	//  (string) (len=5) "arg_1": (string) (len=9) "published",
	//  (string) (len=5) "arg_2": (int) news.ID
	// }

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return news, err
	}

	err = stmt.Get(&news, args)
	if err != nil {
		return news, err
	}

	return news, nil
}
