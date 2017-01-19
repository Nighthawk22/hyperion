rgb-alertmanager: main.go
	env GOOS=linux GOARCH=arm go build -o rgb-alertmanager main.go