#!/usr/bin/env sh
rm ~/.tendermint/config/addrbook.json
tendermint unsafe_reset_all
