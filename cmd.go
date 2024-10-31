package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("no handler function for %s", cmd.Name)
	}
	err := f(s, cmd)
	if err != nil {
		return err
	}
	return nil
}
