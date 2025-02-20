package controllers

import (
	"encoding/json"
	"net/http"
	"weather-api/internal/services"
)

type WeatherController struct {
    WeatherService services.WeatherService
}

func NewWeatherController() *WeatherController {
    return &WeatherController{
        WeatherService: *services.NewWeatherService(),
    }
}

func (wc *WeatherController) GetWeatherAlerts(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    lat := r.URL.Query().Get("lat")
    lon := r.URL.Query().Get("lon")

    if lat == "" || lon == "" {
        http.Error(w, "lat と lon パラメータが必要です", http.StatusBadRequest)
        return
    }

    alerts, err := wc.WeatherService.GetWeatherAlerts(lat, lon)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(alerts)
}
