# Example Certificate Check that uses the endpointmonitor_check_host
# and endpointmonitor_check_group data sources to get the id of the
# group it sits in and the host it is to run on.

data "endpointmonitor_check_host" "controller" {
  search = "controller"
}

data "endpointmonitor_check_group" "websites" {
  search = "Public Websites"
}

resource "endpointmonitor_certificate_check" "example" {
  name                   = "Website Certificate"
  description            = "Website certificate expiry check"
  check_frequency        = 3600
  trigger_count          = 2

  url                    = "https://www.mycompany.com/"
  warning_days_remaining = 7
  alert_days_remaining   = 2

  check_host_id  = data.endpointmonitor_check_host.controller.id
  check_group_id = data.endpointmonitor_check_group.websites.id
}