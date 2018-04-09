package app

import (
	"../types"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tendermint/abci/example/code"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"

	// Tendermint http client
	httpcli "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var (
	stateKey        = []byte("stateKey")
	kvPairPrefixKey = []byte("kvPairKey:")
)

// TODO: make configurable
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

var _ abci.Application = (*PlacechainApp)(nil)

type PlacechainApp struct {
	abci.BaseApplication

	state      types.AppState
	httpClient httpcli.HTTP
}

func NewPlacechainApp() *PlacechainApp {
	state := loadState(dbm.NewMemDB())
	httpClient := httpcli.NewHTTP("tcp://0.0.0.0:46657", "/websocket")
	return &PlacechainApp{state: state, httpClient: *httpClient}
}

func (app *PlacechainApp) Info(req abci.RequestInfo) (resInfo abci.ResponseInfo) {
	return abci.ResponseInfo{Data: fmt.Sprintf("{\"size\":%v}", app.state.Size)}
}

func (app *PlacechainApp) PublishTx(tx types.Transaction) (*ctypes.ResultBroadcastTx, error) {
	bytes, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	return app.httpClient.BroadcastTxSync(bytes)
}

func (app *PlacechainApp) DeliverTx(txBytes []byte) abci.ResponseDeliverTx {
	fmt.Println("========================== DELIVER TX")
	var key []byte

	var tx types.TransactionWithType
	json.Unmarshal(txBytes, &tx) // pass tx to get TxType

	var err error
	if tx.Type == types.PIXEL_TRANSACTION {
		var pt types.PixelTransaction
		json.Unmarshal(txBytes, &pt)
		key = []byte(fmt.Sprintf("%d,%d", pt.X, pt.Y))
	} else if tx.Type == types.REGISTER_TRANSACTION {
		var rt types.RegisterTransaction
		json.Unmarshal(txBytes, &rt)
		key, _ = rt.PubKey.MarshalJSON()
	}

	if err != nil {
		return abci.ResponseDeliverTx{Code: code.CodeTypeEncodingError}
	}

	app.state.Db.Set(prefixKey(key), txBytes)
	app.state.Size += 1
	fmt.Println("========================== SCCESFULLY DELIVERED TX")

	return abci.ResponseDeliverTx{Code: code.CodeTypeOK}
}

func (app *PlacechainApp) CheckTx(txBytes []byte) abci.ResponseCheckTx {
	fmt.Println("========================== CHECK TX")
	tx, err := toTransaction(txBytes)

	if err != nil {
		fmt.Println("CheckTx:", err)
		return abci.ResponseCheckTx{Code: code.CodeTypeEncodingError}
	}
	if !tx.IsValid() {
		fmt.Println("CheckTx: invalid transaction content", err)
		return abci.ResponseCheckTx{Code: code.CodeTypeEncodingError}
	}
	return abci.ResponseCheckTx{Code: code.CodeTypeOK}
}

func (app *PlacechainApp) Commit() abci.ResponseCommit {
	fmt.Println("========================== COMMIT")
	// Using a memdb - just return the big endian size of the db
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.Size)
	app.state.AppHash = appHash
	app.state.Height += 1
	saveState(app.state)
	return abci.ResponseCommit{Data: appHash}
}

func (app *PlacechainApp) GetGrid() *types.Grid {
	grid := make(types.Grid, gridsize)
	for i := range grid {
		grid[i] = make([]types.Pixel, gridsize)
	}
	for x := range grid {
		for y := range grid[x] {
			keyString := fmt.Sprintf("%d,%d", x, y)
			key := prefixKey([]byte(keyString))

			if !app.state.Db.Has(key) {
				continue
			}

			// Get color out of key value store and convert it to int
			bytes := app.state.Db.Get(key)
			var pt types.PixelTransaction
			json.Unmarshal(bytes, &pt)
			pixelPublicKey, err := pt.PubKey.MarshalJSON()

			if err != nil {
				fmt.Printf("Error while decoding Public Key for PixelTransaction %v", pt)
			}

			var rt types.RegisterTransaction
			json.Unmarshal(app.state.Db.Get(prefixKey(pixelPublicKey)), &rt)
			grid[x][y] = types.Pixel{Color: pt.Color, Profile: rt.Profile}
		}
	}
	return &grid
}

func (app *PlacechainApp) LoadGenesis(genesisFile string) error {
	// TODO handle custom content in genesis file.
	return nil
}

func toTransaction(txBytes []byte) (types.Transaction, error) {
	var tx types.TransactionWithType
	json.Unmarshal(txBytes, &tx)
	if tx.Type == types.PIXEL_TRANSACTION {
		var pt types.PixelTransaction
		json.Unmarshal(txBytes, &pt)
		return pt, nil
	} else if tx.Type == types.REGISTER_TRANSACTION {
		var rt types.RegisterTransaction
		json.Unmarshal(txBytes, &rt)
		return rt, nil
	}
	return nil, errors.New("Cannot convert []byte to Transaction, unknown type")
}
