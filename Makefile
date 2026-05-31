.PHONY: build clean linux-amd64 linux-arm64 release version windows

TARGET := $(notdir $(shell go list -m 2>/dev/null))
ifeq ($(TARGET),)
    TARGET := $(notdir $(CURDIR))
endif

export CGO_ENABLED=0

build: cmd/schemas.go
	@go build -ldflags "-s -w"

clean:
	$(RM) $(TARGET) $(TARGET)-linux-amd64 $(TARGET)-linux-arm64 $(TARGET).exe

cmd/schemas.go: cmd/openapi.json
	@oapi-codegen -package cmd -generate models $< > $@
	@sd '\n\n\t//[^\n]*' '' $@
	@sd '^\t*//[^\n]*\n' '' $@
	@sd 'Path\s+\[\]\[\]interface\{\}' 'Path [][2]int' $@
	@gofmt -w $@

linux-amd64: cmd/schemas.go
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(TARGET)-linux-amd64

linux-arm64: cmd/schemas.go
	@GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o $(TARGET)-linux-arm64

release: version linux-amd64 linux-arm64 windows
	@go run ./tools/release/main.go

version:
	@go run ./tools/version/main.go

windows: cmd/schemas.go
	@GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $(TARGET).exe
