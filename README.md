# Практическая работа 13 — товары и PostgreSQL

## Краткий отчёт

Учет товаров на складе в PostgreSQL.

## Выполнение по заданию

### 1. База данных

Таблица `products` и пять товаров создаются скриптом `sql/init.sql`. Доступ к БД выполняется через `database/sql` и pgx.

### 2. Маршруты

GET и POST `/products`, GET, PUT и DELETE `/products/{id}`, фильтр `?category=Electronics`.

### 3. Проверка

Проверены список, фильтр, добавление, изменение цены и количества, удаление и ошибка отрицательной цены.

## Запуск

Локальный PostgreSQL установлен в `%LOCALAPPDATA%\Programs\postgres-labs`, порт 55432. База `inventory_db` должна быть создана; выполните скрипт sql/init.sql. Откройте PowerShell в этой папке:

```powershell
$env:DATABASE_URL = 'postgres://postgres:YOUR_POSTGRES_PASSWORD@127.0.0.1:55432/inventory_db?sslmode=disable'
go run .
```

Замените `YOUR_POSTGRES_PASSWORD` паролем своей базы (специальные символы в URL нужно закодировать). Откройте нужный адрес на `http://localhost:8080`. Если PostgreSQL остановлен, запустите его:

```powershell
& "$env:LOCALAPPDATA\Programs\postgres-labs\pgsql\bin\pg_ctl.exe" -D "$env:LOCALAPPDATA\Programs\postgres-labs\data" -l "$env:LOCALAPPDATA\Programs\postgres-labs\server.log" -o "-h 127.0.0.1 -p 55432" -w start
```

На другом компьютере создайте базу `inventory_db`, выполните `sql/init.sql` через `psql` и задайте свою строку подключения в `DATABASE_URL`.

## Пример добавления

Пока сервер работает, откройте второе окно PowerShell:

```powershell
$body = @{ name="Tablet"; category="Electronics"; price=100000; quantity=5 } | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri http://localhost:8080/products -ContentType "application/json; charset=utf-8" -Body ([System.Text.Encoding]::UTF8.GetBytes($body))
```

## Скриншот выполнения

![Результат работы приложения](docs/result.png)

## Вывод

Все записи хранятся в PostgreSQL, а запросы с параметрами передают значения отдельно от SQL.
