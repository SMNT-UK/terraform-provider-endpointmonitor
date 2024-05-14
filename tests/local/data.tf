data "endpointmonitor_proxy_hosts" "test" {
  search = ""
}

data "endpointmonitor_proxy_host" "test" {
  search = "VLAN10"
}

data "endpointmonitor_dashboard_groups" "test" {
  search = ""
}

data "endpointmonitor_dashboard_group" "test" {
  search = "Terraform"
}

data "endpointmonitor_check_groups" "test" {
  search = ""
}

data "endpointmonitor_check_group" "test" {
  search = "Terraform"
}

data "endpointmonitor_check_hosts" "test" {
  search = ""
}

data "endpointmonitor_check_host" "test" {
  search = "Terraform"
}

data "endpointmonitor_web_journey_common_step" "test" {
  search = "Test Common WebJourney Step"
}

data "endpointmonitor_web_journey_common_steps" "test" {
  search = "Test Common"
}