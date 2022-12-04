# test-project

## Описание задачи

    Python/Go Lang developer test

Create a microservice that collect data from a cryptocompare using its API
for example (https://min-api.cryptocompare.com/data/pricemultifull?fs
y ms=BTC&tsyms=USD,EUR) For instance:
"fsyms" => ["BTC", "XRP", "ETH", "BCH", "EOS", "LTC", "XMR", "DASH"],
"tsyms" => ["USD", "EUR", "GBP", "JPY", "RUR"],

Struct of response currency

{
    CHANGE24HOUR string
    CHANGEPCT24HOUR string
    OPEN24HOUR string
    VOLUME24HOUR string
    VOLUME24HOURTO string
    LOW24HOUR string
    HIGH24HOUR string
    PRICE string
    SUPPLY string
    MKTCAP string
}

Currency pairs should be configurable.
Postgres parameters should be configurable.
Service must store data to postgres by sheduler (rawjson is ok).
Service must work in background.
If cryptocompare is non accessible service must return data from database via own API.
Data in response must be fresh (realtime). 2-3 minutes discrepancy is ok.

Using websockets is a plus. Clean architecture is a plus. Service scalability is a plus.

# Технические детали

-  Язык программирования: Golang 1.19
-  База данных: PostgreSQL 14.1


### Структура проекта

- `config` - содержит конфигурации сервиса, получаемые из переменных окружения.
- `pkg` - содержит вспомогательный функционал для работы .
- `Dockerfile` - содержит настройки для Docker.
- `go.mod` - содержит список всех используемых библиотек.
- `go.sum` - содержит список контрольных сумм всех используемых библиотек.
- `main.go` - содержит функию main(), которая запускается при старте сервиса и содержит в себе функционал для настройки и старта сервера.

### Зависимости

Список всех используемых библиотек можно найти в файле `go.mod`.


### Сборка

Для сборки проекта необходимо выполнить команду:

    docker-compose up
снести
    docker-compose down

### Запуск
запуск сервиса можно осуществить с помощью команды:
    go run main.go

### Swagger
для пересборки swagger необходимо выполнить команду:
    swag init
сам сваггер доступен по адресу:
    http://localhost:3000/swagger


### Использование сервиса по API или через сокет

Для получения данных необходимо отправить GET запрос по адресу:

    curl --location --request GET 'http://localhost:3000/service/price'

Для получения данных необходимо отправить запрос по сокету:

    ws://localhost:3000/service/ws/price

сообщение-запрос должен быть формата JSON

    {"tsyms":"BTC", "fsyms":"EUR"}
