resource "opnsense_acmeclient_validation" "cloudflare" {
  name         = "cloudflare-dns01"
  method       = "dns01"
  dns_service  = "dns_cf"
  dns_cf_token = "replace-with-cloudflare-api-token"
}
