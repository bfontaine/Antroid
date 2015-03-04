all:
	go build .

check: devdeps
	go vet ./...


devdeps:
	go get golang.org/x/tools/cmd/vet
	go get github.com/golang/lint/golint
