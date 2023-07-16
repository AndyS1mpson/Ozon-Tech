build-all:
	cd checkout && GOOS=linux GOARCH=amd64 make build
	cd loms && GOOS=linux GOARCH=amd64 make build
	cd notifications && GOOS=linux GOARCH=amd64 make build

run-all: build-all
	sudo docker compose up --force-recreate --build
	#docker-compose up --force-recreate --build

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notifications && make precommit

run-log-env:
	docker-compose -f ./deployments/docker-compose.logs.yaml up --build

run-services:
	docker-compose -f ./deployments/docker-compose.yml up --build

run-services-for-mac:
	cd checkout && env GOOS=linux GOARCH=arm make build && cd .. && \
	cd loms && env GOOS=linux GOARCH=arm make build && cd .. && \
	cd notifications && env GOOS=linux GOARCH=arm make build && cd .. && \
	docker-compose -f ./deployments/docker-compose.yml up --build

stop:
	docker stop $(docker ps -a -q)
