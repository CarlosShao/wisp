# 票 288 — 票 280 裁完归属之后，那一处"两把尺同声、2 行、已在 CI 名册里"的对齐**成了一张空账**

**立票**：2026-10-09 11:5x 编排者（机主令「及时补票」；来路＝票 280 自己的 `Status` 句：「`panel_inbound_guards_35r3_test.go` 那一处对齐「该改」这件事现在是一张空账……改它要么随下一枚真动 `cmd/wisp/**` 的腿捎带、要么另立落地票」——**本票就是那枚"另立"**）
**性质**：一行注释对齐。**低利害、完全可逆、零功能变化**——⛔ 不上机主清单，归工程内务。

## 现量（引用前先重跑；数字是快照）

- 件＝`cmd/wisp/panel_inbound_guards_35r3_test.go`，差异＝**1 个 hunk／2 行**，形状＝`const` 组里尾随 `//` 注释的**对齐**。
- **两把尺同声**（票 280 `AC#1` 的编排者独立复跑）：裸 `gofmt -d` 14 行 ＝ `gofumpt -d` 14 行，同一个 hunk `@@ -63,8`。⇒ 它不属于"gofumpt 独有加严"那一族，改它**同时** satisfy 两把尺。
- 对照（同批另一枚，⛔ 不在本票射程）：`cmd/wisp/panel_transport_35r2_test.go` ＝ `gofmt` 14 ／ `gofumpt` **101**（gofumpt 独有并组／每行一枚 composite literal／顶层声明补空行三类）⇒ 票 280 判它"留给下一枚真动该文件的票"，**本票不搭这一枚的车**（只挑 gofmt 那 1 处改＝CI 照旧红、要改两遍）。
- ⚠ **CI 那扇门今天已因 `.scratch/**` 里的故意坏样本常红**（票 111／269／276 三处既有账），所以"会进 CI 名册"⛔ 读不成"拦得住"——本票改的是**债本身**，不是让门变绿。

## 要建什么

- [ ] **AC#1 只改那一处对齐**：写面**只许** `cmd/wisp/panel_inbound_guards_35r3_test.go`；尺＝改前后各跑 `gofmt -d <该件>` 两列并排（前＝14 行、后＝**0 行**），且 `gofumpt -d <该件>` 也降到 0（它本就是 14＝同声）。⛔ `-w` 全仓禁用、⛔ 不许批量重排、⛔ 不许碰第二枚文件。
- [ ] **AC#2 证明没伤到测试**：`go test ./cmd/wisp/ -count=1` 改前改后各一次，报 **PASS/FAIL/SKIP 三数＋红名册差集**（⚠ 本机跑 `cmd/wisp` 需要 sherpa/onnxruntime 原生 DLL 的 harness：`GOMODCACHE` 三枚 DLL 拷到测试 exe 旁＋CWD＝包目录；**不带 DLL 时那枚 0xc0000135 是环境红、不是被验物**，基线一律具名写清用了哪支）。⛔ 不许因为"整包本来就红"就不跑改前基线。
- [ ] **AC#3 门禁与越界**：`GOFLAGS= go build ./...` rc=0；`sh scripts/d22scan.sh` 纯净树 rc=0；`gofmt -l`／`gofumpt -l` 对动过的件空；`git show --stat` 名册**只含该件**；`frontend/**`／`design/**`／三枚冻结件／`testdata/**`／golden／`thresholds.go`／`allowlist.txt` 零字节；⛔ 零 push、commit 必带显式 pathspec。

## 禁区

⛔ 不许顺手改 `panel_transport_35r2_test.go`（那是上面具名"留给下一枚真动它的票"那一枚）；⛔ 不许顺手改任何注释文字／`ci.yml` 的 lint 步／`.scratch/**` 里的坏样本；⛔ 不许把"让 CI 那步变绿"当成本票的判据（今天它拦不住）。

**Status:** **未开工**。排程＝随下一枚真动 `cmd/wisp/**` 的写腿**同批捎带**最省（单独一发要独占 Go 面、不值）；若没有那样的腿在飞，本票可以单独排一小发。⛔ 零翻框、零 push。

- [2026-10-14 11:37:32 +08] agent=288-r1（落地腿）did=**停手上报：本票那处对齐在 HEAD 上已不存在**——「前」态（blob `a917e0dc`＝`18d43f8e^:cmd/wisp/panel_inbound_guards_35r3_test.go`）两把尺复现票面逐字（裸 `gofmt -d` rc=1 14 行／`gofumpt -d` rc=1 14 行，同一 hunk `@@ -63,8 +63,8 @@`），而 `18d43f8e`（票 298 的 298-r1，2026-10-10 16:19，"35r3＝const 块内尾随注释列对齐"1 hunk／numstat 2 2）已把它洗完、`18d43f8e:path`＝`HEAD:path`＝`9713dc37`（改后无人再动）⇒ 工作树现量 `gofmt -d`／`gofumpt -d` 同声 **0／0**、件级 `gofmt -l`／`gofumpt -l` 皆空、`git status` 对该件 0 行、`i/lf w/lf` CR=0；⇒ 零写面（⛔ 无字节可改，改任何东西都违反 `AC#1` 自己的"只改那一处"）⇒ 本腿**零 commit 于该件**、按派单 §5「未定义即停」停手；`AC#2` 具名未跑整包（无「改前/改后」对可作差＋`wisp.test.exe` PID 38192 占窗；同一 blob 过渡的等价证据在册：`tr -d '[:space:]'|cmp` rc=0＋编排者 10-12 成对两发新增红 0，件 `probes/orch/20261012/logs/298ac4/{new-red.txt,summary.txt}`）；`AC#3` 树态读数＝build rc=0／d22scan rc=0 clean（⛔ 非写腿门禁）；`AC#1/#2/#3` 一格未勾、`Status` 未动；此债自 298 起归票 298 的账（其残余只剩 `AC#5`＝另一枚件的 `gofumpt` 残留） next=编排者裁本票收口措辞（建议句＝债由 `18d43f8e` 付清）＋若仍要一发整包现色则另派（等 PID 38192 退场）；交件 `probes/288/r1/{00-anchor-and-gates.md,logs/*}`（两列并排＋锚链＋树态 rc，无一 0 字节入册件）
- [2026-10-14 12:05:04 +08] agent=288-v1（非实现者验收腿）did=判语：「**债由 `18d43f8e` 付清**」成立——blob 链现跑三枚 `a917e0dc…`／`9713dc37…`／`9713dc37…`（＝HEAD；`git log 18d43f8e..HEAD -- <件>` 空）；前态两把尺同声 14 行 rc=1（两输出 `cmp` rc=0、hunk `@@ -63,8`）、后态同声 0 行 rc=0；付款笔核成＝diff 头逐字 `index a917e0dc..9713dc37`、numstat `2 2`，且 `gofmt(pre blob)`⇄post blob 逐字节 `cmp` rc=0（净效果恰为 gofmt、剥空白 8301⇄8301、工作树⇄HEAD blob rc=0）；链收口现读＝`docs/evidence/s1/ticket298-298-r1-accept-v1.md`＋票 298 编排者终裁（AC#0–3 翻）＋AC#4 翻勾（10-12 14:52，新增红 0）；目录级两枚（`models.go` CRLF 假枚／35r2 gofumpt 残留＝298 AC#5 的账）与本票无关、288 件在两把尺上皆空；AC 收口建议＝AC#1 债已付清勾形、AC#2/AC#3 零写面 N/A 承接形（措辞在证据件 §二）；背景 1 发整包 rc=1、唯一顶层 FAIL＝`TestAC2SealNoticeLandsInTheRunLegLogFile`（GOTMPDIR=C:\tmp 环境形、台账在册非真红）＋build rc=0／d22scan rc=0 clean（⛔ 非判据）；交件 `docs/evidence/s1/ticket288-288-v1-accept-v1.md`＋`probes/288/v1/`（锚＋logs） next=编排者裁收口措辞并翻框；⛔ 零字节改动、零翻框、零 push
