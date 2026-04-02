BINARY  := go-apod-bg
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.version=$(VERSION)
OUTDIR  := dist

PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

.PHONY: all clean build $(PLATFORMS)

all: $(PLATFORMS)

$(PLATFORMS):
	$(eval GOOS := $(word 1,$(subst /, ,$@)))
	$(eval GOARCH := $(word 2,$(subst /, ,$@)))
	$(eval EXT := $(if $(filter windows,$(GOOS)),.exe,))
	@echo "Building $(GOOS)/$(GOARCH)..."
	@mkdir -p $(OUTDIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" \
		-o $(OUTDIR)/$(BINARY)-$(GOOS)-$(GOARCH)$(EXT) .

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

clean:
	rm -rf $(OUTDIR) $(BINARY)
