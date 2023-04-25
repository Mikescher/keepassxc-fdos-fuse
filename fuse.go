package kxcfuse

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"fmt"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/cmdext"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"os"
	"os/signal"
	"time"
)

type FuseMount struct {
	directory    string
	connection   *fuse.Conn
	server       *fs.Server
	root         *RootDir
	files        []rawFile
	specs        []specFile
	inodeCounter uint64
	maxFSize     int
}

func NewFuseMount(dir string, maxf int) *FuseMount {
	return &FuseMount{
		directory:    dir,
		inodeCounter: 1000,
		maxFSize:     maxf,
	}
}

func (f *FuseMount) Run() error {
	var err error

	f.connection, err = fuse.Mount(f.directory, fuse.FSName("KeePassXC-fdos-fuse"), fuse.Subtype("kxc"), fuse.ReadOnly())
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Info().Msg(fmt.Sprintf("Signal received - Unmounting fs (via umount)"))

		cr, err := cmdext.RunCommand("umount", []string{f.directory}, langext.Ptr(5*time.Second))
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("Unmounting fs (via cmd) failed"))
		}
		if cr.ExitCode != 0 {
			log.Err(err).Str("output", cr.StdCombined).Msg(fmt.Sprintf("Unmounting fs (via cmd) failed"))
		}

		log.Info().Msg(fmt.Sprintf("Signal received - Closing connection"))

		conn := f.connection
		f.connection = nil

		err = conn.Close()
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("Unmounting fs failed"))
		}
	}()

	defer func() {
		if f.connection != nil {
			log.Info().Msg(fmt.Sprintf("Unmounting fs"))
			err = f.connection.Close()
			if err != nil {
				log.Err(err).Msg(fmt.Sprintf("Unmounting fs failed"))
			}
		}
	}()

	f.server = fs.New(f.connection, nil)

	f.root = &RootDir{parent: f}

	log.Info().Msg(fmt.Sprintf("Fuse Filesystem mounted on '%s'", f.directory))

	if err := f.server.Serve(f); err != nil {
		return err
	}

	return nil
}

func (f *FuseMount) Root() (fs.Node, error) {

	log.Info().Msg(fmt.Sprintf("[FUSE]>> reading root node"))

	return f.root, nil
}

func (f *FuseMount) RegisterRawFile(fn string, data []byte) {
	f.inodeCounter += 1

	log.Info().Msg(fmt.Sprintf("Registered (raw) file '%s' (@ %d)", fn, f.inodeCounter))

	f.files = append(f.files, rawFile{
		parent:  f,
		name:    fn,
		inode:   f.inodeCounter,
		content: data,
	})
}

func (f *FuseMount) RegisterSecretServiceSpec(fn string, identkey string, identval string, attr string, extra []string) {
	f.inodeCounter += 1

	log.Info().Msg(fmt.Sprintf("Registered (secret-service) file '%s' (@ %d) -> [[%s :: %s :: %s :: %v]]", fn, f.inodeCounter, identkey, identval, attr, extra))

	f.specs = append(f.specs, specFile{
		parent:     f,
		name:       fn,
		inode:      f.inodeCounter,
		ssIdentKey: identkey,
		ssIdentVal: identval,
		ssAttr:     attr,
		extraFlags: extra,
	})
}
