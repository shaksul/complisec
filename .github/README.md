# GitHub Actions CI/CD

Этот проект использует GitHub Actions для автоматизации тестирования, сборки и деплоя.

## Workflows

### 1. CI Pipeline (`ci.yml`)

Запускается при каждом push и pull request в ветки `master`, `main`, `develop`.

**Включает:**
- **Backend Tests** - Go тесты с покрытием кода
- **Frontend Tests** - Линтинг и сборка фронтенда
- **Integration Tests** - Тесты с использованием Docker Compose
- **Security Scan** - Сканирование уязвимостей с Trivy
- **Code Quality** - Анализ кода через SonarCloud
- **Build & Push** - Сборка и публикация Docker образов (только для master)

### 2. E2E Tests (`e2e-tests.yml`)

Запускается при каждом push/PR и ежедневно по расписанию.

**Включает:**
- Запуск полного стека через Docker Compose
- End-to-End тесты с Playwright
- Smoke тесты основных страниц
- Генерация HTML отчетов

### 3. Deploy (`deploy.yml`)

Запускается при создании тега версии (`v*.*.*`) или вручную.

**Включает:**
- Сборка Docker образов с тегом версии
- Деплой на сервер через SSH
- Резервное копирование БД перед деплоем
- Проверка работоспособности после деплоя
- Уведомления в Slack
- Создание GitHub Release

## Настройка секретов

Для работы CI/CD необходимо настроить следующие секреты в GitHub:

### Обязательные для CI:
- `DOCKER_USERNAME` - логин Docker Hub
- `DOCKER_PASSWORD` - пароль или токен Docker Hub

### Для деплоя (опционально):
- `DEPLOY_HOST` - хост сервера для деплоя
- `DEPLOY_USER` - пользователь SSH
- `DEPLOY_SSH_KEY` - приватный SSH ключ
- `DEPLOY_PORT` - порт SSH (по умолчанию 22)

### Для уведомлений (опционально):
- `SLACK_WEBHOOK` - webhook URL для Slack уведомлений

### Для анализа кода (опционально):
- `SONAR_TOKEN` - токен SonarCloud

## Локальный запуск тестов

### Backend тесты:
```bash
cd apps/backend
go test -v ./...
```

### Frontend тесты:
```bash
cd apps/frontend
npm test
```

### Integration тесты:
```bash
docker-compose up -d
# Дождаться запуска сервисов
docker-compose exec backend go test -v ./internal/integration/...
docker-compose down
```

### E2E тесты:
```bash
docker-compose up -d
npm install -D @playwright/test
npx playwright install chromium
npx playwright test
docker-compose down
```

## Деплой

### Автоматический деплой по тегу:
```bash
git tag v1.0.0
git push origin v1.0.0
```

### Ручной деплой через GitHub UI:
1. Перейти в Actions → Deploy to Production
2. Нажать "Run workflow"
3. Выбрать environment (staging/production)
4. Нажать "Run workflow"

## Мониторинг

- **Test Reports**: доступны в Artifacts каждого workflow run
- **Coverage**: автоматически загружается в Codecov
- **Security**: результаты Trivy доступны в Security tab
- **Deployment Status**: уведомления приходят в Slack

## Структура тестов

```
.github/
├── workflows/
│   ├── ci.yml           # Основной CI pipeline
│   ├── e2e-tests.yml    # End-to-End тесты
│   └── deploy.yml       # Деплой
└── README.md            # Эта документация

apps/backend/
├── *_test.go            # Unit тесты
└── internal/integration/ # Integration тесты

e2e-tests/
├── playwright.config.ts
└── tests/
    └── *.spec.ts        # E2E тесты
```

## Troubleshooting

### Тесты падают с timeout
- Увеличьте `timeout-minutes` в workflow
- Проверьте логи сервисов: `docker-compose logs`

### Deployment failed
- Проверьте доступность сервера
- Убедитесь что SSH ключ добавлен в секреты
- Проверьте права пользователя на сервере

### Docker build failed
- Очистите кэш: `docker builder prune`
- Проверьте Dockerfile синтаксис
- Убедитесь что все зависимости доступны

## Best Practices

1. **Всегда создавайте feature branch** для новых изменений
2. **Дождитесь прохождения CI** перед merge в master
3. **Используйте semver** для версионирования (`v1.2.3`)
4. **Пишите тесты** для новой функциональности
5. **Проверяйте Security scan** перед деплоем

