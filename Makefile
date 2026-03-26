.PHONY: test test-v build run docker-up docker-down lint clean

# запуск всех тестов
test:
	go test ./... -count=1

# тесты с подробным выводом
test-v:
	go test ./... -v -count=1

# тесты с покрытием
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

# собрать бинарь
build:
	go build -o bin/quickslot ./cmd/app

# запустить локально
run:
	go run ./cmd/app

# docker
docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

# очистка
clean:
	rm -rf bin/ coverage.out
