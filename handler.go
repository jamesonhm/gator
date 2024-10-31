package main

import (
	"fmt"
)

func handleLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	user := cmd.Args[0]
	err := s.cfg.SetUser(user)
	if err != nil {
		return err
	}
	fmt.Printf("User set to %s\n", user)
	return nil
}
