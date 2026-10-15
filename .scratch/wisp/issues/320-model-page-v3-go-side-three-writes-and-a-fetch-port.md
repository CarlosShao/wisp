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
- [2026-10-15 11:4x +08] agent=320-a1 did=只读设计核查（锚 `36cf5d06`，零产码；件 `.scratch/wisp/probes/320/a1/census.md` 107 行）。**字段集定稿**：三枚全落 313 名册（`internal/panel/config_write_313.go` 的 `fields313`＋twin 三 map）——`provider_billing`（值 `plan|pay-per-token`，provider+值）、`provider_row_create_custom`〔名待裁，值 `<protocol>,<base_url>`＋可选 model〕、`model_row_create`（provider+model，无值）；248 七枚一枚不许动（两处 `want` 钉）。同批必改测试逐文件点名：`config_313_test.go`（`fields313Cases()` 追加三 case；命名禁令测试读 live 自动覆盖）、`cmd/wisp/panel_config_313_test.go`（并集钉读 live 自动）、`internal/config/settings_313_test.go`（三 setter）。**新方法名**：推荐 `config.models.list`（⛔ 不带 `panel.` 前缀——`git_test.go:385-391` 钉四枚 `panel.*` 字面量）；注册四点＋两点＝bridge.go 常量/guard、composer_dispatch case、handler 分支、`inbound_roster_253_test.go` 名册＋两枚计数（6→7、2→3），**且 `frontend/src/lib/panel.ts` 必须同批出现该字面量**（`bridge_test.go:184` 逐枚读 Method* 常量）⇒ 前端那枚归机主外部 agent，Go 半枚做不到 self-green（与 248 同形）；载荷＝散文＋无点标记行（候选 `wisp-models-list:`），config.get 载荷键钉不动。冻结件 `l2_grant_boundary_test.go` 逐字核过：路由词表 11 枚归一化后 0 命中、字段词表只扫 JSON 解码落点的键、locked 前缀 `models.` 带点不撞 `model_row_create`。**执行面**＝`llm.DiscoverModels`＋`secret.Store.Resolve`（两族 ref）＋`context.WithTimeout`，零 goroutine（ban#1 连具名 `go` 都扫）、零落盘、不回显；socket 建议 `ConfigWriteHandler` 另挂可选接口（不加 `ConfigStore` 方法，保 248 测试件零改）。**AC#5 裁：不落**（`defaults.go:12-14` 逐字只供 protocol/base_url；落地要动 preset merge 面比票大）。**面序**：319-r2 → 320-r1（config/panel 半先、cmd/wisp 半后）→ 135-r1；`cmd/wisp`／`internal/panel` 共享编译面每波一腿。现跑门读数：d22scan clean、长度门 GREEN（两把只读脚本）。未定五项（AC#2 甲乙形、拉取超时值/同步面、上游错误体首行进回执、`models.list` 与 SPEC-08 §5.2:172 相撞、AC#5 若改"落"的合并规则）逐条给候选＋判据缺口。next=编排者裁五枚未定＋派 `320-r1`（先 config/panel 半）；⛔ 本腿零产码、未动台账/框/票名、未碰 313/248/`frontend/**` 与任何产码。
