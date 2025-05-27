run: 
	go run cmd/server/main.go

build-server:
	go build -o bin/server cmd/server/main.go

build-node:
	go build -o bin/node cmd/node/main.go

lint: 
	golangci-lint run