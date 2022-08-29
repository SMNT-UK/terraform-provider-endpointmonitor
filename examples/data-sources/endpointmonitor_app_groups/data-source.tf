# Example Check Group that uses the endpointmonitor_app_groups data source 
# to set up a default Check Group for each App Group found.

data "endpointmonitor_app_groups" "example" {
  search = "Websites"
}

resource "endpointmonitor_check_group" "example" {
  count = length(endpointmonitor_app_group.example.ids)

  name            = "Default Group"
  description     = "Default group."
  check_frequency = 60
  app_group_id    = endpointmonitor_app_group.example.ids[count.index]
}