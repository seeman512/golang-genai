SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help
.PHONY: help test build sam-build sam-invoke-avatar sam-invoke-welcome sam-api sam-smoke sam-deploy plan-tf plan-tf-flocki deploy-tf destroy-tf invoke invoke-avatar invoke-welcome clean

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

test: ## Run tests and static analysis
	go vet ./...
	go test ./...

build: ## Build all Lambda zip packages for Terraform or packaging
	scripts/build.sh

# SAM invokes these targets because the template uses Metadata.BuildMethod=makefile.
build-PresentationFunction:
	@mkdir -p "$(ARTIFACTS_DIR)"
	@source scripts/env.sh && GOOS=linux GOARCH="$$GOARCH_LAMBDA" CGO_ENABLED=0 go build -tags lambda.norpc -ldflags='-s -w' -o "$(ARTIFACTS_DIR)/bootstrap" ./cmd/presentation

build-PresentationWebhookFunction:
	@mkdir -p "$(ARTIFACTS_DIR)"
	@source scripts/env.sh && GOOS=linux GOARCH="$$GOARCH_LAMBDA" CGO_ENABLED=0 go build -tags lambda.norpc -ldflags='-s -w' -o "$(ARTIFACTS_DIR)/bootstrap" ./cmd/presentation-webhook

build-UserAvatarFunction:
	@mkdir -p "$(ARTIFACTS_DIR)"
	@source scripts/env.sh && GOOS=linux GOARCH="$$GOARCH_LAMBDA" CGO_ENABLED=0 go build -tags lambda.norpc -ldflags='-s -w' -o "$(ARTIFACTS_DIR)/bootstrap" ./cmd/user-avatar

build-UserEventsWebhookFunction:
	@mkdir -p "$(ARTIFACTS_DIR)"
	@source scripts/env.sh && GOOS=linux GOARCH="$$GOARCH_LAMBDA" CGO_ENABLED=0 go build -tags lambda.norpc -ldflags='-s -w' -o "$(ARTIFACTS_DIR)/bootstrap" ./cmd/user-events-webhook

sam-build: ## Build the functions with AWS SAM
	@source scripts/env.sh && sam build --template-file sam/template.yaml --parameter-overrides ParameterKey=LambdaArchitecture,ParameterValue=$$LAMBDA_ARCH

sam-invoke-avatar: sam-build ## Invoke the avatar Lambda with SAM
	sam local invoke UserAvatarFunction --template .aws-sam/build/template.yaml --event events/user-avatar.json

sam-invoke-welcome: sam-build ## Invoke the registration webhook with SAM
	sam local invoke UserEventsWebhookFunction --template .aws-sam/build/template.yaml --event events/user-registered.json

sam-api: sam-build ## Start the local API Gateway emulator with SAM
	sam local start-api --template .aws-sam/build/template.yaml --host 127.0.0.1 --port 3000

sam-smoke: ## Run the SAM local smoke test
	scripts/sam-smoke.sh

sam-deploy: sam-build ## Deploy the SAM application to real AWS using guided CloudFormation deployment
	sam deploy --template-file .aws-sam/build/template.yaml --guided

plan-tf: build ## Initialize, validate, and show the Terraform plan for real AWS
	scripts/deploy-tf.sh plan

plan-tf-flocki: build ## Plan Terraform against the local Flocki Docker emulator
	@AWS_ACCESS_KEY_ID=test \
	AWS_SECRET_ACCESS_KEY=test \
	AWS_DEFAULT_REGION=us-east-1 \
	AWS_REGION=us-east-1 \
	AWS_EC2_METADATA_DISABLED=true \
	TF_VAR_local_aws=true \
	TF_VAR_local_aws_endpoint="$${FLOCKI_ENDPOINT:-http://localhost:4566}" \
	TF_STATE_FILE=flocki.tfstate \
	TF_PLAN_FILE=flocki.tfplan \
	scripts/deploy-tf.sh plan

deploy-tf: build ## Apply the Terraform configuration after reviewing the plan
	scripts/deploy-tf.sh apply

destroy-tf: ## Destroy resources in the Terraform state
	scripts/deploy-tf.sh destroy

invoke: sam-build ## Invoke the presentation function with SAM
	sam local invoke PresentationFunction --template .aws-sam/build/template.yaml --event events/presentation.json

invoke-avatar: sam-build ## Generate a deterministic user avatar with SAM
	sam local invoke UserAvatarFunction --template .aws-sam/build/template.yaml --event events/user-avatar.json

invoke-welcome: sam-build ## Emulate the registration welcome email with SAM
	sam local invoke UserEventsWebhookFunction --template .aws-sam/build/template.yaml --event events/user-registered.json

clean: ## Remove generated build and SAM metadata
	rm -rf build .aws-sam .lambda-role.json .lambda-policy.json
