package main

import (
	"os"

	"github.com/Gandalf-Le-Dev/TodoGo/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Logger setup
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cmd.Execute()
}
