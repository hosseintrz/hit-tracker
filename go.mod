module github.com/hosseintrz/hit-tracker

go 1.18

require (
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/cockroachdb/cockroach-go/v2 v2.2.16
	github.com/labstack/echo/v4 v4.9.1
	github.com/labstack/gommon v0.4.0
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	golang.org/x/crypto v0.0.0-20220517005047-85d78b3ac167 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20220513210249-45d2b4557a2a // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
)

//replace github.com/hosseintrz/hit-tracker/internal => ./internal
replace github.com/hosseintrz/hit-tracker/internal/database => ./internal/database

replace github.com/hosseintrz/hit-tracker/internal/handler => ./internal/handler

replace github.com/hosseintrz/hit-tracker/internal/model => ./internal/model
