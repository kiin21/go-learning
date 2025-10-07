# Chạy test và tạo coverage profile

```bash
go test -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
```

# Xem coverage trực tiếp trên terminal

```bash
go test -coverprofile=coverage.out && go tool cover -func=coverage.out
```


# Xem một test function cụ thể đã cover được những dòng code nào

```bash
# html report cho TestComputeChunk
go test -run TestComputeChunk -coverprofile=coverage_chunk.out && go tool cover -html=coverage_chunk.out -o coverage_chunk.html

# in terminal
go test -run TestComputeChunk -coverprofile=coverage_chunk.out && go tool cover -func=coverage_chunk.out
```