package server

import (
	"github.com/outofbits/sesameFS/sesamefs/vault"
	"net/http"
)

// Listen starts a HTTP+JSON API server on the given address.
func Listen(address string, guard vault.Guard, errorListener chan error) {
	http.HandleFunc("/pads", func(writer http.ResponseWriter, request *http.Request) {
		handleEncryptedCertificates(guard, writer, request)
	})
	http.HandleFunc("/key", func(writer http.ResponseWriter, request *http.Request) {
		handleKeySend(guard, writer, request)
	})
	err := http.ListenAndServe(address, nil)
	if err != nil && errorListener != nil {
		errorListener <- err
	}
}