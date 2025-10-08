# Настройка CI/CD для CompliSec

## Шаг 1: Настройка GitHub Secrets

### Обязательные секреты для Docker Hub:

1. Перейдите в **Settings** → **Secrets and variables** → **Actions**
2. Нажмите **New repository secret**
3. Добавьте следующие секреты:

| Имя секрета | Описание | Пример значения |
|-------------|----------|-----------------|
| `DOCKER_USERNAME` | Логин Docker Hub | `myusername` |
| `DOCKER_PASSWORD` | Токен Docker Hub | `dckr_pat_xxxxx` |

**Как получить Docker Hub токен:**
1. Войдите в [Docker Hub](https://hub.docker.com/)
2. Settings → Security → New Access Token
3. Скопируйте токен (он будет показан один раз!)

### Опциональные секреты для деплоя:

| Имя секрета | Описание |
|-------------|----------|
| `DEPLOY_HOST` | IP или домен сервера |
| `DEPLOY_USER` | SSH пользователь |
| `DEPLOY_SSH_KEY` | Приватный SSH ключ |
| `DEPLOY_PORT` | SSH порт (по умолчанию 22) |

### Опциональные секреты для уведомлений:

| Имя секрета | Описание |
|-------------|----------|
| `SLACK_WEBHOOK` | Webhook URL для Slack |

### Опциональные секреты для анализа кода:

| Имя секрета | Описание |
|-------------|----------|
| `SONAR_TOKEN` | Токен SonarCloud |

## Шаг 2: Настройка Docker Hub

### Создание репозиториев:

1. Войдите в Docker Hub
2. Создайте два репозитория:
   - `complisec-backend`
   - `complisec-frontend`
3. Сделайте их публичными или настройте доступ

## Шаг 3: Настройка сервера для деплоя (опционально)

### Подготовка сервера:

```bash
# 1. Установите Docker и Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo apt-get install docker-compose-plugin

# 2. Создайте пользователя для деплоя
sudo useradd -m -s /bin/bash deploy
sudo usermod -aG docker deploy

# 3. Создайте директорию для проекта
sudo mkdir -p /opt/complisec
sudo chown deploy:deploy /opt/complisec

# 4. Настройте SSH ключ
sudo -u deploy mkdir -p /home/deploy/.ssh
sudo -u deploy chmod 700 /home/deploy/.ssh
```

### Генерация SSH ключа:

```bash
# На вашей локальной машине
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/complisec_deploy

# Скопируйте публичный ключ на сервер
ssh-copy-id -i ~/.ssh/complisec_deploy.pub deploy@YOUR_SERVER_IP

# Добавьте приватный ключ в GitHub Secrets как DEPLOY_SSH_KEY
cat ~/.ssh/complisec_deploy
```

### Подготовка файлов на сервере:

```bash
# Войдите на сервер как deploy
ssh deploy@YOUR_SERVER_IP

cd /opt/complisec

# Создайте docker-compose.yml (базовую версию)
cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  backend:
    image: YOUR_DOCKERHUB_USERNAME/complisec-backend:latest
    ports:
      - "3001:8080"
    environment:
      - PORT=8080
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=complisec
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=complisec
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - postgres
    restart: always

  frontend:
    image: YOUR_DOCKERHUB_USERNAME/complisec-frontend:latest
    ports:
      - "3000:80"
    restart: always

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=complisec
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=complisec
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: always

volumes:
  postgres_data:
EOF

# Создайте .env файл
cat > .env << 'EOF'
DB_PASSWORD=YOUR_STRONG_PASSWORD
JWT_SECRET=YOUR_JWT_SECRET_AT_LEAST_32_CHARS
EOF

chmod 600 .env
```

## Шаг 4: Первый деплой

### Вариант 1: Через GitHub Actions (рекомендуется)

```bash
# На вашей локальной машине
git tag v1.0.0
git push origin v1.0.0
```

Деплой запустится автоматически!

### Вариант 2: Вручную

```bash
# На сервере
cd /opt/complisec
docker-compose pull
docker-compose up -d
```

## Шаг 5: Проверка работы CI/CD

### После push в master:

1. Перейдите в **Actions** в вашем GitHub репозитории
2. Найдите запущенный workflow **CI/CD Pipeline**
3. Проверьте, что все шаги проходят успешно:
   - ✅ Backend Tests
   - ✅ Frontend Tests
   - ✅ Integration Tests
   - ✅ Security Scan
   - ✅ Build & Push (только для master)

### После создания тега:

1. Создайте тег: `git tag v1.0.0 && git push origin v1.0.0`
2. В **Actions** найдите workflow **Deploy to Production**
3. Проверьте успешность деплоя
4. Проверьте что создан **Release** в GitHub

## Шаг 6: Настройка Slack уведомлений (опционально)

### Создание Incoming Webhook:

1. Перейдите в [Slack API](https://api.slack.com/apps)
2. Create New App → From scratch
3. Выберите workspace
4. Incoming Webhooks → Activate
5. Add New Webhook to Workspace
6. Выберите канал для уведомлений
7. Скопируйте Webhook URL
8. Добавьте в GitHub Secrets как `SLACK_WEBHOOK`

## Шаг 7: Настройка SonarCloud (опционально)

### Подключение проекта:

1. Перейдите на [SonarCloud](https://sonarcloud.io/)
2. Войдите через GitHub
3. Import Organization
4. Выберите ваш репозиторий
5. Generate Token
6. Добавьте токен в GitHub Secrets как `SONAR_TOKEN`
7. Создайте файл `sonar-project.properties` в корне проекта:

```properties
sonar.projectKey=complisec
sonar.organization=your-org

sonar.sources=apps/backend,apps/frontend/src
sonar.tests=apps/backend
sonar.test.inclusions=**/*_test.go

sonar.go.coverage.reportPaths=apps/backend/coverage.txt

sonar.exclusions=**/vendor/**,**/node_modules/**,**/*_test.go,**/dist/**
```

## Шаг 8: Мониторинг

### Проверка здоровья сервисов:

```bash
# Health check бэкенда
curl http://YOUR_SERVER:3001/health

# Detailed health check
curl http://YOUR_SERVER:3001/health/detailed

# Проверка фронтенда
curl http://YOUR_SERVER:3000
```

### Просмотр логов:

```bash
# На сервере
cd /opt/complisec

# Логи всех сервисов
docker-compose logs

# Логи бэкенда
docker-compose logs backend

# Логи с follow
docker-compose logs -f backend
```

### Резервное копирование БД:

```bash
# Создание бэкапа
docker-compose exec postgres pg_dump -U complisec complisec > backup_$(date +%Y%m%d_%H%M%S).sql

# Восстановление из бэкапа
docker-compose exec -T postgres psql -U complisec complisec < backup_YYYYMMDD_HHMMSS.sql
```

## Troubleshooting

### Workflow fails with "No space left on device"

**Решение:** Очистите Docker кэш в GitHub Actions

Добавьте шаг перед сборкой:
```yaml
- name: Clean Docker
  run: docker system prune -af
```

### SSH connection refused

**Проверьте:**
1. Правильность `DEPLOY_HOST` и `DEPLOY_PORT`
2. Что SSH ключ добавлен правильно
3. Что firewall разрешает SSH соединения
4. Что пользователь `deploy` существует

```bash
# Тест SSH с локальной машины
ssh -i ~/.ssh/complisec_deploy deploy@YOUR_SERVER_IP
```

### Docker pull fails on server

**Решение:** Проверьте что пользователь `deploy` в группе `docker`

```bash
sudo usermod -aG docker deploy
# Перелогиньтесь
```

### Health check fails after deployment

**Проверьте:**
1. Логи контейнеров: `docker-compose logs`
2. Что БД запустилась: `docker-compose ps postgres`
3. Что миграции прошли успешно
4. Переменные окружения в `.env`

```bash
# Проверка статуса
docker-compose ps

# Перезапуск сервисов
docker-compose restart
```

## Best Practices

1. **Всегда тестируйте локально** перед push в master
2. **Используйте feature branches** для разработки
3. **Создавайте теги** только для стабильных релизов
4. **Регулярно проверяйте Security scan** результаты
5. **Настройте автоматические бэкапы БД** через cron
6. **Мониторьте логи** и настройте алерты
7. **Документируйте изменения** в CHANGELOG.md

## Полезные команды

```bash
# Запуск workflow вручную
gh workflow run ci.yml

# Просмотр статуса workflows
gh run list

# Просмотр логов workflow
gh run view <run-id> --log

# Создание релиза
gh release create v1.0.0 --generate-notes

# Просмотр секретов (названия, не значения)
gh secret list

# Добавление секрета
gh secret set DOCKER_PASSWORD
```

## Дополнительные ресурсы

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Hub Documentation](https://docs.docker.com/docker-hub/)
- [Playwright Documentation](https://playwright.dev/)
- [SonarCloud Documentation](https://docs.sonarcloud.io/)

