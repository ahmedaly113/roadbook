.DEFAULT_GOAL := run

install-tools:
	go get -u github.com/containous/go-bindata/...
	go get -u github.com/elazarl/go-bindata-assetfs/...
	go get -u -a github.com/vharitonsky/iniflags

.PHONY: roadbook
roadbook:
	go get ./...
	go build -tags=embed -o roadbook roadbook.go

.PHONY: assets
assets:
	go-bindata-assetfs -o generated/bindata.go -pkg generated -tags embed assets/...

roadbook.linux:
	GOOS=linux GOARCH=arm GOARM=7 go build -tags=embed -o roadbook.linux roadbook.go

.PHONY: test
test:
	go test -v ./...

run: roadbook
	./roadbook -config=${USER}.ini

.PHONY: fmt
fmt:
	go fmt ./...
