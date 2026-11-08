# 票 230 三格收口读数（AC#3 复核·AC#4 归口现状·AC#5 整包终态三发）

- 腿：`230-cells`（写腿；写面＝本件＋票 230 票尾 Progress log 一条；**零产码、零框翻**）。
- 落刀时刻 2026-11-08 16:4x–17:4x +08。主仓 HEAD＝`8b568eccd`（现读），UDF 修复 `e7f420a6c` 经 `git merge-base --is-ancestor` 确认在祖先里。
- 导出树（门禁与整包全在仓外树跑，共享工作树只读）：
  - tree1＝`git clone --shared` → `/d/tmp/wisp/20261108/230-cells-tree`（HEAD 同 `8b568eccd`）；
  - tree2＝同 HEAD 重导＋`git -c core.autocrlf=false -c core.eol=lf archive HEAD | tar -x` 覆盖（消 autocrlf 检出伪红，见 §3）＋补 DLL 6 枚与 `third_party/spike-models/`（均为仓内未追踪的本机件，照实记）。

## 1. 起手三核（全现跑）

1. **票面串/行还在不在**：票 80 `:61` 交接句逐字「AC#3/AC#4 的接手方 = **票 21**（审批队列）」现读仍在（`80-…-done.md:61`）；票 21 现量 `grep -c '^- \[ \]'`＝**0**、`grep -c '^- \[x\]'`＝**5**、`grep -c blacklist_overrides`＝**0**；票 230 未勾＝AC#3（`:27`）／AC#4（`:35`）／AC#5（`:36`）三格。
2. **引文出处**：`docs/PLAN.md:2407` 逐字「**但豁免必须按文件、必须由人点、必须留痕。**」现读在；`docs/evidence/s1/230-ac3-handoff-disposition.md`（9847B）在，结论＝乙。
3. **红/绿是否别人在飞件**：共享工作树 `git status --porcelain` 的 tracked 脏面＝**空**（读数现跑），AC#4 已由在先腿入库（`eacf043d0` → `3c3c516ce` → `0d809934e`）；门禁与测试全部在导出树跑。

## 2. AC#3 独立复核：**乙成立（维持 230-ac3 处置，本腿不重做）**

230-ac3 腿 2026-11-05 已把本格落成乙（作废交接句＋两格接收在票 230 名下），本腿逐条独立复现其判性，均成立：

1. **甲＝复活被否决的静态预授权通道**：owner `A448②`「黑名单这个按你推荐吧」＋`A53②` 裁「取消」＋`PLAN.md:2407`「必须由人点」逐字在 ⇒ 甲的落法不是落地是复活。
2. **票 21 不是接收方**：票 21 现量 5 枚判据全 `[x]`、`blacklist_overrides` 命中 **0**（本腿现跑，见 §1）；该格裁决面（改契约＋owner 点）不在票 21。
3. **交接句自认"被裁决取消"**：`80-…:61-63` 逐字在。

★**新事实（本腿现读，派单转述已过期）**：机主 2026-11-05 17:4x 已裁 **AC#4＝走 B（台账 `A1549`，票 230 票尾 `:63` 逐字在）**，且接线已由 230-AC4 腿入库：`internal/config/unwired.go` 的 `risk.blacklist_overrides` 响亮失败行已摘除（原 `:79-84`，现读该处为注释「used to sit here…owner ruling 2026-11-05, ledger A1549」）、`internal/config/overrides.go` 在（7601B）。⇒ 乙接收的两格（票 80 AC#3/AC#4 字面）其落地门（本票 AC#4）已由 owner 批并落地；**本腿对两包零码、零框翻**，只补这句现状。

## 3. AC#5 整包终态：**cmd/wisp 仍未到终态（如实报新死点）；internal/** 全部到终态**

命令一律照票面 `go test ./cmd/wisp ./internal/... -count=1`（环境：`unset GOTMPDIR TMPDIR TEMP TMP`；`export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`；**未设** `CGO_ENABLED=0`）。三发读数：

| 发 | 树/条件 | 终态 | 四数 | 要点 |
|---|---|---|---|---|
| ① | tree1 裸共享克隆（默认全局 `core.autocrlf=true`、无 DLL） | rc=1，**2m49s 到终态** | ok=26 FAIL=9 包 | 但 cmd/wisp／audio/aecspike／speech 三包死因＝`exit status 0xc0000135`（DLL 未追踪不入克隆＝**环境死**，非产品红）；ball 3 红／panel-composer 红＝autocrlf 把 `.html` 资产检出成 CRLF（`.gitattributes` 只钉了 .go/.md/.sh 的 eol）＝**伪红**（blob 现跑 CR=0、orb-siri blob sha256=`d93da649…e3bc6` 与钉值逐字同） |
| ② | tree2（eol=lf archive＋6 DLL，模型未补） | rc=1，internal 全到终态 | ok=28 FAIL=7 包 | **speech 全家红＝env**（缺 `third_party/spike-models/`，报错逐字「set WISP_SPIKE_MODELS, or run from inside the repo」）；tools `TestTicket252R1AlignmentAddsNoAuthorization`＝**env**（临时根 `D:\tmp\wisp\20261108` 无 8.3 短名别名，前置条件缺失不许 skip）；其余见 §4 红名册 |
| ③ | tree2＋spike-models，cmd/wisp+speech `-timeout 25m` | rc=1，**cmd/wisp 仍撞包级超时（1500.217s）** | speech **ok 0.265s**（坐实②的 speech 全红＝env） | **cmd/wisp 没到终态**：25m 闹钟时 `Test258BridgeRebindsLiveKeysFromConfigEdit` 已跑 17m35s 未归；②发 10m 闹钟时死点＝`TestAC294MuteHotkeyRegistersLiveOnARealBallWindow`（1m20s）。**两发两死点、形态一致＝真窗/热键活体测试逐个无限挂起** ⇒ UDF 修复（`e7f420a6c`）解决了"整包挂在启动面"的死点，但 cmd/wisp 的活体测试仍有一个未修的挂起点，**新死点如实上报** |

## 4. 红名册（以发②为准＝internal 终态权威读数；cmd/wisp 取发②③并集）

**环境红（4 名，非产品，均具名归因）**：
- `internal/speech` 全家（发② 21 名全 0.00s）＝缺未追踪 spike-models；发③补件后 **ok**。
- `internal/tools TestTicket252R1AlignmentAddsNoAuthorization`＝临时根无 8.3 别名（主仓路径含空格有别名，此探针在 D:\tmp 根必红）。

**产品红（HEAD `8b568eccd` 真红，逐名；各归各票，本腿零修）**：
| 包 | 测试 | 现场一句 |
|---|---|---|
| internal/panel | `TestC21DesignTokensFourWayAgree` | `tokens.go` 无 `GlassRing/GlassHi/GlassHiSoft` 字面，docs 表 6 行对不上（球/前端 token 地界） |
| internal/plugin | `Test52SchemaViolationsAreRejected`（3 子测试） | `{"days":0}/{"days":31}/{"days":1.5}` 被接受 |
| internal/plugin | `Test52CommandKindValidates/unknown_kind` | C14 报错未按 `.kind` 指名 |
| internal/plugin/manifestslot | `Test52TimeoutKillsTheProcessTree` | 孙进程 61188 超时后仍活着 |
| internal/proc | `TestCrossVetForLinux` | `internal/ball/orbdebug.go:39 undefined: Ball`（GOOS=linux；**球腿地界**） |
| cmd/wisp | `TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`＋`TestTicket255RosterStillMatchesTheActualReadSites` | 255 名册漂移（如 `SizePx 44..72` 已被 65 腿改成 `44..160`；`ball` 行 roster 空而 scan 两枚） |
| cmd/wisp | `TestTicket198FirstRunCallerIsTheRunEntryOnly` | `ensureFirstRunConfig` 调用者＝`[runTextTask runConfigSet65]`，want `[runTextTask]` |
| cmd/wisp | `Test313WritableUnionEqualsTheTwoRosters` | union 多出 `ball_style`（65 腿新键未进 313 两张名册，A881 尺） |
| cmd/wisp | `TestAC4FocusReturnToPriorWindowGap33r5` | 真窗 Show 未拿到前台（0xb5a0f02 始终是前台）——**活体环境敏感** |
| cmd/wisp | `TestLockedSuffixNeverTakesTheLockItself33r11` | `panelHolder.stampPublishLocked`（resident_windows.go:949）带 Locked 后缀不在名册 |
| cmd/wisp | `TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims` | resident 腿喂 `Options.Answer`，7/11 字段 want 6 |
| cmd/wisp | `TestAC247HandingTheLevelToTheSeamIsNotVisibility` | 装配里 prototypeVisuals 为 ON，ship 默认须 off |
| cmd/wisp | **挂起（非 FAIL 名）** | `TestAC294MuteHotkeyRegistersLiveOnARealBallWindow`（②发 10m 闹钟）/`Test258BridgeRebindsLiveKeysFromConfigEdit`（③发 25m 闹钟，17m35s）＝**新死点** |

对照票面 AC#5 的历史在册（「`internal/ball` 1＋`internal/panel` 4 属别人地界」）：本读数现量 ball **0 红**（发②③绿）、panel **1 红**——历史名册已过期，逐名以本表现量为准；`internal/risk TestResolvePerCallBudget` 本读数**绿**（`ok internal/risk 47.739s`，未复现争用假红）。

## 5. 门禁（tree1，HEAD 同串）

| 尺 | 读数 |
|---|---|
| gofumpt（`/d/work/base/gopath/bin/gofumpt.exe -l`，tracked 1282 枚 xargs） | 本票相关包（internal/agent/approval、internal/config、internal/risk、cmd/wisp、internal/winsec）定向＝**空**；全仓列出 **9 枚全在 `forks/go-webview2/**`**＝票 65 腿 2026-11-07 落的上游 fork 存量（`2bb639e4e`），非本票新增 |
| `go vet ./cmd/... ./internal/...` | rc=0 |
| `go build ./...` | rc=0 |
| `sh scripts/d22scan.sh` | rc=0，clean（internal/=964、cmd/=241、frontend/=121 等全扫） |

## 6. 本腿写面与硬闸

- 改动＝本件＋票 230 票尾 Progress log 一条（append-only）。**零 `.go`、零框翻、零 Status 改、零 push**；`frontend/**`／`internal/ball/**`／`internal/config/**`／`PLAN.md`／specs／thresholds／golden／allowlist 零触碰。
- AC#3/AC#4/AC#5 三框的翻勾权归编排者收件；本件只交读数与归因。
- 原始日志：`/d/tmp/wisp/20261108/230-cells-logs/fullpkg.txt`（发①）、`fullpkg2.txt`（发②）、`fullpkg3.txt`（发③）、`gofumpt-out.txt`、`vet.txt`、`build.txt`、`d22scan.txt`。
