package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Подключение к базе данных
	db, err := sql.Open("pgx", "postgres://complisec:complisec123@localhost:5432/complisec?sslmode=disable")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Fatal("Ошибка ping БД:", err)
	}

	// Очищаем таблицы
	fmt.Println("Очищаем таблицы...")
	_, err = db.Exec("DELETE FROM role_permissions")
	if err != nil {
		log.Fatal("Ошибка удаления role_permissions:", err)
	}
	_, err = db.Exec("DELETE FROM permissions")
	if err != nil {
		log.Fatal("Ошибка удаления permissions:", err)
	}

	// Данные для вставки с правильной UTF-8 кодировкой
	permissions := []struct {
		code        string
		module      string
		description string
	}{
		// AI модуль
		{"ai.providers.view", "ИИ", "Просмотр провайдеров ИИ"},
		{"ai.providers.manage", "ИИ", "Управление провайдерами ИИ"},
		{"ai.queries.view", "ИИ", "Просмотр запросов ИИ"},
		{"ai.queries.create", "ИИ", "Создание запросов ИИ"},

		// Активы
		{"asset.view", "Активы", "Просмотр активов"},
		{"asset.create", "Активы", "Создание активов"},
		{"asset.edit", "Активы", "Редактирование активов"},
		{"asset.delete", "Активы", "Удаление активов"},
		{"asset.assign", "Активы", "Назначение активов"},

		// Документы
		{"document.read", "Документы", "Чтение документов"},
		{"document.upload", "Документы", "Загрузка документов"},
		{"document.edit", "Документы", "Редактирование документов"},
		{"document.delete", "Документы", "Удаление документов"},
		{"document.approve", "Документы", "Утверждение документов"},
		{"document.publish", "Документы", "Публикация документов"},

		// Риски
		{"risk.view", "Риски", "Просмотр рисков"},
		{"risk.create", "Риски", "Создание рисков"},
		{"risk.edit", "Риски", "Редактирование рисков"},
		{"risk.delete", "Риски", "Удаление рисков"},
		{"risk.assess", "Риски", "Оценка рисков"},
		{"risk.mitigate", "Риски", "Управление рисками"},

		// Инциденты
		{"incident.view", "Инциденты", "Просмотр инцидентов"},
		{"incident.create", "Инциденты", "Создание инцидентов"},
		{"incident.edit", "Инциденты", "Редактирование инцидентов"},
		{"incident.close", "Инциденты", "Закрытие инцидентов"},
		{"incident.assign", "Инциденты", "Назначение инцидентов"},

		// Обучение
		{"training.view", "Обучение", "Просмотр обучения"},
		{"training.assign", "Обучение", "Назначение обучения"},
		{"training.create", "Обучение", "Создание курсов"},
		{"training.edit", "Обучение", "Редактирование курсов"},
		{"training.view_progress", "Обучение", "Просмотр прогресса"},

		// Соответствие
		{"compliance.view", "Соответствие", "Просмотр соответствия"},
		{"compliance.manage", "Соответствие", "Управление соответствием"},
		{"compliance.audit", "Соответствие", "Проведение аудитов"},

		// Пользователи
		{"users.view", "Пользователи", "Просмотр пользователей"},
		{"users.create", "Пользователи", "Создание пользователей"},
		{"users.edit", "Пользователи", "Редактирование пользователей"},
		{"users.delete", "Пользователи", "Удаление пользователей"},
		{"users.manage", "Пользователи", "Управление пользователями"},

		// Роли
		{"roles.view", "Роли", "Просмотр ролей"},
		{"roles.create", "Роли", "Создание ролей"},
		{"roles.edit", "Роли", "Редактирование ролей"},
		{"roles.delete", "Роли", "Удаление ролей"},

		// Аудит
		{"audit.view", "Аудит", "Просмотр журнала аудита"},
		{"audit.export", "Аудит", "Экспорт журнала аудита"},

		// Дашборд
		{"dashboard.view", "Дашборд", "Просмотр дашборда"},
		{"dashboard.analytics", "Дашборд", "Просмотр аналитики"},
	}

	fmt.Printf("Вставляем %d прав...\n", len(permissions))

	// Вставляем данные
	for _, perm := range permissions {
		id := uuid.New().String()
		_, err := db.Exec(
			"INSERT INTO permissions (id, code, module, description) VALUES ($1, $2, $3, $4)",
			id, perm.code, perm.module, perm.description,
		)
		if err != nil {
			log.Fatal("Ошибка вставки:", err)
		}
		fmt.Printf("Добавлено: %s - %s\n", perm.code, perm.description)
	}

	fmt.Println("✅ Все данные успешно добавлены!")

	// Проверяем результат
	fmt.Println("\nПроверяем AI права:")
	rows, err := db.Query("SELECT code, description FROM permissions WHERE code LIKE 'ai.%' ORDER BY code")
	if err != nil {
		log.Fatal("Ошибка запроса:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var code, description string
		if err := rows.Scan(&code, &description); err != nil {
			log.Fatal("Ошибка сканирования:", err)
		}
		fmt.Printf("  %s: %s\n", code, description)
	}
}
