---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "endpointmonitor_check_host Resource - endpointmonitor"
subcategory: ""
description: |-
  Create and manage the hosts that checks are to be run on.
---

# endpointmonitor_check_host (Resource)

Create and manage the hosts that checks are to be run on.

## Example Usage

```terraform
resource "endpointmonitor_check_host" "example" {
  name             = "Internal Applications Agent"
  description      = "Agent for monitoring internal facing applications."
  hostname         = "machinename.fqdn.com"
  type             = "AGENT"
  enabled          = true
  max_checks       = 1
  send_check_files = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `description` (String) A place to provide more detail about the host if required.
- `enabled` (Boolean) If disabled checks set to run against this host will be paused.
- `hostname` (String) The hostname of the host. Must match what the host believes it's own hostname is.
- `name` (String) A friendly name to describe the host. This is what the host will be refered to as on all screens and alerts.
- `type` (String) Must be either CONTROLLER or AGENT. CONTROLLER is used for hosts that expose the Web GUI and required database access. AGENT is used for hosts that purely just run checks.

### Optional

- `max_checks` (Number) The maximum number of concurrent Web Journey checks the host can run. Default is 1.
- `send_check_files` (Boolean) For agents only. Indicates if it is to send check files such as screenshots back to the controller through the controller API. Should be enabled if there isn't a common file share between agent and controllers.

### Read-Only

- `id` (String) The ID of this resource.

