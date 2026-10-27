# 票 233 终裁腿 `233-v1` — 六格 6/6 能翻；票 223 `AC#2` 的「谈」＝**能翻**；判**能带具名残余结案**

- 身份＝非实现者终裁腿。被裁对象＝`233-r1`（提交 `bd2617652`，8 件 `540+/105-`，件 `docs/evidence/s1/233r1-impl.md` 79 行）；收件账 `A1290`。起手锚＝`adaa18e22`（HEAD，`bd2617652` 之上仅编排者台账一枚）；被裁树与起手树之间 `internal/config`/`cmd/wisp` 零中间提交（`git log bd2617652..HEAD` 无 .go）。
- 我的全部读数在**仓外**：`D:\tmp\wisp\20261027\233v1\{exp,base,exp-mut,stage,stage-base,stage-mut,logs}`（`git archive` 两棵树：`exp`＝`bd2617652`、`base`＝其父 `1bb1ae961`；导出树 8 枚关键文件 `git hash-object` 对 `git rev-parse bd2617652:<f>` **逐字 SAME**）。主检出全程只读（收工 `git status --porcelain -- internal cmd`＝空）。`cmd/wisp` 测试一律 staging（二进制与三枚 sherpa DLL 同目录，零 0xc0000135）。

## 0. 现跑读数（六格现状，尺与对象逐条具名）

| # | 读数 | 数 | 对象/尺 |
|---|---|---|---|
| 1 | 提交真在 | `bd2617652cc11d7…` 2026-10-27 21:42 +08 | `git show --stat`＝8 件恰 540+/105- |
| 2 | 三常量零赋值点已闭合 | `TierHot` 非测试命中 20＝**17 赋值行**（tiers.go）+3（schema.go 注释+声明）/`TierReload` 23＝**20+3**/`TierRestart` 6＝**3+3**/新增 `TierLocked` 8＝**4+4** | `git grep -n` 树 `bd2617652`，分层数：赋值行 vs 注释/声明行（ruler 清单第 1 例） |
| 3 | `TierRegistry` 44 行、值全为常量 | 44＝17 hot+20 reload+3 restart+4 locked（与总数对上） | `Select-String 'Tier(Hot|Reload|Restart|Locked),'` |
| 4 | reload 档消费者 | `OnReload =` 非测试恰 2：`cmd/wisp/config_reload.go:138`（新，**wisp run 首枚**）+`cmd/balldebug/main.go:244`（旁支） | 同尺；`cmd/wisp` 单包＝1 |
| 5 | 热档活引用 | `run.go:1261 cfg := rt.cfg`＋`:1263 cfg = rt.mgr.Config()`（`rt.mgr != nil` 分支） | 现读 `:1252-1264` |
| 6 | `restartTierKeys` 未删 | `config_reload.go:364`，3 键；其常驻钉 `TestTicket255RestartTierKeysAreBackedByATest` 在（要求每键 `TierRegistry` 登记 restart） | 现读 |
| 7 | 禁区零字节 | `git diff --name-only bd2617652^ bd2617652 -- internal/risk|speech|session|audio|tools|secret|memory testdata frontend design scripts .github docs/PLAN.md docs/specs`＝**空** | 树对树 |

## 1. 票 233 六格逐条判（★=能翻的凭据）

| 格 | 判 | 凭据（我自跑/自核） |
|---|---|---|
| **AC#1 现状表** | ★**能翻** | 233r1 §0 修正表（3 处派单过时逐条具名）+§1 44 行逐键表（键/档常量/D36 依据 `PLAN.md:2738-2757` 逐行对上）；「今天实际按哪档走」由四枚完备性尺（`TestEverySectionHasATierInRegistry` 族）＋255 名册常驻钉背书；三常量命中读数在 §0 表内 |
| **AC#2 落点甲** | ★**能翻** | schema 侧真相源成立（registry 值＝`Tier` 常量，typo 编不过）；`restartTierKeys` 与 registry「两边不许漂」由 `restart_tier_keys_255r2_test.go` 钉住（255-r2 既有、零触碰）；⚠ 新导出名 `TierLocked`（schema.go:48）＝禁区「先落 A##」那一格——**A1290 已具名收下该枚**（编排者收件笔），批准面成立、非越权 |
| **AC#3 reload 接/说** | ★**能翻** | 「接」支：`:138` 真生产赋值点（非测试、非 internal/ball、装配根 `startConfigReload`）；三禁法逐条核：①非静默（stdout+audit 双句）②未塞 restartTierKeys 冒充 ③manager.go 触发条件零触碰（diff 无 internal/config/manager.go） |
| **AC#4 常驻判据钉 stdout** | ★**能翻** | `config_tiers_233_test.go` 三用例：stdout 句（`reloadTierSentence233`）与 audit 行（`auditReloadTier233`）**两常量分流断言**，窗口＝plant 前 mark 的增量（`awaitStdoutSince232`/`awaitAuditSince232`）＝223-v2 拼接串坑未复发；行为读点＝`rt.loopOpt`（execute 自己发布的 Options）非探针字段 |
| **AC#5 突变自证** | ★**能翻** | M1/M2 我自己复刻（§3）：两发全红、红名逐字点名用例自己；还原后原样 3/3 绿 |
| **AC#6 门禁与终态** | ★**能翻** | 包级读数我复跑全绿（§5/§6）；字面明写「整包终态与逐名比红名册归编排者，腿只交包级读数」⇒ 腿侧字面已满足；历史在册红（internal/ball 1＋panel 4）不在本腿射程 |

## 2. 三档读数（我自己造的；stage 二进制，非转抄）

| 档 | 用例 | 我的读数 |
|---|---|---|
| **热（真热）** | `Test233HotTierAgentReachesTheNextTask` | **3/3 PASS**（2.10/1.89/1.90s）：手改 `[agent] per_tool_timeout_ms=2500` → tick `hot=[agent]` → **同进程下一个任务** `rt.loopOpt.Config.PerToolTimeout==2.5s` 且 **boot snapshot 仍 0**——「改了→不重启→行为变了」的字面形状 |
| **reload（逐字告知）** | `Test233ReloadTierEventReachesAProductionConsumer` | **3/3 PASS**（2.05/1.94/1.88s）：plant `[voice.asr] model` → `state=applied … reload=[voice]` → audit 逐字 needle `config: RELOAD-TIER state=received`/`sections=[voice]`/`effect=this-host-reloads-nothing` ＋ stdout 逐字 needle「重载档/[voice]/已生效到本进程内存/不会触发模型重载」；对照 `Test233ReloadTierReceiptOnlyForReloadTier` PASS（hot 改动不触发 reload 句） |
| **重启（未生效告知）** | `TestTicket223RestartTierSaysItWillNotApply`（232/272 窗口尺零改动） | **3/3 PASS**（1.98/1.95/2.03s），`state=restart-pending`+stdout 五 needle |

## 3. 两发反向判据（我自己造的；全部在仓外 `exp-mut`，主检出零突变）

- **M1＝摘 `:138` OnReload 赋值**（注释包住）→ `Test233ReloadTier…` **FAIL 42.09s**，红句逐字 `config_tiers_233_test.go:90: the audit trail never carried "config: RELOAD-TIER state=received" AFTER the plant within 40s`；窗口 dump **真 `[audit]` 行分层**（ruler 第 1 例：第一遍我把红句引用的 needle 误计入 RELOAD-TIER 命中，重数后）＝`state=applied … reload=[voice]` 在、`RELOAD-TIER` 真行 **0**（事件进了无人线）；对照 hot 档用例 **PASS 1.87s**（不误伤）。
- **M2＝execute 回 snapshot**（`:1262 if false { cfg = rt.mgr.Config() }`）→ `Test233HotTier…` **FAIL 1.93s**，红句逐字 `:193: … built with PerToolTimeout=0s, want 2.5s …`；对照 reload 用例 **PASS 1.98s**。
- **还原核（时机＝M2 之后、最后一发变异之后）**：两文件 `git hash-object`（SHA-1，同算法同位数）对 `bd2617652` blob **SAME=True**（config_reload.go `63beea79…`）；`MUT233v1` 标记 0 命中；重建后三枚 233 用例 **3/3 PASS**。

## 4. `TestTicket255Roster` 我复核的（A/B ＋ 零放宽核）

- **A/B**：base 树 **FAIL 0.61s**（红因逐字：`config_receipt_255_test.go:295` SECTION 13/14 `session` scan 命中 `resident_voice_windows.go` 名册为 []、SECTION 14/14 `voice` 名册缺第四件——「something moved … has not been re-adjudicated」）；mine 树 **PASS 0.62s**。＝**HEAD 基线既存红**（wire-r1 新读者未回名册），233-r1 重裁后闭合。
- **「只增读者零放宽」逐行核**（diff 通读）：`session` no-reader→other-process（新 cite `resident_voice_windows.go:1329`，真实读点）＋`sectionReadSites["Session"]` nil→新增＋`["Voice"]` 加第四件＝**追认事实**；`agent` snapshot-only→consumed（cite `run.go:1263`）＝**本腿产码真变的同步**；阳性样本 `[session]`→`[ball]`（ball 扫描零读者现量仍真）；**反向半句收紧**：claimed 不得进 honest bucket 从 `[llm]` 扩到 `[llm]/[agent]`＝净更严。**零断言放宽坐实**。

## 5. 零回归（我自己复跑）

| 尺 | base | mine |
|---|---|---|
| `internal/config -count=3` | rc=0（4.278s） | **rc=0**（4.447s） |
| `cmd/wisp` 定向 14 枚 ×3（233×3/255×7/223×2/227×1） | 同尺 ×1 恰红 1＝roster（既存） | **rc=0，42 PASS/0 FAIL** |
| `internal/risk -count=1`（零触碰对照） | 红恰 1＝`TestC26RewriteAccountIsConsumedAtEverySecurityLeg` | **同枚同红**（A380 存量） |
| `go vet ./cmd/wisp ./internal/config/` | — | **rc=0** |

**fs-card 竞态（第 21 处定式）我自己的 A/B**：mine 20 发 **0 红**；base 30 发 **2 红**（红行 `:436/:439`「console 没渲染卡」逐字，log 16b-RED-4/14）＝`awaitCard` 认账与 `Prompt` 渲染之间的**仪器窗口、HEAD 既存**，与本腿无因果（形状与 233r1 的 3/30 vs 2/60 同）。

## 6. `d22scan` 导出树＋正控

- 独立模块 build（`exp/tools/d22scan`，rc=0）→ `./d22scan.exe -root exp` **rc=0 clean**，逐 scope 计数与 233r1 逐字同（bans#1-5 internal=264 cmd=45、#6 frontend=116、#7 tools=25、#8 design=14 frontend=116 internal=684 cmd=160）。
- **正控**：真种 `≤`（U+2264，字节级 `e289a4` 现量）进 `exp/internal/config/zz_d22_positive_control_233v1.go` 字符串 ⇒ **rc=1 且逐字点名** `zz_d22_positive_control_233v1.go:4: [emoji] ban #8 glyph in scope internal/ is banned (D23): non-comment text, string literals included`；种后即删、复扫 **rc=0 clean**。尺是活的。

## 7. ★票 223 `AC#2` 的「谈」（233 落地之后再谈——现在谈）

- **票面逐字**（`:35`）：「立即生效的那一档要能观测到"改了→不重启→行为变了"；重启后生效的那一档要能观测到"改了→不重启→**明确告知未生效**"。⛔ 不许用"静默不生效"充当第二档。」09-29 钉的翻勾条件＝「**233 落地之后再谈**」（`:84`）。
- **`223-v1` 当时判 ⛔ 翻的三条缺口，逐条对今天**：
  1. **reload 生产消费者＋操作员读数** → **已补**：`:138` 生产赋值点恰 1（cmd/wisp 单包）；操作员 stdout 句＋audit 句双面（§2）；M1 摘除即红 42.09s（有牙）。
  2. **hot 的"行为"（当时只是读面答新值）** → **已补成动作面**：`execute()` 每任务现读活 Manager，同进程**下一个任务**的 loop 装配真带 `PerToolTimeout=2.5s`（M2 回 snapshot 即红 1.93s）——不再是「live reader answers differently」而是**下一个任务的行为输入真变**。
  3. **"233 落地之后再谈"未谈** → **本次即谈**：233 已落地（`bd2617652`，六格全绿），本腿就是那一谈。
- **判：★★★`AC#2` 现在能翻**。字面两档（立即/重启后）各有观测：hot＝行为真变、restart＝明确告知未生效；reload 第三档（233 归口）另有「明确告知生效边界」的读数且非静默。**翻它要动票 223 票面的哪一格**＝只动 **AC#2 那一格 checkbox**（`- [ ]`→`- [x]`）＋按本仓翻框笔法在格后追加一行凭据（判语出处 `docs/evidence/s1/233-v1-verdict.md` ＋ 翻框提交号），**⛔ 不动 AC#7/其余格/Status/票名**。⚠ 翻框权＝编排者，我不翻。

## 8. 残余清单（逐条具名）＋能不能带残余结案

1. `TestTicket223HandEditedFsLooseningCostsAnL2Card` 仪器窗口竞态：**HEAD 既存**（我 base 2/30 vs mine 0/20、同红行），修法＝`:434` 单读改 `awaitStdout` 形（351 预案），归票 351/下一张票。
2. reload 档「接」的深度：本宿主 `wisp run` 无语音链可重载，`reportReloadTier` 是告知形状（票面 AC#3 的「接」字面已闭合）；**语音模型真换装**归常驻进程腿（归口在票面与 commit 正文已具名）。
3. 新导出名 `TierLocked` 的批准面＝A1290 具名收件（已落账，无悬空）。
4. macOS/Linux 门禁未跑：改动为纯 Go 语义层、无平台分支，CI 会跑。
5. 整包 `-count=3` 与逐名比红名册：票面 AC#6 字面归编排者，本腿只交包级（§5）。

- **判：★★★能带具名残余结案**。对照先例：`15-done`（3 格未勾＝等真机类）、`347-done`（`AC#0` 机主追认格未勾）、`26-done`（3 格未勾＝真机/barge-in 归 59/C25 归 19）——三例共同形状＝**未勾项全部「不在本票射程＋具名归口」或「等外部事件」**，且本票字面判据无一格未闭合（§1 六格全绿）。残余 5 条全部具名＋归口明确 ⇒ 同形状成立。翻勾/`-done`/台账 A## ＝编排者权。

## 9. 提交号＋逐笔越界

- 本腿提交恰 1 枚（见下），写面＝`docs/evidence/s1/233-v1-verdict.md`（≤180 行）一枚；`.go`＝**0**；越界＝**0**。正控＝这把尺在 `bd2617652` 上数得 7 枚 `.go`（尺看得见）。`add`/`commit` 两处显式 pathspec；`-F` 真临时文件（`D:\tmp\wisp\20261027\233v1\commitmsg-233v1.txt`，仓外）写在 `--` 之前；⛔ push。
- 日志：`D:\tmp\wisp\20261027\233v1\logs\01..22`（只建不删）。导出树 `exp-mut` 与两枚 stage 目录为本腿工作残留（仓外），不清理（临时件只建不删）。
