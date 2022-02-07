BINARY_NAME=main

all: clean test build

build:
	@go build -o ${BINARY_NAME} -ldflags "-s -w" *.go

run: build
	./${BINARY_NAME}

run-sort: build
	./${BINARY_NAME} -n 10 -s

test:
	@go test ./... -v

coverage:
	@go test ./... -v -coverprofile=coverage.out  


cover: test
	@go tool cover -html=coverage.out

clean:
	@go clean
	rm -f coverage.out 
	rm -f ./${BINARY_NAME}
	rm -rf ./words

