package cmd

import (
	"context"
	"fmt"

	"github.com/mongche/gupload/config"
	"github.com/mongche/gupload/internal"
	"github.com/urfave/cli/v2"
)

var Upload = cli.Command{
	Name:   "upload",
	Usage:  "Upload a file to the server",
	Action: uploadAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Value: "localhost:1313",
			Usage: "the address of the server connect to",
		},
		&cli.IntFlag{
			Name:  "chunk-size",
			Value: 1 << 20,
			Usage: "the size of chunk messages",
		},
		&cli.StringFlag{
			Name:  "file",
			Usage: "relative path of a file",
		},
		&cli.StringFlag{
			Name:  "cert",
			Usage: "path of TLS certificate",
		},
		&cli.BoolFlag{
			Name:  "compress",
			Usage: "determine whether enable file compression",
		},
	},
}

func uploadAction(c *cli.Context) (err error) {
	var (
		chunkSize = c.Int("chunk-size")
		address   = c.String("address")
		file      = c.String("file")
		cert      = c.String("cert")
		compress  = c.Bool("compress")
	)

	if address == "" {
		err = fmt.Errorf("address must be required: %v", address)
		return
	}

	if file == "" {
		err = fmt.Errorf("file must be required: %v", file)
		return
	}

	client, err := internal.NewGrpcClient(config.GrpcClientConfig{
		Address:   address,
		ChunkSize: chunkSize,
		Compress:  compress,
		Cert:      cert,
	})

	if err != nil {
		err = fmt.Errorf("failed to new a client. %v", err)
		return
	}

	defer client.Close()

	err = client.UploadFile(context.Background(), file)

	return
}
