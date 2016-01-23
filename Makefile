
run:
	go run *.go

up:
	docker-compose up -d

stop:
	docker-compose stop

.PHONY: run up stop
