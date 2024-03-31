package models

type ConfigFile struct {
	Token                string `json:"token"`
	GoogleAPIKey         string `json:"googleAPIKey"`
	GoogleSearchEngineID string `json:"googleSearchEngineID"`
}

var (
	GoogleAPIKey         string
	GoogleSearchEngineID string
)
