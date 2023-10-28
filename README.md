[English](./README_en.md)

tfoonはterraformの依存関係を収集します

## CLIのインストール

```bash
go install github.com/wim-web/tfoon@latest
```

### 実行サンプル

```bash
# 実行するterraformが依存しているmodule一覧
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
# moduleが依存されているterraform(エントリーポイント)一覧
# CIでgitの差分から実行すべきterraformを調査するのに使えます
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
