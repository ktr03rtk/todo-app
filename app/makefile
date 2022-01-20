SHELL := /bin/sh
TEST?=./...
TESTARGS=-test.v
TEST_FUNC=
AWS_REGION=ap-northeast-1

default: build

build: fmtcheck
	docker build  -t server .
	docker build --target build-stage -t build-server .

dev-run: build
	docker run --rm -it -p 8080:8080 -v ${PWD}:/opt build-server:latest make hot

run: build
	docker run --rm -p 8080:8080 server:latest /opt/server

test:
	gotest $(TEST) $(TESTARGS) -coverprofile=cover.out.tmp
	@cat cover.out.tmp | grep -v "mock_" > cover.out
	@go tool cover -html=cover.out -o cover.html
	@rm -f cover.out cover.out.tmp

test_up: test_down
	docker run -ti --rm -p 8080:8080 -v $(PWD)/:/opt/ \
		--name build-server -e ENVIRONMENT=development build-server:latest sh

test_down:
	docker rm -f build-server

hottest:
ifndef FUNC
	reflex -r '\.go$$' -- gotest $(TEST) $(TESTARGS)
else
	reflex -r '\.go$$' -- gotest -run $(TESTFUNC) $(TEST) $(TESTARGS)
endif

hot:
	reflex -r '(\.go$$|\.html$$)' -s go run main.go

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w .

fmtcheck:
	sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: default build dev-run run test_up test_down test hottest fmt fmtcheck
