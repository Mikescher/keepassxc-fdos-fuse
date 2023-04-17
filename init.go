package kxcfuse

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"os"
	"runtime/debug"
)

var Version = "1.1"

var CommitHash *string
var VCSType *string
var CommitTime *string
var CommitModified *string

func init() {
	initZerolog()

	rbi, ok := debug.ReadBuildInfo()
	if !ok {
		log.Fatal().Msg("Failed to read BuildInfo")
		return
	}

	VCSType = getBuildInfoSetting(rbi, "vcs")
	CommitTime = getBuildInfoSetting(rbi, "vcs.time")
	CommitHash = getBuildInfoSetting(rbi, "vcs.revision")
	CommitModified = getBuildInfoSetting(rbi, "vcs.modified")

	log.Debug().Msg("Initialized")
}

func getBuildInfoSetting(rbi *debug.BuildInfo, key string) *string {
	for _, v := range rbi.Settings {
		if v.Key == key {
			return langext.Ptr(v.Value)
		}
	}
	return nil
}

func initZerolog() {
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
