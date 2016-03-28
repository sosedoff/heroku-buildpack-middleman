build:
	go build -o ./build/server

release:
	GOOS=linux GOARCH=amd64 go build -o ./build/server_linux_amd64