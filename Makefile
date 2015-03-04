all:
	go build .

check: devdeps
	go vet ./...
	golint ./...


devdeps:
	go get golang.org/x/tools/cmd/vet
	go get github.com/golang/lint/golint
