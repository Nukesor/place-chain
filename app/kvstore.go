package app

import (
	"../types"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	abcicli "github.com/tendermint/abci/client"
	"github.com/tendermint/abci/example/code"
	abci "github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

var (
	stateKey        = []byte("stateKey")
	kvPairPrefixKey = []byte("kvPairKey:")
)

const gridsize int = 20

type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

func loadState(db dbm.DB) State {
	stateBytes := db.Get(stateKey)
	var state State
	if len(stateBytes) != 0 {
		err := json.Unmarshal(stateBytes, &state)
		if err != nil {
			panic(err)
		}
	}
	state.db = db
	return state
}

func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}

func prefixKey(key []byte) []byte {
	return append(kvPairPrefixKey, key...)
}

//---------------------------------------------------

var _ abci.Application = (*KVStoreApplication)(nil)

type KVStoreApplication struct {
	abci.BaseApplication

	state  State
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

func (app *KVStoreApplication) SetPixel(x uint8, y uint8) (res *abci.ResponseDeliverTx, err error) {
	return app.client.DeliverTxSync([]byte("lel"))
}

func (app *KVStoreApplication) Info(req abci.RequestInfo) (resInfo abci.ResponseInfo) {
	return abci.ResponseInfo{Data: fmt.Sprintf("{\"size\":%v}", app.state.Size)}
}

// tx is either "key=value" or just arbitrary bytes
func (app *KVStoreApplication) DeliverTx(tx []byte) abci.ResponseDeliverTx {
	fmt.Println("========================== DELIVER TX")
	var message types.Transaction
	json.Unmarshal(tx, &message)

	keyString := fmt.Sprintf("%d,%d", message.X, message.Y)
	key := []byte(keyString)

	app.state.db.Set(prefixKey(key), []byte(strconv.Itoa(int(message.Color))))
	app.state.Size += 1

	tags := []cmn.KVPair{
		{[]byte("app.creator"), []byte("jae")},
		{[]byte("app.key"), key},
	}
	app.GetGrid()
	return abci.ResponseDeliverTx{Code: code.CodeTypeOK, Tags: tags}
}

func (app *KVStoreApplication) CheckTx(tx []byte) abci.ResponseCheckTx {
	fmt.Println("========================== CHECK TX")
	valid, message := validatePayload(tx)
	if !valid {
		fmt.Println("========================== INVALID TX")
		fmt.Println(message)
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

	if reqQuery.Prove {
		value := app.state.db.Get(prefixKey(reqQuery.Data))
		resQuery.Index = -1 // TODO make Proof return index
		resQuery.Key = reqQuery.Data
		resQuery.Value = value
		if value != nil {
			resQuery.Log = "exists"
		} else {
			resQuery.Log = "does not exist"
		}
		return
	} else {
		value := app.state.db.Get(prefixKey(reqQuery.Data))
		resQuery.Value = value
		if value != nil {
			resQuery.Log = "exists"
		} else {
			resQuery.Log = "does not exist"
		}
		return
	}
}

func (app *KVStoreApplication) GetGrid() *types.Grid {
	grid := make(types.Grid, gridsize)
	for i := range grid {
		grid[i] = make([]types.Color, gridsize)
	}

	fmt.Println(grid)

	for x := 0; x < gridsize; x++ {
		for y := 0; y < gridsize; y++ {
			keyString := fmt.Sprintf("%d,%d", x, y)
			key := []byte(keyString)

			if app.state.db.Has(prefixKey(key)) {
				// Get color out of key value store and convert it to int
				colorBytes := app.state.db.Get(prefixKey(key))
				colorString := string(colorBytes[:])
				color, _ := strconv.Atoi(colorString)
				grid[x][y] = types.DataTypesName[color]
				fmt.Println(types.DataTypesName[color])
			}
		}
	}

	fmt.Println(grid)
	return &grid
}

func validatePayload(tx []byte) (bool, string) {
	txString := string(tx[:])
	decoded, err := base64.StdEncoding.DecodeString(txString)

	if err != nil {
		return false, fmt.Sprintf("Invalid base64 encoding %s", err)
	}

	var message types.Transaction
	err = json.Unmarshal(decoded, &message)

	if err != nil {
		return false, fmt.Sprintf("Invalid json format %s", err)
	}

	if (message.X > gridsize) || (message.X < 0) {
		return false, "X coordinate is not in range."
	}
	if (message.Y > gridsize) || (message.Y < 0) {
		return false, "Y coordinate is not in range."
	}

	return true, ""
}
