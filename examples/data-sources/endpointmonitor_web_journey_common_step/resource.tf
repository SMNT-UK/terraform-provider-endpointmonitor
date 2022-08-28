# Example use of endpointmonitor_web_journey_common_step data source to 
# add a common step to a web journey check.

data "endpointmonitor_web_journey_common_step" "cookies" {
  search = "Cookie Prompt"
}

data "endpointmonitor_check_host" "controller" {
  search = "controller"
}

data "endpointmonitor_check_group" "website" {
  search = "Website Checks"
}

resource "endpointmonitor_web_journey_check" "example" {
  name          = "Website Login"
  description   = "Checks website login is working."
  enabled       = true
  start_url     = "https://www.mycompany.com/"
  trigger_count = 2

  step {
    sequence       = 0
    name           = "Accept Cookie Prompt"
    type           = "COMMON"
    common_step_id = data.endpointmonitor_web_journey_common_step.cookies.id
  }

  check_host_id  = data.endpointmonitor_check_host.controller.id
  check_group_id = endpointmonitor_check_group.website.id
}