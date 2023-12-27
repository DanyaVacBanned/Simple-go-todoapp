BINARY_NAME=todo-app
.DEFAULT_GOAL := run


build:
	go build -o ./build/${BINARY_NAME} ./cmd/main/main.go

run: build
	./build/${BINARY_NAME}

