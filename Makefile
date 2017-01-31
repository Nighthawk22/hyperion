hyperion: main.go
	env GOOS=linux GOARCH=arm GOARM=5 go build -o hyperion main.go
