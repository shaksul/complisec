Проблемы и задачи

AuthService жёстко прошивает jwtSecret

Проблема: ключ берётся из константы, а не из конфигурации, поэтому нельзя менять между окружениями.

Ссылки:

Задача:

В auth_service.go добавить параметр jwtSecret string в NewAuthService.

В main.go передавать cfg.JWTSecret при создании AuthService.


Обновить вызовы AuthService (тесты, провайдеры).

Авторизация жёстко привязана к tenant 000...001

Проблема: DTO LoginRequest не содержит tenant_id, поэтому вход невозможен для других арендаторов.

Ссылки:



Задача:

В auth.go расширить LoginRequest, добавив TenantID string с валидацией.

В auth_handler.go читать tenant_id и передавать в authService.Login.

В auth_service.go принимать tenantID и прокидывать в userRepo.GetByEmail.

Обновить вызовы Login и тесты.

/auth/refresh принимает любой валидный JWT

Проблема: не проверяется claim type, можно обменять access-токен на refresh.

Ссылки:



Задача:

В auth_service.go добавить проверку claim type.

В auth_handler.go использовать проверку перед выдачей токенов.

Добавить тесты: access-токен не проходит refresh.

calculateSHA256 обрезает строку до 16 символов

Проблема: вызывает панику на коротких строках.

Ссылки:

Задача:

Использовать crypto/sha256 и возвращать полную строку.

Для усечённой формы проверять длину перед обрезкой.

Гонка в MemoryCache.Get

Проблема: удаление ключа происходит под RLock.

Ссылки:

Задача:

Сначала освободить RLock, потом взять Lock и удалить ключ.

Добавить конкурентный тест.

Методы /users/:id не реализованы

Проблема: возвращают заглушки.

Ссылки:

Задача:

В user_repo.go добавить методы для удаления/ролей.

В user_service.go реализовать DeleteUser, AssignRole, RemoveRole.

В user_handler.go вызывать сервис и возвращать результат.

Обновить тесты и документацию API.

AuthProvider не восстанавливает сессию

Проблема: после перезагрузки user = null, редиректит на логин.

Ссылки:

Задача:

При монтировании валидировать токен или запросить профиль.

Обновить isLoading и обработку ошибок.

При логине сохранять tenant-id/профиль.

Несовпадение DTO Document API (фронт ↔ бэкенд)

Проблема: фронт не содержит обязательных полей, запросы падают.

Ссылки:



Задача:

В documents.ts расширить DTO по document_dto.go.

Проверить использование (например, CreateDocumentWizard).

Обновить маппинг ответов (response.data.data).

Document API возвращает { data: ... }, фронт ждёт объект напрямую

Проблема: фронт получает undefined.

Ссылки:



Задача:

В documents.ts возвращать response.data.data.

Проверить компоненты (UploadNewVersionDialog, CreateDocumentWizard).

Users API возвращает разные структуры (фронт ↔ бэкенд)

Проблема: фронт ждёт response.data, бэкенд отдаёт { data, pagination }.

Ссылки:



Задача:

В users.ts возвращать response.data.data и response.data.pagination.

Обновить использование (например, CreateDocumentWizard.loadUsers).