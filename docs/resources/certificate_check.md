---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "endpointmonitor_certificate_check Resource - endpointmonitor"
subcategory: ""
description: |-
  Create and manage TLS certificate checks that test a given URL for an expected response.
---

# endpointmonitor_certificate_check (Resource)

Create and manage TLS certificate checks that test a given URL for an expected response.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `alert_days_remaining` (Number) The maximum number of remaining days on a certificate before a failure is triggered.
- `check_group_id` (Number) The id of the Check Group the check belongs to. This also determines check frequency.
- `name` (String) A name to describe in the check, used throughout EndPoint Monitor to describe this check, including in notifications.
- `trigger_count` (Number) The sequential number of failures that need to occur for a check to trigger an alert or notification.
- `url` (String) The URL to check the certificate for.
- `warning_days_remaining` (Number) The maximum number of remaining days on a certificate before an warning is triggered.

### Optional

- `check_date_only` (Boolean) If set to true, then only certificate validity period will be checked and nothing else.
- `check_frequency` (Number) The frequency the check will be run in seconds.
- `check_full_chain` (Boolean) If set to false, only the initially returned certificate from the given URL will be checked, and not the full certificate chain.
- `check_host_group_id` (Number) The id of the Check Host Group to run the check on.
- `check_host_id` (Number) The id of the Check Host to run the check on.
- `description` (String) A space to provide a longer description of the check if needed. Will default to the name if not set.
- `enabled` (Boolean) Allows the enabling/disabling of the check from executing.
- `maintenance_override` (Boolean) If set true then notifications and alerts will be suppressed for the check.
- `result_retention` (Number) The number of days to store historic results of the check.

### Read-Only

- `id` (String) The ID of this resource.

