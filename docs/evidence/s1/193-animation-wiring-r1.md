# 193 — 动画档接线腿 r1（泵源两行 ＋ config.set 第九名册 ＋ ApplySetting case ＋ 端到端齿）

- 票：`.scratch/wisp/issues/193-lightweight-must-no-longer-constrain-frontend-animation-and-visuals-make-the-per-panel-gates-named.md`（AC#3 接线步）
- 派单：`.scratch/wisp/probes/orch/20261120/txt/193-wiring-dispatch.txt`
- 前任交件：Go 侧 `f34210795`（`internal/config/settings_animation_193.go` ＋ `internal/panel` 载体/泵源字段/193 载体齿）；界面侧 `acce6ea13`（panel.ts 声明＋真停半）；清单＝`docs/evidence/s1/193-animation-opt-r1.md` §⑤
- 写面：`cmd/wisp/**`（唯一共享编译面写腿）、`internal/panel/**`（名册/写面）、本件。⛔ 未碰：`frontend/**`、`internal/config/**`、PLAN.md、D43、thresholds/golden/allowlist、他人在飞件（`frontend/fixtures/composer-states.html` M、`frontend/scripts/probe-ballsize-render.tsx` ?? 只报不动）。
- 时刻：2026-11-20 22:0x–22:3x +08。
- commit：本件入库的同一枚 commit（回话首行报 hash）。

---

## ① 三步接线（同一 commit，名册与 store case 不分家）

**1. 泵源两行**（照 167/外观波先例：读面有值才发键、无值不发键、绝不发假 off）：
- `cmd/wisp/run.go`：`PumpSources` 字面量 `PanelFontSize` 旁加
  `PanelAnimation: func() panel.PanelAnimationSection { return panel.PanelAnimationSection{Tier: mgr.PanelAnimationOptions().Tier} }`——同一 Manager（`rt.settings` 写的那一个），档位词逐字过线，本字面量不拼写 full/reduce/off。
- `cmd/wisp/panel_resident_pump.go`：`panelFontSizeSectionFromConfig` 旁加映射 `panelAnimationSectionFromConfig`（`m.PanelAnimationOptions().Tier` 逐字），`residentPumpSources65` 的 `cfg != nil` 块内 `src.PanelAnimation = ...` 同形一行。nil Manager ⇒ 读者不装 ⇒ 键缺席（零值不冒充、无值不冒充 off）。

**2. config.set 第九名册**（照 `config_write_39.go`／`WritableFields39()` 先例）：
- 新件 `internal/panel/config_write_193.go`：`FieldPanelAnimation = "panel_animation"`（平面名，点名拼写不收）、`fields193`／`fields193TakesValue`／`WritableFields193()`／`writableField193()`。
- `internal/panel/config_handlers.go`：`knownWritableField` 并集加 `|| writableField193(field)`（第八→第九名册，注释同步）。
- `internal/panel/config_write_313.go`：`acceptedWritableFields()`（拒绝文案的并集）加 `WritableFields193()...`；`takesValueFor` 加 `fields193TakesValue[field]`；注释同步（八→九）。
- `cmd/wisp/panel_config_313_test.go` `Test313WritableUnionEqualsTheTwoRosters`（A881 名册并集钉）：`want = append(want, panel.WritableFields193()...)` 一行——**同步不是放宽**（A-313 教训：名册与 want 追加同 commit；少这一行 half 1 就红）。

**3. ApplySetting case**（`cmd/wisp/panel_config_store.go`）：
- `case panel.FieldPanelAnimation:` → `s.mgr.SetPanelAnimation(strings.TrimSpace(w.Value))`（接缝单源＝`internal/config/settings_animation_193.go`，词表/拒绝句不在此复制，单源＝validate.go 的 `panelAnimationEnumError`）。
- 档位句：与注册表同源（39 腿规矩），问 `registryEffectiveTier39`——**问的是节名**（`[panel]` 在 TierRegistry 是整节行 `tiers.go:44 "panel": TierHot`，无 per-key 行；`TierOf("panel.animation")` 查不到，第一轮被端到端齿抓红后改正）。hot ⇒ `EffectiveNow`（"这一项立即生效。"），Note 写明真实读者节奏：快照每次发送现读 panel.animation，面板下一次收到快照（每次 show／状态变化都会推送）即按新档运行，无需重启进程。**进的是已立即生效列，不是 restart 列，理由**：hot 档写入即落进程内存（config_tier_39.go 头注的 hot 语义），泵读面每次发包现读，端到端齿实测"写完下一个包就带新档位词"（§③）；这与 [panel] 外观波两行（width/font_size）的"restart＋Note"不同——那两行的读者在建窗时，本键的读者就是快照本身。

---

## ② 齿（三发，全部实跑：红 → 还原 → 绿）

| 牙 | 变异 | 红读数 | 还原 |
|---|---|---|---|
| M1 泵源假 off | `panelAnimationSectionFromConfig` 改成恒答 `Tier:"off"`（无值/有值都冒充 off） | 载具 `-test.run ^Test193`：`--- FAIL: Test193PanelAnimationWriteRidesTheNextSnapshot`，红句＝`the pre-write snapshot's panelAnimation.tier = "off", want the schema default "reduce" read verbatim off the config the chain holds`（写都还没写就冒充 off，逐字档位断言抓死）；`Test65ResidentPumpNilManagerOmitsConfigKeys` 的 nil-Manager 缺席钉同步加 `panelAnimation` 作第二道（无 Manager ⇒ 无键 ⇒ 无假 off） | 还原后载具全绿 |
| M2 名册并集 | `knownWritableField` 删 `|| writableField193(field)`（门不再 admit） | `go test ./internal/panel -run '^Test193|^TestPanelAnimation'`：`--- FAIL: Test193RosterJoinsTheUnionWithoutMovingOlderPins` ＋ `--- FAIL: Test193DoorAdmitsPanelAnimationAndRoutesItToTheLeg`（门在活机器上会拒掉动画开关）；cmd/wisp 侧 A881 钉 half 2 同形状红 | 还原后 `internal/panel` 全包 ok 24.548s |
| M3 写路由错指 | `ApplySetting` 删 `case panel.FieldPanelAnimation`（路由指向不存在的腿，落 default 拒绝） | 载具 `-test.run ^Test193`：`--- FAIL: Test193PanelAnimationWriteRidesTheNextSnapshot`（config.set 被拒，"设置页写 panel_animation=off ⇒ 下一次快照 tier==off" 断链） | 还原后载具 `^Test193|^Test248|^Test313` 全绿 |

三发牙都在证据记录里留了实跑读数（非注释承诺）；还原逐字核过（git diff 只剩交付内容）。

## ③ 端到端齿（真链路，非 config 层）

`cmd/wisp/panel_animation_wiring_193_windows_test.go`：
- `Test193PanelAnimationWriteRidesTheNextSnapshot`：真装配（`newResidentPanelManager` ＝ 常驻靴自己的链）→ 经**生产设置门**（`host.disp` → `configStore.ApplySetting` → `SetPanelAnimation` 受守写路径）写 `panel_animation=off` → 同一 Manager 的真泵源（`residentPumpSources65`）**下一次** publish 的线上字节解出 `panelAnimation.tier == "off"`；盘上断言 `config.toml` 携带 `animation = 'off'`；收据断言含 `panel.animation`、`立即生效`，不含 `需要重启`。反向词表腿：写 `full` ⇒ 下一个包 `tier=="full"`（恒 off 的接线在这条腿红）。负控：`panel_animation_typo`（无名册之名）⇒ 门按 `不在可写名册内` 拒、零字节动、活档位不动。
- `Test193AssemblyRootsCarryTheAnimationReader`：两处 PumpSources 装配点的读者行源码钉（外观波 `TestWaveAssemblyRootsCarryTheReaders` 同款）——删任一行即红。
- `Test193PanelAnimationMutationStrippedReaderIsAbsent`：真源剥掉 `PanelAnimation` 读者 ⇒ 键从线上消失（reader-gated 在接线层的牙，65 ballStyle 剥钉同款）。
- internal/panel 半：`config_write_193_test.go` 两枚（名册精确集/并集 admit/各老名册不相交/取值列/无选择器列＋门路由齿＋空值负控＋点拼写负控）。

---

## ④ 三样现跑读数（最终，还原牙后）

| 尺 | 读数 |
|---|---|
| ① `sh scripts/d22scan.sh` | 末行 `d22scan.sh: OK - control rc=0, scan rc=0`；`d22scan: clean - no D22 ban violations`（0 findings；ban #8 cmd/＝360 Go files, internal/＝1075, frontend/＝191, design/＝39） |
| ② 载具 `scripts/wisp-cli-staged-tests.sh -test.run '^Test193|^Test248|^Test313'` | **RUN=15 PASS=15 FAIL=0 SKIP=0**，`ok github.com/CarlosShao/wisp/cmd/wisp 11.00s`，rc=0（A881 名册并集那条 `Test313WritableUnionEqualsTheTwoRosters` 含在其中，绿） |
| ③ `go build ./...` | rc=0 |
| 附加 `go test ./internal/panel -run '^Test193|^TestPanelAnimation' -count=1` | `ok … 0.048s`，rc=0（RUN=5 PASS=5 FAIL=0 SKIP=0：两枚新 193 名册齿＋三枚既有载体齿） |
| 附加 `go test ./internal/panel -count=1`（全包保险） | `ok … 24.548s`，rc=0 |

载具规矩遵守：每发前 `Get-Process wisp.test` 现查（三轮全部 clear，无并发）；经 `scripts/wisp-cli-staged-tests.sh` 载具（pinned DLL 挨着测试二进制 staging，sha256 验过）。⛔ 未裸跑。

## ⑤ 过程诚实记账

- 第一轮端到端齿红了一次：收据答了 restart——根因＝`TierOf("panel.animation")` 查不到（[panel] 是整节行、无 per-key 行），改成问节名后绿。这是齿在干活，不是放宽（断言未动，代码改对）。
- 首发载具两连"只跑 2 枚"＝pwsh→bash 引号串味（`|` 被 bash 当管道/引号进值），改走 `bash -c '… "^Test193|…"'` 后 15 枚齐跑；属调用方引号事故，非包缺陷。
- `gofmt -l cmd/wisp/run.go` 命中一处**存量**不齐（`snapSeen bool` 注释对齐，HEAD 逐字验证非本腿引入；本腿 diff 纯 +15 行，未触碰）。本腿触碰的其余文件 gofmt 全净。
- 255 读者扫描（`TestTicket255RosterStillMatchesTheActualReadSites`）**未动、绿**：新读点走 `PanelAnimationOptions()` 方法形（词级扫描器不键方法调用，`readsFieldAccess` 的文档明写的残差）；"Panel" 读点列保持 `{"cmd/wisp/panel_resident_windows.go"}`，无需重裁。
- `config_tier_39.go` 只同步了 SCOPE 注释（第二名册消费者具名），helper 本体与 39 字段行为零改动。

## ⑥ 拟 Progress log（供编排者追加票面，本腿写面不含 .scratch）

- 2026-11-20 22:3x `193-wiring`（agent-164）交件：AC#3 接线三步同 commit 落齐——泵源两行（run.go＋常驻泵，逐字档位词、nil 不装读者）、第九名册（config_write_193.go＋门并集＋拒绝文案并集＋取值列＋A881 钉同步）、ApplySetting case（走 SetPanelAnimation，hot⇒已立即生效列，Note 写明快照节奏）；端到端齿真链路绿（写 off⇒下一包 tier==off，写 full⇒full）；三发牙实跑全红后还原。AC#3 接线半齐；页面消费半＝界面腿 acce6ea13 已交件。未动：frontend/**、internal/config/**、PLAN.md、thresholds/golden/allowlist、他人在飞件。
