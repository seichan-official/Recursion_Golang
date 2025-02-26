package repositories

import "strings"

// 都市リスト
var allCities = []string{
	"Tokyo", "Osaka", "Nagoya", "Fukuoka", "Sapporo",
	"New York", "Los Angeles", "Toronto", "Chicago",
	"London", "Paris", "Berlin", "Rome",
	"Shanghai", "Hong Kong", "Seoul", "Bangkok",
	"Sydney", "Dubai", "Cape Town", "Rio de Janeiro",
	"Dubai", "Cape Town",
}

// 都市やエリアを検索する関数
func SearchCities(query string) []string {
	var results []string
	query = strings.ToLower(query)
	for _, city := range allCities {
		if strings.Contains(strings.ToLower(city), query) {
			results = append(results, city)
		}
	}
	return results
}