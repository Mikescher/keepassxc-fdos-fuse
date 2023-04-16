package kxcfuse

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"context"
	"fmt"
	"github.com/godbus/dbus/v5"
	keyring "github.com/ppacher/go-dbus-keyring"
	"github.com/rs/zerolog/log"
	"syscall"
)

type specFile struct {
	parent     *FuseMount
	name       string
	inode      uint64
	ssIdentKey string
	ssIdentVal string
	ssAttr     string
}

func (f *specFile) Attr(ctx context.Context, attr *fuse.Attr) error {

	log.Info().Msg(fmt.Sprintf("[FUSE]>> Reading spec-file attributes of '%s' (@ %d)", f.name, f.inode))

	attr.Inode = 2
	attr.Mode = 0o444
	attr.Size = uint64(f.parent.maxFSize)
	return nil
}

func (f *specFile) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {

	log.Info().Msg(fmt.Sprintf("[FUSE]>> Opening spec-file '%s' (@ %d) by [uid:%d | gid:%d | pid:%d]", f.name, f.inode, req.Uid, req.Gid, req.Pid))

	if !req.Flags.IsReadOnly() {
		return nil, fuse.Errno(syscall.EACCES)
	}
	resp.Flags |= fuse.OpenKeepCache
	return f, nil
}

func (f *specFile) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {

	log.Info().Msg(fmt.Sprintf("[FUSE]>> Reading spec-file '%s' (@ %d) by [uid:%d | gid:%d | pid:%d]", f.name, f.inode, req.Uid, req.Gid, req.Pid))

	bus, err := dbus.SessionBus()
	if err != nil {
		log.Err(err).Msg("failed to open dbus")
		return fuse.Errno(syscall.EIO)
	}

	secrets, err := keyring.GetSecretService(bus)
	if err != nil {
		log.Err(err).Msg("failed to get secret-service")
		return fuse.Errno(syscall.EIO)
	}

	collection, err := secrets.GetDefaultCollection()
	if err != nil {
		log.Err(err).Msg("failed to get secret-service default collection")
		return fuse.Errno(syscall.EIO)
	}

	items, err := collection.SearchItems(map[string]string{f.ssIdentKey: f.ssIdentVal})
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("failed to get item '%s :: %s'", f.ssIdentKey, f.ssIdentVal))
		return fuse.Errno(syscall.EIO)
	}

	if len(items) == 0 {
		log.Err(err).Msg(fmt.Sprintf("failed to get item '%s :: %s' (not found in collection)", f.ssIdentKey, f.ssIdentVal))
		return fuse.Errno(syscall.EIO)
	}
	if len(items) > 1 {
		log.Err(err).Msg(fmt.Sprintf("failed to get item '%s :: %s' (found multiple matches collection)", f.ssIdentKey, f.ssIdentVal))
		return fuse.Errno(syscall.EIO)
	}

	item := items[0]

	ok, err := item.Unlock()
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("failed to unlock item '%s :: %s'", f.ssIdentKey, f.ssIdentVal))
		return fuse.Errno(syscall.EIO)
	}
	if !ok {
		log.Error().Msg(fmt.Sprintf("failed to unlock item '%s :: %s'", f.ssIdentKey, f.ssIdentVal))
		return fuse.Errno(syscall.EIO)
	}

	attr, err := item.GetAttributes()
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("failed to get attributes on item '%s :: %s'", f.ssIdentKey, f.ssIdentVal))
		return fuse.Errno(syscall.EIO)
	}

	if f.ssAttr == "@password" {

		session, err := secrets.OpenSession()
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("failed to open dbus-secretservice session"))
			return fuse.Errno(syscall.EIO)
		}

		sec, err := item.GetSecret(session.Path())
		if err != nil {
			return fuse.Errno(syscall.EIO)
		}

		if len(sec.Value) > f.parent.maxFSize {
			log.Err(err).Msg(fmt.Sprintf("data of item '%s :: %s' is too big (%d)", f.ssIdentKey, f.ssIdentVal, len(sec.Value)))
			return fuse.Errno(syscall.EIO)
		}

		fuseutil.HandleRead(req, resp, sec.Value)
		return nil

	} else {

		if v, ok := attr[f.ssAttr]; ok {
			if len(v) > f.parent.maxFSize {
				log.Err(err).Msg(fmt.Sprintf("data of item '%s :: %s' is too big (%d)", f.ssIdentKey, f.ssIdentVal, len(v)))
				return fuse.Errno(syscall.EIO)
			}

			fuseutil.HandleRead(req, resp, []byte(v))
			return nil
		} else {
			log.Err(err).Msg(fmt.Sprintf("failed to get attribute '%s' on item '%s :: %s' (attr does not exist)", f.ssAttr, f.ssIdentKey, f.ssIdentVal))
			return fuse.Errno(syscall.EIO)
		}
	}

}

func (f *specFile) Dirent() fuse.Dirent {
	return fuse.Dirent{
		Inode: f.inode,
		Type:  fuse.DT_File,
		Name:  f.name,
	}
}
