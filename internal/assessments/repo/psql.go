package repo

import (
	"context"
	"fmt"
	"strconv"

	"assessments/internal/assessments/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) StoreDataInDB(ctx context.Context, tx *gorm.DB, weather *model.Weather, station *model.BikeAvailability) error {
	err := tx.WithContext(ctx).Create(weather).Error
	if err != nil {
		return fmt.Errorf("tx.Create failed to insert weather data: %v", err)
	}

	err = tx.WithContext(ctx).Create(station).Error
	if err != nil {
		return fmt.Errorf("tx.Create failed to insert indego data: %v", err)
	}

	return nil
}

func (r *repo) GetWeatherDataByDate(ctx context.Context, kioskId, timestamp string) (data []model.Weather, err error) {
	var weatherData []model.Weather

	if len(kioskId) > 0 {
		newKioskId, err := strconv.Atoi(kioskId)
		if err != nil {
			return weatherData, fmt.Errorf("errr converting kioskId string to int: %w", err)
		}
		err = r.db.Raw("SELECT * FROM weathers WHERE to_char(timestamp, 'YYYY-MM-DD\"T\"HH24:MI:SS') = ? AND kiosk_id = ?", timestamp, newKioskId).Scan(&weatherData).Error
	} else {
		err = r.db.Raw("SELECT * FROM weathers WHERE to_char(timestamp, 'YYYY-MM-DD\"T\"HH24:MI:SS') = ?", timestamp).Scan(&weatherData).Error
	}

	if err != nil {
		return weatherData, err
	}

	return weatherData, nil
}

func (r *repo) GetStationDataByDate(ctx context.Context, kioskId, timestamp string) (data []model.BikeAvailability, err error) {
	var station []model.BikeAvailability

	if len(kioskId) > 0 {
		newKioskId, err := strconv.Atoi(kioskId)
		if err != nil {
			return station, fmt.Errorf("errr converting kioskId string to int: %w", err)
		}
		err = r.db.Raw("SELECT * FROM bike_availabilities WHERE to_char(timestamp, 'YYYY-MM-DD\"T\"HH24:MI:SS') = ? AND kiosk_id = ?", timestamp, newKioskId).Scan(&station).Error
	} else {
		err = r.db.Raw("SELECT * FROM bike_availabilities WHERE to_char(timestamp, 'YYYY-MM-DD\"T\"HH24:MI:SS') = ?", timestamp).Scan(&station).Error
	}

	if err != nil {
		return station, err
	}

	return station, nil
}

func (r *repo) HealthCheck(ctx context.Context) error {
	alive := false
	res := r.db.WithContext(ctx).Raw("SELECT true").Scan(&alive)
	if res.Error != nil {
		log.Debug().Msgf("database is not reachable: ", res.Error)
		return fmt.Errorf("db.Raw: %w", res.Error)
	}
	return nil
}

func (r *repo) RunTransaction(ctx context.Context, fn func(tx *gorm.DB) error, onCommitFail func() error) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return onCommitFail()
	}
	return nil
}
