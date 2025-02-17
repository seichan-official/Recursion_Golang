package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// OpenWeatherMapのレスポンス構造体
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

// OpenWeatherMap APIを呼び出す関数
func fetchWeather(city string) (*WeatherResponse, error) {
	// .env の読み込みを削除し、環境変数を直接参照
	apiKey := os.Getenv("API_KEY")
	baseURL := os.Getenv("BASE_URL")

	if apiKey == "" || baseURL == "" {
		return nil, fmt.Errorf("APIキーまたはBASE_URLが設定されていません")
	}

	url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric&lang=ja", baseURL, city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weatherData WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, err
	}

	return &weatherData, nil
}

// 天気情報を取得するAPIエンドポイント
func getWeatherHandler(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"エラー": "都市名を指定してください"})
		return
	}

	weather, err := fetchWeather(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"エラー": "天気情報の取得に失敗しました"})
		return
	}

	// 天気情報をJSON形式で返す
	c.JSON(http.StatusOK, gin.H{
		"都市":  weather.Name,
		"気温":  fmt.Sprintf("%.1f°C", weather.Main.Temp),
		"湿度":  fmt.Sprintf("%d%%", weather.Main.Humidity),
		"天気":  weather.Weather[0].Description,
	})
}

// メイン関数
func main() {
	r := gin.Default()

	// GETリクエストで天気情報を取得
	r.GET("/weather", getWeatherHandler)

	// サーバーを起動
	r.Run(":8080") // ポート8080でリッスン
}