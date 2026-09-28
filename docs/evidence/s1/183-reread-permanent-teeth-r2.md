# 183-r2 — 测试面装牙：AC#9 的第 8 枚常驻腿 / AC#2b 的 CLI 接缝常驻判据 / AC#4 那一形 / AC#6 一句归口 / AC#8 门禁重取

- 角色：**实现者（写码腿·测试面）**。预期**零产码改动** —— 实测达成：本程名下只有两枚 `*_test.go`，删除列逐枚 0（尺见 §6 末）。
- 派单：`.scratch/wisp/dispatches/2026-09-28-135x-impl-183-r2-grow-the-two-missing-teeth-and-re-take-the-gates-after-your-own-last-commit.md`
- 依据表：`docs/evidence/s1/183-pointer-exemption-accept-v2.md`（§2 恒真性矩阵 M3 列零枚红＝AC#9 的来由；§3 四发进攻）· `docs/evidence/s1/183-pointer-exemption-r1.md`
- 交付是**增量**的：本文件先落骨架＋step-0（commit `cc12af88`），再落两枚判据（commit `afc766ee`），本版本＝终态。
- 本节以下凡自指文件行号，一律写"符号名＋`grep -n` 尺"，不写死号。
- 读数全部落在 `.scratch/wisp/probes/183/r2/**`（`logs/` 五份＋`mut/` 两份变异副本＋两枚 overlay 名册）。

## 1. step-0 五件（本程现量）

**(1) 时刻／分支／起手锚**

```
$ date
Mon Sep 28 13:47:39 CST 2026        # 2026-09-28 13:47 +0800（本机 CST）

$ git rev-parse --abbrev-ref HEAD
dev

$ git rev-parse --short HEAD | sed 's/./& /g'
5 3 0 5 5 0 3 e                      # HEAD = 5305503e
```

⚠ **起手锚差异（按派单要求登记，不据此改判）**：派单写 `98af6535`，盘上 HEAD＝`5305503e`。
`git log --oneline 98af6535..HEAD` 现量＝**恰好多出 1 枚**，且那一枚正是派单 183-r2 自己（`ledger(A366)＋票 183 翻四格补两格＋派 183-r2`）。
⇒ 派单写的是"派单前"的锚，HEAD 靠后一枚是编排者记账＋派单 commit；**按实际 HEAD `5305503e` 做**，终态尺也以此为比较点。

**(2) 在飞脏件闸门**

```
$ git status --porcelain -- internal/ cmd/
（空）
```

⇒ **空**＝没有别人在飞的产码脏件，我不需要 overlay；终态 §6 复量仍空。

**(3) 三向 md5（起手向＝本程第 3 枚调用现量；交件向见 §6）**

| 尺 | 派单预期值 | 起手实测 | 交件实测（§6） | 判 |
|---|---|---|---|---|
| `sed -n '468,473p' internal/risk/provenance.go \| md5sum` | `858e45116383caa3e7c1dd4b0924fad1` | `858e45116383caa3e7c1dd4b0924fad1` | `858e45116383caa3e7c1dd4b0924fad1` | 两向一致 |
| `sed -n '11,15p' internal/risk/taintmatch.go \| md5sum` | `5680ddd18e2d2ec2a85e485b54f4c12e` | `5680ddd18e2d2ec2a85e485b54f4c12e` | `5680ddd18e2d2ec2a85e485b54f4c12e` | 两向一致 |
| `sed -n '482,506p' internal/risk/provenance.go \| md5sum` | `f89e891e5eee3f3ea2b4f89d921072c4` | `f89e891e5eee3f3ea2b4f89d921072c4` | `f89e891e5eee3f3ea2b4f89d921072c4` | 两向一致 |

⇒ 两支冻结文字与 `MarkWithHostPath` 文档块原 25 行**一字未动**，本程也没往文档块里追加任何东西（追加面见 §12 第 1 条：那一句该由谁写）。

**(4) 判据枚数现量**

| 尺 | 起手 | 终态 |
|---|---|---|
| `grep -c "^func Test" internal/risk/pointer_183_test.go` | **7** | **8**（本程补的那枚＝AC#9 第 8 枚常驻腿） |
| `grep -c "^func Test" internal/tools/pointer_183_cli_seam_test.go` | 文件不存在 | **2**（AC#2b／AC#4） |
| `grep -rln "183" internal/tools/*_test.go`（AC#2b 的硬缺口尺） | **零枚**（＝派单/验收表的读数） | **1 枚**＝`internal/tools/pointer_183_cli_seam_test.go` |

票面"八枚常驻判据"（票 177 族）与我派单里的"六枚"都不是现量值；`183-v2` 与我这一程量到的都是 7（补前）。

**(5) 门禁基线（在册，起手现跑）**

- `sh scripts/d22scan.sh` ⇒ rc=0、`ban #8 internal/` ＝ **432**（与在册一致）。
- `bash .scratch/wisp/probes/154/gate-clauses.sh` ⇒ rc=1、红腿名册尺 `grep "BAD" … | grep -oE "腿=G[0-9a-z]+" | sort -u` ＝ **只有 `腿=G6neg`**（比名册不比退码）。
- 没跑 `probes/161/r6/flip-declaration.sh`（禁面）。

## 2. AC#9 — 第 8 枚常驻腿的"牙"

**补的腿**：`internal/risk/pointer_183_test.go` 里的 `TestPointer183DeclaredPathAbsentFromBodyStillHits`（符号尺 `grep -n "^func TestPointer183DeclaredPathAbsent" internal/risk/pointer_183_test.go`）。
**形状**：新增夹具 `ptr183BodyAbsentPath` —— 同一族宿主 stub，但**正文不含被声明的那条路径**，只拼出它的 ≥8-rune 窗口（`ptr183DirSpelling` 盘符＋目录前缀，以及 `\tool-output-` 那段），断言 `Inspect(fs.read,{path: ptr183Path})` **仍命中 R4** 且归因 `task.output`。
**为什么这一形才分得开**：`contains()` 跳过"拼得出声明串"的候选窗口；只要声明串被无条件挂上，那条路径的**每一枚**候选窗口都拼得-out 它自己 ⇒ 整发变 clean。M3 红就红在这里；HEAD 上因为 `runeIndexOf` 那层守卫（符号尺：`grep -n "declared path is not in this content" internal/risk/provenance.go`），声明没挂、正文窗口照旧命中。
两枚前置断言（反恒真）：正文里声明路径**确实不存在**（`runeIndexOf(...) == -1`）、且正文**确实拼出**声明路径的一枚 ≥`contractMinFragmentChars` rune 窗口——否则"仍命中"会因为根本命中不了而假绿。

**两向读数（这一格测的是"腿存在＋那一支成立"，不是自证缺陷）**

未修码（HEAD `5305503e`，`logs/risk-183-on-HEAD.txt`）＝**绿**，且八枚全绿：

```
--- PASS: TestPointer183DeclaredPathAbsentFromBodyStillHits (0.00s)
ok  	github.com/CarlosShao/wisp/internal/risk	0.074s
```

M3 变异（不要求 `runeIndexOf` 命中就无条件挂 `declaredNorm`；副本 `mut/provenance-M3.go`、`overlay-M3.json`，`logs/risk-183-on-M3.txt`）＝**这一枚单独红，其余七枚照旧全绿**：

```
    pointer_183_test.go:386: AC#9 fail-closed: a body that never spells the declared path must not exempt a reread of it (if this goes red, MarkWithHostPath's absent-path branch became a comment again)
--- FAIL: TestPointer183DeclaredPathAbsentFromBodyStillHits (0.00s)
FAIL	github.com/CarlosShao/wisp/internal/risk	0.076s
```

⇒ **`183-v2` 的"七枚全绿、exit=0"那一发从此不再可能**：那支 fail-closed 现在有一枚常驻判据钉着。**纯测试面**：产码一行没动（`git status --porcelain -- internal/ cmd/` 在 M3 跑完仍空，见下），我也不曾把 M3 的读数当"缺陷在盘上"——它是我自己摘掉守卫才出现的。
变异走 `-overlay`，**不改工作树**（比在跟踪件上就地改再还原更稳；`probes/**` 既有台件只读、零改动）。

## 3. AC#2b — CLI 接缝形状的常驻判据（含回填产生的那条 `ArtifactPath`）

**补的腿**：`internal/tools/pointer_183_cli_seam_test.go` 的 `TestPointer183BackfilledArtifactRereadStaysClean`。
**形状照派单点名的那一支**，不是 175-r2 那对桥级腿的复用：`TaskBackfill.Backfill`（真 `agent.Spiller`、`agent.BudgetsFor(128000)` ⇒ `SpillTokens=4000`，20000 字节真落盘）写名册 → `task.output` 当场从名册读出、在 `rec.ArtifactPath != ""` 那一支里 `hostPathBoxFromCtx(ctx)` + `box.set` → `bridge.mark(dec, res, hostPaths.get())` → 模型照桩里那条指针发真 `fs.read` ⇒ 判 **L0**、`fs.read` 成功、**20000 字节逐字节回全**。
反恒真前置（写死在腿里，缺一枚就 `t.Fatalf`）：回填理由为空、名册里的路径真读得回来、产物名前缀是 `tool-output-agent-task-`（＝真 CLI 那一族名字）、桩里的指针逐字等于名册登记那条、**桩正文里 `-output-` 的出现次数多于指针自身带来的次数**（这一条就是它与票 175-r2 那对载具腿的分界：它们的正文是一串无差别 ascii）、`task.output` 确实盖了戳。

**读数**

⚠ **派单写"今天它应该能红"——没兑现：它在 HEAD 上是绿的。** 我没有为了让它红而放宽任何东西，也没改产码。逐字（`logs/tools-183-on-HEAD.txt`）：

```
--- PASS: TestPointer183BackfilledArtifactRereadStaysClean (0.03s)
--- PASS: TestPointer183SiblingArtifactInTheSameDirectoryStillHits (0.02s)
ok  	github.com/CarlosShao/wisp/internal/tools	0.096s
```

⇒ 缺口是**"常驻面缺席"那一枚硬事实**被这一格补上了（`grep -rln "183" internal/tools/*_test.go` 从零枚变一枚），而"端到端跑通一次"这件事从此不必再靠 `probes/` 台件背书。这一格该判达成还是仍欠，由编排者翻勾（AC 框我一枚没动）。

**它的牙（红的那一发在具名改动之后，不是我造的形状）**：摘掉 `internal/tools/task.go` 里那三行 `box.set`（保留名册条目＝175-r2 文件头给 J-3 命名的那一支变异；副本 `mut/task-boxset-off.go`、`overlay-boxset.json`、`logs/tools-183-on-boxset-off.txt`）⇒ 两枚新腿**双双红**，逐字：

```
    pointer_183_cli_seam_test.go:208: 照宿主指针的续读判成 "L2", want L0（R4 又咬住宿主自己写下的指针了）
    pointer_183_cli_seam_test.go:212: 续读那一发被拒了：{Text:L2 审批通道尚未接入（票 21），已拒绝执行 IsError:true RiskLevel:L2 ErrorClass:user_rejected Truncated:false AppliedSteps:[]}
--- FAIL: TestPointer183BackfilledArtifactRereadStaysClean (0.02s)
    pointer_183_cli_seam_test.go:269: 同一发桩里宿主自己声明的那条路径判成 "L2", want L0：豁免连声明方都不放了
--- FAIL: TestPointer183SiblingArtifactInTheSameDirectoryStillHits (0.02s)
FAIL	github.com/CarlosShao/wisp/internal/tools	0.207s
```

跑这一程用 `-overlay`，工作树零改动（跑完 `git status --porcelain -- internal/ cmd/` ＝空，见 §6）。

## 4. AC#4 那一形 ＋ AC#6 一句归口

**AC#4（同目录兄弟路径 × CLI 形状）**：`TestPointer183SiblingArtifactInTheSameDirectoryStillHits`（同一枚文件里）。
形状：两枚真背景任务各落一份产物到**同一个 artifacts 目录**，任务一的正文逐字拼出任务二那条路径 ⇒ `task.output` 的桩里同时出现两条路径，而 `box.set` 只声明**自己那一条**（载具本来就只装一枚，`hostPathBox.set` 后写覆盖）。断言：

1. 对照半（不许拿"全都挡"当绿）：声明方那条照旧 **L0**；
2. 兄弟那条判成非 L0（读数 L2），且 `prov.Inspect(task,"fs.read",{path:sibling})` 命中并归因 `task.output`＝命中的是宿主正文自己拼出的兄弟路径那一维。

**哪一发改动会让它红**：把豁免从"这一条路径"放宽到**目录／前缀级**（＝`MarkWithHostPath` 文档块明令不属于已批形状的那一支）⇒ 兄弟那条变 L0、`Inspect` 不再命中 ⇒ 这一枚当场红。第二支：`fs.read` 从 C25 名册摘出去"别盖戳"（AC#5③ 禁的那一支）同样红。跑法：`go test -count=1 -overlay=.scratch/wisp/probes/183/r2/overlay-boxset.json ./internal/tools/`（上表里第 269 行那一发就是它）。
⇒ 票面 AC#4 原来只做到"改一个 rune 的近邻路径"（包级 R6），这一格补的是 **W-3 那族（同名／兄弟路径）在 CLI 形状上的对应物**；`183-v2` §6 第 2 条具名登记的"不同目录下同名产物"仍**不在本程读数里**（见 §5 第 3 格）。

**AC#6（`窗口 0.0s`）— 只归口、只取读数，本票不修**

秒数读数（尺：`grep -rn "窗口 " .scratch/wisp/probes/183/r1/logs/ .scratch/wisp/probes/179/r2/logs/`，两处台件各自独立复现同一值）：

```
.scratch/wisp/probes/183/r1/logs/cli-reread-landed.txt:25:          窗口 0.0s
.scratch/wisp/probes/183/r1/logs/e2e-readings.txt:42:  窗口 0.0s
.scratch/wisp/probes/179/r2/logs/e2e-readings.txt:38:  窗口 0.0s
```

（我 §3 那条红读数里同一支通道的措辞：`L2 审批通道尚未接入（票 21），已拒绝执行`。）

**具名归口＝票 162 那一族（审批通道未接），不是"超时常量在此形状下取零"**，两条排除理由都是代码面现量（尺可复制，我没动任何常量）：

- 打印侧：`grep -n "窗口 %.1fs" cmd/wisp/run.go` ⇒ 打的是 `p.Window.Seconds()`，即**回执 Prompt 里那个字段**；
- 填充侧：`grep -n "MinL1Window\|MaxL1Window\|win <= 0" internal/agent/approval/queue.go internal/agent/approval/gate.go` ⇒ `MinL1Window=2s`、`MaxL1Window=3s`，且构造期 `win <= 0` 一律抬到 `DefaultL1Window`；配置侧 `grep -n "L1WindowSec" internal/config/schema.go` ⇒ `default:"2"`。
  ⇒ 走 `approval.Gate` 出来的那张卡**给不出 0.0s**；0.0s 是 Prompt 的 `Window` 未被窗口填充那一支的零值，与"L2 通道尚未接入即被拒"同批发出。
  ⚠ 半句标〔推断〕：`0.0s` 的具体构造点我**没**在真机 CLI 上自己复跑（见 §5 第 1 格），排除法只到"常量侧不可能取零"这一层；`窗口` 秒数读数本身是既有台件日志里的逐字值，我改不了它。
- 冻结件零改动：审批超时常量／`thresholds.go`／golden／`allowlist.txt` 一字节没碰（§6 的 md5 与删除列两把尺代表同一份树）。

## 5. 本程没测什么（逐枚具名＋为什么）

1. **真机 CLI 腿一枚没跑**：没跑 `-overlay` 进 `cmd/wisp` 的那一发，也没跑 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` 前缀那一把与不带前缀的 `0xc0000135` 那一把（两形都无读数）。⇒ 票面 AC#8 里"CLI 那一面走 -overlay 且带在册 PATH 前缀"那一格，本程**没有凭据**；AC#2b 我是按派单 §2 给的形状落在 `internal/tools/` 包级，不是真机。AC#6 的 0.0s 因此只有日志复量＋代码侧排除。
2. **`probes/161/r6/flip-declaration.sh` 没跑**（禁面）⇒ `gate-clauses.sh` 的"基线过期表"那一格我今天只是照读名册，没验它自己那把尺。
3. **票 177 W-3 的"同名／不同目录"那一支**（183-v2 §6 第 2 格同款）仍未裁：我补的是同目录兄弟那一支。
4. **`-race` 没跑**：`go test -race ./internal/risk/ ./internal/tools/` 零读数。
5. **票 185 那一发（同一条路径的第二次有界续读）没测**，也不该由我压安静：AC#3／177 W-2 的边界还在，我没碰 `bridge.go`、没碰 D15 再落盘层。
6. **M3 那支 fail-closed 的真机可利用性没测**（要动 `internal/tools/bridge.go` 现场，越面）⇒ 我只证明"摘掉守卫会有判据红"，没证明"桥会不会真给出正文里没有的 hostPath"（183-v2 §6 第 9 条同一格，仍空）。
7. **`notify` 通道那一发没跑**（183-v2 §4(a) 留的那把尺）；`go test ./internal/memory/`、`./internal/config/` 等其余包没跑（不在派单 §4 的四枚名册里）。
8. **F2 两版等价性（索引侧 map／复核侧逐字）没裁**；频率口径那笔更正（`3 of 10` ⇒ 真文件 2/8＋自造 1/2）不是我的面，没重跑。

## 6. 门禁终态（AC#8）

⚠ 取数点：本程**最后一枚碰 `internal/` 的 commit ＝ `afc766ee`**（两枚判据件）之后逐枚现跑；其后的 `docs/evidence/s1/…` 与本票 Progress log 两枚 commit 不含 Go 件，不改变下面任何一格（下面最后一节是它们在终态树上的复量）。

| 尺 | 读数 | 判 |
|---|---|---|
| `go test -count=1 ./internal/risk/`（**单包**，A359） | `ok github.com/CarlosShao/wisp/internal/risk 3.771s` rc=0 | 绿（`logs/gate-internalrisk.txt`） |
| `go test -count=1 ./internal/tools/` | `ok github.com/CarlosShao/wisp/internal/tools 14.027s` rc=0 | 绿（含我新增 2 枚，`logs/gate-internaltools.txt`） |
| `go test -count=1 ./internal/agent/` | `ok github.com/CarlosShao/wisp/internal/agent 1.691s` rc=0 | 绿（`logs/gate-internalagent.txt`） |
| `sh scripts/d22scan.sh` | rc=0、`ban #8 internal/ examined **433** Go files`、`clean - no D22 ban violations` | 432 基线 **＋1**＝本程新增那枚跟踪判据件 `internal/tools/pointer_183_cli_seam_test.go`，具名登记（`logs/gate-d22scan.txt`） |
| `bash .scratch/wisp/probes/154/gate-clauses.sh` | rc=1；`grep "BAD" … \| grep -oE "腿=G[0-9a-z]+" \| sort -u` ⇒ **只有 `腿=G6neg`** | 与在册名册逐字相同（**比名册不比退码**；`logs/gate-clauses.txt`） |
| `"$(go env GOPATH)/bin/gofumpt" -l <名下两枚件>` | 输出为空 | 名下件全格式化 |
| 三向 md5 交件向 | `858e45116383caa3e7c1dd4b0924fad1`／`5680ddd18e2d2ec2a85e485b54f4c12e`／`f89e891e5eee3f3ea2b4f89d921072c4` | 与起手向、与派单基线**三向两枚皆同** |
| `git status --porcelain -- internal/ cmd/` | 空 | 变异（两发 -overlay）没在树上留任何东西 |

**零产码改动两把尺**：`git diff --numstat 5305503e..HEAD`（起手锚→含表件）＝

```
115	0	docs/evidence/s1/183-reread-permanent-teeth-r2.md
60	0	internal/risk/pointer_183_test.go
296	0	internal/tools/pointer_183_cli_seam_test.go
```

删除列**逐枚为 0**；名下产码件（`provenance.go`／`taintmatch.go`／`internal/agent/**`／`cmd/**`）**不在名册里**＝一行没动。

## 7. 被拒／没成功的调用（逐条）

1. **`build/overlay` 类 rc=1，不是判据红**：第一次做 `box.set` 变异时 python 断言失败（`AssertionError: box.set anchor not unique`，我抄的锚文本缩进不对）⇒ `overlay-boxset.json` 当时**没写出来**，那一发 `go test -overlay=…` 直接 rc=1 且**没有任何 `--- FAIL`/`--- PASS` 行**。那是我的台件失败，不是判据红，我没把它当读数。改用正则匹配后第二发才取到 §3 的真红。
2. **我自己草稿里的两处缺陷在跑之前发现并改掉**（`Edit` 各一发，非工具失败）：一枚写死的 `filepath.Clean("")`（撞禁面"在 `risk.PathResolver` 之外用 `filepath.Clean|Abs`"，`d22scan` 会判违规）与一套自造的 JSON 数字编码助手（包内已有 `jsonEscape`）。终态件里两枚都不在场（尺：`grep -n "filepath.Clean\|filepath.Abs" internal/tools/pointer_183_cli_seam_test.go` ⇒ 空）。
3. **权限系统零拒绝**；无 `pathspec did not match` 之类 commit 失败（两枚 commit 都先 `git add <显式路径>` 再 `git commit -q -F - -- <同一路径>`）。
4. 一次 `grep -A4 -- "--- FAIL"` 取不到红消息（Go 把 detail 打在 `--- FAIL` **之前**）⇒ 换成 `grep -B3` 才取到 §2 的逐字，没影响读数。

## 8. 有没有跑过删除命令

**没有**。全程零 `rm`／`git clean`／`git restore`／`checkout .`／`reset`／`--amend`／`rebase`／`stash`／`switch`／`merge`／`worktree`／push；临时件（`probes/183/r2/**`＝五份 logs＋两份变异副本＋两枚 overlay）**只建不删**，连作废的那一发（§7 第 1 条）也留在盘上。两发变异一律走 `go test -overlay=`，**连"改了再还原"都没发生**，所以没有需要还原的树。

## 9. 工具调用枚数 vs 硬顶 40

**40 / 40**（本表写入＝第 37 枚；票面 Progress log 追加＝第 38 枚；表＋票的 commit＝第 39 枚；终态三把尺复量＝第 40 枚）。
增量交付照做：第 9 枚就落了 step-0＋骨架 commit（`cc12af88`），第 32 枚落判据 commit（`afc766ee`）。最后一发**新现量**在第 30 枚（M3 的逐字红消息），其后全部用于写表与收尾 —— ⚠ 比派单"第 28 枚起停止开新战场"晚两枚，超出的两枚花在 AC#2b 的牙齿验证（`box.set` 变异）与其重跑上，属于交付必答的"红读数"，不是新战场；登记在这里，不按派单预期"今天就能红"私自改判据把它拧红。

## 10. 伪授权两栏

**收到、判为真授权**：派单 `2026-09-28-135x-impl-183-r2-…`（五格、帽 40、前 8 枚内先 commit、落点 `probes/183/r2/**`、终态三把尺）· 票 183 票面（AC 逐字、"本票不解决"清单、AC#2b／AC#9 的格面原文）· `AGENTS.md` §1（禁面、git 纪律、AC 框由编排者翻）· `183-v2` 表的 M3 现量（**依据**，不是授权）。

**遇到、判为不是授权、没照做（具名）**：
1. 起手锚 `98af6535` vs 盘上 HEAD `5305503e` ＝ **锚漂**，按派单"登记差异并按实际 HEAD 做"处理，不改判；
2. 派单 §0(4) 与票面的"六枚／八枚"判据数＝**非现量值**，一律以 `grep -c "^func Test"` 现量为准（7→8）；
3. 派单 §2 "AC#2b 今天应该能红"＝**编排者预测**，我的读数是绿的（§3）⇒ 按现量报，不为了对上一句预测去改松／改紧判据；
4. `probes/183/accept-v1/**`＝死程现场，只当**存在性线索**，不删不动、不引用其读数；
5. `probes/183/r1`、`probes/183/v2` 的读数＝**别人程的自述**，我只有 `窗口 0.0s` 那一发是从它们的日志文件里逐字**复量**（只读），其余每格数字都出自本程 `probes/183/r2/**`；
6. 工作树里不是我的脏件（`.gitignore`、`probes/152/my152.py`、`probes/161/r6/logs/flip-*`、`docs/evidence/s1/152-*.md`、`design/**` 的未提交删除）**没提交、没还原、没评论**；`frontend/**`、`design/**` **没读没写没引**；
7. 任何"顺手加固产码"的冲动＝**没做**：那支持 `if lo := runeIndexOf(...)` 本来就是对的，本程只补判据（AC#9 明令纯测试面）。

## 11. 凭据值零抄录

本程没接触任何 API 密钥／DPAPI 明文／`.env`。表内与件内出现的唯一路径形态是 `C:\Users\swq\AppData\Roaming\wisp\artifacts\…`，那是**已在仓内夹具里**的常量形状（`pointer_183_test.go` 的 `ptr183Path`），我新造的夹具走 `t.TempDir()`，不含用户真实内容；配置文件里只引用变量名（`api_key_ref`）。

## 12. next=

1. **`MarkWithHostPath` 文档块该追加一句"第四枚面"＝人工批准面，我不写**：那句 "those three faces are the permanent legs of `internal/risk/pointer_183_test.go`"（尺：`grep -n "those three faces" internal/risk/provenance.go`）现在**数不足**——本程补的第 8 枚钉的是第四支（声明缺席＝不排除）。文档块属"只许追加不许改写"的那 25 行地界，**要不要追加、由谁追加＝先批再动**；批下来的话追加句该具名 `TestPointer183DeclaredPathAbsentFromBodyStillHits`。
2. **AC#2／AC#2b 翻不翻由编排者判**：常驻面已补（`internal/tools/` 里 `grep -rln 183` 从零枚→一枚）、逐字节回全在包级钉住、牙齿在 `box.set` 那一支变异下验证过；真机 CLI 那一形本程没跑（§5 第 1 格），若 AC#2 的格面坚持"必须真机"，缺口还在**票 183 AC#8 的 CLI 那一格**，不在这一枚腿。
3. **票 185 仍独立**（第二次有界续读），别拿本程的绿顶掉它；`probes/183/r2/logs/tools-183-on-boxset-off.txt` 那条 `L2 审批通道尚未接入（票 21）` 也是票 162／AC#6 归口时可直接复用的同措辞读数。
4. **谁补都行**（本程具名没跑的）：真机 CLI 两形（带／不带 PATH 前缀）、`go test -race ./internal/risk/ ./internal/tools/`、`notify` 通道那一发、W-3 的"同名不同目录"那一支、`flip-declaration.sh` 的基线过期格、M3 那支的真机可利用性。
