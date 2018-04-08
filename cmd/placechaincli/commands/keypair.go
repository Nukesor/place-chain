package commands

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	keys "github.com/tendermint/go-crypto/keys"
	"github.com/tendermint/go-crypto/keys/cryptostore"
	"github.com/tendermint/go-crypto/keys/storage/filestorage"
	"github.com/tendermint/tmlibs/cli"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "keypair",
	Short: "Create a new public/private key pair",
	RunE:  newPassword,
}

func init() {
	RootCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("type", "t", "ed25519", "Type of key (ed25519|secp256k1)")
}

func newPassword(cmd *cobra.Command, args []string) error {
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide a name for the key")
	}
	name := args[0]
	algo := viper.GetString("type")

	pass, err := getCheckPassword("Enter a passphrase:", "Repeat the passphrase:")
	if err != nil {
		return err
	}

	info, err := GetKeyManager().Create(name, pass, algo)
	if err == nil {
		printInfo(info)
	}
	return err
}

const KeySubdir = "keys"

var (
	manager keys.Manager
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "keys",
	Short: "Key manager for tendermint clients",
	Long: `Keys allows you to manage your local keystore for tendermint.

These keys may be in any format supported by go-crypto and can be
used by light-clients, full nodes, or any other application that
needs to sign with a private key.`,
}

// GetKeyManager initializes a key manager based on the configuration
func GetKeyManager() keys.Manager {
	if manager == nil {
		// store the keys directory
		rootDir := viper.GetString(cli.HomeFlag)
		keyDir := filepath.Join(rootDir, KeySubdir)
		// and construct the key manager
		manager = cryptostore.New(
			cryptostore.SecretBox,
			filestorage.New(keyDir),
		)
	}
	return manager
}
