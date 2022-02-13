PROFILE=


build:
	@docker-compose build

authenticatie:
	$(eval AWS_ACCESS_KEY_ID := $(shell aws --profile $(PROFILE) configure get aws_access_key_id))
	$(eval AWS_SECRET_ACCESS_KEY := $(shell aws --profile $(PROFILE) configure get aws_secret_access_key))
	$(eval AWS_SESSION_TOKEN := $(shell aws --profile $(PROFILE) configure get aws_session_token))

init: build
	@docker-compose run --rm terraform sh -c "terraform init"

run: build authenticatie
	@docker-compose run --rm -e AWS_REGION=ap-northeast-1 -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) -e AWS_SESSION_TOKEN=$(AWS_SESSION_TOKEN) terraform bash

plan: build authenticatie
	@docker-compose run --rm -e AWS_REGION=ap-northeast-1 -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) -e AWS_SESSION_TOKEN=$(AWS_SESSION_TOKEN) terraform sh -c "terraform plan"

apply: build authenticatie
	@docker-compose run --rm -e AWS_REGION=ap-northeast-1 -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) -e AWS_SESSION_TOKEN=$(AWS_SESSION_TOKEN) terraform sh -c "terraform apply"

destroy: build authenticatie
	@docker-compose run --rm -e AWS_REGION=ap-northeast-1 -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) -e AWS_SESSION_TOKEN=$(AWS_SESSION_TOKEN) terraform sh -c "terraform destroy"

refresh: build authenticatie
	@docker-compose run --rm -e AWS_REGION=ap-northeast-1 -e AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) -e AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) -e AWS_SESSION_TOKEN=$(AWS_SESSION_TOKEN) terraform sh -c "terraform refresh"

tfcheck:
	reflex -r '\.tf$$' ./tflint.sh

.PHONY: build init run plan apply destroy refresh authenticatie tfcheck
