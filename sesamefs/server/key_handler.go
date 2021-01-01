package server

import (
	"encoding/json"
	"fmt"
	"github.com/outofbits/sesameFS/sesamefs/vault"
	"io/ioutil"
	"net/http"
)

func handleKeySend(guard vault.Guard, w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		data, err := ioutil.ReadAll(req.Body)
		if err == nil {
			var passphrase string
			err = json.Unmarshal(data, &passphrase)
			if err == nil {
				err = guard.SetKey(passphrase)
				if err == nil {
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "text/plain")
					_, _ = w.Write([]byte("ok"))
				} else {
					http.Error(w, fmt.Sprintf("could not write the key phrase to the guard: %s", err.Error()),
						http.StatusBadRequest)
				}
			} else {
				http.Error(w, fmt.Sprintf("could not deserialize the message body: %s", err.Error()),
					http.StatusBadRequest)
			}
		} else {
			http.Error(w, fmt.Sprintf("could not read the message body: %s", err.Error()),
				http.StatusBadRequest)
		}
	} else {
		http.Error(w, "you must use the HTTP POST method", http.StatusBadRequest)
	}
}
