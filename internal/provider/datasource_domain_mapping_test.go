package provider

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceDomainMapping(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_domain_mapping")
		rec.SetRealTransport(&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		})
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getTestProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: providerConfig("", user) + DataSourceDomainMapping("mappings", "cf.eu12.hana.ondemand.com", "0bcb0012-a982-42f9-bda4-0a5cb15f88c8", "testterraforminternaldomain"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.cloudconnector_domain_mapping.mappings", "region_host", "cf.eu12.hana.ondemand.com"),
						resource.TestMatchResourceAttr("data.cloudconnector_domain_mapping.mappings", "subaccount", regexpValidUUID),
						resource.TestCheckResourceAttr("data.cloudconnector_domain_mapping.mappings", "virtual_domain", "testterraformvirtualdomain"),
						resource.TestCheckResourceAttr("data.cloudconnector_domain_mapping.mappings", "internal_domain", "testterraforminternaldomain"),
					),
				},
			},
		})

	})

}

func DataSourceDomainMapping(datasourceName string, regionHost string, subaccountID string, internalDomain string) string {
	return fmt.Sprintf(`
	data "cloudconnector_domain_mapping" "%s" {
	region_host= "%s"
    subaccount= "%s"
	internal_domain= "%s"
	}
	`, datasourceName, regionHost, subaccountID, internalDomain)
}
