package main

import (
	"fmt"

	"github.com/jamesonhm/gator/internal/config"
)

func main() {
	c := config.Read()
	fmt.Println(c)
	err := c.SetUser("jhm")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config.Read())
}
