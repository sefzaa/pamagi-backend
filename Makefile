.PHONY: docker-dev docker-down docker-logs migrate-up migrate-down swagger

# Menjalankan seluruh container (Go app dengan Air, MySQL, Redis)
docker-dev:
	docker-compose up --build -d

# Mematikan seluruh container
docker-down:
	docker-compose down

# Melihat logs dari aplikasi Go secara real-time
docker-logs:
	docker-compose logs -f app

# Menjalankan migrasi database ke atas
migrate-up:
	docker-compose exec app migrate -path database/migrations -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(mysql:3306)/${DB_NAME}" -verbose up

# Menurunkan migrasi database
migrate-down:
	docker-compose exec app migrate -path database/migrations -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(mysql:3306)/${DB_NAME}" -verbose down

# Mengenerate dokumentasi Swagger
swagger:
	swag init -g cmd/main.go