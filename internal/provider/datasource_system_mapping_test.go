package provider

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSystemMapping(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_system_mapping")
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
					Config: providerConfig("", user) + DataSourceSystemMapping("test", "cf.eu12.hana.ondemand.com", "0bcb0012-a982-42f9-bda4-0a5cb15f88c8", "testterraformvirtual", "900"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "region_host", "cf.eu12.hana.ondemand.com"),
						resource.TestMatchResourceAttr("data.cloudconnector_system_mapping.test", "subaccount", regexpValidUUID),

						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "virtual_host", "testterraformvirtual"),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "virtual_port", "900"),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "local_host", "testterraforminternal"),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "local_port", "900"),
						resource.TestMatchResourceAttr("data.cloudconnector_system_mapping.test", "creation_date", regexValidTimeStamp),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "protocol", "HTTP"),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "backend_type", "abapSys"),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "authentication_mode", "KERBEROS"),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "host_in_header", "VIRTUAL"),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "sid", ""),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "total_resources_count", "1"),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "enabled_resources_count", "1"),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "description", ""),
						resource.TestCheckResourceAttr("data.cloudconnector_system_mapping.test", "sap_router", ""),
					),
				},
			},
		})

	})

}

func DataSourceSystemMapping(datasourceName string, regionHost string, subaccount string, virtualHost string, virtualPort string) string {
	return fmt.Sprintf(`
	data "cloudconnector_system_mapping" "%s" {
	region_host= "%s"
	subaccount= "%s"
	virtual_host= "%s"
	virtual_port= "%s"	
	}
	`, datasourceName, regionHost, subaccount, virtualHost, virtualPort)
}
