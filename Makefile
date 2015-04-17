.PHONY: build test release doc vet

default: build

build: vet
	go build -v -o bin/eticket-billing-server eticket-billing-server

run:
	go run eticket-billing-server.go -pidfile=./server.pid -v=2 -environment=development -alsologtostderr=true -config-file=./config.gcfg

test:
	go test -tags=${BUILD_TAGS}
	cd middleware;	 go test -tags=${BUILD_TAGS}
	cd operations;	 go test -tags=${BUILD_TAGS}
	cd request;	 go test -tags=${BUILD_TAGS}
	cd server;	 go test -tags=${BUILD_TAGS}

fmt:
	go fmt eticket-billing-server/config
	go fmt eticket-billing-server/operations
	go fmt eticket-billing-server/request
	go fmt eticket-billing-server/server

release:
	GOOS=linux GOARCH=amd64 go build -o switcher-linux-x64
	GOOS=darwin GOARCH=amd64 go build -o switcher-darwin-x64
	GOOS=windows GOARCH=amd64 go build -o switcher-windows-x64.exe
	tar czvf switcher-linux-x64.tar.gz switcher-linux-x64 README.md LICENSE
	tar czvf switcher-darwin-x64.tar.gz switcher-darwin-x64 README.md LICENSE
	tar czvf switcher-windows-x64.tar.gz switcher-windows-x64.exe README.md LICENSE

doc:
    godoc -http=:6060 -index

vet:
	go vet eticket-billing-server/config
	go vet eticket-billing-server/operations
	go vet eticket-billing-server/request
	go vet eticket-billing-server/server
