package servicecatalog_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfservicecatalog "github.com/hashicorp/terraform-provider-aws/internal/service/servicecatalog"
)

func testAccOrganizationsAccess_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_servicecatalog_organizations_access.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckOrganizationsEnabled(ctx, t)
			acctest.PreCheckOrganizationManagementAccount(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, servicecatalog.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckOrganizationsAccessDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationsAccessConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrganizationsAccessExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testAccCheckOrganizationsAccessDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ServiceCatalogConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_servicecatalog_organizations_access" {
				continue
			}

			output, err := tfservicecatalog.WaitOrganizationsAccessStable(ctx, conn, tfservicecatalog.OrganizationsAccessStableTimeout)

			if err != nil {
				return fmt.Errorf("error describing Service Catalog AWS Organizations Access (%s): %w", rs.Primary.ID, err)
			}

			if output == "" {
				return fmt.Errorf("error getting Service Catalog AWS Organizations Access (%s): empty response", rs.Primary.ID)
			}

			return nil
		}

		return nil
	}
}

func testAccCheckOrganizationsAccessExists(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ServiceCatalogConn()

		output, err := tfservicecatalog.WaitOrganizationsAccessStable(ctx, conn, tfservicecatalog.OrganizationsAccessStableTimeout)

		if err != nil {
			return fmt.Errorf("error describing Service Catalog AWS Organizations Access (%s): %w", rs.Primary.ID, err)
		}

		if output == "" {
			return fmt.Errorf("error getting Service Catalog AWS Organizations Access (%s): empty response", rs.Primary.ID)
		}

		if output != servicecatalog.AccessStatusEnabled && rs.Primary.Attributes["enabled"] == "true" {
			return fmt.Errorf("error getting Service Catalog AWS Organizations Access (%s): wrong setting", rs.Primary.ID)
		}

		if output == servicecatalog.AccessStatusEnabled && rs.Primary.Attributes["enabled"] == "false" {
			return fmt.Errorf("error getting Service Catalog AWS Organizations Access (%s): wrong setting", rs.Primary.ID)
		}

		return nil
	}
}

func testAccOrganizationsAccessConfig_basic() string {
	return `
resource "aws_servicecatalog_organizations_access" "test" {
  enabled = "true"
}
`
}
