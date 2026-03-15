aws_region          = "us-east-1"
project             = "pizza-hub"
eks_cluster_version = "1.35"

# VPC Configuration  
vpc_cidr             = "10.0.0.0/16"
public_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.11.0/24", "10.0.12.0/24"]  