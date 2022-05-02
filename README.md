# gupload
The package provides a unified command line interface to upload file via gRPC protocol.

```console
NAME:
   gupload - the file uploader utility

USAGE:
   gupload [global options] command [command options] [arguments...]

COMMANDS:
   serve    Serve the server listen TCP address
   upload   Upload a file to the server
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

## Getting started
Set up the client and server on the local environment.

1. Install dependencies for go mod.
```console
$ make install
```

2. Regenerate gRPC code via protoc CLI. Please make sure you've installed gRPC and protobuf plugins, [see](https://grpc.io/docs/languages/go/quickstart/#prerequisites)
```console
$ make protoc
```

Upload a file to local gRPC server.
```console
$ gupload serve

$ gupload upload --file=./test.txt
```

If you want to set up TLS connection, please make sure you've settled a self-signed certificate and private key using OpenSSL.
```console
$ make certs SERVICE_NAME=localhost
```
Establish the server connection with TLS certificate.
```console
$ gupload serve --cert=./certs/localhost.cert.pem --key=./certs/localhost.key.pem

$ gupload upload --file=./test.txt --cert=./certs/localhost.cert.pem
```
