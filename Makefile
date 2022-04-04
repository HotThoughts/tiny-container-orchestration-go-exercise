build:
	go build -o ./watcher ./cmd/watcher
	go build -o ./controller ./cmd/controller

lint:
	golangci-lint run
