PROJECT_NAME=go-wss-chat

SERVER_BINARY=server_bin
CLIENT_BINARY=client_bin
LOG_FILE=server.log

GO=go

# 의존성 정리 및 패키지 설치
ensure-deps:
	@echo "🔄 Ensuring dependencies..."
	cd server && $(GO) mod tidy
	cd client && $(GO) mod tidy

run-server: ensure-deps
	sudo $(GO) run server/server.go > $(LOG_FILE) 2>&1 &

run-client: ensure-deps
	$(GO) run client/client.go

build-server: ensure-deps
	$(GO) build -o $(SERVER_BINARY) server/server.go

build-client: ensure-deps
	$(GO) build -o $(CLIENT_BINARY) client/client.go

build: ensure-deps
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

