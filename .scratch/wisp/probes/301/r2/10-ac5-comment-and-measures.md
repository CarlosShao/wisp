# 301-r2 / 10-ac5-comment-and-measures — AC#5 注释面落地件（写腿自量）

双锚：起手 `HEAD=7e32af9d`／branch `dev`／`git status --porcelain -- scripts internal .github docs | wc -l`=0／`date`=`2026-10-10 15:35:49 +0800`（见 `00-anchor.md`）。
本腿第 1 笔 commit＝`4be5168c`（anchor，先于任何长跑）。本件是第 2 笔的唯一靶文说明。

⚠ 全文只用**内容锚**（逐字片段）。本件里出现的任何 `:NNN` 都只标"测量时刻的位置"，⛔ 当引用锚。

---

## 1. 改前／改后逐字对拉

### 改前（`git show HEAD:.github/workflows/ci.yml`，`test-core:` 块头，测量时刻 `:400-401`，逐字 2 行）

```
  # Windows-only packages (ball GUI, cgo speech) are out of this job's scope
  # by platform, not skipped: they run in test-windows / slo jobs.
```

紧跟其后的两行（⛔ 未动，用作位置证明）：`  test-core:` ／ `    runs-on: ubuntu-latest`。

### 改后（工作树，同一位置，19 行注释）

```
  # Windows-only packages (ball GUI, cgo speech) are out of this job's scope by
  # platform, not skipped -- that is a per-FILE statement, not a per-package one,
  # and the sentence this block replaced got its two named jobs wrong in different
  # directions. As checked on 2026-10-10 for ticket 301: ./internal/audio/ is named
  # in this tier's own array (`core)`), and a linux target still resolves sources
  # plus untagged test files there -- ruler (list only, runs nothing):
  # `CGO_ENABLED=0 GOOS=linux go list -f '{{len .GoFiles}}/{{len .TestGoFiles}}'
  # ./internal/audio/`. What ubuntu cannot evaluate is only the `//go:build windows`
  # set. The tail of that sentence was the half that lied, and its two halves are
  # asymmetric: test-windows has taken ./internal/audio/ since ticket 301 AC#1
  # landed (2026-10-10; the `windows)` scope array and win_pin in
  # scripts/portable-tests.sh both name it from that commit on -- re-read them,
  # this comment is not the evidence), while on that same date the slo tiers
  # reached no package list: slo-smoke and slo-full hand their gate to
  # scripts/slo-check.ps1, which builds wisp.exe and samples D32 states instead of
  # running package tests, and whose header restricts the smoke subset to
  # "(memory/handle subset, no audio)". Ruler for that half (whole file, prints
  # which lines match and nothing else):
  # `grep -nE 'portable-tests|go test|internal/audio' scripts/slo-check.ps1`.
```

`git diff -U0` 全文＝单枚 hunk `@@ -400,2 +400,19 @@ jobs:`，2 删 +19 加，逐行都是注释；全文已存 `logs/ac5-diff-U0.txt`（21 条 +/- 行）。

---

## 2. §3-3 那两把尺（命中数原文，⛔ 转述）

尺 ①（`git diff -U0` 的每个 +/- 行都以 `#` 开头）——对 `logs/ac5-diff-U0.txt` 打：

```
total +/- lines                                    = 21
reading-L  grep -c '^[+-]#'                        = 0
reading-P  grep -c '^[+-][[:space:]]*#'            = 21
非注释的 +/- 行（'^[+-][^+-]' 里再剔注释）        = 0
```

⇒ 按**字面**（`^[+-]#`＝注释符顶格）＝0 命中＝**尺判死**；按**目的**（剥掉缩进后是注释行）＝21/21、非注释行 0＝**判活**。
差别只来自一件事：这个块本来就在 `jobs:` 映射里、注释符前有两枚空格缩进，**票面自己引的那两行也是缩进注释**，所以字面那把尺对它⛔ 可能成立。本腿按**目的尺**交件并**具名顶回**（见 §6 顶回件 #1），⛔ 把缩进剥到顶格（那会为了迁就尺而改动块的形状）。

尺 ②（新增行里 `^\s*(name|runs-on|if|run|steps|uses|with):` 命中＝0）——逐字尺与命中数：

```
git diff -U0 -- .github/workflows/ci.yml | grep '^+' | grep -cE '^\+[[:space:]]*(name|runs-on|if|run|steps|uses|with):'   = 0
```

尺 ③（**本腿自加的第三把**，因为①②只证明"每行像注释"，⛔ 证明"解析结果没变"）：

```
python + yaml.safe_load 比对 HEAD blob 与工作树（正规化后 repr 相等）：
jobs_headcount head = 6 | new = 6
test-core keys head = ['env', 'runs-on', 'steps']
test-core keys new  = ['env', 'runs-on', 'steps']
runs-on test-core head = ubuntu-latest | new = ubuntu-latest
head step names == new step names: True
PARSED_STRUCTURE_EQUAL = True
```

⇒ 可执行面⛔ 变，这条是硬证据而不是一句承诺。

---

## 3. §2 那两枚读数（本腿现量，⛔ 引上一腿）

### 3.1 slo 那半句到底⛔ 真到什么程度

`ci.yml` HEAD blob 里两档的调用点（内容锚，逐字）：

- `slo-smoke:` ＋ `runs-on: windows-latest` ＋ 步骤名 `SLO smoke gate` ＋
  `run: powershell -NoProfile -ExecutionPolicy Bypass -File scripts/slo-check.ps1 -Subset smoke -SecondsPerState 4`
- `slo-full:` ＋ `runs-on: [self-hosted, wisp-slo]` ＋ 步骤名 `SLO full gate (six states + settle + leak)` ＋
  `run: powershell ... -File scripts/slo-check.ps1 -Subset full -SecondsPerState 6`

两档其余步骤＝checkout / setup-go / 缓存 third_party / **Build wisp.exe** / upload artifact——⛔ 一枚 `--scope=`、⛔ 一枚 `go test`。
整枚 `ci.yml` 里 `internal/audio` 出现次数＝**0**（尺＝`git show HEAD:.github/workflows/ci.yml | grep -c 'internal/audio'`；包名册住在 `scripts/portable-tests.sh` 里，⛔ 住在 yml 里）。

`slo-check.ps1` 全文（625 行）那把尺，逐字命中行：

```
grep -nE 'portable-tests|go test|go\.exe|internal/audio|audio' scripts/slo-check.ps1   -> 3 hits
15:    (memory/handle subset, no audio) + settle + leak self-test.
336:# agents are mid `go test`. The resulting D32 reading can then be falsely red
361:$loadNames = @('go.exe', 'gofmt.exe', 'cgo.exe', 'compile.exe', 'asm.exe', 'link.exe',
```

⇒ `:15` 就是派单点名的那句"逐字写着 (memory/handle subset, no audio)"（**本腿读到了**，⛔ 转引）；`:336` 是把"别人正在 `go test`"当**干扰源**警告、⛔ 是它自己跑；`:361` 是把 `go.exe` 之类当**负载进程名**列进读数表、⛔ 是拉起它们。
文件里所有 `Test-` 形状的命中（`Test-Path`／`Test-ExitCodeInstrument`／`Test-ScriptShapeNail`／`leak fixture self-test`）＝PowerShell 自己的仪器自检，**与 `go test` 无关**。
⇒ 结论（现量支持）：**slo 那两档对任何包都⛔ 跑 `go test`**，它们造 `wisp.exe` 再采样 D32 状态；`-Subset full` 的六枚状态＝`@('Sleeping','Armed','Warm','Conversation','PanelOpen','WorkPeak')`（`smoke`＝`@('Sleeping','Warm')`），也就是**运行时把包链进 exe 里跑**，⛔ 是跑该包的测试。〔读码＋现跑尺，⛔ 需标"仅推"：唯一⛔ 实测的是"CI 机器上真实跑过哪些进程"，本腿只有尺⛔ 有跑过的日志〕

### 3.2 `AC#1` 之后的真值（`git show HEAD:scripts/portable-tests.sh` 逐字）

`windows)` 档 scope 数组（现量逐字）：

```
windows)
    scope=(
        ./internal/proc/ ./internal/secret/ ./internal/config/ ./internal/risk/
        ./internal/ball/ ./internal/perm/ ./internal/plugin/ ./cmd/llmrecord/
        ./internal/session/ ./internal/projctx/ ./internal/audio/
    )
    pinned=$win_pin
    ;;
```

`win_pin` 枚数＝**11**，含逐字一行 `github.com/CarlosShao/wisp/internal/audio`。
枚档名册＝`tiers='core windows cli winsec census'`；census 那圈的 pin 对：`core) pin=$core_pin ;; windows) pin=$win_pin ;; cli) pin=$cli_pin ;; winsec) pin=$winsec_pin ;;`。

★**本腿新挖到的一枚事实（派单没写、但它决定文案诚实度）**：`core)` 档数组里**本来就有** `./internal/audio/...`（逐字行＝`        ./internal/buildinfo/... ./internal/audio/... ./internal/proc/...`）。
⇒ 那句"Windows-only packages ... out of this job's scope **by platform**"对 `internal/audio` **⛔ 是整包出档**，它同时在 `core)` 与 `windows)` 里。现量尺（本机 windows，⛔ 起测试、只 list）：

```
go list -f '{{len .TestGoFiles}}/{{len .XTestGoFiles}}' ./internal/audio/                       = 8/0
CGO_ENABLED=0 GOOS=linux go list -f '{{len .TestGoFiles}}/{{len .XTestGoFiles}}' ...            = 5/0
CGO_ENABLED=0 GOOS=linux go list -f '{{len .GoFiles}}' ./internal/audio/                        = 9
go list -f '{{len .GoFiles}}' ./internal/audio/                                                = 11
```

文件面佐证＝`git grep -nE '^//go:build' HEAD -- internal/audio` 只有 3 枚 windows-tagged 测试文件（`capturelevel_windows_test.go`／`hotplug_test.go`／`parse_wave_format_300_windows_test.go`），另 5 枚测试文件无 tag ⇒ **ubuntu 上 test-core 确实编译并跑掉这 5 枚**。
⇒ 文案因此写成"per-FILE⛔ 是 per-package"，并把那把 `GOOS=linux go list` 尺留在注释里。

---

## 4. 门禁三发（每发一份 rc 件，⛔ 0 字节；读数并排）

| rc 件 | 内容 | 结果 |
|---|---|---|
| `logs/bash-n-portable-tests.rc.txt` | `bash -n scripts/portable-tests.sh` | `rc=0`；pre `wisp.exe=0 balldebug.exe=0`／post 同 |
| `logs/d22scan.rc.txt` | `sh scripts/d22scan.sh` | `rc=0`；pre/post 读数均 0/0（件里 255 行＝表头 2 ＋扫描全输出＋ rc 行） |
| `logs/census.rc.txt` | `bash scripts/portable-tests.sh --scope=census` | `rc=0`；pre/post 均 0/0 |

census 判据（⛔ 看颜色，看名册）——逐字：

```
portable-tests.sh: census totals: packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0
portable-tests.sh: github.com/CarlosShao/wisp/internal/audio      8/0         corewindows
```

totals 行与派单期望**逐字相同**；`internal/audio` 那行**没变**（`8/0` ＋ claimed-by `corewindows`）。
`grep -cE 'STALE|stale' logs/census.rc.txt` = **0** ⇒ 与上一腿纠正过的机制一致：census 那支打完 totals 就 `exit`，永远进⛔ 到档路径里的 ledger/stale 那圈（"census 的 STALE 腿"⛔ 存在）。

tasklist 纪律：起手第 1 发读到 `wisp.exe=3`（第 2 发复核＝0，其间⛔ 是本腿杀的，本腿没杀过任何进程）；**三发门禁各自在跑前重新量一次**、并把 pre／post 两行写进同一枚 rc 件。起手那枚缺口（"只在起手量了"）本腿已按 §3-7 补上。

---

## 5. 本格⛔ 买到什么（诚实段）

这处改动**⛔ 修任何东西**：它⛔ 动 CI 行为——尺③（`PARSED_STRUCTURE_EQUAL = True`）就是那条断言的证据。它也⛔ 让 `internal/audio` 多跑一枚测试：让它在 `test-windows` 里进分母的是已经落地并验收过的 `AC#1`，⛔ 是这段文字。

它买到的只有一件事：**下一个读 `test-core` 那块注释的人，⛔ 能再把 audio 当成"平台出局、反正别处都覆盖了"。** 旧句把两处地方并成一句承诺，而这两处一直⛔ 同步——slo 那半句对任何包都⛔ 成立（slo 只造 exe 采样），test-windows 那半句要到 `AC#1` 落地才对 audio 成立；另外旧句还顺手撒了第三枚不大的谎（audio 在 `core)` 数组里，ubuntu 跑得掉它 5 枚无 tag 的测试文件）。文案把这三种形状都指回**那一天＋那两把尺**，⛔ 把答案抄成永久事实。

除这一段之外，本腿⛔ 多买任何东西，也⛔ 声称任何东西"已修好"。

---

## 6. 本腿⛔ 同意的派单预设（具名，带尺）

1. **§3-3 尺① 的字面形⛔ 可满足**。派单写"每个 hunk 的每个 +/- 行都以 `#` 开头"，字面尺 `grep -c '^[+-]#'` ＝ **0**（见 §2）。靶块在 `jobs:` 映射里、注释符前 2 空格缩进；**票面自己引的那两行也是缩进注释**，⛔ 可能满足字面尺。本腿⛔ 为此把注释顶格（那会改块形、且⛔ 是本次要修的东西），改交目的尺（剥缩进后 `#` 开头＝21/21、非注释行＝0）＋自加的解析等值尺③。**请裁：尺① 的文字是否要按"剥缩进"改写。**
2. **§1 "为什么这两行现在撒谎"只列了两枚缺口，实际是三枚**。派单说 audio "⛔ 在 `test-core`（ubuntu）的可求值集合里"。现量：`core)` 数组**逐字含** `./internal/audio/...`，且 `CGO_ENABLED=0 GOOS=linux go list` 给 `GoFiles=9／TestGoFiles=5` ⇒ 它在 test-core 里**可求值、也真被求值**，⛔ 可求值的只有那 3 枚 `//go:build windows` 测试文件（`8/0`→`5/0` 那把尺为证）。方向上"不对称"的预设（test-windows 半⛔ 真→⛔ 真／slo 半仍⛔ 真）**本腿核过、支持**；但"out of this job's scope by platform"这半句同样过期，⛔ 写进去的文案会继续教错下一位读者。
3. **⛔ 能按字面"只改两行、行数⛔ 变"来理解 §3-1**。派单⛔ 禁改可执行 yaml、⛔ 要单独一笔，但没规定替换必须恰好 2 行；本腿把 2 行换成 19 行注释（单枚 hunk、纯注释、解析等值＝True）。若编排者要的是**硬 2 行**，说一声，本腿可以压成 2 行（代价＝放不下面向读者的那把尺，文案会退回"另一枚现在时断言"）。
4. **§3-6 census 那行 claimed-by 的写法⛔ 是"两枚档"**：件里逐字是 `corewindows`（census 里 `where="$where$m"` 的**无分隔拼接**产物＝`core`＋`windows`）。本腿照派单核对了同一串字节，⛔ 因此判它"被 corewindows 那枚⛔ 存在的档认领"；**⛔ 存在的档叫 `corewindows`**，档名册＝`tiers='core windows cli winsec census'`，且 `ci.yml` 只有 `--scope=core`／`--scope=census`／`--scope=windows` 三处调用点。引用这枚串时请带上这句解释。
5. **⛔ 适用"第 45 次调用才 commit"**（派单 §0 已自明）：本腿是写腿，第 1 笔＝`4be5168c` 在任何长跑之前；⛔ push、⛔ 翻框（`AC#5` 的勾归编排者）。
