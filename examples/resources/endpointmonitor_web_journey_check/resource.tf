# Example Web Journey Check that uses the endpointmonitor_check_host 
# and endpointmonitor_check_group data sources to get the id of the 
# group it sits in and the host it is to run on.
# Additionally uses the endpointmonitor_web_journey_common_step 
# data source to apply a common journey step.

data "endpointmonitor_check_host" "controller" {
  search = "controller"
}

data "endpointmonitor_check_group" "websites" {
  search = "Public Websites"
}

data "endpointmonitor_web_journey_common_step" "initial" {
  search = "Initial Page Checks"
}

resource "endpointmonitor_web_journey_check" "example" {
  name            = "Website Login"
  description     = "Checks website login is working."
  enabled         = true
  check_frequency = 120
  start_url       = "https://www.mycompany.com/"
  trigger_count   = 2

  monitor_domain {
    domain              = "mycompany.com"
    include_sub_domains = true
  }

  monitor_domain {
    domain              = "auth0.com"
    include_sub_domains = true
  }

  step {
    sequence       = 0
    name           = "Initial Page Load Checks"
    type           = "COMMON"
    common_step_id = data.endpointmonitor_web_journey_common_step.initial.id
  }

  step {
    sequence               = 1
    name                   = "Login"
    type                   = "CUSTOM"
    wait_time              = 5000
    page_load_time_warning = 2000
    page_load_time_alert   = 5000

    page_check {
      description = "Check not already logged in."
      type        = "CHECK_FOR_TEXT"

      check_for_text {
        text_to_find = "Logout"
        state        = "ABSENT"
      }
    }

    action {
      sequence        = 0
      description     = "Enter username"
      always_required = true
      type            = "TEXT_INPUT"

      text_input {
        element_id = "login_username"
        input_text = "my.user@mycompany.com"
      }
    }

    action {
      sequence        = 1
      description     = "Enter password"
      always_required = true
      type            = "PASSWORD_INPUT"

      password_input {
        element_id     = "login_password"
        input_password = var.login_password
      }
    }

    action {
      sequence        = 2
      description     = "Click Login"
      always_required = true
      type            = "CLICK"

      click {
        search_text  = "Login"
        element_type = "button"
      }
    }
  }

  step {
    sequence = 2
    name     = "Check Login Successful"
    type     = "CUSTOM"

    page_check {
      description = "Check login response"
      type        = "CHECK_URL_RESPONSE"

      check_url_response {
        url                   = "https://mywebsite.com/login"
        comparison            = "STARTS_WITH"
        warning_response_time = 1500
        alert_response_time   = 3000
        response_code         = 200
      }
    }

    page_check {
      description = "Check username is displayed"
      type        = "CHECK_FOR_TEXT"

      check_for_text {
        text_to_find = "my.user@mycompany.com"
        state        = "PRESENT"
      }
    }
  }

  check_host_id  = data.endpointmonitor_check_host.controller.id
  check_group_id = endpointmonitor_check_group.websites.id
}