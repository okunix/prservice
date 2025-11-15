# Pull Request Service

Задание на программу стажировки avito.tech

Спецификация, описанная в `openapi.yaml`, реализована

## Сборка и запуск

Для сборки приложения вы можете использовать утилиту GNU Make.

```bash
make
```

Будет создан каталог `bin/` с исполняемым файлом prservice.

Для запуска приложения необходимо предоставить ему файл конфигурации. Путь поиска по умолчанию — `/etc/prservice/prservice.yaml`. Чтобы указать другое расположение, используйте флаг `--config`.

Пример:

```bash
./bin/prservice --config ./config.yaml
```

### Docker

Запуск приложения с docker-compose

**!** Атрибут `version` отсутствует в моих конфигах, поскольку он считается устаревшим в последних версиях docker-compose.

```bash
docker compose up
```

Чтобы пересобрать образ и запустить docker-compose, используйте

```bash
docker compose up --build
```

Чтобы использовать версию для разработки с горячей перезагрузкой, можно использовать конфигурацию `compose.dev.yaml`

```bash
docker compose -f compose.dev.yaml up
```

### Kubernetes

Также добавлен Kubernetes манифест для деплоя приложения. Протестировано на minikube.

```bash
kubectl apply -f manifest.yaml
```

## Конфигурация

Пример содержимого файла конфигурации:

```yaml
postgres:
  addr: localhost:5432 # хост и порт субд
  user: postgres # имя пользователя субд
  password: postgres # пароль пользователя субд
  db: prservice_db # название базы данных в postgres
  sslMode: disable # режим sslmode postgres
adminToken: secret_admin_token # Токен администратора, который используется для деактивации пользователей
addr: 0.0.0.0:80 # Адрес и порт на котором будет прослушивать сервер
```

## Дополнительные задания

Добавлен эндпоинт `POST /team/deactivate`, который "деактивирует" пользователей команды

Он требует авторизации с помощью токена админа в http заголовке: `Authorization: Bearer <admin_token>`

Пример тела запроса:

```json
{
  "team_name": "string"
}
```

Пример ответа:

```json
{
  "team": {
    "team_name": "devops",
    "members": [
      {
        "user_id": "u1",
        "username": "Alice",
        "is_active": false
      },
      {
        "user_id": "u2",
        "username": "Bob",
        "is_active": false
      },
      {
        "user_id": "u3",
        "username": "Claude",
        "is_active": false
      }
    ]
  }
}
```
