# Pizza Hub Infrastructure

This repository contains the Terraform configuration for deploying the Pizza Hub application's infrastructure on AWS.

## Summary

This Terraform setup provisions the following core components:

- **VPC**: A dedicated Virtual Private Cloud to ensure network isolation.
- **EKS**: An Amazon EKS cluster to orchestrate the application's containerized services.
- **RDS**: A managed PostgreSQL database for persistent data storage.
- **Redis**: An ElastiCache for Redis instance for caching.

## Setup

Before you can apply the Terraform configuration, you need to set up a remote backend on AWS to store the Terraform state. This ensures that state is managed centrally and securely.

### 1. Create S3 Bucket for Terraform State

First, create an S3 bucket with versioning and encryption enabled.

```bash
export BUCKET_NAME="pizza-hub-terraform-state-$(date +%s)"
aws s3 mb s3://$BUCKET_NAME

aws s3api put-bucket-versioning \
  --bucket $BUCKET_NAME \
  --versioning-configuration Status=Enabled

# Enable encryption
aws s3api put-bucket-encryption \
  --bucket $BUCKET_NAME \
  --server-side-encryption-configuration '{
    "Rules": [{
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }]
  }'
```

### 2. Create DynamoDB Table for State Locking

Next, create a DynamoDB table to handle state locking and prevent concurrent modifications.

```bash
aws dynamodb create-table \
  --table-name pizza-hub-terraform-locks \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
```

### 3. Configure the Backend

Update the `providers.tf` file with the name of your S3 bucket and DynamoDB table.

## Workspaces

Terraform workspaces are used to manage multiple environments (e.g., `dev`, `staging`, `prod`).

To create a new workspace, run the following command:

```bash
terraform workspace new <workspace_name>
```

## Modules

This infrastructure is organized into the following modules:

### [VPC](./modules/vpc/README.md)

The VPC module creates a secure and isolated network environment for the application. For more details on its architecture, variables, and outputs.

### [EKS](./modules/eks/README.md)

The EKS module provisions a managed Kubernetes cluster. For detailed information on its configuration, including node groups, IAM roles, and security settings.

### [RDS](./modules/rds/README.md)

The RDS module sets up a managed PostgreSQL database. See its README for details on security, high availability, and configuration.

### [Redis](./modules/redis/README.md)

The Redis module provides a high-performance caching layer using ElastiCache for Redis. Refer to its README for architecture and setup details.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
