rgb-alertmanager: main.go
	env GOOS=linux GOARCH=arm GOARM=5 go build -o rgb-alertmanager main.go/usr/local/bin/