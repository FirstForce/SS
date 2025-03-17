.PHONY: start-local-dependencies stop-local-dependencies

start-local-dependencies:
	docker-compose up -d --renew-anon-volumes

stop-local-dependencies:
	docker-compose down
	