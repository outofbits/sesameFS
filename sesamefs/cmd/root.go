package cmd

import "github.com/spf13/cobra"

const VERSION string = "1.0.0-alpha"

var (
	rootCmd = &cobra.Command{
		Use:   "sesamefs",
		Short: "A filesystem to protect your SP certificate.",
		Long: `The sesameFS filesystem can be mounted to a special location
on your host filesystem. It manages a single root directory 
and an access control (over HTTP+JSON API) to the three files 
kes.skey, vrf.skey & node.cert, which are required to operate
a block producer in the Cardano network.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
