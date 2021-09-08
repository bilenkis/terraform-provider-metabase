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
