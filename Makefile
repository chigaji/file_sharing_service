.PHONY: run stop exec_db migrate test 
run: 
	docker-compose up -d
	go run cmd/main.go

# stop the application (db for now)
stop:
	docker-compose down

# Open a psql shell
exec_db: 
	docker-compose exec db psql -U admin -d file_sharing_db
	 
migrate:
	go run cmd/main.go migrate

test:
	go test ./..

clean: 
	docker-compose down -v


