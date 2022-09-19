data "endpointmonitor_check_host" "epm_01" {
  search = "smnt-edin-epm01"
}

data "endpointmonitor_check_host" "app_01" {
  search = "smnt-edin-app01"
}

data "endpointmonitor_check_host" "app_02" {
  search = "smnt-edin-app02"
}

data "endpointmonitor_check_group" "integration_tests" {
  search = "Integration Tests"
}