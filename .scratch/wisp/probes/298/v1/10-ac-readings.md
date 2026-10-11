# 298-v1 — 10-ac-readings：五格逐格复跑读数（非实现者验收腿；⛔ 本腿零翻框、零实现修复）

起手锚＝`00-anchor.md`（HEAD `27059960`；三发期间别的腿交错入库，HEAD 漂到 `5f2289e3`——所有读数一律对象锚：落地腿四笔 `24ae75b7`／`18d43f8e`／`c14fd204`／`1b311587`，改前一棵＝`18d43f8e^`＝`48e7b0a4`）。原始件全部在 `logs/`，每行带 `rc=`；⛔ `.out`、⛔ 0 字节件。临时目录一律仓外。

## AC#0 两枚真格式债的 blob 脏形＝逐枚复跑，与腿件对拉 ⇒ **成立**

- 同一性前提：`cmp <(git show 6414a4bb:<p>) <(git show 48e7b0a4:<p>)` 两枚 **rc=0**（`logs/ac0-anchor-blob-identity.txt`）＝腿锚与我对锚同 blob。
- `35r3`：改前 blob `gofmt -d` ⇒ **1 hunk `@@ -63,8 +63,8 @@`／14 行（内容行 12）／const 块内两行尾随 `//` 注释列对齐**；首笔 `7f9d6e40`（票 35 `35-r3` 腿，`git log -- <p>` 名册 oldest 具名**对**）blob 同尺 ⇒ 同 hunk，且与改前 blob `cmp` **rc=0（9759 字节）**＝首笔即脏、中间无人动过。
- `35r2`：改前 blob ⇒ **1 hunk `@@ -804,8 +804,8 @@`／14 行／两枚单行方法 `{` 前空格 4→1**；首笔 `286a7f30`（票 35 `35-r2` 腿，oldest 具名**对**）⇒ hunk `@@ -772,8 +772,8 @@`、1824 行／52759 字节（`wc -l` 实测 1824）⇄ 改前 2142 行／73476 字节＝位移 **+32 行**，与 `2fc5f5c9` 提交正文自记「:775-776→:807-808」逐字对上 ⇒ 归因链闭合。
- diff 全文对拉：两份 `gofmt -d` 剥路径头后与腿原件 `diff` ⇒ **rc=0、内容行 12⇄12**（`logs/ac0-leg-vs-v1-diffbody.35r3.txt`／`.35r2.txt`；空输出＝相同，件尾有说明行。同两名的 `-inbound-guards/-transport` 后缀文件＝本腿第一发文件名变量写歪的失败尝试产物，只建不删、⛔ 引用）。
- **具名腿的口径笔误**：`gofmt -d` 在 diff 非空时实退 **rc=1**（本腿两枚 pre-blob 发均 rc=1；干净 blob 目录同尺 rc=0 互证），腿件写 "rc=0"＝记录缺陷，⛔ 影响实质。

## AC#1 "只动空白层"＝本腿三把独立尺（⛔ 转抄腿的自述）⇒ **成立**

1. `git diff --numstat 48e7b0a4 18d43f8e -- <两枚>` ⇒ `2 2`＋`2 2`，rc=0。
2. 同尺加 `-w`（ignore-all-space）⇒ **rc=0、0 行**＝每一处差异都在空白层（件尾说明行已注）。
3. 最硬：两版 `git show | tr -d '[:space:]'` 逐字节 `cmp` ⇒ `35r3 rc=0（8301⇄8301 字节）`、`35r2 rc=0（59964⇄59964 字节）`＝**非空白字节一字不差**——标识符／反引号串／注释文字／断言／用例名⛔ 一枚可指为变动。（第一发此尺因本腿变量笔误失效 rc=2＝尺坏⛔ 码坏，第二发即上数，失败件留 `logs/myscript-stray-*.bak`。）
4. 判据末条：改后 blob 那把 `gofmt -l` ⇒ `18d43f8e` 与当前 HEAD 两枚均 **count=0**；`cmp` 改前剥空白 ⇄ HEADnow 剥空白 rc=0＝修后无人再动。`git show --name-only --format= 18d43f8e` 名册只有那两枚。

## AC#2 假枚标注法逐字可重跑；那一问确留给编排者、腿未自裁未清洗 ⇒ **成立**

- 三把尺原样复跑（`cmd/wisp/models.go`）：`git ls-files --eol` ⇒ 逐字 `i/lf    w/crlf  attr/text eol=lf`；CR＝**334**＝行数；`git show HEAD:` 落仓外 `gofmt -l` ⇒ **0 命中**；`git status --porcelain` 0 行 ⇒ 标注形 `<路径>（换行符假枚：i/lf w/crlf，CR=<n>，blob 尺不命中）` 下一位可逐字重跑且当场成立。
- 扩射程：`git ls-files --eol cmd internal tools | grep 'w/crlf'` ⇒ **5 枚**；全仓 tracked 尺的非 `.scratch` 5 枚（`models.go`／`pending_read.go`／`tools.go`／`provenance.go`／`bridge.go`）＝同一集合 ⇒ 票面更正节"1→5"复现。`core.autocrlf=true` 只读。
- "要不要真换成本机 LF"＝腿件 `20-…` 具名留编排者、只转引编排者"本轮不换"；四笔名册里 `models.go`／`.gitattributes` **0 命中**＝⛔ 清洗⛔ checkout 族。本格判据原文是"写清楚"⛔"必须换"——⛔ 把它读成欠清洗。

## AC#3 口径句写进凭据件＋两把尺并排带枚数 ⇒ **成立（附两枚具名文案缺陷）**

票面追加的逐字尺（blob 档用 `for p in $(git ls-files cmd/wisp | grep '\.go$')` 展开形）本腿原样重跑：

| 尺 | 腿声称（改后） | 本腿现跑 |
|---|---|---|
| `gofmt -l cmd/wisp`（工作树） | 1 | **rc=0 count=1**＝models.go 假枚 |
| blob 尺（`git ls-files cmd/wisp \| grep '\.go$'` 展开落仓外） | 0 | **rc=0 count=0**（当日 tracked `.go` 枚数 114→今 115＝别人交错新增一枚，快照漂移、⛔ 影响两把尺零命中） |
| `gofmt -l cmd/wisp internal tools`（工作树） | 5 | **rc=0 count=5** |
| `git archive HEAD …` 导出树（blob） | 0 | **rc=0 count=0**；同树 `gofumpt -l` ⇒ **1**（35r2 残留） |
| `git ls-files '*.go' \| xargs gofmt -l`（stdout/stderr 分离） | 27＝22 `.scratch`＋5 假枚 | **rc=123 count=27，scratch=22，非 .scratch=5、真债 0**；stderr 行＝`fs_broken.go:4:1` |
| `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp`（工作树） | 2 | **rc=0 count=2** |

七把与腿声称**全部同值**＝"两把尺并排带枚数＋射程目录＋rc/流分离"的口径句可被下一位逐字重跑。腿内部 28/23⇄27/22 的污染形已由腿自己在 `20-…` 具名自纠，本腿复跑站在 27/22 一边。
**具名文案缺陷两枚（⛔ 阻格、更正归编排者）**：① `20-…` 的一句话判据写"三件事"，票面同一句写"四件事"（凭据件未同步）；② `20-…` 的 blob 尺逐字⛔ 带 `grep '\.go$'`（票面形才是可原样重跑那把）。

## AC#4 门禁与 `cmd/wisp` 那几发（本腿自跑）⇒ **成立（附ⓐⓑ两枚具名欠账）**

硬要求逐件（每件自落 rc＝`logs/`）：`GOFLAGS= go build ./...` **rc=0**；`sh scripts/d22scan.sh` **rc=0 clean**（23099 字节；`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77` 与 `ban #8 cmd/: 119 files` 与腿声称**逐字同值**；注：该 runtests 只跑 `tools/d22scan` 自身正控，⛔ 覆盖 cmd/wisp——cmd/wisp 的颜色只能靠下面各发）；`gofmt -l cmd/wisp` 工作树⇄blob＝**1⇄0**；gofumpt 全路径＝**rc=0 count=2**。

`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1` 六发，**每发前 `tasklist //FI "IMAGENAME eq wisp.exe"` 均 0 枚**：

| 发 | 内容 | rc | 红名册（顶层） |
|---|---|---|---|
| run1（quiet，改后） | 工作树 | 1 | `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`＋`TestAC14GoSideEvalPushReachesThePage` |
| run2（quiet，改后） | 工作树 | 1 | 同上 2 枚，逐字相同 |
| run3（`-v`，改后） | 工作树 | 1 | 同上 2 枚；顶层 **PASS=285 FAIL=2 SKIP=1**（名册 288 行）；含子测试两把 **PASS=395**＝两把尺差 **110**，本腿数的口径＝**顶层** |
| overlay-v（`-overlay` 载 `48e7b0a4` 两枚件内容＝改前） | 同一工作树 | 1 | 那 2 枚 ⇄ **＋`TestAC228ExitRequestDuringBootStillLeavesThroughD38E`** |
| overlay-quiet（改前） | 同一工作树 | 1 | 那 2 枚 ＋ `TestTicket223HandEditedFsLooseningCostsAnL2Card` |

- **逐名作差＝新增红 0**：run3-v ⇄ overlay-v 名册 288 行 `diff` ⇒ **唯一一处换位**＝`AC228…D38E`（post **PASS** ⇄ pre **FAIL**），其余 287 行逐字同色（`logs/roster-diff-post-vs-pre.txt`）。该枚＝台账逐字在册的 boot 时序竞态（"整包序 3 发 1 红、隔离 5/5 绿"，`pending-and-issues.md:9987`），落点 `cmd/wisp/resident_ball_228_windows_test.go`＝**⛔ 298 两枚件**；`TestTicket223…` 同理在 `config_reload_223_test.go`。**改后侧零新增红、两枚常驻红在改前改后两侧同色**＝票面"作差＝新增红 0 枚"闭合。
- 两枚常驻红与 298 无关的双证：① 定义在 `panel_resident_windows_test.go`（⛔ 本票射程文件）；② `305-a1/20-table2-attribution.md` 整表在册（dist＝探针文档 1044 字节那一族）。
- 同码作差成立：`git log 02053534..c023e8e0 -- cmd/wisp` **空**、`934b350a..5f2289e3 -- cmd/wisp` **空**＝run3-v（post）与 overlay-v（pre）之间 `cmd/wisp` 产码零变动，唯一变量＝那两枚件的 blob 内容。台面漂移具名＝三发期间别的腿（含 `303-r1`）交错入库到仓的其他面，HEAD `27059960`→`5f2289e3`；⛔ 时长读数当凭据（410s 一发＝编排者在跑反形的机器）。
- **ⓐ 腿欠账**：腿对该格**零发**（派单⇄票面冲突具名报回、⛔ 假称跑过）⇒ 改后两发＋名册一发由本腿补齐；改前"同台面严格两发"的历史窗口已关（见下 CI/真机节），本腿以 `-overlay` 同工作树载史 blob 内容代偿并具名形差。
- **ⓑ gofumpt 残留 1 枚**（35r2，HEAD 既有债、⛔ 本票射程）＝编排者裁另立与否。
- 越界尺逐笔（`logs/ac4-scope-per-commit.txt`）：四笔名册越出白名单（两枚件＋`probes/298/r1/**`＋票面 bullet）**0 行**；禁区零命中；r1 件 `-type f -size 0`＝0、⛔ `.out`（本腿一把 `find` 未带 `-type f` 误计目录 1 枚＝验收腿尺口径坑，具名记我）。
- 本腿另一笔自报：脚本变量笔误曾把两枚剥空白中间件落到**仓库根**，当场 `mv` 入 `logs/myscript-stray-*.bak`（只建不删、根计数复原＝实测 `ls panel_*test.go` 0 枚）。
- 腿件 rc 留痕抽检：24 枚 logs 中 8 枚原始 diff 无自落 rc 行（`ac1-diff-U0`／`ac1-numstat`／两枚 `blob-gofmt-d-*`／`d22scan-name-level-*`／`gofumpt-residue-*`／两枚 `*-gofumpt-d-*`）＝"每件自落 rc"声称过宽，具名、⛔ 扣格（md 叙述里各尺均带 rc）。

## 只有真机／只有 CI 才能闭合的格（具名）

1. ★票面归因支"CI 那两步 `go vet`／`gofmt denominator` 长期 **skipped＝从未求值**"⇒ 那两发的**色**只能由 push 后 CI 给（本机绿红⛔ 顶替；归编排者：推送＋`gh` 取日志）。五格本体⛔ 等它，但"落地时没人看见红"的坐实在它身上。
2. ★`AC#4`ⓐ 的严格形（`48e7b0a4` 那棵树在**当时台面**上 ≥2 发颜色）＝窗口已关；本腿交的是 overlay 近似形（同码面已证）。要严格形只能 CI/历史重放——归编排者记账。
3. ★winlive 真窗族⛔ 属票 298 射程（列此只因派单这一问要求具名：真窗那一发⛔ 属于本票，本票的机器面已由上表闭合）。

## 判语汇总（详细推理在裁决表 `docs/evidence/s1/ticket298-298-r1-accept-v1.md`）

`AC#0` 成立／`AC#1` 成立／`AC#2` 成立／`AC#3` 成立（附两枚具名文案缺陷）／`AC#4` 成立（附ⓐⓑ两枚具名欠账）。五格框⛔ 本腿翻，翻勾归编排者。
