#!/usr/bin/env sh
go get -u -v github.com/satori/go.uuid
go get -u -v github.com/tendermint/abci/cmd/abci-cli
go get -u -v github.com/tendermint/tendermint/cmd/tendermint

cd "$GOPATH/src/github.com/tendermint/tendermint/cmd/tendermint" && git checkout v0.16.0 \
	&& go build && mv tendermint "$GOPATH/bin/"
