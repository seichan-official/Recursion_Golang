package services

import (
	"fmt"
	"weather-api/internal/repositories"
)

// GetWeatherInfoService 構造体
type GetWeatherInfoService struct{}

// コンストラクタ
func NewGetWeatherInfoService() *GetWeatherInfoService {
	return &GetWeatherInfoService{}
}

// 単一都市の天気を取得するメソッド
func (s *GetWeatherInfoService) GetWeatherForCity(city string) (repositories.WeatherResponse, error) {
	weather, err := repositories.FetchWeatherExternalAPI(city)
	if err != nil {
		return repositories.WeatherResponse{}, fmt.Errorf("都市 %s の天気情報を取得できませんでした: %w", city, err)
	}
	return weather, nil
}

// 複数都市の天気を取得するメソッド
func (s *GetWeatherInfoService) GetWeatherForCities(cities []string) (map[string]repositories.WeatherResponse, error) {
	weatherData, err := repositories.FetchWeatherMultiple(cities)
	if err != nil {
		return nil, fmt.Errorf("複数都市の天気取得に失敗しました: %w", err)
	}
	return weatherData, nil
}

// エリア（世界規模）の天気を取得するメソッド
func (s *GetWeatherInfoService) GetWeatherForArea(area string) (map[string]repositories.WeatherResponse, error) {
	weatherData, err := repositories.FetchWeatherArea(area)
	if err != nil {
		return nil, fmt.Errorf("エリア %s の天気取得に失敗しました: %w", area, err)
	}
	return weatherData, nil
}