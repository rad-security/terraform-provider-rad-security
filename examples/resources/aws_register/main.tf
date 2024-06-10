resource "rad-security_aws_register" "this" {
  rad_security_assumed_role_arn = "arn:aws:iam::<aws_account_number>:role/rad-security-connector"
  aws_account_id        = "aws_account_id"
}
