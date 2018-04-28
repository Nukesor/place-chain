package commands

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"

	"github.com/tendermint/go-crypto"
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
	FlagChainId = "place-chain"
)

var PubKeys = []map[string]string {
	{
		"name": "jarvis",
		"key": `{"type":"ed25519","data":"58683B7A778D165EDBFA488732FA094AC3B70A6E8FEC81FE945F16BE2BD7A69C"}`,
	},
}

func init() {
	flags := InitCmd.Flags()
	flags.String(FlagChainId, cmn.Fmt("place-chain", cmn.RandStr(6)), "Chain ID")
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

		genDoc.Validators = []types.GenesisValidator{};

		for _, pubKeyInfo := range PubKeys {
			var pubKey crypto.PubKey
			bytes := []byte(pubKeyInfo["key"])
			err = json.Unmarshal(bytes, &pubKey)
			if err != nil {
				return err
			}

			validator := types.GenesisValidator{
				PubKey: pubKey,
				Power:  10,
				Name: pubKeyInfo["name"],
			}
			genDoc.Validators = append(genDoc.Validators, validator)
		}

		if err := genDoc.SaveAs(genFile); err != nil {
			return err
		}
		logger.Info("Generated genesis file", "path", genFile)
	}

	return nil
}
