package models

type WeatherAlert struct {
    SenderName  string `json:"sender_name"`
    Event       string `json:"event"`
    Start       int64  `json:"start"`
    End         int64  `json:"end"`
    Description string `json:"description"`
}

type WeatherResponse struct {
    Alerts []WeatherAlert `json:"alerts"`
}
