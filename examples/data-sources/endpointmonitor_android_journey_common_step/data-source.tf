# Example use of endpointmonitor_android_journey_common_step data source to 
# add a common step to a android journey check.

data "endpointmonitor_android_journey_common_step" "login" {
  search = "Login to Internal App"
}

data "endpointmonitor_check_host" "android" {
  search = "android"
}

data "endpointmonitor_check_group" "apps" {
  search = "Android Apps"
}

resource "endpointmonitor_web_journey_check" "example" {
  name          = "Internal App Login"
  description   = "Checks login is working on internal company app."
  enabled       = true
  trigger_count = 2

  check_host_id  = data.endpointmonitor_check_host.android.id
  check_group_id = data.endpointmonitor_check_group.apps.id

  step {
    sequence       = 1
    type           = "COMMON"
    common_step_id = data.endpointmonitor_android_journey_common_step.login.id
  }
}