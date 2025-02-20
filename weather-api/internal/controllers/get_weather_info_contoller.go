package controllers

import (
	"encoding/json"
	"net/http"
    "weather-api/internal/services"
)

// GetWeatherInfoController 構造体
type GetWeatherInfoController struct {
    GetWeatherInfoService *services.WeatherService 
}

// コンストラクタ
func NewGetWeatherInfoController() *GetWeatherInfoController {
    return &GetWeatherInfoController{
        GetWeatherInfoService: services.NewWeatherService(),
    }
}

// 天気情報を取得するエンドポイント
func (gc *GetWeatherInfoController) GetWeatherInfo(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	weatherInfo, err := gc.GetWeatherInfoService.(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherInfo)
}