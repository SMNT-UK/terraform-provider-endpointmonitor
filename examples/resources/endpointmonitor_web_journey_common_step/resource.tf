# Example Web Journey Common Step of setting up a common set of checks, 
# actions and suppressions that could be used at the start of each 
# web journey check for a site to remove duplication and make changes 
# accross multiple checks easier to accomplish.

resource "endpointmonitor_web_journey_common_step" "example" {
  name                   = "Initial Page Load"
  description            = "Runs check and suppressions for initial website page load."
  wait_time              = 10000
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

  console_message_suppression {
    description = "Ignore error from 3rd party marketing script."
    log_level   = "ERROR"
    comparison  = "STARTS_WITH"
    message     = "Generic Marketing Compnay encountered an error when"
  }

  network_suppression {
    description      = "Ignore errors from advertising provider."
    url              = "https://api.advertisingpeople.com/"
    comparison       = "STARTS_WITH"
    any_client_error = true
    any_server_error = true
  }

  action {
    sequence        = 0
    description     = "Close any marketing pop-up."
    always_required = false
    type            = "CLICK"

    click {
      xpath = "//*[@class='aria-close']"
    }
  }

  action {
    sequence        = 1
    description     = "Accept Cookies"
    always_required = true
    type            = "CLICK"

    click {
      search_text  = "Accept All"
      element_type = "button"
    }
  }
}