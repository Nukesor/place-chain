ln -sf "$(pwd)/.tendermint.config.toml" "$HOME/.tendermint/config/config.toml"
ln -sf "$(pwd)/.tendermint.genesis.json" "$HOME/.tendermint/config/genesis.json"
ln -sf "$(pwd)/pre-commit" ".git/hooks/pre-commit"
