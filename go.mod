module github.com/daemosity/go-gator

go 1.24.2

require github.com/daemosity/go-gator/internal/config v0.0.0

replace (
	github.com/daemosity/go-gator/internal/command v0.0.0 => ./internal/command
	github.com/daemosity/go-gator/internal/config v0.0.0 => ./internal/config
)
