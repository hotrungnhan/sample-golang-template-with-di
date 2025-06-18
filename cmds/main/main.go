package main

import (
	"github.com/hotrungnhan/surl/cmds"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:    "Shorten Url App",
	Usage:   "surl",
	Version: "1.0.0",
}

func main() {

	app.Commands = []*cli.Command{
		cmds.NewHttpServerCommand(),
		cmds.NewSeederCommand(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run the application")
	}
}
