package main

import (
	"fmt"
)

func getCommands() commands {
	commands := initCommands()
	commands.register("login", handlerLogin)

	return commands
}

func handlerLogin(s *state, cmd command) error {
	if !cmd.hasArgs() {
		return fmt.Errorf("error: %s requires one (1) argument: [username]", cmd.name)
	}

	given_user := cmd.args[0]

	err := s.config.SetUser(given_user)
	if err != nil {
		return err
	}

	fmt.Printf("INFO: %s has been set as current user.\n", given_user)

	return nil
}
