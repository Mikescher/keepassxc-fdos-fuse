package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	kxcfuse "keepassxc-fdos-fuse"
	"os"
)

func main() {
	InitZerolog()

	mount := kxcfuse.NewFuseMount("/home/mike/mounts/kxc-fuse/", 1024*1024)

	mount.RegisterRawFile("hello-world1.txt", []byte("hi! 11"))
	mount.RegisterRawFile("hello-world2.txt", []byte("hi! 22"))
	mount.RegisterRawFile("hello-world3.txt", []byte("hi! 33"))

	mount.RegisterSecretServiceSpec("vward-test1.txt", "identifier", "795f02e0-e9d3-4e5b-af6a-2f31040c0791", "xdg:schema")
	mount.RegisterSecretServiceSpec("mongo-test2.txt", "account", "36b8968d-2965-4ff5-a2fe-201cc74bf6df", "service")
	mount.RegisterSecretServiceSpec("mongo-test3.txt", "account", "36b8968d-2965-4ff5-a2fe-201cc74bf6df", "service")
	mount.RegisterSecretServiceSpec("mongo-test4.txt", "account", "36b8968d-2965-4ff5-a2fe-201cc74bf6df", "@password")

	err := mount.Run()
	if err != nil {
		panic(err)
	}
}

func InitZerolog() {
	cw := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05 Z07:00",
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	multi := zerolog.MultiLevelWriter(cw)
	logger := zerolog.New(multi).With().
		Timestamp().
		Caller().
		Logger()

	log.Logger = logger

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	log.Debug().Msg("Initialized")
}
