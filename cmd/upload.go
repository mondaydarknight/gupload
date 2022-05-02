package cmd

import (
	"context"
	"fmt"
	"strings"

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
			Value: "localhost:4443",
			Usage: "the address of the server connect to",
		},
		&cli.IntFlag{
			Name:  "chunk-size",
			Value: 1 << 20,
			Usage: "the size of chunk messages",
		},
		&cli.StringFlag{
			Name:  "cert",
			Usage: "path of TLS certificate",
		},
		&cli.BoolFlag{
			Name:  "compress",
			Usage: "determine whether enable file compression",
		},
		&cli.StringFlag{
			Name:  "file",
			Usage: "relative path of a file",
		},
		&cli.BoolFlag{
			Name:  "http2",
			Usage: "Determine whether switch to HTTP/2.0 protocol",
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
		http2     = c.Bool("http2")
		client    internal.Client
	)

	if address == "" {
		err = fmt.Errorf("address must be required: %v", address)
		return
	}

	if file == "" {
		err = fmt.Errorf("file must be required: %v", file)
		return
	}

	switch {
	case http2:
		var http2Client internal.Http2Client

		if !strings.HasPrefix(address, "https://") {
			address = "https://" + address
		}

		http2Client, err = internal.NewHttp2Client(config.Http2ClientConfig{
			Address: address,
			Cert:    cert,
		})

		client = &http2Client
	default:
		var grpcClient internal.GrpcClient
		grpcClient, err = internal.NewGrpcClient(config.GrpcClientConfig{
			Address:   address,
			ChunkSize: chunkSize,
			Compress:  compress,
			Cert:      cert,
		})

		client = &grpcClient
	}

	if err != nil {
		err = fmt.Errorf("failed to new a client. %v", err)
		return
	}

	defer client.Close()

	err = client.UploadFile(context.Background(), file)

	return
}
