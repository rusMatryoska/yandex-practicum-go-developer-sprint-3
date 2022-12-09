package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	handlers "github.com/rusMatryoska/yandex-practicum-go-developer-sprint-3/internal/handlers"
	middleware "github.com/rusMatryoska/yandex-practicum-go-developer-sprint-3/internal/middleware"
	storage "github.com/rusMatryoska/yandex-practicum-go-developer-sprint-3/internal/storage"
)

func main() {
	var (
		st       storage.Storage
		err      error
		server   = flag.String("a", os.Getenv("SERVER_ADDRESS"), "server address")
		baseURL  = flag.String("b", os.Getenv("BASE_URL"), "base URL")
		filePath = flag.String("f", os.Getenv("FILE_STORAGE_PATH"), "server address")
		connStr  = flag.String("d", os.Getenv("DATABASE_DSN"), "connection url for DB")
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

	mwItem := &middleware.MiddlewareStruct{
		SecretKey: middleware.SecretKey,
		BaseURL:   *baseURL,
		Server:    *server,
	}

	if *connStr != "" {
		log.Println("WARNING: saving will be done through DataBase.")
		log.Println("connStr:", *connStr)

		DBItem := &storage.Database{
			BaseURL:   *baseURL,
			DBConnURL: *connStr,
			CTX:       context.Background(),
		}
		var dbErrorConnect error

		pool, err := DBItem.GetDBConnection()
		defer pool.Close()

		if err != nil {
			log.Println(err)
			dbErrorConnect = err
		}

		DBItem.ConnPool = pool
		DBItem.DBErrorConnect = dbErrorConnect

		if DBItem.DBErrorConnect == nil {
			if err := DBItem.CreateDBStructure(); err != nil {
				log.Fatal("unable to create db structure:", err)
			}
		}
		st = storage.Storage(DBItem)

	} else if *connStr == "" && *filePath != "" {
		log.Println("WARNING: saving will be done through file.")

		fileItem := &storage.File{
			BaseURL:  *baseURL,
			Filepath: *filePath,
			ID:       999,
			URLID:    make(map[string]int),
			IDURL:    make(map[int]string),
			UserURLs: make(map[string][]int),
		}

		if _, err := os.Stat(*filePath); os.IsNotExist(err) {
			middleware.CreateFile(*filePath)
		} else {
			targets := middleware.InitMapByJSON(*filePath)
			fileItem.NewFromFile(*baseURL, targets)
		}
		st = storage.Storage(fileItem)

	} else if *connStr == "" && *filePath == "" {
		log.Println("WARNING: saving will be done through memory.")
		memoryItem := &storage.Memory{
			BaseURL:  *baseURL,
			ID:       999,
			URLID:    make(map[string]int),
			IDURL:    make(map[int]string),
			UserURLs: make(map[string][]int),
		}

		st = storage.Storage(memoryItem)
	}

	if err = http.ListenAndServe(":"+strings.Split(*server, ":")[1],
		handlers.NewRouter(st, *mwItem)); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

}
