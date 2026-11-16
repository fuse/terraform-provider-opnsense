package acmeclient_test

import (
	"regexp"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAcmeclientSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testAccAcmeclientSettingsResourceConfig(false, "prod", "normal"),
				ResourceName:       "opnsense_acmeclient_settings.settings",
				ImportState:        true,
				ImportStateId:      "acmeclient_settings",
				ImportStatePersist: true,
			},
			{
				Config: testAccAcmeclientSettingsResourceConfig(true, "prod", "normal"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_acmeclient_settings.settings", "id", "acmeclient_settings"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_settings.settings", "enabled", "true"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_settings.settings", "auto_renewal", "true"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_settings.settings", "environment", "prod"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_settings.settings", "challenge_port", "43580"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_settings.settings", "tls_challenge_port", "43581"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_settings.settings", "restart_timeout", "600"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_settings.settings", "log_level", "normal"),
				),
			},
			{
				Config: testAccAcmeclientSettingsResourceConfig(true, "stg", "extended"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_acmeclient_settings.settings", "environment", "stg"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_settings.settings", "log_level", "extended"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_settings.settings", "enabled", "true"),
				),
			},
			{
				Config: testAccAcmeclientSettingsResourceConfig(true, "prod", "normal"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_acmeclient_settings.settings", "environment", "prod"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_settings.settings", "log_level", "normal"),
				),
			},
		},
	})
}

func TestAccAcmeclientSettingsResource_CreateBlocked(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccAcmeclientSettingsResourceConfigMinimal(),
				ExpectError: regexp.MustCompile("Cannot Create Singleton Resource"),
			},
		},
	})
}

func testAccAcmeclientSettingsResourceConfig(enabled bool, environment, logLevel string) string {
	enabledStr := "false"
	if enabled {
		enabledStr = "true"
	}
	return `
resource "opnsense_acmeclient_settings" "settings" {
  enabled            = ` + enabledStr + `
  auto_renewal       = true
  environment        = "` + environment + `"
  challenge_port     = "43580"
  tls_challenge_port = "43581"
  restart_timeout    = "600"
  log_level          = "` + logLevel + `"
}
`
}

func testAccAcmeclientSettingsResourceConfigMinimal() string {
	return `
resource "opnsense_acmeclient_settings" "settings" {
}
`
}
