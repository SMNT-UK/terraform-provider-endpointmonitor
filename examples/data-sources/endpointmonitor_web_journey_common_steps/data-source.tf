# Example use of endpointmonitor_web_journey_common_steps data source to 
# add all common steps to a web journey check.
# Can think of no good reason to do this, but best example I could think of.

data "endpointmonitor_web_journey_common_steps" "all" {
  search = ""
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

  dynamic "step" {
    for_each = data.endpointmonitor_web_journey_common_steps.all.ids

    content {
      sequence       = step.key
      name           = "Common Step ${step.key}"
      type           = "COMMON"
      common_step_id = step.value
    }
  }

  check_host_id  = data.endpointmonitor_check_host.controller.id
  check_group_id = endpointmonitor_check_group.website.id
}