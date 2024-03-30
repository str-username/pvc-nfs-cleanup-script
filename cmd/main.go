package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"pvc-cleanup/cleanup"
	"time"
)

func main() {
	runMode := os.Getenv("RUN_MODE")
	runInterval := os.Getenv("RUN_INTERVAL")
	checkFolders := os.Getenv("CHECK_FOLDERS")

	if checkFolders == "" {
		log.Fatal().Msg("environment CHECK_FOLDERS not set")
		os.Exit(10)
	}

	if runMode == "" {
		log.Fatal().Msg("environment RUN_MODE not set")
		os.Exit(10)
	}

	if runMode == "job" {
		log.Info().Msg("run as job mode")
		cleanup.RemovePvcs(checkFolders)
	}

	if runMode == "always" {
		if runInterval == "" {
			log.Fatal().Msg("for always mode set RUN_INTERVAL environment, eg: 5s, 5m")
			os.Exit(10)
		}
		interval, err := time.ParseDuration(runInterval)

		if err != nil {
			log.Fatal().Err(err).Msg("invalid duration format")
			os.Exit(10)
		}
		ticker := time.Tick(interval)

		log.Info().Str("run_mode", "always").Str("interval", interval.String()).Msg("run")

		for range ticker {
			sTime := time.Now()
			cleanup.RemovePvcs(checkFolders)
			eTime := time.Since(sTime)
			log.Info().Str("execution_time", "seconds").Msg(eTime.String())
		}
	}
}
