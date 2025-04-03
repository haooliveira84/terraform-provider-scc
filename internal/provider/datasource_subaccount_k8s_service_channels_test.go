package provider

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountK8SServiceChannels(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_k8s_service_channels")
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
					Config: providerConfig("", user) + DataSourceSubaccountK8SServiceChannels("cc_scs", "cf.eu12.hana.ondemand.com", "0bcb0012-a982-42f9-bda4-0a5cb15f88c8"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_k8s_service_channels.cc_scs", "region_host", "cf.eu12.hana.ondemand.com"),
						resource.TestMatchResourceAttr("data.cloudconnector_subaccount_k8s_service_channels.cc_scs", "subaccount", regexpValidUUID),

						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_k8s_service_channels.cc_scs", "subaccount_k8s_service_channels.#", "1"),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_k8s_service_channels.cc_scs", "subaccount_k8s_service_channels.0.k8s_cluster", "cp.da2b3e1.stage.kyma.ondemand.com:443"),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_k8s_service_channels.cc_scs", "subaccount_k8s_service_channels.0.k8s_service", "bd64665f-060a-47b6-8aba-f406703f0acf"),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_k8s_service_channels.cc_scs", "subaccount_k8s_service_channels.0.port", "8080"),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_k8s_service_channels.cc_scs", "subaccount_k8s_service_channels.0.connections", "1"),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_k8s_service_channels.cc_scs", "subaccount_k8s_service_channels.0.type", "K8S"),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_k8s_service_channels.cc_scs", "subaccount_k8s_service_channels.0.enabled", "true"),
						resource.TestCheckResourceAttrSet("data.cloudconnector_subaccount_k8s_service_channels.cc_scs", "subaccount_k8s_service_channels.0.id"),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_k8s_service_channels.cc_scs", "subaccount_k8s_service_channels.0.state.connected", "true"),
						resource.TestMatchResourceAttr("data.cloudconnector_subaccount_k8s_service_channels.cc_scs", "subaccount_k8s_service_channels.0.state.connected_since_time_stamp", regexp.MustCompile(`^(0|\d{13})$`)),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_k8s_service_channels.cc_scs", "subaccount_k8s_service_channels.0.state.opened_connections", "1"),
					),
				},
			},
		})

	})

}

func DataSourceSubaccountK8SServiceChannels(datasourceName string, regionHost string, subaccountID string) string {
	return fmt.Sprintf(`
	data "cloudconnector_subaccount_k8s_service_channels" "%s" {
	region_host = "%s"
	subaccount = "%s"
	}
	`, datasourceName, regionHost, subaccountID)
}
