# dalpha-terraform-provider

### ~/.terraformrc 설정

```
provider_installation {

  dev_overrides {
      "registry.terraform.io/dalphakr/dalpha" = "{repository 경로}"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

```
cd dalpha-infra/terraform/openai
vim .envrc
```
