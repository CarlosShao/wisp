# 票 266 实现腿 `266-r2` —— 外牙 `scripts/slo-nail-census.sh` 落地（形 ①-A 的 `scripts/` 半场）

- 腿身份＝**实现者**。⛔ 未翻 AC 框、⛔ 未动票面 `Status:`/票名/票尾 bullet（票面不在本腿写面）。
- 锚：起手现读 HEAD＝`fae8ba704`（派单写 `9e84ae42e`；现读多一枚 `fae8ba704`＝wave14 收件腿，`--name-only` 与本腿写面零相交）；
  交件 commit＝`6589a0625`（`git show --name-only` ＝ **仅** `scripts/slo-nail-census.sh` 一枚）；读数窗口 **2026-11-05 18:17–18:36 +08**。
- 写面（派单）＝`scripts/**` ＋ 本件。改动两枚：新件 `scripts/slo-nail-census.sh`（569 行）＋ 本件。
  **零 Go、零前端** ⇒ 派单的 STOP 条件未触发。
- 环境：`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`；`unset GOTMPDIR TMPDIR TEMP TMP`；⛔ 全程未用 `CGO_ENABLED=0`。
- **窗口内 HEAD 漂移（共享工作树，具名）**：`b5ac0a181`（收 wave16）与 `6a46140ea`（230-ac4 part1）在本腿交件前后落地；
  两笔 `git show --name-only` 与本腿写面零相交（config/tools 侧）。`scripts/slo-check.ps1` 自 `cb8dceccc`（263-r1）以来**一字节未动**
  （`git log --oneline -1 -- scripts/slo-check.ps1` 与 `grep -n` 行号在 18:35 现读仍为 `:214/:233/:291/:292`）⇒ 本件引用的行号在交件时仍准。
- **结论先行**：AC#1/AC#2/AC#3/AC#4 已交（读数见下）；**形 ①-A 的另一半（CI 步）未交**，因为 `.github/**` 不在本腿写面
  ⇒ 逐字补丁在 §6。⚠ **所以今天这枚外牙还没有任何 CI 调用点**：它是"能手跑、且可直接接进 lint 作业的一枚牙"，
  ⛔ 不是"已上牙"。这一句请编排者具名收件，别当已闭合。

## §1 AC#1 —— 形状落在被检文件之外（`file:line` 具名）＋ 绕过回答

**形状**：新尺 `scripts/slo-nail-census.sh`（POSIX sh，`set -eu`）。**声明**（= D32 门该有 **2** 枚自锁钉、两枚各叫什么、
各带哪句 ok 句）写在它**自己**的 `:110-112`（`declared_nail_count=2` / `nail_names=…` / `nail_ok=…`）；
被检件 `scripts/slo-check.ps1` 里**一字不提**这些声明 ⇒ "删钉"与"删声明"不在同一枚文件里。

**它判什么（指向被检件的真身行号＝今天 18:24 现读，非照抄票面 §现量 4）**

| # | 判据 | 被检件真身（2026-11-05 18:24 +08 现读） |
|---|---|---|
| 1 | 声明自洽：`declared_nail_count` ＝ 两张表的行数 | 本件 `:110-112`（改一不改二 ⇒ 红 `declaration-drift`） |
| 2 | 每枚钉的定义**恰一条**、顶格 `function <名> {` | `scripts/slo-check.ps1:214`（Test-ExitCodeInstrument）／`:233`（Test-ScriptShapeNail） |
| 3 | 每枚钉**自己的 ok 句**落在**自己函数体内**（定义行 → 下一个顶层 `function` 之间） | ok 句 `:230`（`instrument self-check ok`）／`:258`（`shape nail ok`） |
| 4 | 每枚钉的**调用点恰一条**、顶格、未注释、独立行 | `:291`（Test-ScriptShapeNail）／`:292`（Test-ExitCodeInstrument） |
| 5 | 位置：调用点晚于自身定义、且早于首枚采样调用 `-FilePath $WispExe` | 首枚采样＝`:533`；"两枚钉在任何测量之前跑"（`:286-290` 注释逐字）从此有牙 |

**凭什么发现不了"删检查也删声明"——三层绕过，逐层具名（⛔ 不许把任一层说成关上的门）**

1. **今天能发现的**：把 `scripts/slo-check.ps1` 里两枚调用删掉（票 263 的 `delete-nails` 变异）⇒ **必须响**（§2 发①逐字）。
   机制＝声明住在**另一枚文件**里，删被检件的钉不动声明，所以尺还在喊"该有 2 枚"。
2. **仍然绕得过、且体积最小的一发**：一枚 commit **同时删掉** `scripts/slo-nail-census.sh` 与调用它的 CI 步（§6 那一步）
   ⇒ 今天**没有任何东西会响**。唯一另一枚外牙 `scripts/slo-freshness.sh` 的 P3 判的是**最新 `slo-full-report` 产物年龄**
   （常量 `:137` `SLO_FULL_SAMPLE_MAX_AGE_DAYS:-3`、判定 `:550` `-gt`、默认 3 天），对"钉在不在"**零敏感**；
   且它的自动入口 `slo-fresh.yml:46-52` 只有 `schedule`+`workflow_dispatch`，cron 只在默认分支评估（dev 上 inert）。
   ⇒ 本件是**对一行编辑的牙**，不是关上的门。
3. **静态尺的固有盲区**：把两枚调用**移进一个没人调用的函数**里，本尺看不见（它钉"调用点在主流程里、位置对"，不证可达）。
   "真跑"（那两行 ok 句真被打出来）只能由形 ② 给；本尺反过来说明了这一点，且每次 PASS 都带一行 NOTE 复述，不藏。

**另一条诚实的边界**：本尺只读文本。目标里今天**没有** here-string（现跑 `grep -nE "@.( |)$" scripts/slo-check.ps1` ＝ **10 命中**，逐条都是 `@{`／`@(` 造型；here-string 收尾形 `grep -nE "^[[:space:]]*['\"]@" scripts/slo-check.ps1` ＝ **0 命中**），
但 reader 仍为 here-string 与块注释各留了一枚状态位；运行时拼出来的调用（`Get-Variable` 那族）照旧抓不到（票 263 已具名）。

## §2 AC#2 —— 正控两发（"删钉必响"逐字）＋ 还原证

**变异台**：`mktemp -d` 造台件，里面放**本尺的拷贝**＋**被检件的逐字节拷贝**（`scripts/{slo-nail-census.sh,slo-check.ps1}`）。
被检件路径由本尺自身位置推出（`scripts/d22scan.sh` 同款规矩：不依赖 cwd、不留第二份路径），⇒ 变异台**零测试缝**，
也**零工作树足迹**（真件一字节没改）。

**发①（删两枚调用 ⇒ 必响）—— rc=1，红句逐字（台件 `/tmp/tmp.lNlkWi6EUW/out-5.txt`）**

```
slo-nail-census: nail 1/2 Test-ScriptShapeNail - definition slo-check.ps1:233, ok sentence slo-check.ps1:258, call slo-check.ps1:0
slo-nail-census: FAIL slo-nail-census-nail-call-missing: the declared nail 'Test-ScriptShapeNail' has no column-0, uncommented, standalone call line in slo-check.ps1 (its definition sits at slo-check.ps1:233). Nothing runs it - this is the ticket 263 delete-nails mutant, and the reason this tooth exists (ticket 266 AC#1/AC#2).
slo-nail-census: nail 2/2 Test-ExitCodeInstrument - definition slo-check.ps1:214, ok sentence slo-check.ps1:230, call slo-check.ps1:0
slo-nail-census: FAIL slo-nail-census-nail-call-missing: the declared nail 'Test-ExitCodeInstrument' has no column-0, uncommented, standalone call line in slo-check.ps1 (its definition sits at slo-check.ps1:214). Nothing runs it - this is the ticket 263 delete-nails mutant, and the reason this tooth exists (ticket 266 AC#1/AC#2).
slo-nail-census: FAIL 2 declaration violation(s) - the D32 gate is not the shape this file declares (ticket 266). A green from a gate whose nails are gone is not a reading: ticket 263 measured an all-green screen with both call lines deleted.
```

**发②（还原 ⇒ 转绿）**：同一台件把被检件拷贝回原样 ⇒ **rc=0**，逐字
`slo-nail-census: PASS - 2/2 nails declared here are defined once, carry their own ok sentence, are called once at column 0, and run before the first sample at slo-check.ps1:531`。
（`:531` 而非 `:533` 是因为发①删掉的两行把行号前移了 2；还原后真件读数仍是 `:533`。）

**32 发控制全 PASS（本腿在三种 shell/awk 组合上各现跑一发：bash/sh、dash、dash+`gawk --posix`；末发 18:33:49–18:34:21，读数 `32 control(s) ran, 0 failed`）**，
逐发点名（关键几发）：

- 纯拷贝 → rc=0，且读数点名 `nail 1/2` / `nail 2/2`（**正控的"静"这一半**）；
- **删两条调用** → rc=1 ＋ `nail-call-missing`（发①）；**还原** → rc=0（发②）；
- **只删一条** → rc=1，红句**只骂欠的那枚**，另一枚**不误伤**；
- **行注释**包住两条 → rc=1；**块注释**包住两条（`<#`…`#>`）→ rc=1（这一发只有真做注释剥离才会红，是本尺"不是只数文本"的证据）；
- 删定义 → `nail-definition-missing`；删 ok 句 → `nail-ok-sentence-missing`；
- **缺件** → `exit 2` ＋ `slo-nail-census-cannot-decide` ＋ 逐字 `NOT a pass`；**空件** → 同；
- **CRLF 拷贝 → 仍绿**（`\r` 剥离是载荷：`git ls-files --eol scripts/slo-check.ps1` ＝ `i/lf w/lf attr/text eol=crlf`，属性写 CRLF 而工作树是 LF）；
- 声明计数**单独**改成 3 → `declaration-drift` 红；
- 两枚调用**挪到首枚采样之后** → `ordering` 红（逐字 `not before the first sampling call`）。

**还原证（AC#2 后半）**：`git status --porcelain -- scripts tools .github` **空**（18:30:02 一次；两笔交件 commit 后的 **18:35:28** 再跑一次仍空，见 §5）。
真件 `scripts/slo-check.ps1` 在 `fae8ba704..6589a0625` 区间**一字节未动**（`git diff --stat` 空）。

## §3 AC#3 —— 四条偷懒形状逐条点名（**我没做哪一种**）

1. **动阈值／golden／`internal/observe/thresholds.go`**：**没做**。`git show --name-only 6589a0625` ＝ 1 枚新脚本；
   `internal/observe/**` 与我零相交（`git status --porcelain -- internal` 只列别人的一枚未跟踪件 `internal/tools/ticket23_image_channel_test.go`，与票 266 无关）。
   本尺里没有任何阈值、没有任何 `<=`/`-gt` 判定业务数字（唯一的比较是"计数/行号"，不是 SLO 数值）。
2. **放宽 `Set-StrictMode`**：**没做**。`scripts/slo-check.ps1` 一字未动（见上），`Set-StrictMode -Version 2.0` 仍在 `:82`。
3. **把某一步改成"跳过"**：**没做**。本尺没有可跳过的步；正控失败时它**拒绝**出读数（`positive control RED - the census is not run` ／ exit 1），
   而不是跳过；§6 的新 CI 步逐字**不带** `continue-on-error`、不带 `if: false`。⚠ 但请照 §1 第 2 条读：
   **今天它没有 CI 调用点**（`.github/**` 不在写面）——"已入 CI"这句话本腿没说。
4. **用 `try/catch` 吞**：**没做**。`sh` 里没有 `catch`、没有 `|| true`、没有吞 rc 的构造；全文只有两处 `2>/dev/null`
   （`mktemp` 与 `mktemp -d`，只为不把 util 的噪声混进读数），且两处的失败都被 `||` 接住并转成具名 **CANNOT DECIDE（exit 2）**，
   ⛔ 不是绿。现读：`grep -n "2>/dev/null" scripts/slo-nail-census.sh` ＝ **2 命中**（`:306` `rowsfile=$(mktemp 2>/dev/null)`、`:429` `BENCH=$(mktemp -d 2>/dev/null)`），
   `grep -n "|| true" scripts/slo-nail-census.sh` ＝ **0 命中**；`wc -l` ＝ 569 行（＝commit 的 569 insertions）。

## §4 AC#4 —— 卫生门读数不扩大（带取数时刻，红名逐名比对）

**A. `tools/d22scan`（含票 212 ban #9）**

| 何时 | 命令 | rc | 读数（逐字） |
|---|---|---|---|
| 落地前 18:25:52→18:26:59 | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=363, bans #1-5 cmd/=57, ban #6 frontend/=116, ban #7 internal/tools/=48, ban #8 design/=39, ban #8 frontend/=116, ban #8 internal/=903, ban #8 cmd/=209 …` |
| 落地后 18:28:41→18:30:52 | 同上 | **0** | **与上一行逐字相同**（同一串分母，无增无减） |

ban #9（票 212 幻引）在同一次 scan 内判、0 命中（`clean - no D22 ban violations`）。⚠ 本腿**未**走候选 ③：
`tools/d22scan` 的 scope 名册一枚未动（`grep -c '"scripts' tools/d22scan/main.go` 现跑 ＝ **0**，`rc=1`）⇒ 分母未变、无需解冻。

**B. 票面长度门 `scripts/check-path-length-budget.sh --with-self-test`**

| 何时 | rc | 逐字读数 |
|---|---|---|
| 落地前 18:27:09→18:27:18 | **1** | `denominator read: 18772 tracked paths`；`tracked paths=18772 over-budget=65 covered by roster=64 not in roster=1`；self-test `3/3 PASSED`；`VERDICT RED` |
| 落地后 18:28:41→18:29:05 | **1** | `denominator read: 18773 tracked paths`；`tracked paths=18773 over-budget=65 covered by roster=64 not in roster=1`；self-test `3/3 PASSED`；`VERDICT RED` |

**红名集合逐名比对（两刻同两枚，一名不差）**：

1. `over budget and NOT in the roster: .scratch/wisp/issues/179-loop-deciderisk-…-holds-that-sentence-done.md (relative 175 chars, name 154 chars, … budget 165)`
2. `roster entry that is not an over-budget tracked path any more (the roster has drifted away from the tree): .scratch/wisp/issues/179-loop-deciderisk-…-holds-that-sentence.md`

（另有一枚 `seeded-xxx…-over-the-hat.md` 是本门 **self-test 的植入正控**，两刻都在、属期望输出，不属红名集合。）

**归因（⛔ 不自行翻豁免）**：这枚红**与本腿无关、且在本腿之前就红**——票 179 已被别人改名加 `-done`
（`git log --oneline -2 -- ".scratch/wisp/issues/179*"` ＝ `846c37cfc …`，现名 `…-sentence-done.md` 存在且 tracked），
而门的 roster 仍写旧名 ⇒ "roster 漂离树"。未动 roster、未动阈值、未动 `tools/d22scan/allowlist.txt`（票面禁区），报给编排者。

**分母增量（具名）**：`18772 → 18773`，唯一增量 ＝ 本腿新增的 tracked 路径 `scripts/slo-nail-census.sh`（26 字符，远在 100 字符帽内，未进 over-budget）；
`over-budget` 65→65、roster 覆盖 64→64 **不变** ⇒ 读数不扩大。

## §5 门禁 rc 表（本腿现跑）＋ 交件干净度

| 门 | 命令（要点） | rc | 读数 / 取数时刻 |
|---|---|---|---|
| 语法（bash） | `bash -n scripts/slo-nail-census.sh` | **0** | 18:30:02（交件字节上） |
| 语法（sh） | `sh -n scripts/slo-nail-census.sh` | **0** | 18:30:02 |
| **CI 的 sh 是 dash** | `dash scripts/slo-nail-census.sh` ＋ `--self-test` | **0 / 0** | 18:33:49–18:34:21；真读 `PASS - 2/2 nails … :533`；正控 `32 control(s) ran, 0 failed` |
| **CI 的 awk 是 mawk** | 同上，`PATH` 前置 `gawk --posix` 垫片（POSIX ERE 代理） | **0 / 0** | 18:33:49–18:34:21；同读数。垫片＝`/tmp/tmp.jGY73DTlfg/awk`（`exec gawk --posix "$@"`）；本机无 mawk（`command -v mawk` 空），故这一发是**代理**、不是 mawk 本体，具名存档 |
| 自控＋真读 | `sh scripts/slo-nail-census.sh --with-self-test` | **0** | 18:25:21；末三行＝`nail 2/2 …:214/:230/:292`、`PASS - 2/2 … :533`、NOTE |
| Go 构建 | `go build ./cmd/... ./internal/... ./tools/...` | **0** | 18:25:15→18:25:26，空输出 |
| D22 扫描 | `sh scripts/d22scan.sh`（×2） | **0 / 0** | 见 §4A |
| 路径长度门 | `sh scripts/check-path-length-budget.sh --with-self-test`（×2） | **1 / 1** | 见 §4B（既有红，非本腿） |
| 写面干净 | `git status --porcelain -- scripts tools .github` | 空 | 18:30:02（commit `6589a0625` 后）＋ **18:35:28（两笔交件 commit 后，交件时刻）** |

## §6 形 ①-A 的另一半：CI 步（⛔ 写面外，逐字补丁交编排者）

**位置**：`lint` 作业，紧接 `"Tracked path-length budget (ticket 262)"` 步之后（现 `.github/workflows/ci.yml:176-207` 那一步的下一行）。
**形**：ADDED STEP ONLY —— 不动 `on:` / `paths:` / 并发块 / `runs-on` / 任何既有步，无 `continue-on-error`，无 `needs:`。

```yaml
      - name: "SLO nail census (ticket 266 - the outer tooth over slo-check.ps1's nails)"
        # ADDED STEP ONLY. This is a text reader with no toolchain: it declares in
        # scripts/slo-nail-census.sh that the D32 gate carries exactly two
        # self-locking nails, and it goes red when one of those call lines is
        # deleted or commented out (ticket 263 mutant delete-nails measured an
        # all-green screen with both calls removed; ticket 266 is that hole).
        # It refuses to print a verdict at all - exit 2 - when the target is
        # missing/empty or its sampling anchor is gone, so "did not look" cannot
        # read as green. Runs its own 32-shot positive control first and refuses
        # to print a reading if the control cannot be made red. No threshold, no
        # golden, no allowlist entry, and no existing step touched.
        if: ${{ !cancelled() }}
        run: sh scripts/slo-nail-census.sh --with-self-test
```

**未接线的后果（诚实档）**：今天这枚牙的唯一调用方式是人手
`sh scripts/slo-nail-census.sh --with-self-test`；CI 里零调用点 ⇒ "删 `slo-check.ps1` 的两枚调用"这件事
**在自动门上仍无人发现**（P3 按天钝、且 dev 上 inert，见 §1 第 2 条）。接上 §6 那一步之后，这一形才闭环到"每次 push 都有人看"。
`--with-self-test` 的 32 发控制全在台件上跑（不碰被检件），单跑大约 1 秒级，对 lint 作业 hang bound（20 min）无压力。

## §7 边界与"下一腿"

- **未做**：AC 框勾选（归编排者）；`.github/workflows/**`（写面外，§6 只给补丁）；`Q-78` 的形选择（③ 会改 `tools/d22scan` 的 scope/分母＝人工批准面，本腿未碰）；
  票面 `Status:`/票尾 bullet（票面不在写面——若需要，逐字候选在文末）。
- **未碰**（逐条现读核过）：`docs/PLAN.md`、`docs/specs/**`、`internal/observe/thresholds.go`、任何 golden、`tools/d22scan/allowlist.txt`、
  `frontend/**`、`design/**`、三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）、
  `internal/agent/**`、`cmd/wisp/**`、`internal/ball/**`、**`scripts/slo-check.ps1`**、**`scripts/slo-freshness.sh`**、**`scripts/check-path-length-budget.sh`**。
- **台件留档**：`/tmp/tmp.lNlkWi6EUW`（32 发控制的完整输出；`/tmp` 重启即失，关键红句已逐字抄进 §2）＋ awk 垫片 `/tmp/tmp.hFuHVcZFTa/`。
- **next**（下一腿精确剩余）：
  1. 要 ①-A 完整形 ⇒ 照 §6 加一步（需 `.github/**` 写面；本腿无此写面，⛔ 未擅自加）。
  2. 要"钉必须真跑"的运行时半场 ⇒ 形 ②：改 `scripts/slo-check.ps1` 让两枚钉把"已跑"记进 `slo-report.json`，
     再由 `slo-smoke` 作业尾部**新增**一步读该报告（缺字段／缺报告＝具名红）。代价：本机正控要跑整条 smoke（build＋采样），
     且 precheck 会把同波在飞腿的 `go.exe` 判成 contention 拒采 ⇒ 正控须在 hosted runner 上取读数。
  3. 编排者收件时要的那一发：在 CI（或任一 linux 机）上跑 `sh scripts/slo-nail-census.sh --with-self-test`，把它的 32 行控制读数落档。

**票尾 bullet 候选（⛔ 本腿未写票面，供编排者决定是否追加）**

```
- [2026-11-05T18:3x+08] agent=266-r2 did=impl（形 ①-A 的 scripts/ 半场）：新尺 `scripts/slo-nail-census.sh`（commit `6589a0625`，569 行，POSIX sh）——声明「D32 门该有 2 枚自锁钉」写在**被检件之外**，判①定义×1②自身 ok 句在自己体内③调用点×1（顶格/未注释/独立行）④位置（晚于定义、早于 `:533` 首枚采样）⑤声明自洽（计数=表行数）；读不到/空件/采样锚丢失 ⇒ `exit 2` 具名 `cannot-decide`「NOT a pass」。**AC#1 交**；**AC#2 交**：台件删两枚调用 ⇒ rc=1 逐字 `slo-nail-census-nail-call-missing`（两枚都点名）、还原 ⇒ rc=0 `PASS - 2/2`；32 发控制 0 失败（含行注释/块注释/单删/删定义/删 ok 句/缺件/空件/CRLF/声明漂移/乱序）；`git status --porcelain -- scripts tools .github` 空两次。**AC#3 四条偷懒形状逐条点名未做**（阈/golden/StrictMode/跳过/吞）。**AC#4**：d22scan 前后各 rc=0 且 `clean` 行逐字相同（未动 scope，`grep -c '"scripts'`＝0）；路径门前后各 rc=1，**既有红名集合两刻同名**（票 179 改名 `-done` 而 roster 写旧名，`846c37cfc`），分母 18772→18773（唯一增量＝本腿新件，26 字符、未进 over-budget），over-budget 65→65／roster 64→64 不变。⚠ **形 ①-A 的 CI 步未交（`.github/**` 写面外）**：今天本尺**零 CI 调用点**，逐字补丁在 `docs/evidence/s1/ticket266-266-r2-impl-v1.md` §6；接上之前「删钉」在自动门上仍无人发现（P3 按天钝、dev 上 inert）。next=§6 那一步（须 `.github/**` 写面）→ 形 ② 的"真跑"半场（`slo-report.json` 记账字段＋`slo-smoke` 尾部新步）。
```
