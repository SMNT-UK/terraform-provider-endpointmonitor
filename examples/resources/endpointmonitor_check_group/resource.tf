# Example Check Group that uses the endpointmonitor_app_group data source 
# to get the id of the parent App Group.

data "endpointmonitor_app_group" "example" {
  search = "Public Websites"
}

resource "endpointmonitor_check_group" "example" {
  name            = "Main Company Website"
  description     = "Contains checks for the main company website."
  app_group_id    = endpointmonitor_app_group.example.id
}