package app

import (
	"../types"
	"encoding/binary"
	"encoding/json"
	//"errors"
	"fmt"
	"os"
	"strconv"

	abcicli "github.com/tendermint/abci/client"
	"github.com/tendermint/abci/example/code"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

var (
	stateKey        = []byte("stateKey")
	kvPairPrefixKey = []byte("kvPairKey:")
)

const gridsize int = 20

func loadState(db dbm.DB) types.AppState {
	stateBytes := db.Get(stateKey)
	var state types.AppState
	if len(stateBytes) != 0 {
		err := json.Unmarshal(stateBytes, &state)
		if err != nil {
			panic(err)
		}
	}
	state.Db = db
	return state
}

func saveState(state types.AppState) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.Db.Set(stateKey, stateBytes)
}

func prefixKey(key []byte) []byte {
	return append(kvPairPrefixKey, key...)
}

//---------------------------------------------------

var _ abci.Application = (*KVStoreApplication)(nil)

type KVStoreApplication struct {
	abci.BaseApplication

	state  types.AppState
	client abcicli.Client
}

func NewKVStoreApplication() *KVStoreApplication {
	state := loadState(dbm.NewMemDB())
	return &KVStoreApplication{state: state, client: nil}
}

func (app *KVStoreApplication) StartClient() error {
	client, err := abcicli.NewClient("tcp://0.0.0.0:46658", "socket", false)
	if err != nil {
		return err
	}
	app.client = client
	allowLevel, err := log.AllowLevel("debug")
	logger := log.NewFilter(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), allowLevel)
	client.SetLogger(logger.With("module", "abci-client"))
	if err := client.Start(); err != nil {
		return err
	}
	return nil
}

func (app *KVStoreApplication) PublishTx(tx types.Transaction) (res *abci.ResponseDeliverTx, err error) {
	bytes, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	return app.client.DeliverTxSync(bytes)
}

func (app *KVStoreApplication) Info(req abci.RequestInfo) (resInfo abci.ResponseInfo) {
	return abci.ResponseInfo{Data: fmt.Sprintf("{\"size\":%v}", app.state.Size)}
}

func (app *KVStoreApplication) DeliverTx(tx []byte) abci.ResponseDeliverTx {
	fmt.Println("========================== DELIVER TX")
	var key []byte
	var value []byte

	var message types.Tx
	json.Unmarshal(tx, &message) // pass message to get TxType

	if message.Type == types.PIXEL_TRANSACTION {
		var pt types.PixelTransaction
		json.Unmarshal(tx, &pt)
		key = []byte(fmt.Sprintf("%d,%d", pt.X, pt.Y))
		value = []byte(strconv.Itoa(int(pt.Color)))
	} else if message.Type == types.REGISTER_TRANSACTION {
		var rt types.RegisterTransaction
		json.Unmarshal(tx, &rt)
		key, _ = rt.Acc.PubKey.MarshalJSON()
	}

	app.state.Db.Set(prefixKey(key), value)
	app.state.Size += 1

	return abci.ResponseDeliverTx{Code: code.CodeTypeOK}
}

func (app *KVStoreApplication) CheckTx(tx []byte) abci.ResponseCheckTx {
	fmt.Println("========================== CHECK TX")
	var message types.Transaction
	err := json.Unmarshal(tx, &message)
	if err != nil || !message.IsValid() {
		fmt.Println("Received malformed payload", err)
		return abci.ResponseCheckTx{Code: code.CodeTypeEncodingError}
	}
	return abci.ResponseCheckTx{Code: code.CodeTypeOK}
}

func (app *KVStoreApplication) Commit() abci.ResponseCommit {
	fmt.Println("========================== COMMIT")
	// Using a memdb - just return the big endian size of the db
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.Size)
	app.state.AppHash = appHash
	app.state.Height += 1
	saveState(app.state)
	return abci.ResponseCommit{Data: appHash}
}

func (app *KVStoreApplication) Query(reqQuery abci.RequestQuery) (resQuery abci.ResponseQuery) {
	fmt.Println("========================== QUERY")
	value := app.state.Db.Get(prefixKey(reqQuery.Data))
	if reqQuery.Prove {
		resQuery.Index = -1 // TODO make Proof return index
		resQuery.Key = reqQuery.Data
	}
	resQuery.Value = value
	if value != nil {
		resQuery.Log = "exists"
	} else {
		resQuery.Log = "does not exist"
	}
	return
}

func (app *KVStoreApplication) GetGrid() *types.Grid {
	grid := make(types.Grid, gridsize)
	for i := range grid {
		grid[i] = make([]types.Color, gridsize)
	}

	for x := 0; x < gridsize; x++ {
		for y := 0; y < gridsize; y++ {
			keyString := fmt.Sprintf("%d,%d", x, y)
			key := []byte(keyString)

			if app.state.Db.Has(prefixKey(key)) {
				// Get color out of key value store and convert it to int
				colorBytes := app.state.Db.Get(prefixKey(key))
				colorString := string(colorBytes[:])
				color, _ := strconv.Atoi(colorString)
				grid[x][y] = types.DataTypesName[color]
			}
		}
	}
	return &grid
}
