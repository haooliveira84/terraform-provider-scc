terraform {
  required_providers {
    scc = {
        source = "sap/cloudconnector"
    }
  }
}

provider "scc" {
  ca_certificate = file("rootCA.pem")
}