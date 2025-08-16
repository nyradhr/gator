package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/nyradhr/gator/internal/config"
	"github.com/nyradhr/gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	dbURL := cfg.DbUrl
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	s := &state{
		db:     dbQueries,
		config: &cfg,
	}
	cmds := commands{}
	cmds.list = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	input := os.Args
	if len(input) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	cmdName := input[1]
	cmdArgs := input[2:]
	userCmd := command{Name: cmdName, Args: cmdArgs}
	err = cmds.run(s, userCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
