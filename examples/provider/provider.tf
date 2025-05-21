terraform {
  required_providers {
    scc = {
        source = "sap/scc"
        version = "0.1.0-beta1"
    }
  }
}

provider "scc" {
  username = "Username"
  password = "Password"
  instance_url = "https://localhost:port"
  ca_certificate = file("PATH_TO_FILE/user-cert.pem")
}