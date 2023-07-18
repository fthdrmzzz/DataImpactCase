build-app:
	docker-compose -f ./docker/docker-compose.yaml build app

run:
	docker-compose -f ./docker/docker-compose.yaml up -d	

down:
	docker-compose -f ./docker/docker-compose.yaml down 	
