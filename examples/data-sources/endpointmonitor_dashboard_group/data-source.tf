# Example Check Group that uses the endpointmonitor_dashboard_group 
# data source to get the id of the parent Dashboard Group.

data "endpointmonitor_dashboard_group" "example" {
  search = "Public Websites"
}

resource "endpointmonitor_check_group" "example" {
  name               = "Main Company Website"
  description        = "Contains checks for the main company website."
  check_frequency    = 60
  dashboard_group_id = data.endpointmonitor_dashboard_group.example.id
}