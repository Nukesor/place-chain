CURRENT_HEAD=`git rev-parse HEAD`

## Convenience wrapper

.PHONY: all
all: clean build docker

.PHONY: clean
clean:
	$(RM) -r dist

.PHONY: build
build:
	CGO_ENABLED=0 go build -o dist/place-chain ./cmd/placechainnode/

.PHONY: install
install:
	CGO_ENABLED=0 go install ./cmd/tendermint

.PHONY: container
container: build
	docker build . -t place-chain:${CURRENT_HEAD}