terraform {
  required_providers {
    cloudconnector = {
        source = "SAP/cloudconnector"
        version = "0.1.0-beta1"
    }
  }
}

provider "cloudconnector" {
  username = "Username"
  password = "Password"
  instance_url = "https://localhost:port"
  ca_certificate = file("PATH_TO_FILE/user-cert.pem")
}