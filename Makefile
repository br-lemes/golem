PLATFORMS := linux-amd64 linux-arm64 windows-amd64

.PHONY: build clean release version $(PLATFORMS)

TARGET := $(notdir $(shell go list -m 2>/dev/null))
ifeq ($(TARGET),)
    TARGET := $(notdir $(CURDIR))
endif

ARTIFACTS := $(foreach p,$(PLATFORMS),\
    $(TARGET)-$(p)$(if $(filter windows%,$(p)),.exe))

SEMVER := github.com/br-lemes/semver@latest

build: pkg/schemas/schemas.go
	@go build -ldflags "-s -w"

clean:
	$(RM) $(ARTIFACTS)

pkg/schemas/schemas.go: pkg/database/openapi.json
	@oapi-codegen -package cmd -generate models $< > $@
	@sd '\n\n\t//[^\n]*' '' $@
	@sd '^\t*//[^\n]*\n' '' $@
	@sd 'Path\s+\[\]\[\]interface\{\}' 'Path [][2]int' $@
	@gofmt -w $@

$(PLATFORMS): pkg/schemas/schemas.go
	@$(eval GOOS := $(word 1,$(subst -, ,$@)))
	@$(eval GOARCH := $(word 2,$(subst -, ,$@)))
	@$(eval OUTPUT := $(TARGET)-$@$(if $(filter windows,$(GOOS)),.exe))
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags "-s -w" -o $(OUTPUT)

release: version $(PLATFORMS)
	@go run $(SEMVER) release $(ARTIFACTS)

version:
	@go run $(SEMVER)
