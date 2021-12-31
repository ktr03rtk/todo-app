SHELL := /bin/bash
TEST?=./...
TESTARGS=-test.v
TEST_FUNC=
AWS_REGION=ap-northeast-1
IMAGE=iotlogtransformfunction

default: build

build: fmtcheck
	docker build --target build-stage -t server:build .

run: build
	docker run --rm -p 8080:8080 server:build /opt/server

test_up: test_down
	docker run -ti --rm -p 8080:8080 -v $(PWD)/:/opt/ \
		--name server -e ENVIRONMENT=development server:build bash

test_down:
	docker rm -f server

test: fmtcheck
	gotest $(TEST) $(TESTARGS) -timeout=5m -parallel=8 -coverprofile=cover.out.tmp
	@cat cover.out.tmp | grep -v "mock_" > cover.out
	@go tool cover -html=cover.out -o cover.html
	@rm -f cover.out cover.out.tmp

hottest:
ifndef FUNC
	reflex -r '\.go$$' -- gotest $(TEST) $(TESTARGS) -parallel=8
else
	reflex -r '\.go$$' -- gotest -run $(TESTFUNC) $(TEST) $(TESTARGS)
endif

ci_test:
	docker run --rm --name server -e ENVIRONMENT=development server:build sh -c "gotest ./... -v"

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w .

fmtcheck:
	@sh -c "find . -name '*.go' | grep -v vendor | xargs gofmt -l -s"

.PHONY: default build run test_up test_down test hottest ci_test fmt fmtcheck
