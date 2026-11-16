# 248 — 设置写入路由收口轮（2026-11-16，写腿 agent-78；只核不翻框，框归非实现者）

来路＝派单 `.scratch/wisp/probes/orch/20261116/txt/248-settings-route-dispatch.txt`（HEAD 载入锚 `189d3e35d`）。
本轮**零产码改动**：票 248 的主链（入向白名单扩名、config.get/config.set、受保护写、DPAPI 凭据、哨兵常驻用例）由此前各腿（`248-r1`/`248-r2` 及后续 35/255/256/265 波次）已全部在树上；本轮做的是派单点名的四件核量＋AC#9 的能力读数＋证据落档。

## 1. 设置写入路由（主体）——已落地，file:line 现读

| 件 | 落点 |
|---|---|
| 白名单两枚方法名 | `internal/panel/bridge.go:66-67`（`MethodConfigGet`/`MethodConfigSet`，名字出自 SPEC-08 §5.2，不带 `panel.` 前缀、不含授权词表）；`knownComposerMethod` 名册 `bridge.go:193`（单行 const case，冻结锚改形后仍钉着） |
| 派发表接线 | `cmd/wisp/panel_inbound.go:325-329`（`ConfigWriteHandler` 建链）＋`:397`（`ComposerDispatch.Config` 赋值）；凭据存储同链装配 `:321`（`secret.NewStore(dir)`，DPAPI 唯一持有者） |
| 共用封套不含凭据键 | `bridge.go:157-160` 只有 `configField/configProvider/configModel/configValue` 四枚 selector；凭据值走**写-only 专用封套** `credentialWriteRequest`（`cmd/wisp/panel_config_store.go:60-67`，json 键 `credentialValue` 不在共用封套上） |
| 写 config.toml 档位 | `cmd/wisp/panel_config_store.go:226-441` `ApplySetting`：逐字段白名单 switch，全部经 `config.Manager` 的受保护 setter（`internal/config/writeguard.go:111` `mergeWrite` 家族：先重读、单键替换、校验）落盘；档位＝`res.Tier = panel.EffectiveRestart`（`:412`，凭据写入同款 `:487`）——回执答「需重启」，`[ball]` 两枚＋`ball_style` 是实测 now（`:419-433`），`[panel]` 两枚带 Note「关窗再开即跟上」仍答 restart（`:438-440`） |
| 凭据落点 | `StoreCredential`（`panel_config_store.go:450-491`）：值→`secret.NewRef()`＋`secrets.Store(ref, value)`（DPAPI，`internal/secret/dpapi_windows.go`）；config.toml 只落 `api_key_ref = 'dpapi:<id>'`；⛔ 明文不进 config.toml、不进回执、不进审计行 |
| 写前校验（AC#7） | 非法值在进 setter 前被字段级解析拒（`:246-252` 引用形状、`:254-266` 数值形状）；深层校验由 config 侧 validate 单源，`TestAC7InvalidSettingsValueIsRefusedBeforeTheFile` 钉磁盘逐字节不变 |

## 2. AC#9 —— 能力读数（本轮实做）

- 构建：`npm --prefix frontend run build`（node 24.21.0 / npm 11.19.0，`D:\work\server\node24`）⇒ **rc=0**，848ms，产物 `dist/index.html`＋`dist/assets/index-CZoXrb9N.css`（65,264B）＋`dist/assets/index-QUtCIbHm.js`（640,318B）。
- Go 侧读数（`go test ./internal/panel -run TestBuiltinBundleIsCompleteOrAbsentNeverHalf -count=1 -v`）⇒ **PASS**，判语逐字：
  `verdict: bundle complete - entry 1044 bytes, 2 asset refs, 4 embedded files`
  ＝ `frontend.Dist()` 经 `panel.BuiltinAssets()`（`internal/panel/assets.go:43-49`）打开后 **4 条目**（index.html＋两条 hashed assets＋.gitkeep 锚），入口 HTML 1044B、**2 枚 asset 引用全部可解析**——`index.html`／`assets/*.js`／`assets/*.css` 三类关键入口全在，非 `.gitkeep` 占位形态。
- 跑完已删本地产物：`frontend/dist` 现量只有 `.gitkeep`（`git status --porcelain -- frontend/dist`＝0 行），**dist 字节零提交**。
- 注：`internal/panel` 整包本轮 **ok 23.134s（全绿）**，其中票 248 §6 时点的四枚常红之一 `TestComposerContractTypesMatchFrontend`（Go 侧键声明与前端 TS 的双向对账门）**现绿**——AC#2 的界面侧对账半已闭合。

## 3. AC#2 —— 哨兵喂写路径＋全仓 grep

- 常驻用例（已在仓，本轮经载具实测）：
  - `TestAC2CredentialSentinelAppearsInNoArtifact`（`cmd/wisp/panel_config_248_test.go:154`）——哨兵走**真写腿**（真 Manager＋真 secret.NewStore，零 fake），扫六面：receipt／[audit] 行／slog sink／快照／ledger 行／数据根全部文件名＋字节 ⇒ **0 命中**；并正向断言写真发生（回执带 `dpapi:`、config.toml 落 `api_key_ref = 'dpapi:`）。
  - 正控 `TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere`（`:218`）——同一把尺指到真带哨兵的产物必红，尺的视力有背书。
- 载具读数（`scripts/wisp-cli-staged-tests.sh`，一次一发，**exit 0**；日志 `.scratch/wisp/probes/orch/20261116/txt/248-carrier-acpins-run2.txt`）：9 枚全 PASS——
  `TestAC2CredentialSentinelAppearsInNoArtifact`／`TestAC2LeakRulerFiresWhenTheCanaryIsReallyThere`／`TestAC2SharedEnvelopeCannotCarryTheCredentialValue`／`TestAC2SnapshotReportsRefWithoutBlobAsAPartState`／`TestAC1InboundLegAnswersSettingsRouteEndToEnd`／`TestAC1SettingsKeysAreTheOnlyNewOnesOnTheSharedEnvelope`／`TestAC1OtherDoorsStillRefuseByNameAfterTheSettingsSocket`／`TestAC3CredentialLivesInOneStoreAndConfigKeepsOnlyARef`／`TestAC7InvalidSettingsValueIsRefusedBeforeTheFile`。
- 全仓 grep（哨兵值按其原值做 fixed-string 搜，tracked＋untracked）⇒ **5 处命中，全部是哨兵常量自身的定义／镜像／文书引用，无一来自写路径产物**：
  `cmd/wisp/panel_config_248_test.go:39`（定义）、`internal/panel/inbound_raw_leak_35r7_test.go:51`（票 35 的镜像尺，有意复用同一哨兵）、`.scratch/wisp/probes/35/a6/verdict.md:211`、`.scratch/wisp/probes/35/v7/logs/q1-test-copy.txt:51`（源码抄本）、`docs/evidence/s1/35-recon-20261115.md:396`（文书引用）。`logs/`、任何数据根、任何快照产物 **0 命中**。
- ⚠ 载具操作坑（记给下一枚腿）：经 pwsh 调 `bash.exe` 时 `-test.run` 会被拆成 `-test`＋`.run`（症状＝测试二进制打整页 flag help＋`flag provided but not defined: -test`）——**flag 与 pattern 都要单引号成参**；另 `Select-Object -First/-Last` 截断管道会中途杀掉 bash，读输出要么落文件要么全量收。run1 的失败现场在 `248-carrier-acpins-run1.txt`。

## 4. AC#8 —— 回执文案判据现读

- 生产面 **`保存即生效` 0 命中**（全仓 *.go grep 仅 4 处，全在测试的**反向断言**里：`cmd/wisp/panel_config_248_test.go:436`、`cmd/wisp/panel_config_313_test.go:82`、`internal/panel/config_route_248_test.go:324/:380`）。
- 回执答「需重启」：`TestAC1InboundLegAnswersSettingsRouteEndToEnd` 断言 stdout 含「重启」、不含「保存即生效／立即生效」，PASS（见上）。
- 档位同源那一半（票面 §8 判据：`config.TierOf` 驱动而非硬填）：票 255-r2 已落 `cmd/wisp/config_readers_255.go`（TierOf 首个生产调用者）＋`restart_tier_keys_255r2_test.go`（restart 句与 TierRegistry 单源对钉）；设置回执的 tier 现量仍是 `panel_config_store.go:412/:487` 静态盖 restart（`[ball]/[panel]` 两类有实测/Note 例外）。**AC#8 整格仍不勾**——「到达页面可见面」那半的归属与翻转在票面 §8 具名（Go→页返回值无人接＝`frontend/**`；请求↔回执配对＝票 35），非本腿射程，照实留给非实现者裁。

## 5. AC#10 —— 归口具名（ⓑ/ⓘ 之争已被后续票落定）

票面 AC#10 的前提「常驻腿吃不到 `[risk]` 两项」**今天不再成立**：
- `cmd/wisp/resident_approval_windows.go:601` `newResidentApprovalWithConfig`——经 `residentRiskGateValues(dataDir)`（`:606`）真读 `[risk]`，`Window`／`ApprovalTimeout` 两枚喂进 `approval.New`（`:608-612`），唯一生产调用点在 resident boot；**票 256 AC#1 落地**（测试 `cmd/wisp/resident_approval_risk_256_windows_test.go` 在整包跑）。
- `Options.Grants`＝`residentGrantHolder`（`:605`，票 265 落地、裁定 A601 §4 形 ⓐ-Ⅰ）：「本会话内允许」有真落点，`GRANT-DROPPED` 的常驻形态已废。
- 限制照登记表写死（`TestTicket256ResidentGateRiskValuesAreConstructionTimeOnly` 钉着）：门在**本次启动构造时**跟 `[risk]`，改 config.toml 不热跟——`[risk]` 是 locked tier，与 AC#8 的「不许说保存即生效」同方向。
- 设置页射程照 J5：`risk.*` 全族拒写（`internal/panel/config_handlers.go` 拒写族），设置回执无 `[risk]` 两项可写 ⇒ ⓑ 那句「常量 300s／3s」的回执文案前提已被 256/265 消掉。**AC#10 整格的翻勾归非实现者**，归口＝票 256（[risk] 消费）＋票 265（Grants holder）＋台账 A601。

## 6. 禁区自查

- `git status --porcelain -- cmd internal`＝**0 行**（本轮零产码改动，`Remove-Object` 失败未伤任何文件；frontend 源码一字未动——`frontend/fixtures/composer-states.html` 的既有 M 是 2026-11-11 16:59:50 的旧痕，非本轮，未触碰）。
- ⛔ 未动 thresholds/golden/allowlist/PLAN/specs/D43/design；⛔ 未放宽断言；⛔ 未 push；⛔ 未动 C17 白名单形状（`internal/panel/l2_grant_boundary_test.go` 冻结件一字未改，:2202 起的 guardRosterOf 形是此前具名解冻波次的存量）。
- 本轮写面＝本文件＋两枚载具日志（`.scratch/wisp/probes/orch/20261116/txt/248-carrier-acpins-run{1,2}.txt`）。
