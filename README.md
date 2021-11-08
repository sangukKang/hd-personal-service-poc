# API Server PoC

AWS S3 Upload, Download 를 테스트 해볼 수 있는 간단한 API 서버.

## Run

```
go run ./cmd/server
```

## Proto

### Generate Proto files

```
cd ./proto
go generate
```

### IDL

### REST API

| Method | Path         | Params |
| ------ | ------------ | ------ |
| POST   | /v1/upload   |        |
| GET    | /v1/download |        |
|        |              |        |





