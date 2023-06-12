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
go get github.com/Snegniy/testTaskResponseApi/cmd
```


```bash
# Build and Run
cd github.com/Snegniy/testTaskResponseApi/cmd
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
  12 threads and 1999 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.57ms    3.80ms  66.03ms   86.23%
    Req/Sec    37.37k    16.29k  106.30k    66.13%
  Latency Distribution
     50%  722.00us
     75%    3.82ms
     90%    7.79ms
     99%   16.52ms
  53565794 requests in 2.00m, 9.08GB read
  Socket errors: connect 983, read 0, write 0, timeout 0
Requests/sec: 446018.74
Transfer/sec:     77.38MB
```