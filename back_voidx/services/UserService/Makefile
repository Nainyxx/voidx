BUILD_DIR := build
BIN       := $(BUILD_DIR)/userservice
TEST_BIN  := $(BUILD_DIR)/userservice.test

.PHONY: all build-dir build test test-build run clean

all: build test-build

build-dir:
	mkdir -p $(BUILD_DIR)

build: build-dir
	go build -o $(BIN) .

test-build: build-dir
	go test -c -o $(TEST_BIN) ./userService

# запустить тесты
test:
	go test -race -v ./...


run: build
	./$(BIN)

clean:
	rm -rf $(BUILD_DIR)
