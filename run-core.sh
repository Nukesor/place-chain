#!/usr/bin/env sh
tendermint node --p2p.seeds "46.4.89.126:46656" --log_level "main:info,state:info,p2p:debug,*:error" --consensus.create_empty_blocks false
