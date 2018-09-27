package main

import (
	"log"

	"github.com/joho/godotenv"
	cli "github.com/sotatek-dev/heta/cli/commands"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cli := cli.CLI{}
	cli.Run()
}
