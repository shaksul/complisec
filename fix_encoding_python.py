#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import psycopg2
import uuid

# Подключение к базе данных
conn = psycopg2.connect(
    host="localhost",
    port="5432",
    database="complisec",
    user="complisec",
    password="complisec123"
)

cur = conn.cursor()

try:
    # Очищаем таблицы
    print("Очищаем таблицы...")
    cur.execute("DELETE FROM role_permissions")
    cur.execute("DELETE FROM permissions")
    
    # Вставляем правильные данные с UTF-8 кодировкой
    permissions_data = [
        # AI модуль
        ('ai.providers.view', 'ИИ', 'Просмотр провайдеров ИИ'),
        ('ai.providers.manage', 'ИИ', 'Управление провайдерами ИИ'),
        ('ai.queries.view', 'ИИ', 'Просмотр запросов ИИ'),
        ('ai.queries.create', 'ИИ', 'Создание запросов ИИ'),
        
        # Активы
        ('asset.view', 'Активы', 'Просмотр активов'),
        ('asset.create', 'Активы', 'Создание активов'),
        ('asset.edit', 'Активы', 'Редактирование активов'),
        ('asset.delete', 'Активы', 'Удаление активов'),
        ('asset.assign', 'Активы', 'Назначение активов'),
        
        # Документы
        ('document.read', 'Документы', 'Чтение документов'),
        ('document.upload', 'Документы', 'Загрузка документов'),
        ('document.edit', 'Документы', 'Редактирование документов'),
        ('document.delete', 'Документы', 'Удаление документов'),
        ('document.approve', 'Документы', 'Утверждение документов'),
        ('document.publish', 'Документы', 'Публикация документов'),
        
        # Риски
        ('risk.view', 'Риски', 'Просмотр рисков'),
        ('risk.create', 'Риски', 'Создание рисков'),
        ('risk.edit', 'Риски', 'Редактирование рисков'),
        ('risk.delete', 'Риски', 'Удаление рисков'),
        ('risk.assess', 'Риски', 'Оценка рисков'),
        ('risk.mitigate', 'Риски', 'Управление рисками'),
        
        # Инциденты
        ('incident.view', 'Инциденты', 'Просмотр инцидентов'),
        ('incident.create', 'Инциденты', 'Создание инцидентов'),
        ('incident.edit', 'Инциденты', 'Редактирование инцидентов'),
        ('incident.close', 'Инциденты', 'Закрытие инцидентов'),
        ('incident.assign', 'Инциденты', 'Назначение инцидентов'),
        
        # Обучение
        ('training.view', 'Обучение', 'Просмотр обучения'),
        ('training.assign', 'Обучение', 'Назначение обучения'),
        ('training.create', 'Обучение', 'Создание курсов'),
        ('training.edit', 'Обучение', 'Редактирование курсов'),
        ('training.view_progress', 'Обучение', 'Просмотр прогресса'),
        
        # Соответствие
        ('compliance.view', 'Соответствие', 'Просмотр соответствия'),
        ('compliance.manage', 'Соответствие', 'Управление соответствием'),
        ('compliance.audit', 'Соответствие', 'Проведение аудитов'),
        
        # Пользователи
        ('users.view', 'Пользователи', 'Просмотр пользователей'),
        ('users.create', 'Пользователи', 'Создание пользователей'),
        ('users.edit', 'Пользователи', 'Редактирование пользователей'),
        ('users.delete', 'Пользователи', 'Удаление пользователей'),
        ('users.manage', 'Пользователи', 'Управление пользователями'),
        
        # Роли
        ('roles.view', 'Роли', 'Просмотр ролей'),
        ('roles.create', 'Роли', 'Создание ролей'),
        ('roles.edit', 'Роли', 'Редактирование ролей'),
        ('roles.delete', 'Роли', 'Удаление ролей'),
        
        # Аудит
        ('audit.view', 'Аудит', 'Просмотр журнала аудита'),
        ('audit.export', 'Аудит', 'Экспорт журнала аудита'),
        
        # Дашборд
        ('dashboard.view', 'Дашборд', 'Просмотр дашборда'),
        ('dashboard.analytics', 'Дашборд', 'Просмотр аналитики'),
    ]
    
    print(f"Вставляем {len(permissions_data)} прав...")
    
    for code, module, description in permissions_data:
        permission_id = str(uuid.uuid4())
        cur.execute(
            "INSERT INTO permissions (id, code, module, description) VALUES (%s, %s, %s, %s)",
            (permission_id, code, module, description)
        )
        print(f"Добавлено: {code} - {description}")
    
    # Коммитим изменения
    conn.commit()
    print("✅ Все данные успешно добавлены!")
    
    # Проверяем результат
    print("\nПроверяем AI права:")
    cur.execute("SELECT code, description FROM permissions WHERE code LIKE 'ai.%' ORDER BY code")
    for row in cur.fetchall():
        print(f"  {row[0]}: {row[1]}")
        
except Exception as e:
    print(f"❌ Ошибка: {e}")
    conn.rollback()
finally:
    cur.close()
    conn.close()
