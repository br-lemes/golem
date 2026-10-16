PLATFORMS := linux-amd64 linux-arm64 windows-amd64

.PHONY: all assets clean coverage dev lint release serve test version \
	$(PLATFORMS)

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
BUN_SCRIPT ?= build

API_GENERATED_FILES := \
	pkg/catalog/enums.json \
	pkg/schemas/params.go \
	pkg/schemas/schemas.go

OPENAPI_FILES := sandbox.json standard.json

DEV_HOST ?= localhost
BUILD_TIME := $(shell date '+%Y-%m-%d %H:%M:%S')

sandbox.json: URL := https://api.sandbox.artifactsmmo.com/openapi.json
standard.json: URL := https://api.artifactsmmo.com/openapi.json

all: $(PLATFORMS)

assets: node_modules/.bun-install
	@bun run $(BUN_SCRIPT)

bin/golangci-lint:
	@curl -sSfL https://golangci-lint.run/install.sh | sh -s v2.14.0

bun.lock: package.json
	@bun install --lockfile-only
	@touch $@

clean:
	$(RM) $(ARTIFACTS) $(OPENAPI_FILES)

coverage: lint
	@go test ./... -coverprofile=coverage.out && \
		sed -i '/cmd\/api/d' coverage.out && \
		go run $(GOCOVER) full --cover-profile=coverage.out

dev:
	go tool templ generate --watch \
		--cmd="make serve" \
		--open-browser=false \
		--proxy="http://localhost:8080" \
		--proxybind="$(DEV_HOST)" \
		--watch-pattern='(.+\.(go|templ)|.*/assets/src/.*\.(ts|css))$$' \
		--ignore-pattern='(.+_templ\.go$$)'

$(API_GENERATED_FILES) $(OPENAPI_FILES): node_modules/.bun-install

lint: assets $(API_GENERATED_FILES) bin/golangci-lint
	@if ! git merge-base --is-ancestor master HEAD; then \
		echo "Branch is behind or has diverged from master."; \
		exit 1; \
	fi
	@go tool templ fmt .
	@bun run biome format --write .
	@gofmt -w $$(go list -f '{{.Dir}}/*.go' ./...)
	@./bin/golangci-lint run

node_modules/.bun-install: bun.lock
	@bun install --frozen-lockfile
	@touch $@

pkg/catalog/enums.json: pkg/catalog/openapi.json enums.jq
	@jq -f enums.jq pkg/catalog/openapi.json > $@
	@bun run biome format --write $@

pkg/catalog/openapi.json: $(OPENAPI_FILES) check.jq merge.jq
	@jq -e -s -f check.jq $(OPENAPI_FILES) > /dev/null
	@jq -s -f merge.jq $(OPENAPI_FILES) > $@
	@bun run biome format --write $@

pkg/schemas/params.json: pkg/catalog/openapi.json params.jq
	@jq -f params.jq pkg/catalog/openapi.json > $@
	@bun run biome format --write $@

pkg/schemas/schemas.json: pkg/catalog/openapi.json schemas.jq
	@jq -f schemas.jq pkg/catalog/openapi.json > $@
	@bun run biome format --write $@

$(OPENAPI_FILES):
	@curl $(URL) -o $@
	@sd -F '"anyOf":[{"type":"boolean"},{"type":"null"}]' \
		'"type":"boolean","nullable":true' $@
	@sd -F '"exclusiveMinimum":0.0' \
		'"minimum":0.0,"exclusiveMinimum":true' $@
	@bun run biome format --write $@

$(PLATFORMS): lint test
	@GOOS=$(word 1,$(subst -, ,$@)) \
		GOARCH=$(word 2,$(subst -, ,$@)) \
		go build -ldflags "-s -w" \
			-o $(TARGET)-$@$(if $(filter windows-%,$@),.exe)

release: version $(PLATFORMS)
	@GOLEM_RELEASE=1 go run $(SEMVER) release $(ARTIFACTS)

serve: BUN_SCRIPT = dev
serve: assets $(API_GENERATED_FILES)
	@go build -ldflags="-X 'main.version=$(BUILD_TIME)'"
	@./golem serve

test: assets $(API_GENERATED_FILES)
	@go test ./...

version: test
	@go run $(SEMVER)

%.go: %.json
	@go run $(CODEGEN) -generate models,skip-prune -package schemas -o $@ $<
	@sd -A '\n\n\t//[^\n]*' '' $@
	@sd -A '^\t*//[^\n]*\n' '' $@
	@sd 'Path\s+\[\]\[\]int' 'Path [][2]int' $@
	@gofmt -w $@
