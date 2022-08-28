resource "endpointmonitor_check_host" "example" {
  name             = "Internal Applications Agent"
  description      = "Agent for monitoring internal facing applications."
  hostname         = "machinename.fqdn.com"
  type             = "AGENT"
  enabled          = true
  max_checks       = 1
  send_check_files = true
}