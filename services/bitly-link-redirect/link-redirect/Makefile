
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

docker-build: 
	docker build -t link_redirect_server .
	docker images

docker-run:
	docker run --network api_network \
			--name link_redirect_server -p 3002:3002 -td link_redirect_server
	docker ps

rabbitmq-run:
	docker run --name rabbitmq --network api_network \
			   -p 8080:15672 -p 4369:4369 -p 5672:5672 \
			   -d rabbitmq:3-management

network-create:
	docker network create --driver bridge api_network

network-delete:
	docker network rm api_network

docker-network:
	docker network ls

docker-network-inspect:
	docker network inspect host

docker-shell:
	docker exec -it link_redirect_server bash 

docker-clean:
	docker stop rabbitmq
	docker rm rabbitmq
	docker stop link_redirect_server
	docker rm link_redirect_server
	docker rmi link_redirect_server

docker-ip:
	docker-machine ip

docker-ps:
	 docker ps --all --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t"

docker-ps-ports:
	 docker ps --all --format "table {{.Names}}\t{{.Ports}}\t"

test-ping:
	curl localhost:3002/ping

test-gumball:
	curl localhost:3002/e4259
