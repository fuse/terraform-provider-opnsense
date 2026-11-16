package acmeclient_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAcmeclientActionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAcmeclientActionConfig("acc-reload-gui"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_acmeclient_action.test", "name", "acc-reload-gui"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_action.test", "type", "configd_restart_gui"),
					resource.TestCheckResourceAttr("opnsense_acmeclient_action.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("opnsense_acmeclient_action.test", "id"),
				),
			},
			{
				ResourceName:      "opnsense_acmeclient_action.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAcmeclientActionConfig("acc-reload-gui-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_acmeclient_action.test", "name", "acc-reload-gui-updated"),
				),
			},
		},
	})
}

func testAccAcmeclientActionConfig(name string) string {
	return fmt.Sprintf(`
resource "opnsense_acmeclient_action" "test" {
  name = %[1]q
  type = "configd_restart_gui"
}
`, name)
}
