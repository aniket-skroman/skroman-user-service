run:
	go run application.go 

liveserver:
	nodemon --exec go run application.go --signal SIGTERM

migratecreate:
	migrate create -ext sql -dir db/migrations/ -seq init_skroman_client_contact_unique

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:support12@localhost:5432/skroman_users?sslmode=disable" --verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:root@localhost:5432/skroman_users?sslmode=disable" --verbose down

migratefix:
	migrate -path db/migrations/ -database postgres://postgres:root@localhost:5432/skroman_users?sslmode=disable force 2

migrateup_maindb:
	migrate -path db/migrations -database "postgresql://postgres:support12@skroman-user.ckwveljlsuux.ap-south-1.rds.amazonaws.com:5432/skroman_users" --verbose up

migratedown_maindb:
	migrate -path db/migrations -database "postgresql://postgres:support12@skroman-user.ckwveljlsuux.ap-south-1.rds.amazonaws.com:5432/skroman_users"  down

migratefix_maindb:
	migrate -path db/migrations -database "postgresql://postgres:support12@skroman-user.ckwveljlsuux.ap-south-1.rds.amazonaws.com:5432/skroman_users" force 4

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

PHONY:
	run, liveserver, migratecreate, migrateup, migratedown
