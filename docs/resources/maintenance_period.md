---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "endpointmonitor_maintenance_period Resource - endpointmonitor"
subcategory: ""
description: |-
  Create and manage scheduled maintenance periods to prevent checks from alerting during certain periods of the day or week.
---

# endpointmonitor_maintenance_period (Resource)

Create and manage scheduled maintenance periods to prevent checks from alerting during certain periods of the day or week.

## Example Usage

```terraform
# Example Maintenance Period that suppresses alerts between 01h and 03h every night 
# for all checks with 'Website' in their name or description.

data "endpointmonitor_checks" "websites" {
  search = "Website"
}

resource "endpointmonitor_maintenance_period" "example" {
  description = "Suppress alerts during backup period."
  enabled     = true
  day_of_week = "ALL"
  start_time  = "01:00"
  end_time    = "03:00"
  check_ids   = data.endpointmonitor_checks.websites.ids
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `day_of_week` (String) The day of week the maintenance period applies to. Set as ALL for every day of the week. Must otherwise be SUNDAY, MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY or SATURDAY.
- `description` (String) Space for a description of the maintenance periods purpose.
- `enabled` (Boolean) Enable or disable the maintenance period from applying to attached checks.
- `end_time` (String) The end of time the maintenance period in format 24HH:MM.
- `start_time` (String) The start time of the maintenance period in format 24HH:MM.

### Optional

- `check_group_ids` (List of Number) A list of ids of Check Groups that are directly linked to the maintenance period.
- `check_ids` (List of Number) A list of ids of Checks that are directly linked to the maintenance period.
- `dashboard_group_ids` (List of Number) A list of ids of Dashboard Groups that are linked to this maintenance period.

### Read-Only

- `id` (Number) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# Maintenance Periods can be imported using their numeric id, which can be see in the address bar when editing a Maintenance Period in the web interface.
terraform import endpointmonitor_maintenance_period.example 123
```
