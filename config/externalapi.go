package config

type ExternalAPIConfig struct {
	OpenWeatherMapApiKey       string
	OpenWeatherMapApiUrl       string
	OpenWeatherMapApiCityParam string
	IndegoApiUrl               string
	IndegoApiAbbreviationParam string
	JobDelayMin                int
}

func getExternalAPIConfig() ExternalAPIConfig {
	return ExternalAPIConfig{
		OpenWeatherMapApiKey:       getString("OPEN_WEATHER_MAP_API_KEY"),
		OpenWeatherMapApiUrl:       getString("OPEN_WEATHER_MAP_API_URL"),
		OpenWeatherMapApiCityParam: getString("OPEN_WEATHER_MAP_API_CITY_PARAM"),
		IndegoApiUrl:               getString("INDEGO_API_URL"),
		IndegoApiAbbreviationParam: getString("INDEGO_API_ABBREVIATION_PARAM"),
		JobDelayMin:                getIntOrDefault("JOB_DELAY_MIN", 60),
	}
}
