hyperion: **/*.go
	env GOOS=linux GOARCH=arm GOARM=5 go build -o hyperion ./cmd/main.go
test:
	go test -v