# Test Task - SitesResponseAPI
"Задача:
Написать программу, которая будет проверять список сайтов на доступность.

Раз в минуту нужно проверять доступны ли сайты из списка и засекать время доступа к ним.
Есть большое кол-во пользователей, которые хотят знать время доступа к сайтам.
У пользователей есть 3 варианта запросов (эндпойнта):
1. Получить время доступа к определенному сайту.
2. Получить имя сайта с минимальным временем доступа.
3. Получить имя сайта с максимальным временем доступа.

И есть администраторы, которые хотят получать статистику количества запросов пользователей по трем вышеперечисленным эндпойнтам."

## Installation & Run
```bash
# Download this project
go get github.com/Snegniy/testTaskResponseApi
```


```bash
# Build and Run
cd testTaskResponseApi
go build -o testapp
./testapp

# API Endpoint : http://127.0.0.1:8000
# Test JWT : eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjN9._FCTegWdDGIUkKF6vz8ikeiuUi9r0lOxurginXhY9fQ
```

## Structure
```
├── cmd
│   ├── app.go
├── internal
│   ├── config
│   │   ├── config.go   // инициализация конфигурации приложения 
│   ├── cronjob
│   │   ├── cronjob.go // обходчик сайтов
│   ├── handlers
│   │   ├── handlers.go // хэндлеры
│   │   ├── writeJSON.go // отправка ответа в формате JSON
│   ├── model
│   │   ├── model.go // структура хранения кэша с обработанными сайтами
│   ├── repository
│   │   ├── repository.go // инициализация кэша и извлечение данных из него
│   ├── service
│   │   ├── service.go // бизнес-логика
├── pkg
│   ├── graceful
│   │   ├── server.go  // запуск graceful сервера
│   ├── jwt
│   │   ├── jwtcheck.go // создание тестового токена для админов
│   ├── logger
│   │   ├── logger.go // инициализация логгера
├── config.yml  // конфигурационные установки по умолчанию
├── go.mod
├── sites.txt // список сайтов для проверки их ответа
```

## API

#### /{site url}
* `GET` : Получить время ответа к определенному сайту из списка

#### /min
* `GET` : Получить имя сайта с минимальным временем ответа

#### /max
* `GET` : Получить имя сайта с максимальным временем ответа

#### /stat/{site url}
* `GET` : Получить статистику запросов по определенному сайту (доступно без авторизации (по умолчанию) или через авторизацию JWT)

#### /stat/min
* `GET` : Получить статистику запросов по сайту с минимальным временем доступа (доступно без авторизации (по умолчанию) или через авторизацию JWT)

#### /stat/max
* `GET` : Получить статистику запросов по сайту с максимальным временем доступа (доступно без авторизации (по умолчанию) или через авторизацию JWT)

## Performance
wrk tests
```
Running 2m test @ http://127.0.0.1:8000/min
  12 threads and 1600 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.78ms    3.82ms  81.56ms   86.12%
    Req/Sec    37.15k    13.11k   90.06k    65.04%
  Latency Distribution
     50%  843.00us
     75%    4.41ms
     90%    7.98ms
     99%   16.42ms
  53237054 requests in 2.00m, 8.88GB read
  Socket errors: connect 587, read 0, write 0, timeout 0
Requests/sec: 443297.07
Transfer/sec:     75.69MB
```