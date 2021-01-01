package main

import (
	"fmt"
	"github.com/outofbits/sesameFS/sesamefs-client/cmd"
	"os"
)

// main is the entry point of this sesamefs application.
func main() {
	err := cmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}