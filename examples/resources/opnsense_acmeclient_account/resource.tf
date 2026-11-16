resource "opnsense_acmeclient_account" "staging" {
  name  = "letsencrypt-staging"
  email = "admin@example.com"
  ca    = "letsencrypt_test"
}

resource "opnsense_acmeclient_account" "production" {
  name  = "letsencrypt-production"
  email = "admin@example.com"
  ca    = "letsencrypt"
}
