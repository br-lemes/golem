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
GOCOVER := github.com/Azure/gocover@latest

all: $(PLATFORMS)

clean:
	$(RM) $(ARTIFACTS)

coverage:
	@go test ./... -coverprofile=coverage.out && \
		sed -i '/cmd\/api/d' coverage.out && \
		go run $(GOCOVER) full --cover-profile=coverage.out

custom-gcl:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint custom -v; \
	else \
		echo "Warning: 'golangci-lint' is not installed. Skipping."; \
	fi

lint: custom-gcl
	@if [ -f ./custom-gcl ]; then ./custom-gcl run; fi

pkg/database/openapi.json:
	@curl https://api.artifactsmmo.com/openapi.json -o $@
	@sd -F '"anyOf":[{"type":"boolean"},{"type":"null"}]' \
		'"type":"boolean","nullable":true' $@
	@biome format --write $@

pkg/schemas/schemas.go: pkg/database/openapi.json
	@oapi-codegen -package schemas -generate models $< > $@
	@sd -A '\n\n\t//[^\n]*' '' $@
	@sd -A '^\t*//[^\n]*\n' '' $@
	@sd 'Path\s+\[\]\[\]int' 'Path [][2]int' $@
	@gofmt -w $@

$(PLATFORMS): pkg/schemas/schemas.go lint test
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
