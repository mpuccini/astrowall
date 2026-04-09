BINARY  := astrowall
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

INSTALL_DIR   := /usr/local/bin
SYSTEMD_DIR   := $(HOME)/.config/systemd/user

.PHONY: all clean build install install-systemd uninstall $(PLATFORMS)

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

install: build
	install -Dm755 $(BINARY) $(INSTALL_DIR)/$(BINARY)
	$(MAKE) install-systemd
	@echo "Installed $(BINARY) to $(INSTALL_DIR) and enabled astrowall.timer"

install-systemd:
	@command -v $(INSTALL_DIR)/$(BINARY) >/dev/null 2>&1 || \
		{ echo "Error: $(INSTALL_DIR)/$(BINARY) not found. Run 'make install' or copy the binary first."; exit 1; }
	@mkdir -p $(SYSTEMD_DIR)
	sed 's|@@BINARY@@|$(INSTALL_DIR)/$(BINARY)|g' systemd/astrowall.service > $(SYSTEMD_DIR)/astrowall.service
	cp systemd/astrowall.timer $(SYSTEMD_DIR)/astrowall.timer
	systemctl --user daemon-reload
	systemctl --user enable --now astrowall.timer
	@echo "Systemd units installed and astrowall.timer enabled"

uninstall:
	systemctl --user disable --now astrowall.timer || true
	rm -f $(SYSTEMD_DIR)/astrowall.service $(SYSTEMD_DIR)/astrowall.timer
	rm -f $(INSTALL_DIR)/$(BINARY)
	systemctl --user daemon-reload
	@echo "Uninstalled $(BINARY) and removed systemd units"

clean:
	rm -rf $(OUTDIR) $(BINARY)
