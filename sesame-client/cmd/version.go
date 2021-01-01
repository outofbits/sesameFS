package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of sesame-sesame",
	Long:  `All software has versions. This is sesame-client's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sesame-sesame %s\n", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}