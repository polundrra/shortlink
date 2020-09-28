# shortlink

Сервис для сокращения ссылок. В качестве персистентного хранилища используется postgres.
Поднять - docker-compose up

### Api:
Генерация короткой ссылки:
```
POST /get-short-link example body: {"Url":"google.com", "CustomEnd": "foo"}
Response: {"Url":"foo"}
```
CustomEnd - опциональное поле для генерации собственной ссылки

Редирект по короткой ссылке:
```
GET /{code}
Response: 308 Location: google.com
```

Запуск тестов - go test ./... -cover Для интеграции с базой поднимется контейнер - должен умереть после выполнения тестов
