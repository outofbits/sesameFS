package sesame

import (
	"encoding/json"
	"io/ioutil"
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

// Bytes serializes the Package into a JSON string and returns the corresponding
// byte array. An error will be returned, if the serialization fails.
func (pkg *Package) Bytes() ([]byte, error) {
	data, err := json.Marshal(pkg)
	return data, err
}

// read reads the certificate details and packages them into a Package.
// An error will be returned if the files could not be read in.
func read(kesFilePath string, vrfFilePath string, certFilePath string) (*Package, error) {
	// KES
	kesData, err := ioutil.ReadFile(kesFilePath)
	if err != nil {
		return nil, err
	}
	var kesKey KeyJSON
	err = json.Unmarshal(kesData, &kesKey)
	if err != nil {
		return nil, err
	}
	// VRF
	vrfData, err := ioutil.ReadFile(vrfFilePath)
	if err != nil {
		return nil, err
	}
	var vrfKey KeyJSON
	err = json.Unmarshal(vrfData, &vrfKey)
	if err != nil {
		return nil, err
	}
	// Node Certificate
	nodeDate, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		return nil, err
	}
	var nodeCert KeyJSON
	err = json.Unmarshal(nodeDate, &nodeCert)
	if err != nil {
		return nil, err
	}
	pkg := Package{
		KES:  kesKey,
		VRF:  vrfKey,
		Cert: nodeCert,
	}
	return &pkg, nil
}
