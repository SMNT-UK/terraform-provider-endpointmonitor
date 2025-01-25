# Example of a common shared step of an Android Journey check.
# These can be used to reduce duplication of configuration and also 
# easier maintenance of common actions needed across multiple checks.

resource "endpointmonitor_android_journey_common_step" "example" {
  name        = "Login to Internal App"
  description = "Log test user into internal company app."
  wait_time   = 5000

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