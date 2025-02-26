package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"weather-api/internal/config"
	"weather-api/internal/models"
)

// GeocodeRepository 構造体
type GeocodeRepository struct{}

// コンストラクタ
func NewGeocodeRepository() *GeocodeRepository {
    return &GeocodeRepository{}
}

// Geocoding API から緯度経度を取得
func (r *GeocodeRepository) FetchCoordinates(cityName string) (*models.GeocodeResponse, error) {
    apiKey := config.GetEnv("OPEN_WEATHER_API_KEY", "")
    if apiKey == "" {
        return nil, fmt.Errorf("APIキーが設定されていません")
    }

    url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", cityName, apiKey)

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var geocodeResults []models.GeocodeResponse
    if err := json.Unmarshal(body, &geocodeResults); err != nil {
        return nil, err
    }

    if len(geocodeResults) == 0 {
        return nil, fmt.Errorf("都市が見つかりませんでした")
    }

    return &geocodeResults[0], nil
}
