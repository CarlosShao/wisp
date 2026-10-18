# 票 322 收尾腿 322-r2 — 读数件（`cmd/wisp` 那一发＋残余两格）

> 本件＝读数件，⛔ 判语。AC 框一枚没翻（起手与交件同一把尺都是 **1 未勾／5 勾**＝AC#4）。
> 写面三处：`internal/winsec/sealhandle_privateroute_120_other_test.go`（仅注释＋一枚出向字符串的词）、
> `.scratch/wisp/probes/322/r2/**`、票 322 尾**一条**追加 bullet。⛔ push、⛔ 动 `privateFile` 产码、⛔ 动 `winsec_other.go` 一字。

## 0. 起手与合规

- 锚＝自跑 `git log -1` = `c11163c896b09e7b6863d347dadbcb65702585cc`（10-18 10:02:56 +0800）。
- **已改落点＝`D:\tmp\wisp\20261018\322r2`**（机主 10-1x 新铁律覆盖原台面配方）：树 `head.RJNBgS`
  （8.3 别名 `HEAD~1.RJN` 实测在位）、`GOTMPDIR=gtD` 先 `mkdir -p`、`mine.NkOEKC`＝同一棵树＋我那枚 `.go`。
- 台面自证：导出树四枚 winsec 字节 sha256 ＝ 锚上 blob 同一回事；漂移尺
  `git rev-parse c11163c8:internal/winsec` = `HEAD:internal/winsec` = **`8ddcb8cb`**（读数对锚与现 HEAD 同时有效）；
  `cmd/wisp` 树在锚与 HEAD 也**同一枚** `54efbbbf`；`.go=760` 数过才信。盘符根零污染自检与起手台件作废见 [00-anchor.md](00-anchor.md) §1。
- 起手 `git status --porcelain -- internal cmd` = **0 行**（⛔ 本腿写这两面）。

## 1. ★头号活：`cmd/wisp` 那一发＝**已跑完**（编排者 11:0x 放行后起跑）

**四数逐枚表、逐名归位、死亡尺、build-tag 因果尺全在 [cmdwisp-readings.md](cmdwisp-readings.md)（138 行）；正文＝`logs/cmdwisp-*`。**
本节只留最硬的几行（裁定②⒜已拍＝**AC#4 以定向 `-run` 的 estate 色为准，整包色只作台面终态记录**）：

- 起跑前照裁定①注意①复跑漂移尺＝`cmd/wisp` 树 `c11163c8`＝现 `HEAD`＝**`54efbbbf`** 同一枚、导出树⇄现 HEAD
  **139 枚文件差 0 行** ⇒ 旧树仍等于现锚字节（⛔ 拿旧树冒充现锚）。
- 整包成对两台：`cwisp-full-D` rc=1 `448/322/8/3`（D 卷 GOTMPDIR，483.7s）⇄ `cwisp-full-C` rc=1 `448/323/7/3`
  （默认 TMP＝C 卷，445.6s）；差集只一枚 `TestAC2SealNoticeLandsInTheRunLegLogFile`（D 红 C 绿）、SKIP 逐字相同。
- estate：`estate-resident` rc=0 `31/22/0/0`；`estate-deathpoints` rc=0 `2/2/0/0`＝228-r3b 具名那两枚死亡点
  定向**都绿**；`estate-live` 第一发 `0/0/0/0`＋`[no tests to run]`＝**台面故障⛔ 计入判据**（根因＝`TestLive*`
  六件首行全 `//go:build windows && winlive`，我台件⛔ 带 `-tags`）⇒ 补正 `estate-live-winelive` rc=1 `12/6/6/0`。
- **八枚整包红逐名归位＝本票新增红零枚**；**这一回没有整进程死亡**（`=== RUN` 448＝顶格 333＋子测试 115＝判定行
  448＝顶格 333＋缩进 115，未判定 0 枚）。★机制级因果尺＝我那枚 `.go` 是 `//go:build !windows` 件（GOOS=windows
  进 `IgnoredGoFiles` 1 枚、`TestGoFiles` 0），且 `_test.go` 只编进自己包（`cmd/wisp` 89 枚测试件里来自
  `internal/winsec` 的 **0 枚**）⇒ 本腿那枚字节**机制上无法**改 `cmd/wisp` 的颜色。

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

**证色（决定性）＝容器成对两发**（该件 `//go:build !windows` ⇒ 本机 Windows 档根本不编译它）：
`cw2-head`（pristine `fb6c54ed…`）⇄ `cw2-mine`（我的字节 `f87035bd…`）各 **`RUN=142 PASS_TOP=86 FAIL_TOP=0 SKIP=0` rc=0
一字不变**；分支尺两台都 `holds-its-own=2 / R-322-SHAPE=0 / payload-IN=0`；AC#1 成对两样本各自
`damagedWitnesses=0` 且 `leafLegRefusals=43/50`（head）、`25/60`（mine）
⇒ 我那一枚 `.go` 对判据**零影响、且没把任何判据放宽**。
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
⇒ ⚠ **具名上报一处与既有判语不符**：`322-r1b` §6 写"Linux 档五包全 ok"，本腿现量 `internal/risk` 在 Linux 档 **rc=1**；
本票⛔ 一字节未动 `internal/risk`／`internal/panel` ⇒ 变色非本票引入（红名册两台一字不变背书）。
★编排者 10:4x 裁定③收单＝他自己建容器台件复跑那两枚（⛔ 抄本腿的数）＝红句 `anchor_miss_reason_334_test.go:131/:156`
`fixture guard: the fail-closed shape must still be A here, got class=B …`＋`:193/:197/:201` 三句、rc=1 ⇒ 判"成立、
**且不是 322 造成的**＝票 334 自己的新仪器在 Linux 档红"，已立**票 339**（`issues/339-ticket334-anchor-miss-instrument-runs-on-linux-but-its-fake-home-is-windows-only.md`），
⛔ 记在 322 头上；本腿⛔ 动 `internal/risk/**`（＝336-r1 写面）。中途多出的 `Test285LoopDispatchesOneCorrPerCall`
（328 的件）按单臂 `-count=3` 两台复跑**都 ok rc=0**＝一发摆动、⛔ 计入判据。

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

## 7. 本腿自己的偏差与台面故障（自抓，⛔ 计入判据；明细＝`cmdwisp-readings.md` §6/§7、`00-anchor.md` §3）

1. **★同一笔形状违规犯了两次**＝`git add` 逐枚字面 pathspec 合规，但 `git commit` 的 pathspec 写成从**共享 index**
   反推的形状。`7f805583` 越界 1 枚（`328-v1` 的 `start-rulers-v2.txt`，8/0 只追加、blob `b6cae95c`→`9d1ac93e`
   ＝工作树现量）＝编排者裁定①收单；`9c372800` 越界 2 枚（编排者自己的 `probes/orch/20261018/txt/{commit-40em.sh,
   commit-corrections.sh}`，两枚在父笔不存在、我笔 blob＝工作树逐字相同）＝⚠ 我那笔正文刚写着"两处都字面写死"，
   那句话当场是假的。**定式升级＝`commit` 的 pathspec 逐枚手打完，两处都⛔ 来自任何 `git diff --cached`／
   `git status` 输出**；`17e67cca` 起照此执行（复尺 out=0）。⛔ amend/reset/rebase/回改历史（＝`A1011` 第⑧节同族）。
2. 四枚台面故障（全部零读数、⛔ 计入判据、件与正文只建不删留底）＝`host-winsec.sh` 用 MSYS 路径形 ⇒ `git -C` rc=128；
   `gates.ps1` 前两版把 `mine` 树归一坏成 `D:\d\tmp\…` ⇒ 6 枚 `*-mine` 的 `rc=0` **作废**（第三版加 `Test-Path`
   硬闸 `exit 91` 后重跑＝§6 那批数）；`host-consumers.ps1` 首版 `Trim('/')` 被开头的 `.` 挡住 ⇒ 五包零读数；
   `cmd-wisp.ps1` 前两版函数取名 `Go` 撞上外部命令 ⇒ 递归调自己、五发零读数（`0xc0000135=2` 是自指产物，两串 dll
   现量都在树里），第三版改名 `Run-Shot`＋`& 'go'`＋`ParseFile` 自检。另：我用 `Set-Content -Encoding ascii` 改那枚
   台件时把它自己的 CJK 注释打成 `?`（已 UTF-8 重写）。
3. `e73b8c35` **误复用了 `msg-c5.txt`** ⇒ 那一笔正文写着"票尾一条 bullet"、与内容（台件改进）不符＝描述面缺陷
   （name-only 尺合格＝1 枚在我写面内）；真实内容以本节＋`cmd-wisp.ps1` 文件头为准。
4. 起手台件曾建在 `D:\tmp\probe322r2\`（在 `D:\tmp` 下、⛔ 盘符根），新铁律到达后全部作废、改用具名落点；⛔ 删。
5. `report.md` 多次越过 160 行帽（逐次现量最高读到 **172**）⇒ 每一轮都当场压回（本件终态＝**160 行＝正好顶帽**）；
   压缩只动叙述与重复段（明细改指 `cmdwisp-readings.md`），⛔ 改任何读数——复尺＝`448/322/8/3`、`448/323/7/3`、
   `31/22/0/0`、`2/2/0/0`、`12/6/6/0`、`142/86/0/0` 六枚逐名仍在位。★另一次公共块 §shell 我自己犯了：量行数的
   命令把反引号写进双引号串 ⇒ PowerShell ParserError、那一发零输出＝自己的台面故障（原文＝⛔ 双引号串里出现反引号或 `$`）。

## 8. 落点索引

- 正文件四枚＝`00-anchor.md`(45)／`report.md`(本件 160)／[report-sec9.md](report-sec9.md)(41，必答)／
  [cmdwisp-readings.md](cmdwisp-readings.md)(138，头号活明细)；台件七枚＝`mk-tree.sh`、`host-winsec.ps1`、
  `container.ps1`、`container-mine.ps1`、`gates.ps1`、`host-consumers.ps1`、`cmd-wisp.ps1`（＋〔故障件〕那枚 `.sh`）。
- `logs/`＝逐枚原始正文＋summary（`host-winsec-`／`gates-`＋`gate-*.txt`／`container-`／`container-mine-`／
  `host-consumers-`／`cmdwisp-`／`hc-*`），故障留底＝`cwisp-*.FAULT-recursion.log`／`cmdwisp-summary.FAULT-recursion.txt`／`test.FAULT-recursion.log`。
- 只建不删的树＝`D:\tmp\wisp\20261018\322r2\{head.RJNBgS,mine.NkOEKC,gtD}`＋作废的 `D:\tmp\probe322r2\*`；产码面＝那一枚 `.go` 的 7 行注释＋2 行字符串（`ab49e9f9`）；票面＝尾一条追加 bullet；⛔ push。