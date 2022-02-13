# Session 02
## Tech stack for rewrite
- Go 1.17
- Gin framework for HTTP
- gin/contrib for sessions
- go-sqlite3

## Changes
We have attempted to do it as 1:1 as possible, but in moving from Python to Go, there are some things that have to be different.

- `queryDb` no longer has a `one` parameter because that would require returning different types.
- `@app.after_request` is not a separate function. It is handled by `defer db.Close()` and the `c.Next()`.
- `login` is split into two separate `GET` and `POST` handlers
- Instead of `url_for`, we are using string constants

## Things to improve/fix later
- `useDbAndUser` should be split into two separate middlewares
- `r.Use(beforeRequest)` connects and queries the database for all users - should only on a request-basis or have a global connection
- Use a ORM, forexample GORM
    - Or for making our simple crap code better, helper functions to create structs from a Row map
- Currently, status code 307 is used in many redirects - investigate if this is proper usage of `c.Redirect` (or if it matches the old minitwit)
- Use an external frontend. The html template stuff is not working well
- Split things into multiple files, one large main.go file is confusing
- Storing any cookies in the containerized app, results in a runtime error the next time the container is run
    - This happens when you create a new user, and sign in

## Other notes
- Remember to remove `ping`, `pingHandler` etc. that was used for testing
