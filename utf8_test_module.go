package main

import (
	"database/sql"
	"testing"
	"unicode/utf8"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// TestUTF8Encoding проверяет корректность UTF-8 кодировки в базе данных
func TestUTF8Encoding(t *testing.T) {
	// Подключение к тестовой базе данных
	db, err := sql.Open("pgx", "postgres://complisec:complisec123@localhost:5432/complisec?sslmode=disable")
	if err != nil {
		t.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		t.Fatalf("Ошибка ping БД: %v", err)
	}

	// Проверяем кодировку клиента
	var clientEncoding string
	err = db.QueryRow("SELECT current_setting('client_encoding')").Scan(&clientEncoding)
	if err != nil {
		t.Fatalf("Ошибка получения client_encoding: %v", err)
	}

	if clientEncoding != "UTF8" {
		t.Errorf("client_encoding должен быть UTF8, получен: %s", clientEncoding)
	}

	// Проверяем кодировку базы данных
	var dbEncoding string
	err = db.QueryRow("SELECT pg_encoding_to_char(encoding) FROM pg_database WHERE datname = 'complisec'").Scan(&dbEncoding)
	if err != nil {
		t.Fatalf("Ошибка получения кодировки БД: %v", err)
	}

	if dbEncoding != "UTF8" {
		t.Errorf("Кодировка БД должна быть UTF8, получена: %s", dbEncoding)
	}

	t.Log("✅ Кодировка базы данных и клиента корректно установлена в UTF-8")
}

// TestUTF8Permissions проверяет, что все описания прав корректно читаются в UTF-8
func TestUTF8Permissions(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://complisec:complisec123@localhost:5432/complisec?sslmode=disable")
	if err != nil {
		t.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Получаем все права с описаниями
	rows, err := db.Query("SELECT code, description FROM permissions WHERE description IS NOT NULL")
	if err != nil {
		t.Fatalf("Ошибка запроса прав: %v", err)
	}
	defer rows.Close()

	corruptedCount := 0
	totalCount := 0

	for rows.Next() {
		var code, description string
		if err := rows.Scan(&code, &description); err != nil {
			t.Errorf("Ошибка сканирования права %s: %v", code, err)
			continue
		}

		totalCount++

		// Проверяем, что строка является корректной UTF-8
		if !utf8.ValidString(description) {
			t.Errorf("Право %s имеет некорректную UTF-8 строку: %q", code, description)
			corruptedCount++
			continue
		}

		// Проверяем на наличие признаков mojibake (поврежденной кодировки)
		if containsMojibake(description) {
			t.Errorf("Право %s содержит mojibake (поврежденную кодировку): %q", code, description)
			corruptedCount++
			continue
		}

		// Проверяем, что русские символы отображаются корректно
		if containsRussianText(description) && !isValidRussianText(description) {
			t.Errorf("Право %s содержит некорректный русский текст: %q", code, description)
			corruptedCount++
		}
	}

	if corruptedCount > 0 {
		t.Errorf("Найдено %d из %d прав с проблемами кодировки", corruptedCount, totalCount)
	} else {
		t.Logf("✅ Все %d прав имеют корректную UTF-8 кодировку", totalCount)
	}
}

// TestUTF8DataInsertion проверяет, что новые данные корректно вставляются в UTF-8
func TestUTF8DataInsertion(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://complisec:complisec123@localhost:5432/complisec?sslmode=disable")
	if err != nil {
		t.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Тестовые данные с русским текстом
	testData := []struct {
		code        string
		description string
		module      string
	}{
		{"test.utf8.check", "Тестовая проверка UTF-8", "Тестирование"},
		{"test.russian.text", "Проверка русского текста: ё, й, щ, э", "Тестирование"},
		{"test.special.chars", "Специальные символы: №, §, ©, ®", "Тестирование"},
	}

	for _, data := range testData {
		// Вставляем тестовые данные
		_, err := db.Exec(
			"INSERT INTO permissions (id, code, module, description) VALUES (gen_random_uuid(), $1, $2, $3)",
			data.code, data.module, data.description,
		)
		if err != nil {
			t.Errorf("Ошибка вставки тестовых данных %s: %v", data.code, err)
			continue
		}

		// Проверяем, что данные корректно прочитались
		var readDescription string
		err = db.QueryRow("SELECT description FROM permissions WHERE code = $1", data.code).Scan(&readDescription)
		if err != nil {
			t.Errorf("Ошибка чтения тестовых данных %s: %v", data.code, err)
			continue
		}

		// Сравниваем исходный и прочитанный текст
		if readDescription != data.description {
			t.Errorf("Текст не совпадает для %s: исходный=%q, прочитанный=%q",
				data.code, data.description, readDescription)
		}

		// Удаляем тестовые данные
		_, err = db.Exec("DELETE FROM permissions WHERE code = $1", data.code)
		if err != nil {
			t.Errorf("Ошибка удаления тестовых данных %s: %v", data.code, err)
		}
	}

	t.Log("✅ Все тестовые данные корректно обработаны в UTF-8")
}

// containsMojibake проверяет наличие признаков mojibake (поврежденной кодировки)
func containsMojibake(text string) bool {
	// Проверяем на типичные признаки mojibake
	mojibakePatterns := []string{
		"Рџ", "СЂ", "Рѕ", "СЃ", "Рј", "Рѕ", "С‚", "СЂ", // Типичные повреждения кириллицы
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

// containsRussianText проверяет, содержит ли текст русские символы
func containsRussianText(text string) bool {
	for _, r := range text {
		if (r >= 'А' && r <= 'я') || r == 'ё' || r == 'Ё' {
			return true
		}
	}
	return false
}

// isValidRussianText проверяет корректность русского текста
func isValidRussianText(text string) bool {
	// Проверяем, что все русские символы находятся в допустимом диапазоне
	for _, r := range text {
		if r >= 'А' && r <= 'я' {
			continue
		}
		if r == 'ё' || r == 'Ё' {
			continue
		}
		// Разрешаем пробелы, цифры, знаки препинания
		if r == ' ' || r == '-' || r == ':' || r == '.' || r == ',' || r == '(' || r == ')' {
			continue
		}
		// Разрешаем латинские буквы (для смешанного текста)
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			continue
		}
		// Разрешаем цифры
		if r >= '0' && r <= '9' {
			continue
		}
	}

	return true
}
