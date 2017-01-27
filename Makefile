rgb-alertmanager: main.go
	env GOOS=linux GOARCH=arm GOARM=5 go build -o led-alertmanager main.go
