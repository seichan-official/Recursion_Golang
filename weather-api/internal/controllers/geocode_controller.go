package controllers

import (
    "encoding/json"
    "weather-api/internal/services"
    "net/http"
)

// GeocodeController 構造体
type GeocodeController struct {
    GeocodeService services.GeocodeService
}

// コンストラクタ
func NewGeocodeController() *GeocodeController {
    return &GeocodeController{
        GeocodeService: *services.NewGeocodeService(),
    }
}

// 緯度経度を取得するエンドポイント
func (gc *GeocodeController) GetCoordinates(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    city := r.URL.Query().Get("city")
    if city == "" {
        http.Error(w, "City parameter is required", http.StatusBadRequest)
        return
    }

    coordinates, err := gc.GeocodeService.GetCoordinates(city)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(coordinates)
}
