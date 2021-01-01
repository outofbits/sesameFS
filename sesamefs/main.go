package main

import (
	"fmt"
	"github.com/outofbits/sesameFS/sesamefs/cmd"
	log "github.com/sirupsen/logrus"
	"os"
)

// main is the entry point of this sesamefs application.
func main() {
	log.SetFormatter(&log.JSONFormatter{})
	err := cmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
