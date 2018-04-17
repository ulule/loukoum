// Package loukoum provides a simple SQL Query Builder.
// At the moment, only PostgreSQL is supported.
//
// If you have to generate complex queries, which rely on various contexts, loukoum is the right tool for you.
// It helps you generate SQL queries from composable parts.
// However, keep in mind it's not an ORM or a Mapper so you have to use a SQL connector
// (like "database/sql" or "sqlx", for example) to execute queries.
//
// If you're afraid to slip a tiny SQL injection manipulating fmt (or a byte buffer...) when you append
// conditions, loukoum is here to protect you against yourself.
//
// For further informations, you can read this documentation:
// https://github.com/ulule/loukoum/blob/master/README.md
//
// Or you can discover loukoum with these examples.
// An "insert" can be generated like that:
//
//   builder := loukoum.Insert("comments").
//       Set(
//           loukoum.Pair("email", comment.Email),
//           loukoum.Pair("status", "waiting"),
//           loukoum.Pair("message", comment.Message),
//           loukoum.Pair("created_at", loukoum.Raw("NOW()")),
//       ).
//       Returning("id")
//
// Also, if you need an upsert, you can define a "on conflict" clause:
//
//   builder := loukoum.Insert("comments").
//       Set(
//           loukoum.Pair("email", comment.Email),
//           loukoum.Pair("status", "waiting"),
//           loukoum.Pair("message", comment.Message),
//           loukoum.Pair("created_at", loukoum.Raw("NOW()")),
//       ).
//       OnConflict("email", loukoum.DoUpdate(
//           loukoum.Pair("message", comment.Message),
//           loukoum.Pair("status", "waiting"),
//           loukoum.Pair("created_at", loukoum.Raw("NOW()")),
//           loukoum.Pair("deleted_at", nil),
//       )).
//       Returning("id")
//
// Updating a news is also simple:
//
//   builder := loukoum.Update("news").
//       Set(
//           loukoum.Pair("published_at", loukoum.Raw("NOW()")),
//           loukoum.Pair("status", "published"),
//       ).
//       Where(loukoum.Condition("id").Equal(news.ID)).
//       And(loukoum.Condition("deleted_at").IsNull(true)).
//       Returning("published_at")
//
// You can remove a specific user:
//
//   builder := loukoum.Delete("users").
//       Where(loukoum.Condition("id").Equal(user.ID))
//
// Or select a list of users...
//
//   builder := loukoum.Select("id", "first_name", "last_name", "email").
//       From("users").
//       Where(loukoum.Condition("deleted_at").IsNull(true))
//
package loukoum
