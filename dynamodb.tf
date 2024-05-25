resource "aws_dynamodb_table" "cam-users" {
  name           = "Users"
  billing_mode   = "PROVISIONED"
  read_capacity  = 10
  write_capacity = 10
  hash_key       = "ID"

  attribute {
    name = "ID"
    type = "N"
  }

}
