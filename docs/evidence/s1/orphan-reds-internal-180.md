# 180 — internal/** 孤儿红修复（98 归因的射程核对 ＋ 两枚孤儿红转绿）

**票/轮**：ORPHAN-RED REPAIR · 180 轮（派单只给「internal/**，⛔ cmd/wisp 及其测试」）。
**锚定树**：`3ff1b1395`（HEAD）。**时间**：2026-11-05 13:00–13:2x +08。
**Status**：证据件（非实现者裁决另派）。本件记录的改动只有三处：
`internal/proc/treemetrics_windows.go`、`internal/agent/reminder_timer.go` 两处生产代码，
`internal/proc/snapshot_growth_180_test.go`、`internal/agent/reminder_disarm_race_180_test.go` 两枚新测试，
以及本文件。**未碰** `cmd/wisp/**`、`internal/risk/**`（冻结）、`internal/panel/**`、`frontend/**`、`design/**`。

## 0. 一句话结论

1. **`docs/evidence/s1/98-attribution.md` 归因的 81 枚红，没有一枚住在 `internal/**`** —— 81/81 的「红在哪一行」
   都是 `cmd/wisp/*_test.go`（§1，可逐字重跑）。它文里出现的 `internal/...` 全是**根因站点**（谁真读了 config、
   谁的桩拒绝），不是红用例。其中唯一落在本腿写面外沿的 C6 根因是 `internal/risk`（**冻结**）。
2. 本机 `go test ./internal/... -count=1` 在 HEAD 上**确实有红**：4 个包 5 枚（§2）。逐枚归因后
   **2 枚是本腿可修的孤儿红、3 枚撞 STOP**（`frontend/**` 或 `internal/risk` 冻结）。
3. 两枚孤儿红**修绿**，断言零放宽（§3）：
   - `internal/proc` `TestExternalSamplerTracksSubjectGrowth` —— 根因是**票 66 残口①**（票面上逐字写「未修」）；
     修的是 `NtQuerySystemInformation` 的**重试规则**（照抄 required ⇒ 系统在长就永远差一步）。
   - `internal/agent` `Test31DisarmStopsTheWait` —— 账册 A1477 逐字写「**未在任何账册/票面出现**，疑计时敏感/待核」；
     实测**不是计时敏感，是 select 在 `stop`＋`ready` 同时就绪时随机取**（`-count=4` 得 3 红 1 绿）。
4. 门禁（§5）：`go build ./internal/...` rc=0、两包 `go test -count=1` ok、`gofumpt -l` 空、`sh scripts/d22scan.sh` rc=0。

## 1. 前提核对：98 归因的 81 枚红，住在哪？

尺（逐字重跑）：把 §5 那张 81 行表的「红在哪一行」列抽出来，看有没有以 `internal/` 开头的。

```bash
cd "/d/work/workspace/projects plans/Wisp"
sed -n '113,193p' docs/evidence/s1/98-attribution.md | grep -c '^| [0-9]'      # 81 行
sed -n '113,193p' docs/evidence/s1/98-attribution.md | grep -c 'internal/'     # 0
sed -n '113,193p' docs/evidence/s1/98-attribution.md | awk -F'|' '{print $4}' | grep -c 'cmd/wisp/'
```

读数：**81 行、`internal/` 命中 0、`cmd/wisp/` 命中 81**。

⇒ 98 号件对本腿的写面**没有直接指名任何一枚红**。它文里的 `internal/` 六处，逐处核过，都不是红用例：

| 98 号件的 `internal/` 出现处 | 那是什么 | 对本腿意味着什么 |
|---|---|---|
| `internal/secret/protect_other.go:18`（C1，76 枚） | POSIX 上 DPAPI 的 **documented fail-closed 桩**；红的是 `cmd/wisp` 夹具写 `dpapi:` ref 后 `t.Fatal` | 平台规则归票 06/08；修法（换 `env:` ref 或加 `//go:build windows`）归**各用例所在的票**，且写面全在 `cmd/wisp` ⇒ 不在本腿 |
| `internal/agent/costbook_44.go`、`quota_44r4.go`、`internal/speech/cloudcatalog.go`（C2） | **真读者**（名册该登的站点），不是红 | C2 的红是 `cmd/wisp/config_receipt_255_test.go:299`，且 HEAD 上已由 `50b540fad` 修绿 |
| `internal/risk/pathresolver.go:36`（C6） | 拒绝**确实发生**，只是 POSIX 上哨兵名不同（`ErrReparseDenied` vs `winsec.ErrUnresolvedPath`） | **冻结包** ⇒ 本腿 STOP；红用例本身也在 `cmd/wisp/secret_dataroot_119b_test.go` |
| `internal/winsec`（C6 尾注） | 票 94 已 done 的 POSIX 半边 | 无红用例 |
| `internal/memory/l1_warm_roster_29r2_test.go`、`internal/tools/spawn_roster_census_229_test.go`（文首） | 普查期间**别人**新增的两枚测试，用来证明快照与 HEAD 的 `cmd/wisp` 逐字节相同 | 与本腿无关 |

## 2. 实测：HEAD `3ff1b1395` 上 `./internal/...` 的红（逐枚）

仪器：本机 Windows `go test ./internal/... -count=1 -timeout 900s`，
`PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`，落 `/tmp/internal-all.log`。
读数：36 个包中 **29 ok / 4 FAIL / 3 无测试文件**（`internal/agent/scheduler`、`internal/audio/aec`、
`internal/streamkey`）；`SKIP` 行 **0**；顶层红 **5 枚**。

| # | 红 | 包 | 第一条红（逐字） | 根因 | 归谁 | 本腿 |
|---|---|---|---|---|---|---|
| 1 | `Test31DisarmStopsTheWait` | `internal/agent` | `reminder_timer_test.go:370: 取消后仍发了 1 条通知` | `run()` 的 `select`：Disarm 关掉 `stop` 之后、协程仍未进 select 时 Advance 又放了 `ready` ⇒ 两个 case 同时就绪，Go 随机取 ⇒ 一半的跑法照样 `fire`，而假行还是 armed ⇒ 通知发出 | 票 31 **已 done**；账册 A1477 逐字：「**未在任何账册/票面出现**，疑计时敏感/待核，⛔ 未证为回归」⇒ **孤儿** | **修**（§3.2） |
| 2 | `TestExternalSamplerTracksSubjectGrowth` | `internal/proc` | `externalsampler_windows_test.go:86: first ReadTree: resource: proc: external system process snapshot: NtQuerySystemInformation: buffer never sufficient (last 1143584 bytes)` | 重试规则照抄 `RequiredLength`；系统进程/线程数在两发之间还在长 ⇒ 每一发都比自己要的尺寸晚一步 ⇒ 8 发（旧 6 发）全 mismatch | 票 66 **已 done**，票面残口①逐字：「`NtQuerySystemInformation: buffer never sufficient`（200 次读 1 次）**未修**，归票 71 之后的仪器批次」——那批从没来 ⇒ **孤儿** | **修**（§3.1） |
| 3 | `TestComposerContractTypesMatchFrontend` | `internal/panel` | `composer_test.go:78: Go ComposerState emits [files] that interface ComposerState does not declare` | Go 侧发 `files`，`frontend/**` 的 `interface ComposerState` 没声明 | 前端（114-v2 已列「live 5 红既有集」） | **STOP**：要改 `frontend/**` |
| 4 | `TestC21DesignTokensFourWayAgree` | `internal/panel` | `tokens_fourway_test.go:726: docs/evidence/s1/c21-native-tokens.md:58: the table names Palette.GlassRing, but internal/ball/tokens.go holds no such colour literal in DarkPalette()` | 四向对账：证据表 78 色行 vs `internal/ball/tokens.go` 37/35 vs `frontend/scripts/gen-tokens.mjs` 43 | 设计/前端（D29 视觉面）＋证据表 | **STOP**：测试的四条腿有两条读 `frontend/**`，缺的两枚色值属设计面 |
| 5 | `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` | `internal/risk` | `pathresolver_rewrite_account_test.go:200: ..\..\cmd\wisp\panel_pump.go:283 reads Result.Canonical without reading the rewrite account (Actable/Rewritten)` | 新消费者读 `Canonical` 不看 rewrite account（票 102 的洞） | A380／票 189+190 | **STOP**：测试住在**冻结**的 `internal/risk`，且四个违规站点里两个在 `cmd/wisp/panel_pump.go`（本腿禁写） |

⇒ 可修 = 第 1、2 枚（都是"孤儿"：主人票已 done、残口无人接、账册点名"待核"），本腿全修。

## 3. 两枚修了什么（RED-before / GREEN-after，断言零放宽）

### 3.1 `internal/proc` —— 票 66 残口①：快照重试永远差一步

**机制**（`NtQuerySystemInformation(SystemProcessInformation)` 的契约）：`STATUS_INFO_LENGTH_MISMATCH`
回填的 `RequiredLength` 是**那一个瞬间**需要的字节数。旧规则在 `retLen > size` 时直接 `size = retLen`：
机器上的进程/线程数在两发之间继续长，下一发就又差一点 ⇒ **永远差一步**，6 发用尽后报
`buffer never sufficient`。这就是票面上「200 次读 1 次」的来源：只有系统正在长的时候才发作，
所以单跑 20 发全绿（`ok … 7.679s`），整包并发时（本仓自己的并行测试腿就是那种负载）才红。
票 66 的读数 `last 1143584 bytes` 与本轮复现的 `last 1245184 bytes` 是**同一条**路径。

**改法**（`internal/proc/treemetrics_windows.go`）：

- `sysProcNextSize(size, required)`（新，`:242`）：**至少翻倍**（psutil 的 Windows `NtQuerySystemInformation`
  循环就是翻倍），若 required 已超过翻倍值则给 `required + required/8` 的**余量**，再夹到上限
  `sysProcMaxSize = 64 MiB`；到顶返回 `ok=false` ⇒ 调用方 fail-closed，绝不无限涨。
- 尝试数 `sysProcAttempts = 8`（旧 6）：翻倍路径下 8 发足够从 1 MiB 摸到 64 MiB 顶；
  **真正的拒绝原因是上限，不是次数**。
- 把内核调用抽成 `sysProcQuery` 注入（`SystemProcessSnapshot` 只负责把 syscall 包成 query），
  于是「系统在长」这件事可以**确定性地**建模（`systemProcessSnapshot(query)`，`:282`）——
  旧规则下这条路不可能被稳定复现（它只随负载发作），这是它躲过 8 个月的原因。

**新测试** `internal/proc/snapshot_growth_180_test.go`（3 枚）：

| 用例 | 钉住什么 |
|---|---|
| `TestSystemProcessSnapshotKeepsUpWithAGrowingSystem`（`:47`） | 每发 need 长 24 KiB 的模型下，**必须取得到快照**（≤3 发），且解出来的 pid/私有工作集/线程/句柄/CPU 逐字对；**测试内置反向对照**：同一模型下旧规则（照抄 required）必须**不收敛**，否则这个用例什么也证明不了 |
| `TestSystemProcessSnapshotFailsClosedAtTheCap`（`:98`） | required 超过上限 ⇒ **报错**（≤2 发），不无限分配 |
| `TestSysProcNextSizeNeverShrinksOrExceedsTheCap`（`:115`） | 每一步都变大、绝不越过上限、required 领先翻倍时**必须留余量**（旧规则恰好落在 required 上） |

**RED-before（把 `sysProcNextSize` 临时换回旧规则，overlay，逐字）**：

```
=== RUN   TestSystemProcessSnapshotKeepsUpWithAGrowingSystem
    snapshot_growth_180_test.go:78: a snapshot whose system is still growing must still be read:
        NtQuerySystemInformation: buffer never sufficient (last 1245184 bytes)
--- FAIL: TestSystemProcessSnapshotKeepsUpWithAGrowingSystem (0.00s)
=== RUN   TestSysProcNextSizeNeverShrinksOrExceedsTheCap
    snapshot_growth_180_test.go:122: next size = 1048577, want >= one doubling (2097152)
    snapshot_growth_180_test.go:132: next size = 8388608, want > reported need 8388608 (headroom for growth between two calls)
--- FAIL: TestSysProcNextSizeNeverShrinksOrExceedsTheCap (0.00s)
FAIL	github.com/CarlosShao/wisp/internal/proc	0.057s
```

注意第 78 行那句：**和线上那条红逐字同形**（同一个 `buffer never sufficient`），只是尺寸不同。

**GREEN-after**：

```
=== RUN   TestSystemProcessSnapshotKeepsUpWithAGrowingSystem
--- PASS: TestSystemProcessSnapshotKeepsUpWithAGrowingSystem (0.00s)
=== RUN   TestSystemProcessSnapshotFailsClosedAtTheCap
--- PASS: TestSystemProcessSnapshotFailsClosedAtTheCap (0.02s)
=== RUN   TestSysProcNextSizeNeverShrinksOrExceedsTheCap
--- PASS: TestSysProcNextSizeNeverShrinksOrExceedsTheCap (0.00s)
=== RUN   TestExternalSamplerTracksSubjectGrowth
--- PASS: TestExternalSamplerTracksSubjectGrowth (0.31s)
ok  	github.com/CarlosShao/wisp/internal/proc	0.464s
```

### 3.2 `internal/agent` —— `Test31DisarmStopsTheWait`：取消输给"已放行的定时器"

**机制**（不是计时敏感）：`Arm` 在调用方协程上先 `t.clock.After(delta)` 拿到 `ready` 再 spawn 等待协程。
`Disarm` 关 `stop`；`clock.Advance` 给 `ready` 塞值。若协程**还没进 `select`**，醒来时 `stop` 与 `ready`
**同时就绪**，Go 的 `select` 在多个就绪 case 之间**随机**挑一个 ⇒ 约一半的跑法走 `fire`：
假行的状态仍是 `armed` ⇒ `MarkReminderFired` 回 `applied=true` ⇒ **通知照发**。
原注释说「这个竞态由 fire 里的 fired-stamp 守卫兜住」——那只兜**反方向**（取消在 fire 之后到达、
行已被标 fired 的那种），兜不了「取消先到、放行后到」。

**改法**（`internal/agent/reminder_timer.go` `run()` `:323-338`）：`<-ready` 分支之后**再对 `stop` 做一次非阻塞检查**，
命中就按 DISARMED 退出（同一个 emit 文案）。这不是"第二权威"：`stop` 是**这一个等待**的权威
（取消先到就绝不投递），行上的 fired-stamp 仍是**跨进程/跨路径**的第二道保险（取消在 fire 飞行途中到达）。
文件头 `:45-55` 的说明同步改掉了旧的「故意不复查」表述。

**新测试** `internal/agent/reminder_disarm_race_180_test.go`：`Test180DisarmBeatsAnAlreadyReleasedWait`（`:25`）
把同一序列跑 **48 轮**（Arm → Disarm → Advance → 等协程退出），**任何一轮**发出通知即判红，
并同时断言行仍是 `armed`、审计行 0 枚（与票 31 原用例同口径）。单发 ~50% 的红 ⇒ 48 轮漏网概率 ~1e-15，
所以这条用例的判据是**确定性**的，不再是丢硬币。

**RED-before（overlay：把复查那段删掉，逐字）**：

```
=== RUN   Test180DisarmBeatsAnAlreadyReleasedWait
    reminder_disarm_race_180_test.go:57: round 1: 取消后仍发了 1 条通知
--- FAIL: Test180DisarmBeatsAnAlreadyReleasedWait (0.01s)
=== RUN   Test180DisarmBeatsAnAlreadyReleasedWait
    reminder_disarm_race_180_test.go:57: round 0: 取消后仍发了 1 条通知
--- FAIL: Test180DisarmBeatsAnAlreadyReleasedWait (0.00s)
FAIL	github.com/CarlosShao/wisp/internal/agent	0.055s
```

（两发都在第 0/1 轮就红 —— 说明这条红一直是**高概率**的，不是"偶发"。）

**GREEN-after**（`-count=3`，两枚用例交替跑）：

```
=== RUN   Test180DisarmBeatsAnAlreadyReleasedWait
--- PASS: Test180DisarmBeatsAnAlreadyReleasedWait (0.12s)
=== RUN   Test31DisarmStopsTheWait
--- PASS: Test31DisarmStopsTheWait (0.05s)
…（×3，6/6 PASS）
ok  	github.com/CarlosShao/wisp/internal/agent	…
```

**原用例的旧读数**（改前，`-count=4 -v`，同一条命令）：

```
=== RUN   Test31DisarmStopsTheWait
    reminder_timer_test.go:370: 取消后仍发了 1 条通知
--- FAIL: Test31DisarmStopsTheWait (0.05s)          ← ×3
--- PASS: Test31DisarmStopsTheWait (0.05s)          ← ×1
```

⇒ **3 红 1 绿**：这就是 A1477 那位前辈"疑计时敏感"的真相，也是"整包有时红有时绿"的原因。

## 4. 三枚 STOP（写明为什么不在本腿写面）

- **#3 `TestComposerContractTypesMatchFrontend`**：红是 `frontend/**` 的 `interface ComposerState` 少了 Go 侧
  `files` 字段。改法只有两条，都在本腿写面外（改前端声明 / 改 Go 侧不发 `files` 又等于撤已交付的能力）。
- **#4 `TestC21DesignTokensFourWayAgree`**：测试四条腿里的两条读 `frontend/scripts/gen-tokens.mjs` 与
  `frontend/src/styles/tokens.generated.css`；缺的三枚色值（`Palette.GlassRing/GlassHi/GlassHiSoft`）是
  **设计面**（D29 视觉/token 归外部设计 agent）。本腿不造色值（宁缺毋造）。
- **#5 `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`**：红用例住在 **`internal/risk`（冻结）**；
  它的四个违规站点里两个在 **`cmd/wisp/panel_pump.go:283/:291`（本腿禁写）**。按派单红线 STOP。
  ⇒ 这枚仍是孤儿，建议另开票（票 189+190 已 done，A380 只归了 `TestC26` 这个名字，没有接残口）。

## 5. 门禁读数（本腿自验，逐条可重跑）

```bash
cd "/d/work/workspace/projects plans/Wisp"
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"
go build ./internal/...                                  # rc=0
go test ./internal/agent ./internal/proc -count=1 -timeout 240s   # ok / ok
go test ./internal/... -count=1 -timeout 900s            # 只剩 §4 那 3 枚（STOP）
/d/work/base/gopath/bin/gofumpt.exe -l internal/agent/reminder_timer.go \
  internal/agent/reminder_disarm_race_180_test.go \
  internal/proc/treemetrics_windows.go internal/proc/snapshot_growth_180_test.go   # 输出空
sh scripts/d22scan.sh                                    # 见下
```

| 门 | 读数 |
|---|---|
| `go build ./internal/...` | **rc=0** |
| `go test ./internal/agent ./internal/proc -count=1 -timeout 240s` | **ok 7.401s / ok 1.586s**（GATE_TEST_RC=0） |
| 整包同一发里这两包 | agent **ok 28.453s**（改前同一条命令 FAIL 36.571s）、proc **ok 5.897s**（改前 FAIL 13.992s） |
| `go test ./internal/... -count=1 -timeout 900s` | **31 ok / 2 FAIL / 3 无测试文件**（36 包）；剩 `internal/panel`(2 枚，STOP)、`internal/risk`(1 枚，STOP)；**`internal/agent`、`internal/proc` 转绿** |
| `gofumpt -l`（四个改动文件） | **空** |
| `gofumpt -l internal/`（全包） | 只剩 **1 枚不是我改的**：`internal/agent/approval/pending_read.go`（`git ls-files --eol` = `i/lf w/crlf`，HEAD 里 2026-10-04 `d03d166fe` 就在，工作树被谁写成了 CRLF；`git status` 对它为空 ⇒ **会话前既有、未被我碰**，也不在本腿写面）。 |
| `sh scripts/d22scan.sh` | **rc=0**：`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=357, bans #1-5 cmd/=57, ban #6 frontend/=116, ban #7 internal/tools/=48, ban #8 design/=39, ban #8 frontend/=116, ban #8 internal/=882, ban #8 cmd/=206`（正控 + 真扫两段都过） |

## 6. 本件没有建立的事（宁缺毋造）

- **没有重跑 POSIX/ubuntu 腿**：本腿只在本机 Windows 上取数；C1 那 76 枚、C6 那 1 枚的处置**不因本件改变**。
- **没有把 `go test ./cmd/wisp` 当门**：那是派单的禁写面，也是 98 号件的主读数；本件只引用、不触碰。
- **`internal/proc` 的修只证明了"规则收敛" + "整包不再红"**：真机上"每 200 发红 1 次"这种长尾，
  本件没有跑 200 发 × N 组去统计（时间/收益不成比）。新用例把**机制**变成确定性判据，这是本件的证法。
- **没有动 `internal/ball`**：C21 缺的三枚色值是设计面取值，不是"补个常量"。
- **没有裁 `Test31DisarmStopsTheWait` 原用例本身**：它的断言一字未改，改的是 `run()` 的行为。
  若 owner 认为"取消先到也可能投递"是允许的，那要改的是**用例与票 31 的措辞**（改契约文字 ⇒ 摆机主），
  不是本腿能自决的。
- **没有修 `internal/agent/approval/pending_read.go` 的 CRLF**：那是工作树里既有的 `w/crlf`（HEAD blob 是 LF），
  `gofumpt -l internal/` 因此不是空。它不在本腿写面；照实记在这里，不顺手改别人的文件。

## 7. 复跑口令（逐字可重跑）

```bash
cd "/d/work/workspace/projects plans/Wisp"
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"

# 1) 两枚新用例（绿）
go test ./internal/proc  -run 'TestSystemProcessSnapshotKeepsUpWithAGrowingSystem|TestSystemProcessSnapshotFailsClosedAtTheCap|TestSysProcNextSizeNeverShrinksOrExceedsTheCap' -count=1 -v
go test ./internal/agent -run 'Test180DisarmBeatsAnAlreadyReleasedWait|^Test31DisarmStopsTheWait$' -count=3 -v

# 2) 反向（红）：见 §3 两段 overlay。proc 的改法是 sysProcNextSize 里
#    `next := size*2 … required+required/8` 换回 `next := required`；agent 的改法是删掉 `case <-ready:`
#    之后那段对 w.stop 的非阻塞复查。改完各自跑一次上面第一条命令即红。

# 3) 整包对照
go test ./internal/... -count=1 -timeout 900s
```

## 8. 时刻账

- `13:00` 读派单 §98 号件；`13:01` 起 `go build ./internal/...` rc=0、整包 `./internal/...` 首跑（红 5 枚）。
- `13:04` 两枚候选单跑（都 PASS）⇒ 判为负载相关；`13:05` `-count=4` 把 agent 那枚钉成 3 红 1 绿。
- `13:06–13:10` `internal/proc` 改（重试规则 + 注入缝 + 3 枚新用例）；overlay 反向红（逐字见 §3.1）。
- `13:11–13:15` `internal/agent` 改（复查 `stop`）＋ 48 轮新用例；overlay 反向红（逐字见 §3.2）。
- `13:16–13:2x` 整包复跑（agent/proc 转绿，剩 3 枚 STOP）、gofumpt 空、d22scan、本件。
