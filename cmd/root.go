package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/subosito/gotenv"

	"github.com/Bnei-Baruch/mdb-links/common"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "mdb-links",
	Short: "Backend for links from mdb logical files to physical files service",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is config.toml)")
}

func initConfig() {
	gotenv.Load()
	common.Init()
}
