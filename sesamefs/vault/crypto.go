package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
)

// Package wraps all the details required for operating stake pool.
type Package struct {
	KES  KeyJSON `json:"kes"`
	VRF  KeyJSON `json:"vrf"`
	Cert KeyJSON `json:"cert"`
}

// KeyJSON is the human-friendly format in which keys and certs are
// communicated in the Cardano network.
type KeyJSON struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	CBOR        string `json:"cborHex"`
}

// PadKey are the key details needed to decrypt pads in the Vault.
type PadKey struct {
	Key []byte
	IV  []byte
}

// ParseKeyPhrase parses the given key phrase and returns a PadKey, if
// the phrase was given in a valid format. Otherwise an error will be
// returned.
func ParseKeyPhrase(phrase string) (*PadKey, error) {
	if len(phrase) < 72 {
		return nil, errors.New("invalid key phrase passed")
	}
	phrase = phrase[4:] // ignore prefix
	iv, err := base64.StdEncoding.DecodeString(phrase[:24])
	if err != nil {
		return nil, err
	}
	key, err := base64.StdEncoding.DecodeString(phrase[24:])
	if err != nil {
		return nil, err
	}
	pk := PadKey{
		Key: key,
		IV:  iv,
	}
	return &pk, nil
}

// Decrypt uses this PadKey to decrypt the given encrypted entry. The entry is
// decoded with Base64 into a byte array and then decrypted further using the
// IV and AES key of this PadKey. If the decryption has been successful a Package
// with the key details will be returned.
//
// Should the entry not be in a valid Base64 format, then an error will be returned.
// Moreover, the decoded byte array must be a multiple of aes.BlockSize, otherwise an
// error will be returned.
//
// The CBC block cipher using AES decrypts the ciphertext with any key and IV. However,
// in case of a wrong key, we are going to get gibberish that is not going to serialize
// to valid JSON and an error will be returned. There is no distinction made between
// wrong key and corrupted vault entry.
func (pk *PadKey) Decrypt(entry string) (*Package, error) {
	plainData, err := base64.StdEncoding.DecodeString(entry)
	if err != nil {
		return nil, err
	}
	if (len(plainData) % aes.BlockSize) != 0 {
		return nil, errors.New("the vault entry was corrupted")
	}
	block, err := aes.NewCipher(pk.Key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, pk.IV)
	plainText := make([]byte, len(plainData))
	mode.CryptBlocks(plainText, plainData)
	var pkg Package
	err = json.Unmarshal(plainText, &pkg)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}
