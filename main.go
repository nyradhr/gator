package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nyradhr/gator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	s := &state{
		config: &cfg,
	}
	cmds := commands{}
	cmds.list = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)
	input := os.Args
	if len(input) < 2 {
		log.Fatal("Username required for login")
	}
	cmdName := input[1]
	cmdArgs := input[2:]
	userCmd := command{name: cmdName, args: cmdArgs}
	err = cmds.run(s, userCmd)
	if err != nil {
		log.Fatalf("Error during login: %v", err)
	}
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}
