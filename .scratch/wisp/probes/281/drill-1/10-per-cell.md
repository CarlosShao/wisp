# 281-drill-1 逐格件：10 枚补勾行的两栏原文对照（承载表行 vs 裁决表行）

全部现读自对象层 `git show HEAD:<path>`，HEAD `69c2bfd9`，取数 2026-10-09 11:1x +0800。
"承载"行原文取自 `.scratch/wisp/probes/281/r1/10-reconcile.md`（十枚数据行真号 `:18`–`:27`）。
"裁决"行原文取自各 `docs/evidence/s1/*.md` 的对应单行。
命中标记＝对该行内被引短语跑 `grep -nF`，输出的行号（期望＝承载表写的那一枚）。

---

## 格 1 — 票 92·AC#6 ｜ 承载 `:18` → `92b-adversarial-acceptance.md:138` ｜ 判语：**对上**

承载表 `:18` 现读（逐字）：
```
| 1 | 92·AC#6 | `docs/evidence/s1/92b-adversarial-acceptance.md:138`「AC#6 台账与门禁 \| **通过（数字全复算）** \| 〔独立复现〕 \| …`go test -count=2 -v ./internal/panel/ ./internal/tools/` = **360/238/0/0 rc=0**（与它报的逐位相等）、`sh scripts/d22scan.sh` 工作树 rc=0 且 `ban #6 frontend/=43`…」 | `93bbb54d` |
```
裁决表 `:138` 现读（逐字）：
```
| AC#6 台账与门禁 | **通过（数字全复算）** | 〔独立复现〕 | `gofmt`/`gofumpt` 空、`go vet ./internal/panel/ ./internal/tools/ ./internal/memory/` rc=0、`go test -count=2 -v ./internal/panel/ ./internal/tools/` = **360/238/0/0 rc=0**（与它报的逐位相等）、`sh scripts/d22scan.sh` 工作树 rc=0 且 `ban #6 frontend/=43`、`ban #8 internal/=373`（它报 371，+2 是邻居在飞文件）；⚠ **整步 `portable-tests.sh` 它未跑**（`R-92b-4`） |
```
逐短语：`AC#6 台账与门禁` 中／`**通过（数字全复算）**` 中／`〔独立复现〕` 中／`360/238/0/0 rc=0`（与它报的逐位相等）中／`工作树 rc=0 且 ban #6 frontend/=43` 中。承载用 `…` 省略的部分（gofmt/gofumpt/vet 前置、`ban #8 internal/=373` 尾注）在裁决行内存在，非编造。
`grep -nF` 命中＝**138**（唯一）。

对照警示（同票另一张表，见 `00-findings.md` §四#2）：`92-adversarial-acceptance.md:147` 现读含 `**结论：FAIL —— 附我本机实测数字。**（AC#6 原文要求 gofumpt -l 空，它不空；且 POSIX 读数不可复现。）`——返工前那轮判的是 FAIL；凭据引的是 `92b`，方向正确。

---

## 格 2 — 票 92·AC#7 ｜ 承载 `:19` → `92b-adversarial-acceptance.md:139` ｜ 判语：**对上**

承载表 `:19` 现读：
```
| 2 | 92·AC#7 | `:139`「AC#7 不做 git 切换 \| **通过** \| 〔独立复现〕 \| 包内用例 PASS + 我把 `fixtures/`、`dist/` 也 grep 了一遍 0 命中」 | `93bbb54d` |
```
裁决表 `:139` 现读：
```
| AC#7 不做 git 切换 | **通过** | 〔独立复现〕 | 包内用例 PASS + 我把 `fixtures/`、`dist/` 也 grep 了一遍 0 命中 |
```
整行逐字等同（承载只是把行首尾的 `|` 加了转义）。`grep -nF` 命中＝**139**（唯一）。
⚠ 承载此处**只写 `:139` 未写文件名**（承前 `:18` 的 `92b-`）——沿用到哪张表需读上一行才能定，本腿按承前解释并已在 00 件里具名提示 `<NN>` 模板陷阱。

---

## 格 3 — 票 97·AC#5 ｜ 承载 `:20` → `97-adversarial-acceptance.md:17` ｜ 判语：**对上**

承载表 `:20` 现读：
```
| 3 | 97·AC#5 | `docs/evidence/s1/97-adversarial-acceptance.md:17`「AC#5 …\| **通过（PASS 相加≠RUN 的写法要更正，数字自洽）** \| 〔独立复现〕 \| …`go test -count=2 -v ./internal/agent/approval/` **rc=0**…`sh scripts/d22scan.sh`…**rc=0 / clean**」 | `226d3717` |
```
裁决表 `:17` 现读（该行极长，此处保留承载引到的片段与首尾；整行已在会话现读输出中）：
```
| AC#5 | 按包门禁：gofmt/gofumpt 空、`go vet` rc=0、`go test -count=2` 四数逐条点名（2 条 SKIP 必须点名）、收尾 `sh scripts/d22scan.sh` | **通过（PASS 相加≠RUN 的写法要更正，数字自洽）** | 〔独立复现〕 | 全部在 `f140079` 纯净快照：…`go test -count=2 -v ./internal/agent/approval/` **rc=0**，四数我自己数：`=== RUN` **98**、顶层 `--- PASS` **58**、顶层 `--- FAIL` **0**、顶层 `--- SKIP` **2**…⚠ **PASS 两个数相加 == RUN 吗？不等于**…⇒ `R-97-3`…`sh scripts/d22scan.sh`（**未**从仓根 `go run ./tools/d22scan`）**rc=0 / clean**…
```
`grep -nF` 命中：`| AC#5 | 按包门禁：gofmt/gofumpt 空`→**17**；`**通过（PASS 相加≠RUN 的写法要更正，数字自洽）**`→**17**。
**先例核查（编排者点名的 ±1 漂）**：97 表全长 51 行；`:18` 现读＝**空行**（`## 第三条路…` 在 `:19`）。⇒ 承载表写的 `:17` 就是真身，那一枚 `:18` 的旧引出自上游判语件 `done-class-b-1`（`10-reconcile.md:66` 具名），不是承载表自己。本轮无漂。

---

## 格 4 — 票 104·AC#3 ｜ 承载 `:21` → `104-adversarial-acceptance.md:115` ｜ 判语：**对上**

承载表 `:21` 现读：
```
| 4 | 104·AC#3 | `docs/evidence/s1/104-adversarial-acceptance.md:115`「**结论：通过。〔独立复现〕**（票面要求的三腿变异 + 我自加两发全部我自己在 `/tmp` 仓外快照重抽、重打、重量；另有一发 HEAD 漂移核对，非变异）」 | `d1c7e1aa` |
```
裁决表 `:115` 现读：
```
**结论：通过。〔独立复现〕**（票面要求的三腿变异 + 我自加两发全部我自己在 `/tmp` 仓外快照重抽、重打、重量；另有一发 HEAD 漂移核对，非变异）
```
整行逐字等同（`:114`、`:116` 皆空行⇒该行是独立成段的结论句，不是表格行）。`grep -nF` 命中＝**115**（唯一）。

---

## 格 5 — 票 104·AC#4 ｜ 承载 `:22` → `104-adversarial-acceptance.md:162` ｜ 判语：**对上**

承载表 `:22` 现读：
```
| 5 | 104·AC#4 | `:162`「**结论：通过。〔独立复现〕**（每个数字都是我在 `/tmp/ac104-fk`（`4d43447` 纯净树）本机重跑的，未照抄自述）」 | `d1c7e1aa` |
```
裁决表 `:162` 现读：
```
**结论：通过。〔独立复现〕**（每个数字都是我在 `/tmp/ac104-fk`（`4d43447` 纯净树）本机重跑的，未照抄自述）
```
整行逐字等同。`grep -nF` 命中＝**162**（唯一）。
⚠ 同一枚文件里 `**结论：通过。〔独立复现〕**` 这个前缀**不止一处**（`:115`/`:162` 同形），承载靠**括号内的限定语**（"票面要求的三腿变异…"/"每个数字都是我在 /tmp/ac104-fk…"）把两枚区分开——两枚限定语各自唯一命中，**不构成歧义**；但若有人只抄前缀不抄限定语，就会撞上"两行同句"。登记为形状风险，不是错。

---

## 格 6 — 票 105·AC#5 ｜ 承载 `:23` → `105-adversarial-acceptance.md:16` ＋ `:153` ＋ `:194` ｜ 判语：**对上 ×3**

承载表 `:23` 现读：
```
| 6 | 105·AC#5 | `docs/evidence/s1/105-adversarial-acceptance.md:16`「AC#5 \| 门禁…\| **通过** — 分包四数复现，d22scan 三枚树 rc=0、`internal/` 363→370 升；两处交件数字口径差异已点名 \| 〔独立复现〕」＋`:153`「**口径**：全部跑在 `f6818f2` 的**纯净快照** `/tmp/ac105b-gates`」＋`:194`「"各 scope 不降"（A64②）**在纯净树复算成立**：`internal/` 363→370…rc=0 clean」 | `2a2dad22` |
```
裁决表 `:16` 现读：
```
| AC#5 | 门禁：go test/gofmt/gofumpt/vet/GOOS=linux vet/d22scan 纯净快照 rc=0，台账不降 | **通过** — 分包四数复现，d22scan 三枚树 rc=0、`internal/` 363→370 升；两处交件数字口径差异已点名 | 〔独立复现〕 |
```
裁决表 `:153` 现读：
```
**口径**：全部跑在 `f6818f2` 的**纯净快照** `/tmp/ac105b-gates`（`git archive f6818f2 | tar -x -C …`，容器/测试内 `ls -l /wisp/go.mod` 之类的"文件在"证明逐条附）。
```
裁决表 `:194` 现读：
```
⇒ "各 scope 不降"（A64②）**在纯净树复算成立**：`internal/` 363→370（本票新增 2 枚测试）升、其余各 scope 持平，rc=0 clean。
```
`grep -nF` 命中：`| AC#5 | 门禁：go test/gofmt/gofumpt/vet`→**16**；`**通过** — 分包四数复现，d22scan 三枚树 rc=0`→**16**；`**口径**：全部跑在 f6818f2 的**纯净快照**`→**153**；`（A64②）**在纯净树复算成立**`→**194**。四处全部单行唯一命中。
署名补充（见 00 件 §五）：本枚三行的作者按该文件 `:19` 自陈＝接续代理 `acceptor-ticket105b`，而文件标题署 `acceptor-ticket105`。

---

## 格 7 — 票 110·AC#3 ｜ 承载 `:24` → `110-adversarial-acceptance.md:16` ｜ 判语：**对上**

承载表 `:24` 现读：
```
| 7 | 110·AC#3 | `docs/evidence/s1/110-adversarial-acceptance.md:16`「AC#3 \| 新步自己会红 + 不是空仪器…\| **通过（三发全部独立复现）** \| 〔独立复现〕 \| 全部在 `/tmp/wisp-ac110-snap`；仓库树 winsec 全程未动」 | `c9348a99` |
```
裁决表 `:16` 现读（行首尾片段，中段为三发变异 M1/M2/M3 读数）：
```
| AC#3 | 新步自己会红 + 不是空仪器（同链 `grep -n` 证落地、先 `go build` rc=0） | **通过（三发全部独立复现）** | 〔独立复现〕 | 全部在 `/tmp/wisp-ac110-snap`；仓库树 winsec 全程未动。**M1 安全断言弄坏**：…`bash scripts/winsec-tests.sh` **rc=1，=== RUN=68 / PASS=9 / FAIL=26 / SKIP=0**…**M2 仪器弄空**…⇒ **rc=2** + `GUARD 1 …`**M3 包被掏空**…**rc=1** + `GUARD 2 …`…
```
`grep -nF` 命中：`| AC#3 | 新步自己会红 + 不是空仪器`→**16**；`**通过（三发全部独立复现）** | 〔独立复现〕`→**16**。

---

## 格 8 — 票 110·AC#5 ｜ 承载 `:25` → `110-adversarial-acceptance.md:18` ｜ 判语：**对上**

承载表 `:25` 现读：
```
| 8 | 110·AC#5 | `:18`「AC#5 \| `bash -n` rc=0；…\| **通过（附一条后果登记）** \| 〔独立复现〕 \| `bash -n scripts/winsec-tests.sh` rc=0、`bash -n scripts/portable-tests.sh` rc=0…」 | `c9348a99` |
```
裁决表 `:18` 现读（片段）：
```
| AC#5 | `bash -n` rc=0；`sh scripts/d22scan.sh` 纯净快照 rc=0 且各 scope 不降；新步不得排在会失败的步之后 | **通过（附一条后果登记）** | 〔独立复现〕 | `bash -n scripts/winsec-tests.sh` rc=0、`bash -n scripts/portable-tests.sh` rc=0（`9e9a2f5` 快照）。`sh scripts/d22scan.sh` 在 `9e9a2f5` 纯净快照 **rc=0**：`bans #1-5 internal/=202、cmd/=20、ban #6 frontend/=40…`…**步骤位置我自己判一次**：…⚠ 代价当场兑现：见 `R-110-2`。 |
```
`grep -nF` 命中：`| AC#5 | `bash -n` rc=0`→**18**；`**通过（附一条后果登记）** | 〔独立复现〕`→**18**。
⚠ 该行自陈的"附一条后果登记"＝`R-110-2`，与承载表 `:30` 的"110#5＝R-110-2 后果仍在开放池"同向——**勾与残账在两层都写着，没被抹掉**（本腿不判它合不合规）。

---

## 格 9 — 票 113·AC#4 ｜ 承载 `:26` → `113-adversarial-acceptance.md:16` ｜ 判语：**对上**

承载表 `:26` 现读：
```
| 9 | 113·AC#4 | `docs/evidence/s1/113-adversarial-acceptance.md:16`「AC#4 变异四发 \| MUT-A 9 红 / MUT-B 恰 1 红 / MUT-C 9 红 / MUT-D2 2 红 \| **四发全部我自己下刀、自己复量**… \| 〔独立复现〕 \| **绿**」 | `acea1708` |
```
裁决表 `:16` 现读（片段）：
```
| AC#4 变异四发 | MUT-A 9 红 / MUT-B 恰 1 红 / MUT-C 9 红 / MUT-D2 2 红 | **四发全部我自己下刀、自己复量**（`git archive 3c5d1c3` → `/d/tmp/mut-ac113-{A,B,C,D2}`，仓库内无 worktree）…读数一字不差：A `RUN=16 PASS=7 FAIL=9`；B `FAIL=1`…C `FAIL=9`…D2 `FAIL=2`…⇒ 全集主张成立，但**叶子这一维在 `SealFile` 方向没有交付用例**，见下"攻#4" | 〔独立复现〕 | **绿** |
```
`grep -nF` 命中：`| AC#4 变异四发 | MUT-A 9 红 / MUT-B 恰 1 红`→**16**；`**四发全部我自己下刀、自己复量**`→**16**。

---

## 格 10 — 票 113·AC#5 ｜ 承载 `:27` → `113-adversarial-acceptance.md:17` ｜ 判语：**对上**

承载表 `:27` 现读：
```
| 10 | 113·AC#5 | `:17`「AC#5 门禁四数 + 三门 + d22 \| 540/532/0/8…\| …**rc=0、`=== RUN` 540 / `PASS:` 532 / `FAIL:` 0 / `SKIP:` 8**…**d22：`sh scripts/d22scan.sh` rc=0，但 ban#8 `internal/=371`，不是它报的 368**… \| 〔独立复现〕 \| **绿（四数与三门），但 ban#8 那格的 sha↔数字配对错了一格**」 | `acea1708` |
```
裁决表 `:17` 现读：
```
| AC#5 门禁四数 + 三门 + d22 | 540/532/0/8、ban#8 internal/=368 | 五包 `-count=2 -v`（winsec/memory/risk/secret/models，容器 `golang:1.27`、`CGO_ENABLED=0`、uid 0）⇒ **rc=0、`=== RUN` 540 / `PASS:` 532 / `FAIL:` 0 / `SKIP:` 8**，逐包 `ok winsec 0.273s / ok memory 24.499s / ok risk 7.822s / ok secret 0.009s / ok models 0.530s`…`gofmt -l .` **0 行**；容器内 `go vet` 五包 rc=0、`GOOS=windows go vet ./internal/winsec/` rc=0、`GOOS=darwin go vet ./internal/winsec/` rc=0。**d22：`sh scripts/d22scan.sh` rc=0，但 ban#8 `internal/=371`，不是它报的 368** ⇒ 见"攻#5" | 〔独立复现〕 | **绿（四数与三门），但 ban#8 那格的 sha↔数字配对错了一格** |
```
`grep -nF` 命中：`| AC#5 门禁四数 + 三门 + d22`→**17**；`不是它报的 368`→**17**。
⚠ 这一枚是"**通过但当场点名交件方报错数**"的形状（裁决行自己写着 ban#8 371≠368），承载表 `:30` 也把它登记为 `R-113-F`→`A89④`——两层一致，不是凭据与残账打架。

---

## 附：本件用到的取数命令形状（可复核）

- 承载表行号：`git show HEAD:.scratch/wisp/probes/281/r1/10-reconcile.md | grep -nE '^\| ([0-9]+) \|'` → `18..27`
- 单行现读：`git show HEAD:<path> | awk 'NR==<M> {print NR": "$0}'`
- 命中定位：`git show HEAD:<path> | grep -nF -- '<被引短语>'`（固定串、含空格与反引号，逐枚期望唯一命中）
- 越界检查：`git show HEAD:<path> | wc -l`（92b=347 / 97=51 / 104=306 / 105=232 / 110=43 / 113=134）
- 署名：`git show HEAD:<path> | awk 'NR<=12'` ＋ `grep -oE 'acceptor-[a-z0-9-]+' | sort -u`
