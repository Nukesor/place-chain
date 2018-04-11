GOTOOLS = \
	github.com/golang/dep/cmd/dep
CURRENT_HEAD=`git rev-parse HEAD`
APP_BASE=${HOME}/.place-chain

## Convenience wrapper

.PHONY: all
all: get-tools install-deps clean build

.PHONY: clean
clean:
	@echo "--> Clean dist"
	$(RM) -r dist

.PHONY: build
build:
	@echo "--> Go build place-chain"
	CGO_ENABLED=0 go build -o dist/place-chain ./cmd/placechainnode/

.PHONY: install
install:
	@echo "--> Go install place-chain"
	CGO_ENABLED=0 go install ./cmd/placechainnode

.PHONY: container
container: build
	@echo "--> Docker build"
	docker build . -t place-chain:${CURRENT_HEAD}

.PHNOY: get-tools
get-tools:
	@echo "--> Updating tools"
	go get -u -v $(GOTOOLS)

.PHONY: install-deps
install-deps:
	@rm -rf vendor/
	@echo "--> Running dep"
	@dep ensure -v