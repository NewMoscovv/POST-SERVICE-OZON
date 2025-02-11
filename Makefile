graph.gen:
	go run github.com/99designs/gqlgen generate --verbose

docker.rebuild:
	docker compose --env-file ./docker.env up -d --build app
docker.run:
	docker compose --env-file ./docker.env up -d
docker.run.db:
	docker compose --env-file ./docker.env up -d postgres
docker.run.migrate:
	docker compose --env-file ./docker.env up -d migrate
docker.down:
	docker compose down
migrate.up:
	migrate -path ./migrations -database "postgres://postgress-user:postgress-password@localhost:5432/pgdb?sslmode=disable" up
migrate.down:
	migrate -path ./migrations -database "postgres://postgress-user:postgress-password@localhost:5432/pgdb?sslmode=disable" down