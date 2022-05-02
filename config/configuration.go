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
	// The public cert path of TLS certificate
	Cert string
	// The private key path of TLS certificate
	Key string
	// The server port
	Port int
}

type Http2ClientConfig struct {
	// TCP address connect to the server
	Address string
	// The root certificate path for CA
	Cert string
}

type Http2ServerConfig struct {
	// The public cert path of TLS certificate
	Cert string
	// The private key path of TLS certificate
	Key string
	// the server port connect to the server
	Port int
}
