tfoon collects Terraform dependencies.

## CLI Installation

```bash
go install github.com/wim-web/tfoon@latest
```

### Execution Example

```bash
# List of modules that the Terraform being executed depends on
tfoon ./testdata/terraform/caller1/,./testdata/terraform/caller3 | jq

[
  {
    "path": "testdata/terraform/caller1",
    "children": [
      {
        "path": "testdata/terraform/modules/noop",
        "children": []
      },
      {
        "path": "testdata/terraform/modules/nest_noop",
        "children": [
          {
            "path": "testdata/terraform/modules/nest_noop/modules/inner_noop",
            "children": []
          }
        ]
      }
    ]
  },
  {
    "path": "testdata/terraform/caller3",
    "children": [
      {
        "path": "testdata/terraform/modules/noop",
        "children": []
      }
    ]
  }
]
```

---

```bash
# List of Terraform (entry points) that modules depend on
# Useful for investigating which Terraform should be executed based on Git diffs in a CI/CD pipeline
tfoon -m2e ./testdata/terraform/caller1/,./testdata/terraform/caller3 | jq

{
  "testdata/terraform/modules/nest_noop": [
    "testdata/terraform/caller1"
  ],
  "testdata/terraform/modules/nest_noop/modules/inner_noop": [
    "testdata/terraform/caller1"
  ],
  "testdata/terraform/modules/noop": [
    "testdata/terraform/caller1",
    "testdata/terraform/caller3"
  ]
}
```
