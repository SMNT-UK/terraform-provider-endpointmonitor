# Example Android Journey Check that uses the endpointmonitor_check_host 
# and endpointmonitor_check_group data sources to get the id of the 
# group it sits in and the host it is to run on.
# Additionally uses the endpointmonitor_android_journey_common_step 
# data source to apply a common journey step.

data "endpointmonitor_check_host_group" "agent_group" {
  search = "android"
}

data "endpointmonitor_check_group" "apps" {
  search = "Internal Apps"
}

data "endpointmonitor_android_journey_common_step" "initial" {
  search = "Accept T&Cs and Privacy Agreement"
}

resource "endpointmonitor_android_journey_check" "example" {
  name            = "Company Internal App"
  description     = "Checks company internal app is working."
  enabled         = true
  check_frequency = 120
  trigger_count   = 2

  apk = filebase64("${path.module}/company_app.apk")

  common_step {
    sequence       = 1
    name           = "Initial App Load"
    common_step_id = data.endpointmonitor_android_journey_common_step.initial.id
  }

  custom_step {
    sequence               = 2
    name                   = "Login"
    wait_time              = 5000

    step_check {
      description = "Check not already logged in."
      type        = "CHECK_FOR_TEXT"

      check_for_text {
        text_to_find = "Logout"
        state        = "ABSENT"
      }
    }

    step_interaction {
      sequence        = 1
      description     = "Enter username"
      always_required = true
      type            = "INPUT_TEXT"

      text_input {
        component_id = "login_username"
        input_text   = "my.user@mycompany.com"
      }
    }

    step_interaction {
      sequence        = 2
      description     = "Enter password"
      always_required = true
      type            = "INPUT_PASSWORD"

      password_input {
        component_id   = "login_password"
        input_password = var.login_password
      }
    }

    step_interaction {
      sequence        = 3
      description     = "Click Login"
      always_required = true
      type            = "CLICK"

      click {
        search_text  = "Login"
        element_type = "button"
      }
    }
  }

  custom_step {
    sequence = 2
    name     = "Check Login Successful"

    step_check {
      description = "Check logout button showing."
      type        = "CHECK_FOR_TEXT"

      check_for_text {
        text_to_find = "Logout"
        state        = "PRESENT"
      }
    }

    step_check {
      description = "Check username is displayed"
      type        = "CHECK_FOR_TEXT"

      check_for_text {
        text_to_find = "my.user@mycompany.com"
        state        = "PRESENT"
      }
    }
  }

  host_group_id  = data.endpointmonitor_check_host_group.agent_group.id
  check_group_id = data.endpointmonitor_check_group.apps.id
}