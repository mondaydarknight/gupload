package cmd

import (
	"fmt"

	"github.com/mongche/gupload/config"
	"github.com/mongche/gupload/internal"
	"github.com/urfave/cli/v2"
)

var Serve = cli.Command{
	Name:   "serve",
	Usage:  "Serve the server listen TCP address",
	Action: serveAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "cert",
			Usage: "cert path of TLS certificate",
		},
		&cli.StringFlag{
			Name:  "key",
			Usage: "key path of TLS certificate",
		},
		&cli.IntFlag{
			Name:  "port",
			Value: 1313,
			Usage: "the server port",
		},
	},
}

func serveAction(c *cli.Context) (err error) {
	server, err := internal.NewGrpcServer(config.GrpcServerConfig{
		Cert: c.String("cert"),
		Key:  c.String("key"),
		Port: c.Int("port"),
	})

	if err != nil {
		err = fmt.Errorf("failed to set up the server: %v", err)
		return
	}

	defer server.Close()

	return
}
