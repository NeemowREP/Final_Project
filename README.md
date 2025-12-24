TODO Planner
Описание проекта

Это веб-приложение-планировщик задач. Пользователь может добавлять задачи, редактировать их, отмечать выполнение, удалять, а также использовать повторяющиеся задачи.
Проект реализован на Go с фронтендом на HTML/JS, база данных — SQLite.

Выполненные задания со звёздочкой:

- Поиск задач по заголовку, комментарию и дате
- Редактирование задач через API
- Отметка выполнения задач с пересчётом следующей даты для повторяющихся задач
- Аутентификация через JWT-токен
- Сборка и запуск через Docker


Запуск локально

1. Склонируйте репозиторий и перейдите в папку проекта:

git clone https://github.com/NeemowREP/Final_Project.git
cd Final_Project

2. Создайте файл .env с примером:

TODO_PORT=7540
LOG_LEVEL=info
READ_TIMEOUT=5s
WRITE_TIMEOUT=10s
WEB_DIR=./web
TODO_DBFILE=./scheduler.db
TODO_PASSWORD=12345

3. Запустите сервер:

go run main.go

4. Откройте браузер и перейдите на:

http://localhost:7540

Работа с API (curl-примеры)

- Получение токена (авторизация)
curl -X POST -H "Content-Type: application/json" \
  -d '{"password":"12345"}' \
  http://localhost:7540/api/signin

Тестирование

1. Получите токен через /api/signin и вставьте его в tests/settings.go:
var Token = "тут_вставить_значение_токена"

2. Запустите тесты:

go test ./tests

Сборка и запуск через Docker

1. Сборка Docker-образа:

docker build -t todo_app .

2. Запуск контейнера с пробросом порта и подключением базы данных:

docker run -it --rm -p 7540:7540 -v $(pwd)/scheduler.db:/app/scheduler.db \
  -e TODO_PASSWORD=12345 \
  todo_app

3. Перейдите в браузере по адресу:

http://localhost:7540