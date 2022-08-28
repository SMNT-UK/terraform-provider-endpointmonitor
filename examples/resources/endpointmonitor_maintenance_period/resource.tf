# Example Maintenance Period that suppresses alerts between 01h and 03h every night 
# for all checks with 'Website' in their name or description.

data "endpointmonitor_checks" "websites" {
  search = "Website"
}

resource "endpointmonitor_maintenance_period" "example" {
  description = "Suppress alerts during backup period."
  enabled     = true
  day_of_week = "ALL"
  start_time  = "01:00"
  end_time    = "03:00"
  check_ids   = data.endpointmonitor_checks.websites.ids
}