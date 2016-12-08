# REPO
REPO?=github.com/anarcher/infrakit.gcp

# Package list
PKGS_AND_MOCKS := $(shell go list ./... | grep  ^${REPO}/cmd/)
PKGS := $(shell echo $(PKGS_AND_MOCKS) | tr ' ' '\n' | grep -v /mock$)

build:
	@echo "+ $@"
	@go build ${GO_LDFLAGS} $(PKGS)

clean:
	@echo "+ $@"
	rm -rf build
	mkdir -p build
