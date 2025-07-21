# Variables
GO = go
CLIENT_DIR = cmd/client
SERVER_DIR = cmd/server
CLIENT_MAIN = $(CLIENT_DIR)/main.go
SERVER_MAIN = $(SERVER_DIR)/main.go

# Run the server
run-server:
	@echo "Running the server..."
	$(GO) run $(SERVER_MAIN)

# Run the client
run-client:
	@echo "Running the client..."
	$(GO) run $(CLIENT_MAIN)

# Generate protobuf
protoc:
	@echo "Generating protobuf files..."
	# root proto path
	protoc \
		--proto_path=./ \
		--proto_path=./proto/libs \
		--proto_path=./proto/ \
		--proto_path=./proto/system \
		--go_out=./protogen/golang \
		--go_opt=paths=source_relative \
		--go-grpc_out=./protogen/golang \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=./protogen/golang \
		--grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt=generate_unbound_methods=true \
		./proto/*/*.proto