package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type WeatherResponse struct {
	Name    string  `json:"name"`
	Temp    float64 `json:"temp"`
	Weather string  `json:"weather"`
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	cities := r.URL.Query().Get("cities")
	if cities == "" {
		http.Error(w, `{"error": "citiesパラメータが必要です"}`, http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		http.Error(w, `{"error": "APIキーが設定されていません"}`, http.StatusInternalServerError)
		return
	}

	var weatherResults []WeatherResponse

	for _, city := range strings.Split(cities, ",") {
		city = strings.TrimSpace(city)
		url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=ja", city, apiKey)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("天気情報の取得に失敗しました（%s）: %v", city, err)
			continue
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("レスポンスボディを閉じる際にエラーが発生しました（%s）: %v", city, err)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			log.Printf("天気APIのレスポンスが異常（%s）: ステータスコード %d", city, resp.StatusCode)
			continue
		}

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
			log.Printf("JSONデコードエラー（%s）: %v", city, err)
			continue
		}

		weatherDescription := "データなし"
		if len(data.Weather) > 0 {
			weatherDescription = data.Weather[0].Description
		}

		weatherResults = append(weatherResults, WeatherResponse{
			Name:    data.Name,
			Temp:    data.Main.Temp,
			Weather: weatherDescription,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherResults)
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
		<h1>Weather API Documentation</h1>
		<p><strong>エンドポイント:</strong> <code>/weather</code></p>
		<p><strong>使用方法:</strong> <code>/weather?cities=Tokyo,Osaka,NewYork</code></p>
		<p><strong>レスポンス形式 (JSON):</strong></p>
		<pre>
[
  {
    "name": "Tokyo",
    "temp": 15.0,
    "weather": "clear sky"
  },
  {
    "name": "Osaka",
    "temp": 18.5,
    "weather": "few clouds"
  }
]
		</pre>
	`)
}

func main() {
	http.HandleFunc("/weather", weatherHandler)
	http.HandleFunc("/docs", docsHandler)

	fmt.Println("サーバー起動: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
