build:
	CGO_ENABLED=0 go build -o bin/ak820pro ./cmd/ak820pro

test:
	go test ./...

run: build
	./bin/ak820pro

clean:
	rm -rf bin/
