# cmd/shortener

go run cmd/shortener/main.go -a localhost:33303 -b 1000 -f /home/victoria/Desktop/yandex_practicum_increments/storage/URL_STORAGE.json


# Обновление шаблона
https://github.com/Yandex-Practicum/go-autotests

cd cmd/shortener
go build -o shortener *.go
shortenertest -test.v -test.run=^TestIteration1$ -binary-path=shortener

/usr/local/shortenertest/shortenertest -test.v -test.run=^TestIteration1$ -binary-path=/home/victoria/Desktop/yandex_practicum_increments/yandex-practicum-go-developer-sprint-3/cmd/shortener/shortener
/usr/local/shortenertest/shortenertest -test.v -test.run=^TestIteration9$ -binary-path=/home/victoria/Desktop/yandex_practicum_increments/yandex-practicum-go-developer-sprint-3/cmd/shortener/shortener -source-path=.
Чтобы иметь возможность получать обновления автотестов и других частей шаблона выполните следующую команду:

```
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

# Endpoints

post:   
http://localhost:8080/
http://localhost:8080/api/shorten
http://localhost:8080/sign_in
http://localhost:8080/api/shorten/batch

get:    
http://localhost:8080/1001
http://localhost:8080/api/user/urls
http://localhost:8080/ping

#PostgreSQL
docker run --name habr-pg-13.3 -p 5432:5432 -e POSTGRES_USER=pguser -e POSTGRES_PASSWORD=pgpwd -e POSTGRES_DB=db -d postgres:13.3

url = postgresql://pguser:pgpwd@127.0.0.1:5432/db

//ctx, cancel := context.WithTimeout(storage.CTX, 5*time.Second)
//defer cancel()
//
//tx, err := storage.ConnPool.Begin(storage.CTX)
//if err != nil {
//	log.Fatal(err)
//}
//defer tx.Rollback(storage.CTX)