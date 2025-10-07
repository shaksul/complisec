package main

import (
	"database/sql"
	"fmt"
	"log"

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

	// Получаем ID ролей
	var adminRoleID, userRoleID, managerRoleID, userRoleID2 string

	err = db.QueryRow("SELECT id FROM roles WHERE name = 'Admin' LIMIT 1").Scan(&adminRoleID)
	if err != nil {
		log.Fatal("Ошибка получения Admin роли:", err)
	}

	err = db.QueryRow("SELECT id FROM roles WHERE name = 'User' LIMIT 1").Scan(&userRoleID)
	if err != nil {
		log.Fatal("Ошибка получения User роли:", err)
	}

	err = db.QueryRow("SELECT id FROM roles WHERE name = 'Manager' LIMIT 1").Scan(&managerRoleID)
	if err != nil {
		log.Fatal("Ошибка получения Manager роли:", err)
	}

	err = db.QueryRow("SELECT id FROM roles WHERE name = 'Пользователь' LIMIT 1").Scan(&userRoleID2)
	if err != nil {
		log.Fatal("Ошибка получения Пользователь роли:", err)
	}

	fmt.Printf("Admin ID: %s\n", adminRoleID)
	fmt.Printf("User ID: %s\n", userRoleID)
	fmt.Printf("Manager ID: %s\n", managerRoleID)
	fmt.Printf("Пользователь ID: %s\n", userRoleID2)

	// Получаем все права
	rows, err := db.Query("SELECT id, code FROM permissions ORDER BY code")
	if err != nil {
		log.Fatal("Ошибка получения прав:", err)
	}
	defer rows.Close()

	var permissions []struct {
		id   string
		code string
	}

	for rows.Next() {
		var id, code string
		if err := rows.Scan(&id, &code); err != nil {
			log.Fatal("Ошибка сканирования прав:", err)
		}
		permissions = append(permissions, struct {
			id   string
			code string
		}{id, code})
	}

	fmt.Printf("Найдено %d прав\n", len(permissions))

	// Назначаем права ролям
	// Admin - все права
	fmt.Println("Назначаем все права роли Admin...")
	for _, perm := range permissions {
		_, err := db.Exec(
			"INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			adminRoleID, perm.id,
		)
		if err != nil {
			log.Printf("Ошибка назначения права %s роли Admin: %v", perm.code, err)
		}
	}

	// User - базовые права просмотра
	fmt.Println("Назначаем базовые права роли User...")
	userPermissions := []string{
		"dashboard.view",
		"users.view",
		"asset.view",
		"risk.view",
		"incident.view",
		"document.read",
		"training.view",
		"compliance.view",
	}

	for _, permCode := range userPermissions {
		var permID string
		err := db.QueryRow("SELECT id FROM permissions WHERE code = $1", permCode).Scan(&permID)
		if err != nil {
			log.Printf("Право %s не найдено: %v", permCode, err)
			continue
		}

		_, err = db.Exec(
			"INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			userRoleID, permID,
		)
		if err != nil {
			log.Printf("Ошибка назначения права %s роли User: %v", permCode, err)
		}

		// Также назначаем роли "Пользователь"
		_, err = db.Exec(
			"INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			userRoleID2, permID,
		)
		if err != nil {
			log.Printf("Ошибка назначения права %s роли Пользователь: %v", permCode, err)
		}
	}

	// Manager - права создания и редактирования
	fmt.Println("Назначаем права роли Manager...")
	managerPermissions := []string{
		"dashboard.view",
		"users.view",
		"asset.view", "asset.create", "asset.edit",
		"risk.view", "risk.create", "risk.edit",
		"incident.view", "incident.create", "incident.edit",
		"document.read", "document.upload", "document.edit",
		"training.view", "training.create", "training.edit",
		"compliance.view", "compliance.manage",
	}

	for _, permCode := range managerPermissions {
		var permID string
		err := db.QueryRow("SELECT id FROM permissions WHERE code = $1", permCode).Scan(&permID)
		if err != nil {
			log.Printf("Право %s не найдено: %v", permCode, err)
			continue
		}

		_, err = db.Exec(
			"INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			managerRoleID, permID,
		)
		if err != nil {
			log.Printf("Ошибка назначения права %s роли Manager: %v", permCode, err)
		}
	}

	fmt.Println("✅ Связи ролей и прав восстановлены!")
}
