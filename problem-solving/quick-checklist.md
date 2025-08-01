# Быстрый чек-лист критических проблем

## 🔴 КРИТИЧНО - Блокирует продакшн

### Безопасность
- [ ] **Hardcoded credentials** - Заменить на environment variables
  - [ ] app-go/read_products.go:18
  - [ ] app-node/create_product.js:4
  - [ ] scripts/create_app_user.py:9,14
- [ ] **Secret management** - Настроить централизованное управление
- [ ] **.env files** - Создать примеры конфигурации

### Инфраструктура
- [ ] **Docker Compose** - Создать единый файл запуска
- [ ] **Health checks** - Добавить проверки состояния
- [ ] **Logging** - Настроить структурированное логирование

### Качество кода
- [ ] **Tests** - Добавить unit/integration тесты
- [ ] **Error handling** - Улучшить обработку ошибок
- [ ] **CI/CD** - Настроить автоматизацию

### Документация
- [ ] **README.md** - Создать подробную документацию
- [ ] **Architecture docs** - Описать архитектуру системы

## 🟡 ВАЖНО - Влияет на качество

### Дополнительная безопасность
- [ ] **Input validation** - Добавить валидацию данных
- [ ] **Rate limiting** - Настроить защиту от DDoS

## 📊 Быстрая проверка

### Команды для проверки
```bash
# 1. Проверить hardcoded credentials
python3 scripts/check_passwords.py app-go/read_products.go app-node/create_product.js

# 2. Проверить docker-compose
docker-compose up -d
docker-compose ps

# 3. Проверить health checks
curl http://localhost:8080/health  # Go app
curl http://localhost:3000/health  # Node.js app

# 4. Проверить тесты
go test ./app-go/
npm test  # в app-node/

# 5. Проверить pre-commit hooks
pre-commit run --all-files
```

### Ожидаемые результаты
- [ ] Password detection: ✅ No hardcoded passwords detected
- [ ] Docker Compose: ✅ All services healthy
- [ ] Health checks: ✅ 200 OK responses
- [ ] Tests: ✅ All tests passing
- [ ] Pre-commit: ✅ All hooks passed

## 🚨 Критические файлы для изменения

### Приоритет 1 (КРИТИЧНО)
1. `app-go/read_products.go` - Заменить hardcoded URI
2. `app-node/create_product.js` - Заменить hardcoded URI
3. `scripts/create_app_user.py` - Заменить hardcoded passwords
4. `docker-compose.yml` - Создать в корне проекта

### Приоритет 2 (ВАЖНО)
1. `README.md` - Создать подробную документацию
2. `app-go/Dockerfile` - Добавить health checks
3. `app-node/Dockerfile` - Добавить health checks
4. `tests/` - Создать папку с тестами

## ⏱️ Временные рамки

### Сегодня (4-6 часов)
- [ ] Заменить hardcoded credentials
- [ ] Создать .env.example
- [ ] Добавить error handling

### На этой неделе (2-3 дня)
- [ ] Создать docker-compose.yml в корне
- [ ] Добавить health checks
- [ ] Написать базовые тесты
- [ ] Создать README.md

### В течение месяца
- [ ] Полная CI/CD pipeline
- [ ] Secret management
- [ ] Input validation
- [ ] Rate limiting
- [ ] Мониторинг

## 🎯 Критерии готовности к продакшну

### Безопасность
- [ ] 0 hardcoded credentials
- [ ] Environment variables используются
- [ ] Secret management настроен
- [ ] Security scans проходят

### Инфраструктура
- [ ] Один docker-compose запускает всё
- [ ] Health checks работают
- [ ] Logging настроен
- [ ] Мониторинг активен

### Качество
- [ ] Тесты проходят
- [ ] CI/CD работает
- [ ] Error handling есть
- [ ] Документация готова

## 📝 Заметки

### Команды для быстрого исправления
```bash
# 1. Создать .env.example
cp env.example .env.example

# 2. Заменить credentials в коде
# Использовать os.Getenv() в Go
# Использовать process.env в Node.js

# 3. Создать docker-compose.yml
# Объединить все сервисы в один файл

# 4. Добавить health checks
# /health endpoint в каждом приложении

# 5. Написать тесты
# go test для Go
# npm test для Node.js
```

### Полезные ссылки
- [Pre-commit hooks](https://pre-commit.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [MongoDB Security](https://docs.mongodb.com/manual/security/)
- [Go Testing](https://golang.org/pkg/testing/)
- [Node.js Testing](https://jestjs.io/) 