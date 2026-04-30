build:
	go build -o bin/ak820pro ./cmd/ak820pro

test:
	CGO_ENABLED=0 go test ./...

run: build
	./bin/ak820pro

clean:
	rm -rf bin/