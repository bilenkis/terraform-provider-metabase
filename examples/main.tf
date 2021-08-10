terraform {
  required_providers {
    metabase = {
      version = "0.2.3"
      source  = "bilenkis/metabase"
    }
  }
}

variable "metabase_username" {
  type = string
}

variable "metabase_password" {
  type = string
}

variable "metabase_url" {
  type    = string
  default = "http://localhost:3000"
}

variable "base_name" {
  type    = string
  default = "Sample Dataset"
}

provider "metabase" {
  url      = var.metabase_url
  username = var.metabase_username
  password = var.metabase_password
}

data "metabase_bases" "all" {}

# Returns all databases
output "all_bases" {
  value = data.metabase_bases.all.databases
}

# Only returns "Sample Dataset" database
output "id" {
  value = {
    for base in data.metabase_bases.all.databases :
    base.name => base.id
    if base.name == var.base_name
  }
}

# Returns parameters for the database with id=1
data "metabase_base" "one" {
  id = 1
}

output "base_name" {
  value = data.metabase_base.one
}

# Creates a new database
resource "metabase_database" "my" {
  name     = "test"
  engine   = "postgres"
  host     = "pg"
  port     = 5432
  db       = "postgres"
  user     = "postgres"
  password = "2I1dnzeYCIefM8Ru6Rxj"
}

output "my_id" {
  value = metabase_database.my.id
}
