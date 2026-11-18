# 票 40 AC#5 — Diagnostics preview＝精确清单；导出包实证排除 音频/密钥/完整转写（写腿 agent-134，2026-11-18）

**票面原句**（`40-security-privacy-cost-pages.md:52-53`）：`Diagnostics preview lists exactly the bundle
contents; exported bundle verified to exclude audio/keys/transcript-full.`

**交件 commit**：见回话第一行（本件随码同笔入库）。
**写面**：`internal/panel/**`（名册＋派发表＋wave-5 牙）、`cmd/wisp/**`（面板诊断腿＋装配＋测试＋255 读点名册＋35g1 变异表）、
`frontend/src/**`＋`frontend/scripts/render-*`＋`frontend/package.json`（诊断页＋仪器）、本证据件。
⛔ 未动 thresholds/golden/allowlist/PLAN/specs/D43/design；⛔ 未 push；⛔ 未动 `frontend/fixtures/composer-states.html`
（它保持 M，只报不改）；⛔ 探针全为假串（`WISP40*` 标记＋`sk-wisp40probe9f3` 形状），零真凭据进仓或进盘。

---

## 1. 真清单：一份文件列表来源，不是两套

现有诊断包装配点（派单第 1 步的现量，`grep -rn "diagnostic|bundle" --include=*.go internal/ cmd/ | grep -v _test`）：

- **引擎**（票 45）：`internal/observe/diagnostics.go` — `PreviewDiagnosticsBundle`／`BuildDiagnosticsBundle`
  是**兄弟函数**，读同一份 `BundleOptions`（同一 LogDir/ConfigPath、同一 skip 规则、同一 per-item 同意规则），
  preview 的 `Entries` 就是 builder 的 zip 条目清单（`guards_45_test.go` 的
  `TestAC1BundleContentsAreExactlyAsPreviewed` 已钉引擎层）。
- **CLI 面**（deferred disposition 20261110 #5）：`cmd/wisp/diagnostics_cmd.go` — `diagnosticsBundleOptions180P`
  是唯一 options 组装器，`[privacy] diagnostics_opt_in` 总闸。
- **面板面**（本腿，此前零命中）：`internal/panel/bridge.go` 新增
  `MethodDiagnosticsPreview = "diagnostics.preview"`／`MethodDiagnosticsExport = "diagnostics.export"`
  （**export 用的是 SPEC-08 §5.2:173 自己的拼写**；`diagnostics.preview` 按 CLI 子命令动词命名）→
  `composer_dispatch.go`（`DiagnosticsHandler` 两方法＋两 case）→ `cmd/wisp/panel_inbound.go` 装配
  `newDiagnosticsDoor40(dir, …)` → **cmd/wisp/panel_diagnostics_40.go**。

**⛔ 预览写死、导出另一套的那条线是这样堵死的**：面板门自己**零枚举**——两个 handler 都调
`diagnosticsBundleOptions180P(h.dir, c, outPath, "")`（与 CLI 同一个组装器、同一份数据根），再交给同一个
`observe.PreviewDiagnosticsBundle`／`observe.BuildDiagnosticsBundle` 对。门文件里出现
`os.ReadDir`／`zip.NewWriter`／`zip.OpenReader`／`filepath.Walk*` 任一个都会被源级牙
`Test40DiagnosticsOneEnumerateSourceNoSecondFileList` 当场打红（正控：三个必须存在的调用字面量逐一钉住）。
导出 handler 先跑 preview、失败即中止（builder 本身 fail-closed 且不留半成品 zip），页面上看到的清单
**就是** builder 出货的那份。

## 2. 排除实证（牙②）：探针明文喂进源目录，zip 解包逐字节 grep

判据照派单：**不是看白名单，是解包 grep 探针明文**。夹具 `seedDiagnosticsDataDir40` 种进数据根：

| 类 | 负探针（必须 0 命中） | 种在哪 | 正控（必须 ≥1 命中，证 grep 真能咬这类字节） |
|---|---|---|---|
| 音频 | `WISP40AUDIOPROBE-9f3` | log 目录里的 `wisp-audioprobe-9f3.wav`（**builder 唯一会走的目录**，prefix/ext 过滤必须拦住它）＋ 数据根 `audio-capture-probe.opus` | `WISP40CTRL-LOG-9f3` 写进 `wisp-*.jsonl`（随包进 zip ⇒ 命中） |
| 密钥（Store 值/落盘散料） | `WISP40SECRET-PROBE-9f3` | 数据根 `secrets/probe-note.txt`（builder 只读 log 目录与 config 两个输入） | 同上行（字节级搜索能命中包内字节这件事由两条正控共同证明） |
| 密钥（config 明文 key 形） | `sk-wisp40probe9f3` | `config.toml` 注释行（脱敏器把它掩成 `sk-e9f3`，明文不得出货） | `WISP40CTRL-CONFIG-9f3` 同文件注释（随 `config.redacted.toml` 进 zip ⇒ 命中） |
| 完整转写 | `WISP40TRANSFULL-9f3` | 数据根 `transcripts/session-full.txt`，**未逐项同意**（面板装配 TranscriptDir 恒零，引擎对种了目录但未 opt-in 的形也必须排除） | **同一探针**走显式逐项 opt-in 构建 ⇒ 在 `transcripts/session-full.txt` 条目里命中（证 grep 咬得住转写本文字节） |

读数：`Test40DiagnosticsProbesDoNotComeOutTheOtherSide` **PASS**（仓内载具形 0.19s；4 条负读数全 0 命中、
3 条正控全命中——正控哪条失守都写明“grep 无法证明该类排除”）。

## 3. 三枚牙

**牙① 预览与实包文件名集合相等**：`Test40DiagnosticsPreviewEntriesEqualTheZip` ——
经**生产链**（`newComposerDispatchChain` → `ComposerDispatch.Handle` 全封套）先 preview 后 export，
payload `entries` 集合 vs 真实 zip 解包条目集合双向差集=0；export payload 的 entries 也对同一 zip 复核。
**PASS 0.19s**。

**牙② 音频/密钥/完整转写任一放行 ⇒ 红**：见 §2 四条负读数；变异仓外两发都咬在这条牙上：
- **M2**（log 过滤改 copy-all）：红在 `牙② 红 (exclusion broken): the audio probe plaintext
  "WISP40AUDIOPROBE-9f3" came out the other side`（`Test40DiagnosticsProbesDoNotComeOutTheOtherSide`
  FAIL，rc=1）。
- **M2b**（builder 无 opt-in 也打包 transcript 目录）：红在 `牙② 红 (exclusion broken): the un-opted-in
  transcript's probe came out the other side - the default exclusion is gone`（同测试 FAIL，rc=1）。

**牙③ 预览画“空清单=正常” ⇒ 红**：页面上零条目预览必须**大声说**「条目 0 个：这个包当前不携带任何文件」
＋把载荷里各类引擎自报的未收集原因逐类逐字列出（音频那句逐字上屏），并且**不画** zip 条目清单表格；
这些断言钉在仪器 `render:diagnostics40`（`frontend/scripts/render-diagnostics-40.tsx`）。
变异 **M3**（把组件零条目分支的条件掏成 `false`）：仪器 **FAIL 3 枚**，逐名：
`the zero-entry preview must refuse to read as a normal empty list`／
`the zero-entry preview must draw NO entry-list chrome`／
`the zero-entry branch must draw no entry-list chrome`，rc=1。
另有一条源级护栏：页面渲染的是**载荷自己的 detail**，组件源码不许拼引擎的排除句
（`!cardSrc.includes("never collected")`）——页面自己写排除句＝伪造读数。

**牙① 的变异（M1）**（preview 的 addEntry 掏成 no-op）：
红在 `Test40DiagnosticsPreviewEntriesEqualTheZip`（`--- FAIL`，rc=1）——预览清单与 zip 条目集合一脱钩，牙①当场响。

## 4. 变异实验室口径（仓外，零污染工作树）

- 快照：`git archive HEAD`（`9467a91f3`）解到 `%TEMP%\wisp-d40-lab`，覆盖本腿 9 枚 Go 件＋6 枚前端件；
  `third_party/sherpa-onnx` 经 junction 借给快照（仓内未建 worktree、未 checkout）。
- 每发变异：同链 `Select-String` 打印被改后整行（落地证明）→ `go build ./internal/observe/`（或
  `npm run typecheck`）**rc=0 先量到**（编译失败不算变异）→ 跑责任测试读红名 → 逐发从 `.pristine`
  覆盖还原，还原后 build rc=0。
- 读数汇总：**M1 → 牙①测试 FAIL rc=1**；**M2 → 牙②音频句 FAIL rc=1**；**M2b → 牙②转写句 FAIL rc=1**；
  **M3 → render:diagnostics40 FAIL 3 枚 rc=1**。基线（同快照未变异）`Test40Diagnostics*` 5/5 PASS rc=0、
  `render:diagnostics40` rc=0。
- 诚实登记一发废弹：M3 第一发写成 `false || entries.length === 0` ＝ 逻辑恒等的 no-op（rc=0），
  **作废重打**为 `false ?` 后才计红。变异只存在于快照，仓库树从未变异。

## 5. 门禁读数（全部现跑）

- `sh scripts/d22scan.sh`：**clean，0 findings**（control rc=0, scan rc=0）。⚠ 途中抓到并修掉一枚
  phantom-citation：我在 bridge.go 注释里写了 `transcripts/watchdog/audio` 这串斜杠名，仪器当场红
  （`[phantom-citation]`），改为顿号列举后过——仪器照常咬人，记一笔。
- `go build ./...`：**rc=0**。
- 载具（`sh scripts/wisp-cli-tests.sh`，整包 `./cmd/wisp/`，跑前 `Get-Process wisp.test` 确认无在飞）：
  **rc=0 全绿**（日志 `%TEMP%\wisp-d40-vehicle.log`；定向先行：`^Test40Diagnostics` 5/5 PASS rc=0）。
- `go test ./internal/panel/`：**ok 28.5s rc=0**（名册 18→20／非前缀 14→16 两枚摩擦常数随路线同笔移动；
  双向镜像门 `TestEveryWhitelistedInboundNameIsAlsoSentByThePage` 先红后绿——红句点名
  `panel.ts never spells diagnostics.export, diagnostics.preview`，正是它该咬的形状）。
- 前端：`build` rc=0；`typecheck` rc=0；`lint` rc=0（45 条既有 warning，0 error）；`seeds` 0 findings；
  `render:diagnostics40` rc=0；回归 `render:security40`／`render:privacy40`／`render:sessionrules40`／
  `render:nav`／`render:composer`／`render:shell77`／`render:cfggolden` 全 rc=0；dist 未入库（构建产物
  不在 git status）。

## 6. 名册与仪器的摩擦（同笔移动，零断言放宽）

- `inbound_roster_253_test.go`：closed roster 18→20、非 `panel.` 前缀 14→16、roster 切片两枚、
  guardLine 两处、wave-5 注释。
- `bridge_test.go` 的 `rosterPageGapLedger194`：**零条目**——两枚页面字面量
  （`requestDiagnosticsPreview`／`requestDiagnosticsExport`）与 Go 常量**同 commit** 落 `panel.ts`，
  镜像门直接判。
- `cmd/wisp/panel_inbound_fuzz_35g1_test.go`：禁表里的 `diagnostics.export` 按 privacy.purge/privacy.export
  先例**移出**（白名单答的名字不是 fuzz 件），换近形 `diagnostics.export.all`，表保持恰好 50
  （`:318` 的 50 枚钉子现跑仍绿）。
- `cmd/wisp/config_readers_255.go`：`Privacy` 读点名册加入 `panel_diagnostics_40.go`（同一把总闸
  `diagnostics_opt_in` 的面板读点，带 2026-11-18 重裁注释；[privacy] 行的 verdict 不动——仍是
  boot/request-time LoadFile，不读 live Manager）。

## 7. 诚实边界（未勾 AC 格，裁决者≠实现者）

- 票面 AC#5 的框不由本腿勾；对活宿主的真窗 e2e（面板真点击 → 真桥 → 真落盘）不在本腿射程，
  行为面由载具测试＋SSR 仪器覆盖。
- 导出 zip 落在数据根（与 CLI 同默认），面板不传 zip 字节——页面上写明路径，用户从盘上打开；
  面板门没有逐项 opt-in 选择器（CLI 也没有：本构建无磁盘转写目录可诚实指名），逐项同意仍走引擎接口，
  manifest 如实记录。
- 票 40 其余未勾格（AC#6 视觉签收、AC#6b mod 管理行）不在本派单。
