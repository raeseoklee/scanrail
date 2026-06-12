resource "aws_s3_bucket" "public_logs" {
  bucket = "scanrail-demo-public-logs"
}

resource "aws_s3_bucket_public_access_block" "public_logs" {
  bucket = aws_s3_bucket.public_logs.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}
