---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "endpointmonitor_socket_check Resource - endpointmonitor"
subcategory: ""
description: |-
  Create and manage socket checks which test to ensure a hostname is listening on a pre-defined port.
---

# endpointmonitor_socket_check (Resource)

Create and manage socket checks which test to ensure a hostname is listening on a pre-defined port.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `check_group_id` (Number) The id of the Check Group the check belongs to. This also determines check frequency.
- `check_host_id` (Number) The id of the Check Host to run the check on.
- `hostname` (String) The hostname to check the socket is open on.
- `name` (String) A name to describe in the check, used throughout EndPoint Monitor to describe this check, including in notifications.
- `port` (Number) The TCP port to check is listening.
- `trigger_count` (Number) The sequential number of failures that need to occur for a check to trigger an alert or notification.

### Optional

- `description` (String) A space to provide a longer description of the check if needed. Will default to the name if not set.
- `enabled` (Boolean) Allows the enabling/disabling of the check from executing.
- `maintenance_override` (Boolean) If set true then notifications and alerts will be suppressed for the check.
- `result_retention` (Number) The number of days to store historic results of the check.

### Read-Only

- `id` (String) The ID of this resource.

