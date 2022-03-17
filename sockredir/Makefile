GO := go
GO_BUILD = CGO_ENABLED=0 $(GO) build
GO_GENERATE = $(GO) generate
GO_TAGS ?=
TARGET=BPFDemo
INSTALL = $(QUIET)install
BINDIR ?= /usr/local/bin
VERSION=$(shell git describe --tags --always)

TEST_TIMEOUT ?= 5s
RELEASE_UID ?= $(shell id -u)
RELEASE_GID ?= $(shell id -g)

$(TARGET):
	$(GO_GENERATE)
	$(GO_BUILD) $(if $(GO_TAGS),-tags $(GO_TAGS)) \
		-ldflags "-w -s \
		-X 'github.com/SimpCosm/godemo/sockredir/SockRedir.Version=${VERSION}'"


release:
	docker run \
		--env "RELEASE_UID=$(RELEASE_UID)" \
		--env "RELEASE_GID=$(RELEASE_GID)" \
		--rm \
		--workdir /demo \
		--volume `pwd`:/demo docker.io/library/golang:1.17.6-alpine3.15 \
		sh -c "apk add --no-cache make git clang llvm && make local-release VERSION=${VERSION}"

local-release: clean
	OS=linux; \
	ARCHS='amd64 arm64'; \
	for ARCH in $$ARCHS; do \
		echo Building release binary for $$OS/$$ARCH...; \
		test -d release/$$OS/$$ARCH|| mkdir -p release/$$OS/$$ARCH; \
		$(GO_GENERATE) main_$$ARCH.go; \
		env GOOS=$$OS GOARCH=$$ARCH $(GO_BUILD) $(if $(GO_TAGS),-tags $(GO_TAGS)) -ldflags "-w -s -X 'github.com/SimpCosm/godemo/ebpf/BPFDemo.Version=${VERSION}'" -o release/$$OS/$$ARCH/$(TARGET) ; \
		tar -czf release/$(TARGET)-$$OS-$$ARCH.tar.gz -C release/$$OS/$$ARCH $(TARGET); \
		(cd release && sha256sum $(TARGET)-$$OS-$$ARCH.tar.gz > $(TARGET)-$$OS-$$ARCH.tar.gz.sha256sum); \
		rm -r release/$$OS; \
	done; \
	if [ $$(id -u) -eq 0 -a -n "$$RELEASE_UID" -a -n "$$RELEASE_GID" ]; then \
		chown -R "$$RELEASE_UID:$$RELEASE_GID" release; \
	fi

clean:
	rm -f sockredir 
	rm -f sockredir_bpf*
	rm -rf ./release
