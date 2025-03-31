.PHONY: start-local-dependencies stop-local-dependencies

build-app-backend:
	docker build -t app-backend .

start-local-dependencies:
	docker compose up -d --renew-anon-volumes

stop-local-dependencies:
	docker compose down
	