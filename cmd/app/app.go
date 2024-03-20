package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"assessments/config"
	"assessments/internal/assessments"
	"assessments/internal/assessments/repo"
	"assessments/server"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/rs/zerolog/log"
)

type App struct {
	conf          *config.Config
	assessmentSvc assessments.Service
	server        *server.Server
	cancelFunc    context.CancelFunc
	repoSvc       assessments.Repo
}

func (a *App) Run() error {
	done := make(chan error, 1)
	// Waits for CTRL-C or os SIGINT for server shutdown.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		done <- a.server.GracefulShutdown()
	}()

	if err := a.server.Run(); err != nil {
		return err
	}
	return <-done
}

func (a *App) closeFN() error {
	if err := a.server.Close(); err != nil {
		return fmt.Errorf("could not close server %w", err)
	}
	a.cancelFunc()
	return nil
}

func (a *App) Setup() error {
	srv := server.New(
		a.conf,
		a.assessmentSvc,
		a.repoSvc,
	)
	a.server = srv

	return nil
}

func StartAPP() {
	log.Info().Msg("app: loading config...")
	_, cancelFunc := context.WithCancel(context.Background())
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

	err = RunMigration(sqlDB, c.MigrationPath)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("migration failed %s", err.Error()))
		return
	}

	newRepo := repo.NewRepo(db)
	svc := assessments.NewAssessmentsService(newRepo, db)
	app := &App{conf: c, assessmentSvc: svc}
	if err = app.Setup(); err != nil {
		log.Fatal().Msg(fmt.Sprintf("could not load config %s", err.Error()))
	}

	if err = app.Run(); err != nil {
		log.Fatal().Msg(fmt.Sprintf("could not load config %s", err.Error()))
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	<-sigterm

	cancelFunc()
	err = app.closeFN()
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("could not close %s", err.Error()))
	}
	log.Info().Msg("app: exited!")
}

func RunMigration(db *sql.DB, path string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path),
		"postgres", driver)

	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
