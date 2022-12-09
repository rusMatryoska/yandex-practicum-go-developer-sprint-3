package storage

import (
	"encoding/json"
	m "github.com/rusMatryoska/yandex-practicum-go-developer-sprint-2/internal/middleware"
	"log"
	"os"
	"strconv"
	"sync"
)

type StorageStruct struct {
	MU    sync.Mutex
	ID    int
	URLID map[string]int
	IDURL map[int]string
}

type StorageInterface interface {
	AddURL(url string, filePath string, baseURL string) (string, error)
	SearchURL(id int) string
}

func (storage *StorageStruct) AddURL(url string, filePath string, baseURL string) (string, error) {

	storage.MU.Lock()
	defer storage.MU.Unlock()

	if _, found := storage.URLID[url]; !found {

		storage.ID = storage.ID + 1
		storage.URLID[url] = storage.ID
		storage.IDURL[storage.ID] = url

		if filePath != "" {
			var URLSToWrite m.JSONStruct

			URLSToWrite.FullURL = url
			URLSToWrite.ShortenURL = storage.URLID[url]

			m.URLJSONList = append(m.URLJSONList, URLSToWrite)
			jsonString, err := json.Marshal(m.URLJSONList)
			if err != nil {
				return "", err
			}
			os.WriteFile(filePath, jsonString, 0644)
		}
		log.Println("url", url, "added to storage, you can get access by shorten:", baseURL+strconv.Itoa(storage.ID))
	}

	return baseURL + strconv.Itoa(storage.URLID[url]), nil
}

func (storage *StorageStruct) SearchURL(id int) string {
	storage.MU.Lock()
	defer storage.MU.Unlock()
	return storage.IDURL[id]
}
