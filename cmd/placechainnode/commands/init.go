package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"

	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

// InitCmd initialises a fresh Tendermint Core + Place-Chain instance.
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize full TM place-chain node",
	RunE:  initFiles,
}

const (
	FlagChainId = "chain-id"
)

func init() {
	flags := InitCmd.Flags()
	flags.String(FlagChainId, cmn.Fmt("test-chain-%v", cmn.RandStr(6)), "Chain ID")
	tcmd.AddNodeFlags(InitCmd)
}

func initFiles(cmd *cobra.Command, args []string) error {
	// private validator
	config, err := tcmd.ParseConfig()
	if err != nil {
		return err
	}
	privValFile := config.PrivValidatorFile()
	var privValidator *types.PrivValidatorFS
	if cmn.FileExists(privValFile) {
		privValidator = types.LoadPrivValidatorFS(privValFile)
		logger.Info("Found private validator", "path", privValFile)
	} else {
		privValidator = types.GenPrivValidatorFS(privValFile)
		privValidator.Save()
		logger.Info("Generated private validator", "path", privValFile)
	}

	// genesis file
	genFile := config.GenesisFile()
	if cmn.FileExists(genFile) {
		logger.Info("Found genesis file", "path", genFile)
	} else {
		genDoc := types.GenesisDoc{
			ChainID: viper.GetString(FlagChainId),
		}
		genDoc.Validators = []types.GenesisValidator{{
			PubKey: privValidator.GetPubKey(),
			Power:  10,
		}}

		if err := genDoc.SaveAs(genFile); err != nil {
			return err
		}
		logger.Info("Generated genesis file", "path", genFile)
	}
	return nil
}
