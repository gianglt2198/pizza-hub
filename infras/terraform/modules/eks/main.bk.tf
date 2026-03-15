# module "eks" {
#   source = "terraform-aws-modules/eks/aws"
#   version = "~> 21.0"

#   name    = var.cluster_name
#   kubernetes_version = var.cluster_version

#   endpoint_private_access = var.cluster_endpoint_private_access
#   endpoint_public_access  = var.cluster_endpoint_public_access

#   vpc_id     = var.vpc_id
#   subnet_ids = var.subnet_ids

#     # EKS Addons
#   addons = {
#     coredns = {}
#     eks-pod-identity-agent = {
#       before_compute = true
#     }
#     kube-proxy = {}
#     vpc-cni = {
#       before_compute = true
#     }
#   }

#   enable_irsa = var.enable_irsa

#   compute_config = {
#     enabled    = true
#     node_pools = ["general-purpose"]
#   }

#   enable_cluster_creator_admin_permissions = true

#   tags = var.tags
# }