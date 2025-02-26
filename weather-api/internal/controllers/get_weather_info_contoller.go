package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"weather-api/internal/services"
)

// GetWeatherInfoController 構造体
type GetWeatherInfoController struct {
	GetWeatherInfoService *services.GetWeatherInfoService
}

// コンストラクタ
func NewGetWeatherInfoController() *GetWeatherInfoController {

	return &GetWeatherInfoController{
		GetWeatherInfoService: services.NewGetWeatherInfoService(),
	}
}

// 天気情報取得エンドポイント
func (gc *GetWeatherInfoController) GetWeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	cities := r.URL.Query().Get("cities")
	area := r.URL.Query().Get("area")

	var response interface{}
	var err error

	if city != "" {
		response, err = gc.GetWeatherInfoService.GetWeatherForCity(city)
	} else if cities != "" {
		citiesList := strings.Split(cities, ",")
		response, err = gc.GetWeatherInfoService.GetWeatherForCities(citiesList)
	} else if area != "" {
		response, err = gc.GetWeatherInfoService.GetWeatherForArea(area)
	} else {
		http.Error(w, `{"error": "city, cities, area のいずれかが必要です"}`, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// 正常なレスポンスを JSON 形式で返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
} 