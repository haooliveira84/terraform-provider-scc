package provider

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSystemMappingResource(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_system_mapping_resource")
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
					Config: providerConfig(user) + DataSourceSystemMappingResource("test", "cf.eu12.hana.ondemand.com", "0bcb0012-a982-42f9-bda4-0a5cb15f88c8", "testterraformvirtual", "900", "/google.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping_resource.test", "region_host", "cf.eu12.hana.ondemand.com"),
						resource.TestMatchResourceAttr("data.cloudconnector_system_mapping_resource.test", "subaccount", regexpValidUUID),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping_resource.test", "virtual_host", "testterraformvirtual"),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping_resource.test", "virtual_port", "900"),

						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping_resource.test", "id", "/google.com"),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping_resource.test", "enabled", "true"),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping_resource.test", "exact_match_only", "true"),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping_resource.test", "websocket_upgrade_allowed", "false"),
						resource.TestMatchResourceAttr("data.cloudconnector_system_mapping_resource.test", "creation_date", regexValidTimeStamp),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping_resource.test", "description", ""),
					),
				},
			},
		})

	})

}

func DataSourceSystemMappingResource(datasourceName string, regionHost string, subaccount string, virtualHost string, virtualPort string, systemMappingResource string) string {
	return fmt.Sprintf(`
	data "cloudconnector_system_mapping_resource" "%s" {
	region_host= "%s"
	subaccount= "%s"
	virtual_host= "%s"
	virtual_port= "%s"
	id= "%s"
	}
	`, datasourceName, regionHost, subaccount, virtualHost, virtualPort, systemMappingResource)
}
