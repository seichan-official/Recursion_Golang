package models

type WeatherResponse struct{
	Name string `json:"name"`
	Main struct{
	Temp float64 `json:"temp"`
	Humidity int `json:"humidity"`
	} `json:"main"`
	Weather []struct{
	Description string `json:"description"`
	} `json:"weather"`
}