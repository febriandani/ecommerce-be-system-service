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
		--proto_path=./proto/libs/grpc-gateway \
		--go_out=./protogen/golang \
		--go_opt=paths=source_relative \
		--go-grpc_out=./protogen/golang \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=./protogen/golang \
		--grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt=generate_unbound_methods=true \
		./proto/*/*.proto

docker-run:
	@echo "Running Docker container..."
	docker run -d --name backend-system-service \
	  --network fluent-network \
	  -p 8080:8080 \
	  -p 4442:4442 \
	  -p 9992:9992 \
	  -v $(shell pwd)/log:/app/log \
	  -e ELASTIC_APM_SERVER_URL=http://apm-server:8200 \
	  -e ELASTIC_APM_SERVICE_NAME=backend_system_service \
	  -e ELASTIC_APM_ENVIRONMENT=staging \
	  system-service:latest

docker-build:
	@echo "Building Docker image..."
	docker build --no-cache -t system-service . 