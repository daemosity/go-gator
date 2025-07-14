package main

import (
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

	err = configObj.SetUser("S C")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	configObj2, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", configObj2)
}
