package jobs

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"assessments/cmd/app"
	"assessments/config"
	"assessments/internal/assessments"
	"assessments/internal/assessments/repo"
	"github.com/rs/zerolog/log"
)

type DataSyncDaemon struct {
	s    assessments.Service
	repo assessments.Repo
	cnf  *config.Config
}

func NewDataSyncDaemon(
	s assessments.Service,
	repo assessments.Repo,
	cnf *config.Config,
) *DataSyncDaemon {
	return &DataSyncDaemon{
		s:    s,
		repo: repo,
		cnf:  cnf,
	}
}

func SetupDataSyncWorker() {
	ctx := context.Background()
	c, err := config.LoadDefault()
	if err != nil {
		log.Error().Msg(fmt.Sprintf("could not load config %s", err.Error()))
		return
	}

	db, err := c.DbConfig.GetDB()
	if err != nil {
		log.Error().Msg(fmt.Sprintf("could not init db %s", err.Error()))
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Msg(fmt.Sprintf("could not init sql db %s", err.Error()))
		return
	}

	err = app.RunMigration(sqlDB, c.MigrationPath)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("migration failed %s", err.Error()))
		return
	}

	newRepo := repo.NewRepo(db)
	svc := assessments.NewAssessmentsService(newRepo, db)

	externalApiDaemon := NewDataSyncDaemon(svc, newRepo, c)
	externalApiDaemon.StartReconcilerDaemon(ctx)
}

func (d *DataSyncDaemon) StartReconcilerDaemon(ctx context.Context) {
	log.Info().Msg("Data sync daemon started...... ")

	baseURL := d.cnf.ExternalAPIConfig.OpenWeatherMapApiUrl
	apiKey := d.cnf.ExternalAPIConfig.OpenWeatherMapApiKey
	city := d.cnf.ExternalAPIConfig.OpenWeatherMapApiCityParam
	units := "metric"
	openWeatherMapURL, err := url.Parse(baseURL)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("error parsing URL: %v", err))
		return
	}

	q := openWeatherMapURL.Query()
	q.Set("q", city)
	q.Set("appid", apiKey)
	q.Set("units", units)
	openWeatherMapURL.RawQuery = q.Encode()

	indegoApiUrl := d.cnf.ExternalAPIConfig.IndegoApiUrl
	abbreviation := d.cnf.ExternalAPIConfig.IndegoApiAbbreviationParam
	newIndegoApiUrl := indegoApiUrl + abbreviation

	log.Info().Msg(fmt.Sprintf("Open Weather Map URL: %s", openWeatherMapURL.String()))
	log.Info().Msg(fmt.Sprintf("Indego Api URL: %s", newIndegoApiUrl))
	log.Info().Msg(fmt.Sprintf("Job delay: %v", time.Duration(d.cnf.ExternalAPIConfig.JobDelayMin)*time.Minute))

	externalApiDaemon := assessments.NewDataSyncWorker(d.s, d.repo, d.cnf)
	ctx, cancel := context.WithCancel(ctx)

	delay := time.Duration(d.cnf.ExternalAPIConfig.JobDelayMin) * time.Minute
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				err := externalApiDaemon.DataSync(ctx, openWeatherMapURL, newIndegoApiUrl, city)
				if err != nil {
					log.Debug().Msg(fmt.Sprintf("Error executing DataSync: %v", err))
				} else {
					log.Info().Msg(fmt.Sprintf("Successfully inserted data!"))
				}
			case <-ctx.Done():
				log.Info().Msg("Context cancelled. Exiting DataSync Job...")
				return
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	cancel()
	wg.Wait()
}
