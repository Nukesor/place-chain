package app

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	cfg "github.com/tendermint/tendermint/config"
	"place-chain/types"

	"github.com/tendermint/abci/example/code"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	cmn "github.com/tendermint/tmlibs/common"
	tmTypes "github.com/tendermint/tendermint/types"


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

func prefixKey(key []byte) []byte {
	return append(kvPairPrefixKey, key...)
}

//---------------------------------------------------

var _ abci.Application = (*PlacechainApp)(nil)

type PlacechainApp struct {
	abci.BaseApplication

	state      types.AppState
	httpClient httpcli.HTTP
	config     *cfg.Config
}

func NewPlacechainApp(config *cfg.Config) *PlacechainApp {
	var state types.AppState
	state.Db = dbm.NewMemDB()
	httpClient := httpcli.NewHTTP("tcp://0.0.0.0:46657", "/websocket")
	return &PlacechainApp{state: state, httpClient: *httpClient, config: config}
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
		key = []byte(rt.TwitterHandle)
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

	if !app.IsTransactionValid(tx) {
		return abci.ResponseCheckTx{Code: code.CodeTypeEncodingError}
	}

	return abci.ResponseCheckTx{Code: code.CodeTypeOK}
}

func (app *PlacechainApp) IsTransactionValid(tx types.Transaction) bool {
	twitterHandle := tx.GetTwitterHandle()
	if twitterHandle == "" {
		return false
	}

	bytes, err := tx.SignedBytes()
	if err != nil {
		fmt.Println("Transaction: Could not serialize transaction bytes for verifying signature")
		return false
	}

	if pt, isPt := tx.(types.PixelTransaction); isPt {
		pubKey, err := app.GetPubKey(pt.GetTwitterHandle())
		if err != nil {
			return false
		}

		return pubKey.VerifyBytes(bytes, pt.Signature)
	} else if rt, isRt := tx.(types.RegisterTransaction); isRt {
		// TODO verify validator pubkey
		return rt.ValidatorPubKey.VerifyBytes(bytes, rt.Signature)
	}

	fmt.Printf("Transaction: Can't verify transaction %v\n", tx)
	return false
}

func (app *PlacechainApp) Commit() abci.ResponseCommit {
	fmt.Println("========================== COMMIT")
	// Using a memdb - just return the big endian size of the db
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.Size)
	app.state.AppHash = appHash
	app.state.Height += 1
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

			grid[x][y] = types.Pixel{Color: pt.Color, TwitterHandle: pt.TwitterHandle}
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

func (app *PlacechainApp) GetPubKey(twitterHandle string) (crypto.PubKey, error) {
	bytes := app.state.Db.Get(prefixKey([]byte(twitterHandle)))
	var pubKey crypto.PubKey
	err := json.Unmarshal(bytes, &pubKey)
	return pubKey, err
}

func (app *PlacechainApp) RegisterUser(rr types.RegisterRequest) error {
	privValFile := app.config.PrivValidatorFile()
	if !cmn.FileExists(privValFile) {
		return errors.New("I'm not a validator, can't create user")
	}
	if len(app.state.Db.Get(prefixKey([]byte(rr.TwitterHandle)))) > 0 {
		return errors.New("This twitterHandle is already registered")
	}
	if err := callTwitter(); err != nil {
		return err
	}
	privValidator := tmTypes.LoadPrivValidatorFS(privValFile)
	data := struct {
		TwitterHandle string
		UserPubKey    crypto.PubKey
	}{
		rr.TwitterHandle, rr.PubKey,
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	signature := privValidator.PrivKey.Sign(bytes)
	tx := rr.ToTransaction(privValidator.PubKey, signature)
	app.PublishTx(tx)
	return nil
}

func callTwitter() error {
	return nil
}
