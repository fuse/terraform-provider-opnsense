package acmeclient_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAcmeclientAccountResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAcmeclientAccountConfig("acc-test-account", "admin@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_acmeclient_account.test", "name", "acc-test-account"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_account.test", "email", "admin@example.com"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_account.test", "ca", "letsencrypt_test"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_account.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_acmeclient_account.test", "id"),
				),
			},
			{
				ResourceName:            "opnsense_acmeclient_account.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_key", "eab_hmac", "status_code", "status_last_update"},
			},
			{
				Config: testAccAcmeclientAccountConfig("acc-test-account-updated", "admin@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_acmeclient_account.test", "name", "acc-test-account-updated"),
				),
			},
		},
	})
}

func testAccAcmeclientAccountConfig(name, email string) string {
	return fmt.Sprintf(`
resource "opnsense_acmeclient_account" "test" {
  name  = %[1]q
  email = %[2]q
  ca    = "letsencrypt_test"
}
`, name, email)
}
