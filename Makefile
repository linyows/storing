default: build

BUCKET?=storing

build:
	go build -ldflags "-X main.defaultBucket=$(BUCKET)" -o storing ./main.go

release:
	goreleaser release --snapshot --rm-dist
