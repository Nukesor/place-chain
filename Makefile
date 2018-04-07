.PHONY: setup
setup:
	./install-dotfiles.sh
	go get -u -v github.com/satori/go.uuid
	go get -u -v github.com/tendermint/abci/cmd/abci-cli
	go get -u -v github.com/tendermint/tendermint/cmd/tendermint
	@echo "==== ✅  installation successful ===="
