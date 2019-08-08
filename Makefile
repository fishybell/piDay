build:
	GOOS=linux GOARCH=arm GOARM=5 go build -o app ./...

copy:
	scp -P 20 app pi@192.168.43.124:~/
