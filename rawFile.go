package kxcfuse

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"syscall"
)

type rawFile struct {
	parent  *FuseMount
	name    string
	inode   uint64
	content []byte
}

func (f *rawFile) Attr(ctx context.Context, attr *fuse.Attr) error {

	log.Info().Msg(fmt.Sprintf("[FUSE]>> Reading raw-file attributes of '%s' (@ %d)", f.name, f.inode))

	attr.Inode = 2
	attr.Mode = 0o444
	attr.Size = uint64(len(f.content))
	return nil
}

func (f *rawFile) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {

	log.Info().Msg(fmt.Sprintf("[FUSE]>> Opening raw-file '%s' (@ %d) by [uid:%d | gid:%d | pid:%d]", f.name, f.inode, req.Uid, req.Gid, req.Pid))

	if !req.Flags.IsReadOnly() {
		return nil, fuse.Errno(syscall.EACCES)
	}
	resp.Flags |= fuse.OpenKeepCache
	return f, nil
}

func (f *rawFile) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {

	log.Info().Msg(fmt.Sprintf("[FUSE]>> Reading raw-file '%s' (@ %d) by [uid:%d | gid:%d | pid:%d]", f.name, f.inode, req.Uid, req.Gid, req.Pid))

	fuseutil.HandleRead(req, resp, f.content)
	return nil
}

func (f *rawFile) Dirent() fuse.Dirent {
	return fuse.Dirent{
		Inode: f.inode,
		Type:  fuse.DT_File,
		Name:  f.name,
	}
}
