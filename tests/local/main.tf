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
  #url = "http://10.10.0.11:8080/api"
  url = "http://192.168.44.99:8080/api"
  key = "7412f79a-70a4-450c-8c10-0b5e958a726d" # This is safe for leaving in, it's just a key on a local dev instance.
}

resource "endpointmonitor_proxy_host" "test" {
  name        = "Test Terraform ProxyHost"
  description = "Testing Terraform Descrition"
  hostname    = "test-terraform-host"
  port        = 3318
}

resource "endpointmonitor_check_host" "test" {
  name             = "Terraform CheckHost 2"
  description      = "Terraform CheckHost Description"
  hostname         = "test.terraform.host2"
  type             = "AGENT"
  enabled          = true
  max_checks       = 1
  send_check_files = true
}

resource "endpointmonitor_app_group" "test" {
  name        = "Terraform App Group 2"
  description = "Terraform Descrition of an App Group"
}

resource "endpointmonitor_check_group" "test" {
  name            = "Terraform Check Group 2"
  description     = "Terraform description of a check group"
  check_frequency = 60
  app_group_id    = endpointmonitor_app_group.test.id
}

resource "endpointmonitor_url_check" "test" {
  name                   = "Terraform URL Check"
  description            = "Terraform URL check descrtiption"
  url                    = "https://www.bbc.co.uk/weather"
  trigger_count          = 2
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
  check_host_id  = data.endpointmonitor_check_host.test.id
  check_group_id = endpointmonitor_check_group.test.id
}

resource "endpointmonitor_dns_check" "test" {
  name               = "Terraform DNS Check"
  description        = "Terraform DNS check descrtiption"
  hostname           = "smnt-read-sql01.net.smnt.co.uk"
  expected_addresses = ["10.20.0.31", "10.20.0.32"]
  trigger_count      = 5
  check_host_id      = endpointmonitor_check_host.test.id
  check_group_id     = endpointmonitor_check_group.test.id
}

resource "endpointmonitor_web_journey_check" "test" {
  name                 = "Terraform WebJourney Check"
  description          = "Terraform WebJourney check descrtiption"
  enabled              = false
  maintenance_override = true
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