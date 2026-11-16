package acmeclient_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAcmeclientCertificateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAcmeclientCertificateConfig("acc.example.com", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_acmeclient_certificate.test", "name", "acc.example.com"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_certificate.test", "enabled", "false"),
					resource.TestCheckResourceAttrSet("opnsense_acmeclient_certificate.test", "id"),
					resource.TestCheckResourceAttrSet("opnsense_acmeclient_certificate.test", "account"),
					resource.TestCheckResourceAttrSet("opnsense_acmeclient_certificate.test", "validation_method"),
				),
			},
			{
				ResourceName:            "opnsense_acmeclient_certificate.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status_code", "status_last_update", "last_update", "cert_ref_id"},
			},
			{
				Config: testAccAcmeclientCertificateConfig("acc-updated.example.com", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_acmeclient_certificate.test", "name", "acc-updated.example.com"),
				),
			},
		},
	})
}

func testAccAcmeclientCertificateConfig(name string, enabled bool) string {
	return fmt.Sprintf(`
resource "opnsense_acmeclient_account" "test" {
  name  = "acc-cert-account"
  email = "admin@example.com"
  ca    = "letsencrypt_test"
}

resource "opnsense_acmeclient_validation" "test" {
  # Must be enabled: certificate.validationMethod only lists enabled validations.
  enabled      = true
  name         = "acc-cert-validation"
  method       = "dns01"
  dns_service  = "dns_cf"
  dns_cf_token = "dummy"
}

resource "opnsense_acmeclient_certificate" "test" {
  enabled           = %[2]t
  name              = %[1]q
  account           = opnsense_acmeclient_account.test.id
  validation_method = opnsense_acmeclient_validation.test.id
  auto_renewal      = false
}
`, name, enabled)
}
