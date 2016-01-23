
run:
	go run *.go

up:
	docker-compose up -d

stop:
	docker-compose stop

mobile:
	gomobile bind -target=ios -o todo/Item.framework github.com/zemirco/todo/item

.PHONY: run up stop mobile
