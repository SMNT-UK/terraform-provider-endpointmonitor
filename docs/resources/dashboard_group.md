---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "endpointmonitor_dashboard_group Resource - endpointmonitor"
subcategory: ""
description: |-
  Create and manage Dashboard Groups, the top-level organisational groups of checks running in EndPoint Monitor.
---

# endpointmonitor_dashboard_group (Resource)

Create and manage Dashboard Groups, the top-level organisational groups of checks running in EndPoint Monitor.

## Example Usage

```terraform
resource "endpointmonitor_dashboard_group" "example" {
  name        = "Public Websites"
  description = "Contains monitors for all public facing websites we host."
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `description` (String) Space to provide a longer description of this Dashboard Group.
- `name` (String) The name of the Dashboard Group. This will be used in alerts and notifications.

### Read-Only

- `id` (Number) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# Dashboard Groups can be imported using their numeric id, which can be see in the address bar when editing a Dashboard Group in the web interface.
terraform import endpointmonitor_dashboard_group.example 123
```
