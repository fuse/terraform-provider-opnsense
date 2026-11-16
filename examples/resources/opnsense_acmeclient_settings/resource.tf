// Import the singleton resource before managing it:
// terraform import opnsense_acmeclient_settings.settings acmeclient_settings

resource "opnsense_acmeclient_settings" "settings" {
  enabled            = true
  auto_renewal       = true
  environment        = "prod"
  challenge_port     = "43580"
  tls_challenge_port = "43581"
  restart_timeout    = "600"
  log_level          = "normal"
}
