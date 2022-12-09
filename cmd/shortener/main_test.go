package main

import (
	h "github.com/rusMatryoska/yandex-practicum-go-developer-sprint-2/internal/handlers"
	s "github.com/rusMatryoska/yandex-practicum-go-developer-sprint-2/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (int, string) {
	t.Helper()
	r := strings.NewReader(body)
	req, err := http.NewRequest(method, ts.URL+path, r)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func TestRouter(t *testing.T) {
	storageItem := &s.StorageStruct{
		ID:    1000,
		URLID: make(map[string]int),
		IDURL: make(map[int]string),
	}

	h.SetValues("", "http://localhost:8080/", "localhost:8080")
	r := h.NewRouter(s.StorageInterface(storageItem))

	ts := httptest.NewServer(r)
	defer ts.Close()

	status, body := testRequest(t, ts, http.MethodGet, "/1001", "")
	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, "There is no url with this id\n", body)

	status, body = testRequest(t, ts, http.MethodGet, "/1111a", "")
	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "ID parameter must be Integer type\n", body)

	status, body = testRequest(t, ts, http.MethodPost, "/", "https://golang-blog.blogspot.com")
	assert.Equal(t, http.StatusCreated, status)
	assert.Equal(t, "http://localhost:8080/1001", body)

	status, _ = testRequest(t, ts, http.MethodGet, "/1001", "")
	assert.Equal(t, http.StatusOK, status)

	status, body = testRequest(t, ts, http.MethodPost, "/api/shorten", "{\"url\":\"https://e.mail.ru/inbox/23445\"}")
	assert.Equal(t, http.StatusCreated, status)
	assert.Equal(t, "{\"result\":\"http://localhost:8080/1002\"}\n", body)

	status, _ = testRequest(t, ts, http.MethodGet, "/1002", "")
	assert.Equal(t, http.StatusOK, status)

}
