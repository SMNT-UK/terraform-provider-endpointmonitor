# Example DNS Check that uses the endpointmonitor_check_host and
# endpointmonitor_check_group data sources to get the id of the 
# group it sits in and the host it is to run on.

data "endpointmonitor_check_host" "controller" {
  search = "controller"
}

data "endpointmonitor_check_group" "dns" {
  search = "DNS Checks"
}

resource "endpointmonitor_dns_check" "example" {
  name               = "Loadbalancer Failure Check"
  description        = "Checks the load balancer hasn't automatically failed over."
  hostname           = "primary-website.mycompany.com"
  expected_addresses = ["1.2.3.4", "1.2.3.5"]
  trigger_count      = 3
  check_host_id      = data.endpointmonitor_check_host.controller.id
  check_group_id     = data.endpointmonitor_check_group.dns.id
}