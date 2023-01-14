include dev.env

scheme?=new_scheme
package_name=github.com/ProstoyVadila/simple_bank


check_env:
	echo $(DB_SOURCE)

run_postgres:
	@docker run --name simple_bank_db -p 5432:$(PGPORT) -e POSTGRES_USER=$(PGUSER) -e POSTGRES_PASSWORD=$(PGPASSWORD) -e POSTGRES_DB=$(PGBASE) -d postgres:14-alpine
start_postgres:
	@docker start simple_bank_db
stop_postgres:
	@docker stop simple_bank_db
psql:
	@PGPASSWORD=$(PGPASSWORD) psql -U $(PGUSER) -h $(PGHOST) -p $(PGPORT) -d $(PGBASE)
create_db:
	docker exec -it simple_bank_db createdb --username=$(PGUSER) $(PGBASE)
drop_db:
	docker exec -it simple_bank_db dropdb --username=$(PGUSER) $(PGBASE)

migrate_create:
	migrate create -ext sql -dir db/migrations -seq $(scheme)
migrate_up:
	@migrate -path db/migrations -database '$(DB_SOURCE)' -verbose up
migrate_up_last:
	@migrate -path db/migrations -database '$(DB_SOURCE)' -verbose up 1
migrate_down:
	@migrate -path db/migrations -database '$(DB_SOURCE)' -verbose down
migrate_down_last:
	@migrate -path db/migrations -database '$(DB_SOURCE)' -verbose down 1


gen_sqlc:
	sqlc generate

sqlc: gen_sqlc mocks fieldalignment

tests:
	go test -v -cover ./...

fieldalignment:
	fieldalignment -fix ./... 

server:
	@GIN_MODE=release go run main.go

mocks:
	mockgen -build_flags=--mod=mod -package mockdb -destination db/mock/store.go $(package_name)/db/sqlc Store


.PHONY: run_postgres start_postgres createdb dropdb recreate_db psql sqlc migrate_create migrate_up mgirate_down fieldalignment server gen_mocks gen_sqlc
