package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Bnei-Baruch/mdb-links/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of mdb-links",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mdb-links version %s\n", version.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
