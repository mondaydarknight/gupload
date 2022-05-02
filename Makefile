install:
	go install -v

certs:
	openssl genrsa -out ./certs/${SERVICE_NAME}.key.pem 2048
	openssl req \
		-new -x509 \
		-days 3650 \
		-key ./certs/${SERVICE_NAME}.key.pem \
		-out ./certs/${SERVICE_NAME}.cert.pem \
		-subj /CN=localhost \
		-addext "subjectAltName = DNS:localhost"

protoc:
	protoc --proto_path=./messages \
		--go_out=./messages \
		--go_opt=paths=source_relative \
		--go-grpc_out=./messages \
		--go-grpc_opt=paths=source_relative \
		./messages/*.proto

.PHONY: fmt install grpc certs
