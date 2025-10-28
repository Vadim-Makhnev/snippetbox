# SnippetBox 📝

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?style=for-the-badge&logo=mysql)
![GitHub](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)

**SnippetBox** — это веб-приложение для хранения и обмена текстовыми сниппетами кода, аналогичное Pastebin. Разработано на чистом Go с использованием лучших практик веб-разработки.

## ✨ Возможности

- 🚀 **Создание сниппетов** - Быстрое добавление кода с подсветкой синтаксиса
- 👁️ **Просмотр сниппетов** - Удобный просмотр с красивым форматированием
- ⏳ **Время жизни** - Настройка срока отображения сниппетов (1, 7, 365 дней)
- 👤 **Аутентификация** - Регистрация и вход пользователей
- 🔒 **Безопасность** - Хеширование паролей и защита от CSRF
- 📱 **Адаптивный дизайн** - Работает на всех устройствах

## Демо
<img width="2494" height="1323" alt="Screenshot 2025-10-27 at 23-45-25 Home - Snippetbox" src="https://github.com/user-attachments/assets/198a5882-c4c1-4439-aac6-3b3f8ba75277" />
<img width="2494" height="1323" alt="Screenshot 2025-10-27 at 23-45-36 Create a New Snippet - Snippetbox" src="https://github.com/user-attachments/assets/bbdbec6b-2d72-40d3-859e-b4ee6d83a80c" />
<img width="2494" height="1323" alt="Screenshot 2025-10-27 at 23-45-47 Signup - Snippetbox" src="https://github.com/user-attachments/assets/7463b435-00d4-4e63-b73c-e10fc2da59e7" />
<img width="2494" height="1323" alt="Screenshot 2025-10-27 at 23-45-53 Login - Snippetbox" src="https://github.com/user-attachments/assets/9bb98642-6b08-4f57-9b08-09235e06063b" />

## Swagger
<img width="2494" height="1002" alt="Screenshot 2025-10-27 at 23-44-53 Swagger UI" src="https://github.com/user-attachments/assets/8c2b420f-cd78-4d5a-b427-5e1e8ac4cb9d" />

## 🛠️ Технологический стек

### Backend
- **Язык**: Go 1.21+
- **HTTP сервер**: Standard library `net/http`
- **Шаблоны**: `html/template`
- **Валидация**: Кастомная система валидации
- **Сессии**: Secure cookie-based sessions

### База данных
- **СУБД**: MySQL 8.0+
- **Миграции**: `golang-migrate`
- **Драйвер**: `github.com/go-sql-driver/mysql`

### Тестирование
- **Фреймворк**: `testify/assert`
- **Покрытие**: Модульные и интеграционные тесты

## 🚀 Быстрый старт

### Предварительные требования

- Go 1.21 или выше
- MySQL 8.0+
- Git

### Установка и настройка
1. **Клонирование репозитория**
```bash
git clone https://github.com/Vadim-Makhnev/snippetbox.git
cd snippetbox
```
2. **Настройка базы данных и запуск миграций**
```sql
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'web'@'localhost' IDENTIFIED BY 'pass';
GRANT ALL PRIVILEGES ON snippetbox.* TO 'web'@'localhost';
```
```bash
migrate -path=./migrations -database="mysql://web:pass@tcp(localhost:3306)/snippetbox" up
```
3. **Запуск приложения**
```go
go run ./cmd/web
```


