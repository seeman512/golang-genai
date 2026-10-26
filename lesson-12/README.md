# Lesson 12 — Go Lambda functions with AWS SAM

This lesson contains four Go Lambda examples:

- `cmd/presentation` — direct Lambda invocation with a typed request and response;
- `cmd/presentation-webhook` — generic API Gateway Lambda-proxy webhook;
- `cmd/user-avatar` — deterministic initials-based SVG avatar generation;
- `cmd/user-events-webhook` — validates a `user.registered` event and emulates
  the greeting email that would be sent after registration.

The avatar generator uses only the standard library. It returns a stable SVG
for the same user ID and full name, so the demo does not need an image service,
SMTP provider, database migration, or external AI API.

See [`specs/lesson-12.md`](specs/lesson-12.md) for the requirements and
acceptance criteria.

## Local checks

```bash
make test
make build
```

The build creates four zip packages under `build/`. Each archive contains an
executable named `bootstrap`, as required by `provided.al2023`.

## Local testing with AWS SAM

The local workflow uses AWS SAM CLI and Docker. No AWS account is needed for
`sam local` commands.

```bash
make sam-build
make sam-invoke-avatar
make sam-invoke-welcome
make sam-smoke
```

To start the local API Gateway emulator interactively:

```bash
make sam-api
```

The local API exposes:

- `POST http://127.0.0.1:3000/webhook`;
- `POST http://127.0.0.1:3000/user-events`.

For example:

```bash
curl -i -X POST http://127.0.0.1:3000/user-events \
  -H 'Content-Type: application/json' \
  --data-binary @events/user-registered-body.json
```

SAM runs the custom `provided.al2023` functions in Docker containers. The first
invocation may download the Lambda runtime image.

## Deploying to real AWS

SAM can deploy the application through CloudFormation:

```bash
make sam-deploy
```

The command uses `sam deploy --guided` and requires AWS credentials, a region,
and permission to create Lambda, API Gateway, IAM, and CloudFormation resources.
Review the generated IAM role and stack before confirming deployment.

Terraform is also available as an optional real-AWS infrastructure-as-code
path:

```bash
make plan-tf
make deploy-tf
make destroy-tf
```

For an AWS-compatible Docker emulator named `flocki` (or Flocki), Terraform can
plan against the container without AWS credentials:

```bash
FLOCKI_ENDPOINT=http://localhost:4566 make plan-tf-flocki
```

Use `FLOCKI_ENDPOINT=http://flocki:4566` when Terraform runs in another
container on the same Docker network. The local plan uses separate
`terraform/flocki.tfstate` and `terraform/flocki.tfplan` files. `make
plan-tf-flocki` only creates a plan; it does not apply resources.

Use either SAM deployment or Terraform for a given AWS stack, not both. Terraform
configuration is documented in [`terraform/README.md`](terraform/README.md).

## Cleanup

Remove local build artifacts with:

```bash
make clean
```

Remove a SAM-deployed CloudFormation stack with the stack name selected during
`sam deploy --guided`, for example:

```bash
sam delete --stack-name lesson-12
```
