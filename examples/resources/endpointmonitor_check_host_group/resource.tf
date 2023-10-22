data "endpointmonitor_check_hosts" "agents" {
  search = "agents"
}

resource "endpointmonitor_check_host_group" "example" {
  name             = "ECS Agent Cluster"
  description      = "ECS agent cluster for website checks"
  enabled          = true
  check_host_ids   = endpointmonitor_check_hosts.agents.ids
}