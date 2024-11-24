docker_build:
	sudo docker compose up --build

docker_up:
	sudo docker compose up -d

docker_down:
	sudo docker compose down

psql:
	sudo docker compose exec db psql -U postgres -d musiclib

createdb:
	sudo docker compose exec db createdb -U postgres musiclib

dropdb:
	sudo docker compose exec db dropdb -U postgres musiclib