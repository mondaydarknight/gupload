package main

import (
	"log"
	"os"

	"github.com/mongche/gupload/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gupload",
		Usage: "the file uploader utility",
		Commands: []*cli.Command{
			&cmd.Serve,
			&cmd.Upload,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
