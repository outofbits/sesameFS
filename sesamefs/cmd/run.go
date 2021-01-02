package cmd

import (
	"errors"
	"fmt"
	"github.com/outofbits/sesameFS/sesamefs/fs"
	"github.com/outofbits/sesameFS/sesamefs/server"
	"github.com/outofbits/sesameFS/sesamefs/vault"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"
)

const (
	DefaultUID uint32 = 1000
	DefaultGID uint32 = 1000
)

var (
	mountPoint        string
	uid               uint32
	gid               uint32
	dataDirectoryPath string
	serverAddress     string
	allowOther        bool
	runCommand        = &cobra.Command{
		Use:   "run",
		Short: "run the sesameFS filesystem",
		Long:  `Runs a sesameFS filesystem at the given mountpoint
and the HTTP+JSON API at the given address. At the moment, the
encrypted certificate details can either be stored in an in-memory
vault or in an file vault on the host filesystem. In the former mode
the encrypted certificate details will be lost after the sesameFS
filesystem is closed down.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// check the mount point argument
			if mountPoint == "" {
				return errors.New("a mountpoint must be given")
			}
			mountPointInfo, err := os.Stat(mountPoint)
			if err != nil {
				if os.IsNotExist(err) {
					err = os.MkdirAll(mountPoint, 0770)
					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "error: could not create directory '%s':\n", mountPoint)
						os.Exit(1)
					}
				} else {
					_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
					os.Exit(1)
				}
			}
			if mountPointInfo != nil && !mountPointInfo.IsDir() {
				_, _ = fmt.Fprintf(os.Stderr, "error: the mountpoint '%s' must be a directory, but wasn't. ", mountPoint)
				os.Exit(1)
			}
			// check the data directory argument
			if dataDirectoryPath != "" {
				dataDirectoryInfo, err := os.Stat(dataDirectoryPath)
				if err != nil {
					if !os.IsNotExist(err) {
						_, _ = fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
						os.Exit(1)
					}
					err = os.MkdirAll(dataDirectoryPath, 0770)
					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "error: could not create directory '%s':\n", dataDirectoryPath)
						os.Exit(1)
					}
				}
				if dataDirectoryInfo != nil && !dataDirectoryInfo.IsDir() {
					_, _ = fmt.Fprintf(os.Stderr, "error: the data directory '%s' must be a directory, but wasn't. ", dataDirectoryPath)
					os.Exit(1)
				}
			}
			// run the filesystem
			fsConfig := fs.FilesystemConfig{
				UID:        uid,
				GID:        gid,
				AllowOther: allowOther,
			}
			err = run(mountPoint, fsConfig, dataDirectoryPath, serverAddress)
			if err != nil {
				_, _ = fmt.Printf("error: %s\n", err.Error())
				os.Exit(1)
			}
			return nil
		},
	}
)

func init() {
	runCommand.Flags().StringVar(&mountPoint, "mountpoint", "", "location on the host filesystem where this sesamefs shall be mounted")
	_ = runCommand.MarkFlagRequired("mountpoint")
	runCommand.Flags().StringVar(&dataDirectoryPath, "data-dir", "", "optional directory on the host filesystem where to store encrypted certificates (default: in-memory vault)")
	runCommand.Flags().StringVar(&serverAddress, "server-address", "localhost:13456", "address on which the HTTP+JSON API server shall listen for requests")
	// optional uid, gid and file mode
	defaultUID, defaultGID := getUserDetails()
	runCommand.Flags().Uint32Var(&uid, "uid", defaultUID, "the user id of the certificate files in sesameFS")
	runCommand.Flags().Uint32Var(&gid, "gid", defaultGID, "the group id  of the certificate files in sesameFS")
	runCommand.Flags().BoolVar(&allowOther, "allow-other", false, "whether others shall be allowed to mount this filesystem")
	// add command to root
	rootCmd.AddCommand(runCommand)
}

// getUserDetails aims to get details about the user that is running this process. This method returns the user ID
// and group ID of the fetched user details. If it fails to fetch details, then the default UID and GID are
// returned.
func getUserDetails() (uint32, uint32) {
	processUser, err := user.Current()
	if err == nil {
		var uid uint32
		parsedUID, err := strconv.ParseUint(processUser.Uid, 10, 32)
		if err != nil {
			log.Warnf("could not detect the UID of the user running this process")
			uid = DefaultUID
		} else {
			uid = uint32(parsedUID)
		}
		var gid uint32
		parsedGID, err := strconv.ParseUint(processUser.Gid, 10, 32)
		if err != nil {
			log.Warnf("could not detect the GID of the user running this process")
			gid = DefaultGID
		} else {
			gid = uint32(parsedGID)
		}
		return uid, gid
	}
	return DefaultUID, DefaultGID
}

func run(mountPoint string, config fs.FilesystemConfig, dataDirectoryPath string, serverAddress string) error {
	// create or open vault
	var err error
	var myVault vault.Vault
	if dataDirectoryPath != "" {
		myVault, err = vault.NewFileVault(dataDirectoryPath)
		if err != nil {
			return err
		}
		log.Infof("opened a vault at '%s'", dataDirectoryPath)
	} else {
		myVault = vault.NewInMemoryVault()
		log.Info("opened an im-memory vault")
	}
	guard := vault.NewVaultGuard(myVault)
	defer myVault.Close()

	errorListener := make(chan error)
	// mount sesameFS
	sesameFS := &fs.SesameFS{
		Guard:         guard,
		Config:        config,
		ErrorListener: errorListener,
	}
	err = sesameFS.Mount(mountPoint)
	if err != nil {
		return err
	}
	defer sesameFS.Umount(mountPoint)
	//start the API
	go server.Listen(serverAddress, myVault, guard, errorListener)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// handle shutdown
	select {
	case err = <-errorListener:
		return err
	case sig := <-sigs:
		log.Infof("detected '%s' signal, shutting down ...", sig.String())
		return nil
	}
}
