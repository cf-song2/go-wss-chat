PROJECT_NAME=go-wss-chat

SERVER_BINARY=server_bin
CLIENT_BINARY=client_bin

GO=go

run-server:
	sudo $(GO) run server/server.go

run-client:
	$(GO) run client/client.go

build-server:
	$(GO) build -o $(SERVER_BINARY) server/server.go

build-client:
	$(GO) build -o $(CLIENT_BINARY) client/client.go

build:
	$(MAKE) build-server
	$(MAKE) build-client

start-server: build-server
	sudo ./$(SERVER_BINARY)

start-client: build-client
	./$(CLIENT_BINARY)

clean:
	rm -f $(SERVER_BINARY) $(CLIENT_BINARY)

