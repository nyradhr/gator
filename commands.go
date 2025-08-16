package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nyradhr/gator/internal/database"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	list map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user does not exist")
		}
		return err
	}
	err = s.config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("name missing")
	}
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	u, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("user with given name already in db: %w", err)
	}
	err = s.config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been created")
	fmt.Printf("User data: %#v", u)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handler := c.list[cmd.Name]
	if handler == nil {
		return errors.New("no handler found for given command")
	}
	err := handler(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.list[name] = f
}
