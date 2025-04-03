package provider

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceSubaccountServiceChannelK8S(t *testing.T) {

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_channel_k8s")
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
					Config: providerConfig("", user) + ResourceSubaccountServiceChannelK8S("test", "cf.eu12.hana.ondemand.com", "d3bbbcd7-d5e0-483b-a524-6dee7205f8e8", "cp.ace9fb5.stage.kyma.ondemand.com:443", "29d4e6f6-8e7f-4882-b434-21a52bb75e0f", 6000, 1, true),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("cloudconnector_subaccount_service_channel_k8s.test", "credentials.region_host", "cf.eu12.hana.ondemand.com"),
						resource.TestMatchResourceAttr("cloudconnector_subaccount_service_channel_k8s.test", "credentials.subaccount", regexpValidUUID),
						resource.TestCheckResourceAttr("cloudconnector_subaccount_service_channel_k8s.test", "subaccount_service_channel_k8s.k8s_cluster", "cp.ace9fb5.stage.kyma.ondemand.com:443"),
						resource.TestCheckResourceAttr("cloudconnector_subaccount_service_channel_k8s.test", "subaccount_service_channel_k8s.k8s_service", "29d4e6f6-8e7f-4882-b434-21a52bb75e0f"),
						resource.TestCheckResourceAttr("cloudconnector_subaccount_service_channel_k8s.test", "subaccount_service_channel_k8s.port", "6000"),
						resource.TestCheckResourceAttr("cloudconnector_subaccount_service_channel_k8s.test", "subaccount_service_channel_k8s.connections", "1"),
						resource.TestCheckResourceAttr("cloudconnector_subaccount_service_channel_k8s.test", "subaccount_service_channel_k8s.type", "K8S"),
						resource.TestCheckResourceAttr("cloudconnector_subaccount_service_channel_k8s.test", "subaccount_service_channel_k8s.enabled", "true"),
						resource.TestCheckResourceAttrSet("cloudconnector_subaccount_service_channel_k8s.test", "subaccount_service_channel_k8s.id"),
						resource.TestCheckResourceAttr("cloudconnector_subaccount_service_channel_k8s.test", "subaccount_service_channel_k8s.state.connected", "true"),
						resource.TestMatchResourceAttr("cloudconnector_subaccount_service_channel_k8s.test", "subaccount_service_channel_k8s.state.connected_since_time_stamp", regexp.MustCompile(`^(0|\d{13})$`)),
						resource.TestCheckResourceAttr("cloudconnector_subaccount_service_channel_k8s.test", "subaccount_service_channel_k8s.state.opened_connections", "1"),
					),
				},
			},
		})

	})

}

func ResourceSubaccountServiceChannelK8S(datasourceName string, regionHost string, subaccount string, k8sCluster string, k8sService string, port int64, connections int64, enabled bool) string {
	return fmt.Sprintf(`
	resource "cloudconnector_subaccount_service_channel_k8s" "%s" {
	credentials = {
		region_host = "%s"
		subaccount = "%s"
	}
	subaccount_service_channel_k8s = {
		k8s_cluster =  "%s",
		k8s_service =  "%s",
		port = "%d",
		connections = "%d"
		enabled = "%t"
	}
	}
	`, datasourceName, regionHost, subaccount, k8sCluster, k8sService, port, connections, enabled)
}
