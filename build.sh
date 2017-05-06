GOOS=linux GOARCH=amd64; go build -o builds/rkd-$GOOS-$GOARCH
GOOS=linux GOARCH=386; go build -o builds/rkd-$GOOS-$GOARCH
GOOS=solaris GOARCH=amd64; go build -o builds/rkd-$GOOS-$GOARCH
GOOS=darwin GOARCH=amd64; go build -o builds/rkd-$GOOS-$GOARCH
GOOS=darwin GOARCH=386; go build -o builds/rkd-$GOOS-$GOARCH
GOOS=freebsd GOARCH=amd64; go build -o builds/rkd-$GOOS-$GOARCH
GOOS=freebsd GOARCH=386; go build -o builds/rkd-$GOOS-$GOARCH
