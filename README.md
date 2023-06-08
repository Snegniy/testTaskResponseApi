# Test Task - SitesResponseAPI
"Программа, которая будет проверять список сайтов на доступность.

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
go build
./testTaskResponseApi

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
│   ├── handlers
│   │   ├── writeJSON.go // отправка ответа в формате JSON
│   │   ├── handlers.go // хэндлеры
│   ├── middleware
│   │   ├── handlers.go // создание тестового токена для админов
│   ├── model
│   │   ├── model.go // структура хранения кэша с обработанными сайтами
│   ├── repository
│   │   ├── repository.go // инициализация кэша и извлечение данных из него
│   ├── responser
│   │   ├── response.go // обходчик сайтов
│   ├── service
│   │   ├── service.go // бизнес-логика
├── pkg
│   ├── graceful
│   │   ├── server.go  // запуск graceful сервера
│   ├── logger
│   │   ├── logger.go // инициализация логгера
├── config.yml  // конфигурационные установки по умолчанию
├── go.mod
├── sites.txt // список сайтов для проверки их ответа
```

## API

#### /url/{site url}
* `GET` : Получить время ответа к определенному сайту из списка

#### /min
* `GET` : Получить имя сайта с минимальным временем ответа

#### /max
* `GET` : Получить имя сайта с максимальным временем ответа

#### /stat/url{site url}
* `GET` : Получить статистику запросов по определенному сайту (доступно только через авторизацию JWT)

#### /stat/min
* `GET` : Получить статистику запросов по сайту с минимальным временем доступа (доступно только через авторизацию JWT)
* 
#### /stat/max
* `GET` : Получить статистику запросов по сайту с максимальным временем доступа (доступно только через авторизацию JWT)