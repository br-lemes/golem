PLATFORMS := linux-amd64 linux-arm64 windows-amd64

.PHONY: all clean coverage dev release test version $(PLATFORMS)

TARGET := $(notdir $(shell go list -m 2>/dev/null))
ifeq ($(TARGET),)
	TARGET := $(notdir $(CURDIR))
endif

export CGO_ENABLED=0

ARTIFACTS := $(foreach p,$(PLATFORMS),\
	$(TARGET)-$(p)$(if $(filter windows%,$(p)),.exe))

GOCOVER := github.com/Azure/gocover@latest
CODEGEN := github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
SEMVER := github.com/br-lemes/semver@latest

GENERATED_FILES := pkg/database/enums.json \
	pkg/schemas/params.go \
	pkg/schemas/schemas.go

OPENAPI_FILES := sandbox.json standard.json

sandbox.json: URL := https://api.sandbox.artifactsmmo.com/openapi.json
standard.json: URL := https://api.artifactsmmo.com/openapi.json

all: $(PLATFORMS)

clean:
	$(RM) $(ARTIFACTS) $(OPENAPI_FILES)

coverage: $(GENERATED_FILES) lint
	@go test ./... -coverprofile=coverage.out && \
		sed -i '/cmd\/api/d' coverage.out && \
		go run $(GOCOVER) full --cover-profile=coverage.out

custom-gcl: .custom-gcl.yml
	@golangci-lint custom -v

dev: $(GENERATED_FILES) lint test
	@go build -ldflags "-X 'main.version=$$(date '+%Y-%m-%d %H:%M:%S')'"

lint: custom-gcl
	@if ! git merge-base --is-ancestor master HEAD; then \
		echo "Branch is behind or has diverged from master."; \
		exit 1; \
	fi
	@gofmt -w $$(go list -f '{{.Dir}}/*.go' ./...)
	@./custom-gcl run

pkg/database/enums.json: pkg/database/openapi.json enums.jq
	@jq -f enums.jq pkg/database/openapi.json > $@
	@biome format --write --indent-width 4 $@

pkg/database/openapi.json: $(OPENAPI_FILES) check.jq merge.jq
	@jq -e -s -f check.jq $(OPENAPI_FILES) > /dev/null
	@jq -s -f merge.jq $(OPENAPI_FILES) > $@
	@biome format --write --indent-width 4 $@

pkg/schemas/params.json: pkg/database/openapi.json params.jq
	@jq -f params.jq pkg/database/openapi.json > $@
	@biome format --write --indent-width 4 $@

pkg/schemas/schemas.json: pkg/database/openapi.json schemas.jq
	@jq -f schemas.jq pkg/database/openapi.json > $@
	@biome format --write --indent-width 4 $@

$(OPENAPI_FILES):
	@curl $(URL) -o $@
	@sd -F '"anyOf":[{"type":"boolean"},{"type":"null"}]' \
		'"type":"boolean","nullable":true' $@
	@sd -F '"exclusiveMinimum":0.0' \
		'"minimum":0.0,"exclusiveMinimum":true' $@
	@biome format --write --indent-width 4 $@

$(PLATFORMS): $(GENERATED_FILES) lint test
	@$(eval GOOS := $(word 1,$(subst -, ,$@)))
	@$(eval GOARCH := $(word 2,$(subst -, ,$@)))
	@$(eval OUTPUT := $(TARGET)-$@$(if $(filter windows,$(GOOS)),.exe))
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w" -o $(OUTPUT)

release: version $(PLATFORMS)
	@GOLEM_RELEASE=1 go run $(SEMVER) release $(ARTIFACTS)

test:
	@go test ./...

version: test
	@go run $(SEMVER)

%.go: %.json
	@go run $(CODEGEN) -generate models,skip-prune -package schemas -o $@ $<
	@sd -A '\n\n\t//[^\n]*' '' $@
	@sd -A '^\t*//[^\n]*\n' '' $@
	@sd 'Path\s+\[\]\[\]int' 'Path [][2]int' $@
	@gofmt -w $@
