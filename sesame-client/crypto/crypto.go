package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	b64 "encoding/base64"
)

// PadKey is the key package for the operational certificate details.
type PadKey struct {
	Key []byte
	IV  []byte
}

func (pk PadKey) String() string {
	return "100X" + b64.StdEncoding.EncodeToString(pk.IV) + b64.StdEncoding.EncodeToString(pk.Key)
}

// GeneratePads generates for the given plain data the given number of PadKey and
// then encrypts the plain data with those keys. The PadKey and corresponding cipher
// data are returned. Should the generation fail, an error will be returned.
func GeneratePads(number int, data []byte) ([]PadKey, []string, error) {
	modulo := len(data) % aes.BlockSize
	if modulo > 0 { // todo: this padding is not secure
		paddingSize := aes.BlockSize - modulo
		paddingArray := make([]byte, paddingSize)
		var whitespace byte = ' '
		for i := 0; i < paddingSize; i++ {
			paddingArray[i] = whitespace
		}
		data = append(data, paddingArray...)
	}
	// generate keys
	pks := make([]PadKey, number)
	for i := 0; i < number; i++ {
		key := make([]byte, 32)
		_, err := rand.Read(key)
		if err != nil {
			return nil, nil, err
		}
		iv := make([]byte, aes.BlockSize)
		_, err = rand.Read(iv)
		if err != nil {
			return nil, nil, err
		}
		pk := PadKey{
			Key: key,
			IV:  iv,
		}
		pks[i] = pk
	}
	// encrypt data
	pads := make([]string, number)
	for i := 0; i < number; i++ {
		block, err := aes.NewCipher(pks[i].Key)
		if err != nil {
			return nil, nil, err
		}
		mode := cipher.NewCBCEncrypter(block, pks[i].IV)
		cipherText := make([]byte, len(data))
		mode.CryptBlocks(cipherText, data)
		pads[i] = b64.StdEncoding.EncodeToString(cipherText)
	}
	return pks, pads, nil
}
