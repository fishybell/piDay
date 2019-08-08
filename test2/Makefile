default: build copy

build:
	GOOS=linux GOARCH=arm GOARM=7 go build -o app ./...

copy:
	scp -P 20 app pi@192.168.43.124:~/
