package main

import (
	"database/sql"

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

// FindUsers retrieves non-deleted users.
func FindUsers(db *sql.DB) ([]User, error) {
	builder := lk.Select("id", "first_name", "last_name", "email").
		From("users").
		Where(lk.Condition("deleted_at").IsNull(true))

	// query: SELECT id, first_name, last_name, email FROM users WHERE (deleted_at IS NULL)
	//  args: []interface{}{
	//
	//        }
	query, args := builder.Query()

	stmt, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	users := []User{}

	for stmt.Next() {
		user := User{}
		err = stmt.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	err = stmt.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// DeleteUser deletes a user.
func DeleteUser(db *sql.DB, user User) error {
	builder := lk.Delete("users").
		Where(lk.Condition("id").Equal(user.ID))

	// query: DELETE FROM users WHERE (id = $1)
	//  args: []interface{}{
	//            int64(user.ID),
	//        }
	query, args := builder.Query()

	_, err := db.Exec(query, args...)
	return err
}
