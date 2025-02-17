package services

import (
    "weather-api/internal/models"
    "weather-api/internal/repositories"
)

// GeocodeService 構造体
type GeocodeService struct {
    GeocodeRepo repositories.GeocodeRepository
}

// コンストラクタ
func NewGeocodeService() *GeocodeService {
    return &GeocodeService{
        GeocodeRepo: *repositories.NewGeocodeRepository(),
    }
}

// 緯度経度を取得する処理
func (s *GeocodeService) GetCoordinates(city string) (*models.GeocodeResponse, error) {
    return s.GeocodeRepo.FetchCoordinates(city)
}
