# 票 322 — 非实现者终裁腿 322-v2 判语件（只裁 AC#4：门禁＋消费者不推红）

> 身份：裁决者≠实现者。产码＝`322-r1`/`322-r1b`；上轮判语＝`322-v1`；补齐读数＝`322-r2`（交料⛔ 判语）。
> 本腿不是以上任何一程。**六枚 AC 框一枚没翻**（起手尺 `- [ ]`=1/`- [x]`=5，`logs/00-start-rulers.txt:4`；交件侧复量见自证）。
> ⛔ push／⛔ 改票名／Status／`-done`／⛔ 动任何 `.go`（本腿全部提交 `.go` 命中 0）。
> 起手锚＝本腿 `git log -1` 现认 `aaaedbb4`（12:22:13 +0800，本腿第一发现跑）；取证期间 HEAD 在漂，各读数件各自带时刻
> （13:1x 终量侧＝`925481ea`/`fd009fa1`）；读数取自【具名提交树】与【已入库原始 log 逐枚复数＋blob 对撞（14/14 MATCH，`logs/00-start-rulers.txt`）】，
> ⛔ 抄腿的数、⛔ 抄编排者的数（编排者的 448/322/8/3 与 448/323/7/3 一律当未验证，本腿自己从 `cwisp-full-D/C.log` 重数坐实＝逐字同，见第 1 节）。

## 0. AC#4 总判语

**成立**（三件必答口径：estate＝以定向 `-run` 为准＋整包两台背书；winlive 档＝**不进**本格数、登记为环境存量＋CI 无 test 档读数；
estate 缺 C 台那一发＝**不必补**，整包两台已覆盖该维）。附三条具名残余（第 6 节），都不把本格翻成"部分成立"：
它们是〔别票在册／已登记形状／台面摆动〕，⛔ 一形属"本票把门禁或消费者推红"。

## 1. 四数复数（本腿自己现跑，尺＝^=== RUN ⇄ 顶格 ^--- PASS/FAIL/SKIP，⛔ 混用⛔ 相减；逐枚指名单文件⛔ 通配）

| 发（指名原始 log） | RUN | PASS_TOP | FAIL_TOP | SKIP | 故障尺(135/bf/ntr/sf) | 来源 |
|---|---|---|---|---|---|---|
| cwisp-full-D.log | 448 | 322 | 8 | 3 | 0/0/0/0 | r2 入库·blob=HEAD 对撞 MATCH·本腿重数 |
| cwisp-full-C.log | 448 | 323 | 7 | 3 | 0/0/0/0 | 同上 |
| cwisp-estate-resident.log | 31 | 22 | 0 | 0 | 0/0/0/0 | 同上 |
| cwisp-estate-deathpoints.log | 2 | 2 | 0 | 0 | 0/0/0/0 | 同上 |
| cwisp-estate-live-winelive.log | 12 | 6 | 6 | 0 | 0/0/0/0 | 同上 |
| cwisp-estate-live.log＝台面故障发 | 0 | 0 | 0 | 0 | ntr=2 ⇒ ⛔ 计入判据 | 同上 |
| cw2-head/mine.body.log 各 | 142 | 86 | 0 | 0 | 0/0/0/0 ×2 | r2 入库·blob MATCH |
| cc2-head/mine.body.log 各 | 1160 | 702 | 6 | 6 | 0/0/0/0 ×2 | r2 入库·blob MATCH |
| orch-322v1-ac4-cmdwisp-full.txt | 430 | 305 | 6 | 3 | 0/0/0/0 | 编排者 10-16 件·本腿自数 |
| v2head.body.log（**本腿现跑**） | 142 | 86 | 0 | 0 | 0/0/0/0 | 本腿导出台 tree.v2B(c11163c8) |
| v2mine.body.log（**本腿现跑**） | 142 | 86 | 0 | 0 | 0/0/0/0 | 本腿导出台 tree.v2A(HEAD 含 r2 字节) |
| v2winsec.body.log（**本腿现跑**） | 142 | 86 | 0 | 0 | 0/0/0/0 | 同 v2A |
| v2cons.body.log（**本腿现跑**） | 1172 | 714 | 6 | 6 | 0/0/0/0 | Linux 五包·HEAD 档 |

明细与逐枚 sha＝`logs/recount.txt`。`cmd/wisp` 整包那两发本腿**未重跑**（写腿 `338-r1` 12:22 后仍可能占 `go test ./cmd/wisp` 车道，
本腿起收两次 tasklist 均无 go.exe＝当时车道空，但⛔ 押注别腿收工时刻）⇒ **标"待车道"复跑发**；其凭据＝入库原始 log＋逐枚 blob 对撞，有效。

## 2. 红名册逐名归因（8 枚整包红逐名，⛔ 只看枚数；作差尺本腿独立跑＝`logs/roster-diff.txt`＋`roster/*.txt`）

10-16 名册 6 枚**全部 ⊂ D 台 8 枚**（in-O-not-D＝0 枚）；D∪C ⊖ O＝2 枚，逐名归因：
- `TestCleanCheckoutBuilds_AC11`＝**导出台无 `.git` 装置红**：D 台红句 `panel_host_gate_test.go:404: git rev-parse --show-toplevel: exit status 128`；
  同一枚在 10-16（共享树台面）＝`--- PASS`（该 log `:958` 本腿逐字验）⇒ 台面形状所致、非产品字节。
- `TestAC2SealNoticeLandsInTheRunLegLogFile`＝**随台面摆动**：两台差集唯一一枚（D 红 C 绿；红句 `logsink_windows_test.go:175: records … = 2, want 1`）；
  本腿独立尺＝全 `.scratch` 顶格 FAIL 句命中 24 枚件（`logs/sealnum-files.txt`）、最早入库 10-13（`186/r1`）⇄ 10-14（`133/r1`），`228/r3b` surface c/d 两台都红（`:865/:864`）
  ＝别票台面基线既有形状，全早于 322 产码 `8baec646`（10-16 09:06）。
其余六名（258×3／AC14／AC4Focus／AC246）＝⊂ 10-16 名册的别票存量；`AC246` 红句自陈 `resident_task_source_246_windows_test.go:434:
AC#7 RED (desktop state, not our code): a dev Wisp is already running`＝本机常驻实例 PID 21892（现量 StartTime 2026-10-15 19:03:14，起于所有腿之前）。
⇒ **本票新增红＝零枚 站得住**（尺＝两台并集 ⊖ 10-16 名册 ＝ {装置枚, 摆动枚}，⛔ 属本票字节；必答 (ii) 同文）。

## 3. 两台成对那一维（派单③）

整包 D⇄C 两台在盘（第 1 节）；estate 三档只有 D 台。本腿作名册尺（非推理，`logs/estate-dC-rulers.txt`）：定向 resident 22 枚判定名**全部**出现在
full-D **与** full-C 的顶格判定名（`comm` 差＝0/0，量的）；
deathpoints 两枚在两台各有判定行；resident 族两台的 `=== RUN`＝31/判定＝22/FAIL＝0 一字不变 ⇒
**"estate 的色随临时根台面摆动"这一维已被整包两台覆盖，缺的 estate C 台一发不改变任何判语**⇒ 不判部分成立／不补发（必答(v)）。

## 4. 因果尺（复认 `322-r2` 第 4 节，本腿在导出树上独立现跑 `go list -e -f`，`logs/causal-golist.txt`）

`sealhandle_privateroute_120_other_test.go` 首行 `//go:build !windows`；GOOS=windows：该件进 `IgnoredGoFiles`(1)⇄`TestGoFiles`(0)；
GOOS=linux 反之(0/1)。`cmd/wisp` TestGoFiles＝89（windows）/49（linux）、其中来自 winsec 的＝0；`cmd/wisp` 的 Imports⇄TestImports
含 `internal/winsec`＝0/0（前导空格假零已更正复跑）。⇒ 复认成立：那枚 `.go`（注释＋一句出向字符串）**机制上改不了 `cmd/wisp` 颜色**。

## 5. 门禁四数＋消费者（本腿现跑，共享树先 `unset GOTMPDIR TMPDIR`＝票 311 shape B GUARD 定式）

`logs/gates-shared.txt`：gofumpt `-l internal/winsec/` **rc=0 列 0**；全仓形 `-l internal cmd` 列出 4＝票 338 存量假脏族（A1024 形，⛔ 归本票）；
`sh scripts/d22scan.sh` **rc=0** clean（ban#8 internal=590/cmd=143）；长度门 `--with-self-test` **rc=0** VERDICT GREEN；
`go vet ./internal/winsec/`⇄GOOS=darwin⇄GOOS=linux⇄`go vet ./cmd/wisp/` 四发 **rc=0**。
容器成对 142/86/0/0：r2 两台＋本腿 fresh 两台（第 1 节，共 4 发一字不变）。消费者：本腿 fresh Linux 五包 v2cons 六枚顶格 FAIL＝三枚唯一名×2 样本
（`anchor_miss_reason_334_test.go` 两枚＝票 339 在册；`TestC26RewriteAccount…`＝台账 A380／历史 `200/r2`·`216/*` 件在册）；
RUN 1160→1172 的 12 行增量＝227-r1 新落 config 件（`logs/cons-delta.txt`），红名册一字不变 ⇒ **risk 变色非本票**。
本机五包 hc-* 复认＝risk 顶格 FAIL 1（C26 同一枚）、其余四包 0 红（r2 入库件 blob MATCH）。

## 6. 残余与具名上报（都⛔ 阻断 AC#4）

① `witnessTruncated=true witnessSize=35` 复认成立：本腿尺＝`wt-D/wh-D/wh-C` 各 1 枚 trunc=true、size=35 各 1、`witnessSize=0` **0 枚**
（r2 报的 4/4 复现其第 4 发在腿正文件、非入库 log；本腿按可核 3/3 复认＋=0 形未再现＝按已登记形状成立；属票 120 Windows 半格外射程）。
② `cmd/wisp` 整包本腿未重跑＝待车道（第 1 节末），凭据＝入库原始 log；`cmd/wisp` 树 12:22 起漂移（338-r1，测试件重排、`diff -w` 仍 24/10＝非纯空白）
⇒ 编排者若要"现 HEAD 整包色"须自排车道复跑，⛔ 用 r2 发冒充现锚。
③ winlive 六枚红＝环境存量＋CI 无 test 档读数（必答(iv)）；三档 estate 只有 D 台（必答(v)＝不补）。

## 7. 六格状态（本格之外⛔ 重裁）

`AC#0/#1/#2/#3/#5`＝照 `322-v1` 五格成立判语维持（本腿⛔ 重裁、⛔ 推翻凭据未见）；**`AC#4`＝本件改判：部分成立 → 成立**
（缺的那发已由 `322-r2` 补料、本腿全部复数＋fresh 现跑坐实）。翻框权归编排者。

## 8. 落点

判语件＝本件（必答 (i)–(v) 全文＝具名分件 `.scratch/wisp/probes/322/v2/mandatory.md`，84 行，本件 第 1 到第 5 节 为其尺的落点索引）；
读数件＝`.scratch/wisp/probes/322/v2/`（`00-start-rulers.txt`／`recount.sh+logs/recount.txt`／`roster-diff.sh+logs/roster-diff.txt+roster/`／
`causal-golist.txt`／`gates-shared.txt`＋`gates-vet.txt`（首发 vet 三档，并入终量、件头具名）／`pair-fresh.sh+logs/pair-fresh-summary.txt+v2head/v2mine.body.log`／`container.sh+logs/container-summary.txt+
v2winsec/v2cons.body.log+logs/v2cons-delta.txt+logs/roster/{run,fail,skip}-*.txt+new-tests.txt`／`estate-dC-rulers.sh+logs/estate-dC-rulers.txt`／`roster/{o,d,c,union-dc,verd-d,verd-c,est-res,est-dp,skip-*}.txt`／`logs/sealnum-files.txt`／
`logs/container-summary.FAULT-pathconv.txt`＝台面故障件只建不删⛔ 计入判据；`logs/96-selfcheck-v2.txt`＝交件侧三把尺改正版；`msg/`＝逐笔提交正文入库；
具名**不入库**两枚＝`logs/95-selfcheck.txt`（首版自尺：拼尺行笔误"七枚 g: command not found"＋把分节号字面写进尺文件＝描述它的字占了一枚，
336-r1 同坑再犯、当场改正＝needle 改字节构造；件在盘只建不删⛔ 计入任何尺）与 `ticket-before.txt`（票面追加前快照＝原字节遗产、
内含编排者原文枚里的分节号字符，⛔ 属我的新字节））；
导出树＝`D:\tmp\wisp\20261018\322v2\{tree.v2A,tree.v2B,gtD}`（盘符根零污染，`tree.v2A.provenance.txt`）。票尾一条 bullet。⛔ push。
