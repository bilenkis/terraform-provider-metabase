terraform {
  required_providers {
    metabase = {
      version = "0.2"
      source  = "bilenkis.io/bilenkis/metabase"
    }
  }
}

variable "base_name" {
  type    = string
  default = "Sample Dataset"
}

data "metabase_bases" "all" {}

# Returns all coffees
output "all_bases" {
  value = data.metabase_bases.all.databases
}

# Only returns packer spiced latte
output "id" {
  value = {
    for base in data.metabase_bases.all.databases :
    base.name => base.id
    if base.name == var.base_name
  }
}
