# Terraform deployment

This configuration is an optional real-AWS deployment path. It manages:

- the shared Lambda execution role and logging policy;
- `presentation`;
- `presentation-webhook`;
- `user-avatar`;
- `user-events-webhook`;
- API Gateway `POST /webhook`;
- API Gateway `POST /user-events`;
- Lambda invoke permissions, deployment, and stage.

Local testing uses AWS SAM. The default Terraform configuration targets a
real AWS account and requires credentials configured through the normal AWS SDK
mechanisms, such as environment variables, `~/.aws/credentials`, or an attached
instance role.

Build the packages from the repository root before planning:

```bash
make build
make plan-tf
make deploy-tf
```

Terraform state is local to this directory. `make deploy-tf` uses
`terraform apply -auto-approve`; review the plan with `make plan-tf` first.
Destroy resources with:

```bash
make destroy-tf
```

## Planning against Flocki

The provider also supports an AWS-compatible Docker emulator. Start the emulator
with its AWS API endpoint available at port `4566`, then run:

```bash
FLOCKI_ENDPOINT=http://localhost:4566 make plan-tf-flocki
```

When Terraform itself runs in Docker, use the emulator's container hostname
instead:

```bash
FLOCKI_ENDPOINT=http://flocki:4566 make plan-tf-flocki
```

The target uses fake credentials, routes the AWS provider services to the local
endpoint, and stores state and the generated plan separately as
`flocki.tfstate` and `flocki.tfplan`. It performs `init`, `validate`, and
`plan`, but does not apply anything. The emulator must implement the IAM,
Lambda, API Gateway, STS, and CloudWatch Logs APIs used by this configuration.

The local plan is not proof that the configuration will be accepted by real AWS;
use `make plan-tf` with real credentials for that check.

The Lambda architecture must match the compiled binaries. The build scripts
detect the host architecture automatically. Override it explicitly when needed:

```bash
TF_VAR_lambda_architecture=arm64 make build
TF_VAR_lambda_architecture=arm64 make plan-tf
```
