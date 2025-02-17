package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// モックAPIのレスポンスデータ
var mockResponse = `{ 
	"name": "Tokyo",
 	"main": {
		 "temp": 20.0, 
		 "humidity": 50 
	}, 
	"weather": [ { 
		"description": "晴れ" } ] 
}`

// モックサーバーを作成する関数
func setupMockServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	})
	return httptest.NewServer(handler)
}

// `fetchWeather` のテスト
func TestFetchWeather(t *testing.T) {
	// モックサーバーを作成
	mockServer := setupMockServer()
	defer mockServer.Close()

	// 環境変数を変更（モックサーバーのURLに設定）
	os.Setenv("BASE_URL", mockServer.URL)
	os.Setenv("API_KEY", "bcaa618n4b046e69bdcff9d630b76c493") 

	// 関数を実行
	weather, err := fetchWeather("Tokyo")

	// エラーがないことを確認
	assert.NoError(t, err)

	// 期待通りのデータが取得できたか確認
	assert.Equal(t, "Tokyo", weather.Name)
	assert.Equal(t, 20.0, weather.Main.Temp)   
	assert.Equal(t, 50, weather.Main.Humidity) 
	assert.Equal(t, "晴れ", weather.Weather[0].Description)
}