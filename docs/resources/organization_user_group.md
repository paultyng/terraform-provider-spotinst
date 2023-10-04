---
layout: "spotinst"
page_title: "Spotinst: organization_user_group"
subcategory: "Organization"
description: |-
  Provides a Spotinst user-group of your Spot organization.
---

# spotinst\_organization\_user\_group

Provides a Spotinst user-group of your Spot organization.

## Example Usage

```hcl
resource "spotinst_organization_user_group" "terraform_user_group" {
  name = "test_user_group"
  description = "user group by terraform"
  user_ids = ["u-372gf6ae"]
  policies {
    account_ids = ["act-a1b2c3d4"]
    policy_id = "pol-vv7d8c06"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) User group name.
* `description` - (Optional) User group description.
* `user_ids` - (Optional) The users to register under the created group
   (should be existing users only).
* `policies` - (Optional) The policies to register under the given group
   (should be existing policies only).
  * `account_ids` - (Required) A list of accounts to register with the assigned under the
     given group (should be existing accounts only).
  * `policy_id` - (Required) A policy to register under the given group
     (should be existing policy only).

## Attributes Reference

The following attributes are exported:

* `id` - The Spotinst User Group ID.
