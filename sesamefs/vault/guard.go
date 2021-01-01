package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
)

// InfoType refers to the type (KES, VRF, and Cert) of information required to operate
//a block producer in the Cardano network.
type InfoType string

const (
	KES  InfoType = "kes"
	VRF  InfoType = "vrf"
	Cert InfoType = "cert"
)

// Guard controls the access to the three InfoType (keys and certificate) for operating
// a stake pool that are stored in the Vault in an encrypted format.
//
// A guard allows access only once to each InfoType for the same Vault entry and PadKey.
// The client has to set the correct key using the SetKey method and then all of the
// three InfoType can be accessed exactly one time. The decrypted Vault entry as well as
// the PadKey are thrown away immediately afterwards.
type Guard interface {

	// SetKey sets the given key phrase for the guard, if the same key phrase hasn't already
	// been set. The given key phrase must not be empty. If setting the key fails. then an
	// error will be returned.
	SetKey(keyPhrase string) error

	// Access aims to decrypt the info of a given InfoType of the first entry of a Vault by
	// using the currently set key (method: SetKey). If the decryption is successful, then
	// the JSON data of InfoType is returned. Otherwise an error will be returned.
	Access(typ InfoType) ([]byte, error)
}

type vaultGuard struct {
	vault     Vault
	lock      *sync.Mutex
	accessMem AccessMemory
}

// NewVaultGuard creates a new Guard controlling the access to the three InfoType that are
// stored in the given vault.
func NewVaultGuard(myVault Vault) Guard {
	return &vaultGuard{
		vault: myVault,
		lock:  &sync.Mutex{},
	}
}

func (vg *vaultGuard) Access(typ InfoType) ([]byte, error) {
	vg.lock.Lock()
	defer vg.lock.Unlock()
	if vg.accessMem == nil {
		return nil, errors.New("key passphrase has not been set")
	}
	if vg.accessMem.HasBeenAccessed(typ) {
		return nil, errors.New("file has already been accessed")
	}
	pads, err := vg.vault.Read()
	if err != nil {
		return nil, err
	}
	if len(pads) == 0 {
		return nil, errors.New("no pads in the vault")
	}
	pk, err := ParseKeyPhrase(vg.accessMem.Key())
	if err != nil {
		return nil, err
	}
	pkg, err := pk.Decrypt(pads[0])
	if err != nil {
		return nil, err
	}
	// successful decryption
	vg.accessMem.SetAccessFlag(typ)
	if vg.accessMem.AllAccessed() {
		vg.accessMem = nil
		err = deleteVaultEntry(vg.vault)
		if err != nil {
			log.Errorf("couldn't delete the vault entry: %s", err.Error())
		}
	}
	return readInfo(typ, pkg)
}

func readInfo(typ InfoType, pkg *Package) ([]byte, error) {
	var keyJSON KeyJSON
	switch typ {
	case KES:
		keyJSON = pkg.KES
	case VRF:
		keyJSON = pkg.VRF
	case Cert:
		keyJSON = pkg.Cert
	default:
		panic(fmt.Sprintf("unknown info type '%v'", typ))
	}
	data, err := json.Marshal(keyJSON)
	return data, err
}

func deleteVaultEntry(vault Vault) error {
	pads, err := vault.Read()
	if err != nil {
		return err
	}
	if len(pads) > 0 {
		err = vault.Write(pads[1:])
		return err
	}
	return nil
}

func (vg *vaultGuard) SetKey(keyPhrase string) error {
	vg.lock.Lock()
	defer vg.lock.Unlock()
	if keyPhrase == "" {
		return errors.New("the given key phrase must not be empty")
	}
	if len(keyPhrase) != 72 || keyPhrase[:4] != "100X" {
		return errors.New("the format of the key is invalid")
	}
	if vg.accessMem != nil {
		if keyPhrase == vg.accessMem.Key() {
			return errors.New("same key phrase has already been set")
		}
		if vg.accessMem.PartiallyAccessed() {
			err := deleteVaultEntry(vg.vault)
			if err != nil {
				log.Errorf("could not delete the first entry in vault: %s", err.Error())
			}
		}
	}
	vg.accessMem = NewAccessMemory(keyPhrase)
	return nil
}

// AccessMemory is an object that keeps a memory on which InfoType has already
// been accessed for a certain key.
type AccessMemory interface {

	// Key returns the stored key phrase.
	Key() string

	// SetAccessFlag sets the access flag of the given InfoType to true.
	SetAccessFlag(typ InfoType)

	// HasBeenAccessed returns true, if the access flag has previously been set for
	// the given InfoType.
	HasBeenAccessed(typ InfoType) bool

	// PartiallyAccessed returns true, if not all InfoType has been accessed, but at
	// least one has.
	PartiallyAccessed() bool

	// AllAccessed returns true, if all three certificate InfoType have been accessed
	// at least once.
	AllAccessed() bool
}

type accessMemory struct {
	key  string
	flag map[string]bool
}

// NewAccessMemory creates a simple AccessMemory using a map.
func NewAccessMemory(key string) AccessMemory {
	flagMap := map[string]bool{
		"kes":  false,
		"vrf":  false,
		"cert": false,
	}
	mem := accessMemory{
		key:  key,
		flag: flagMap,
	}
	return &mem
}

func (mem *accessMemory) Key() string {
	return mem.key
}

func (mem *accessMemory) SetAccessFlag(typ InfoType) {
	mem.flag[string(typ)] = true
}

func (mem *accessMemory) HasBeenAccessed(typ InfoType) bool {
	val, found := mem.flag[string(typ)]
	if !found {
		panic(fmt.Sprintf("unknown info type '%v'", typ))
	}
	return val
}

func (mem *accessMemory) PartiallyAccessed() bool {
	kes, _ := mem.flag[string(KES)]
	vrf, _ := mem.flag[string(VRF)]
	cert, _ := mem.flag[string(Cert)]
	return !mem.AllAccessed() || (kes || vrf || cert)
}

func (mem *accessMemory) AllAccessed() bool {
	kes, _ := mem.flag[string(KES)]
	vrf, _ := mem.flag[string(VRF)]
	cert, _ := mem.flag[string(Cert)]
	return kes && vrf && cert
}
