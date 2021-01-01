package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of sesameFS",
	Long:  `All software has versions. This is sesameFS's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sesamefs %s\n", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}