module github.com/daemosity/go-gator

go 1.24.2

require (
	github.com/daemosity/go-gator/internal/config v0.0.0
	github.com/daemosity/go-gator/internal/database v0.0.0
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
)

require (
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	golang.org/x/net v0.26.0 // indirect
)

replace (
	github.com/daemosity/go-gator/internal/config v0.0.0 => ./internal/config
	github.com/daemosity/go-gator/internal/database v0.0.0 => ./internal/database
)
