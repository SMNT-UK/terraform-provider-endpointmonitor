# Example use of endpointmonitor_check data source to grab the
# id of a check for a Maintenance Period.

data "endpointmonitor_check" "example" {
  search = "Website Homepage"
}

resource "endpointmonitor_maintenance_period" "example" {
  description = "Suppress alerts during backup period."
  enabled     = true
  day_of_week = "ALL"
  start_time  = "01:00"
  end_time    = "03:00"
  check_ids   = [data.endpointmonitor_check.example.id]
}