package server

import (
	"encoding/json"
	"fmt"
	"github.com/outofbits/sesameFS/sesamefs/vault"
	"io/ioutil"
	"net/http"
)

func handleEncryptedCertificates(vault vault.Vault, w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		data, err := ioutil.ReadAll(req.Body)
		if err == nil {
			var otpArray []string
			err = json.Unmarshal(data, &otpArray)
			if err == nil {
				err = vault.Write(otpArray)
				if err == nil {
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "text/plain")
					_, _ = w.Write([]byte("ok"))
				} else {
					http.Error(w, fmt.Sprintf("could not write the OTP array to the vault: %s", err.Error()),
						http.StatusInternalServerError)
				}
			} else {
				http.Error(w, fmt.Sprintf("could not deserialize the message body: %s", err.Error()),
					http.StatusInternalServerError)
			}
		} else {
			http.Error(w, fmt.Sprintf("could not read the message body: %s", err.Error()),
				http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "you must use the HTTP POST method", http.StatusBadRequest)
	}
}
