package fs

import (
	"bazil.org/fuse"
	fuseFs "bazil.org/fuse/fs"
	"context"
	"github.com/outofbits/sesameFS/sesamefs/vault"
	"os"
	"syscall"
)

// SesameFSDir describes directories in the sesamefs filesystem. At the moment, only one
// root directory is going to exist over the life span of a mounted filesystem instance.
type SesameFSDir struct {
	Guard  vault.Guard
	config FilesystemConfig
}

var (
	directoryContent = []fuse.Dirent{
		{Inode: 2, Name: string(NodeCertFile), Type: fuse.DT_File},
		{Inode: 3, Name: string(KESSecretFile), Type: fuse.DT_File},
		{Inode: 4, Name: string(VRFSecretFile), Type: fuse.DT_File},
	}
)

// ReadDirAll lists always the same three files (NodeCertFile, KESSecretFile, VRFSecretFile) as
// content of the root directory. This implementation only works, if it is assumed that only
// one root directory exists over the complete life span of a filesystem instance.
func (SesameFSDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return directoryContent, nil
}

// Attr lists the attributes of the root directory. It has the mode 770, and thus allows read
// and write access for the owner and group members, but not for others. The owner and group is
// specified with the construction of the filesystem.
func (sf SesameFSDir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0o775
	a.Uid = sf.config.UID
	a.Gid = sf.config.GID
	return nil
}

// Lookup checks whether file name that was asked for by the caller is matching one of the three
// constant file names (NodeCertFile, KESSecretFile, VRFSecretFile) of this filesystem.
//
// If no match can be found, then an ENOENT error code "No such file or directory" will be returned.
// If the name is matching, then the corresponding SesameFSFile is constructed and returned to
// the caller.
func (sf SesameFSDir) Lookup(ctx context.Context, name string) (fuseFs.Node, error) {
	switch FileType(name) {
	case NodeCertFile:
		nodeFile := SesameFSFile{
			name:   NodeCertFile,
			guard: sf.Guard,
			config: sf.config,
		}
		return nodeFile, nil
	case KESSecretFile:
		kesFile := SesameFSFile{
			name:   KESSecretFile,
			guard: sf.Guard,
			config: sf.config,
		}
		return kesFile, nil
	case VRFSecretFile:
		vrfFile := SesameFSFile{
			name:   VRFSecretFile,
			guard: sf.Guard,
			config: sf.config,
		}
		return vrfFile, nil
	default:
		return nil, syscall.ENOENT
	}
}
