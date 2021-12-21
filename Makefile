build:
	(rm -r bin || true) && mkdir bin
	go build -o bin/hotomata-inventory cmd/hotomata-inventory/main.go
	go build -o bin/hotomata cmd/hotomata/*.go
