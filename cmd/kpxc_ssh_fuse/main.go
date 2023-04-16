package kpxc_ssh_fuse

import (
	"fmt"
	kxcfuse "keepassxc-fdos-fuse"
	"os"
)

func main() {
	conf := kxcfuse.ParseArgs(os.Args[1:])

	mount := kxcfuse.NewFuseMount(conf.MountDir, 1024*1024)

	if kxcfuse.CommitHash != nil {
		mount.RegisterRawFile("@COMMIT_HASH", []byte(*kxcfuse.CommitHash))
	}
	if kxcfuse.VCSType != nil {
		mount.RegisterRawFile("@VCS_TYPE", []byte(*kxcfuse.VCSType))
	}
	if kxcfuse.CommitTime != nil {
		mount.RegisterRawFile("@COMMIT_TIME", []byte(*kxcfuse.CommitTime))
	}
	if kxcfuse.CommitModified != nil {
		mount.RegisterRawFile("@COMMIT_MODIFIED", []byte(*kxcfuse.CommitModified))
	}
	mount.RegisterRawFile("@PID", []byte(fmt.Sprintf("%d", os.Getpid())))

	for _, spec := range conf.Spec {
		mount.RegisterSecretServiceSpec(spec.Filename, spec.IdentKey, spec.IdentVal, spec.Attr)
	}

	err := mount.Run()
	if err != nil {
		panic(err)
	}
}
