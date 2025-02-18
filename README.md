# Song Library API 🎶

Реализация онлайн-библиотеки песен с REST API, интеграцией с внешним API и документацией Swagger.

## Требования

Для работы с проектом вам понадобятся следующие инструменты:
- **Go** (версия 1.23 или выше)
- **Docker** (для запуска PostgreSQL и приложения)
- **PostgreSQL** (если вы хотите запустить базу данных локально)

## Структура проекта
```SongLibraryApi/
├── api/ # Обработчики HTTP-запросов
├── models/ # Модели данных (Song)
├── repositories/ # Работа с базой данных
├── services/ # Логика взаимодействия с внешним API
├── utils/ # Вспомогательные утилиты (логирование, конфигурация)
├── migrations/ # Миграции для PostgreSQL
├── .env.example # Пример конфигурационного файла
├── docker-compose.yml # Конфигурация Docker
├── Dockerfile # Dockerfile для сборки приложения
├── go.mod # Зависимости Go
├── main.go # Точка входа в приложение
└── README.md # Документация проекта
```


## Установка и запуск

### 1. Клонирование репозитория

Склонируйте репозиторий на ваш компьютер:

```bash
git clone https://github.com/itocode21/SongLibraryApi.git
cd SongLibraryApi
```
## Генерация документации Swagger

Для генерации документации Swagger выполните команду:

```bash
swag init --parseDependency --parseInternal


## 2. Настройка переменных окружения
Создайте файл .env на основе .env.example:
```bash
cp .env.example .env
```
Заполните файл .env своими данными. Пример:
```env
PORT=8080
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=music_library
DB_HOST=localhost
DB_PORT=5432
EXTERNAL_API_URL=http://localhost:9090/info
LOG_LEVEL=info
```

## 3. Запуск через Docker
Убедитесь, что Docker установлен и запущен. Запустите приложение с помощью ```docker-compose```:
```bash
docker-compose up --build
```
После запуска:
* Приложение будет доступно по адресу: ```http://localhost:8080```.
* Документация Swagger будет доступна по адресу: ```http://localhost:8080/swagger/index.html```.

## 4. Локальный запуск(без Docker)
Если хочется запустить приложение локальнео:
1. Убедитесь, что PostgreSQL запущен и доступен.
2. Примените миграции для создания таблиц:
```bash
go run main.go migrate-up
```
3. Запустите приложение:
```bash
go run main.go
```
## 5. Тестирование API

Вы можете протестировать API с помощью Swagger UI (```http://localhost:8080/swagger/index.html```) или Postman.
1. Получение списка песен с фильтрацией и пагинацией
Запрос:
```http
GET /api/v1/songs?group=Muse&song=Supermassive%20Black%20Hole&limit=10&offset=0
```
* Параметры запроса :
    * group (опционально): Фильтр по группе.
    * song (опционально): Фильтр по названию песни.
    * limit (опционально): Количество результатов на странице (по умолчанию: 10).
    * offset (опционально): Смещение для пагинации (по умолчанию: 0).
Пример ответа:
```json
[
  {
    "id": 1,
    "group_name": "Muse",
    "song_title": "Supermassive Black Hole",
    "release_date": "2006-07-16T00:00:00Z",
    "text": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?",
    "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
  }
]
```

2. Добавление новой песни
Запрос:
```http
POST /api/v1/songs
Content-Type: application/json
```
Тело запроса:
```json
{
  "group": "Muse",
  "song": "Supermassive Black Hole"
}
```
Пример успешного ответа:
```json
{
  "message": "Song created successfully"
}
```
Пример ошибки(если песня уже существует):
```json
{
  "error": "Song already exists"
}
```

3. Обновление данных песни
Запрос:
```http
PUT /api/v1/songs/1
Content-Type: application/json
```
Тело запроса:
```json
{
  "group_name": "Muse",
  "song_title": "Supermassive Black Hole (Updated)",
  "release_date": "2006-07-16T00:00:00Z",
  "text": "Updated lyrics...",
  "link": "https://www.youtube.com/watch?v=updatedLink"
}
```
Пример успешного ответа:
```json
{
  "message": "Song updated successfully"
}
```
Пример ошибки(если песня не найдена):
```json
{
  "error": "Song not found"
}
```

4. Удаление песни
Запрос:
```http
DELETE /api/v1/songs/1
```
Пример успешного ответа:
```json
{
  "message": "Song deleted successfully"
}
```
Пример ошибки(если песня не найдена):
```json
{
  "error": "Song not found"
}
```


5. Получение текста песни с пагинацией по куплетам
Запрос:
```http
GET /api/v1/songs/1/verses?limit=1&offset=0
```
Пример успешного ответа:
```json
{
  "message": "Song deleted successfully"
}
```
* Параметры запроса:
    * ```limit```(опционально): Количество куплетов на странице(default: 1).
    * ```offset```(опциональано): Смещение для пагинации (default: 0).
Пример ответа:
```json
{
  "verses": [
    "Ooh baby, don't you know I suffer?"
  ]
}
```
Пример ошибки(если песня не найдена):
```json
{
  "error": "Song not found"
}
```
6. Интеграция с внешним API
При добавлении новой песни выполняется запрос к внешнему API для получения 
дополнительной информации о песне. Например:
Запрос к внешнему API:
```http
GET /info?group=Muse&song=Supermassive%20Black%20Hole
```
Пример ответа от внешнего API:
```json
{
  "releaseDate": "16.07.2006",
  "text": "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?",
  "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
}
```
7. Ошибки
Если что-то пошло не так, API возвращает соответствующие HTTP-статусы и сообщения об ошибках. 
Примеры:
Неверный формат данных:
```json
{
  "error": "Invalid input"
}
```
Внутренняя ошибка сервера:
```json
{
  "error": "Internal server error"
}
```

### Так же методы можно посмотреть по адресу:
* ```http://localhost:8080/swagger/index.html```
* Выберите нужный метод и выполните запрос прямо из интерфейса.

