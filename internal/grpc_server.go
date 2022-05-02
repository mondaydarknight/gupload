package internal

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/mongche/gupload/config"
	pb "github.com/mongche/gupload/messages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// The gRPC server is used for file uploading.
type GrpcServer struct {
	pb.UnimplementedGuploadServiceServer
	server *grpc.Server
}

// Upload a file received from streaming.
func (s *GrpcServer) Upload(stream pb.GuploadService_UploadServer) (err error) {
	var chunk *pb.Chunk

	for {
		chunk, err = stream.Recv()

		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.New("failed to read file chunks from the stream unexpectadely")
			return
		}

		log.Printf("received binary data from streaming, %v", chunk.Content)
	}

	err = stream.SendAndClose(&pb.UploadStatus{
		Message: "File uploaded sucessfully",
		Code:    pb.UploadStatusCode_OK,
	})

	if err != nil {
		err = errors.New("failed to respond to the client")
		return
	}

	return
}

// Establish the gRPC server connection to listen for requests.
func NewGrpcServer(config config.GrpcServerConfig) (s GrpcServer, err error) {
	var (
		listener  net.Listener
		grpcCreds credentials.TransportCredentials
		grpcOpts  = []grpc.ServerOption{}
	)

	if config.Port == 0 {
		err = errors.New("port must be specified")
		return
	}

	listener, err = net.Listen("tcp", fmt.Sprintf(":%d", config.Port))

	if err != nil {
		err = fmt.Errorf("failed to listen. %v", err)
		return
	}

	log.Printf("Server listening at %v", listener.Addr())

	if config.Cert != "" && config.Key != "" {
		grpcCreds, err = credentials.NewServerTLSFromFile(config.Cert, config.Key)

		if err != nil {
			err = fmt.Errorf("failed to set up TLS certificate on server using cert: %s, key: %s", config.Cert, config.Key)
			return
		}

		grpcOpts = append(grpcOpts, grpc.Creds(grpcCreds))
	}

	s.server = grpc.NewServer(grpcOpts...)
	pb.RegisterGuploadServiceServer(s.server, &s)
	err = s.server.Serve(listener)

	if err != nil {
		err = fmt.Errorf("failed to listen grpc connection, %v", err)
		return
	}

	return
}

// Close the server connection.
func (s *GrpcServer) Close() {
	if s.server != nil {
		s.server.Stop()
	}
}
