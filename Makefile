SHELL = /bin/bash -o pipefail

BENCHSTAT := $(GOPATH)/bin/benchstat
BUMP_VERSION := $(GOPATH)/bin/bump_version
RELEASE := $(GOPATH)/bin/github-release
WRITE_MAILMAP := $(GOPATH)/bin/write_mailmap
UNAME = $(shell uname -s)

test:
	go test ./...

ci-install:
	go install honnef.co/go/tools/cmd/staticcheck@latest

ci: lint race-test

race-test: lint
	go test -trimpath -race ./...

lint:
	go vet ./...
	staticcheck ./...

bench: | $(BENCHSTAT)
	go list ./... | grep -v vendor | xargs go test -benchtime=2s -bench=. -run='^$$' 2>&1 | $(BENCHSTAT) /dev/stdin

$(BUMP_VERSION):
	go get -u github.com/kevinburke/bump_version

$(BENCHSTAT):
	go get golang.org/x/perf/cmd/benchstat

$(RELEASE):
	go get -u github.com/aktau/github-release

$(GOPATH)/bin:
	mkdir -p $(GOPATH)/bin

release: race-test | $(BUMP_VERSION) $(RELEASE)
ifndef version
	@echo "Please provide a version"
	exit 1
endif
ifndef GITHUB_TOKEN
	@echo "Please set GITHUB_TOKEN in the environment"
	exit 1
endif
	$(BUMP_VERSION) --version=$(version) main.go
	git push origin --tags
	mkdir -p releases/$(version)
	GOOS=linux GOARCH=amd64 go build -o releases/$(version)/tss-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o releases/$(version)/tss-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o releases/$(version)/tss-windows-amd64 .
	# These commands are not idempotent, so ignore failures if an upload repeats
	$(RELEASE) release --user kevinburke --repo tss --tag $(version) || true
	$(RELEASE) upload --user kevinburke --repo tss --tag $(version) --name tss-linux-amd64 --file releases/$(version)/tss-linux-amd64 || true
	$(RELEASE) upload --user kevinburke --repo tss --tag $(version) --name tss-darwin-amd64 --file releases/$(version)/tss-darwin-amd64 || true
	$(RELEASE) upload --user kevinburke --repo tss --tag $(version) --name tss-windows-amd64 --file releases/$(version)/tss-windows-amd64 || true

force: ;

AUTHORS.txt: force | $(WRITE_MAILMAP)
	$(WRITE_MAILMAP) > AUTHORS.txt

authors: AUTHORS.txt
