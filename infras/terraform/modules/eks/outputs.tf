output "cluster_id" {
  description = "EKS cluster ID"
  value       = aws_eks_cluster.cluster.id
}

output "cluster_arn" {
  description = "EKS cluster ARN"
  value       = aws_eks_cluster.cluster.arn
}

output "cluster_name" {
  description = "EKS cluster name"
  value       = aws_eks_cluster.cluster.name
}

output "cluster_endpoint" {
  description = "EKS cluster endpoint"
  value       = aws_eks_cluster.cluster.endpoint
}

output "cluster_version" {
  description = "EKS cluster Kubernetes version"
  value       = aws_eks_cluster.cluster.version
}

output "cluster_ca_certificate" {
  description = "Base64 encoded certificate data required to communicate with cluster"
  value       = aws_eks_cluster.cluster.certificate_authority[0].data
}

output "cluster_oidc_issuer_url" {
  description = "The URL on the EKS cluster OIDC Issuer"
  value       = aws_eks_cluster.cluster.identity[0].oidc[0].issuer
}

output "oidc_provider_arn" {
  description = "ARN of the OIDC Provider for EKS"
  value       = aws_iam_openid_connect_provider.cluster.arn
}

output "node_groups" {
  description = "EKS node groups"
  value = {
    for k, v in aws_eks_node_group.node : k => {
      arn           = v.arn
      status        = v.status
      capacity_type = v.capacity_type
      instance_types = v.instance_types
    }
  }
}

output "cluster_security_group_id" {
  description = "ID of the EKS cluster security group"
  value       = aws_security_group.cluster.id
}

output "node_security_group_id" {
  description = "ID of the EKS node security group"  
  value       = aws_security_group.node.id
}
