package services

import (
	"weather-api/internal/models"
	"weather-api/internal/repositories"
)

type WeatherService struct {
    WeatherRepo repositories.WeatherRepository
}

func NewWeatherService() *WeatherService {
    return &WeatherService{
        WeatherRepo: *repositories.NewWeatherRepository(),
    }
}

func (s *WeatherService) GetWeatherAlerts(lat, lon string) ([]models.WeatherAlert, error) {
    return s.WeatherRepo.FetchWeatherAlerts(lat, lon)
}
