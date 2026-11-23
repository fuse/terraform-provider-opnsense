package routing_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoutingGatewayResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingGatewayResourceConfig("192.0.2.254"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_routing_gateway.test", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_routing_gateway.test", "name", "tf_acc_gw"),
					resource.TestCheckResourceAttr("opnsense_routing_gateway.test", "interface", "wan"),
					resource.TestCheckResourceAttr("opnsense_routing_gateway.test", "ip_protocol", "inet"),
					resource.TestCheckResourceAttr("opnsense_routing_gateway.test", "gateway", "192.0.2.254"),
					resource.TestCheckResourceAttr("opnsense_routing_gateway.test", "default_gateway", "false"),
					resource.TestCheckResourceAttr("opnsense_routing_gateway.test", "far_gateway", "true"),
					resource.TestCheckResourceAttr("opnsense_routing_gateway.test", "monitor_disable", "true"),
					resource.TestCheckResourceAttrSet("opnsense_routing_gateway.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_routing_gateway.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRoutingGatewayResourceConfig("192.0.2.253"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_routing_gateway.test", "gateway", "192.0.2.253"),
				),
			},
		},
	})
}

func testAccRoutingGatewayResourceConfig(gatewayIP string) string {
	return fmt.Sprintf(`
resource "opnsense_routing_gateway" "test" {
  enabled          = true
  name             = "tf_acc_gw"
  description      = "Acceptance test gateway"
  interface        = "wan"
  ip_protocol      = "inet"
  gateway          = %q
  default_gateway  = false
  far_gateway      = true
  monitor_disable  = true
  monitor_noroute  = false
  force_down       = false
  priority         = "255"
  weight           = "1"
}
`, gatewayIP)
}
