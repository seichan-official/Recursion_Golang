package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"weather-api/internal/config"
	"weather-api/internal/models"
)

type WeatherRepository struct{}

func NewWeatherRepository() *WeatherRepository {
    return &WeatherRepository{}
}

func (r *WeatherRepository) FetchWeatherAlerts(lat, lon string) ([]models.WeatherAlert, error) {
    apiKey := config.GetEnv("OPEN_WEATHER_API_KEY", "")
    if apiKey == "" {
        return nil, fmt.Errorf("APIキーが設定されていません")
    }

    url := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%s&lon=%s&exclude=current,minutely,hourly,daily&appid=%s", lat, lon, apiKey)
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("APIリクエストに失敗しました: %v", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("レスポンスの読み取りに失敗しました: %v", err)
    }

    var weatherResponse models.WeatherResponse
    if err := json.Unmarshal(body, &weatherResponse); err != nil {
        return nil, fmt.Errorf("JSONの解析に失敗しました: %v", err)
    }

    return weatherResponse.Alerts, nil
}
