package main

import (
	"errors"
	"fmt"
	"os"

	"assessments/cmd/app"
	"assessments/cmd/jobs"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	SERVER = "SERVER"
	WORKER = "WORKER"
	FILE   = ".env"
)

func main() {
	if err := setupViper(FILE); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal().Msg(fmt.Sprintf("failed to setup viper err: %s", err.Error()))
		return
	}

	process := getString("PROCESS")

	if process == SERVER {
		app.StartAPP()
	}
	if process == WORKER {
		jobs.SetupDataSyncWorker()
	}
}

func setupViper(file string) error {
	viper.SetConfigFile(file)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func getString(key string) string {
	if !viper.IsSet(key) {
		log.Fatal().
			Str("key", key).
			Msg("Unable to find config value for key")
	}
	return viper.GetString(key)
}
