package cmd

import "github.com/spf13/cobra"

const VERSION string = "1.0.0-alpha"

var (
	address string
	rootCmd = &cobra.Command{
		Use:   "sesame-client",
		Short: "A sesame to interact with sesameFS",
		Long: `This sesame client connects to the API of a sesameFS instance and makes
it possible to execute the supported operations such as sending the
next throw-away key.`,
	}
)

func init() {
	rootCmd.Flags().StringVar(&address, "address", "localhost:13456", "address of the API endpoint of sesameFS instance")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
