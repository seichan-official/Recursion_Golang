package repositories

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

// WeatherResponse 構造体
type WeatherResponse struct {
	Name    string  `json:"name"`
	Temp    float64 `json:"temp"`
	Weather string  `json:"weather"`
}

// 世界中のエリアごとの都市リスト
var areaCities = map[string][]string{
	"japan":         {"Tokyo", "Osaka", "Nagoya", "Fukuoka", "Sapporo"},
	"north_america": {"New York", "Los Angeles", "Toronto", "Chicago"},
	"europe":        {"London", "Paris", "Berlin", "Rome"},
	"asia":          {"Shanghai", "Hong Kong", "Seoul", "Bangkok"},
	"australia":     {"Sydney"},
	"others":        {"Dubai", "Cape Town", "Rio de Janeiro"},
}

// 外部APIから天気情報を取得
func FetchWeatherExternalAPI(city string) (WeatherResponse, error) {
	apiKey := os.Getenv("OPEN_WEATHER_API_KEY")
	if apiKey == "" {
		return WeatherResponse{}, fmt.Errorf("APIキ-が設定されていません")
	}

	apiURL := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=ja", city, apiKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("外部APIへのリクエスト失敗 %s: %v", city, err)
	}
	defer resp.Body.Close()

	var data struct {
		Name string `json:"name"`
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return WeatherResponse{}, fmt.Errorf("外部APIレスポンスデータのデコード失敗 %s: %v", city, err)
	}

	return WeatherResponse{
		Name:    data.Name,
		Temp:    data.Main.Temp,
		Weather: data.Weather[0].Description,
	}, nil
}

// 複数都市の天気を取得
func FetchWeatherMultiple(cities []string) (map[string]WeatherResponse, error) {
    results := make(map[string]WeatherResponse)
    var wg sync.WaitGroup
    var mu sync.Mutex

    for _, city := range cities {
        wg.Add(1)
        go func(city string) {
            defer wg.Done()
            weather, err := FetchWeatherExternalAPI(city)
            if err != nil {
                log.Printf("Error fetching weather for %s: %v\n", city, err)
                mu.Lock()
                results[city] = WeatherResponse{
                    Name:    city,
                    Temp:    0,
                    Weather: "データなし",
                }
                mu.Unlock()
                return
            }
            mu.Lock()
            results[city] = weather
            mu.Unlock()
        }(city)
    }

    wg.Wait()
    return results,nil
}

// エリアごとの天気取得
func FetchWeatherArea(area string) (map[string]WeatherResponse, error) {
    results := make(map[string]WeatherResponse)
    
	// エリアごとの主要都市リスト
    areaCities := map[string][]string{
       	"japan":         {"Tokyo", "Osaka", "Nagoya", "Fukuoka", "Sapporo"},
		"north_america": {"New York", "Los Angeles", "Toronto", "Chicago"},
		"europe":        {"London", "Paris", "Berlin", "Rome"},
		"asia":          {"Shanghai", "Hong Kong", "Seoul", "Bangkok"},
		"australia":     {"Sydney"},
		"others":        {"Dubai", "Cape Town", "Rio de Janeiro"},
    }

    cities, exists := areaCities[area]
    if !exists {
        return nil, fmt.Errorf("指定されたエリア '%s' の天気情報は取得できません", area)
    }

    // 各都市の天気情報を取得
    for _, city := range cities {
        weather, err := FetchWeatherExternalAPI(city)
        if err != nil {
            log.Printf("Error fetching weather for %s: %v\n", city, err)
            results[city] = WeatherResponse{Name: city, Temp: 0, Weather: "データなし"}
            continue
        }
        results[city] = weather
    }

    return results, nil
}