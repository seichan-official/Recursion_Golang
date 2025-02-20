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
		// OpenWeatherMap API の URL を生成（指定した都市の現在の気象情報を取得、単位は摂氏）
		url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			continue
		}
		resp.Body.Close()
		//resp.Bodyを受け取り終わったら、それ以上のデータは受け取らないためにresp.Bodyを閉じる

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
			continue
		}

		weatherResults = append(weatherResults, WeatherResponse{
			Name:    data.Name,
			Temp:    data.Main.Temp,
			Weather: data.Weather[0].Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherResults)
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
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
