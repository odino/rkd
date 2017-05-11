build:
	go build -o rkd main.go
release:
	rm -rf builds/*
	GOOS=linux GOARCH=amd64; go build -o builds/rkd-$$GOOS-$$GOARCH
	GOOS=linux GOARCH=386; go build -o builds/rkd-$$GOOS-$$GOARCH
	GOOS=darwin GOARCH=386; go build -o builds/rkd-$$GOOS-$$GOARCH
	GOOS=darwin GOARCH=amd64; go build -o builds/rkd-$$GOOS-$$GOARCH
	GOOS=windows GOARCH=amd64; go build -o builds/rkd-$$GOOS-$$GOARCH
	GOOS=windows GOARCH=386; go build -o builds/rkd-$$GOOS-$$GOARCH
	GOOS=freebsd GOARCH=386; go build -o builds/rkd-$$GOOS-$$GOARCH
	GOOS=freebsd GOARCH=amd64; go build -o builds/rkd-$$GOOS-$$GOARCH
	ls -la builds/*
