# YouTube Audio Download API

Python FastAPI сервис для скачивания аудио из YouTube видео с использованием yt-dlp.

## Возможности

- Скачивание аудио из YouTube видео
- Поддержка форматов: MP3, M4A, WAV, OPUS
- Стриминг файлов напрямую клиенту (без сохранения на диск)
- Автоматическая конвертация в нужный формат
- Health check endpoints
- Swagger документация

## API Endpoints

### GET /
Health check endpoint
```bash
curl http://localhost:8001/
```

### GET /health
Health check endpoint
```bash
curl http://localhost:8001/health
```

### GET /download
Скачать аудио из YouTube видео

**Параметры:**
- `url` (обязательный) - URL YouTube видео
- `format` (опциональный) - Формат аудио (mp3, m4a, wav, opus). По умолчанию: mp3

**Примеры:**
```bash
# Скачать в MP3 (по умолчанию)
curl "http://localhost:8001/download?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ" -o audio.mp3

# Скачать в M4A
curl "http://localhost:8001/download?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ&format=m4a" -o audio.m4a

# С использованием wget
wget "http://localhost:8001/download?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ" -O audio.mp3
```

### POST /download
Альтернативный endpoint с POST методом

**Body:**
```json
{
  "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "format": "mp3"
}
```

**Пример:**
```bash
curl -X POST "http://localhost:8001/download" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ", "format": "mp3"}' \
  -o audio.mp3
```

## Локальная разработка

### Установка зависимостей
```bash
pip install -r requirements.txt
```

### Запуск
```bash
python main.py
```

Сервис будет доступен по адресу: http://localhost:8001

### Swagger документация
После запуска доступна по адресу: http://localhost:8001/docs

## Docker

### Сборка образа
```bash
docker build -t yt-dl-api .
```

### Запуск контейнера
```bash
docker run -p 8001:8001 yt-dl-api
```

## Docker Compose

Запуск через docker-compose (из корневой директории проекта):

```bash
# Запустить все сервисы
docker-compose up

# Запустить только yt-dl-api
docker-compose up yt-dl-api

# Запустить в фоновом режиме
docker-compose up -d

# Остановить сервисы
docker-compose down
```

## Технологии

- **FastAPI** - современный веб-фреймворк для построения API
- **yt-dlp** - инструмент для скачивания видео/аудио из YouTube
- **FFmpeg** - конвертация аудио форматов
- **Uvicorn** - ASGI сервер

## Примечания

- Сервис использует временные директории для обработки файлов
- Файлы автоматически удаляются после отправки клиенту
- FFmpeg необходим для конвертации аудио форматов
- Сервис работает на порту 8001
