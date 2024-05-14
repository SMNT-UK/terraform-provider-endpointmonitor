# Example use of endpointmonitor_check_hosts to setup a heartbeat
# check for each host with agent in the name.

data "endpointmonitor_check_hosts" "agents" {
  search = "agent"
}

data "endpointmonitor_check_group" "agent" {
  search = "Agent Checks"
}

resource "endpointmonitor_url_check" "example" {
  count = length(data.endpointmonitor_check_hosts.agents)

  name                   = "Agent Heartbeat"
  description            = "Agent heartbeat check to check it can communicate with controller."
  url                    = "https://endpointmonitor.internal.mycompany.com/"
  trigger_count          = 2
  request_method         = "GET"
  expected_response_code = 200
  alert_response_time    = 5000
  warning_response_time  = 3000
  timeout                = 10000

  check_host_id  = data.endpointmonitor_check_hosts.agents.id[0]
  check_group_id = data.endpointmonitor_check_group.agent.id
}