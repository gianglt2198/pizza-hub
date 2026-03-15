data "aws_eks_cluster" "cluster" {
  name = var.cluster_name

  depends_on = [aws_eks_cluster.cluster]
}

data "aws_eks_cluster_auth" "cluster" {
  name = var.cluster_name

  depends_on = [aws_eks_cluster.cluster]
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  token                  = data.aws_eks_cluster_auth.cluster.token
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority[0].data)
}