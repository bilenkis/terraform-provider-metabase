terraform {
  required_providers {
    metabase = {
      version = "0.2"
      source  = "bilenkis.io/bilenkis/metabase"
    }
  }
}

variable "metabase_username" {
  type = string
}

variable "metabase_password" {
  type = string
}

variable "base_name" {
  type = string
  default = "Sample Dataset"
}

provider "metabase" {
  username = var.metabase_username
  password = var.metabase_password
}

module "base" {
  source = "./base"

  base_name = var.base_name
}

output "all_bases" {
  value = module.base.all_bases
}

output "id" {
  value = module.base.id[var.base_name]
}
