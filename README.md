# S3 web interface written on go and vuejs

## Migrations

Creting new migration:

```bash
migrate create -ext sql -dir migrations -seq users_table
```

Apply migrations:

```bash
migrate -database 'postgres://almazilaletdinov@localhost:5432/web_s3?sslmode=disable' -path migrations up
```

## Contributing

[go-migrations guide](https://github.com/golang-migrate/migrate/blob/v4.18.1/GETTING_STARTED.md)
[sqlx guide](https://jmoiron.github.io/sqlx/)
