.PHONY: build test release

build:
	go build -o switcher *.go

test:
	go test -tags=${BUILD_TAGS}
	cd middleware;	 go test -tags=${BUILD_TAGS}
	cd operations;	 go test -tags=${BUILD_TAGS}
	cd performers;	 go test -tags=${BUILD_TAGS}
	cd request;	 go test -tags=${BUILD_TAGS}
	cd server;	 go test -tags=${BUILD_TAGS}

release:
	GOOS=linux GOARCH=amd64 go build -o switcher-linux-x64
	GOOS=darwin GOARCH=amd64 go build -o switcher-darwin-x64
	GOOS=windows GOARCH=amd64 go build -o switcher-windows-x64.exe
	tar czvf switcher-linux-x64.tar.gz switcher-linux-x64 README.md LICENSE
	tar czvf switcher-darwin-x64.tar.gz switcher-darwin-x64 README.md LICENSE
	tar czvf switcher-windows-x64.tar.gz switcher-windows-x64.exe README.md LICENSE
