
all: clean

clean: 
	find . -name 'goapi' -type f -exec rm -f {} \; 
	go clean

go-get:
	rm -rf src/github.com
	go get -v github.com/codegangsta/negroni
	go get -v github.com/gorilla/mux
	go get -v github.com/unrolled/render
	go get -v github.com/satori/go.uuid
	go get -v database/sql
	go get -v github.com/go-sql-driver/mysql

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
	docker build -t create_link_server .
	docker images

network-create:
	docker network create --driver bridge api_network

network-inspect:
	docker network inspect api_network

mysql-run:
	docker run -d --name mysql --network api_network -td -p 3306:3306 -e MYSQL_ROOT_PASSWORD=cmpe281 mysql:5.5

rabbitmq-run:
	docker run --name rabbitmq --network api_network \
			   -p 8080:15672 -p 4369:4369 -p 5672:5672 \
			   -d rabbitmq:3-management

docker-run:
	docker run -d --name create_link_server --network api_network -td -p 3001:3001 create_link_server


kong-database:
	docker run -d --name kong-database --network api_network -p 9042:9042 cassandra:2.2

kong-run:
	docker run -d --name kong \
	          --network api_network \
              -e "KONG_DATABASE=cassandra" \
              -e "KONG_CASSANDRA_CONTACT_POINTS=kong-database" \
              -e "KONG_PG_HOST=kong-database" \
              -p 8000:8000 \
              -p 8443:8443 \
              -p 8001:8001 \
              -p 7946:7946 \
              -p 7946:7946/udp \
              kong:0.9.9

docker-shell:
	docker exec -it create_link_server bash 

kong-shell:
	docker exec -it kong bash 

docker-network:
	docker network ls

docker-network-inspect:
	docker network inspect host

docker-network-delete:
	docker network rm api_network

docker-clean:
	docker stop mysql
	docker rm mysql
	docker stop rabbitmq
	docker rm rabbitmq
	docker stop create_link_server
	docker rm create_link_server
	docker rmi create_link_server

docker-ip:
	docker-machine ip

docker-ps:
	 docker ps --all --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t"

docker-ps-ports:
	 docker ps --all --format "table {{.Names}}\t{{.Ports}}\t"

test-ping:
	curl localhost:3001/ping

test-create-link:
	curl -X POST \
  	http://localhost:3000/createlink \
  	-H 'Content-Type: application/json' \
  	-d '{ \
  		"uri":"https://medium.com/@elliotchance/goroutines-and-channels-a-real-ly-simple-server-in-go-93ba49ff7c5c" }'


