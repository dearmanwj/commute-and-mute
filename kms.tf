data "aws_caller_identity" "current" {}

resource "aws_kms_key" "cam_idp" {
  description              = "Asymmetric key for signing CAM jwts"
  customer_master_key_spec = "ECC_NIST_P256"
  key_usage                = "SIGN_VERIFY"
}

resource "aws_kms_key_policy" "root_key_policy" {
    key_id = aws_kms_key.cam_idp.key_id
    policy = jsonencode({
    Id = "cam_idp_policy"
    Statement = [
      {
        Action = "kms:*"
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
        }

        Resource = "*"
        Sid      = "Enable IAM User Permissions"
      },
      {
        Action = "kms:*"
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:user/commute-and-mute-terraform"
        }

        Resource = "*"
        Sid      = "Enable terraform sa permissions"
      },
    {
        Action = "kms:*"
        Effect = "Allow"
        Principal = {
          AWS = "${aws_iam_role.iam_for_lambda.arn}"
        }
        "Action": [
                "kms:DescribeKey",
                "kms:GetPublicKey",
                "kms:Sign",
                "kms:Verify"
            ]
        Resource = "*"
        Sid      = "Enable key usage by lambda"
      },
    ]
    Version = "2012-10-17"
  })
}
