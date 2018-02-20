package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	lk "github.com/ulule/loukoum"
)

// Comment model
type Comment struct {
	ID        int64
	Email     string      `db:"mail"`
	Status    string      `db:"status"`
	Message   string      `db:"message"`
	CreatedAt pq.NullTime `db:"deleted_at"`
	DeletedAt pq.NullTime `db:"deleted_at"`
}

// CreateComment creates a comment.
func CreateComment(db *sqlx.DB, comment Comment) (Comment, error) {
	builder := lk.Insert("comments").
		Set(
			lk.Pair("email", comment.Email),
			lk.Pair("status", "waiting"),
			lk.Pair("message", comment.Message),
			lk.Pair("created_at", lk.Raw("NOW()")),
		).
		Returning("id")

	query, args := builder.Prepare()
	// query: INSERT INTO comments (created_at, email, message, status) VALUES (NOW(), :arg_1, :arg_2, :arg_3) RETURNING id
	// args: (map[string]interface {}) (len=3) {
	// (string) (len=5) "arg_1": (string) (len=5) comment.Email,
	// (string) (len=5) "arg_2": (string) (len=7) comment.Message,
	// (string) (len=5) "arg_3": (string) (len=7) "waiting"
	// }

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return comment, err
	}

	err = stmt.Get(&comment, args)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

// UpsertComment insert or update a comment based on email attribute.
func UpsertComment(db *sqlx.DB, comment Comment) (Comment, error) {
	builder := lk.Insert("comments").
		Set(
			lk.Pair("email", comment.Email),
			lk.Pair("status", "waiting"),
			lk.Pair("message", comment.Message),
			lk.Pair("created_at", lk.Raw("NOW()")),
		).
		OnConflict("email", lk.DoUpdate(
			lk.Pair("message", comment.Message),
			lk.Pair("status", "waiting"),
			lk.Pair("created_at", lk.Raw("NOW()")),
			lk.Pair("deleted_at", nil),
		)).
		Returning("id")

	query, args := builder.Prepare()
	// query: INSERT INTO comments (created_at, email, message, status) VALUES (
	//		NOW(), :arg_1, :arg_2, :arg_3
	// ) ON CONFLICT (email) DO UPDATE SET created_at = NOW(), deleted_at = NULL, message = :arg_4, status = :arg_5 RETURNING id
	// args: (map[string]interface {}) (len=5) {
	// (string) (len=5) "arg_1": (string) (len=5) comment.Email,
	// (string) (len=5) "arg_2": (string) (len=7) comment.Message,
	// (string) (len=5) "arg_3": (string) (len=7) "waiting",
	// (string) (len=5) "arg_4": (string) (len=7) comment.Message,
	// (string) (len=5) "arg_5": (string) (len=7) "waiting"
	// }

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return comment, err
	}

	err = stmt.Get(&comment, args)
	if err != nil {
		return comment, err
	}

	return comment, nil
}
