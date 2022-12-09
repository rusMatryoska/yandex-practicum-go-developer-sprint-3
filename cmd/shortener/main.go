package main

import (
	"flag"
	h "github.com/rusMatryoska/yandex-practicum-go-developer-sprint-2/internal/handlers"
	m "github.com/rusMatryoska/yandex-practicum-go-developer-sprint-2/internal/middleware"
	s "github.com/rusMatryoska/yandex-practicum-go-developer-sprint-2/internal/storage"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	var (
		server   = flag.String("a", os.Getenv("SERVER_ADDRESS"), "server address")
		baseURL  = flag.String("b", os.Getenv("BASE_URL"), "base URL")
		filePath = flag.String("f", os.Getenv("FILE_STORAGE_PATH"), "server address")
	)
	flag.Parse()

	if *server == "" || *baseURL == "" {
		*server = "localhost:8080"
		*baseURL = "http://" + *server + "/"
	}

	if len(strings.Split(*server, ":")) != 2 {
		log.Fatal("Need address in a form host:port")
	}

	if bu := *baseURL; bu[len(bu)-1:] != "/" {
		*baseURL = *baseURL + "/"
	}

	storageItem := &s.StorageStruct{
		ID:    1000,
		URLID: make(map[string]int),
		IDURL: make(map[int]string),
	}
	if *filePath == "" {
		log.Println("WARNING: saving is done to a buffer, not a file! \n" +
			"Therefore all saves will be lost upon restart!\n" +
			"To save to a file, set the correct values for FILE_STORAGE_PATH environment variables.")
	} else {
		if _, err := os.Stat(*filePath); os.IsNotExist(err) {
			m.CreateFile(*filePath)
		} else {
			targets := m.InitMapByJSON(*filePath)
			storageItem.MU.Lock()
			for _, t := range targets {
				storageItem.URLID[t.FullURL] = t.ShortenURL
				storageItem.IDURL[t.ShortenURL] = t.FullURL
				storageItem.ID = t.ShortenURL
				m.URLJSONList = append(m.URLJSONList, t)
				log.Println("url", t.FullURL, "added to storage, you can get access by shorten:", *baseURL+strconv.Itoa(t.ShortenURL))
			}
			storageItem.MU.Unlock()
		}
	}

	h.SetValues(*filePath, *baseURL, *server)

	err := http.ListenAndServe(":"+strings.Split(*server, ":")[1], h.NewRouter(s.StorageInterface(storageItem)))
	if err != nil {
		log.Fatal(err)
	}

}
