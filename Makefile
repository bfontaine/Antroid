all:
	go build .

check: deps
	go vet ./...
	go test -v -cover ./...

lint: deps
	golint ./...

deps:
	go get -v -d -t ./...
	go get golang.org/x/tools/cmd/vet
	go get github.com/golang/lint/golint
