package middleware

import (
	"encoding/json"
	"log"
	"os"
)

type MiddlewareStruct struct {
	FilePath string
	BaseURL  string
	Server   string
}

type JSONStruct struct {
	FullURL    string `json:"fullURL"`
	ShortenURL int    `json:"shortenURL"`
}

type URLFull struct {
	URLFull string `json:"url"`
}

type URLShorten struct {
	URLShorten string `json:"result"`
}

var URLJSONList []JSONStruct

func CreateFile(filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
}

func InitMapByJSON(filePath string) []JSONStruct {
	jsonString, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	targets := []JSONStruct{}

	err = json.Unmarshal(jsonString, &targets)
	if err != nil {
		log.Fatal(err)
	}
	return targets

}

func (m *MiddlewareStruct) InitMiddlewareStruct(filePath string, baseURL string, server string) {
	m.BaseURL = baseURL
	m.FilePath = filePath
	m.Server = server
}
