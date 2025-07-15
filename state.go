package main

import (
	"github.com/daemosity/go-gator/internal/config"
	"github.com/daemosity/go-gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}
