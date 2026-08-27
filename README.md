docker run command

```
 docker run --rm -p 8080:8080 --env-file .env -e LOCAL_SYSTEM_DB=false myapp
```

Run test cases

```
go run ./cmd/test_dashboard/main.go
```