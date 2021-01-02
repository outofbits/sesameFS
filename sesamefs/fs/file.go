package fs

import (
	"bazil.org/fuse"
	"context"
	"github.com/outofbits/sesameFS/sesamefs/vault"
	log "github.com/sirupsen/logrus"
	"syscall"
)

type FileType string

const (
	NodeCertFile  FileType = "node.cert"
	KESSecretFile FileType = "kes.skey"
	VRFSecretFile FileType = "vrf.skey"
)

const fileSize = 2048

type SesameFSFile struct {
	name   FileType
	guard  vault.Guard
	config FilesystemConfig
}

func (f SesameFSFile) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Infof("details about '%v' file requested", f.name)
	a.Inode = 2
	a.Uid = f.config.UID
	a.Gid = f.config.GID
	a.Mode = 0o400
	a.Size = fileSize
	return nil
}

func getPadding(length int) []byte {
	if length == 0 {
		return []byte{}
	}
	var newLine byte = '\n'
	if length == 1 {
		return []byte{newLine}
	}
	var whitespace byte = ' '
	padding := make([]byte, length)
	for i := 0; i < length-1; i++ {
		padding[i] = whitespace
	}
	padding[length-1] = newLine
	return padding
}

func (f SesameFSFile) ReadAll(ctx context.Context) ([]byte, error) {
	log.Infof("read of '%v' file requested", f.name)
	var data []byte
	var err error
	switch f.name {
	case KESSecretFile:
		data, err = f.guard.Access(vault.KES)
	case VRFSecretFile:
		data, err = f.guard.Access(vault.VRF)
	case NodeCertFile:
		data, err = f.guard.Access(vault.Cert)
	default:
		return nil, syscall.EACCES
	}
	if err != nil {
		log.Warnf("could not read '%v': %s", f.name, err.Error())
		return nil, syscall.EACCES
	}
	if len(data) > fileSize {
		log.Errorf("the operational certificate data is unusually large")
		return nil, syscall.EACCES
	}
	msg := append(data, getPadding(fileSize-len(data))...)
	return msg, nil
}
