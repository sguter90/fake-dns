TARGET_OSARCHS="linux/amd64 windows/amd64 darwin/amd64 darwin/arm64"
HOST_GOOS=$(shell go env GOHOSTOS)
HOST_GOARCH=$(shell go env GOHOSTARCH)
GOBUILDFLAGS?=-v
OSARCHS?=$(HOST_GOOS)/$(HOST_GOARCH)
BINARY_NAME=fake-dns


build:
	@for arch in $(OSARCHS); do \
		export EXE_SUFFIX=""; \
	    /bin/bash -c "[[ $$arch =~ ^windows/ ]]" && export EXE_SUFFIX=.exe || true; \
		export GOOS=$$(echo $$arch | sed -e 's|/.*||g'); \
		export GOARCH=$$(echo $$arch | sed -e 's|.*/||g'); \
		go install -v ./... && \
		echo "[$$arch] Building binary..."; \
		(go build $(GOBUILDFLAGS) \
			-o=$(CURDIR)/bin/$(BINARY_NAME)_$${GOOS}_$${GOARCH}$${EXE_SUFFIX} \
			./ || exit 1); \
	done

build_all: clean
	$(MAKE) build OSARCHS=$(TARGET_OSARCHS)

run:
	go build -o ./bin/dev ./
	bin/dev root_zones.json ./cache

clean:
	@rm -rf $(PWD)/bin
