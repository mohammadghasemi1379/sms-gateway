# build the app
build-up:
	docker compose up -d --build

# run the app
up:
	docker compose up -d

# destroy and remove all volumes
destroy:
	docker compose down -v --remove-orphans

logs:
	docker-compose logs -f app