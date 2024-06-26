build-backend-dev:
	docker build ./backend -f backend/Dockerfile.dev -t wdml_mtg_backend:dev

run-backend-dev:
	docker run --rm -p 8080:8080 wdml_mtg_backend:dev

backend-dev: build-backend-dev run-backend-dev

dev:
	docker-compose -f dev.docker-compose.yaml up --build --force-recreate --remove-orphans

e2e-test:
	# cd backend; (E2E=true go run main.go &)
	# sleep 2
	cd e2e; godotenv -f ./.env go test ./test