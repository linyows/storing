default: build

GOOS?=darwin
GOARCH?=arm64
BUCKET?=linyows-storing

build:
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.defaultBucket=$(BUCKET)" -o storing ./main.go
