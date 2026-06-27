PLATFORMS := linux-amd64 linux-arm64 windows-amd64

.PHONY: all build clean coverage release test version $(PLATFORMS)

TARGET := $(notdir $(shell go list -m 2>/dev/null))
ifeq ($(TARGET),)
	TARGET := $(notdir $(CURDIR))
endif

export CGO_ENABLED=0

ARTIFACTS := $(foreach p,$(PLATFORMS),\
	$(TARGET)-$(p)$(if $(filter windows%,$(p)),.exe))

SEMVER := github.com/br-lemes/semver@latest

build: pkg/schemas/schemas.go test
	@go build -ldflags "-s -w"

all: $(PLATFORMS)

clean:
	$(RM) $(ARTIFACTS)

coverage:
	@go test ./... -coverprofile=coverage.out && \
		go tool cover -html=coverage.out

pkg/schemas/schemas.go: pkg/database/openapi.json
	@oapi-codegen -package schemas -generate models $< > $@
	@sd -A '\n\n\t//[^\n]*' '' $@
	@sd -A '^\t*//[^\n]*\n' '' $@
	@sd 'Path\s+\[\]\[\]interface\{\}' 'Path [][2]int' $@
	@gofmt -w $@

$(PLATFORMS): pkg/schemas/schemas.go test
	@$(eval GOOS := $(word 1,$(subst -, ,$@)))
	@$(eval GOARCH := $(word 2,$(subst -, ,$@)))
	@$(eval OUTPUT := $(TARGET)-$@$(if $(filter windows,$(GOOS)),.exe))
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w" -o $(OUTPUT)

release: version $(PLATFORMS)
	@go run $(SEMVER) release $(ARTIFACTS)

test:
	@go test ./...

version: test
	@go run $(SEMVER)
