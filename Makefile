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
	mkdir -p vendor data

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

shell:
	docker-compose exec app /bin/bash

tests:
	docker-compose run --entrypoint test.sh

#############################
# Argument fix workaround
#############################
%:
	@: