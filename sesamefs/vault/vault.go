package vault

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// Vault is a storage for the encrypted operational details for operating a
// block producer in the Cardano network. The operations on the vault storage
// must be thread-safe.
type Vault interface {

	//Read reads the complete pad array from the vault storage. Should the
	// read not be successful, then an error will be returned.
	Read() ([]string, error)

	// Write writes the given pad array to the vault and overrides
	// content that might have been in the vault before. Should the write
	// not be successful, then an error will be returned.
	Write(pads []string) error

	// Close frees resources that are occupied by this vault.
	Close()
}

// fileVault is an implementation of Vault that uses a file on the host filesystem as storage.
type fileVault struct {
	f    *os.File
	lock *sync.RWMutex
}

// NewFileVault creates a new vault using the given data directory on the host filesystem as
// storage. An error will be returned, if accessing the data directory fails.
func NewFileVault(dataDirectoryPath string) (Vault, error) {
	vaultFile, err := os.OpenFile(filepath.Join(dataDirectoryPath, "vault.db"),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0660)
	if err != nil {
		return nil, err
	}
	vault := fileVault{
		f:    vaultFile,
		lock: &sync.RWMutex{},
	}
	return &vault, nil
}

func (vault *fileVault) Read() ([]string, error) {
	vault.lock.RLock()
	defer vault.lock.RUnlock()
	_, err := vault.f.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(vault.f)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []string{}, nil
	}
	var otpArray []string
	err = json.Unmarshal(data, &otpArray)
	if err != nil {
		return nil, err
	}
	return otpArray, nil
}

func (vault *fileVault) Write(otpArray []string) error {
	vault.lock.Lock()
	defer vault.lock.Unlock()
	data, err := json.Marshal(otpArray)
	if err != nil {
		return err
	}
	err = vault.f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = vault.f.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = vault.f.Write(data)
	if err != nil {
		return err
	}
	err = vault.f.Sync()
	return err
}

func (vault *fileVault) Close() {
	vault.lock.Lock()
	defer vault.lock.Unlock()
	_ = vault.f.Close()
}

// inMemoryVault is an implementation of Vault that just uses a simple
// array in memory to store the vault entries.
type inMemoryVault struct {
	arr  []string
	lock *sync.RWMutex
}

// NewInMemoryVault creates a new Vault that just uses a simple array in memory
// to store the vault entries
func NewInMemoryVault() Vault {
	return &inMemoryVault{
		arr:  []string{},
		lock: &sync.RWMutex{},
	}
}

func (vault *inMemoryVault) Read() ([]string, error) {
	vault.lock.RLock()
	defer vault.lock.RUnlock()
	return vault.arr, nil
}

func (vault *inMemoryVault) Write(otpArray []string) error {
	vault.lock.Lock()
	defer vault.lock.Unlock()
	vault.arr = otpArray
	return nil
}

func (vault *inMemoryVault) Close() {
	// nothing to do
}
