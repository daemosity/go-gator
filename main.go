package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/daemosity/go-gator/internal/config"
)

func main() {
	configObj, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	s := state{
		config: &configObj,
	}

	cmds := getCommands()

	userCmd, err := getUserCommand()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	if err := cmds.run(&s, userCmd); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

}

func getUserInput() []string {
	return os.Args
}

func validateUserInput(rawInput []string) ([]string, error) {
	if len(rawInput) < 2 {
		return nil, errors.New("error: program requires a command name")
	}

	userInput := rawInput[1:]

	return userInput, nil
}

func buildUserCommand(userInput []string) command {
	userCmd := command{
		name: userInput[0],
		args: userInput[1:],
	}

	return userCmd
}

func getUserCommand() (command, error) {
	rawInput := getUserInput()
	userInput, err := validateUserInput(rawInput)
	if err != nil {
		return command{}, err
	}

	userCmd := buildUserCommand(userInput)
	return userCmd, nil

}
