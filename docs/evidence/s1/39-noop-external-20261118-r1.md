# 票 39 AC#4 收件（no-op 保存零 diff ＋ 外部改/未知键行号报错 ＋ 保存冲突 refused-drift）— r1，2026-11-18

写腿：agent-130（GLM 5.3 Flash 路由组）。派单：`.scratch/wisp/probes/orch/20261118/txt/39-noop-external-dispatch.txt`。
写面：`internal/config/**`、`internal/panel/**`、`cmd/wisp/**`（唯一）、`frontend/src/**`、`frontend/scripts/render-*`。

## 三条各自读数与牙名

### 1. no-op 保存零 diff（判据：写前写后 sha256 相同 ＋ mtime 不变；⛔墙钟猜想，判据全是可判事实）

- **机制（既有，本轮钉牙）**：`internal/config/writeguard.go:117` `mergeWrite` 的第 2 分支——`diffKeyPaths(disk, base)` 为空（盘上已带该值）时**直接 return，不调 SaveFile**，文件一个字节都不动（日志行逐字：`config.toml already carries this key; this write added nothing (the file's bytes and mtime stayed put …)`）。回执侧 `internal/config/sectionwrite.go:360` `writtenKeyPaths` 答空列表 ⇒ `panel.SettingWriteResult.Written` 空 ⇒ `internal/panel/config_handlers.go:659` 收据逐字 `本次没有任何键路径发生变化（文件里已经是这个值）。`
- **牙 1（config 层）**：`internal/config/fingerprint_39_test.go` `Test39NoopSectionWriteLeavesBytesAndMtimeUntouched`——夹具带手写注释＋空行，`SetBallSize(90)`（盘上已是 90）⇒ 收据空、sha256 前后相同、mtime 相同、字节逐位相同、两条注释仍在。**同案内正控**：`SetBallSize(120)` 真改 ⇒ sha 变（证明该 harness 分得清 no-op 与真写——no-op 走成真重写⇒此案红）。
- **牙 2（路由层，真门→真腿→真写路径）**：`cmd/wisp/noop_external_39_test.go` `Test39NoopSaveLeavesFileByteIdenticalAndSaysSo`——`panel.ConfigWriteHandler` 收 `config.set ball_size=90`（盘上 90）⇒ 收据含 `本次没有任何键路径发生变化`、事件 `written=[]`、字节逐位相同＋mtime 相同＋注释保留；正控同键真改 120 ⇒ 事件 `written=[ball.size]`、字节变。
- **载具读数**：`--- PASS: Test39NoopSaveLeavesFileByteIdenticalAndSaysSo (0.02s)`。

### 2. 外部改＋未知键 ⇒ errors[]{line,message} 带 1-based 行号并点名该键

- **机制（readface wave 1 已落，本轮复核接线＋补牙）**：`internal/config/loaderrors_456.go`（`ParseErrors*` → `decodeErrorRows`：StrictMissingError 逐键一行带其行号；DecodeError 一行带 line+col）→ `internal/panel/config_handlers.go` `SettingsView.LoadErrors`/`configErrorPayload{line,message}` → `cmd/wisp/panel_config_store.go` `configLoadErrors456`。前端 `ConfigErrorList`（`frontend/src/components/config-editor.tsx:176`）按行号列出，line=0 画 `（无行号）`，不折第 1 行（agent-109 已落，本轮未动）。
- **牙（路由层，双向）**：`cmd/wisp/noop_external_39_test.go` `Test39ExternalEditErrorsCarryRealLinesAndNameTheKeys`——真 rig＝**好配置先 NewManager、盘上文件后被外部改坏**（进程照常跑在 last good config 上）：
  - 未知键 `bogus_key = 1` 在**第 5 行** ⇒ 读面与 wire payload（`"errors":[…]`）都答 `line=5` 且 message 点名 `bogus_key`；**每个解码行 `line==0 || line==1` 即红**（行号吞成 0/1⇒红）；**零错误行即红**（未知键吞掉⇒红）；Section 树照常答（运行面未死）。
  - 类型不符 `size = "not-a-number"` 在**第 4 行** ⇒ 恰有行指向 4。
- **载具读数**：`--- PASS: Test39ExternalEditErrorsCarryRealLinesAndNameTheKeys (0.01s)`（既有 `Test39ConfigStoreReadFaceSectionsErrors` 同发 PASS 0.00s）。

### 3. 外部改＋保存冲突 ⇒ 拒绝静默覆盖（`refused-drift`，明示「盘上已变」），给读数

- **机制（本轮新落，主流 harness 的乐观并发 / If-Match 形状）**：
  - 读半：`internal/config/fingerprint_39.go` `FileSHA256()`（盘上原始字节的 sha256 十六进制；缺文件答 `("", nil)`＝无指纹；读不了答错）。`config.get` 载荷带 `baseSHA`（omitempty：腿没答＝键不上线，从不造 hash）。
  - 写半：信封 `configBaseSHA`（`internal/panel/bridge.go`）→ `SettingWrite.BaseSHA` → `cmd/wisp/panel_config_store.go` ApplySetting 开头核对：**非空且 ≠ 盘上现 hash ⇒ `observe.ClassConfig` 拒绝**，逐字 `config: refused-drift：盘上已变——config.toml 在这次读取之后又被改过（外部编辑或另一处写入），这次保存拿的是旧视图，拒绝静默覆盖：一个键都没写。请重读配置确认盘上现值后再保存。` 空 echo＝页没持有指纹（旧页形）⇒不核对，pre-39 形逐字不变。
  - 闭环：接受写的事件 `configReloadEvent` 带 `baseSHA`＝**写后指纹**（`SettingWriteResult.BaseSHAAfter`），前端缓存随事件前移 ⇒ 连续两发保存不会撞上自己前一发造成的「漂」；被拒写零事件，缓存停在原处（被拒本身就是旧视图）。
  - 前端：`frontend/src/lib/config-base.ts`（单源缓存）＋`config-view.ts`（读时记忆）＋`config-actions.ts`（四枚写助手都回声 echo）＋`panel.ts`（信封键）＋`config-feedback.ts`（拒绝句逐字匹配 `盘上已变` ⇒ 一句人话：`配置文件刚被外部修改过，这次保存没有执行。请重新打开本页读取最新配置后再保存。`）。
- **牙（路由层）**：`cmd/wisp/noop_external_39_test.go` `Test39ExternalDriftRefusesTheStaleSave`——真 rig：读 ⇒ 拿 `sha0` ⇒ **外部手改**（ball.size 90→80＋手写注释行）⇒ 页拿旧指纹旧值（90）保存 ⇒ **拒绝**（err 含 `盘上已变`＋`refused-drift`），**盘上字节逐位仍是外部版本、mtime 不动**（外部改没被吞）⇒ 页重读拿新指纹再存 80 ⇒ 放行，事件 baseSHA＝写后真 hash；旧形（无指纹信封）照常写（旧页兼容，钉成故意）。面板层三枚牙在 `internal/panel/drift_39_test.go`（信封→store 透传／payload 只在腿答了才带 baseSHA／事件同）。
- **载具读数**：`--- PASS: Test39ExternalDriftRefusesTheStaleSave (0.02s)`。

## 三样现跑读数

1. **d22scan**：`d22scan.sh: OK - control rc=0, scan rc=0`；末行 `clean - no D22 ban violations`（0 findings）。
2. **载具**（`sh scripts/wisp-cli-staged-tests.sh`，跑前 `Get-Process wisp.test` 空）：五案一发全绿——本腿三枚新牙＋`Test39ConfigStoreReadFaceSectionsErrors`＋`Test37VerdictBatchFieldsReachThePacket`，`ok github.com/CarlosShao/wisp/cmd/wisp`。包测直跑：`go test ./internal/config ./internal/panel -count=1` ＝ 2×`ok`。
3. **go build ./...** rc=0。
4. **前端**：`typecheck` rc=0｜`lint` rc=0（45 警告 0 错误，**触过的文件零警告**，全为存量他票）｜`build` rc=0｜仪器：`render-config-screen` rc=0（vite --ssr＋node 手动跑，无 npm script，历史如此）——新增 drift/no-op 牙全绿（缓存 hold/clear、baseSHA 解析透传、事件带写后指纹并前移缓存、信封回声 configBaseSHA／无指纹时键不上线、盘上已变拒绝句、no-op 没有变化句）；`render:cfggolden` rc=0（18 节 × 115 行 byte-verified）｜`render:composer`／`render:model-settings` rc=0（回归核）。

## 变更清单（11 文件）

- `internal/config/fingerprint_39.go`（新）＋ `fingerprint_39_test.go`（新，两枚牙）
- `internal/panel/config_handlers.go`（SettingsView/SettingWrite/SettingWriteResult 的 BaseSHA 三处＋载荷 baseSHA＋事件 baseSHA）
- `internal/panel/bridge.go`（信封 `configBaseSHA`）
- `internal/panel/drift_39_test.go`（新，三枚牙）
- `cmd/wisp/panel_config_store.go`（读半填 BaseSHA／写半 refused-drift／写后 BaseSHAAfter）
- `cmd/wisp/noop_external_39_test.go`（新，三枚路由牙）
- `frontend/src/lib/config-base.ts`（新）＋ `config-view.ts` ＋ `config-actions.ts` ＋ `config-feedback.ts` ＋ `panel.ts`
- `frontend/scripts/render-config-screen.tsx`（仪器牙＋OK 句）

## 边界与诚实声明

- 真写（非 no-op）仍按 ticket 226 的既定代价**整文件重排 canonical 形**（注释丢失）——mergeWrite 自己的头注写明，本票不动它；no-op 恰恰因为不写所以格式必保（牙 1 钉死）。
- 指纹核对只武装**持有指纹的页**（读到过 baseSHA 的）；旧信封（无 configBaseSHA）走 pre-39 形，测试把它钉成故意而不是漏。
- AC#4 的勾格仍**留白给裁决者**（裁决者≠实现者）；本轮交付的是机制＋双向牙＋读数。
- ⛔ 未 push；⛔ 未动 thresholds/golden/allowlist/PLAN/specs/D43/design；⛔ 未触碰他人在飞文件（git status 全程只有本腿写面变更）。
