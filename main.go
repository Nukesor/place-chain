package main

import (
	app "./app"
	"github.com/tendermint/abci/server"
	cmn "github.com/tendermint/tmlibs/common"
)

func main() {
	kvStoreApp := app.NewKVStoreApplication()
	startApp(kvStoreApp)
}

func startApp(kvStoreApp *app.KVStoreApplication) error {
	// Start the ABCI listener
	svr, err := server.NewServer("tcp://localhost:46657", "socket", kvStoreApp)
	go (&app.WebServer{}).LaunchHTTP()
	if err != nil {
		panic(err)
	}
	svr.Start()

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		svr.Stop()
	})
	return nil
}
