#############################
# Credit: https://github.com/webdevops/php-docker-boilerplate/blob/master/Makefile
#############################

ARGS = $(filter-out $@,$(MAKECMDGOALS))
MAKEFLAGS += --silent

list:
	sh -c "echo; $(MAKE) -p no_targets__ | awk -F':' '/^[a-zA-Z0-9][^\$$#\/\\t=]*:([^=]|$$)/ {split(\$$1,A,/ /);for(i in A)print A[i]}' | grep -v '__\$$' | grep -v 'Makefile'| sort"

#############################
# Docker machine states
#############################

init:
	mkdir -p vendor data/mysql data/redis data/rabbitmq

up:
	docker-compose up -d

start:
	docker-compose start

stop:
	docker-compose stop

state:
	docker-compose ps

rebuild:
	docker-compose stop
	docker-compose pull
	docker-compose rm --force app
	docker-compose build --no-cache --pull
	docker-compose up -d --force-recreate

#############################
# General
#############################
bash: shell

shell:
	docker-compose exec app /bin/bash

protobuf:
	docker-compose exec app /bin/bash -c "protoc -I/usr/local/include -I. -I/go/src -I/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. ./src/app/proto/visit_service.proto"
	docker-compose exec app /bin/bash -c "protoc -I/usr/local/include -I. -I/go/src -I/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. ./src/app/proto/visit_service.proto"
	docker-compose exec app /bin/bash -c "protoc -I/usr/local/include -I. -I/go/src -I/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. ./src/app/proto/visit_service.proto"
	docker-compose exec app /bin/bash -c "docker-compose exec app chown -R 1000:1000 ./src/app/proto"

sqlboiler:
	docker-compose exec app /bin/bash -c "sqlboiler --add-global-variants --add-panic-variants --wipe -d -c ./sqlboiler.yaml -o ./src/app/mysql/models -p models mysql"
	docker-compose exec app /bin/bash -c "chown -R 1000:1000 ./src/app/mysql/models"

grift:
	docker-compose exec app /bin/bash -c "cd src/ && grift $(filter-out $@,$(MAKECMDGOALS))"

tests:
	docker-compose run --entrypoint test.sh

#############################
# Argument fix workaround
#############################
%:
	@: