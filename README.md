# KZ Salary Calculator

Веб-приложение для расчёта зарплаты в Казахстане с учётом всех налогов и отчислений.

## Стек

- **Backend:** Go
- **Database:** PostgreSQL
- **Frontend:** HTML + Tailwind CSS + Vanilla JS
- **Infrastructure:** Docker Compose

## Быстрый старт
```bash
git clone https://github.com/your-username/salary-calculator.git
cd salary-calculator

cp .env.example .env
# при необходимости отредактируй .env

docker compose up --build
```

Открой в браузере: **http://localhost:8080**

## Переменные окружения

| Переменная    | Описание                  | Пример       |
|---------------|---------------------------|--------------|
| `DB_HOST`     | Хост PostgreSQL            | `postgres`   |
| `DB_PORT`     | Порт PostgreSQL            | `5432`       |
| `DB_USER`     | Пользователь БД            | `postgres`   |
| `DB_PASSWORD` | Пароль БД                  | `postgres`   |
| `DB_NAME`     | Имя базы данных            | `salary_db`  |
| `MCI`         | МРП (МЗП) на текущий год  | `3932`       |
| `SERVER_PORT` | Порт сервера               | `8080`       |

## API

### Рассчитать зарплату
```
POST /api/calculate
Content-Type: application/json

{
  "salary": 500000,
  "mode": "gross"
}
```

| Поле     | Тип     | Описание                        |
|----------|---------|---------------------------------|
| `salary` | float64 | Сумма зарплаты                  |
| `mode`   | string  | `gross` — оклад, `net` — на руки |

### История расчётов
```
GET /api/history?limit=10&offset=0
```

### Текущий МРП
```
GET /api/mci
```

## Бизнес-логика

### Вычеты из зарплаты сотрудника

| Удержание | Ставка | База                              |
|-----------|--------|-----------------------------------|
| ОПВ       | 10%    | от оклада                         |
| ВОСМС     | 2%     | от оклада                         |
| ИПН       | 10%    | оклад − ОПВ − ВОСМС − 14 × МРП   |
| **NET**   |        | оклад − ОПВ − ВОСМС − ИПН        |

### Отчисления работодателя

| Отчисление | Ставка | База                              |
|------------|--------|-----------------------------------|
| СО         | 3.5%   | оклад − ОПВ                       |
| ООСМС      | 3%     | от оклада                         |
| СН         | 9.5%   | оклад − ОПВ − ВОСМС, минус СО    |

МРП 2025 — **3 932 ₸** (конфигурируется через `MCI` в `.env`)

## Структура проекта
```
├── cmd/server/
│   ├── main.go         # точка входа
│   └── app.go          # инициализация, роутер, graceful shutdown
├── internal/
│   ├── calculator/     # формулы расчёта
│   ├── config/         # загрузка переменных окружения
│   ├── handlers/       # HTTP handlers
│   ├── middleware/      # логирование запросов
│   ├── models/         # структуры данных
│   ├── repository/     # работа с БД
│   └── services/       # бизнес-логика
├── migrations/
│   └── 001_create_table.sql
├── web/
│   ├── index.html
│   ├── app.js
│   └── style.css
├── docker-compose.yml
├── Dockerfile
├── .env.example
└── CLAUDE.md
```
```

