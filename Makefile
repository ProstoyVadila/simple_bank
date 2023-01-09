postgres:
	@docker run --name simple_bank_db -p 5432:5432 -e POSTGRES_USER=admni -e POSTGRES_PASSWORD=stopmining -e POSTGRES_DB=data -d postgres:14-alpine
psql:
	@PGPASSWORD=stopmining psql -U admni -h localhost -p 5432 -d data
create_db:
	docker exec -it simple_bank_db createdb --username=admni data
drop_db:
	docker exec -it simple_bank_db dropdb --username=admni data

migrate_up:
	@migrate -path db/migrations -database 'postgresql://admni:stopmining@localhost:5432/data?sslmode=disable' -verbose up
mgirate_down:
	@migrate -path db/migrations -database 'postgresql://admni:stopmining@localhost:5432/data?sslmode=disable' -verbose down

gen_sqlc:
	sqlc generate

sqlc: gen_sqlc fieldalignment

tests:
	go test -v -cover ./...

fieldalignment:
	fieldalignment -fix ./... 

server:
	go run main.go

mocks:
	mockgen -build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/ProstoyVadila/simple_bank/db/sqlc Store


.PHONY: postgres createdb dropdb recreate_db psql sqlc migrate_up mgirate_down fieldalignment server gen_mocks gen_sqlc
