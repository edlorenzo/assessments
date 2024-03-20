package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"assessments/config"
	"assessments/internal/assessments"
	"assessments/internal/assessments/model"
	"assessments/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type assessmentsHandler struct {
	s    assessments.Service
	repo assessments.Repo
	cnf  *config.Config
}

func SetupAssessmentsHandler(router gin.IRouter, s assessments.Service, repo assessments.Repo, cnf *config.Config) {
	hh := &assessmentsHandler{
		s:    s,
		repo: repo,
		cnf:  cnf,
	}

	router.GET("/health/readiness", hh.ReadinessProbe)
	router.GET("/health/liveness", hh.ReadinessProbe)

	endpoints := router.Group("/api/v1")

	endpoints.POST("/data-fetch-and-store-it-db", hh.fetchAndStore)
	endpoints.GET("/stations", hh.getStationsData)
	endpoints.GET("/stations/:kioskId/", hh.getStationsData)
}

func (hh *assessmentsHandler) ReadinessProbe(c *gin.Context) {
	fmt.Print("Entering health check...")
	err := hh.s.HealthCheck(c)

	if err != nil {
		response.EncodeJSONResp(c, gin.H{
			"status":  "failed",
			"message": err.Error,
		}, http.StatusServiceUnavailable)
		return
	}

	response.EncodeJSONResp(c, gin.H{
		"status": "success",
	}, http.StatusOK)
}

// fetchAndStore POST endpoint to trigger data fetch and store
func (hh assessmentsHandler) fetchAndStore(ctx *gin.Context) {
	baseURL := hh.cnf.ExternalAPIConfig.OpenWeatherMapApiUrl
	apiKey := hh.cnf.ExternalAPIConfig.OpenWeatherMapApiKey
	city := hh.cnf.ExternalAPIConfig.OpenWeatherMapApiCityParam
	units := "metric"
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("error parsing URL: %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error parsing JSON."})
		return
	}

	q := u.Query()
	q.Set("q", city)
	q.Set("appid", apiKey)
	q.Set("units", units)
	u.RawQuery = q.Encode()

	log.Info().Msg(fmt.Sprintf("Open Weather Map URL: %s", u.String()))

	weatherData, err := hh.s.FetchWeatherData(ctx, u.String())
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("Failed to fetch weather data: %v", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weather data."})
		return
	}

	indegoApiUrl := hh.cnf.ExternalAPIConfig.IndegoApiUrl
	abbreviation := hh.cnf.ExternalAPIConfig.IndegoApiAbbreviationParam
	newIndegoApiUrl := indegoApiUrl + abbreviation
	log.Info().Msg(fmt.Sprintf("Indego Api URL: %s", newIndegoApiUrl))

	indegoData, err := hh.s.FetchStationData(ctx, newIndegoApiUrl)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("Failed to fetch Indego data: %v", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Indego data."})
		return
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

	err = hh.s.ConsolidatedData(ctx, &weatherObj, &indegoObj)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("Failed to store data in PostgreSQL: %v", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to store data in PostgreSQL. - %v", err)})
		return
	}

	log.Info().Msg("Data fetched and stored successfully...")
	responseData := gin.H{
		"stations": indegoData,
		"weather":  weatherData,
	}

	response.EncodeJSONResp(ctx, responseData, http.StatusOK)
}

// getStationsData Handler to retrieve station and weather data for a specific timestamp
func (hh assessmentsHandler) getStationsData(c *gin.Context) {
	kioskId := c.Params.ByName("kioskId")
	timestampStr := c.Query("at")
	weatherData, err := hh.s.WeatherSnapShotByDate(c.Request.Context(), kioskId, timestampStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve weather data"})
		return
	}

	stationData, err := hh.s.StationSnapShotByDate(c.Request.Context(), kioskId, timestampStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve station data"})
		return
	}

	responseData := gin.H{
		"at":       timestampStr,
		"stations": stationData,
		"weather":  weatherData,
	}

	response.EncodeJSONResp(c, responseData, http.StatusOK)
}
