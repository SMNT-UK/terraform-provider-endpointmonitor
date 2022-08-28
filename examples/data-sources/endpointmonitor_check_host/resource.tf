# Example use of endpointmonitor_check_host to find the check 
# host id to add to a URL Check.

data "endpointmonitor_check_host" "controller" {
  search = "controller"
}

data "endpointmonitor_check_group" "websites" {
  search = "Website Checks"
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

  check_host_id  = data.endpointmonitor_check_host.controller.id
  check_group_id = data.endpointmonitor_check_group.websites.id
}