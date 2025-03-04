PROJECT_NAME=go-wss-chat

SERVER_BINARY=server_bin
CLIENT_BINARY=client_bin
LOG_FILE=server.log

GO=go

run-server:
	sudo $(GO) run server/server.go > $(LOG_FILE) 2>&1 &

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
	sudo ./$(SERVER_BINARY) > $(LOG_FILE) 2>&1 & echo $$! > server.pid

start-client: build-client
	./$(CLIENT_BINARY)

stop-server:
	@if [ -f server.pid ]; then \
		echo "Stopping server..."; \
		sudo kill `cat server.pid`; \
		rm -f server.pid; \
	else \
		echo "No server running."; \
	fi

clean:
	rm -f $(SERVER_BINARY) $(CLIENT_BINARY) $(LOG_FILE) server.pid

