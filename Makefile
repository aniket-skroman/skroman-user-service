run:
	go run application.go 

liveserver:
	nodemon --exec go run application.go --signal SIGTERM

migratecreate:
	migrate create -ext sql -dir db/migrations/ -seq alter_unique_constraints

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:root@localhost:5432/skroman_users?sslmode=disable" --verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:root@localhost:5432/skroman_users?sslmode=disable" --verbose down

migratefix:
	migrate -path db/migrations/ -database postgres://postgres:root@localhost:5432/skroman_users?sslmode=disable force 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

PHONY:
	run, liveserver, migratecreate, migrateup, migratedown
