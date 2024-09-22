data "endpointmonitor_check" "test" {
  search = "koolness.co.uk/test"
}

data "endpointmonitor_checks" "test" {
  search = "Development"
}

data "endpointmonitor_check_host_group" "test" {
  search = "Test Host Group"
}

data "endpointmonitor_check_host_groups" "test" {
  search = "Terraform"
}

data "endpointmonitor_maintenance_period" "test" {
  search = "Period Test"
}

data "endpointmonitor_maintenance_periods" "test" {
  search = "Terraform"
}

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
  search = "test.terraform.host2"
}

data "endpointmonitor_web_journey_common_step" "test" {
  search = "Initial Page Load"
}

data "endpointmonitor_web_journey_common_steps" "test" {
  search = "Test Common"
}