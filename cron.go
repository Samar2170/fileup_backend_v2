package main

import (
	"fileupbackendv2/config"
	"fileupbackendv2/internal/dirManager"
	"fileupbackendv2/internal/storage/image"
	"fileupbackendv2/pkg/logging"
	"time"

	"github.com/go-co-op/gocron"
)

func StartCronServer() {
	t := time.Now()
	logging.AuditLogger.Info().Msgf("Starting cron server at %s", t.Format(time.RFC3339))
	s := gocron.NewScheduler(time.UTC)
	s.Every(2).Hour().Do(func() {
		err := image.MarkImages()
		if err != nil {
			logging.Errorlogger.Error().Msgf("Error in Marking images: %s", err.Error())
		}
	})
	s.Every(1).Hour().Do(func() {
		err := image.CompressImages(config.CompressionQuality)
		if err != nil {
			logging.Errorlogger.Error().Msgf("Error in Compressing images: %s", err.Error())
		}
	})
	s.Every(1).Hour().Do(func() {
		dirManager.UpdateDirsData()
		dirManager.UpdateUserDirsData()

	})
	s.StartBlocking()
}
