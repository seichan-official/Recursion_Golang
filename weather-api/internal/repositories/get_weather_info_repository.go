package repositories

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

// 都市名と国名のマップ
var cityCountryMap = map[string]string{
	"New York":    "New York,US",
	"Los Angeles": "Los Angeles,US",
	"Toronto":     "Toronto,CA",
	"London":      "London,GB",
	"Paris":       "Paris,FR",
	"Berlin":      "Berlin,DE",
	"Rome":        "Rome,IT",
	"Shanghai":    "Shanghai,CN",
	"Hong Kong":   "Hong Kong,HK",
	"Seoul":       "Seoul,KR",
	"Bangkok":     "Bangkok,TH",
	"Sydney":      "Sydney,AU",
	"Dubai":       "Dubai,AE",
}

// 外部APIから単一都市の天気情報を取得
func FetchWeatherExternalAPI(city string) (WeatherResponse, error) {
	apiKey := os.Getenv("OPEN_WEATHER_API_KEY")
	if apiKey == "" {
		return WeatherResponse{}, fmt.Errorf("APIキーが設定されていません")
	}

	// 都市名をマッピング
	if mappedCity, exists := cityCountryMap[city]; exists {
		city = mappedCity
	}

	// URL エンコード
	encodedCity := url.QueryEscape(city)

	apiURL := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=ja", encodedCity, apiKey)
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

	//天気情報が空の場合の処理
	weatherDescription := "データなし"
	if len(data.Weather) > 0 {
		weatherDescription = data.Weather[0].Description
	}

	return WeatherResponse{
		Name:    data.Name,
		Temp:    data.Main.Temp,
		Weather: weatherDescription,
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

			// 都市名をマッピング
			if mappedCity, exists := cityCountryMap[city]; exists {
				city = mappedCity
			}

			weather, err := FetchWeatherExternalAPI(city)
			if err != nil {
				log.Printf("Error fetching weather for %s: %v\n", city, err)
				mu.Lock()
				results[city] = WeatherResponse{Name: city, Temp: 0, Weather: "データなし"}
				mu.Unlock()
				return
			}
			mu.Lock()
			results[city] = weather
			mu.Unlock()
		}(city)
	}

	wg.Wait()
	return results, nil
}

// エリアごとの天気取得
func FetchWeatherArea(area string) (map[string]WeatherResponse, error) {
	results := make(map[string]WeatherResponse)

	// エリアが定義されているか確認
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