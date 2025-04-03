package provider

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountServiceChannelK8S(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_channel_k8s")
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
					Config: providerConfig("", user) + DataSourceSubaccountServiceChannel("cc_sc", "cf.eu12.hana.ondemand.com", "0bcb0012-a982-42f9-bda4-0a5cb15f88c8", 183),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_service_channel_k8s.cc_sc", "credentials.region_host", "cf.eu12.hana.ondemand.com"),
						resource.TestMatchResourceAttr("data.cloudconnector_subaccount_service_channel_k8s.cc_sc", "credentials.subaccount", regexpValidUUID),

						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_service_channel_k8s.cc_sc", "subaccount_service_channel_k8s.k8s_cluster", "cp.da2b3e1.stage.kyma.ondemand.com:443"),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_service_channel_k8s.cc_sc", "subaccount_service_channel_k8s.k8s_service", "bd64665f-060a-47b6-8aba-f406703f0acf"),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_service_channel_k8s.cc_sc", "subaccount_service_channel_k8s.port", "8080"),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_service_channel_k8s.cc_sc", "subaccount_service_channel_k8s.connections", "1"),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_service_channel_k8s.cc_sc", "subaccount_service_channel_k8s.type", "K8S"),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_service_channel_k8s.cc_sc", "subaccount_service_channel_k8s.enabled", "true"),
						resource.TestCheckResourceAttrSet("data.cloudconnector_subaccount_service_channel_k8s.cc_sc", "subaccount_service_channel_k8s.id"),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_service_channel_k8s.cc_sc", "subaccount_service_channel_k8s.state.connected", "true"),
						resource.TestMatchResourceAttr("data.cloudconnector_subaccount_service_channel_k8s.cc_sc", "subaccount_service_channel_k8s.state.connected_since_time_stamp", regexp.MustCompile(`^(0|\d{13})$`)),
						resource.TestCheckResourceAttr("data.cloudconnector_subaccount_service_channel_k8s.cc_sc", "subaccount_service_channel_k8s.state.opened_connections", "1"),
					),
				},
			},
		})

	})

}

func DataSourceSubaccountServiceChannel(datasourceName string, regionHost string, subaccountID string, id int64) string {
	return fmt.Sprintf(`
	data "cloudconnector_subaccount_service_channel_k8s" "%s" {
	credentials = {
		region_host = "%s"
		subaccount = "%s"
	}
	subaccount_service_channel_k8s = {
		id = "%d"
	}
	}
	`, datasourceName, regionHost, subaccountID, id)
}
