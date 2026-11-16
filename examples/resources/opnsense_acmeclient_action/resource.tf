resource "opnsense_acmeclient_action" "reload_gui" {
  name        = "reload-web-gui"
  description = "Reload the OPNsense web GUI after certificate renewal"
  type        = "configd_restart_gui"
}
