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

Local testing uses AWS SAM. Terraform requires an AWS account and credentials
configured through the normal AWS SDK mechanisms, such as environment variables,
`~/.aws/credentials`, or an attached instance role.

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

The Lambda architecture must match the compiled binaries. The build scripts
detect the host architecture automatically. Override it explicitly when needed:

```bash
TF_VAR_lambda_architecture=arm64 make build
TF_VAR_lambda_architecture=arm64 make plan-tf
```
