COVERFILE=count.out

all:
	go build .

check: deps
	go vet ./...
	go test -v -cover ./...

covercheck: deps
	go test -covermode=count -coverprofile=$(COVERFILE) ./api
	go tool cover -html=$(COVERFILE)

lint: deps
	golint ./...

deps:
	go get -v -d -t ./...
	go get golang.org/x/tools/cmd/vet
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover

clean:
	$(RM) *~ $(COVERFILE)


# Report

REPORT = rapport.pdf
REPORT_SOURCE = docs/rapport.md
REPORT_TPL = docs/template.tex

report: $(REPORT)

$(REPORT): $(REPORT_SOURCE) $(REPORT_TPL)
	pandoc --template $(REPORT_TPL) --latex-engine=xelatex $(REPORT_SOURCE) -o $@
