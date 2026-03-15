# Data sources  
data "aws_caller_identity" "current" {}  
data "aws_partition" "current" {}  

################################################################################
# KMS Key for EKS secrets encryption  
################################################################################
resource "aws_kms_key" "eks" {
  description             = "${var.cluster_name} EKS Secret Encryption Key"  
  deletion_window_in_days = 7  
  enable_key_rotation     = true  

  tags = merge(var.tags, {  
    Name = "${var.cluster_name}-eks-kms-key"  
  })  
}

resource "aws_kms_alias" "eks" {
  name          = "alias/${var.cluster_name}-eks-kms-key"
  target_key_id = aws_kms_key.eks.key_id
}

################################################################################
# EKS Cluster IAM Role  
################################################################################
resource "aws_iam_role" "cluster" {  
  name = "${var.cluster_name}-cluster-role"  

  assume_role_policy = jsonencode({  
    Version = "2012-10-17"  
    Statement = [  
      {  
        Action = "sts:AssumeRole"  
        Effect = "Allow"  
        Principal = {  
          Service = "eks.amazonaws.com"  
        }  
      }  
    ]  
  })  
tags = merge(var.tags, {  
    Name = "${var.cluster_name}-cluster-role"  
  })   
}  

resource "aws_iam_role_policy_attachment" "cluster_AmazonEKSClusterPolicy" {
  role       = aws_iam_role.cluster.name
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonEKSClusterPolicy"
}

################################################################################
# EKS Node Group IAM Role  
################################################################################
resource "aws_iam_role" "node" {
  name = "${var.cluster_name}-node-group-role"  

  assume_role_policy = jsonencode({  
    Version = "2012-10-17"  
    Statement = [  
      {  
        Action = "sts:AssumeRole"  
        Effect = "Allow"  
        Principal = {  
          Service = "ec2.amazonaws.com"  
        }  
      }  
    ]  
  })  

 tags = merge(var.tags, {  
    Name = "${var.cluster_name}-node-group-role"  
  })  
}

resource "aws_iam_role_policy_attachment" "node_AmazonEKSWorkerNodePolicy" {
  role       = aws_iam_role.node.name
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonEKSWorkerNodePolicy"
}

resource "aws_iam_role_policy_attachment" "node_AmazonEKS_CNI_Policy" {
  role       = aws_iam_role.node.name
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonEKS_CNI_Policy"
}

resource "aws_iam_role_policy_attachment" "node_AmazonEC2ContainerRegistryReadOnly" {
  role = aws_iam_role.node.name
  policy_arn = "arn:${data.aws_partition.current.partition}:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

################################################################################
# Security Group for EKS Cluster  
################################################################################
resource "aws_security_group" "cluster" {
  name        = "${var.cluster_name}-cluster-sg"
  description = "Security group for EKS cluster control plane"
  vpc_id      = var.vpc_id

    ## Allow all outbound traffic from the cluster security group
  egress {  
    from_port   = 0  
    to_port     = 0  
    protocol    = "-1"  
    cidr_blocks = ["0.0.0.0/0"]  
  }  

tags = merge(var.tags, {  
    Name = "${var.cluster_name}-cluster-sg"  
  })  
}

################################################################################
# Security Group for EKS Node Groups
################################################################################
resource "aws_security_group" "node" {  
  name        = "${var.cluster_name}-node-sg"  
  description = "Security group for EKS node groups"  
  vpc_id      = var.vpc_id  

    ## Allow all outbound traffic from the node security group
  egress {  
    from_port   = 0  
    to_port     = 0  
    protocol    = "-1"  
    cidr_blocks = ["0.0.0.0/0"]  
  }  

  tags = merge(var.tags, {  
    Name = "${var.cluster_name}-node-sg"  
  })  
}  

################################################################################
# Security Group Rules  
################################################################################
resource "aws_security_group_rule" "cluster_ingress_node_https" {
    description = "Allow nodes to communicate with cluster API Server"  
  security_group_id = aws_security_group.cluster.id
  source_security_group_id = aws_security_group.node.id
  type              = "ingress"
  from_port                = 443  
  protocol                 = "tcp" 
 to_port                  = 443  
}

resource "aws_security_group_rule" "node_ingress_self" {  
  description       = "Allow node to communicate with each other"  
  from_port         = 0  
  protocol          = "-1"  
  security_group_id = aws_security_group.node.id  
  self              = true  
  to_port           = 65535  
  type              = "ingress"  
}  

resource "aws_security_group_rule" "node_ingress_cluster_https" {  
  description              = "Allow cluster control plane to communicate with worker node kubelet and kube-proxy"  
  from_port                = 443  
  protocol                 = "tcp"  
  security_group_id        = aws_security_group.node.id  
  source_security_group_id = aws_security_group.cluster.id  
  to_port                  = 443  
  type                     = "ingress"  
}  

resource "aws_security_group_rule" "node_ingress_cluster_kubelet" {  
  description              = "Allow cluster control plane to communicate with worker node kubelet"  
  from_port                = 10250  
  protocol                 = "tcp"  
  security_group_id        = aws_security_group.node.id  
  source_security_group_id = aws_security_group.cluster.id  
  to_port                  = 10250  
  type                     = "ingress"  
}  

################################################################################
# EKS Cluster  
################################################################################
resource "aws_eks_cluster" "cluster" {
    name = var.cluster_name
    version = var.cluster_version
    role_arn = aws_iam_role.cluster.arn

    vpc_config {
      subnet_ids = var.subnet_ids
      endpoint_private_access = var.cluster_endpoint_private_access
      endpoint_public_access = var.cluster_endpoint_public_access
      public_access_cidrs = var.public_access_cidrs
      security_group_ids = [aws_security_group.cluster.id]
    }


  encryption_config {  
    provider {  
      key_arn = aws_kms_key.eks.arn  
    }  
    resources = ["secrets"]  
  }  

  enabled_cluster_log_types = [ "api", "audit", "authenticator", "controllerManager", "scheduler" ]

   depends_on = [  
    aws_iam_role_policy_attachment.cluster_AmazonEKSClusterPolicy,  
  ]  

  tags = merge(var.tags, {  
    Name = "${var.cluster_name}-eks-cluster"
    })
}

################################################################################
# OIDC Identity Provider  
################################################################################
data "tls_certificate" "cluster" {  
  url = aws_eks_cluster.cluster.identity[0].oidc[0].issuer  
}  

resource "aws_iam_openid_connect_provider" "cluster" {
   client_id_list  = ["sts.amazonaws.com"]  
  thumbprint_list = [data.tls_certificate.cluster.certificates[0].sha1_fingerprint]  
  url             = aws_eks_cluster.cluster.identity[0].oidc[0].issuer  

    tags = merge(var.tags, {  
     Name = "${var.cluster_name}-oidc-provider"  
      })
}

################################################################################
# EKS Managed Node Group  
################################################################################
resource "aws_eks_node_group" "node" {
  for_each = var.node_groups  

    cluster_name    = aws_eks_cluster.cluster.name
    node_group_name = each.key
    node_role_arn   = aws_iam_role.node.arn
    subnet_ids      = var.subnet_ids

    scaling_config {
         desired_size = each.value.desired_capacity  
    max_size     = each.value.max_capacity  
    min_size     = each.value.min_capacity  
    }

  update_config {  
    max_unavailable = 1  
  }  

instance_types = each.value.instance_types
capacity_type  = lookup(each.value, "capacity_type", "ON_DEMAND")  
  disk_size      = lookup(each.value, "disk_size", 20)  

  # Kubernetes labels  
    labels = merge(  
    {  
      "node-group" = each.key  
    },  
    lookup(each.value, "labels", {})  
  )  

    # Kubernetes taints  
  dynamic "taint" {  
    for_each = lookup(each.value, "taints", [])  
    content {  
      key    = taint.value.key  
      value  = taint.value.value  
      effect = taint.value.effect  
    }  
  }  


  depends_on = [  
    aws_iam_role_policy_attachment.node_AmazonEKSWorkerNodePolicy,  
    aws_iam_role_policy_attachment.node_AmazonEKS_CNI_Policy,  
    aws_iam_role_policy_attachment.node_AmazonEC2ContainerRegistryReadOnly,  
  ]  

 tags = merge(var.tags, {
    Name = "${var.cluster_name}-node-group-${each.key}"
 })

  lifecycle {  
    ignore_changes = [scaling_config[0].desired_size]  
  }  
}

################################################################################
# EKS Addons
################################################################################
resource "aws_eks_addon" "vpc_cni" {  
for_each = var.add_ons

  cluster_name = aws_eks_cluster.cluster.name  
  addon_name   = each.key
  resolve_conflicts_on_create = "OVERWRITE"  
  resolve_conflicts_on_update = "OVERWRITE"

  depends_on = [aws_eks_node_group.node]  

  tags = var.tags  
}  

