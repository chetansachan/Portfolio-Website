docker run command

```
docker build -t myapp .

 docker run --rm -p 8080:8080 --env-file .env -e LOCAL_SYSTEM_DB=false myapp
```

Run test cases

```
go run ./cmd/test_dashboard/main.go
```

to create database 

```
createdb -U postgres portfolio_website

migrate -path migration -database "postgres://postgres:chetan156@localhost:5432/portfolio_website?sslmode=disable" up
```