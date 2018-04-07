target="$(pwd)/.tendermint.config.toml"
ln -sf "$target" "$HOME/.tendermint/config/config.toml"
ln -sf "$(pwd)/pre-commit" ".git/hooks/pre-commit"
