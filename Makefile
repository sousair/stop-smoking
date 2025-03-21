run:
	docker compose -f build/docker-compose.yml up 

rund:
	docker compose -d -f build/docker-compose.yml up 

down:
	docker compose -f build/docker-compose.yml down
