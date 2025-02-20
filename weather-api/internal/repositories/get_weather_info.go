package repositories

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "sync"
)

// 天気APIのレスポンス構造体
type WeatherResponse struct {}

// func FetchWeather(city string) (*WeatherResponse, error) {
//     apiKey := os.Getenv("OPEN_WEATHER_API_KEY")
//     baseURL := os.Getenv("BASE_URL")

//     if apiKey == "" || baseURL == "" {
//         log.Println("[WARNING] APIキーまたはBASE_URLが設定されていません。仮データを返します。")
//         return &WeatherResponse{
//             Name: city,
//             Main: struct {
//                 Temp     float64 `json:"temp"`
//                 Humidity int     `json:"humidity"`
//             }{Temp: 25.0, Humidity: 60},
//             Weather: []struct {
//                 Description string `json:"description"`
//             }{{Description: "clear sky"}},
//         }, nil
//     }

//     url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric&lang=ja", baseURL, city, apiKey)
//     log.Printf("[DEBUG] Requesting weather data from: %s\n", url)

//     resp, err := http.Get(url)
//     if err != nil {
//         log.Printf("[ERROR] Failed to fetch weather data for %s: %v\n", city, err)
//         log.Println("[WARNING] API未提供のため、仮データを返します。")
//         return &WeatherResponse{
//             Name: city,
//             Main: struct {
//                 Temp     float64 `json:"temp"`
//                 Humidity int     `json:"humidity"`
//             }{Temp: 25.0, Humidity: 60},
//             Weather: []struct {
//                 Description string `json:"description"`
//             }{{Description: "clear sky"}},
//         }, nil
//     }
//     defer resp.Body.Close()

//     if resp.StatusCode != http.StatusOK {
//         log.Printf("[ERROR] API returned non-200 status for %s: %d\n", city, resp.StatusCode)
//         return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
//     }

//     var weatherData WeatherResponse
//     if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
//         log.Printf("[ERROR] Failed to decode JSON for %s: %v\n", city, err)
//         return nil, err
//     }

//     log.Printf("[DEBUG] Successfully fetched weather data for %s: %+v\n", city, weatherData)
//     return &weatherData, nil
// }

// 単一都市の天気を取得する関数
func FetchWeather(city string) (*WeatherResponse, error) {
    apiKey := os.Getenv("OPEN_WEATHER_API_KEY")
    baseURL := os.Getenv("BASE_URL")

    if apiKey == "" || baseURL == "" {
        return nil, fmt.Errorf("APIキーまたはBASE_URLが設定されていません")
    }

    url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric&lang=ja", baseURL, city, apiKey)
    // APIリクエストのURLをログ出力
    log.Printf("[DEBUG] Requesting weather data from: %s\n", url)

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    // APIレスポンスのステータスコードをログ出力
    log.Printf("[DEBUG] Received response for %s with status: %d\n", city, resp.StatusCode)

    var weatherData WeatherResponse
    if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
        return nil, err
    }

    return &weatherData, nil
}

// 複数都市の天気を取得する関数
func FetchWeatherMultiple(cities []string) (map[string]WeatherResponse, error) {
    results := make(map[string]WeatherResponse)
    var wg sync.WaitGroup
    var mu sync.Mutex

    for _, city := range cities {
        wg.Add(1)

        go func(city string) {
            defer wg.Done()

			 // デバッグログ
			log.Printf("[DEBUG] Fetching weather for city: %s\n", city)


            weather, err := FetchWeather(city)
            if err != nil {
                log.Printf("Error fetching weather for %s: %v\n", city, err)
                return
					
            }
			// デバッグログ
			log.Printf("[DEBUG] Weather data received for %s: %+v\n", city, weather)
            // 結果をスレッドセーフに保存
            mu.Lock()
            results[city] = *weather
            mu.Unlock()
        }(city)
    }

    wg.Wait()
    return results, nil
}

// 天気情報を取得するエンドポイント
func GetWeatherHandler(w http.ResponseWriter, r *http.Request) {
    citiesParam := r.URL.Query().Get("cities")
    if citiesParam == "" {
        http.Error(w, `{"error": "cities parameter is required"}`, http.StatusBadRequest)
        return
    }

    // カンマ区切りの都市リストをスライスに変換
    cities := strings.Split(citiesParam, ",")

    // 複数都市の天気情報を取得
    weatherData, err := FetchWeatherMultiple(cities)
    if err != nil {
        http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }

    // JSON でレスポンスを返す
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(weatherData)
}