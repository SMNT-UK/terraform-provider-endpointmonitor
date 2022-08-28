# Example URL Check that uses the endpointmonitor_check_host and
# endpointmonitor_check_group data sources to get the id of the 
# group it sits in and the host it is to run on.

data "endpointmonitor_check_host" "controller" {
  search = "controller"
}

data "endpointmonitor_check_group" "websites" {
  search = "Public Websites"
}

resource "endpointmonitor_url_check" "example" {
  name                   = "Home Page Check"
  description            = "Check home page loads as expected"
  url                    = "https://www.mycompany.com/"
  trigger_count          = 2
  request_method         = "GET"
  expected_response_code = 200
  alert_response_time    = 5000
  warning_response_time  = 3000
  timeout                = 10000
  allow_redirects        = false

  request_header {
    name  = "Agent"
    value = "EndPoint Monitor"
  }

  check_host_id  = data.endpointmonitor_check_host.controller.id
  check_group_id = data.endpointmonitor_check_group.websites.id
}