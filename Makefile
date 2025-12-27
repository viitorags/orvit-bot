-include .env

APP_NAME=$(BOT_NAME)

.PHONY: build run clean

build:
	@echo "Construindo o binário: $(APP_NAME)..."
	@go build -o $(APP_NAME) cmd/bot/main.go

run: build
	@echo "Executando..."
	@./$(APP_NAME)

clean:
	@rm -f $(APP_NAME)
	@echo "Limpeza concluída."
