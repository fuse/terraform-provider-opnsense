package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfacesVipProxyArpResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccVipResourceConfig("proxyarp", "Proxy ARP VIP test", "wan", "192.168.2.22/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "mode", "proxyarp"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "description", "Proxy ARP VIP test"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "interface", "wan"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "network", "192.168.2.22/32"),
					resource.TestCheckResourceAttrSet("opnsense_interfaces_vip.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "opnsense_interfaces_vip.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update and Read testing
			{
				Config: testAccVipResourceConfig("proxyarp", "Updated Proxy ARP VIP", "wan", "192.168.2.23/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "mode", "proxyarp"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "description", "Updated Proxy ARP VIP"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "interface", "wan"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "network", "192.168.2.23/32"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccInterfacesVipIpAliasResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccVipResourceConfig("ipalias", "IP Alias VIP test", "wan", "192.168.2.24/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "mode", "ipalias"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "description", "IP Alias VIP test"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "interface", "wan"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "network", "192.168.2.24/32"),
					resource.TestCheckResourceAttrSet("opnsense_interfaces_vip.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "opnsense_interfaces_vip.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update and Read testing
			{
				Config: testAccVipResourceConfig("ipalias", "Updated IP Alias VIP", "wan", "192.168.2.25/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "mode", "ipalias"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "description", "Updated IP Alias VIP"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "interface", "wan"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "network", "192.168.2.25/32"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccInterfacesVipCarpResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVipCarpResourceConfig("CARP VIP test", "wan", "192.168.2.26/32", "1", "0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "mode", "carp"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "description", "CARP VIP test"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "interface", "wan"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "network", "192.168.2.26/32"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "vhid", "1"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "advskew", "0"),
					resource.TestCheckResourceAttrSet("opnsense_interfaces_vip.test", "id"),
				),
			},
			{
				ResourceName:            "opnsense_interfaces_vip.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccVipCarpResourceConfig("Updated CARP VIP", "wan", "192.168.2.26/32", "1", "100"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "description", "Updated CARP VIP"),
					resource.TestCheckResourceAttr("opnsense_interfaces_vip.test", "advskew", "100"),
				),
			},
		},
	})
}

func testAccVipResourceConfig(mode, description, interf, network string) string {
	return fmt.Sprintf(`
resource "opnsense_interfaces_vip" "test" {
  mode         = %[1]q
  description  = %[2]q
  interface    = %[3]q
  network      = %[4]q
}
`, mode, description, interf, network)
}

func testAccVipCarpResourceConfig(description, interf, network, vhid, advskew string) string {
	return fmt.Sprintf(`
resource "opnsense_interfaces_vip" "test" {
  mode        = "carp"
  description = %[1]q
  interface   = %[2]q
  network     = %[3]q
  password    = "tf-acc-carp-password"
  vhid        = %[4]q
  advbase     = "1"
  advskew     = %[5]q
}
`, description, interf, network, vhid, advskew)
}
