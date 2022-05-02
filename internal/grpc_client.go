package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mongche/gupload/config"
	pb "github.com/mongche/gupload/messages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ServerName   = "localhost"
	MaxChunkSize = (1 << 24)
)

type GrpcClient struct {
	client    pb.GuploadServiceClient
	conn      *grpc.ClientConn
	chunkSize int
}

// Create a new gRPC client that provides to upload a file
func NewGrpcClient(config config.GrpcClientConfig) (c GrpcClient, err error) {
	var (
		grpcCreds credentials.TransportCredentials
		grpcOpts  = []grpc.DialOption{}
	)

	if config.Compress {
		grpcOpts = append(grpcOpts, grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	}

	if config.Cert != "" {
		grpcCreds, err = credentials.NewClientTLSFromFile(config.Cert, ServerName)

		if err != nil {
			err = fmt.Errorf("failed to set up TLS client via root-certificate [%s],", config.Cert)
			return
		}

		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcCreds))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	switch {
	case config.ChunkSize == 0:
		err = errors.New("chunk size must be specified")
		return
	case config.ChunkSize > MaxChunkSize:
		err = errors.New("chunk size must be less than 16MB")
		return
	default:
		c.chunkSize = config.ChunkSize
	}

	c.conn, err = grpc.Dial(config.Address, grpcOpts...)

	if err != nil {
		err = fmt.Errorf("failed to start grpc connection with address, %v", err)
		return
	}

	c.client = pb.NewGuploadServiceClient(c.conn)

	return
}

// Upload a file to the server by the given file path
func (c *GrpcClient) UploadFile(ctx context.Context, file string) (err error) {
	var (
		n      int
		status *pb.UploadStatus
	)

	f, err := os.Open(file)
	if err != nil {
		err = fmt.Errorf("failed to open file, %v", err)
		return
	}
	defer f.Close()

	stream, err := c.client.Upload(ctx)
	if err != nil {
		err = fmt.Errorf("failed to establish stream to upload file, %v", err)
		return
	}
	defer stream.CloseSend()

	buffer := make([]byte, c.chunkSize)

	for {
		n, err = f.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			err = fmt.Errorf("failed to copy file to buffer, %v", err)
			return
		}

		err = stream.Send(&pb.Chunk{Content: buffer[:n]})

		if err != nil {
			err = fmt.Errorf("failed to send chunk via stream, %v", err)
			return
		}
	}

	status, err = stream.CloseAndRecv()

	if err != nil {
		err = fmt.Errorf("failed to receive upstream response, %v", err)
		return
	}

	log.Printf("received response: %v", status)
	return
}

// Close gRPC client connection
func (c *GrpcClient) Close() {
	c.conn.Close()
}
