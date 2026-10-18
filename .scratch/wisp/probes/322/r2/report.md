# 票 322 收尾腿 322-r2 — 读数件（`cmd/wisp` 那一发＋残余两格）

> 本件＝读数件，⛔ 判语。AC 框一枚没翻（起手与交件同一把尺都是 **1 未勾／5 勾**＝AC#4）。
> 写面三处：`internal/winsec/sealhandle_privateroute_120_other_test.go`（仅注释＋一枚出向字符串的词）、
> `.scratch/wisp/probes/322/r2/**`、票 322 尾**一条**追加 bullet。⛔ push、⛔ 动 `privateFile` 产码、⛔ 动 `winsec_other.go` 一字。

## 0. 起手与合规

- 锚＝自跑 `git log -1` = `c11163c896b09e7b6863d347dadbcb65702585cc`（10-18 10:02:56 +0800）。
- **已改落点＝`D:\tmp\wisp\20261018\322r2`**（机主 10-1x 新铁律覆盖原台面配方）：树 `head.RJNBgS`
  （8.3 别名 `HEAD~1.RJN` 实测在位）、`GOTMPDIR=gtD` 先 `mkdir -p`、`mine.NkOEKC`＝同一棵树＋我那枚 `.go`。
  盘符根零污染自检见 [00-anchor.md](00-anchor.md) §1；旧 `D:\tmp\probe322r2\*` 三枚作废目录只建不删、⛔ 引用。
- 台面自证：导出树四枚 winsec 字节 sha256 ＝ 锚上 blob 同一回事；漂移尺
  `git rev-parse c11163c8:internal/winsec` = `HEAD:internal/winsec` = **`8ddcb8cb`**（读数对锚与现 HEAD 同时有效）；
  `cmd/wisp` 树在锚与 HEAD 也**同一枚** `54efbbbf`；`.go=760` 数过才信。
- 起手 `git status --porcelain -- internal cmd` = **0 行**（⛔ 本腿写这两面）。

## 1. ★头号活：`cmd/wisp` 整包那一发＝〔待编排者放行＝228-r3c 交件并核验后〕

派单 §9＋编排者 10:0x 回复①：本波 `go test ./cmd/wisp` 唯一持跑者＝`228-r3c`，放行只由编排者主动发消息。
现量（10:4x）：`228-r3c` 已交 `8adac478`（report）＋`7b1b22d8`（票尾一条），`list_agents` 仍 `running`，
我**尚未收到放行消息** ⇒ 本腿⛔ 起跑（含定向 `-run`）。编排者 10:4x 裁定⑤＝缺的这一发按 `A1021` 那形记账、
**归编排者排程、⛔ 算本腿欠账**。台件已就绪＝`cmd-wisp.ps1`（`full` ⇄ `estate` 两档；PATH 带
`third_party/sherpa-onnx`＋`build` 两串 dll、`-count=1 -timeout 30m -v`、四数逐枚点名，
`0xc0000135`／`build failed`／`no tests to run`／`setup failed` 另有单独一把尺＝命中即判台面故障、⛔ 计入判据）。
★定向 `-run` 串逐枚对过实名册（⛔ 指不存在的名字当凭据）：`Resident`＝**28 枚顶层**、`TestLive`＝**12 枚**、
`TestPanelThreadIsSTAAndExitsCleanly`＝1（`panel_resident_windows_test.go:958`）、
`TestAC13BringUpSurvivesAReusedThreadQuit`＝1（同件 `:523`）；
尺＝`git grep -h '^func Test' -- cmd/wisp | Where-Object { $_ -match <pat> } | Measure-Object`。
本腿已跑的 `cmd/wisp` 相关**只有一维**＝`go vet ./cmd/wisp/`（导出台两枚台 head⇄mine，各 **rc=0**）＝vet⛔ test，不占车道。

## 2. 残余①：过期注释更正（＝唯一那一枚 `.go`）

改动＝`sealhandle_privateroute_120_other_test.go` 现读 `:230-238`（派单写的 `:232/:237` 行号只当读数）：
`7 insertions / 7 deletions`、**净行数 0**（文件仍 **260 行**）、**0 枚 CR 字节**＝LF。
- 注释那半：`stated as unfixed at this anchor in winsec_other.go:78-82` → 带条件事实句
  （POSIX **叶子**那一支已由票 322 的 `oNoFollowPlatform` 关上、**祖先**那一支仍开、谁都⛔ 读成窗口已关）。
- 字符串那半（编排者授权，四条硬边界逐条走）：`(ticket 322, unfixed at this anchor; winsec_other.go:78-82)`
  → `(ticket 322 closed only the POSIX leaf leg; the ancestor leg stays open)`。**只改词。**
- (a) 读者尺：`git grep -l 'R-322 SHAPE OBSERVED'` 全仓 **11 枚文件**（8 份归档 `.log`＋2 份判语文本＋本件），
  仪器形（`.sh/.ps1/.py/.js/.pl`）命中 **0** ⇒ 没人数这句话，未触发停手条款。
  `git grep -n 'unfixed at this anchor' -- '*.go'` 起手恰 **2 枚**＝就是这两句；交件侧复量 **0 枚**。
- (b) 码位：全文件非 ASCII 字节现量 **0**、我新增 7 行也 0 ⇒ 纯 ASCII（⛔ `§ ≤ ✓ ⚠ ✗ →`）。`d22scan` rc=0。
- (c) `git diff -U0` 里 `Errorf|Fatalf|switch|case|want` 命中 **0 行**；`%q`＋`body` 实参、`t.Logf` 存在性、
  `case` 结构一字未动；该用例判据文字未动。
- (d) `winsec_other.go` 本体未动（复量 `19f3c211055011cc…` ＝ r1b 记录的 pristine 同串）；只把过期指针 `:78-82` 改 `:103-113`。
- 全称声明尺：新句里 `window is closed` 肯定式命中 **0**。

**证色（决定性）＝容器成对两发**，因为该件 `//go:build !windows` ⇒ 本机 Windows 档**根本不编译它**：

```
cw2-head（pristine fb6c54ed…）rc=0  RUN=142 PASS_TOP=86 FAIL_TOP=0 SKIP=0
cw2-mine （我的字节 f87035bd…）rc=0  RUN=142 PASS_TOP=86 FAIL_TOP=0 SKIP=0
分支尺（同发同行）：两台都 holds-its-own=2 / R-322-SHAPE=0 / payload-IN=0
AC#1 成对两样本（两台各自）：damagedWitnesses=0 且 leafLegRefusals=43/50（head）、25/60（mine）
```

⇒ 四数与分支色一字不变＝我那一枚 `.go` 对判据**零影响、且没把任何判据放宽**。
⚠ 那一枚 `R-322 SHAPE OBSERVED` 分支**在已修 POSIX 码上今天不再走**（`cw-winsec` 两样本 2/2 走 `:225`
`witness holds its own bytes`）＝与 `120-v2-verdict` §交付缺陷 2 那句"本腿两发盲跑里那一形都没出现"相符。

## 3. 本机 `internal/winsec` 取证（estate 色取定向 `-run`，⛔ 拿整包当全绿）

成对两枚临时根台面（D 台 `gtD` 预建 ⇄ C 台＝默认 `TMP=C:\Users\swq\AppData\Local\Temp`）：

```
wh-D 整包 -count=1 -v  rc=1  RUN=110 PASS_TOP=59 FAIL_TOP=2 SKIP=0
wh-C 整包 -count=1 -v  rc=0  RUN=110 PASS_TOP=61 FAIL_TOP=0 SKIP=0
wd-D 定向 estate 两枚  rc=1  RUN=5  PASS_TOP=0  FAIL_TOP=2 SKIP=0
wd-C 定向 estate 两枚  rc=0  RUN=5  PASS_TOP=2  FAIL_TOP=0 SKIP=0
wt-D 定向 TestAC3Windows rc=0 RUN=2 PASS_TOP=2 FAIL_TOP=0 SKIP=0
```

两枚红＝`TestGateJudgesThePrivateSetByResolvedSID`、`TestSealNarrowsAndNamesThePrincipalItRemovedBySID`
（逐字红句 `private_set_sid_windows_test.go:314: planted as strangers: [S-1-1-0 S-1-5-32-546 …-500]; …`＋`:193/:197` 同族），
**D 台红／C 台绿随卷确定性换色＝票 310/311 那一族既有形状**，两枚用例一字未动、⛔ 计入本票、归属交编排者。
`walk=` 命中具名：`wh-D` 9 枚／`wh-C` 7 枚**全部**是用例名与 WARN 正文里的 `…WalkSkipsLinks`／`…IsNotWalked` 子串，
⛔ 台面故障；`build failed`／`no tests to run`／`setup failed`／`0xc0000135` 全 0。

## 4. 残余②：Windows `witnessTruncated=true` 那半格——只读复认（⛔ 为它动任何字节）

`wt-D`／`wh-D`／`wh-C` 三发逐字复现 **`witnessTruncated=true witnessSize=35`（3/3，加首发共 4/4）**：

```
openleaf_guard_322_windows_test.go:251: AC#3 local paired reading … attempts=400 flips=341
guardRefusals=27 otherRefusals=260 returnedNil=113 witnessTruncated=true witnessSize=35
witnessMode=-rw-rw-rw- firstGuardVerbatim="winsec: … the name carries a reparse point, so the open that
follows would go through that link and truncate whatever it targets …"
```

- 守卫**在产线路由上答话**（`guardRefusals>0`＋`flips>0`＝仪器行自检绿，三发都 PASS）；
- `=35`＝外来见证被换上 raced payload（`"raced private payload 322r1 windows"` 恰 35 字节）＝`322-v1` §AC#3 具名的
  **票 120 Windows 半格外射程**（`sealHandle` 无同一对象闸门，`winsec_windows.go` 那格 `120-v2` 缺陷 3 在册）；
- `witnessSize=0` 那一形本腿**未再现**＝按 `322-r1b` §5.4／`322-v1` §4 已登记形状复认，⛔ 更正、⛔ 动字节。

## 5. AC#4 消费者面成对读数（Linux 档 ⇄ 本机；五包＝risk/secret/config/agent/memory）

```
容器 cc2-head（pristine）  rc=1  RUN=1160 PASS_TOP=702 FAIL_TOP=6 SKIP=6
容器 cc2-mine （我的字节）  rc=1  RUN=1160 PASS_TOP=702 FAIL_TOP=6 SKIP=6   ← 逐字相同
容器 cc-consumers（首发）    rc=1  RUN=1160 PASS_TOP=702 FAIL_TOP=6 SKIP=6
本机 hc-*（Windows 档 -count=1）：
  risk   rc=1  RUN=230 PASS_TOP=139 FAIL_TOP=1 SKIP=1   FAIL	github.com/…/internal/risk
  secret rc=0  RUN=29  PASS_TOP=21  FAIL_TOP=0 SKIP=0   ok
  config rc=0  RUN=185 PASS_TOP=104 FAIL_TOP=0 SKIP=0   ok
  agent  rc=0  RUN=114 PASS_TOP=90  FAIL_TOP=0 SKIP=0   ok
  memory rc=0  RUN=80  PASS_TOP=43  FAIL_TOP=0 SKIP=2   ok
```

红名册**两台同三枚、全在 `internal/risk`**：`TestAnchorMissReasonNamesCandidateAndMissingAnchor`、
`TestAnchorMissKeepsNonOverridableClassAndNeedL2`（＝票 334 `anchor_miss_reason_334_test.go`，10-17 11:47 落在我锚**之前**），
`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（红句逐字 `internal/panel/composer.go:127 … display_scrub_216.go:26/:37/:218
reads Result.Canonical without reading the rewrite account … ticket 102's hole re-opened by a new consumer`＝**台账 A380／200 已在册**）。
⇒ ⚠ **具名上报一处与既有判语不符**：`322-r1b` §6 写"Linux 档五包全 ok"，本腿现量 `internal/risk` 在 Linux 档 **rc=1**。
本票⛔ 一字节未动 `internal/risk`／`internal/panel` ⇒ 那一支变色非本票引入（红名册两台一字不变背书）；归编排者。
★编排者 10:4x 裁定③收单：他**自己建容器台件复跑**那两枚（⛔ 抄本腿的数）＝红句逐字 `anchor_miss_reason_334_test.go:131/:156`
`fixture guard: the fail-closed shape must still be A here, got class=B …`＋`:193/:197/:201` 三句、rc=1 ⇒ 判"成立、
**且不是 322 造成的**＝票 334 自己的新仪器在 Linux 档红"，已立**票 339**
（`.scratch/wisp/issues/339-ticket334-anchor-miss-instrument-runs-on-linux-but-its-fake-home-is-windows-only.md`），
⛔ 记在 322 头上；本腿⛔ 动 `internal/risk/**`（那是 336-r1 写面）。
中途一枚多出的 `Test285LoopDispatchesOneCorrPerCall`（328 的件）按单臂 `-count=3` 两台复跑**都 ok rc=0**＝一发摆动、⛔ 计入判据。

## 6. 门禁（逐枚 rc=；格式化/d22scan/长度门＝共享工作树，build/vet＝仓外导出台）

`gofumpt.exe -l`（真身 `/d/work/base/gopath/bin/`）对我的件 **rc=0 列出 0 枚**；`gofmt -l` 同件 **rc=0**；
`sh scripts/d22scan.sh` **rc=0** clean（ban #8 `design/=39 frontend/=116 internal/=586 cmd/=143`）；
长度门 `--with-self-test` **rc=0** VERDICT GREEN（over=60/roster=60）；
`go build ./...` 两台各 **rc=0**；`go vet ./internal/winsec/` 本机＋linux＋darwin＋windows × 两台 ＋
`go vet ./cmd/wisp/` × 两台 ＝ **十二发全 rc=0**。
★存量假脏具名（⛔ 顺手重排别腿文件）：`gofumpt -l internal cmd` 列出 **6 枚**＝
`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`、`internal/risk/provenance.go`、
`internal/tools/task_output_pointer_notice_errtext_174r7_test.go`、`cmd/wisp/models.go`、`cmd/wisp/panel_transport_35r2_test.go`
（台账 A1024 那笔＝6 枚现量；票 338 已立＝真脏 2 枚＋CRLF 假脏 4 枚）。

## 7. 本腿自己的偏差与台面故障（自抓，⛔ 计入判据）

1. **★一笔形状违规（越界 1 枚路径）＝编排者已收单、不追责、⛔ 回改历史**：`7f805583` 里
   `git commit -F msg -- $(git diff --cached --name-only)`，那个 `$(…)` 取自**共享 index**，把 `leg-328-v1`
   已 `git add` 未 commit 的 `.scratch/wisp/probes/328/v1/logs/start-rulers-v2.txt` 卷进我这一笔（`--numstat 8/0` 只追加）。
   尺：blob 演进 `ce16d7bf`=`b6cae95c` → `7f805583^`=`b6cae95c` → `7f805583`=`9d1ac93e`＝工作树现量
   （`git hash-object`）⇒ **那 8 行逐字是 328-v1 自己写的，零枚我的字节；内容没丢、没改坏，只是代提在我名下**。
   定式（编排者入台账）＝`git add` 逐枚但 `commit` 从 index 取名＝**只守了一半**；本腿其后每一笔 pathspec
   **字面逐枚写死**，⛔ 再用任何"从 index 反推"的形状；⛔ amend/reset/rebase/回改历史（＝`A1011` 第⑧节同族）。
2. `host-winsec.sh` 首版把仓路径写成 MSYS 形 ⇒ `git -C` 当场 rc=128 零读数＝台面故障，件标〔故障件〕只建不删；
   有效版＝`host-winsec.ps1`。同坑亦见 `container.ps1` 用 `D:/` 形挂载路径。
3. `gates.ps1` 前两版把 `mine` 树路径归一坏成 `D:\d\tmp\…`（不存在）⇒ 那一版 **6 枚 `*-mine` 发 Push-Location
   报错、`rc=0` 是无效读数**；已具名作废，第三版加 `Test-Path` 硬闸（不在即 `exit 91`）后重跑＝§6 那批数。
4. `host-consumers.ps1` 首版 `$p.Trim('/')` 被开头的 `.` 挡住 ⇒ tag=`./risk.`、正文写进不存在的 `hc-\risk.log`、
   五包**零读数**＝台面故障、⛔ 计入判据（作废 summary 只建不删）；本版显式配对 pkg/tag。
5. 起手台件曾建在 `D:\tmp\probe322r2\`（`D:\tmp` 下、不在盘符根），新铁律到达后全部作废、改用具名落点；⛔ 删。

## 8. 落点索引

- 本件＝`.scratch/wisp/probes/322/r2/report.md`；起手件＝`00-anchor.md`；台件＝`mk-tree.sh`（＋标〔故障件〕那枚 `.sh`）、
  `host-winsec.ps1`、`container.ps1`、`container-mine.ps1`、`gates.ps1`、`host-consumers.ps1`、`cmd-wisp.ps1`。
- `logs/`＝逐枚原始正文＋六枚 summary（`host-winsec-`／`gates-`＋`gate-*.txt`／`container-`／`container-mine-`／
  `host-consumers-`），正文件＝`wh-*`／`wd-*`／`wt-*`、`cw-winsec`、`cc-*`、`cc2-*`、`cw2-*`、`cmdwisp-*`。
- 只建不删的树＝`D:\tmp\wisp\20261018\322r2\{head.RJNBgS,mine.NkOEKC,gtD}`＋作废的 `D:\tmp\probe322r2\*`。
- 产码面＝那一枚 `.go` 的 7 行注释＋2 行字符串（`ab49e9f9`）；票面＝尾一条追加 bullet；⛔ push。
- **必答 (i)(ii)(iii) 另立一件**＝[report-sec9.md](report-sec9.md)（41 行，同写面），本件到 160 行帽为止。
