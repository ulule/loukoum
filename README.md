# Loukoum

[![Documentation][godoc-img]][godoc-url]
![License][license-img]

*A simple SQL Query Builder for PostgreSQL.*

[![Loukoum][loukoum-img]][loukoum-url]

## Introduction

Loukoum is a simple SQL Query Builder for **PostgreSQL**.

If you have to generate complex queries, which relies on various context, **loukoum** is the right tool for you.

Afraid to slip a tiny **SQL injection** manipulating `fmt` to append conditions? **Fear no more**, loukoum is here to protect you against yourself.

Just a few examples when and where loukoum can become handy:

 * Remove user anonymity if user is an admin.
 * Display news draft for an author.
 * Add filters in query if there is query parameters in request.
 * Add a `ON CONFLICT` clause for resource's owner.
 * And so on...

## Installation

```bash
$ dep ensure -add github.com/ulule/loukoum@master
```

## Usage

Loukoum helps you generate SQL queries from composable parts.

However, this is not an ORM or a Mapper.

### Select

```go
import lk "github.com/ulule/loukoum"

stmt := lk.Select("id", "first_name", "last_name", "email").
    From("users").
    Where(lk.Condition("deleted_at").IsNull(true))

users := []User{}
query, args := stmt.Prepare()

err := sqlx.FromContext(ctx).Select(query, args, &users)
if err != nil {
    return nil, err
}

return users, nil
```

### Insert

```go
import lk "github.com/ulule/loukoum"

stmt := lk.Insert("comments").
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

query, args := stmt.Prepare()
return sqlx.FromContext(ctx).Insert(query, args, comment)
```


### Update

```go
import lk "github.com/ulule/loukoum"

stmt := lk.Update("news").
    Set(
        lk.Pair("published_at", lk.Raw("NOW()")),
        lk.Pair("status", "published"),
    ).
    Where(lk.Condition("id").Equal(news.ID)).
    And(lk.Condition("deleted_at").IsNull(true)).
    Returning("published_at")

query, args := stmt.Prepare()
return sqlx.FromContext(ctx).Update(query, args, news)
```

### Delete

```go
import lk "github.com/ulule/loukoum"

stmt := lk.Delete("users").
    Where(lk.Condition("email").Equal(email))

query, args := stmt.Prepare()
return sqlx.FromContext(ctx).Exec(query, args)
```

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
