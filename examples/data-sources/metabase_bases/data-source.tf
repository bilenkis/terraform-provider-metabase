data "metabase_bases" "all" {}

# Returns all databases
output "all_bases" {
  value = data.metabase_bases.all.databases
}
