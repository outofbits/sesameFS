package cmd

import (
	"errors"
	"fmt"
	"github.com/outofbits/sesameFS/sesamefs-client/sesame"
	"github.com/spf13/cobra"
	"os"
)

var (
	sendKeyCmd = &cobra.Command{
		Use:   "send-key <key-phrase>",
		Short: "Send given key phrase to sesameFS instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you must specify exactly one key phrase to send")
			}
			client, err := sesame.NewClient(address)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error: could not reach API endpoint: %s\n", err.Error())
				os.Exit(1)
			}
			err = client.SendKeyPhrase(args[0])
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error: could not send key: %s\n", err.Error())
				os.Exit(1)
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(sendKeyCmd)
}
