# Loukoum

[![CircleCI][circle-img]][circle-url]
[![Documentation][godoc-img]][godoc-url]
![License][license-img]

*A simple SQL Query Builder.*

[![Loukoum][loukoum-img]][loukoum-url]

## Introduction

Loukoum is a simple SQL Query Builder, only **PostgreSQL** is supported at the moment.

If you have to generate complex queries, which rely on various contexts, **loukoum** is the right tool for you.

Afraid to slip a tiny **SQL injection** manipulating `fmt` to append conditions? **Fear no more**, loukoum is here to protect you against yourself.

Just a few examples when and where loukoum can become handy:

 * Remove user anonymity if user is an admin
 * Display news draft for an author
 * Add filters in query based on request parameters
 * Add a `ON CONFLICT` clause for resource's owner
 * And so on...

## Installation

Using [dep](https://github.com/golang/dep)

```console
dep ensure -add github.com/ulule/loukoum@master
```

or `go get`

```console
go get -u github.com/ulule/loukoum
```

## Usage

Loukoum helps you generate SQL queries from composable parts.

However, keep in mind it's not an ORM or a Mapper so you have to use a SQL connector ([database/sql][sql-url], [sqlx][sqlx-url], etc.)
to execute queries.

### INSERT

Insert a new `Comment` and retrieve its `id`.

```go
import lk "github.com/ulule/loukoum"

// Comment model
type Comment struct {
	ID        int64
	Email     string      `db:"email"`
	Status    string      `db:"status"`
	Message   string      `db:"message"`
	UserID    int64       `db:"user_id"`
	CreatedAt pq.NullTime `db:"created_at"`
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

### INSERT on conflict (UPSERT)

```go
import lk "github.com/ulule/loukoum"

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

### UPDATE

Publish a `News` by updating its status and publication date.

```go
// News model
type News struct {
	ID          int64
	Status      string      `db:"status"`
	PublishedAt pq.NullTime `db:"published_at"`
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
```

### SELECT

#### Basic SELECT with an unique condition

Retrieve non-deleted users.

```go
import lk "github.com/ulule/loukoum"

// User model
type User struct {
	ID int64

	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
	IsStaff   bool        `db:"is_staff"`
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

	err = stmt.Select(&users, args)
	if err != nil {
		return nil, err
	}

	return users, nil
}
```

#### SELECT IN with subquery

Retrieve comments only sent by staff users, the staff users query will be a subquery
as we don't want to use any JOIN operations.

```go
// FindStaffComments retrieves comments by staff users.
func FindStaffComments(db *sqlx.DB, comment Comment) ([]Comment, error) {
	builder := lk.Select("id", "email", "status", "user_id", "message", "created_at").
		From("comments").
		Where(lk.Condition("deleted_at").IsNull(true)).
		Where(lk.Condition("user_id").In(lk.Select("id").From("users").Where(lk.Condition("is_staff").Is(true))))

	// query: SELECT id, email, status, user_id, message, created_at FROM comments WHERE ((deleted_at IS NULL) AND (user_id IN (SELECT id FROM users WHERE (is_staff IS :arg_1))))
	// args: (map[string]interface {}) (len=1) {
	// (string) (len=5) "arg_1": (bool) true
	// }
	query, args := builder.Prepare()

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	comments := []Comment{}

	err = stmt.Select(&comments, args)
	if err != nil {
		return comments, err
	}

	return comments, nil
}
```

#### SELECT with JOIN

Retrieve non-deleted comments sent by a user with embedded user in results.

First, we need to update the `Comment` struct to embed `User`.

```go
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
```

Let's create a `FindComments` method to retrieve these comments.

In this scenario we will use an `INNER JOIN` but loukoum also supports `LEFT JOIN` and `RIGHT JOIN`.

```go
// FindComments retrieves comments by users.
func FindComments(db *sqlx.DB, comment Comment) ([]Comment, error) {
	builder := lk.Select("id", "email", "status", "user_id", "message", "created_at").
		From("comments").
		Join(lk.Table("users"), lk.On("comments.user_id", "users.id")).
		Where(lk.Condition("deleted_at").IsNull(true))

	// query: SELECT id, email, status, user_id, message, created_at FROM comments INNER JOIN users ON comments.user_id = users.id WHERE (deleted_at IS NULL)
	// args: (map[string]interface {}) {
	// }
	query, args := builder.Prepare()

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	comments := []Comment{}

	err = stmt.Select(&comments, args)
	if err != nil {
		return comments, err
	}

	return comments, nil
}
```

### DELETE

Delete a user based on ID.

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

## Inspiration

* [squirrel](https://github.com/Masterminds/squirrel)
* [goqu](https://github.com/doug-martin/goqu)
* [sqlabble](https://github.com/minodisk/sqlabble)

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

[loukoum-url]: https://github.com/ulule/loukoum
[loukoum-img]: docs/images/banner.png
[godoc-url]: https://godoc.org/github.com/ulule/loukoum
[godoc-img]: https://godoc.org/github.com/ulule/loukoum?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[software-license-url]: LICENSE
[artwork-license-url]: docs/images/LICENSE
[sql-url]: https://golang.org/pkg/database/sql/
[sqlx-url]: https://github.com/jmoiron/sqlx
[circle-url]: https://circleci.com/gh/ulule/loukoum/tree/master
[circle-img]: https://circleci.com/gh/ulule/loukoum.svg?style=shield&circle-token=1de7bc4fd603b0df406ceef4bbba3fb3d6b5ed10
