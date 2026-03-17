locals {
  # Just keep essential common tags and minimal environment logic  
  common_tags = {
    Project     = var.project
    Environment = terraform.workspace
    ManagedBy   = "gianglt1"
    Workspace   = terraform.workspace
    CreatedDate = formatdate("YYYY-MM-DD", timestamp())
  }
}
