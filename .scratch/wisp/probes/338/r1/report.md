# 票 338 落地腿 338-r1 · 交件 report（写代理·格式归位）

- 起手锚＝我自己 `git log -1` 现认：`d19288a8`（2026-10-18 11:53，dev）。
  落 AC#0 时 HEAD 已漂到 `90dafe77`→`631a04ba`（工作树在飞，⛔ 拿编排者号当锚）。
- 我的落笔序列：`0b1d1f0b`（AC#0 尺＋logs）→ `87521cb4`（AC#1 tools 枚归位）→
  `3d26240c`（AC#2 改前四道门）→ `93d0f39c`（AC#1b panel 枚归位）→
  `ee75c23f`（AC#2 交件侧门＋lead 条件④）→ AC#3 logs → 本 report → 票尾一条 bullet。
- 尺＝`/d/work/base/gopath/bin/gofumpt.exe v0.12.0 (go1.27.1)` ⇄ 同版 `gofmt`；
  go＝`go1.27.1 windows/amd64`。落点一律 `D:\tmp\wisp\20261018\338r1\**`（先 `mkdir -p`，
  ⛔ 盘符根、⛔ 删任何旧树）。
- **四格一枚没翻**（AC#0–AC#4 全部仍是 `[ ]`，框由非实现者裁）；票名／`Status:` 未动；⛔ push。

## AC#0 起手现跑（票面四行逐枚重跑＋非空白尺自证）

- `gofumpt -l internal cmd`＝**7 行**（票面记 6）：多的那枚
  `internal/config/writeguard_handedit_227_test.go`＝`??` 未跟踪（227-r1 在写），
  不在 CI 分母（`ci.yml:189` 那道＝跟踪集）⇒ 按跟踪集过滤＝**6 行**＝与票面同。
- blob 定性（`git show HEAD:<f>` → 仓外 LF → 同尺）：真脏 2＝本票两枚；
  假脏 4＝CRLF 族（blob `gofumpt -l rows=0`）。逐枚 bytes／crlf_lines／rows 见
  `logs/ac0-start.txt` [3]。
- `gofmt -d` 那两枚逐字：tools 枚有 diff（注释列对齐）；panel 枚＝**0 字节**。见 `logs/ac0-start.txt` [3]。
- **非空白尺自证**（⛔ 交一把两种形状都放行的尺）：
  正控（纯空白）`rc=0` 空输出；负控（改一枚数字）`rc=1`＋逐字 `2c2 / <varb=2//x --- >varb=3//x`。

★与票面不符＝当场具名更正（lead 已在票尾裁语收下，凭据台账 `A1031`，撤销口令「撤 338 甲形」）：
票面"差的形状＝只有注释对齐空白"对两枚成立＝**不成立**。
- tools 枚：成立。`gofumpt -w` 一发收敛，非空白尺 `rc=0` 空输出，非空白字节 4966→4966。
- panel 枚：不成立。`gofmt -d`=0 字节（脏只属 gofumpt 档），`gofumpt -d` 要动 token 级三形
  ＝`type ( … )` 分组＋`jsPuncts` 切片拆行补尾逗号＋函数间补空行；非空白字节 62106→62104、
  `diff` 31 行、尺 `rc=1`；★一发 `-w` 不收敛（round1 后仍脏、round2 才净）。

## AC#1 归位（两枚）

- tools 枚（`87521cb4`）：`gofumpt -w` 一发，`rc_w=0`；**非空白尺＝`RULER_RC=0`＋空输出**
  （`diff <(tr -d ' \t' <改前>) <(tr -d ' \t' <改后>)`＝零字节零行）；
  `git diff --numstat`＝`5 5`；`git ls-files --eol`＝`i/lf w/lf attr/text eol=lf`；
  新字节 CR 行数＝0；复量 `gofumpt -l`＝0 行；二发 `-w` 幂等＝YES。逐字见 `logs/ac1-ac2.txt` [4][5][6]。
- panel 枚（`93d0f39c`＝lead 甲形放行）：逐发 `-w` 到收敛（round1 `rc_w=0`/l=1 → round2 `rc_w=0`/l=0）；
  条件① 改动只来自 `-w` 本身；条件⑤ 交件侧 `gofmt -d`=0 字节＋`gofmt -l`=0 行；
  numstat＝`34 20`；`i/lf w/lf`＋CR 0 行；非空白尺记录＝`RULER_RC=1`（47 行）62106→62107
  （甲形下只对这一枚让路，逐字存 `logs/ac1b.txt`）。
- ★另四枚 CRLF 假脏＝票面禁区逐字未动：`internal/agent/approval/pending_read.go`／
  `internal/agent/tools.go`／`internal/risk/provenance.go`／`cmd/wisp/models.go`
  **仍脏、⛔ 本腿射程**。尺＝wt sha256 三点对读（`logs/ac1-ac2.txt` [2]⇄[8]⇄`logs/ac1b.txt` [7]）
  逐枚同哈希：`ff26e84c…`／`14b2c067…`／`958364fa…`／`ad839c10…`。

## AC#2 门（成对＋四道门逐枚 rc）

名册逐名作差（同一台面，跟踪集分母）：改前 **6 行** ⇄ 交件侧 **4 行**
＝少的恰好本票两枚（tools＋panel），⛔ 多动（`comm -13`＝空输出，逐字见 `logs/ac1-ac2.txt` [7]）。
交件侧 4 枚＝上面那四枚 CRLF 假脏，逐名：`cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／
`internal/agent/tools.go`／`internal/risk/provenance.go`。

四道门（共享工作树，跑前 `unset GOTMPDIR TMPDIR`＝票 311 AC#1 shape B 卷闸定式，起手 echo 坐实）：
`sh scripts/d22scan.sh` rc=0 clean｜`sh scripts/check-path-length-budget.sh` rc=0 VERDICT GREEN
（15516 跟踪／over 60／roster 60）｜`go vet ./internal/tools/` rc=0｜`go vet ./cmd/wisp/` rc=0。
lead 条件④另两档：`go build ./...` rc=0（★派单禁共享树 build ⇒ 仓外导出树）；
`GOOS=linux go vet ./cmd/wisp/` **rc=1＝存量**，逐字红句
`imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in …@v1.13.8`，
归位尺＝改前锚 `0b1d1f0b` 同配方树复跑＝同 rc 同逐字红句、同树本机档正控 rc=0
（`logs/ac2c-ctrl-linuxvet.txt`）⇒ 与本票字节无关，⛔ 动 `go.mod`/`go.sum`/`ci.yml`。

## AC#3 颜色（两包成对，四数逐枚点名⛔ 混用⛔ 相减）

台面＝`trees-at-121746\**`（★嵌一层＝父目录名不合 8.3⇒有别名；v1 落点 `338r1` 本身合 8.3＝
无别名，把 252 族探针打成台面红＝我的落点形状，非判据响，具名更正）。
- `./internal/tools/`（0b1d1f0b ⇄ 87521cb4，同台面）：
  before `RUN=323 PASS_TOP=235 FAIL_TOP=7 SKIP=0` ⇄ after 同四数；名册作差 0／0。
  我那枚 `TestCanonicalizeErrorNoticeRelaysNoUpstreamErrorText174r7` 两侧 PASS。
  7 枚红逐名归位＝别票／本机存量：5 枚真 NTFS junction 造不出（`mklink /J` 非管理员 exit 1）＋
  1 枚跨卷写 `C:\.wisp-tmp-*` Access denied＋1 枚 `dirOf` 斜杠期望（票 252 平台分隔符族）。
- `./cmd/wisp/` 定向成对（lead 甲形条件③，我那枚文件里的 9 枚顶层测试）：
  台面＝同一棵树只换 panel 那一枚字节＝真·单差异对；
  before ⇄ after 都 `RUN=9 PASS_TOP=9 FAIL_TOP=0 SKIP=0` rc=0，名册作差 0／0
  ⇒ token 级改动语义零位移（0 枚新增红）。
  ★两处我的尺的错，具名更正：① v1 那对树差 3 枚（两锚之间夹别腿提交）＝台面缺陷，v1 整包读数不当成对凭据；
  ② v1 的 `-run` 串以 `/` 开头 ⇒ Go 按 `/` 切段首段为空＝匹配一切＝**根本没过滤**（RUN=154/153 是整包）。
- `./cmd/wisp/` 整包（已知真窗区整进程死亡＝旁证，非判据）：死亡点在 MARK299／panel 驻留族；
  整包名册两侧差 5 枚＝**进程死亡位置不同造成的截断差**，⛔ 读成"我跑红／跑没"。
  三枚红的定向归位复跑（两侧同 rc=1／RUN=3／PASS_TOP=0／FAIL_TOP=3／SKIP=0、逐字红因同串）：
  `TestAC2SealNoticeLandsInTheRunLegLogFile`＝`logsink_windows_test.go:175` 记录数 2 want 1（winsec 驻留族）；
  `TestCleanCheckoutBuilds_AC11`＝`panel_host_gate_test.go:404` `git rev-parse --show-toplevel: exit status 128`（导出树无 `.git`＝台面形状）；
  `TestAC14GoSideEvalPushReachesThePage`＝`panel_resident_windows_test.go:944` `title=""` want `PUSHED-33R5-OK`
  （真窗＋驻留族＝本机常驻 `wisp.exe` **PID 21892** 握全局热键＋单实例锁 ⇒ 存量，⛔ 杀那枚 PID）。
  单差异整包对（`logs/ac3-full2.txt`，同一对 ctrl 台面）＝
  before `RUN=207 PASS_TOP=150 FAIL_TOP=2 SKIP=0` rc=1 ⇄ after `RUN=203 PASS_TOP=146 FAIL_TOP=2 SKIP=0` rc=1；
  FAIL 名册两侧同两枚（上面 `TestAC2SealNotice…`／`TestCleanCheckoutBuilds_AC11`）＝**0 枚新增**；
  名册作差 `only_before=4`／`only_after=0`＝那 4 枚（AC13 三枚＋AC14AwaitedBindingReply 一枚）
  只在 before 侧留下结果行＝after 侧进程死亡更早的**截断差**（两侧末行分别是
  `TestAC14AwaitedBindingReplyReachesThePage`／`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`
  ＝死亡点在驻留窗族），⛔ 读成"我这枚字节让 4 枚测试消失"；四数逐枚点名、⛔ 相减。
- 红名册合计：**本票新增红 0 枚**（逐名都归到别票／本机常驻／台面形状）。

## AC#4 纪律格（自证）

- 写面＝三处，逐笔越界尺 0 行：
  `git show --name-only --format= <每笔> | grep -v -E '^(三处写面)'` ⇒ OVERREACH_ROWS=0（逐笔见下）。
- `git add` 与 `git commit` **两处都逐枚手写字面路径**；⛔ `add -A`／`commit -- .`／
  `$(git diff --cached --name-only)`／从 `git status` 反推／通配／`xargs`。
  ⛔ amend／reset／rebase／stash／checkout .／clean／在仓内建 worktree／push。
- ⛔ 翻任何 AC 框、⛔ 改票名／`Status:`／加 `-done`；票尾只追加一条 bullet。
- 禁区逐条未碰：`ci.yml`／`thresholds.go`／golden／`allowlist.txt`／三枚冻结件／`go.mod`／`go.sum`／
  `.gitignore`／D1–D47／C1–C32／R1–R9／D43／`internal/panel/**`／`internal/risk/**`／
  `internal/config/**`／`internal/session/**`／`frontend/**`／`design/**`。
- 中文正文含反引号／`$` ⇒ 全部 Write 成 `D:\tmp\wisp\20261018\338r1\msgs\*.txt` 再 `git commit -F`，
  ⛔ 内联 printf 双引号串。本 report 与 logs 内那枚分节号字符 0 命中、CR 行数 0（尺＝`grep -c`／`awk '/\r/'`；
  ★具名踩坑＝我上一稿在本行里**打了那个字符本身**去描述它⇒尺命中 1＝假读数，本版改为指名不写字面，
  与 336-r1 那枚同形更正同理）。

## 待定／上报（未定义即停，⛔ 自行假设）

1. panel 枚甲形已按 lead 票面裁语落地；`GOOS=linux go vet` 那枚存量红要不要另立票，归编排者拍。
2. 四枚 CRLF 假脏的"谁把工作树写成 CRLF"机制仍未归到底（票面禁区）＝本票外。
3. `338r1` 这类**合 8.3 的落点名**会让 252 族短名探针打红＝台面配方的一味，
   建议写进派单公共块（本 report 已具名给出现跑读数）。
