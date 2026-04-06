# base +record-batch-set

> **前置条件：** 先阅读 [`../lark-shared/SKILL.md`](../../lark-shared/SKILL.md) 了解认证、全局参数和安全规则。

批量更新记录。

## 推荐命令

```bash
lark-cli base +record-batch-set \
  --base-token app_xxx \
  --table-id tbl_xxx \
  --json '{"patch":{"状态":"Done"},"filter":{"logic":"and","conditions":[["状态","==","Open"]]},"offset":0,"limit":200}'
```

```bash
lark-cli base +record-batch-set \
  --base-token app_xxx \
  --table-id tbl_xxx \
  --json '{"records":["rec_xxx"],"fields":["fld_xxx","fld_yyy"],"cells":[["xx","yy"]]}'
```

## 参数

| 参数 | 必填 | 说明 |
|------|------|------|
| `--base-token <token>` | 是 | Base Token |
| `--table-id <id_or_name>` | 是 | 表 ID 或表名 |
| `--json <body>` | 是 | 批量更新请求体，必须是 JSON 对象；CLI 不改写内容，直接发送给 API |

## API 入参详情

**HTTP 方法和路径：**

```
PATCH /open-apis/base/v3/bases/:base_token/tables/:table_id/records/batch
```

## `--json` Raw JSON Schema

```json
{"type":"object","properties":{"filter":{"type":"object","properties":{"logic":{"type":"string","enum":["and","or"],"default":"and","description":"Filter Condition Logic"},"conditions":{"type":"array","items":{"type":"array","minItems":3,"maxItems":3,"items":[{"type":"string","minLength":1,"maxLength":100,"description":"Field id or name"},{"type":"string","enum":["==","!=",">",">=","<","<=","intersects","disjoint","empty","non_empty"],"description":"Condition operator"},{"anyOf":[{"not":{}},{"anyOf":[{"anyOf":[{"type":"string","description":"text & formula & location field support string as filter value"},{"type":"number","description":"number & auto_number(the underfly incremental_number) field support number as filter value"},{"type":"array","items":{"type":"string","description":"option name"},"description":"select field support one option: [\"option1\"] or multiple options: `[\"option1\", \"option2\"]` as filter value."},{"type":"array","items":{"type":"object","properties":{"id":{"type":"string","description":"record id"}},"required":["id"],"additionalProperties":false},"description":"link field support record id list as filter value"},{"type":"string","description":"\ndatetime & create_at & updated_at field support relative and absolute filter value.\nabsolute:\n- \"ExactDate(yyyy-MM-dd)\"\nrelative:\n- Today\n- Tomorrow\n- Yesterday\n"},{"type":"array","items":{"type":"object","properties":{"id":{"type":"string","description":"user id"}},"required":["id"],"additionalProperties":false},"description":"user field support user id list as filter value"},{"type":"boolean","description":"checkbox field support boolean as filter value"}]},{"type":"null"}]}]}],"description":"one condition expression. shape: [field_id, filter_operator, value]. when operator is \"empty\" or \"non_empty\", the value is not required."},"default":[]}},"additionalProperties":false},"offset":{"type":"integer","minimum":0,"default":0},"limit":{"type":"integer","minimum":1,"maximum":200},"patch":{"type":"object","additionalProperties":{"anyOf":[{"anyOf":[{"type":"string","description":"text field cell, example: \"one string and [one url](https://foo.bar)\""},{"type":"number","description":"number field cell, can be any float64 value"},{"type":"array","items":{"type":"string","description":"option name"},"description":"select field cell, example: [\"option_1\", \"option_2\"]"},{"type":"string","description":"datetime field cell. accepts common datetime strings and timestamp-like values. Prefer \"YYYY-MM-DD HH:mm:ss\" in requests because it is the most stable format and matches the API output. Example: \"2026-01-01 19:30:00\""},{"type":"array","items":{"type":"object","properties":{"id":{"type":"string","description":"record id"}},"required":["id"],"additionalProperties":false},"description":"link field cell, example: [{\"id\": \"rec_123\"}]"},{"type":"array","items":{"type":"object","properties":{"id":{"type":"string","description":"user id"}},"required":["id"],"additionalProperties":false},"description":"user field cell, example: [{\"id\": \"ou_123\"}]"},{"type":"object","properties":{"lng":{"type":"number","description":"Longitude"},"lat":{"type":"number","description":"Latitude"}},"required":["lng","lat"],"additionalProperties":false,"description":"location field cell, example: {\"lng\": 113.94765, \"lat\": 22.528533}"},{"type":"boolean","description":"checkbox field cell"},{"type":"array","items":{"type":"object","properties":{"file_token":{"type":"string","minLength":0,"maxLength":50},"name":{"type":"string","minLength":1,"maxLength":255},"mime_type":{"type":"string","maxLength":255,"description":"deprecated field"},"size":{"type":"integer","minimum":0,"description":"deprecated field"},"image_width":{"type":"integer","minimum":0,"description":"deprecated field"},"image_height":{"type":"integer","minimum":0,"description":"deprecated field"},"deprecated_set_attachment":{"type":"boolean","description":"deprecated field"}},"required":["file_token","name"],"additionalProperties":false},"description":"attachment field cell. temporary compatibility for attachment writes."},{"type":"null"}]},{"type":"null"}]},"propertyNames":{"minLength":1,"maxLength":100}}},"required":["filter","limit","patch"],"additionalProperties":false,"$schema":"http://json-schema.org/draft-07/schema#"}
```

## 返回重点

- 返回对象键：
  - `has_more`
  - `record_id_list`
  - `update`
  - `ignored_fields`（可选）

## 坑点

- ⚠️ `--json` 必须是对象。
- ⚠️ `limit` 超过 `200` 会被 OpenAPI 拒绝。
- ⚠️ 命令不会自动做字段/行映射转换，传什么就发什么。

## 参考

- [lark-base-record.md](lark-base-record.md) — record 索引页
- [lark-base-view-set-filter.md](lark-base-view-set-filter.md) — 视图筛选配置参考
