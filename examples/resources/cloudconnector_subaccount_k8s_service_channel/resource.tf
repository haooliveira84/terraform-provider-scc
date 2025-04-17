resource "cloudconnector_subaccount_k8s_service_channel" "cc_sc" {
  region_host = "cf.eu12.hana.ondemand.com"
  subaccount = "12345678-90ab-cdef-1234-567890abcdef"
  k8s_cluster =  "cp.app.cluster.kyma.ondemand.com"
  k8s_service =  "12345678-90ab-cdef-1234-567890abcdef"
  port = 3000
  connections = 1
  enabled = true
}