package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	lk "github.com/ulule/loukoum"
)

// User model.
type User struct {
	ID        int64       `db:"id"`
	FirstName string      `db:"first_name"`
	LastName  string      `db:"last_name"`
	Email     string      `db:"email"`
	IsStaff   bool        `db:"is_staff"`
	DeletedAt pq.NullTime `db:"deleted_at"`
}

// FindUsers retrieves non-deleted users
func FindUsers(db *sqlx.DB) ([]User, error) {
	builder := lk.Select("id", "first_name", "last_name", "email").
		From("users").
		Where(lk.Condition("deleted_at").IsNull(true))

	// query: SELECT id, first_name, last_name, email FROM users WHERE (deleted_at IS NULL)
	//  args: map[string]interface{}{
	//
	//        }
	query, args := builder.NamedQuery()

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	users := []User{}

	err = stmt.Select(&users, args)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// DeleteUser deletes a user.
func DeleteUser(db *sqlx.DB, user User) error {
	builder := lk.Delete("users").
		Where(lk.Condition("id").Equal(user.ID))

	// query: DELETE FROM users WHERE (id = :arg_1)
	//  args: map[string]interface{}{
	//            "arg_1": int64(user.ID),
	//        }
	query, args := builder.NamedQuery()

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args)

	return err
}
