SHELL := /bin/sh
TEST?=./...
TESTARGS=-test.v
FUNC=
AWS_REGION=ap-northeast-1
SEQ=
VER=
MYSQL_URL="mysql://root:password@tcp(db:3306)/todo?x-tls-insecure-skip-verify=false"
TZ?=Asia/Tokyo

default: build

build: fmtcheck
	docker build --target build-stage -t build-server --build-arg TZ=$(TZ) .
	docker build -t server --build-arg TZ=$(TZ) .

dev-run:
	COMPOSE_FILE="docker-compose.yml:docker-compose-dev.yml" \
		TARGET=build-stage ENTRYPOINT="/opt/todo/scripts/wait-for-mysql.sh make hot" \
		docker-compose up --build

debug-run:
	COMPOSE_FILE="docker-compose.yml:docker-compose-debug.yml" \
		TARGET=debug-stage ENTRYPOINT="/opt/todo/scripts/wait-for-mysql.sh make hot-debug" \
		docker-compose up --build

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

hot-debug:
	reflex -r '(\.go$$|\.html$$)' -s make debug

debug:
	go build -gcflags "all=-N -l" -o debug_app
	dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./debug_app

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
	migrate -database=$(MYSQL_URL) -path=db/migrations/ up 1

migrate_down:
	migrate -database=$(MYSQL_URL) -path=db/migrations/ down 1

migrate_force:
	migrate -database=$(MYSQL_URL) -path=db/migrations/ force $(VER)

sql:
	mysql -h db -u root -ppassword todo

sqldump:
	mysqldump -h db -u root -ppassword --databases todo > db/dump.sql

drop_table:
	mysql -h db -u root -ppassword todo -e'SET FOREIGN_KEY_CHECKS = 0; DROP TABLE IF EXISTS tasks; DROP TABLE IF EXISTS users; DROP TABLE IF EXISTS sessions;'

restore_table:
	mysql -h db -u root -ppassword < db/dump.sql

.PHONY: default build dev-run degub-run run test_up test_down test hottest hot hot-debug debug generate fmt fmtcheck create_migration migrate_up migrate_down migrate_force sql sqldump drop_table restore_table
