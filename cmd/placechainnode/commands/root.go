package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tmlibs/cli"
	tmflags "github.com/tendermint/tmlibs/cli/flags"
	"github.com/tendermint/tmlibs/log"
)

const (
	defaultLogLevel = "error"
	FlagLogLevel    = "log_level"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
)

var RootCmd = &cobra.Command{
	Use:   "Place chain",
	Short: "Draw beautiful things using a tendermint blockchain.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		level := viper.GetString(FlagLogLevel)
		logger, err = tmflags.ParseLogLevel(level, logger, defaultLogLevel)
		if err != nil {
			return err
		}
		if viper.GetBool(cli.TraceFlag) {
			logger = log.NewTracingLogger(logger)
		}
		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().String(FlagLogLevel, defaultLogLevel, "Log level")
}
