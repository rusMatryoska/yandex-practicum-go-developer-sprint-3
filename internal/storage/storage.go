package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	middleware "github.com/rusMatryoska/yandex-practicum-go-developer-sprint-3/internal/middleware"
)

const (
	schema     = "public"
	table      = "storage"
	tableUsers = "users"
	sequence   = "id_serial"
)

type Storage interface {
	AddURL(url string, user string) (string, error)
	SearchURL(id int) (string, error)
	GetAllURLForUser(user string) ([]middleware.JSONStructForAuth, error)
	Ping() error
}

//MEMORY PART//

type Memory struct {
	BaseURL  string
	mu       sync.Mutex
	ID       int
	URLID    map[string]int
	IDURL    map[int]string
	UserURLs map[string][]int
}

func (m *Memory) AddURL(url string, user string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, found := m.URLID[url]; !found {

		m.ID = m.ID + 1
		m.URLID[url] = m.ID
		m.IDURL[m.ID] = url
		m.UserURLs[user] = append(m.UserURLs[user], m.ID)

		//log.Println("url", url, "added to storage, you can get access by shorten:", m.BaseURL+strconv.Itoa(m.ID),
		//	" by user ", user)
	}

	return m.BaseURL + strconv.Itoa(m.URLID[url]), nil
}

func (m *Memory) SearchURL(id int) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.IDURL[id] != "" {
		return m.IDURL[id], nil
	} else {
		return "", errors.New("No URL with this ID")
	}

}

func (m *Memory) GetAllURLForUser(user string) ([]middleware.JSONStructForAuth, error) {

	var (
		JSONStructList []middleware.JSONStructForAuth
		JSONStruct     middleware.JSONStructForAuth
	)

	if len(m.UserURLs[user]) == 0 {
		return JSONStructList, middleware.ErrNoContent
	} else {
		for i := range m.UserURLs[user] {
			JSONStruct.ShortURL = m.BaseURL + strconv.Itoa(m.UserURLs[user][i])
			JSONStruct.OriginalURL, _ = m.SearchURL(m.UserURLs[user][i])
			JSONStructList = append(JSONStructList, JSONStruct)

		}
		return JSONStructList, nil
	}
}

func (m *Memory) Ping() error {
	return errors.New("there is no connection to DB")
}

//FILE PART//

type File struct {
	BaseURL        string
	Filepath       string
	mu             sync.Mutex
	ID             int
	URLID          map[string]int
	IDURL          map[int]string
	UserURLs       map[string][]int
	URLSToWrite    middleware.JSONStruct
	JSONStructList []middleware.JSONStruct
}

func (f *File) NewFromFile(baseURL string, targets []middleware.JSONStruct) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.JSONStructList = targets

	for _, t := range targets {
		f.URLID[t.FullURL] = t.ShortenURL
		f.IDURL[t.ShortenURL] = t.FullURL
		f.UserURLs[t.User] = append(f.UserURLs[t.User], t.ShortenURL)
		f.ID = t.ShortenURL
		log.Println("url", t.FullURL, "added to storage, you can get access by shorten:", baseURL+strconv.Itoa(t.ShortenURL))
	}
}

func (f *File) AddURL(url string, user string) (string, error) {

	f.mu.Lock()
	defer f.mu.Unlock()

	log.Println("AddURL")
	log.Println(user)
	if _, found := f.URLID[url]; !found {

		f.ID = f.ID + 1
		f.URLID[url] = f.ID
		f.IDURL[f.ID] = url
		f.UserURLs[user] = append(f.UserURLs[user], f.ID)

		f.URLSToWrite.FullURL = url
		f.URLSToWrite.ShortenURL = f.URLID[url]
		f.URLSToWrite.User = user

		f.JSONStructList = append(f.JSONStructList, f.URLSToWrite)
		jsonString, err := json.Marshal(f.JSONStructList)
		if err != nil {
			return "", err
		}
		os.WriteFile(f.Filepath, jsonString, 0644)

		log.Println("url", url, "added to storage, you can get access by shorten:", f.BaseURL+strconv.Itoa(f.ID))
	}
	return f.BaseURL + strconv.Itoa(f.URLID[url]), nil
}

func (f *File) SearchURL(id int) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.IDURL[id], nil
}

func (f *File) GetAllURLForUser(user string) ([]middleware.JSONStructForAuth, error) {
	var (
		JSONStructList []middleware.JSONStructForAuth
		JSONStruct     middleware.JSONStructForAuth
	)

	log.Println(user)
	log.Println(f.UserURLs)
	if len(f.UserURLs[user]) == 0 {
		return JSONStructList, middleware.ErrNoContent
	} else {
		for i := range f.UserURLs[user] {
			JSONStruct.ShortURL = f.BaseURL + strconv.Itoa(f.UserURLs[user][i])
			JSONStruct.OriginalURL, _ = f.SearchURL(f.UserURLs[user][i])
			JSONStructList = append(JSONStructList, JSONStruct)

		}
		return JSONStructList, nil
	}
}

func (f *File) Ping() error {
	return errors.New("there is no connection to DB")
}

//DATABASE PART//

type Database struct {
	BaseURL        string
	DBConnURL      string
	CTX            context.Context
	ConnPool       *pgxpool.Pool
	DBErrorConnect error
}

func (db *Database) GetRows(query string) (pgx.Rows, error) {
	// ctx, cancel := context.WithTimeout(db.CTX, 5*time.Second)
	// defer cancel()

	rows, err := db.ConnPool.Query(db.CTX, query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (db *Database) Exec(query string) (pgconn.CommandTag, error) {
	res, err := db.ConnPool.Exec(db.CTX, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (db *Database) GetDBConnection() (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(db.CTX, db.DBConnURL)
	if err != nil {
		return nil, err
	} else {
		return pool, nil
	}
	// ..
}

func (db *Database) Ping() error {
	if db.DBErrorConnect != nil {
		return db.DBErrorConnect
	} else {
		err := db.ConnPool.Ping(db.CTX)
		return err
	}
}

func (db *Database) CreateDBStructure() error {
	tableExistsRes, err := db.GetRows(fmt.Sprintf("select true "+
		"from pg_catalog.pg_tables "+
		"where schemaname = '%s' and tablename = '%s'", schema, table))
	if err != nil {
		return err
	}
	tableExists := tableExistsRes.Next()

	tableUsersExistsRes, err := db.GetRows(fmt.Sprintf("select true "+
		"from pg_catalog.pg_tables "+
		"where schemaname = '%s' and tablename = '%s'", schema, tableUsers))
	if err != nil {
		return err
	}
	tableUsersExists := tableUsersExistsRes.Next()

	seqExistsRes, err := db.GetRows(fmt.Sprintf("select true "+
		"from pg_catalog.pg_sequences "+
		"where schemaname = '%s' and sequencename = '%s'", schema, sequence))
	if err != nil {
		return err
	}

	seqExists := seqExistsRes.Next()

	createTable := fmt.Sprintf("CREATE TABLE %s.%s"+
		"(\n  \"id\" SERIAL,"+
		"\n  \"full_url\" TEXT, "+
		"\n  \"user_id\" TEXT, "+
		"CONSTRAINT pk_storage\n\tPRIMARY KEY(id))", schema, table)

	createTableUsers := fmt.Sprintf("CREATE TABLE %s.%s"+
		"(\"user_id\" TEXT, CONSTRAINT pk_users\n\tPRIMARY KEY(user_id))", schema, tableUsers)

	createTableIndex := fmt.Sprintf("CREATE UNIQUE INDEX index_name ON %s.%s USING btree (full_url)", schema, table)

	createSequence := fmt.Sprintf("CREATE SEQUENCE %s.%s START ", schema, sequence)

	dropSequence := fmt.Sprintf("DROP SEQUENCE %s.%s", schema, sequence)

	if !tableUsersExists {
		db.Exec(createTableUsers)
	}
	if !tableExists {
		if seqExists {
			db.Exec(dropSequence)
		}
		db.Exec(createTable)
		db.Exec(createTableIndex)
		db.Exec(createSequence + strconv.Itoa(1000))
	} else if tableExists && !seqExists {
		row, err := db.GetRows(fmt.Sprintf("select max(id) from %s.%s", schema, table))
		if err != nil {
			return err
		}

		defer row.Close()
		row.Next()

		value, err := row.Values()
		if err != nil {
			return err
		}

		if value[0] == nil {
			db.Exec(createSequence + strconv.Itoa(1000))
		} else {
			initID := strconv.FormatInt(int64(value[0].(int32))+1, 10)
			db.Exec(createSequence + initID)
		}
	}
	log.Println("Structure for DB is ready.")
	return nil
}

func (db *Database) AddURL(url string, user string) (string, error) {
	var newID int64

	row := db.ConnPool.QueryRow(db.CTX, "INSERT INTO public.storage (full_url, user_id) VALUES ($1, $2) RETURNING id",
		url, user)
	if err := row.Scan(&newID); err != nil {
		return "", fmt.Errorf("318: row.Scan: %w", err)
	}
	return db.BaseURL + strconv.FormatInt(newID, 10), nil
}

func (db *Database) SearchURL(id int) (string, error) {
	var url string
	query := fmt.Sprintf("select full_url from %s.%s where id = %v", schema, table, id)
	row, err := db.GetRows(query)
	if err != nil {
		return "", err
	}
	defer row.Close()

	for row.Next() {
		value, err := row.Values()
		if err != nil {
			return "", err
		}

		if value[0] == nil {
			url = ""
		} else {
			url = value[0].(string)
		}
	}

	return url, nil

}

func (db *Database) GetAllURLForUser(user string) ([]middleware.JSONStructForAuth, error) {
	var (
		JSONStructList []middleware.JSONStructForAuth
		JSONStruct     middleware.JSONStructForAuth
		returnErr      error
		userID         string
	)

	queryUser := fmt.Sprintf("select user_id from %s.%s", schema, tableUsers)
	rowUser, err := db.GetRows(queryUser)
	if err != nil {
		return nil, err
	}
	defer rowUser.Close()
	//for rowUser.Next() {
	//	value, err := rowUser.Values()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	if fmt.Sprintf("%x", middleware.SetSign(value[0].(string), SecretKey)) == cookieUser {
	//		userID = value[0].(string)
	//	}
	//}

	query := fmt.Sprintf("select id, full_url from %s.%s where user_id = '%s'", schema, table, userID)
	row, err := db.GetRows(query)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if !row.Next() {
		returnErr = middleware.ErrNoContent
	} else {
		value, err := row.Values()
		if err != nil {
			return nil, err
		}

		JSONStruct.ShortURL = db.BaseURL + strconv.FormatInt(int64(value[0].(int32)), 10)
		JSONStruct.OriginalURL = value[1].(string)
		JSONStructList = append(JSONStructList, JSONStruct)
	}

	for row.Next() {
		value, err := row.Values()
		if err != nil {
			return nil, err
		}

		JSONStruct.ShortURL = db.BaseURL + strconv.FormatInt(int64(value[0].(int32)), 10)
		JSONStruct.OriginalURL = value[1].(string)
		JSONStructList = append(JSONStructList, JSONStruct)
	}
	fmt.Sprintln(JSONStructList)
	return JSONStructList, returnErr
}
