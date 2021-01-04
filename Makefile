SHELL := /bin/bash
BUILDPATH=$(CURDIR)
UTILS_PATH=/ldaputils
UTILS_BIN_PATH=$(CURDIR)

# docker parameters
DOCKERCMD=$(shell which docker)

GOBUILDPATHINCONTAINER=/ldaputils
GOBUILDIMAGE=golang:1.14.7


compile_config_utils:
	@echo "compiling binary for config_utils..."
	@echo $(GOBUILDPATHINCONTAINER)
	@$(DOCKERCMD) run --rm -v $(BUILDPATH):$(GOBUILDPATHINCONTAINER) -w $(UTILS_PATH) $(GOBUILDIMAGE) go build -o ldaputils
	@echo "Done."

container:
	@echo "build container"
	#@wget https://github.com/cloudfoundry/bosh-cli/releases/download/v6.4.1/bosh-cli-6.4.1-linux-amd64
	@$(DOCKERCMD) build -t ldaputils:1.0 .
	@echo "Done."

all: compile_config_utils container
