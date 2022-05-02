package cmd

import (
	"fmt"
	"log"

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
			Usage: "private key path of TLS certificate",
		},
		&cli.IntFlag{
			Name:  "port",
			Value: 4443,
			Usage: "the server port",
		},
		&cli.BoolFlag{
			Name:  "http2",
			Usage: "Determine whether enable HTTP/2.0 protocol",
		},
	},
}

func serveAction(c *cli.Context) (err error) {
	var (
		cert   = c.String("cert")
		key    = c.String("key")
		port   = c.Int("port")
		http2  = c.Bool("http2")
		server internal.Server
	)

	log.Printf("Server listening at %d", port)

	switch {
	case http2:
		var http2Server internal.Http2Server
		http2Server, err = internal.NewHttp2Server(config.Http2ServerConfig{
			Cert: cert,
			Key:  key,
			Port: port,
		})

		server = &http2Server
	default:
		var grpcServer internal.GrpcServer
		grpcServer, err = internal.NewGrpcServer(config.GrpcServerConfig{
			Cert: cert,
			Key:  key,
			Port: port,
		})

		server = &grpcServer
	}

	if err != nil {
		err = fmt.Errorf("failed to set up the server: %v", err)
		return
	}

	defer server.Close()

	return
}
