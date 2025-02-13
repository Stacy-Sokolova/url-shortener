build:
	docker-compose build app

run:
	docker-compose up app

migrate:
	migrate -path ./migrations -database 'postgres://postgres:Stacy@0.0.0.0:5432/postgres?sslmode=disable' up