# Loukoum

[![Documentation][godoc-img]][godoc-url]
![License][license-img]

*A simple SQL Query Builder.*

[![Loukoum][loukoum-img]][loukoum-url]

## Introduction

Loukoum is a simple SQL Query Builder, only **PostgreSQL** is supported at the moment.

If you have to generate complex queries, which relies on various context, **loukoum** is the right tool for you.

Afraid to slip a tiny **SQL injection** manipulating `fmt` to append conditions? **Fear no more**, loukoum is here to protect you against yourself.

Just a few examples when and where loukoum can become handy:

 * Remove user anonymity if user is an admin
 * Display news draft for an author
 * Add filters in query based on request parameters
 * Add a `ON CONFLICT` clause for resource's owner
 * And so on...

## Installation

```bash
$ dep ensure -add github.com/ulule/loukoum@master
```

## Usage

Loukoum helps you generate SQL queries from composable parts.

However, keep in mind it's not an ORM or a Mapper so you have to use a SQL connector ([database/sql][sql-url], [sqlx][sqlx-url], etc.).

### Select

```go
import lk "github.com/ulule/loukoum"

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
```

### Insert

```go
import lk "github.com/ulule/loukoum"

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
	// (string) (len=5) "arg_1": (string) comment.Email,
	// (string) (len=5) "arg_2": (string) comment.Message,
	// (string) (len=5) "arg_3": (string) "waiting"
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
```

### Insert on conflict (upsert)

```go
import lk "github.com/ulule/loukoum"

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
	// (string) (len=5) "arg_1": (string) comment.Email,
	// (string) (len=5) "arg_2": (string) comment.Message,
	// (string) (len=5) "arg_3": (string) "waiting",
	// (string) (len=5) "arg_4": (string) comment.Message,
	// (string) (len=5) "arg_5": (string) "waiting"
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
```

### Update

```go
// News model
type News struct {
	ID          int64
	Status      string      `db:"status"`
	PublishedAt pq.NullTime `db:"deleted_at"`
	DeletedAt   pq.NullTime `db:"deleted_at"`
}

// UpdateNews update a news.
func UpdateNews(db *sqlx.DB, news News) (News, error) {
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
```

### Delete

```go
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
```

See [examples](examples) directory for more information.

## License

This is Free Software, released under the [`MIT License`][software-license-url].

Loukoum artworks are released under the [`Creative Commons BY-SA License`][artwork-license-url].

## Contributing

* Ping us on twitter:
  * [@novln_](https://twitter.com/novln_)
  * [@oibafsellig](https://twitter.com/oibafsellig)
  * [@thoas](https://twitter.com/thoas)
* Fork the [project](https://github.com/ulule/loukoum)
* Fix [bugs](https://github.com/ulule/loukoum/issues)

**Don't hesitate ;)**

[sql-url]: https://golang.org/pkg/database/sql/
[sqlx-url]: https://github.com/jmoiron/sqlx
[loukoum-url]: https://github.com/ulule/loukoum
[loukoum-img]: docs/images/banner.png
[godoc-url]: https://godoc.org/github.com/ulule/loukoum
[godoc-img]: https://godoc.org/github.com/ulule/loukoum?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[software-license-url]: LICENSE
[artwork-license-url]: docs/images/LICENSE
