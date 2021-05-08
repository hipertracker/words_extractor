all: build

build:
	@go build -o sort-me-out *.go

run: build
	@./sort-me-out

test:
	@go test ./... -v -race -cover