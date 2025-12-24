APP_NAME := geode
CMD_PATH := ./cmd/geode
BINARY := ./bin/$(APP_NAME)

PORT ?= 3001
DIR  ?= content

.PHONY: all build run serve build-site clean

all: build

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BINARY) $(CMD_PATH)/main.go
	@echo "Done."

run: build
	@$(BINARY)

serve: build
	@$(BINARY) serve -dir $(DIR) -port $(PORT)

build-site: build
	@$(BINARY) build -dir $(DIR)

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY)