resource "opnsense_interfaces_vip" "proxyarp_example_vip" {
  mode        = "proxyarp"
  interface   = "wan"
  network     = "192.168.0.189/32"
  description = "proxyarp example vip"
}

resource "opnsense_interfaces_vip" "ipalias_example_vip" {
  mode        = "ipalias"
  interface   = "wan"
  network     = "192.168.0.166/32"
  description = "ipalias example vip"
}

resource "opnsense_interfaces_vip" "carp_example_vip" {
  mode        = "carp"
  interface   = "opt1"
  network     = "192.168.10.254/24"
  description = "carp example vip"
  password    = "changeme"
  vhid        = "1"
  advbase     = "1"
  advskew     = "0"
}