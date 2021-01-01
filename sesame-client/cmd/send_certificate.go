package cmd

import (
	"errors"
	"fmt"
	"github.com/outofbits/sesameFS/sesamefs-client/sesame"
	"github.com/spf13/cobra"
	"os"
)

var (
	number             uint8
	kesFilePath        string
	vrfFilePath        string
	certFilePath       string
	sendCertificateCmd = &cobra.Command{
		Use:   "send-certificate",
		Short: "Send encrypted certificate to sesameFS instance",
		Long: `This command parses the given certificate details and creates n 
one-time-pads for your certificate. The n differently encrypted
versions of your certificate are send to the given sesameFS instance. 
Please don't send it over the public internet. You could use a SSH 
tunnel or Wireguard.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if number == 0 {
				return errors.New("the given number must not be 0")
			}
			_, err := os.Stat(kesFilePath)
			if err != nil {
				return err
			}
			_, err = os.Stat(vrfFilePath)
			if err != nil {
				return err
			}
			_, err = os.Stat(certFilePath)
			if err != nil {
				return err
			}
			client, err := sesame.NewClient(address)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error: could not reach API endpoint: %s", err.Error())
				os.Exit(1)
			}
			err = client.SendCertificate(number, kesFilePath, vrfFilePath, certFilePath)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error: could not send certificates: %s", err.Error())
				os.Exit(1)
			}
			return nil
		},
	}
)

func init() {
	sendCertificateCmd.Flags().Uint8Var(&number, "n", 1, "the number of one-time-pads to generate")
	sendCertificateCmd.PersistentFlags().StringVar(&kesFilePath, "kes-signing-key", "", "path to the kes signing key")
	sendCertificateCmd.PersistentFlags().StringVar(&vrfFilePath, "vrf-signing-key", "", "path to the vrf signing key")
	sendCertificateCmd.PersistentFlags().StringVar(&certFilePath, "node-cert", "", "path to the node certificate")
	rootCmd.AddCommand(sendCertificateCmd)
}
