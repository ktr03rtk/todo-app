SHELL := /bin/sh
TEST?=./...
TESTARGS=-test.v
FUNC=
AWS_REGION=ap-northeast-1
SEQ=
MYSQL_URL="mysql://root:password@tcp(db:3306)/todo?x-tls-insecure-skip-verify=false"
TZ?=Asia/Tokyo

default: build

build: fmtcheck
	docker build --target build-stage -t build-server --build-arg TZ=$(TZ) .
	docker build -t server --build-arg TZ=$(TZ) .

dev-run:
	COMPOSE_FILE="docker-compose.yml:docker-compose-dev.yml" \
		TARGET=build-stage ENTRYPOINT="make hot" docker-compose up --build

run: build
	docker-compose up --build

test:
	gotest $(TEST) $(TESTARGS) -coverprofile=cover.out.tmp
	@cat cover.out.tmp | grep -v "mock_" > cover.out
	@go tool cover -html=cover.out -o cover.html
	@rm -f cover.out cover.out.tmp

test_up: test_down build
	docker run -ti --rm -p 8080:8080 -v $(PWD)/:/opt/todo/ \
		--name build-server -e ENVIRONMENT=development build-server:latest

test_down:
	docker rm -f build-server

hottest:
ifndef FUNC
	reflex -r '\.go$$' -- gotest $(TEST) $(TESTARGS)
else
	reflex -r '\.go$$' -- gotest -run $(FUNC) $(TEST) $(TESTARGS)
endif

hot:
	reflex -r '(\.go$$|\.html$$)' -s go run main.go

generate: fmtcheck
	go generate ./...

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w .

fmtcheck:
	sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

create_migration:
	migrate create -ext sql -dir db/migrations -seq $(SEQ)

migrate_up:
	migrate -database=$(MYSQL_URL) -path=db/migrations/ up

migrate_down:
	migrate -database=$(MYSQL_URL) -path=db/migrations/ down

sql:
	mysql -h db -u root -ppassword todo

sqldump:
	mysqldump -h db -u root -ppassword --databases todo > db/dump.sql

drop_table:
	mysql -h db -u root -ppassword todo -e'DROP TABLE tasks;'

restore_table:
	mysql -h db -u root -ppassword < db/dump.sql

.PHONY: default build dev-run run test_up test_down test hottest generate fmt fmtcheck create_migration migrate_up migrate_down sql sqldump drop_table restore_table
