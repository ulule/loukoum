package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	lk "github.com/ulule/loukoum"
)

// Comment model
type Comment struct {
	ID        int64
	Email     string      `db:"email"`
	Status    string      `db:"status"`
	Message   string      `db:"message"`
	UserID    int64       `db:"user_id"`
	User      *User       `db:"users"`
	CreatedAt pq.NullTime `db:"created_at"`
	DeletedAt pq.NullTime `db:"deleted_at"`
}

// FindComments retrieves comments by users.
func FindComments(db *sqlx.DB, comment Comment) ([]Comment, error) {
	builder := lk.Select("id", "email", "status", "user_id", "message", "created_at").
		From("comments").
		Join(lk.Table("users"), lk.On("comments.user_id", "users.id")).
		Where(lk.Condition("deleted_at").IsNull(true))

	// query: SELECT id, email, status, user_id, message, created_at FROM comments WHERE ((deleted_at IS NULL) AND (user_id IN (SELECT id FROM users WHERE (is_staff IS :arg_1))))
	// args: (map[string]interface {}) (len=1) {
	// (string) (len=5) "arg_1": (bool) true
	// }
	query, args := builder.Prepare()

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	comments := []Comment{}

	err = stmt.Select(&comments, args)
	if err != nil {
		return comments, err
	}

	return comments, nil
}

// FindStaffComments retrieves comments by staff users.
func FindStaffComments(db *sqlx.DB, comment Comment) ([]Comment, error) {
	builder := lk.Select("id", "email", "status", "user_id", "message", "created_at").
		From("comments").
		Where(lk.Condition("deleted_at").IsNull(true)).
		Where(lk.Condition("user_id").In(
			lk.Select("id").
				From("users").
				Where(lk.Condition("is_staff").
					Is(true))))

	// query: SELECT id, email, status, user_id, message, created_at FROM comments WHERE ((deleted_at IS NULL) AND (user_id IN (SELECT id FROM users WHERE (is_staff IS :arg_1))))
	// args: (map[string]interface {}) (len=1) {
	// (string) (len=5) "arg_1": (bool) true
	// }
	query, args := builder.Prepare()

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	comments := []Comment{}

	err = stmt.Select(&comments, args)
	if err != nil {
		return comments, err
	}

	return comments, nil
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
	// (string) (len=5) "arg_1": (string) comment.Email,
	// (string) (len=5) "arg_2": (string) comment.Message,
	// (string) (len=5) "arg_3": (string) "waiting"
	// }

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return comment, err
	}
	defer stmt.Close()

	err = stmt.Get(&comment, args)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

// UpsertComment inserts or updates a comment based on the email attribute.
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
	defer stmt.Close()

	err = stmt.Get(&comment, args)
	if err != nil {
		return comment, err
	}

	return comment, nil
}
