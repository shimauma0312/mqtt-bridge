# read-srv

```sh
# 起動 / 停止 / 作り直し
docker compose up -d --build --wait
docker compose down
docker compose down -v

# サンプル投入
docker compose exec -T mysql mysql -ubridge_writer -pbridge_writer bridge < mysql/sample/sample_records.sql

# DB
docker compose exec mysql mysql -ureader -preader bridge

# API
curl -s localhost:8080/api/devices
curl -s localhost:8080/api/devices/dev-001
curl -s "localhost:8080/api/devices/dev-001/track?from=2026-10-05T00:00:00Z&to=2026-10-06T00:00:00Z"
curl -s localhost:8080/api/devices/dev-001/trips
curl -s "localhost:8080/api/reports/daily?date=2026-10-05&tz=Asia/Tokyo"

# テスト
docker run --rm -v "$PWD/api":/src -w /src golang:1.23 sh -c 'go vet ./... && go test ./...'
```
