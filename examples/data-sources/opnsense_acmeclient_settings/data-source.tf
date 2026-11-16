data "opnsense_acmeclient_settings" "current" {}

output "acmeclient_enabled" {
  value = data.opnsense_acmeclient_settings.current.enabled
}
