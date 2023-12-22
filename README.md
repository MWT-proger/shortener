# Shortener - Сервис сокращения URL



## Развертывание проекта

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.

```bash
git clone https://github.com/MWT-proger/shortener.git
```


2. Скопируйте шаблон файла с переменным окружения

```bash
  cp deployments/.env.example deployments/.env
```

3. Укажите верные переменные окружения в только что созданный файл [.env](deployments/.env)

*Доступны следующие переменные*
```bash
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=testDB
POSTGRES_PORT=5432
```
4. Запустите БД Postgres следующей командой

```bash
  docker compose -f deployments/docker-compose.yaml --env-file deployments/.env up -d
```

5. Запустите cервис сокращения URL

```
go run ./cmd/shortener -a "localhost:7000" -d "user=postgres password=postgres host=localhost port=5432 dbname=testDB sslmode=disable" -l debug
```

________________________________________________
- [Подробней по автотестам](docs/auto_tests.md)
- [launch.json для vscode](docs/vscode.md)