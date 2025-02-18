package main

import (
    "fmt"
    "net/http"
    "strings"

    "weather-api/internal/config"
    "weather-api/internal/repositories"
    "weather-api/internal/routes"

    "github.com/gin-gonic/gin"
)

// 天気情報を取得するエンドポイント
func getWeatherHandler(c *gin.Context) {
    citiesParam := c.Query("cities")
    if citiesParam == "" {  
        c.JSON(http.StatusBadRequest, gin.H{"error": "cities parameter is required"}) 
        return
    }

    cities := strings.Split(citiesParam, ",")

    weatherData, err := repositories.FetchWeatherMultiple(cities)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, weatherData)
}

func main() {
    config.LoadEnv()

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