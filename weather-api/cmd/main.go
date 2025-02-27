package main

import (
	"fmt"
	"net/http"
	"weather-api/internal/config"
	"weather-api/internal/routes"
    "os"
    "log"
    
)

func main() {
    // APIキーが取得されているかの確認
    apiKey := os.Getenv("OPEN_WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("APIキーが設定されていません")
	} else {
		fmt.Println("API Key:", apiKey)
	}

    // 設定を取得
    port := config.GetEnv("PORT", "8080")

    if port == "" {
        fmt.Println("PORT is not set. Using default port 8080")
    }

    mux := http.NewServeMux()

	// ルーティングを設定
    routes.SetupRouter(mux)
    fmt.Println("Server is running on:" + port)
    http.ListenAndServe(":"+port, mux)
}