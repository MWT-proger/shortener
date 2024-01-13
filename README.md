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

## Тестирование проекта

1. Benchmark 

```bash
go tool pprof -http=":9090" -seconds=30 http://localhost:7000/debug/pprof/profile 
```

```bash
go build ./cmd/
go tool pprof shortener -seconds=30 http://localhost:7000/debug/pprof/profile 
go test  -bench=. -cpuprofile=cpu.out -coverpkg=./../../...

go test -bench=. -memprofile=base.out
go tool pprof -http=":9090" bench.test base.out 
goimports -local "github.com/MWT-proger/shortener" -w main.go 
```


________________________________________________
- [Подробней по автотестам](docs/auto_tests.md)
- [launch.json для vscode](docs/vscode.md)

export PATH=$(go env GOPATH)/bin:$PATH
swag init -g internal/shortener/handlers/handlers.go -o ./docs

mwtech@mwtech-G3-3579:~/Projects/shortener$ export GOPATH="$HOME/go" 
mwtech@mwtech-G3-3579:~/Projects/shortener$ export PATH="$GOPATH/bin:$PATH"