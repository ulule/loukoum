package main

import (
	"database/sql"

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
func FindComments(db *sql.DB) ([]Comment, error) {
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
	//  args: []interface{}{
	//
	//        }
	query, args := builder.Query()

	stmt, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	comments := []Comment{}

	for stmt.Next() {
		comment := Comment{}
		err = stmt.Scan(
			&comment.ID, &comment.Email, &comment.Status,
			&comment.UserID, &comment.Message, &comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	err = stmt.Err()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// FindStaffComments retrieves comments by staff users.
func FindStaffComments(db *sql.DB) ([]Comment, error) {
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
	//        (user_id IN (SELECT id FROM users WHERE (is_staff = $1))))
	//  args: []interface{}{
	//            bool(true),
	//        }
	query, args := builder.Query()

	stmt, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	comments := []Comment{}

	for stmt.Next() {
		comment := Comment{}
		err = stmt.Scan(
			&comment.ID, &comment.Email, &comment.Status,
			&comment.UserID, &comment.Message, &comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	err = stmt.Err()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// CreateComment creates a comment.
func CreateComment(db *sql.DB, comment Comment) (Comment, error) {
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
	//        VALUES (NOW(), $1, $2, $3, $4) RETURNING id
	//  args: []interface{}{
	//            string(comment.Email),
	//            string(comment.Message),
	//            string("waiting"),
	//            int64(comment.UserID),
	//        }
	query, args := builder.Query()

	stmt, err := db.Query(query, args...)
	if err != nil {
		return comment, err
	}
	defer stmt.Close()

	for stmt.Next() {
		err = stmt.Scan(&comment.ID)
		if err != nil {
			return comment, err
		}
	}

	err = stmt.Err()
	return comment, err
}

// UpsertComment inserts or updates a comment based on the email attribute.
func UpsertComment(db *sql.DB, comment Comment) (Comment, error) {
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

	// query: INSERT INTO comments (created_at, email, message, status, user_id) VALUES (NOW(), $1, $2, $3, $4)
	//        ON CONFLICT (email) DO UPDATE SET created_at = NOW(), deleted_at = NULL, message = $5, status = $6,
	//        user_id = $7 RETURNING id, created_at
	//  args: []interface{}{
	//            string(comments.Email),
	//            string(comments.Message),
	//            string("waiting"),
	//            int64(comment.UserID),
	//            string(comments.Message),
	//            string("waiting"),
	//            int64(comment.UserID),
	//        }
	query, args := builder.Query()

	stmt, err := db.Query(query, args...)
	if err != nil {
		return comment, err
	}
	defer stmt.Close()

	for stmt.Next() {
		err = stmt.Scan(&comment.ID, &comment.CreatedAt)
		if err != nil {
			return comment, err
		}
	}

	err = stmt.Err()
	return comment, err
}
