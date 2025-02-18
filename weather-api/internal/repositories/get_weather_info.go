package repositories

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "sync"

    "github.com/gin-gonic/gin"
)

// 天気APIのレスポンス構造体
type WeatherResponse struct {
    Name string `json:"name"`
    Main struct {
        Temp     float64 `json:"temp"`
        Humidity int     `json:"humidity"`
    } `json:"main"`
    Weather []struct {
        Description string `json:"description"`
    } `json:"weather"`
}

// 単一都市の天気を取得する関数
func fetchWeather(city string) (*WeatherResponse, error) {
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
func fetchWeatherMultiple(cities []string) (map[string]WeatherResponse, error) {
    results := make(map[string]WeatherResponse)
    var wg sync.WaitGroup
    var mu sync.Mutex

    for _, city := range cities {
        wg.Add(1)

        // goroutine で並列処理
        go func(city string) {
            defer wg.Done()

			 // デバッグログ
			log.Printf("[DEBUG] Fetching weather for city: %s\n", city)


            weather, err := fetchWeather(city)
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
func getWeatherHandler(c *gin.Context) {
    // クエリパラメータ "cities" を取得
    citiesParam := c.Query("cities")
    if citiesParam == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "cities parameter is required"})
        return
    }
	

    // カンマ区切りの都市リストをスライスに変換
    cities := strings.Split(citiesParam, ",")

    // 複数都市の天気情報を取得
    weatherData, err := fetchWeatherMultiple(cities)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // JSON でレスポンスを返す
    c.JSON(http.StatusOK, weatherData)
}