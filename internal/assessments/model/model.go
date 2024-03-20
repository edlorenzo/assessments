package model

import (
	"time"
)

type (
	Weather struct {
		ID          uint      `gorm:"primary_key"`
		Timestamp   time.Time `gorm:"index"`
		City        string
		Temperature float64
		Humidity    float64
		WindSpeed   float64
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
)
