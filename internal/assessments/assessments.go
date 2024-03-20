package assessments

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"assessments/internal/assessments/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type (
	assessmentsService struct {
		repo Repo
		db   *gorm.DB
	}
)

func NewAssessmentsService(
	repo Repo,
	db *gorm.DB,
) Service {
	return &assessmentsService{
		repo: repo,
		db:   db,
	}
}

func (s *assessmentsService) HealthCheck(ctx context.Context) error {
	return s.repo.HealthCheck(ctx)
}

func (s *assessmentsService) FetchStationData(ctx context.Context, url string) (*bikeAvailabilityResponse, error) {
	var rs bikeAvailabilityResponse
	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Debug().Msg(fmt.Sprintf("error: %s", resp.Status))
		return nil, fmt.Errorf("errore: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&rs); err != nil {
		return nil, err
	}

	return &rs, nil
}

func (s *assessmentsService) FetchWeatherData(ctx context.Context, url string) (*weatherResponse, error) {
	var rs weatherResponse
	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Debug().Msg(fmt.Sprintf("error: %s", resp.Status))
		return nil, fmt.Errorf("error: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&rs); err != nil {
		return nil, err
	}

	return &rs, nil
}

func (s *assessmentsService) ConsolidatedData(ctx context.Context, weather *model.Weather, availability *model.BikeAvailability) error {
	transactionFn := func(tx *gorm.DB) error {
		err := s.repo.StoreDataInDB(ctx, tx, weather, availability)
		if err != nil {
			log.Debug().Msg(fmt.Sprintf("insert data failed: %v", err))
			return err
		}
		return nil
	}

	err := s.repo.RunTransaction(ctx, transactionFn, func() error {
		log.Debug().Msg("inserting consolidated data failed")
		return errors.New("inserting consolidated data failed")
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *assessmentsService) ConsolidatedData2(ctx context.Context, weather *model.Weather, availability *model.BikeAvailability) error {
	transactionFn := func(tx *gorm.DB) error {
		err := s.repo.StoreDataInDB(ctx, tx, weather, availability)
		if err != nil {
			log.Debug().Msg(fmt.Sprintf("insert data failed: %v", err))
			return err
		}
		return nil
	}

	err := s.repo.RunTransaction(ctx, transactionFn, func() error {
		log.Debug().Msg("inserting consolidated data failed")
		return errors.New("inserting consolidated data failed")
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *assessmentsService) WeatherSnapShotByDate(ctx context.Context, kioskId, date string) (*[]weatherResponse, error) {
	var response []weatherResponse
	timestamps, err := time.Parse(time.RFC3339, date)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("error parsing weather date: %v", err))
		return &response, err
	}

	formattedTimestamp := timestamps.Format("2006-01-02T15:04:05")
	rs, err := s.repo.GetWeatherDataByDate(ctx, kioskId, formattedTimestamp)
	if err != nil || len(rs) == 0 {
		return nil, err
	}

	for i, w := range rs {
		if i >= len(response) {
			response = append(response, weatherResponse{}) // Initialize response if needed
		}
		response[i].Weather = append(response[i].Weather, struct {
			Main string `json:"main"`
		}{Main: w.City})
		response[i].Main.Temp = w.Temperature
		response[i].Main.Humidity = w.Humidity
		response[i].Wind.Speed = w.WindSpeed
	}

	return &response, err
}

func (s *assessmentsService) StationSnapShotByDate(ctx context.Context, kioskId, date string) (*[]bikeAvailabilityResponses, error) {
	var response []bikeAvailabilityResponses
	timestamps, err := time.Parse(time.RFC3339, date)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("error parsing station date: %v", err))
		return &response, err
	}

	formattedTimestamp := timestamps.Format("2006-01-02T15:04:05")
	rs, err := s.repo.GetStationDataByDate(ctx, kioskId, formattedTimestamp)
	if err != nil {
		return &response, err
	}

	for _, b := range rs {
		properties := bikeProperties{
			ID:                     b.ID,
			Name:                   b.Name,
			TotalDocks:             b.TotalDocks,
			DocksAvailable:         b.DocksAvailable,
			BikesAvailable:         b.BikesAvailable,
			ClassicBikesAvailable:  b.ClassicBikesAvailable,
			SmartBikesAvailable:    b.SmartBikesAvailable,
			ElectricBikesAvailable: b.ElectricBikesAvailable,
			RewardBikesAvailable:   b.RewardBikesAvailable,
			RewardDocksAvailable:   b.RewardDocksAvailable,
			KioskStatus:            b.KioskStatus,
			KioskPublicStatus:      b.KioskPublicStatus,
			KioskConnectionStatus:  b.KioskConnectionStatus,
			KioskType:              b.KioskType,
			AddressStreet:          b.AddressStreet,
			AddressCity:            b.AddressCity,
			AddressState:           b.AddressState,
			AddressZipCode:         b.AddressZipCode,
			IsEventBased:           b.IsEventBased,
			IsVirtual:              b.IsVirtual,
			KioskID:                b.KioskID,
			TrikesAvailable:        b.TrikesAvailable,
			Latitude:               b.Latitude,
			Longitude:              b.Longitude,
		}
		feature := bikeFeature{
			Properties: properties,
			Type:       "Feature",
		}
		responseObj := bikeAvailabilityResponses{
			LastUpdated: b.Timestamp,
			Features:    []bikeFeature{feature},
			Type:        "FeatureCollection",
		}
		response = append(response, responseObj)
	}

	return &response, err
}
