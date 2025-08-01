# Локальная разработка

Эта папка содержит скрипты для автоматизации локальной настройки и очистки окружения проекта.

## Скрипты

### 📦 `setup.sh` - Установка и инициализация

Автоматический скрипт для установки и настройки локального окружения проекта.

#### Возможности:
- ✅ Проверка и установка Docker и Docker Compose
- ✅ Проверка и установка Go и Node.js
- ✅ Установка Python зависимостей
- ✅ Настройка переменных окружения
- ✅ Установка зависимостей приложений
- ✅ Создание MongoDB ключа
- ✅ Запуск проекта через Docker Compose
- ✅ Инициализация MongoDB replica set
- ✅ Тестирование приложений
- ✅ Проверка health checks

#### Использование:

```bash
# Базовая установка
./local-development/setup.sh

# Пропустить установку зависимостей
./local-development/setup.sh --skip-deps

# Пропустить инициализацию MongoDB
./local-development/setup.sh --skip-mongo

# Принудительная переустановка
./local-development/setup.sh --force

# Показать справку
./local-development/setup.sh --help
```

#### Что делает скрипт:

1. **Проверка окружения**
   - Проверяет наличие Docker, Docker Compose, Go, Node.js
   - Устанавливает недостающие компоненты

2. **Настройка проекта**
   - Создает файл `.env` из `env.example`
   - Устанавливает зависимости приложений
   - Создает MongoDB ключ

3. **Запуск системы**
   - Останавливает существующие контейнеры
   - Запускает проект через Docker Compose
   - Инициализирует MongoDB replica set

4. **Тестирование**
   - Проверяет health checks приложений
   - Тестирует подключения к MongoDB
   - Выводит статус системы

### 🧹 `teardown.sh` - Очистка окружения

Скрипт для полной очистки локального окружения проекта.

### 📊 `status.sh` - Проверка статуса

Скрипт для быстрой проверки статуса локального окружения проекта.

#### Возможности:
- ✅ Остановка и удаление контейнеров
- ✅ Удаление Docker образов
- ✅ Удаление Docker volumes
- ✅ Удаление Docker сетей
- ✅ Очистка Docker системы
- ✅ Удаление локальных файлов
- ✅ Очистка логов
- ✅ Полная очистка Docker данных (опционально)

#### Использование:

```bash
# Обычная очистка с подтверждением
./local-development/teardown.sh

# Принудительная очистка без подтверждения
./local-development/teardown.sh --force

# Полная очистка (включая Docker данные)
./local-development/teardown.sh --full

# Показать справку
./local-development/teardown.sh --help
```

#### Возможности:
- ✅ Проверка Docker и Docker Compose
- ✅ Проверка контейнеров и их статуса
- ✅ Проверка MongoDB replica set
- ✅ Проверка health checks приложений
- ✅ Проверка файлов и зависимостей
- ✅ Проверка логов и ресурсов
- ✅ Финальная сводка статуса

#### Использование:

```bash
# Полная проверка
./local-development/status.sh

# Быстрая проверка
./local-development/status.sh --quick

# Подробная проверка
./local-development/status.sh --verbose

# Показать справку
./local-development/status.sh --help
```

#### Что удаляет скрипт:

1. **Docker ресурсы**
   - Контейнеры проекта
   - Образы проекта
   - Volumes проекта
   - Сети проекта

2. **Локальные файлы**
   - Файл `.env`
   - MongoDB ключ
   - `node_modules` (опционально)
   - Go кэш (опционально)

3. **Системные ресурсы**
   - Неиспользуемые Docker ресурсы
   - Логи Docker

## Быстрый старт

### Первая установка:

```bash
# Сделать скрипты исполняемыми
chmod +x local-development/setup.sh
chmod +x local-development/teardown.sh

# Запустить установку
./local-development/setup.sh
```

### Повторная установка:

```bash
# Очистить окружение
./local-development/teardown.sh

# Установить заново
./local-development/setup.sh
```

### Проверка статуса:

```bash
# Быстрая проверка
./local-development/status.sh --quick

# Полная проверка
./local-development/status.sh

# Проверить контейнеры
docker ps

# Проверить логи
docker logs app-node
docker logs app-go

# Проверить MongoDB
docker exec mongo-0 mongo --eval "rs.status()"
```

## Решение проблем

### Проблемы с установкой:

1. **Docker не установлен**
   ```bash
   sudo apt update
   sudo apt install docker.io
   sudo systemctl start docker
   sudo usermod -aG docker $USER
   ```

2. **Проблемы с правами**
   ```bash
   sudo chown $USER:$USER -R .
   chmod +x local-development/*.sh
   ```

3. **Проблемы с портами**
   ```bash
   # Проверить занятые порты
   sudo netstat -tlnp | grep -E "(3000|8080|27030|27031|27032)"
   ```

### Проблемы с MongoDB:

1. **Replica set не инициализируется**
   ```bash
   docker exec mongo-0 mongo --eval "rs.initiate({_id: 'rs0', members: [{_id: 0, host: 'mongo-0:27017'}, {_id: 1, host: 'mongo-1:27017'}, {_id: 2, host: 'mongo-2:27017'}]})"
   ```

2. **Проблемы с ключом**
   ```bash
   openssl rand -base64 756 > mongo/mongo-keyfile
   sudo chmod 400 mongo/mongo-keyfile
   ```

### Проблемы с приложениями:

1. **Health checks не проходят**
   ```bash
   # Проверить логи
   docker logs app-node --tail 20
   docker logs app-go --tail 20
   
   # Перезапустить приложения
   docker-compose restart app-node app-go
   ```

2. **Проблемы с подключением к MongoDB**
   ```bash
   # Проверить переменные окружения
   docker exec app-node env | grep MONGO
   docker exec app-go env | grep MONGO
   ```

## Логи и отладка

### Просмотр логов:

```bash
# Логи всех контейнеров
docker-compose logs

# Логи конкретного контейнера
docker logs app-node
docker logs app-go
docker logs mongo-0

# Логи с последними строками
docker logs app-node --tail 50
```

### Отладка:

```bash
# Войти в контейнер
docker exec -it app-node bash
docker exec -it app-go bash
docker exec -it mongo-0 bash

# Проверить процессы
docker exec app-node ps aux
docker exec app-go ps aux
```

## Автоматизация

### Добавление в .bashrc:

```bash
# Добавить алиасы в ~/.bashrc
echo 'alias dev-setup="./local-development/setup.sh"' >> ~/.bashrc
echo 'alias dev-clean="./local-development/teardown.sh"' >> ~/.bashrc
echo 'alias dev-status="./local-development/status.sh"' >> ~/.bashrc
source ~/.bashrc
```

### Использование алиасов:

```bash
# Установка
dev-setup

# Очистка
dev-clean

# Очистка с подтверждением
dev-clean --force

# Проверка статуса
dev-status --quick
```

## Безопасность

### Важные замечания:

1. **Скрипт setup.sh требует sudo для установки пакетов**
2. **Скрипт teardown.sh удаляет все данные проекта**
3. **Полная очистка (--full) удаляет все Docker данные**
4. **Всегда делайте резервные копии перед очисткой**

### Рекомендации:

1. **Используйте виртуальную машину для разработки**
2. **Регулярно делайте резервные копии важных данных**
3. **Проверяйте логи перед очисткой**
4. **Используйте --force только при необходимости** 