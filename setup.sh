go get -u -v github.com/satori/go.uuid
go get -u -v github.com/tendermint/abci/cmd/abci-cli
go get -u -v github.com/tendermint/tendermint/cmd/tendermint
cd "$GOPATH/src/github.com/tendermint/tendermint" && git checkout v0.17.1
ln -sf "$(pwd)/.tendermint.genesis.json" "$HOME/.tendermint/config/genesis.json"
ln -sf "$(pwd)/pre-commit" ".git/hooks/pre-commit"

echo "==== ✅  installation successful ===="
