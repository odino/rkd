build:
	go build -o rkd main.go
release:
	rm -rf builds/*
	GOOS=linux GOARCH=amd64; go build -o builds/rkd-$$GOOS-$$GOARCH-`./rkd version`
	GOOS=linux GOARCH=386; go build -o builds/rkd-$$GOOS-$$GOARCH-`./rkd version`
	GOOS=darwin GOARCH=386; go build -o builds/rkd-$$GOOS-$$GOARCH-`./rkd version`
	GOOS=darwin GOARCH=amd64; go build -o builds/rkd-$$GOOS-$$GOARCH-`./rkd version`
	GOOS=windows GOARCH=amd64; go build -o builds/rkd-$$GOOS-$$GOARCH-`./rkd version`
	GOOS=windows GOARCH=386; go build -o builds/rkd-$$GOOS-$$GOARCH-`./rkd version`
	GOOS=freebsd GOARCH=386; go build -o builds/rkd-$$GOOS-$$GOARCH-`./rkd version`
	GOOS=freebsd GOARCH=amd64; go build -o builds/rkd-$$GOOS-$$GOARCH-`./rkd version`
	cd builds && find . -type f ! -name '*.gz' -exec gzip "{}" \;
	ls -la builds/*
install: build
	sudo mv rkd /usr/local/bin
