package routes

import (
	"net/http"
	"weather-api/internal/controllers"
)

// ルーティングを設定
func SetupRouter(mux *http.ServeMux) {

	// コントローラーのインスタンスを生成
	geocodeController := controllers.NewGeocodeController()

	NewGetWeatherInfoController := controllers.NewGetWeatherInfoController()

	// ルーティングを設定
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	mux.HandleFunc("/geocode", geocodeController.GetCoordinates)
	mux.HandleFunc("/weather", NewGetWeatherInfoController.GetWeatherHandler)

}
