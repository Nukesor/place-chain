package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/abci/server"
	cmn "github.com/tendermint/tmlibs/common"

	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"

	"../../../app"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start place chain node",
	RunE:  startCmd,
}

//nolint
const (
	FlagAddress           = "address"
	FlagWithoutTendermint = "without-tendermint"
)

func init() {
	flags := StartCmd.Flags()
	flags.String(FlagAddress, "tcp://0.0.0.0:46658", "Listen address")
	flags.Bool(FlagWithoutTendermint, false, "Only run place-chain abci app, assume external tendermint process")
	// add all standard 'tendermint node' flags
	tcmd.AddNodeFlags(StartCmd)
}

func startCmd(cmd *cobra.Command, args []string) error {
	placeChainApp := app.NewKVStoreApplication()
	startApp(placeChainApp)
	return nil
}

func startApp(placeChainApp *app.KVStoreApplication) error {
	// Start the ABCI listener
	addr := viper.GetString(FlagAddress)
	go (&app.WebServer{}).LaunchHTTP()
	svr, err := server.NewServer(addr, "socket", placeChainApp)
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