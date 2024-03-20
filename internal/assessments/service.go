package assessments

import (
	"context"
	"time"

	"assessments/internal/assessments/model"
	"gorm.io/gorm"
)

type (
	Service interface {
		HealthCheck(ctx context.Context) error
		FetchWeatherData(ctx context.Context, url string) (*weatherResponse, error)
		FetchStationData(ctx context.Context, url string) (*bikeAvailabilityResponse, error)
		ConsolidatedData(ctx context.Context, weather *model.Weather, availability *model.BikeAvailability) error
		StationSnapShotByDate(ctx context.Context, kioskId, date string) (*[]bikeAvailabilityResponses, error)
		WeatherSnapShotByDate(ctx context.Context, kioskId, date string) (*[]weatherResponse, error)
	}

	Repo interface {
		StoreDataInDB(ctx context.Context, tx *gorm.DB, weather *model.Weather, station *model.BikeAvailability) error
		GetWeatherDataByDate(ctx context.Context, kioskId, date string) ([]model.Weather, error)
		GetStationDataByDate(ctx context.Context, kioskId, date string) ([]model.BikeAvailability, error)
		RunTransaction(ctx context.Context, fn func(tx *gorm.DB) error, onCommitFail func() error) error
		HealthCheck(ctx context.Context) error
	}

	DataSyncDaemon interface {
		StartReconcilerDaemon(ctx context.Context)
	}
)

type (
	weatherResponse struct {
		Weather []struct {
			Main string `json:"main"`
		} `json:"weather"`
		Main struct {
			Temp     float64 `json:"temp"`
			Humidity float64 `json:"humidity"`
		} `json:"main"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
	}

	bikeAvailabilityResponse struct {
		LastUpdated time.Time `json:"last_updated"`
		Features    []struct {
			Geometry struct {
				Coordinates []float64 `json:"coordinates"`
				Type        string    `json:"type"`
			} `json:"geometry"`
			Properties struct {
				ID                     int           `json:"id"`
				Name                   string        `json:"name"`
				Coordinates            []float64     `json:"coordinates"`
				TotalDocks             int           `json:"totalDocks"`
				DocksAvailable         int           `json:"docksAvailable"`
				BikesAvailable         int           `json:"bikesAvailable"`
				ClassicBikesAvailable  int           `json:"classicBikesAvailable"`
				SmartBikesAvailable    int           `json:"smartBikesAvailable"`
				ElectricBikesAvailable int           `json:"electricBikesAvailable"`
				RewardBikesAvailable   int           `json:"rewardBikesAvailable"`
				RewardDocksAvailable   int           `json:"rewardDocksAvailable"`
				KioskStatus            string        `json:"kioskStatus"`
				KioskPublicStatus      string        `json:"kioskPublicStatus"`
				KioskConnectionStatus  string        `json:"kioskConnectionStatus"`
				KioskType              int           `json:"kioskType"`
				AddressStreet          string        `json:"addressStreet"`
				AddressCity            string        `json:"addressCity"`
				AddressState           string        `json:"addressState"`
				AddressZipCode         string        `json:"addressZipCode"`
				Bikes                  []interface{} `json:"bikes"`
				CloseTime              interface{}   `json:"closeTime"`
				EventEnd               interface{}   `json:"eventEnd"`
				EventStart             interface{}   `json:"eventStart"`
				IsEventBased           bool          `json:"isEventBased"`
				IsVirtual              bool          `json:"isVirtual"`
				KioskID                int           `json:"kioskId"`
				Notes                  interface{}   `json:"notes"`
				OpenTime               interface{}   `json:"openTime"`
				PublicText             string        `json:"publicText"`
				TimeZone               interface{}   `json:"timeZone"`
				TrikesAvailable        int           `json:"trikesAvailable"`
				Latitude               float64       `json:"latitude"`
				Longitude              float64       `json:"longitude"`
			} `json:"properties"`
			Type string `json:"type"`
		} `json:"features"`
		Type string `json:"type"`
	}

	BikeAvailability struct {
		ID                     int       `gorm:"primaryKey"`
		Timestamp              time.Time `gorm:"index"`
		Name                   string
		TotalDocks             int
		DocksAvailable         int
		BikesAvailable         int
		ClassicBikesAvailable  int
		SmartBikesAvailable    int
		ElectricBikesAvailable int
		RewardBikesAvailable   int
		RewardDocksAvailable   int
		KioskStatus            string
		KioskPublicStatus      string
		KioskConnectionStatus  string
		KioskType              int
		AddressStreet          string
		AddressCity            string
		AddressState           string
		AddressZipCode         string
		IsEventBased           bool
		IsVirtual              bool
		KioskID                int
		TrikesAvailable        int
		Latitude               float64
		Longitude              float64
	}

	bikeProperties struct {
		ID                     int     `json:"id"`
		Name                   string  `json:"name"`
		TotalDocks             int     `json:"totalDocks"`
		DocksAvailable         int     `json:"docksAvailable"`
		BikesAvailable         int     `json:"bikesAvailable"`
		ClassicBikesAvailable  int     `json:"classicBikesAvailable"`
		SmartBikesAvailable    int     `json:"smartBikesAvailable"`
		ElectricBikesAvailable int     `json:"electricBikesAvailable"`
		RewardBikesAvailable   int     `json:"rewardBikesAvailable"`
		RewardDocksAvailable   int     `json:"rewardDocksAvailable"`
		KioskStatus            string  `json:"kioskStatus"`
		KioskPublicStatus      string  `json:"kioskPublicStatus"`
		KioskConnectionStatus  string  `json:"kioskConnectionStatus"`
		KioskType              int     `json:"kioskType"`
		AddressStreet          string  `json:"addressStreet"`
		AddressCity            string  `json:"addressCity"`
		AddressState           string  `json:"addressState"`
		AddressZipCode         string  `json:"addressZipCode"`
		IsEventBased           bool    `json:"isEventBased"`
		IsVirtual              bool    `json:"isVirtual"`
		KioskID                int     `json:"kioskId"`
		TrikesAvailable        int     `json:"trikesAvailable"`
		Latitude               float64 `json:"latitude"`
		Longitude              float64 `json:"longitude"`
	}

	bikeFeature struct {
		Properties bikeProperties `json:"properties"`
		Type       string         `json:"type"`
	}

	bikeAvailabilityResponses struct {
		LastUpdated time.Time     `json:"last_updated"`
		Features    []bikeFeature `json:"features"`
		Type        string        `json:"type"`
	}
)
