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
  type    = string
  default = "Sample Dataset"
}

provider "metabase" {
  username = var.metabase_username
  password = var.metabase_password
}

#module "base" {
#source = "./base"

#base_name = var.base_name
#}

#output "all_bases" {
#value = module.base.all_bases
#}

#output "base_name_to_id" {
#value = module.base.id[var.base_name]
#}

data "metabase_base" "one" {
  id = 1
}

output "base_name" {
  value = data.metabase_base.one
}

#data "metabase_bases" "all" {}

## Returns all coffees
#output "all_bases" {
  #value = data.metabase_bases.all.databases
#}

## Only returns packer spiced latte
#output "id" {
  #value = {
    #for base in data.metabase_bases.all.databases :
    #base.name => base.id
    #if base.name == var.base_name
  #}
#}
