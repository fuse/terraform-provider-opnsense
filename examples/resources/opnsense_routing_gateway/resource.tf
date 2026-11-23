resource "opnsense_routing_gateway" "wan" {
  enabled          = true
  name             = "WAN_GW"
  description      = "Interface WAN Gateway"
  interface        = "wan"
  ip_protocol      = "inet"
  gateway          = "192.0.2.1"
  default_gateway  = true
  far_gateway      = true
  monitor_disable  = true
  priority         = "255"
  weight           = "1"
}

resource "opnsense_routing_gateway" "peer_site" {
  enabled          = true
  name             = "gw-peer"
  interface        = "opt2"
  ip_protocol      = "inet"
  gateway          = "10.0.0.10"
  default_gateway  = false
  far_gateway      = false
  monitor_disable  = true
}
