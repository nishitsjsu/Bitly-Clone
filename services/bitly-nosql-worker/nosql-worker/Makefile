
all: clean

clean: 
	find . -name 'goapi' -type f -exec rm -f {} \; 
	go clean

run:
	go run src/app/$(app).go

main:
	go run src/app/main.go

format:
	go fmt goapi

install:
	go install goapi

build:
	go build goapi

start:
	./goapi 

test-ping:
	curl localhost:3000/ping

rabbitmq-run:
	docker run --name rabbitmq --network api_network \
			   -p 8080:15672 -p 4369:4369 -p 5672:5672 \
			   -d rabbitmq:3-management

docker-build: 
	docker build -t nosql_worker .
	docker images

docker-run:
	docker run --network api_network \
			--name nosql_worker -p 3004:3004 -td nosql_worker
	docker ps

network-create:
	docker network create --driver bridge api_network

network-delete:
	docker network rm api_network

docker-network:
	docker network ls

docker-network-inspect:
	docker network inspect host

docker-shell:
	docker exec -it nosql_worker bash 

docker-clean:
	docker stop nosql_worker
	docker rm nosql_worker
	docker rmi nosql_worker

docker-ip:
	docker-machine ip

docker-ps:
	 docker ps --all --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t"

docker-ps-ports:
	 docker ps --all --format "table {{.Names}}\t{{.Ports}}\t"

