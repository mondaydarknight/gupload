package internal

import "context"

// The client interface defines a group of file upload functionalities
type Client interface {
	// Upload a file to the server by the given file path
	UploadFile(ctx context.Context, file string) (err error)
	// Close the client connection
	Close()
}
