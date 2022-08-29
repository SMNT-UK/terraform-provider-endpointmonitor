# Example use of endpointmonitor_check_groups to grab a list of check 
# group ids to add to a maintenance period.

data "endpointmonitor_check_groups" "websites" {
  search = "Website"
}

resource "endpointmonitor_maintenance_period" "example" {
  description     = "Suppress alerts during backup period."
  enabled         = true
  day_of_week     = "ALL"
  start_time      = "01:00"
  end_time        = "03:00"
  check_group_ids = data.endpointmonitor_check_groups.websites.ids
}