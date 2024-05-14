# Example Check Group that uses the endpointmonitor_dashboard_groups data source 
# to set up a default Check Group for each Dashboard Group found.

data "endpointmonitor_dashboard_groups" "example" {
  search = "Websites"
}

resource "endpointmonitor_check_group" "example" {
  count = length(endpointmonitor_dashboard_group.example.ids)

  name               = "Default Group"
  description        = "Default group."
  check_frequency    = 60
  dashboard_group_id = endpointmonitor_dashboard_group.example.ids[count.index]
}