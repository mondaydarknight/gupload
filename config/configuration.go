package config

type GrpcClientConfig struct {
	// TCP address connect to the server
	Address string
	// The chunk size splitted from the file in bytes
	ChunkSize int
	// Determine whether enable compress the given file
	Compress bool
	// The root certificate path for CA
	Cert string
}

type GrpcServerConfig struct {
	Key  string
	Cert string
	Port int
}
