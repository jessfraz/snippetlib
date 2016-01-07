# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)

.PHONY: clean all fmt vet lint build test install dbuild drun db
.DEFAULT: default

all: clean build fmt lint test vet install

build:
	@echo "+ $@"
	@go build ./...

fmt:
	@echo "+ $@"
	@gofmt -s -l .

lint:
	@echo "+ $@"
	@golint ./...

test: fmt lint vet
	@echo "+ $@"
	@go test -v ./...

vet:
	@echo "+ $@"
	@go vet ./...

clean:
	@echo "+ $@"
	@rm -rf cliaoke

install:
	@echo "+ $@"
	@go install -v .

dbuild:
	@docker build --rm --force-rm -t jess/snippetlib .

db:
	@docker run -d \
		--name snippets-db \
		-p 3000:3000 \
		-v $(HOME)/snippets-db:/var/lib/postgresql/data \
		postgres:9.3

drun: dbuild
	@docker run --rm -it \
		--net container:snippets-db \
		jess/snippetlib -d --mailchimp-apikey "$(MAILCHIMP_APIKEY)" --mailchimp-listid "$(MAILCHIMP_LISTID)"
