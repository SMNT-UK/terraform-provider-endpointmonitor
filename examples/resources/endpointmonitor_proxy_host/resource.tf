resource "endpointmonitor_proxy_host" "test" {
  name        = "Internal HTTP Proxy"
  description = "HTTP proxy host used for internet access"
  hostname    = "squid-01.internal.mycompany.com"
  port        = 3128
}