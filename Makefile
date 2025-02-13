build:
	docker-compose build app

run:
	STORAGE_TYPE=memdb docker-compose up app

migrate:
	migrate -path migrations -database '$(PG_URL_LOCALHOST)?sslmode=disable' up