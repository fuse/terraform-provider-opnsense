resource "opnsense_acmeclient_account" "example" {
  name  = "letsencrypt-staging"
  email = "admin@example.com"
  ca    = "letsencrypt_test"
}

resource "opnsense_acmeclient_validation" "example" {
  name         = "cloudflare-dns01"
  method       = "dns01"
  dns_service  = "dns_cf"
  dns_cf_token = "replace-with-cloudflare-api-token"
}

resource "opnsense_acmeclient_action" "reload_gui" {
  name = "reload-web-gui"
  type = "configd_restart_gui"
}

resource "opnsense_acmeclient_certificate" "www" {
  name              = "www.example.com"
  account           = opnsense_acmeclient_account.example.id
  validation_method = opnsense_acmeclient_validation.example.id
  key_length        = "2048"
  auto_renewal      = true
  renew_interval    = "60"

  alt_names = [
    "example.com",
  ]

  restart_actions = [
    opnsense_acmeclient_action.reload_gui.id,
  ]
}
