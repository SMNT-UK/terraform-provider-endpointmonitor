# Example Socket Check that uses the endpointmonitor_check_host and
# endpointmonitor_check_group data sources to get the id of the 
# group it sits in and the host it is to run on.

data "endpointmonitor_check_host" "controller" {
  search = "controller"
}

data "endpointmonitor_check_group" "database" {
  search = "Database Checks"
}

resource "endpointmonitor_socket_check" "example" {
  name           = "Primary Database Listening"
  hostname       = "db01.internal.mycompany.com"
  port           = 5432
  trigger_count  = 2
  check_host_id  = data.endpointmonitor_check_host.controller.id
  check_group_id = data.endpointmonitor_check_group.database.id
}