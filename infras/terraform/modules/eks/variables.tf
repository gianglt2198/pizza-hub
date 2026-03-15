variable "cluster_name" {
  description = "Name of the EKS cluster"
  type        = string
}

variable "cluster_version" {
  description = "Kubernetes version"
  type        = string
   default     = "1.35"  

  validation {  
    condition = contains([  
      "1.35"
    ], var.cluster_version)  
    error_message = "Kubernetes version must be 1.35."  
  }  
}

variable "cluster_endpoint_private_access" {
  description = "Indicates whether or not the Amazon EKS private API server endpoint is enabled"
  type        = bool
  default     = true
}

variable "cluster_endpoint_public_access" {
  description = "Indicates whether or not the Amazon EKS public API server endpoint is enabled"
  type        = bool
  default     = true
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "subnet_ids" {
  description = "List of subnet IDs"
  type        = list(string)

    validation {  
    condition     = length(var.subnet_ids) >= 2  
    error_message = "At least 2 subnets are required for high availability."  
  }  
}

variable "enable_irsa" {
  description = "Enable IAM roles for service accounts"
  type        = bool
  default     = true
}

variable "node_groups" {  
  description = "Map of EKS node group configurations"  
  type = map(object({  
    desired_capacity = number  
    max_capacity     = number  
    min_capacity     = number  
    instance_types   = list(string)  
    capacity_type    = optional(string, "ON_DEMAND")  
    disk_size        = optional(number, 20)  
    labels           = optional(map(string), {})  
    taints = optional(list(object({  
      key    = string  
      value  = string  
      effect = string  
    })), [])  
  }))  

  validation {  
    condition = alltrue([  
      for k, v in var.node_groups : v.min_capacity <= v.desired_capacity && v.desired_capacity <= v.max_capacity  
    ])  
    error_message = "Node group capacity must satisfy: min <= desired <= max."  
  }  
}  

variable "public_access_cidrs" {  
  description = "List of CIDR blocks that can access the EKS cluster endpoint"  
  type        = list(string)  
  default     = ["0.0.0.0/0"]  
}  


variable "tags" {
  description = "A map of tags to add to all resources"
  type        = map(string)
  default     = {}
}

variable "add_ons" {
  description = "Map of EKS add-on configurations"
  type = map(any)
  default = {
     coredns = {}
     eks-pod-identity-agent = {}
     kube-proxy = {}
     vpc-cni = { }
  }
}