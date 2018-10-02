BINARIES := bin/hyperion, bin/hyperion-arm
all: $(BINARIES)

clean:
	rm -rf bin

bin/hyperion: $(shell find . -name '*.go')
	go generate ./...
	cd cmd/hyperion && go build -o ../../$@

bin/hyperion-arm: $(shell find . -name '*.go')
	go generate ./...
	cd cmd/hyperion && env GOOS=linux GOARCH=arm GOARM=5 go build -o ../../$@

test:
	go generate ./...
	go test ./...

prepare:
	go generate ./...
	go mod vendor


.PHONY: all
.PHONY: clean
.PHONY: test