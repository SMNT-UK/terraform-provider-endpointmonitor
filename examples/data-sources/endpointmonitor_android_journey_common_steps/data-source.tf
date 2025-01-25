# Example use of endpointmonitor_android_journey_common_steps data source to 
# add a common step to a android journey check.

data "endpointmonitor_android_journey_common_steps" "initial" {
  search = "Initial check"
}

data "endpointmonitor_check_host" "android" {
  search = "android"
}

data "endpointmonitor_check_group" "apps" {
  search = "Android App Checks"
}

resource "endpointmonitor_android_journey_check" "example" {
  name          = "Company Internal App Loading"
  description   = "Checks company internal app is loading correctly."
  enabled       = true
  trigger_count = 2

  apk = filebase64("${path.module}/company_internal_app.apk")

  dynamic "step" {
    for_each = data.endpointmonitor_android_journey_common_steps.all.ids

    content {
      sequence       = step.key
      type           = "COMMON"
      common_step_id = step.value
    }
  }

  check_host_id  = data.endpointmonitor_check_host.android.id
  check_group_id = data.endpointmonitor_check_group.apps.id
}