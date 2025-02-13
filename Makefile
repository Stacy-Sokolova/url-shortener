build:
	docker-compose build app

run:
	docker-compose up app

migrate:
	migrate -path migrations -database '$(PG_URL_LOCALHOST)?sslmode=disable' up