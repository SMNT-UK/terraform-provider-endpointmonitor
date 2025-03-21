---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "endpointmonitor_web_journey_common_step Resource - endpointmonitor"
subcategory: ""
description: |-
  Create and manage web journey common steps which are used to provide common checks and actions to take for web journey checks.
---

# endpointmonitor_web_journey_common_step (Resource)

Create and manage web journey common steps which are used to provide common checks and actions to take for web journey checks.

## Example Usage

```terraform
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
    sequence        = 1
    description     = "Close any marketing pop-up."
    always_required = false
    type            = "CLICK"

    click {
      xpath = "//*[@class='aria-close']"
    }
  }

  action {
    sequence        = 2
    description     = "Accept Cookies"
    always_required = true
    type            = "CLICK"

    click {
      search_text  = "Accept All"
      element_type = "button"
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `description` (String) Space to provide a longer description of what this common step can be used for.
- `name` (String) A name to describe what the step is doing. This will be included in any alerts and notifications.

### Optional

- `action` (Block List) The set of actions to perform at the end of the step such as clicking on elements or enterting text. (see [below for nested schema](#nestedblock--action))
- `console_message_suppression` (Block List) Suppress one or more cosole log messages from creating a warning or failure for a Web Journey Step. (see [below for nested schema](#nestedblock--console_message_suppression))
- `network_suppression` (Block List) Suppress one or more network calls from causing any warnings or failures. (see [below for nested schema](#nestedblock--network_suppression))
- `page_check` (Block List) The set of checks to run against the currently loaded content. (see [below for nested schema](#nestedblock--page_check))
- `page_load_time_alert` (Number) The maximum number of milliseconds that any discovered network call can take before an alert is created for it, and the check is set to a failed status.
- `page_load_time_warning` (Number) The maximum number of milliseconds that any discovered network call can take before a warning is created for it and the check is set to a warning status.
- `wait_time` (Number) The number of milliseconds to wait for any page load / actions on the page to complete before any checks on this step are started.

### Read-Only

- `id` (Number) The ID of this resource.

<a id="nestedblock--action"></a>
### Nested Schema for `action`

Required:

- `description` (String) Space for a description of what this action does.
- `sequence` (Number) This defines the order that actions will be taken, from number lowest first to highest number last.
- `type` (String) The type of action to perform. Options are: CLICK, DOUBLE_CLICK, RIGHT_CLICK, TEXT_INPUT, PASSWORD_INPUT, CHANGE_WINDOW_BY_ORDER, CHANGE_WINDOW_BY_TITLE, NAVIGATE_URL, WAIT, REFRESH_PAGE, CLOSE_WINDOW, CHANGE_IFRAME_BY_ORDER, CHANGE_IFRAME_BY_XPATH, SCROLL_TO_ELEMENT, TAKE_SCREENSHOT, SAVE_DOM or SELECT_OPTION.

Optional:

- `always_required` (Boolean) If true the the action given must be able to be completed against the current page, and if it can't the check will be marked as failed. If false, and the action can't complete, for example because the element is missing, the step will continue onto the next action regardless.
- `click` (Block, Optional) The additional details needed for a CLICK, DOUBLE_CLICK or RIGHT_CLICK action type. (see [below for nested schema](#nestedblock--action--click))
- `iframe_id` (Number) The order number of the iframe to set focus to for the CHANGE_IFRAME_BY_ORDER action type. Set to 0 if you need to move focus back to the main page.
- `iframe_xpath` (String) The xpath of the iframe to set focus to for the CHANGE_IFRAME_BY_XPATH action type.
- `navigate_url` (String) The URL to navigate to for the NAVIGATE_URL action type.
- `password_input` (Block, Optional) The additional details needed for a PASSWORD_INPUT action type. (see [below for nested schema](#nestedblock--action--password_input))
- `scroll_to_element` (Block, Optional) The additional details needed for the SCROLL_TO_ELEMENT action type. (see [below for nested schema](#nestedblock--action--scroll_to_element))
- `select_option` (Block, Optional) Additional details needed for SELECT_OPTION action type, used to choose a value from a select element. (see [below for nested schema](#nestedblock--action--select_option))
- `text_input` (Block, Optional) The additional details needed for a TEXT_INPUT action type. (see [below for nested schema](#nestedblock--action--text_input))
- `wait_time` (Number) The number of milliseconds to wait for the WAIT action type.
- `window_id` (Number) The opening order number of the window to change focus to for CHANGE_WINDOW_BY_ORDER action types.
- `window_title` (String) The title of the window to change focus to for CHANGE_WINDOW_BY_TITLE action types.

Read-Only:

- `id` (Number)

<a id="nestedblock--action--click"></a>
### Nested Schema for `action.click`

Optional:

- `element_type` (String) Only to be used alongside search_text. The element type/name to help target the given search_text.
- `search_text` (String) The text on the page to click on. If this has multiple matches then the first will be used. Can not be used with xpath.
- `xpath` (String) The xpath of the element to click on. If multiple matches, the first will be used. Can not be used with search_text.


<a id="nestedblock--action--password_input"></a>
### Nested Schema for `action.password_input`

Optional:

- `element_id` (String) The id of the element to input the password into. Not to be used with xapth or element_name.
- `element_name` (String) The name of the element to input the password into. Not to be used with xapth or element_id.
- `input_password` (String) The password to input. This will not be stored in your Terraform state and ideally should be passed in to your Terraform as a environment variable rather than statically stored in your Terraform code.
- `xpath` (String) The xpath of the element to input the password into. If multiple matches, the first will be used. Not to be used with element_id or element_name.


<a id="nestedblock--action--scroll_to_element"></a>
### Nested Schema for `action.scroll_to_element`

Optional:

- `element_type` (String) Only to be used alongside search_text. The element type/name to help target the given search_text.
- `search_text` (String) The text on the page to scroll to. If this has multiple matches then the first will be used. Can not be used with xpath.
- `xpath` (String) The xpath of the element to scroll to. If multiple matches, the first will be used. Can not be used with search_text.


<a id="nestedblock--action--select_option"></a>
### Nested Schema for `action.select_option`

Optional:

- `element_id` (String) The id of the select element to select a value from. Not to be used when xpath is set.
- `option_index` (Number) Choose the option to select by the order it is shown in the list, starting from 0.
- `option_name` (String) Choose the option to select by its name shown in the list.
- `option_value` (String) Choose the option to select by its form value.
- `xpath` (String) The xpath of the select element to select a value from. Not to be used when element_id is set.


<a id="nestedblock--action--text_input"></a>
### Nested Schema for `action.text_input`

Optional:

- `element_id` (String) The id of the element to input text into. Not to be used with xapth or element_name.
- `element_name` (String) The name of the element to input text into. Not to be used with xapth or element_id.
- `input_text` (String) The text to input.
- `xpath` (String) The xpath of the element to input text into. If multiple matches, the first will be used. Not to be used with element_id or element_name.



<a id="nestedblock--console_message_suppression"></a>
### Nested Schema for `console_message_suppression`

Required:

- `comparison` (String) Must be EQUALS, STARTS_WITH, ENDS_WITH or CONTAINS. The way to compare the given full or part message to the console logs made after the previous step.
- `description` (String) Space for a description of what this is supressing.
- `log_level` (String) The log level to suppress. Must be ANY, WARNING or ERROR.
- `message` (String) The full log message or part of the log message to suppress.

Read-Only:

- `id` (Number)


<a id="nestedblock--network_suppression"></a>
### Nested Schema for `network_suppression`

Required:

- `comparison` (String) Must be EQUALS, STARTS_WITH, ENDS_WITH or CONTAINS. The way to compare the given url to the network calls made after the last step.
- `description` (String) Space for a description of what this is supressing.
- `url` (String) The full or part URL to suppress.

Optional:

- `any_client_error` (Boolean) Suppress any 400-499 response code for the given url.
- `any_server_error` (Boolean) Suppress any 500-599 response code for the given url.
- `response_code` (Number) The response code for the given url that is to be suppressed for warnings or alerts.

Read-Only:

- `id` (Number)


<a id="nestedblock--page_check"></a>
### Nested Schema for `page_check`

Required:

- `description` (String) A description of what this check is doing. This will be used in alerts and notifications.
- `type` (String) The type of check to execute. Options are: CHECK_FOR_TEXT - Check for any string on or not on the current page. CHECK_FOR_ELEMENT - Check for an element and it's properties on the current page. CHECK_CURRENT_URL - Check the current url. CHECK_URL_RESPONSE - Check for specific network calls made after the last step. CHECK_CONSOLE_LOG - Check for console logs made after the last step.

Optional:

- `check_console_log` (Block, Optional) Check for a log entry made after the past step. (see [below for nested schema](#nestedblock--page_check--check_console_log))
- `check_current_url` (Block, Optional) Check the URL of the current page. (see [below for nested schema](#nestedblock--page_check--check_current_url))
- `check_element_on_page` (Block, Optional) Check for a specific element and it's attributes on the current page. (see [below for nested schema](#nestedblock--page_check--check_element_on_page))
- `check_for_text` (Block, Optional) Check a specific stirng is present or absent on the current page. (see [below for nested schema](#nestedblock--page_check--check_for_text))
- `check_url_response` (Block, Optional) Check a network request made after the previous step. (see [below for nested schema](#nestedblock--page_check--check_url_response))
- `warning_only` (Boolean) If true then if this check fails, then it will only produce a warning, not a full check failure. Default is false.

Read-Only:

- `id` (Number)

<a id="nestedblock--page_check--check_console_log"></a>
### Nested Schema for `page_check.check_console_log`

Optional:

- `comparison` (String) Must be one of EQUALS, DOES_NOT_EQUAL, STARTS_WITH, ENDS_WITH, CONTAINS or DOES_NOT_CONTAIN. The way to compare the given message against the console logs.
- `log_level` (String) Must be one of: ANY, NORMAL, WARNING or ERROR. The level of the log to check for.
- `message` (String) The full or partial log message to check for.

Read-Only:

- `id` (Number)


<a id="nestedblock--page_check--check_current_url"></a>
### Nested Schema for `page_check.check_current_url`

Optional:

- `comparison` (String) Must be one of EQUALS, DOES_NOT_EQUAL, STARTS_WITH, ENDS_WITH, CONTAINS or DOES_NOT_CONTAIN. The way to compare the given url against the current URL of the page.
- `url` (String) The URL to compare against the current URL of the page.

Read-Only:

- `id` (Number)


<a id="nestedblock--page_check--check_element_on_page"></a>
### Nested Schema for `page_check.check_element_on_page`

Optional:

- `attribute_name` (String) Filter element matches out by those only containing a given attribute name.
- `attribute_value` (String) Further filter element matches out by having a given attribute value too.
- `elemenet_id` (String) The id of the element to check.
- `elemenet_name` (String) The name of the element to check.
- `element_content` (String) Filter element matches out by their content.
- `state` (String) Must be either PRESENT or ABSENT. PRESENT means the element must be found on the page for the check to succeed. ABSENT means the element must not be on the page for the check to succeed.

Read-Only:

- `id` (Number)


<a id="nestedblock--page_check--check_for_text"></a>
### Nested Schema for `page_check.check_for_text`

Optional:

- `element_type` (String) Limit the search to specific elements.
- `state` (String) Must be either PRESENT or ABSENT. PRESENT means the text_to_find must be found on the page for the check to succeed. ABSENT mesns the text_to_find must not be on the page for the check to succeed.
- `text_to_find` (String) The string to search for for on the page.

Read-Only:

- `id` (Number)


<a id="nestedblock--page_check--check_url_response"></a>
### Nested Schema for `page_check.check_url_response`

Optional:

- `alert_response_time` (Number) The response time in milliseconds that will trigger the check to fail.
- `any_client_error_response` (Boolean) Accept any response code from 400-499.
- `any_info_response` (Boolean) Accept any response code from 100-199.
- `any_redirect_response` (Boolean) Accept any response code from 300-399.
- `any_server_error_response` (Boolean) Accept any response code from 500-599.
- `any_success_response` (Boolean) Accept any response code from 200-299.
- `comparison` (String) Must be one of EQUALS, DOES_NOT_EQUAL, STARTS_WITH, ENDS_WITH, CONTAINS or DOES_NOT_CONTAIN. The way to compare the given url against the current URL of the page.
- `response_code` (Number) The response code required for the check to be successful.
- `url` (String) The URL to search for.
- `warning_response_time` (Number) The response time in milliseconds that will trigger a warning.

Read-Only:

- `id` (Number)

## Import

Import is supported using the following syntax:

```shell
# Web Journey Common Steps can be imported using their numeric id, which can be see in the address bar when editing a Common Step in the web interface.
terraform import endpointmonitor_web_journey_common_step.example 123
```
