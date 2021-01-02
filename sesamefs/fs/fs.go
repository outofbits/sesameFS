package fs

import (
	"bazil.org/fuse"
	fuseFs "bazil.org/fuse/fs"
	"github.com/outofbits/sesameFS/sesamefs/vault"
)

type SesameFS struct {
	conn          *fuse.Conn
	Guard         vault.Guard
	Config        FilesystemConfig
	ErrorListener chan error
}

type FilesystemConfig struct {
	UID        uint32
	GID        uint32
	AllowOther bool
}

func (sfs *SesameFS) Mount(mountPoint string) error {
	fuse.ReadOnly()
	fuse.DefaultPermissions()
	if sfs.Config.AllowOther {
		fuse.AllowOther()
	}
	c, err := fuse.Mount(mountPoint, fuse.FSName("sesame"), fuse.Subtype("sesamefs"))
	if err != nil {
		return err
	}
	sfs.conn = c
	go sfs.serve()
	return nil
}

func (sfs *SesameFS) serve() {
	err := fuseFs.Serve(sfs.conn, SesameFS{Guard: sfs.Guard, Config: sfs.Config})
	if err != nil {
		sfs.ErrorListener <- err
	}
}

func (sfs *SesameFS) Umount(mountPoint string) {
	_ = fuse.Unmount(mountPoint)
	_ = sfs.conn.Close()
}

func (sfs SesameFS) Root() (fuseFs.Node, error) {
	return SesameFSDir{Guard: sfs.Guard, config: sfs.Config}, nil
}
