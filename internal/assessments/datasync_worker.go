package assessments

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"assessments/config"
	"assessments/internal/assessments/model"
	"github.com/rs/zerolog/log"
)

type DataSyncWorker struct {
	s    Service
	repo Repo
	cnf  *config.Config
}

func NewDataSyncWorker(
	s Service,
	repo Repo,
	cnf *config.Config,
) *DataSyncWorker {
	return &DataSyncWorker{
		s:    s,
		repo: repo,
		cnf:  cnf,
	}
}

func (d *DataSyncWorker) DataSync(ctx context.Context, u *url.URL, newIndegoApiUrl, city string) error {
	weatherData, err := d.s.FetchWeatherData(ctx, u.String())
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("Failed to fetch weather data: %v", err))
		return err
	}

	indegoData, err := d.s.FetchStationData(ctx, newIndegoApiUrl)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("Failed to fetch Indego data: %v", err))
		return err
	}

	timestamp := time.Now()
	weatherObj := model.Weather{
		Timestamp:   timestamp,
		City:        city,
		Temperature: weatherData.Main.Temp,
		Humidity:    weatherData.Main.Humidity,
		WindSpeed:   weatherData.Wind.Speed,
	}

	indegoObj := model.BikeAvailability{
		Timestamp:              timestamp,
		Name:                   indegoData.Features[0].Properties.Name,
		TotalDocks:             indegoData.Features[0].Properties.TotalDocks,
		DocksAvailable:         indegoData.Features[0].Properties.DocksAvailable,
		BikesAvailable:         indegoData.Features[0].Properties.BikesAvailable,
		ClassicBikesAvailable:  indegoData.Features[0].Properties.ClassicBikesAvailable,
		SmartBikesAvailable:    indegoData.Features[0].Properties.SmartBikesAvailable,
		ElectricBikesAvailable: indegoData.Features[0].Properties.ElectricBikesAvailable,
		RewardBikesAvailable:   indegoData.Features[0].Properties.RewardBikesAvailable,
		RewardDocksAvailable:   indegoData.Features[0].Properties.RewardDocksAvailable,
		KioskStatus:            indegoData.Features[0].Properties.KioskStatus,
		KioskPublicStatus:      indegoData.Features[0].Properties.KioskPublicStatus,
		KioskConnectionStatus:  indegoData.Features[0].Properties.KioskConnectionStatus,
		KioskType:              indegoData.Features[0].Properties.KioskType,
		AddressStreet:          indegoData.Features[0].Properties.AddressStreet,
		AddressCity:            indegoData.Features[0].Properties.AddressCity,
		AddressState:           indegoData.Features[0].Properties.AddressState,
		AddressZipCode:         indegoData.Features[0].Properties.AddressZipCode,
		IsEventBased:           indegoData.Features[0].Properties.IsEventBased,
		IsVirtual:              indegoData.Features[0].Properties.IsVirtual,
		KioskID:                indegoData.Features[0].Properties.KioskID,
		TrikesAvailable:        indegoData.Features[0].Properties.TrikesAvailable,
		Latitude:               indegoData.Features[0].Properties.Latitude,
		Longitude:              indegoData.Features[0].Properties.Longitude,
	}

	err = d.s.ConsolidatedData(ctx, &weatherObj, &indegoObj)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("attempt to save data in the database was unsuccessful: %v", err))
		return err
	}
	return nil
}
