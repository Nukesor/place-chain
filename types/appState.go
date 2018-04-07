package types

import (
	dbm "github.com/tendermint/tmlibs/db"
)

type AppState struct {
	Db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}
