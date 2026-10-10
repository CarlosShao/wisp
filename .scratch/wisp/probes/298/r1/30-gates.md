# 298-r1 — 30-gates：门禁读数＋越界逐笔（AC#1 判据 / AC#4 全格）

锚点起手＝`6414a4bb`（`git log -1 --format='%h %ad' --date=iso-strict`，`rc=0`），详见 `00-anchor.md`。
本件只交**读数**；五格判语归非实现者 `298-v1`，⛔ 本腿自勾任何一格。

## 1. 停手条件（票面 AC#1 起手那条，逐字执行）

```
ruler: git status --porcelain -- cmd/wisp/panel_inbound_guards_35r3_test.go cmd/wisp/panel_transport_35r2_test.go
（无输出）
rc=0 lines=0
```
⇒ **未触发** ⇒ 放行。旁证（防"跟踪状态≠干净"那一形）：
```
ruler: diff <(git show HEAD:cmd/wisp/<两枚>) cmd/wisp/<两枚>
两枚皆 IDENTICAL          rc=0
```
⇒ 本腿 `gofmt -w` 落下去覆盖的就是 blob 自身，⛔ 不可能吞别人未提交的改动。
进程尺＝`tasklist //FI "IMAGENAME eq wisp.exe"` 0 枚、`…balldebug.exe` 0 枚（只记录，⛔ 未杀）。

## 2. AC#1 判据：每个 hunk 都落在空白／换行／对齐上

尺＝`git diff -U0 -- cmd/wisp/panel_inbound_guards_35r3_test.go cmd/wisp/panel_transport_35r2_test.go`
（原件全文＝`logs/ac1-diff-U0.txt`；`gofmt -w` 前先跑过 `gofmt -d` 于 HEAD blob，形状见 `10-ac0-dirt-shape.md`）

| 枚 | hunk 数 | 落点（blob 号） | 形状 | numstat | 语义面 |
|---|---|---|---|---|---|
| `panel_inbound_guards_35r3_test.go` | **1** | `:66-67`（blob `@@ -63,8 +63,8 @@`） | `const` 块内两行尾随 `//` 注释的**前导空格数** | `2 2` | 标识符／反引号字符串／注释正文**逐字未变** |
| `panel_transport_35r2_test.go` | **1** | `:807-808`（blob `@@ -804,8 +804,8 @@`） | 两枚单行方法的**`{` 前空格数** 4→1 | `2 2` | 接收者／方法名／返回类型／函数体**逐字未变** |

⇒ 两枚各 1 hunk、各只换 2 行，**每个 hunk 的每一行差异都在空白层**；本腿逐段过目后⛔ 发现任何标识符、字符串、断言、用例名变动。

尺＝`git diff --numstat -- cmd/wisp/panel_inbound_guards_35r3_test.go cmd/wisp/panel_transport_35r2_test.go` ⇒ `2 2` ＋ `2 2`，`rc=0`
尺＝`git diff --numstat HEAD~2 HEAD -- '*.go'` ⇒ **只有那两枚路径**（第三枚 `.go` 路径＝0 枚 ⇒ 没顺洗整包），`rc=0`
尺＝AC#1 末条「改完 `git show HEAD:<path> | gofmt -l` 式的那一把必须空」⇒ 提交后重跑：
```
ruler: git show HEAD:<三枚路径> 落仓外目录; gofmt -l <目录>      # 原件 logs/post-blob-rulers.txt
rc=0 count=0        ⇒ 空，达成
```

## 3. AC#4 门禁逐条（每条自落 `rc=`；⛔ 既有红＝无，改后必逐字相同这一条因此自动成立）

| 尺（逐字） | 改前 | 改后 | 判据形状 |
|---|---|---|---|
| `GOFLAGS= go build ./...` | `rc=0`（`logs/pre-build.txt`） | `rc=0`（`logs/post-rulers.txt`） | 不变 |
| `sh scripts/d22scan.sh` | `rc=0` clean（`logs/pre-d22scan.txt`，23093 字节） | `rc=0` clean（`logs/post-d22scan.txt`，23091 字节） | 名册级作差＝**逐字相同** |
| `gofmt -l cmd/wisp`（工作树） | `rc=0` 3 枚 | `rc=0` **1 枚** | 只减不增 |
| `gofmt -l cmd/wisp internal tools` | 7 枚＝票面更正节读数（⛔ 本腿改前跑的是**同射程的 `gofumpt`＝7 枚**，`rc=0`；`gofmt` 那一把我⛔ 没在改前跑，欠在此具名） | `rc=0` **5 枚**（`logs/ac3-four-rulers.txt`） | 只减不增（5 ⊂ 7，差额＝本票那两枚） |
| `gofmt -l <HEAD blob 落仓外>` | `rc=0` 2 枚（`logs/../00-anchor.md` G5.2） | `rc=0` **0 枚**（`logs/post-blob-rulers.txt`；宽域 658 枚同树亦 0，`logs/ac3-blob-wide-ruler.txt`） | 只减不增 |
| `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp` | `rc=0` 3 枚（`logs/pre-gofumpt.txt`，⛔ 裸 `gofumpt`＝rc=127 那一形本腿没用） | `rc=0` **2 枚**（`logs/post-rulers.txt`） | 只减不增 |
| `gofumpt -l cmd/wisp internal tools` | `rc=0` 7 枚 | `rc=0` **6 枚**（5 假枚＋35r2 的 gofumpt 残留） | 只减不增 |
| `git ls-files '*.go' \| xargs gofmt -l` | 29（编排者更正节，⛔ 本腿未复跑改前那一发） | `rc=123` **27 枚**＝22 `.scratch`＋5 假枚 | 只减不增；`rc=123` 根因具名见 `20-…-wording.md` |
| `PATH=…sherpa-onnx… go test ./cmd/wisp/ -count=1` | ⛔ **未跑**（派单明令；与票面 `AC#4` 冲突，具名报回见下 §6） | ⛔ 未跑 | 欠的读数已具名交回编排者 |

### 减少了哪几枚、⛔ 哪几枚本来就没动（名册级）

- **减少 2 枚**＝`cmd/wisp/panel_inbound_guards_35r3_test.go`、`cmd/wisp/panel_transport_35r2_test.go`
  （工作树那一把与 HEAD blob 那一把都少这两枚；`gofumpt -l cmd/wisp` 少的是 35r3 那一枚，35r2 因 gofumpt 加严规则仍在，见 §5）。
- **⛔ 没动、改后仍在册 1 枚**＝`cmd/wisp/models.go`（工作树那一把与宽域那一把都还在 ⇒ 它是换行符假枚，本腿零字节，`AC#2`）。
- **⛔ 没动、改后仍在册 4 枚**（宽域那一把）＝`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`、
  `internal/risk/provenance.go`、`internal/tools/bridge.go` ⇒ ⛔ 本票写面（`internal/**`），逐枚 `i/lf w/crlf`＝假枚。
- **⛔ 没动、改后仍在册 22 枚**＝`.scratch/**` 下故意脏格式的 tracked 变异拷贝／正控夹具（⛔ 是真码、⛔ 该洗）。
- **⛔ 新增 0 枚**：改后每一把尺的名册都是改前名册的**子集**（逐档对拉过）。

## 4. d22scan 名册级作差（⛔ 裸 diff 会被计时污染，本腿剥掉计时与临时目录号后再比）

```
ruler: sed -E 's/\([0-9]+\.[0-9]+s\)//g; …' pre → /tmp/pre.norm，同形 post → /tmp/post.norm；diff pre.norm post.norm
rc=0 difflines=0        # 原件 logs/d22scan-name-level-pre-vs-post.diff.txt
ruler: grep -c PASS pre / post ；grep -c FAIL pre / post
pre PASS=77 post PASS=77   pre FAIL-line=1 post FAIL-line=1
# 那唯一一行 FAIL 字样出自同一条自测汇总，两把逐字相同：
#   runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77
```
⇒ 新增红 **0 枚**、缺失 PASS **0 枚**；`d22scan: clean - no D22 ban violations`，
`ban #8 cmd/ examined 119 Go files`（改前后同为 119 ⇒ 本票⛔ 增删文件，只改 4 行空白）。

## 5. gofumpt 残留 1 枚（HEAD 既有债，⛔ 本票射程，本腿只量⛔ 修）

```
ruler: $(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp/panel_transport_35r2_test.go      rc=0 count=1（仍在册）
ruler: … gofumpt.exe -d cmd/wisp/panel_transport_35r2_test.go                          # logs/post-gofumpt-d-35r2.txt：93 行／5 hunk
ruler: 同尺跑改前 blob（仓外 $B/blob）                                                   # logs/pre-gofumpt-d-35r2blob.txt：101 行／5 hunk
ruler: 两者剥行号后 diff                                                                 # logs/gofumpt-residue-pre-vs-post.diff.txt：13 行
```
⇒ 那 13 行差额**只有两处**：① 头两行的文件路径不同（一个是仓外 blob、一个是工作树），
② 本腿那一枚 `:804` 对齐 hunk 的正文——gofumpt 改前也报它、改后不报了 ⇒ **抱怨只减不增、新增 0 枚**。
残留内容＝`type ( … )` 成组、`var jsPuncts` 复合字面量换行、相邻 `func` 之间补空行 ⇒ **`gofmt -w` 管不到**，
`AC#1` 只许 `gofmt -w` ⇒ ⛔ 本腿不动；要不要另立一票归编排者裁。

## 6. 两处具名报回（本腿不自裁）

1. **★派单 vs 票面 `AC#4` 冲突（第一句就报）**：派单写「`cmd/wisp` 的整包测试⛔ 跑，那枚包本机要 sherpa DLL，`go build` 就够」，
   票面 `AC#4` 逐字要求 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ -count=1` **改前改后各 ≥2 发取交集、逐名作差＝新增红 0 枚**。
   本腿按派单禁令执行 ⇒ **那一格读数欠着**（欠的是"两把名册作差"，⛔ 是"某格门禁红"：build＋d22scan 都 rc=0）。
   ⇒ 请编排者二选一：补跑那一发，或订正 `AC#4` 措辞。**这条欠账本腿⛔ 自记入台账**（`docs/reports/**` 本腿只读）。
2. **gofumpt 残留 1 枚**（§5）⇒ 归编排者裁要不要另立一票。

## 7. 越界逐笔（尺＝`git show --name-only --format=<sha>`，⛔ 区间尺 `git diff --name-only A..B`）

| 笔 | sha | 名册（全部路径） | 枚数 | 越界 |
|---|---|---|---|---|
| 1（锚） | `24ae75b7` | `.scratch/wisp/probes/298/r1/00-anchor.md`／`…/logs/commit-msg-1.txt` | 2 | **0** |
| 2（唯一碰 `.go`） | `18d43f8e` | `cmd/wisp/panel_inbound_guards_35r3_test.go`／`cmd/wisp/panel_transport_35r2_test.go` | 2 | **0** |
| 3（凭据件） | `c14fd204` | `.scratch/wisp/probes/298/r1/10-ac0-dirt-shape.md`／`20-ac2-ac3-annotation-and-wording.md`／`logs/*`（20 枚 `.txt`） | 22 | **0** |
| 4（本件＋票面追加） | 本笔 sha 由 `git log -1 --format=%h` 现取 | `.scratch/wisp/probes/298/r1/30-gates.md`／`…/logs/commit-msg-4.txt`／`.scratch/wisp/issues/298-…fake-cell.md`（票面追加） | 3 | **0** |

★**为什么必须逐笔、⛔ 区间尺（本腿现量撞上的，不是转述）**：
```
ruler: git log --format='%h %ad %an %s' --date=iso-strict 6414a4bb..HEAD
rc=0 count=7      # 我的 3 笔（当时）之外，还有编排者交错落进去的 4 笔：
                  #   36efbd80 probes(302-a1) 起手锚／48e7b0a4 台账 A817／6da0e688 台账 A818／514d3c84 301-v3 verdict
```
⇒ 任何 `git diff --name-only 6414a4bb..HEAD` 都会把**别人的 4 笔**算到本腿头上（正是派单点名"这条本波刚被验证过一次"那一形）。
⇒ 本表的越界判据＝`for c in <我的 sha 逐枚列出>; do git show --name-only --format= $c; done`（原件＝`logs/ac4-forbidden-and-scope.txt`）。

## 8. 禁区零改动自检（逐条已真跑；原件＝`logs/ac4-forbidden-and-scope.txt`）

```
ruler: for c in 24ae75b7 18d43f8e c14fd204; do git show --name-only --format= $c; done
       | grep -Ec 'models\.go|^\.gitattributes|frontend/|^design/|internal/|^scripts/|\.github/|^docs/|thresholds\.go|allowlist\.txt'
hits=0                      # grep 无命中故 rc=1，那是"零命中"的口径，⛔ 不是失败
ruler: git diff --numstat 6414a4bb HEAD -- cmd/wisp/models.go
（无输出）rc=0              # ⇒ 连别人交错的 4 笔一起算，models.go 在整个区间上都是零字节改动
ruler: git config core.autocrlf
true    rc=0                # 与改前逐字同值，本腿⛔ 写任何 git 配置、⛔ 动 .gitattributes
ruler: git status --porcelain -- cmd/wisp
（无输出）rc=0 lines=0       # 提交后 cmd/wisp 工作树干净
ruler: find .scratch/wisp/probes/298/r1 -type f -size 0 -print
（无输出）rc=0              # ⛔ 0 字节证据件
ruler: find .scratch/wisp/probes/298/r1 -name '*.out' -o -name '*.sh' -o -name '*.ps1'
（无输出）rc=0              # 证据件只 .md／.txt
```
三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／
`internal/perm/ticket90_persist_test.go`）⛔ 出现在我任何一笔名册里 ⇒ 零字节。
SLO 阈值／golden／D43 转移表／C1–C32／D1–D47 ⛔ 未触及；⛔ 为变绿放宽任何断言（本票零断言改动）。
git 纪律：4 笔各带**显式 pathspec 且写在 `$( … )` 之外**（第 1 笔第一次直接用 `git commit -- pathspec` 被 git 拒
`error: pathspec … did not match any file(s) known to git`＝新建件未入库前必须先 `git add -- <显式路径>`，本腿改正后 `rc=0`，如实记这一手）；
⛔ `add -A`／`add .`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`；⛔ push（推送归编排者）；
⛔ 仓内建 worktree；临时件只建不删（`logs/ac3-wide-ruler-fixed.txt` 末节是一发被超时中止的慢尺留下的半截件，
本腿⛔ 删、⛔ 复写，同一读数改由 `logs/ac3-blob-wide-ruler.txt` 交）。
票面追加＝只追加（上面各节原句一字不改，⛔ 改票 292 任何字，⛔ 翻票 35／292／298 任何框）。
