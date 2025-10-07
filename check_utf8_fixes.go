package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	fmt.Println("🔍 Проверка улучшений UTF-8 кодировки")
	fmt.Println("=====================================")

	// Подключение к базе данных
	db, err := sql.Open("pgx", "postgres://complisec:complisec123@localhost:5432/complisec?sslmode=disable")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// 1. Проверяем кодировку клиента
	fmt.Println("\n1. 🔒 Проверка кодировки клиента:")
	var clientEncoding string
	err = db.QueryRow("SELECT current_setting('client_encoding')").Scan(&clientEncoding)
	if err != nil {
		log.Fatal("Ошибка получения client_encoding:", err)
	}
	fmt.Printf("   client_encoding: %s", clientEncoding)
	if clientEncoding == "UTF8" {
		fmt.Println(" ✅")
	} else {
		fmt.Println(" ❌")
	}

	// 2. Проверяем кодировку базы данных
	fmt.Println("\n2. 🗄️ Проверка кодировки базы данных:")
	var dbEncoding string
	err = db.QueryRow("SELECT pg_encoding_to_char(encoding) FROM pg_database WHERE datname = 'complisec'").Scan(&dbEncoding)
	if err != nil {
		log.Fatal("Ошибка получения кодировки БД:", err)
	}
	fmt.Printf("   database encoding: %s", dbEncoding)
	if dbEncoding == "UTF8" {
		fmt.Println(" ✅")
	} else {
		fmt.Println(" ❌")
	}

	// 3. Проверяем права на наличие mojibake
	fmt.Println("\n3. 🧠 Проверка прав на mojibake:")
	rows, err := db.Query("SELECT code, description FROM permissions WHERE description IS NOT NULL LIMIT 10")
	if err != nil {
		log.Fatal("Ошибка запроса прав:", err)
	}
	defer rows.Close()

	corruptedCount := 0
	totalCount := 0

	for rows.Next() {
		var code, description string
		if err := rows.Scan(&code, &description); err != nil {
			log.Printf("Ошибка сканирования: %v", err)
			continue
		}

		totalCount++

		// Проверяем на mojibake
		if containsMojibake(description) {
			fmt.Printf("   ❌ %s: %q (содержит mojibake)\n", code, description)
			corruptedCount++
		} else {
			fmt.Printf("   ✅ %s: %q\n", code, description)
		}
	}

	if corruptedCount == 0 {
		fmt.Printf("\n🎉 Отлично! Все %d проверенных прав имеют корректную кодировку!\n", totalCount)
	} else {
		fmt.Printf("\n⚠️ Найдено %d из %d прав с проблемами кодировки\n", corruptedCount, totalCount)
	}

	// 4. Тестируем вставку новых данных
	fmt.Println("\n4. 🧪 Тест вставки UTF-8 данных:")
	testDescription := "Тестовая проверка UTF-8: ё, й, щ, э, №, §"

	// Вставляем тестовые данные
	_, err = db.Exec(
		"INSERT INTO permissions (id, code, module, description) VALUES (gen_random_uuid(), $1, $2, $3)",
		"test.utf8.verification", "Тестирование", testDescription,
	)
	if err != nil {
		log.Printf("Ошибка вставки тестовых данных: %v", err)
	} else {
		fmt.Printf("   ✅ Вставка успешна: %q\n", testDescription)
	}

	// Читаем обратно
	var readDescription string
	err = db.QueryRow("SELECT description FROM permissions WHERE code = 'test.utf8.verification'").Scan(&readDescription)
	if err != nil {
		log.Printf("Ошибка чтения тестовых данных: %v", err)
	} else {
		if readDescription == testDescription {
			fmt.Printf("   ✅ Чтение успешно: %q\n", readDescription)
		} else {
			fmt.Printf("   ❌ Данные повреждены при чтении: %q\n", readDescription)
		}
	}

	// Удаляем тестовые данные
	_, err = db.Exec("DELETE FROM permissions WHERE code = 'test.utf8.verification'")
	if err != nil {
		log.Printf("Ошибка удаления тестовых данных: %v", err)
	}

	fmt.Println("\n=====================================")
	fmt.Println("✅ Проверка завершена!")
}

// containsMojibake проверяет наличие признаков mojibake
func containsMojibake(text string) bool {
	mojibakePatterns := []string{
		"Рџ", "СЂ", "Рѕ", "СЃ", "Рј", "Рѕ", "С‚", "СЂ",
		"РІ", "Рѕ", "Р·", "Рј", "Рѕ", "Р¶", "РЅ", "Рѕ", "СЃ", "СЃ",
		"РЅ", "Р°", "СЃ", "С‚", "СЂ", "Рѕ", "Р№", "Рє",
	}

	for _, pattern := range mojibakePatterns {
		if len(pattern) <= len(text) {
			for i := 0; i <= len(text)-len(pattern); i++ {
				if text[i:i+len(pattern)] == pattern {
					return true
				}
			}
		}
	}

	return false
}
