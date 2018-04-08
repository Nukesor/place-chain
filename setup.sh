#!/usr/bin/env sh
tendermint init
sed -i "s/auth_enc\ =\ true/auth_enc\ =\ false/" $HOME/.tendermint/config/config.toml
sed -i "s/create_empty_blocks\ =\ true/create_empty_blocks\ =\ false/" $HOME/.tendermint/config/config.toml
ln -sf "$(pwd)/tendermint.genesis.json" "$HOME/.tendermint/config/genesis.json"
ln -sf "$(pwd)/pre-commit" ".git/hooks/pre-commit"

echo "==== ✅  installation successful ===="
