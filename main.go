package main

import (
    app "./app"
    cmn "github.com/tendermint/tmlibs/common"
    "github.com/tendermint/abci/server"

)
func main() {
    kvStoreApp := app.NewKVStoreApplication()
    startApp(kvStoreApp)
}

func startApp(kvStoreApp *app.KVStoreApplication) error {
    // Start the ABCI listener
    svr, err := server.NewServer("tcp://localhost:46657", "socket", kvStoreApp)
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