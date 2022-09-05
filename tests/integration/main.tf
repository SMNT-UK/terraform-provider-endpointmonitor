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
  url = "https://smnt-read-epm.net.smnt.co.uk/api"
}

resource "endpointmonitor_app_group" "integration_tests" {
  name        = "Integration Tests"
  description = "Integration test checks. Managed by Terraform."
}

resource "endpointmonitor_check_group" "integration_tests" {
  name            = "Integration Tests"
  description     = "Integration test checks. Managed by Terraform."
  check_frequency = 120
  app_group_id    = endpointmonitor_app_group.integration_tests.id
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