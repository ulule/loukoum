package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	lk "github.com/ulule/loukoum/v3"
)

// Comment model.
type Comment struct {
	ID        int64       `db:"id"`
	Email     string      `db:"email"`
	Status    string      `db:"status"`
	Message   string      `db:"message"`
	UserID    int64       `db:"user_id"`
	CreatedAt pq.NullTime `db:"created_at"`
	DeletedAt pq.NullTime `db:"deleted_at"`
}

// FindComments retrieves comments by users.
func FindComments(db *sqlx.DB) ([]Comment, error) {
	builder := lk.
		Select(
			"comments.id", "comments.email", "comments.status",
			"comments.user_id", "comments.message", "comments.created_at",
		).
		From("comments").
		Join(lk.Table("users"), lk.On("comments.user_id", "users.id")).
		Where(lk.Condition("comments.deleted_at").IsNull(true))

	// query: SELECT comments.id, comments.email, comments.status, comments.user_id, comments.message,
	//        comments.created_at FROM comments INNER JOIN users ON comments.user_id = users.id
	//        WHERE (comments.deleted_at IS NULL)
	//  args: map[string]interface{}{
	//
	//        }
	query, args := builder.NamedQuery()

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	comments := []Comment{}

	err = stmt.Select(&comments, args)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// FindStaffComments retrieves comments by staff users.
func FindStaffComments(db *sqlx.DB) ([]Comment, error) {
	builder := lk.Select("id", "email", "status", "user_id", "message", "created_at").
		From("comments").
		Where(lk.Condition("deleted_at").IsNull(true)).
		Where(
			lk.Condition("user_id").In(
				lk.Select("id").
					From("users").
					Where(lk.Condition("is_staff").Equal(true)),
			),
		)

	// query: SELECT id, email, status, user_id, message, created_at
	//        FROM comments WHERE ((deleted_at IS NULL) AND
	//        (user_id IN (SELECT id FROM users WHERE (is_staff = :arg_1))))
	//  args: map[string]interface{}{
	//            "arg_1": bool(true),
	//        }
	query, args := builder.NamedQuery()

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	comments := []Comment{}

	err = stmt.Select(&comments, args)
	if err != nil {
		return nil, err
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
			lk.Pair("user_id", comment.UserID),
			lk.Pair("created_at", lk.Raw("NOW()")),
		).
		Returning("id")

	// query: INSERT INTO comments (created_at, email, message, status, user_id)
	//        VALUES (NOW(), :arg_1, :arg_2, :arg_3, :arg_4) RETURNING id
	//  args: map[string]interface{}{
	//            "arg_1": string(comment.Email),
	//            "arg_2": string(comment.Message),
	//            "arg_3": string("waiting"),
	//            "arg_4": string(comment.UserID),
	//        }
	query, args := builder.NamedQuery()

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
			lk.Pair("user_id", comment.UserID),
			lk.Pair("created_at", lk.Raw("NOW()")),
		).
		OnConflict("email", lk.DoUpdate(
			lk.Pair("message", comment.Message),
			lk.Pair("user_id", comment.UserID),
			lk.Pair("status", "waiting"),
			lk.Pair("created_at", lk.Raw("NOW()")),
			lk.Pair("deleted_at", nil),
		)).
		Returning("id, created_at")

	// query: INSERT INTO comments (created_at, email, message, status, user_id)
	//        VALUES (NOW(), :arg_1, :arg_2, :arg_3, :arg_4)
	//        ON CONFLICT (email) DO UPDATE SET created_at = NOW(), deleted_at = NULL, message = :arg_5,
	//        status = :arg_6, user_id = :arg_7 RETURNING id, created_at
	//  args: map[string]interface{}{
	//            "arg_1": string(comment.Email),
	//            "arg_2": string(comment.Message),
	//            "arg_3": string("waiting"),
	//            "arg_4": string(comment.UserID),
	//            "arg_5": string(comment.Message),
	//            "arg_6": string("waiting"),
	//            "arg_7": string(comment.UserID),
	//        }
	query, args := builder.NamedQuery()

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
