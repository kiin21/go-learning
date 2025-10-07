# Chạy test và tạo coverage profile

```bash
go test -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
```

# Xem coverage trực tiếp trên terminal

```bash
go test -coverprofile=coverage.out && go tool cover -func=coverage.out
```
