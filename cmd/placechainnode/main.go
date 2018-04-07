package main

import (
	commands "./commands"
	"github.com/tendermint/tmlibs/cli"
	"os"
)

func main() {
	rt := commands.RootCmd

	rt.AddCommand(
		commands.StartCmd,
	)

	cmd := cli.PrepareMainCmd(rt, "BC", os.ExpandEnv("$HOME/.place-chain"))
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
