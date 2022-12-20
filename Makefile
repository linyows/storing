default: test

BUCKET?=storing

build:
	go build -ldflags "-X main.defaultBucket=$(BUCKET)" -o storing ./...

test:
	go test ./...

release:
	goreleaser release --snapshot --rm-dist
