target="$(pwd)/.tendermint.config.toml"
echo "$target"
ln -sf "$target" "$HOME/.tendermint/config/config.toml"
