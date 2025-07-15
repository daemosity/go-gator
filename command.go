package main

import (
	"errors"
	"fmt"
)

type command struct {
	name string
	args []string
}

func (c *command) hasArgs() bool {
	return len(c.args) > 0
}

func (c *command) hasName() bool {
	return len(c.name) > 0
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if !cmd.hasName() {
		return errors.New("error: A command must be given")
	}

	cmdName := cmd.name

	f, exists := c.handlers[cmdName]
	if !exists {
		return fmt.Errorf("%s is not a valid command", cmdName)
	}

	err := f(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) error {
	if len(name) < 1 {
		return errors.New("error: A command name must be given")
	}

	c.handlers[name] = f

	return nil
}

func initCommands() commands {
	handlers := make(map[string]func(*state, command) error)
	commands := commands{
		handlers: handlers,
	}

	return commands
}
