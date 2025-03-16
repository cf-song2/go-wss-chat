PROJECT_NAME=go-wss-chat

SERVER_BINARY=server_bin
LOG_FILE=server.log

GO=go

# 의존성 정리 및 패키지 설치
ensure-deps:
	@echo "🔄 Ensuring dependencies..."
	cd server && $(GO) mod tidy

run-server: ensure-deps
	$(GO) run server/server.go

build-server: ensure-deps
	$(GO) build -o $(SERVER_BINARY) server/server.go

start-server: build-server
	@echo "🟢 Starting server..."
	sudo ./$(SERVER_BINARY) > $(LOG_FILE) 2>&1 & echo $$! > server.pid
	@echo "📝 Logging to $(LOG_FILE)"

stop-server:
	@if [ -f server.pid ]; then \
		echo "🛑 Stopping server..."; \
		sudo kill `cat server.pid`; \
		rm -f server.pid; \
	else \
		echo "⚠️ No server running."; \
	fi

view-log:
	@echo "📜 Viewing server log (Ctrl+C to exit)"
	@tail -f $(LOG_FILE)

clean:
	rm -f $(SERVER_BINARY) $(LOG_FILE) server.pid

