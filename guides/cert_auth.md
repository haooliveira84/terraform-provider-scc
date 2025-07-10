### <u> X.509 Certificate Authentication </u>

You would require a valid **X.509 Client Certificate** and the corresponding **Client Key** of an [Administrator](https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/logon-to-cloud-connector-via-client-certificate) to get authenticated.
 
You can only configure the credentials as part of the provider configuration as shown below:

 ```hcl
provider "scc" {
    instance_url = <your_instance_url>
    client_certificate = <your_client_certificate>
    client_key = <your_client_key>
}
```

Ensure to paste the ***content*** of your client certificate rather than the ***file path***.
You can even use the function `file("path_to_certificate.pem")` to load the file content. 
