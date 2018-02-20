package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	lk "github.com/ulule/loukoum"
)

// User model
type User struct {
	ID int64

	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
	DeletedAt pq.NullTime `db:"deleted_at"`
}

// FindUsers retrieves non-deleted users
func FindUsers(db *sqlx.DB) ([]User, error) {
	builder := lk.Select("id", "first_name", "last_name", "email").
		From("users").
		Where(lk.Condition("deleted_at").IsNull(true))

	users := []User{}

	// query: SELECT id, first_name, last_name, email FROM users WHERE (deleted_at IS NULL)
	// args: map[string]interface{}{}
	query, args := builder.Prepare()

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	err = stmt.Select(users, args)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// DeleteUser deletes a user.
func DeleteUser(db *sqlx.DB, user User) error {
	builder := lk.Delete("users").
		Where(lk.Condition("id").Equal(user.ID))

	query, args := builder.Prepare()
	// query: DELETE FROM users WHERE (id = :arg_1)
	// args: (map[string]interface {}) (len=1) {
	//  (string) (len=5) "arg_1": (int) user.ID
	// }

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(args)

	return err
}
