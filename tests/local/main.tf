terraform {
  backend "s3" {
    region = "eu-central-1"
    bucket = "smnt-terraform-states"
    key    = "endpointmonitor/local/state.tfstate"
  }

  required_providers {
    endpointmonitor = {
      source = "smnt/endpointmonitor"
    }
  }
}

provider "endpointmonitor" {
  url = "http://localhost:8080/api"
  key = "5c4f665d-c6ec-4ab7-8bf2-a333af82c182" # This is safe for leaving in, it's just a key on a local dev instance.
}

resource "endpointmonitor_socket_check" "test" {
  name            = "Local Terraform Socket Test"
  description     = "Local Terraform Socket Test. Managed by Terraform."
  check_frequency = 120
  hostname        = "lttstore.co.uk"
  port            = 443
  trigger_count   = 2
  check_host_id   = endpointmonitor_check_host.test.id
  check_group_id  = endpointmonitor_check_group.test.id
}

resource "endpointmonitor_proxy_host" "test" {
  name        = "Test Terraform ProxyHost"
  description = "Testing Terraform Descrition"
  hostname    = "test-terraform-host"
  port        = 3318
}

resource "endpointmonitor_check_host" "test" {
  description      = "Terraform CheckHost Description"
  hostname         = "test.terraform.host2"
  type             = "AGENT"
  enabled          = true
  max_checks       = 1
  send_check_files = true
}

resource "endpointmonitor_dashboard_group" "test" {
  name        = "Terraform Dashboard Group 2"
  description = "Terraform Descrition of an Dashboard Group"
}

resource "endpointmonitor_check_group" "test" {
  name               = "Terraform Check Group 2"
  description        = "Terraform description of a check group"
  dashboard_group_id = endpointmonitor_dashboard_group.test.id
}

resource "endpointmonitor_url_check" "test" {
  name                   = "Terraform URL Check"
  description            = "Terraform URL check descrtiption"
  url                    = "https://www.bbc.co.uk/weather"
  trigger_count          = 2
  check_frequency        = 60
  request_method         = "GET"
  expected_response_code = 200
  alert_response_time    = 5000
  warning_response_time  = 3000
  timeout                = 10000
  allow_redirects        = true

  request_header {
    name  = "Agent"
    value = "EndPoint Monitor"
  }

  request_header {
    name  = "Content-Type"
    value = "application/json"
  }
  check_host_id  = data.endpointmonitor_check_host.test.id
  check_group_id = endpointmonitor_check_group.test.id
}

resource "endpointmonitor_dns_check" "test" {
  name               = "Terraform DNS Check"
  description        = "Terraform DNS check descrtiption"
  check_frequency    = 300
  hostname           = "smnt-edin-sql.net.smnt.co.uk"
  expected_addresses = ["10.20.0.31", "10.20.0.32"]
  trigger_count      = 5
  check_host_id      = endpointmonitor_check_host.test.id
  check_group_id     = endpointmonitor_check_group.test.id
}

resource "endpointmonitor_certificate_check" "test" {
  name            = "Terraform Certificate Check"
  description     = "Terraform Certificate Check Test"
  trigger_count   = 2
  check_frequency = 60

  url                    = "https://epm.smnt.co.uk/"
  warning_days_remaining = "7"
  alert_days_remaining   = "2"

  check_host_id  = endpointmonitor_check_host.test.id
  check_group_id = endpointmonitor_check_group.test.id
}

resource "endpointmonitor_web_journey_check" "test" {
  name                 = "Terraform WebJourney Check"
  description          = "Terraform WebJourney check descrtiption"
  enabled              = false
  check_frequency      = 60
  maintenance_override = false
  start_url            = "https://koolness.co.uk/test"
  trigger_count        = 2

  monitor_domain {
    domain              = "koolness.co.uk"
    include_sub_domains = true
  }
  monitor_domain {
    domain              = "smnt.co.uk"
    include_sub_domains = true
  }

  step {
    sequence       = 0
    name           = "Terraform test step 0"
    type           = "COMMON"
    common_step_id = data.endpointmonitor_web_journey_common_step.test.id
  }

  step {
    sequence               = 1
    name                   = "Terraform Test step 1"
    type                   = "CUSTOM"
    wait_time              = 10000
    page_load_time_warning = 2000
    page_load_time_alert   = 5000

    page_check {
      description = "Terraform Test Step 1 Check 1"
      type        = "CHECK_FOR_TEXT"

      check_for_text {
        text_to_find = "Testing Text to Find"
        state        = "PRESENT"
      }
    }
  }

  step {
    sequence               = 2
    name                   = "Terraform Test step 2"
    type                   = "CUSTOM"
    wait_time              = 10000
    page_load_time_warning = 2000
    page_load_time_alert   = 5000

    action {
      sequence        = 1
      description     = "Terraform Test Step 2 Action 1"
      always_required = true
      type            = "PASSWORD_INPUT"

      password_input {
        element_id     = "login_password"
        input_password = "thisshouldbeavariable"
      }
    }
  }

  check_host_id  = data.endpointmonitor_check_hosts.test.ids[0]
  check_group_id = endpointmonitor_check_group.test.id
}

resource "endpointmonitor_maintenance_period" "test" {
  description = "Terraform Maintenance Period Test"
  enabled     = true
  day_of_week = "ALL"
  start_time  = "01:00"
  end_time    = "03:00"
}

resource "endpointmonitor_web_journey_common_step" "test" {
  name                   = "Test Common Step"
  description            = "Generic Test Common Step"
  wait_time              = 10000
  page_load_time_warning = 2000
  page_load_time_alert   = 5000

  page_check {
    description = "Terraform Test Common Step Check 1"
    type        = "CHECK_FOR_TEXT"

    check_for_text {
      text_to_find = "Testing Text to Find"
      state        = "PRESENT"
    }
  }
}

resource "endpointmonitor_web_journey_check" "integration_test" {
  name            = "Integration Web Journey Test"
  description     = "Integration Web Journey Test. Managed by Terraform."
  enabled         = true
  start_url       = "https://koolness.co.uk/test/"
  trigger_count   = 3
  check_frequency = 60

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
    common_step_id = endpointmonitor_web_journey_common_step.test.id
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
      sequence        = 1
      description     = "Enter username"
      always_required = true
      type            = "TEXT_INPUT"

      text_input {
        element_id = "login_username"
        input_text = "my.user@mycompany.com"
      }
    }

    action {
      sequence        = 2
      description     = "Enter password"
      always_required = true
      type            = "PASSWORD_INPUT"

      password_input {
        element_id     = "login_password"
        input_password = "keepthissafe"
      }
    }

    action {
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

  check_host_id  = endpointmonitor_check_host.test.id
  check_group_id = endpointmonitor_check_group.test.id
}

resource "endpointmonitor_check_host_group" "test" {
  name           = "Terraform Test Host Group"
  description    = "Testing Terraform host group resource."
  enabled        = true
  check_host_ids = data.endpointmonitor_check_hosts.test.ids
}

resource "endpointmonitor_android_journey_check" "android" {
  name            = "Integration Android Journey Test"
  description     = "Integration Android Journey Test. Managed by Terraform."
  enabled         = true
  trigger_count   = 3
  check_frequency = 60

  check_host_id  = endpointmonitor_check_host.test.id
  check_group_id = endpointmonitor_check_group.test.id

  apk = filebase64("${path.module}/test.apk")

  common_step {
    sequence       = 1
    common_step_id = endpointmonitor_android_journey_common_step.android_common.id
  }

  custom_step {
    sequence  = 2
    name      = "Handle Your Data Prompt"
    wait_time = 5000

    step_check {
      description = "Check Your Data prompt"
      type        = "CHECK_FOR_TEXT"

      check_for_text {
        text_to_find = "Your Data"
        state        = "PRESENT"
      }
    }

    step_interaction {
      sequence        = 1
      description     = "Close \"Which Sports\" onboarding question."
      always_required = false
      type            = "CLICK"

      click {
        search_text = "Skip"
      }
    }

    step_interaction {
      sequence        = 2
      description     = "Close Bookmaker Selection"
      always_required = false
      type            = "CLICK"

      click {
        search_text = "+"
      }
    }

    step_interaction {
      sequence        = 3
      description     = "Accept Your Data prompt"
      always_required = true
      type            = "CLICK"

      click {
        search_text = "Continue"
      }
    }
  }

  custom_step {
    sequence  = 3
    name      = "Racecards Displayed"
    wait_time = 5000

    step_check {
      description = "Tomorrow's Racecards text shown"
      type        = "CHECK_FOR_TEXT"

      check_for_text {
        text_to_find = "TOMORROW'S RACECARDS"
        state        = "PRESENT"
      }
    }

    step_check {
      description = "Free Bets button shown"
      type        = "CHECK_FOR_ELEMENT"

      check_for_element {
        xpath           = "//*[@content-desc='FREE BETS' and @text='']"
        state           = "PRESENT"
        attribute_name  = "clickable"
        attribute_value = "true"
      }
    }

    step_interaction {
      sequence        = 1
      description     = "Click on first meeting"
      always_required = true
      type            = "CLICK"

      click {
        component_id = "home-screen-todays-racecards-meeting"
      }
    }

    step_interaction {
      sequence        = 2
      description     = "Wait 2 seconds"
      always_required = true
      type            = "WAIT"
      wait_time       = 2000
    }

    step_interaction {
      sequence        = 3
      description     = "Take Screenshot"
      always_required = true
      type            = "SCREENSHOT"
    }
  }
}

resource "endpointmonitor_android_journey_common_step" "android_common" {
  name        = "Integration Android Common Step"
  description = "Integration Test Android Common Step - Managed by Terraform"
  wait_time   = 5000

  step_check {
    description = "Check for Continue button"
    type        = "CHECK_FOR_TEXT"

    check_for_text {
      text_to_find = "Continue"
      state        = "PRESENT"
    }
  }
  
  step_interaction {
    sequence        = 1
    description     = "Click Continue"
    always_required = false
    type            = "CLICK"

    click {
      search_text = "Continue"
    }
  }

  step_interaction {
    sequence        = 2
    description     = "Click Skip"
    always_required = false
    type            = "CLICK"

    click {
      search_text = "SKIP"
    }
  }

  step_interaction {
    sequence        = 3
    description     = "Click Other Skip"
    always_required = false
    type            = "CLICK"

    click {
      search_text = "Skip"
    }
  }
}