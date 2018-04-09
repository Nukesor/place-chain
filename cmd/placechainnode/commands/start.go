package commands

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path"

	"github.com/tendermint/abci/server"
	cmn "github.com/tendermint/tmlibs/common"

	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	tmNode "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"

	"../../../app"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start place chain node",
	RunE:  startCmd,
}

//nolint
const (
	FlagAddress = "address"
	FullNode    = "full-node"
	LoadGenesis = "load-genesis"
)

func init() {
	flags := StartCmd.Flags()
	flags.String(FlagAddress, "tcp://0.0.0.0:46658", "Abci Server listen address")
	flags.Bool(FullNode, false, "Run full node with internal tendermint node routine")
	flags.Bool(LoadGenesis, false, "(Full node only) load the initial genesis file")
	// add all standard 'tendermint node' flags
	tcmd.AddNodeFlags(StartCmd)
}

func startCmd(cmd *cobra.Command, args []string) error {
	placechainApp := app.NewPlacechainApp()

	// launch placechain webserver
	go (&app.WebServer{placechainApp}).LaunchHTTP()

	if !viper.GetBool(FullNode) {
		logger.Info("Starting placechain abci only")
		return startAbciOnly(placechainApp)
	} else {
		logger.Info("Starting placechain with full tendermint node integrated")
		// start the app with tendermint in-process
		return startFullNode(placechainApp)
	}
}

// Start the ABCI Server that hosts our app.
// External tendermint process can connect to our app via socket
func startAbciOnly(placechainApp *app.PlacechainApp) error {
	addr := viper.GetString(FlagAddress)
	svr, err := server.NewServer(addr, "socket", placechainApp)
	if err != nil {
		return err
	}
	svr.Start()

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		svr.Stop()
	})
	return nil
}

func startFullNode(placechainApp *app.PlacechainApp) error {
	// Create & start node
	cfg, err := tcmd.ParseConfig()
	if err != nil {
		return err
	}


	privValidatorFile := cfg.PrivValidatorFile()
	privValidator := types.LoadOrGenPrivValidatorFS(privValidatorFile)
	node, err := tmNode.NewNode(cfg,
		privValidator,
		proxy.NewLocalClientCreator(placechainApp),
		tmNode.DefaultGenesisDocProviderFunc(cfg),
		tmNode.DefaultDBProvider, logger.With("module", "node"))
	if err != nil {
		return err
	}

	if viper.GetBool(LoadGenesis) {
		if err := blockchainGenesis(placechainApp); err != nil {
			return err
		}
	}
	err = node.Start()
	if err != nil {
		return err
	}

	// Trap signal, run forever.
	node.RunForever()
	return nil
}

func blockchainGenesis(placechainApp *app.PlacechainApp) error {
	// If genesis file exists, set key-value options
	usr, err := user.Current()
	if err != nil {
		return err
	}
	genesisFile := path.Join(usr.HomeDir, ".place-chain", "config", "genesis.json")
	if _, err := os.Stat(genesisFile); err == nil {
		err := placechainApp.LoadGenesis(genesisFile)
		if err != nil {
			return err
		}
	} else {
		errors.New(fmt.Sprintf("No genesis file at %s, skipping...\n", genesisFile))
	}
	return nil
}
