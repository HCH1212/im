package config

type JWT struct {
	AccessKey  string `json:"access_key"`
	RefreshKey string `json:"refresh_key"`
}
