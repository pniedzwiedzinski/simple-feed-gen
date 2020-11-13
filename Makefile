GITV != git describe --tags
GITC != git rev-parse --verify HEAD
SRC  != find . -type f -name '*.go' ! -name '*_test.go'
TEST != find . -type f -name '*_test.go'

PREFIX  ?= /usr/local
VERSION ?= $(GITV)
COMMIT  ?= $(GITC)
BUILDER ?= Makefile

GO      := go
INSTALL := install
RM      := rm

sfg: $(SRC)
	$(GO) build -o $@

.PHONY: clean
clean:
	$(RM) -f amfora

.PHONY: install
install: sfg
	install -Dm 755 amfora $(PREFIX)/bin/sfg

.PHONY: uninstall
uninstall:
	$(RM) -f $(PREFIX)/bin/sfg

# Development helpers
.PHONY: fmt
fmt:
	go fmt ./...
