package app

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"

	abcicli "github.com/tendermint/abci/client"
	"github.com/tendermint/abci/example/code"
	"github.com/tendermint/abci/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
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

type Message struct {
	x      int
	y      int
	color  int
	nonnce string
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

var _ types.Application = (*KVStoreApplication)(nil)

type KVStoreApplication struct {
	types.BaseApplication

	state  State
	client abcicli.Client
}

func NewKVStoreApplication() *KVStoreApplication {
	state := loadState(dbm.NewMemDB())
	client, err := abcicli.NewClient("tcp://0.0.0.0:46658", "socket", true)
	if err != nil {
		panic(err)
	}
	return &KVStoreApplication{state: state, client: client}
}

func (app *KVStoreApplication) SetPixel(x int, y int) {
	app.client.DeliverTxSync([]byte("LELMAO"))
}

func (app *KVStoreApplication) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{Data: fmt.Sprintf("{\"size\":%v}", app.state.Size)}
}

// tx is either "key=value" or just arbitrary bytes
func (app *KVStoreApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {
	fmt.Println("========================== DELIVER TX")
	var message Message
	json.Unmarshal(tx, &message)

	keyString := fmt.Sprintf("%d,%d", message.x, message.y)
	key := []byte(keyString)

	app.state.db.Set(prefixKey(key), []byte(strconv.Itoa(message.color)))
	app.state.Size += 1

	tags := []cmn.KVPair{
		{[]byte("app.creator"), []byte("jae")},
		{[]byte("app.key"), key},
	}
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Tags: tags}
}

func (app *KVStoreApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	fmt.Println("========================== CHECK TX")
	valid, message := validatePayload(tx)
	if !valid {
		fmt.Println("========================== INVALID TX")
		fmt.Println(message)
		return types.ResponseCheckTx{Code: code.CodeTypeEncodingError}
	}
	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

func (app *KVStoreApplication) Commit() types.ResponseCommit {
	fmt.Println("========================== COMMIT")
	// Using a memdb - just return the big endian size of the db
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.Size)
	app.state.AppHash = appHash
	app.state.Height += 1
	saveState(app.state)
	return types.ResponseCommit{Data: appHash}
}

func (app *KVStoreApplication) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
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

func validatePayload(tx []byte) (bool, string) {
	tx_string := string(tx[:])
	decoded, err := base64.StdEncoding.DecodeString(tx_string)

	if err != nil {
		return false, fmt.Sprintf("Invalid base64 encoding %s", err)
	}

	var message Message
	err = json.Unmarshal(decoded, &message)

	if err != nil {
		return false, fmt.Sprintf("Invalid json format %s", err)
	}

	if (message.x > gridsize) || (message.x < 0) {
		return false, "X coordinate is not in range."
	}
	if (message.y > gridsize) || (message.y < 0) {
		return false, "Y coordinate is not in range."
	}

	return true, ""
}
