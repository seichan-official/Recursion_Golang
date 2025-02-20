package services

import (
    "weather-api/internal/models"
    "weather-api/internal/repositories"
)

// GetWeatherInfoService 構造体
type GetWeatherInfoService struct {
    WeatherRepo *repositories.WeatherRepository
}

// コンストラクタ
func NewGetWeatherInfoService() *GetWeatherInfoService {
    return &GetWeatherInfoService{
        WeatherRepo: repositories.NewWeatherRepository(),
    }
}

func (s *GetWeatherInfoService) GetWeatherInfo(city string) (*models.WeatherResponse, error) {
    return s.WeatherRepo.FetchWeather(city)
}