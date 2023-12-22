# Launch.json

```
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Go",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/shortener",
            "env": {
                "SERVER_ADDRESS": ":7000",
                "DATABASE_DSN": "user=postgres password=postgres host=localhost port=5432 dbname=testDB sslmode=disable",
                "FILE_STORAGE_PATH": "./db.json",
                "BASE_URL": "localhost:7000/short/",
            },
            "args": [
                "-a=${env.SERVER_ADDRESS}",
                "-d=${env.DATABASE_DSN}",
                "-f=${env.FILE_STORAGE_PATH}",
                "-b=${env.BASE_URL}",
                "-l=debug"
            ]
        }
    ]
}
```