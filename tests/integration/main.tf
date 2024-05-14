terraform {
  backend "s3" {
    bucket = "smnt-terraform-states"
    key    = "endpointmonitor/integration/state.tfstate"
  }

  required_providers {
    endpointmonitor = {
      source = "smnt/endpointmonitor"
    }
  }
}

provider "endpointmonitor" {
  url = "https://smnt-edin-epm.net.smnt.co.uk/api"
}

resource "endpointmonitor_dashboard_group" "integration_tests" {
  name        = "Integration Tests"
  description = "Integration test checks. Managed by Terraform."
}

resource "endpointmonitor_check_group" "integration_tests" {
  name               = "Integration Tests"
  description        = "Integration test checks. Managed by Terraform."
  dashboard_group_id = endpointmonitor_dashboard_group.integration_tests.id
}

resource "endpointmonitor_url_check" "integration_test" {
  name        = "Integration URL Test"
  description = "Integration URL Test. Managed by Terraform."

  check_group_id = endpointmonitor_check_group.integration_tests.id
  check_host_id  = data.endpointmonitor_check_host.epm_01.id

  url                    = "https://punchlinetheatre.co.uk/"
  request_method         = "GET"
  allow_redirects        = false
  expected_response_code = 200
  timeout                = 10000
  alert_response_time    = 5000
  warning_response_time  = 3000
  trigger_count          = 3
}

resource "endpointmonitor_dns_check" "integration_test" {
  name               = "Integration DNS Test"
  description        = "Integeration DNS Test. Managed by Terraform."
  check_frequency    = 300
  hostname           = "one.one.one.one"
  expected_addresses = ["1.0.0.1", "1.1.1.1"]
  trigger_count      = 2
  check_host_id      = data.endpointmonitor_check_host.epm_01.id
  check_group_id     = data.endpointmonitor_check_group.integration_tests.id
}

resource "endpointmonitor_ping_check" "integration_test" {
  name            = "Intgration Ping Test"
  description     = "Integration Ping Test. Managed by Terraform."
  check_frequency = 30
  hostname        = "bbc.co.uk"
  trigger_count   = 3
  check_host_id   = data.endpointmonitor_check_host.epm_01.id
  check_group_id  = data.endpointmonitor_check_group.integration_tests.id

  warning_response_time = 2000
  timeout_time          = 5000
}

resource "endpointmonitor_socket_check" "integration_test" {
  name            = "Integration Socket Test"
  description     = "Integration Socket Test. Managed by Terraform."
  check_frequency = 120
  hostname        = "lttstore.co.uk"
  port            = 443
  trigger_count   = 2
  check_host_id   = data.endpointmonitor_check_host.epm_01.id
  check_group_id  = data.endpointmonitor_check_group.integration_tests.id
}

resource "endpointmonitor_web_journey_check" "integration_test" {
  name            = "Integration Web Journey Test"
  description     = "Integration Web Journey Test. Managed by Terraform."
  check_frequency = 120
  enabled         = true
  start_url       = "https://koolness.co.uk/test/"
  trigger_count   = 3

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
    common_step_id = endpointmonitor_web_journey_common_step.initial.id
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

  check_host_id  = data.endpointmonitor_check_host.epm_01.id
  check_group_id = endpointmonitor_check_group.integration_tests.id

  depends_on = [endpointmonitor_web_journey_common_step.initial]
}

resource "endpointmonitor_web_journey_common_step" "initial" {
  name                   = "Initial Page Load Suppressions"
  description            = "Suppresses known errors for initial page load"
  wait_time              = 5000
  page_load_time_warning = 2500
  page_load_time_alert   = 5000

  console_message_suppression {
    description = "Suppress Toastie202 not resolved error"
    log_level   = "ERROR"
    message     = "https://toastie202.co.uk/something.jpg - Failed to load resource: net::ERR_NAME_NOT_RESOLVED"
    comparison  = "EQUALS"
  }

  console_message_suppression {
    description = "Suppress Test Error"
    log_level   = "ERROR"
    message     = "This is a test error, but"
    comparison  = "CONTAINS"
  }

  console_message_suppression {
    description = "Supress CSP Error"
    log_level   = "ERROR"
    message     = "Content Security Policy directive: \"frame-ancestors 'self'\"."
    comparison  = "ENDS_WITH"
  }

  network_suppression {
    description   = "Suppress test/nonexistent.js"
    url           = "https://koolness.co.uk/test/nonexistent.js"
    response_code = 404
    comparison    = "EQUALS"
  }

  network_suppression {
    description   = "Suppress test/thisisnsright.png"
    url           = "test/thisisnsright.png"
    response_code = 404
    comparison    = "ENDS_WITH"
  }

  network_suppression {
    description   = "Suppress smnt.co.uk/someUnknownImage.png"
    url           = "https://smnt.co.uk/someUnknownImage.png"
    response_code = 404
    comparison    = "EQUALS"
  }
}