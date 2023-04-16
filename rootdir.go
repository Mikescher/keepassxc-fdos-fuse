package kxcfuse

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"os"
	"syscall"
)

type RootDir struct {
	parent *FuseMount
}

func (d *RootDir) Attr(ctx context.Context, attr *fuse.Attr) error {

	log.Info().Msg(fmt.Sprintf("[FUSE]>> reading root-dir attributes"))

	attr.Inode = 1
	attr.Mode = os.ModeDir | 0o555
	return nil
}

func (d *RootDir) Lookup(ctx context.Context, name string) (fs.Node, error) {

	for _, f := range d.parent.files {
		if f.name == name {

			log.Info().Msg(fmt.Sprintf("[FUSE]>> looking up file '%s' (found raw @ %d)", name, f.inode))

			return &f, nil
		}
	}

	for _, f := range d.parent.specs {
		if f.name == name {

			log.Info().Msg(fmt.Sprintf("[FUSE]>> looking up file '%s' (found spec @ %d)", name, f.inode))

			return &f, nil
		}
	}

	log.Error().Msg(fmt.Sprintf("[FUSE]>> looking up file '%s' (not-found)", name))

	return nil, syscall.ENOENT
}

func (d *RootDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {

	log.Info().Msg(fmt.Sprintf("[FUSE]>> Listing entries of root-dir"))

	e1 := langext.ArrMap(d.parent.files, func(v rawFile) fuse.Dirent { return v.Dirent() })
	e2 := langext.ArrMap(d.parent.specs, func(v specFile) fuse.Dirent { return v.Dirent() })

	return langext.ArrConcat(e1, e2), nil
}
