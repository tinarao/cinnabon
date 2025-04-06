package main

import (
	"cinnabon/config"
)

func main() {
	c := config.New("Secret", ":8080")
	c.Load()
}
