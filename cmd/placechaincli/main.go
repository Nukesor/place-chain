package main

import (
	"os"

	"github.com/tendermint/go-crypto/cmd"
	"github.com/tendermint/tmlibs/cli"
)

func main() {
	root := cli.PrepareMainCmd(cmd.RootCmd, "TM", os.ExpandEnv("$HOME/.tlc"))
	root.Execute()
}
