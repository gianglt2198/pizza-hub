module "vpc" {
  source = "../../modules/vpc"

  project              = "pizza-hub"
  environment          = "dev"
  vpc_cidr             = "10.0.0.0/16"
  public_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24"]
  private_subnet_cidrs = ["10.0.11.0/24", "10.0.12.0/24"]
}

module "eks_cluster" {
  source = "../../modules/eks"

  node_ami_type = "AL2023_x86_64_STANDARD"
  project       = "pizza-hub"

  vpc_id          = module.vpc.vpc_id
  private_subnets = module.vpc.private_subnet_ids

  tags = {
    Environment = "dev"
    Project     = "pizza-hub"
    ManagedBy   = "gianglt1"
  }
}
