# 320 — 模型页 v3 的 Go 侧：**三个写入能力**（计费模式／自定义供应商行／模型条目）＋**"拉取模型列表"端口**（机主 2026-10-14 22:5x 口头批准＝接口面人工批准）

- **Status**：**待派**（设计腿 `320-a1` 起手；写腿面＝`internal/config`＋`internal/panel`＋`cmd/wisp`，与在飞写腿串行）。
- **来源**：机主 2026-10-14 22:5x 原话逐字「**模型列表要允许拉，三个写入能力可以做**」＋方向「别为了录入 key 就搞个简陋的页面……按照票上来」。**批准落账＝`A968`**（本条即"接口面变更获人工批准"的记录；改 C17/config.set 白名单由此获准）。
- **背景证据**：`docs/evidence/s1/model-page-v3-brief.md`（§2 Go 侧普查：`billing` 枚举已存在、非预设名合法〔`validate.go:209-212`〕、`provider_row_create` 选择器只吃预设名、无 billing 写入字段）＋台账 `A960`/`A968`。
- **AC（1:1；裁决表由非实现者出）**：
  - [ ] **AC#1** `provider_billing` 写入（值 `plan|pay-per-token`，枚举校验；落 `llm.providers.<name>.billing`；选择列与既有 248/313 表机制一致）。
  - [ ] **AC#2** **自定义供应商行**：非预设名＋显式 `protocol`＋`base_url`（数据层已合法）；"复制 stepfun 为 stepfun-api"形由 UI 组合；值格式/字段名由设计腿定稿。
  - [ ] **AC#3** **模型条目写入**（provider＋model id；窗口/价格复用既有三字段或并入，设计腿定稿）。
  - [ ] **AC#4** **"拉取模型列表"端口**：用已存密钥 GET `<base>/models`（超时/错误分类/不回显密钥/不落盘明文）；返回 id 列表；**新方法名由设计腿定稿并过白名单检查**。
  - [ ] **AC#5**（可选）预设默认模型目录：仅在拿到真名单后落；或明确"不落、一律拉"（设计腿裁）。
  - [ ] **AC#6** 门禁＋契约：313/248 两表与 `config_313_test.go` 同批对齐；`go test ./internal/panel/ ./internal/config/`＋相关 `cmd/wisp` 用例；d22scan；长度门；⛔ **冻结件 `internal/panel/l2_grant_boundary_test.go` 零改动**（新字段/方法名不许带审批词汇——设计腿先验）。
  - [ ] **AC#7** 证据：每格两向三态＋变异自证；证据件 `docs/evidence/s1/320-*.md`。
- **禁区**：`docs/PLAN.md`、`docs/specs/**`、`internal/observe/thresholds.go`、golden、`tools/d22scan/allowlist.txt`、`.github/**`、`frontend/**`、`design/**`；冻结三件不动。
- **撤销口令**：**「撤 320」**（⇒ 新字段/端口一并退回，页面回到只读展示）。
- **关联**：313（v1 底层）· `docs/evidence/s1/model-page-v3-brief.md`（前端件）· 319（面板宿主修复，在飞）。

## Progress log (append-only, newest last)

- [2026-10-14 22:5x +08] 编排者立票：机主授权逐字落账（`A968`）；设计腿 `320-a1` 起手（只读：字段集/值格式/新方法名/契约测试同批计划/与在飞腿的面序）。
