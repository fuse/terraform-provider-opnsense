package acmeclient_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAcmeclientValidationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAcmeclientValidationConfig("acc-cf-dns", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_acmeclient_validation.test", "name", "acc-cf-dns"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_validation.test", "method", "dns01"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_validation.test", "dns_service", "dns_cf"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_validation.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("opnsense_acmeclient_validation.test", "id"),
				),
			},
			{
				ResourceName:            "opnsense_acmeclient_validation.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dns_cf_key", "dns_cf_token"},
			},
		},
	})
}

func testAccAcmeclientValidationConfig(name string, enabled bool) string {
	return fmt.Sprintf(`
resource "opnsense_acmeclient_validation" "test" {
  enabled      = %[2]t
  name         = %[1]q
  method       = "dns01"
  dns_service  = "dns_cf"
  dns_cf_token = "dummy"
}
`, name, enabled)
}
