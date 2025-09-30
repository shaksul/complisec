# RiskNexus - Система управления рисками и ИБ-документацией

Модульная система учёта рисков, активов, документов и инцидентов с поддержкой динамического управления ролями и правами, мультитенантности и ИИ-помощника.

## 🚀 Быстрый старт

### Требования
- Docker и Docker Compose
- Go 1.21+ (для разработки backend)
- Node.js 18+ (для разработки frontend)

### Запуск через Docker

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd CompliSec
```

2. Запустите все сервисы:
```bash
docker-compose up -d
```

3. Откройте приложение:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432

### Демо-аккаунт
- Email: admin@demo.local
- Пароль: admin123

## 📋 Функциональность

### Основные модули
- **Пользователи и роли** - Динамическое управление ролями и правами (RBAC 2.0)
- **Активы** - Инвентаризация IT-активов с журналами событий
- **Риски** - Матрица рисков, GAP-анализ, статусы жизненного цикла
- **Документы** - Версионирование, утверждение, ознакомление сотрудников
- **Инциденты** - Управление инцидентами с таймлайном расследования
- **Обучение** - Материалы, назначения, тесты и контроль прогресса
- **AI-помощник** - Анализ документов и подсказки по соответствию стандартам

### Архитектура
- **Backend**: Go + Fiber + PostgreSQL
- **Frontend**: React + TypeScript + Vite + Material-UI
- **База данных**: PostgreSQL с миграциями
- **Аутентификация**: JWT + RBAC
- **AI**: Модульная система провайдеров

## 🛠 Разработка

### Backend
```bash
cd apps/backend
go mod download
go run main.go
```

### Frontend
```bash
cd apps/frontend
npm install
npm run dev
```

### База данных
Миграции автоматически применяются при запуске через Docker Compose.

## 📁 Структура проекта

```
CompliSec/
├── apps/
│   ├── backend/          # Go + Fiber API
│   └── frontend/         # React + TypeScript
├── docs/                 # Документация
├── docker-compose.yml    # Docker конфигурация
└── README.md
```

## 🔧 Конфигурация

### Переменные окружения

**Backend:**
- `DATABASE_URL` - URL подключения к PostgreSQL
- `JWT_SECRET` - Секретный ключ для JWT
- `PORT` - Порт сервера (по умолчанию 8080)

**Frontend:**
- `VITE_API_URL` - URL API сервера

## 📊 API Endpoints

### Аутентификация
- `POST /api/auth/login` - Вход в систему
- `POST /api/auth/refresh` - Обновление токена

### Пользователи и роли
- `GET /api/users` - Список пользователей
- `POST /api/users` - Создание пользователя
- `GET /api/roles` - Список ролей
- `POST /api/roles` - Создание роли

### Основные модули
- `GET /api/assets` - Список активов
- `GET /api/risks` - Список рисков
- `GET /api/docs` - Список документов
- `GET /api/incidents` - Список инцидентов
- `GET /api/training/materials` - Материалы обучения

### AI
- `GET /api/ai/providers` - Список AI провайдеров
- `POST /api/ai/providers` - Создание AI провайдера
- `POST /api/ai/query` - Запрос к AI

## 🚀 Развертывание

### Production
1. Настройте переменные окружения
2. Соберите образы: `docker-compose build`
3. Запустите: `docker-compose up -d`

### Масштабирование
- Backend можно масштабировать горизонтально
- Frontend статичен, можно развернуть на CDN
- PostgreSQL рекомендуется использовать с репликацией

## 📝 Лицензия

Проект разработан для внутреннего использования.
