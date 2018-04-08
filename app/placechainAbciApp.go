package app

import (
	"../types"
	"encoding/binary"
	"encoding/json"
	//"errors"
	"fmt"
	"os"

	abcicli "github.com/tendermint/abci/client"
	"github.com/tendermint/abci/example/code"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	// Http client stuff
	httpcli "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
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

var _ abci.Application = (*PlacechainApp)(nil)

type PlacechainApp struct {
	abci.BaseApplication

	state      types.AppState
	client     abcicli.Client
	httpclient httpcli.HTTP
}

func NewPlacechainApp() *PlacechainApp {
	state := loadState(dbm.NewMemDB())
	httpClient := httpcli.NewHTTP("tcp://0.0.0.0:46657", "/websocket")
	return &PlacechainApp{state: state, client: nil, httpclient: *httpClient}
}

func (app *PlacechainApp) StartClient() error {
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

func (app *PlacechainApp) PublishTx(tx types.Transaction) (*ctypes.ResultBroadcastTx, error) {
	bytes, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	return app.httpclient.BroadcastTxSync(bytes)
}

func (app *PlacechainApp) Info(req abci.RequestInfo) (resInfo abci.ResponseInfo) {
	return abci.ResponseInfo{Data: fmt.Sprintf("{\"size\":%v}", app.state.Size)}
}

func (app *PlacechainApp) DeliverTx(tx []byte) abci.ResponseDeliverTx {
	fmt.Println("========================== DELIVER TX")
	var key []byte
	var value []byte

	var message types.Tx
	json.Unmarshal(tx, &message) // pass message to get TxType
	var err error
	if message.Type == types.PIXEL_TRANSACTION {
		var pt types.PixelTransaction
		json.Unmarshal(tx, &pt)
		key = []byte(fmt.Sprintf("%d,%d", pt.X, pt.Y))
		value, err = pt.MarshalJSON()
	} else if message.Type == types.REGISTER_TRANSACTION {
		var rt types.RegisterTransaction
		json.Unmarshal(tx, &rt)
		key, _ = rt.Acc.PubKey.MarshalJSON()
		value, err = rt.MarshalJSON()
	}
	if err != nil {
		return abci.ResponseDeliverTx{Code: code.CodeTypeEncodingError}
	}

	app.state.Db.Set(prefixKey(key), value)
	app.state.Size += 1
	fmt.Println("========================== SCCESFULLY DELIVERED TX")

	return abci.ResponseDeliverTx{Code: code.CodeTypeOK}
}

func (app *PlacechainApp) CheckTx(tx []byte) abci.ResponseCheckTx {
	fmt.Println("========================== CHECK TX")
	if !validateTransactionBytes(tx) {
		fmt.Println("Received malformed transaction payload")
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

func (app *PlacechainApp) Query(reqQuery abci.RequestQuery) (resQuery abci.ResponseQuery) {
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

func (app *PlacechainApp) GetGrid() *types.Grid {
	grid := make(types.Grid, gridsize)
	for i := range grid {
		grid[i] = make([]types.Pixel, gridsize)
	}
	for x := 0; x < gridsize; x++ {
		for y := 0; y < gridsize; y++ {
			keyString := fmt.Sprintf("%d,%d", x, y)
			key := []byte(keyString)

			if app.state.Db.Has(prefixKey(key)) {
				// Get color out of key value store and convert it to int
				bytes := app.state.Db.Get(prefixKey(key))
				var transaction types.PixelTransaction
				json.Unmarshal(bytes, &transaction)
				profileKey, err := transaction.PubKey.MarshalJSON()

				if err != nil {
					fmt.Printf("Error while decoding Public Key for %v", transaction)
				}

				var rt types.RegisterTransaction
				_ = json.Unmarshal(app.state.Db.Get(prefixKey(profileKey)), &rt)
				grid[x][y] = types.Pixel{Color: transaction.Color, Profile: rt.Acc.Profile}
			}
		}
	}
	return &grid
}

func validateTransactionBytes(txBytes []byte) bool {
	var tx types.Tx
	json.Unmarshal(txBytes, &tx)
	isValid := false
	fmt.Println("validate", string(txBytes))
	if tx.Type == types.PIXEL_TRANSACTION {
		var pt types.PixelTransaction
		json.Unmarshal(txBytes, &pt)
		isValid = pt.IsValid()
	} else if tx.Type == types.REGISTER_TRANSACTION {
		var rt types.RegisterTransaction
		json.Unmarshal(txBytes, &rt)
		isValid = rt.IsValid()
	}
	return isValid
}
