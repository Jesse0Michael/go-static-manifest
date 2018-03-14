GO ?= go
COVERAGEDIR = coverage
ifdef CIRCLE_ARTIFACTS
	COVERAGEDIR=$(CIRCLE_ARTIFACTS)/coverage
endif

LDFLAGS = -ldflags '-X main.gitSHA=$(shell git rev-parse HEAD)'

all: build test cover
install-deps:
	glide install
build:
	if [ ! -d bin ]; then mkdir bin; fi
	$(GO) build $(LDFLAGS) -v -o bin/go-static-manifest
fmt:
	find . -not -path "./vendor/*" -name '*.go' -type f | sed 's#\(.*\)/.*#\1#' | sort -u | xargs -n1 -I {} bash -c "cd {} && goimports -w *.go && gofmt -w -l -s *.go"
test:
	if [ ! -d $(COVERAGEDIR) ]; then mkdir $(COVERAGEDIR); fi
	$(GO) test -v ./builder -cover -coverprofile=$(COVERAGEDIR)/builder.coverprofile
cover:
	if [ ! -d $(COVERAGEDIR) ]; then mkdir $(COVERAGEDIR); fi
	$(GO) tool cover -html=$(COVERAGEDIR)/builder.coverprofile -o $(COVERAGEDIR)/builder.html
assert-no-diff:
	test -z "$(shell git status --porcelain)"
clean:
	$(GO) clean
	rm -f bin/go-static-manifest
	rm -rf coverage/
	rm -rf vendor/
