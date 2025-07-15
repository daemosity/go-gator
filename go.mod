module github.com/daemosity/go-gator

go 1.24.2

require (
	github.com/daemosity/go-gator/internal/config v0.0.0
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/daemosity/go-gator/internal/database v0.0.0
)

replace (
	github.com/daemosity/go-gator/internal/config v0.0.0 => ./internal/config
	github.com/daemosity/go-gator/internal/database v0.0.0 => ./internal/database
)
