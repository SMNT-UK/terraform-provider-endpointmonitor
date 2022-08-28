# Example Ping Check that uses the endpointmonitor_check_host and
# endpointmonitor_check_group data sources to get the id of the 
# group it sits in and the host it is to run on.

data "endpointmonitor_check_host" "controller" {
  search = "controller"
}

data "endpointmonitor_check_group" "network" {
  search = "Network Checks"
}

resource "endpointmonitor_ping_check" "example" {
  name           = "Ping Access Point 1"
  hostname       = "ap01.internal.mycompany.com"
  trigger_count  = 3
  check_host_id  = data.endpointmonitor_check_host.controller.id
  check_group_id = data.endpointmonitor_check_group.network.id

  warning_response_time = 2000
  timeout_time          = 5000
}