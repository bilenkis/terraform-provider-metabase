# Returns parameters for the database with id=1
data "metabase_base" "one" {
  id = 1
}

output "base_name" {
  value = data.metabase_base.one
}
