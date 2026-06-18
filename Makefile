dev:
	docker compose -f dev.docker-compose.yaml up --build --force-recreate --remove-orphans

dev-down:
	docker compose -f dev.docker-compose.yaml down -v --remove-orphans

prod:
	docker compose -f docker-compose.yaml up --build -d

prod-down:
	docker compose -f docker-compose.yaml down -v --remove-orphans

build-backend-dev:
	docker build ./backend -f backend/Dockerfile.dev -t wdml_mtg_backend:dev

e2e-test:
	cd e2e; godotenv -f ./.env go test ./test
