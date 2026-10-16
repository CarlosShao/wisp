# 165-v2 读数件（裁决腿·票 165 `AC#6` 的尺与颜色）

腿＝`165-v2`（非实现者终裁）｜起手锚＝具名提交 `44f51fc8745c3ca551f76dc2c0fe78f93ca2b86d`｜
被裁物＝`internal/agent/plan_present_doorbell_165r3_test.go`（221 行／9378 B，1 枚顶层、0 子测试、0 `t.Skip`）
＋`internal/agent/testdata/golden/plan-present.sse`（28 行／1684 B）｜实现者＝另一条腿 `165-r3`。

## 0. 台面（我自己搭的，⛔ 抄他的颜色）

- 树＝`git archive 44f51fc8 go.mod go.sum internal cmd tools scripts frontend models | tar -x -C /d/tmp/165v2/t1`（**八件套**，D 卷），
  再拷 `build/*.dll`（3 枚）＋`third_party/sherpa-onnx/*.dll`（3 枚）＋`frontend/dist`。
- 非空尺＝`find <树> -name '*.go' | wc -l` ＝ **747**；同射程跟踪件尺＝`git ls-tree -r --name-only 44f51fc8 | grep -E '^(internal|cmd|tools|scripts|frontend|models)/' | grep -c '\.go$'` ＝ **747**；
  两列名册逐枚 `diff` ＝ **0 行**（`logs/tracked-go.txt`／`logs/tree-go.txt`）。⇒ 我这棵树⛔ 少件。
- ⚠ 与实现者分母的**具名差**：他报 744（`c9ba5254`）／746（`f7d6dbc98`），我用同一把尺在 `c9ba5254` 现量 **745**、在 `f7d6dbc98` 现量 **747** ⇒ 他那把尺在两枚锚点上各**少 1 枚**（口径差，⛔ 判据差；他的树已不在，⛔ 复现他那发）。
- `GOTMPDIR=D:\tmp\165v2\gotmp`（先 `mkdir -p`，D 卷）。所有 `go test` 一律 `-count=1 -timeout 30m`；⛔ 跑 `./cmd/wisp` 任何形、⛔ 跑 `go test ./...`。
- 门禁四数只在**工作树**跑（导出树缺 `design/` 会让 d22scan ban #8 射程 0 枚⇒`exit 2` 假红）。

## 1. 两枚基线字节（我 `sha256sum` 现量；每发之后复量）

| 件 | sha256 | 行数 |
|---|---|---|
| `plan-present.sse` | `7d26fa889b59fc3a2b39ff18883ee331f8ff1c06ddb45a5a23bc6820af9079be` | 28 |
| `plan_present_doorbell_165r3_test.go` | `00c32049994c441c8927066bc4d4099c1b2977fc20ee9e602bb98f3670d37b75` | 221 |

两串与他报的前缀（`7d26fa889b59fc3a`／`00c32049994c441c`）逐位相符。**每一发（含 10 枚突变）跑完当场还原，还原后两串全部回到上表**；
最后一次复量记在 `logs/f03-fullpkg-anchor.txt` 同轮的 `sha256sum` 输出里。交件时 `git status --porcelain -- internal cmd` ＝ **0 行**（本腿在编译面零字节）。

## 2. 三发主尺（我自己重跑）

尺口径统一：`RUN=grep -c '^=== RUN'`（含子测试）／`PASS_TOP=grep -c '^--- PASS'`（只数顶层）／`FAIL_TOP=同尺的 FAIL`。

| 发 | 命令（scoped 均为 `-run '^TestAC|^TestPlanPresent'`） | rc | RUN/PASS/FAIL | 红句逐字 |
|---|---|---|---|---|
| f01 正向 | 导出树 44f51fc8，⛔ 改任何字节 | 0 | 2/2/0 | —（`--- PASS: TestPlanPresentCallBooksTheToolCallRow165r3 (0.07s)`） |
| f02 去牙 | `sed -i '12d' internal/agent/testdata/golden/plan-present.sse`（28→27，`"name":"plan.present"` 命中 1→0） | 1 | 2/1/1 | `plan_present_doorbell_165r3_test.go:138: doorbell sink reading = 0 rows for "plan.present" over 0 tool_call rows: the replay never reached internal/agent/loop.go:637-640` |
| f03 整包 | `go test ./internal/agent/ -count=1 -timeout 30m` | 0 | `ok 1.790s` | — |
| f04 整包 `-v`（证常驻） | 同上＋`-v` | 0 | RUN=110／PASS_TOP=86／FAIL_TOP=0／SKIP=0 | —（他那枚在 `logs/f04-fullpkg-v.txt:136` 默认跑到，⛔  build tag、⛔ `t.Skip`） |

⚠ `-run '^TestAC|^TestPlanPresent'` 在本包命中 **2 枚**（他那枚＋别腿的 `TestAC3SpillArtifactLandsPrivate`）＝我 `-list` 现量，⛔ 把 2/2 读成"他一人两枚"。
真名尺＝`grep -n '^func Test' internal/agent/plan_present_doorbell_165r3_test.go` → `117:func TestPlanPresentCallBooksTheToolCallRow165r3`。

## 3. 定向突变（10 发，全在导出树、逐枚还原）

口径收窄：其中**他没做过的形＝7 枚**（M-a／M-a2／M-a3／M-b／M-b2／M-e／M-f）；**复算他的＝3 枚**（M-c-A＝他的 M-A、M-c-B＝他的 M-B、M-g＝他的 M-C）。
⚠ 本腿票尾那条 bullet 写的是"九枚定向突变（他⛔ 做过的形）"——那句口径偏松（M-c-A／M-c-B 本是派单第 3(c) 项点名的复算），
更正记在这里与判语件 §四；派单只许票尾**一条**bullet ⇒ 本腿⛔ 追加第二条、⛔ 改已落那一句。

| 形 | 打什么 | rc | 颜色与红句 |
|---|---|---|---|
| M-a | 用例期望名 `plan.present`→`plan.presentX`（常量 `:38`，夹具不动） | 1 | `:138: doorbell sink reading = 0 rows for "plan.presentX" over 1 tool_call rows` ⇒ 名有牙 |
| **M-a2** | **两头同改**：常量→`plan.presentX` ＋ 夹具 `"name":"plan.present"`→`plan.presentX` | **0 绿** | ⇒ **缺口 G1**：名字整对可换、用例不响 |
| M-a3 | 只改夹具线上名（`plan.presentZ`），常量不动 | 1 | `:138: ... 0 rows for "plan.present" over 1 tool_call rows` ⇒ 读数跟**线上名**走 |
| M-b | 镜像条目 `Resident: true`→`false` | **0 绿** | ⇒ **缺口 G2**：机理见 §4 |
| M-b2 | 镜像档位 `RiskL2`→`RiskL1`（他做的是→L0） | 1 | `:157: row.risk_level = "L1", want "L2" (the frozen tier must survive the trip to the row)`（唯一红） ⇒ `Declared` 有牙 |
| M-c-A | `n := sinkCount165r3(...)`→`n := 1`，牙完好 | 0 绿 | ⇒ 那一枚门**不携带信息**（⛔ 当凭据） |
| M-c-B | 同处常量计数 ＋ `sed -i '12d'` 摘牙 | 1 | `:143: no tool_call row was booked for "plan.present" (rows read back = 0): the call never reached internal/agent/loop.go:637-640` ⇒ 假计数顶不掉缺牙 |
| M-e | `Resident→false` ＋ 描述换成与查询无共同字（`先交付方案待批`） | 1 | `:219: round 1 body: missing "\"name\":\"plan.present\"" in "{\"model\":\"mock-small\",…"` ⇒ "请求里带没带"那一半有牙 |
| M-f | `planPresentCallID165r3`→`call_ppX` | 1 | `:199: tool log entry = {CallID:call_pp1 …}, want a rejected plan.present/call_ppX` ＋ `:220: round 2 body: missing "\"tool_call_id\":\"call_ppX\""` ⇒ 调用号真从线上来 |
| M-g | 镜像档位 `RiskL2`→`RiskL0`＝**复跑他的 M-C** | 1 | 五枚红逐行相符：`:157 …"L0", want "L2"`／`:184 row.decision = "allow", want "reject"`／`:191 row.error_class = "tool", want "permission_denied"`／`:199 …Rejected:false`／`:208: executed tool calls = 1, want 0 (an ungated L2 call must never execute)` ⇒ "工具体一枚没跑"那一格我**自己**有凭据 |

## 4. 两枚字段缺口的机理（具名到行）

- **G2（`Resident` 在本仪器里是装饰）**：`internal/agent/inject.go:44-49` 把目录切成"驻留⇒必注入／非驻留⇒BM25 候选"，而 `bm25Select`（`inject.go:172-224`）只要 `score>0` 就入选。本用例目录**只有 1 枚候选**，查询"帮我把这个项目收拾一下"与描述"先交一份「我要干什么」的清单并停下等点头"共用 unigram `我`/`一`/`下` ⇒ 非驻留照样被注入 ⇒ `:219` 的 `mustContain` 察觉不到档位在换。M-e 证明"注入与否"这件事本身有牙，缺的是"驻留"那一档的区分力。
- **G1（名字可整对换）**：`planPresentName165r3`（`:38`）一枚常量同时喂镜像条目（`:67`）、读数靶（`:136`/`:143`）与请求体断言（`:219`）⇒ 把它与夹具里的 `"name"` 一起改掉，四处自洽、用例仍绿。票面硬约束（导入环：`internal/tools/bridge.go:13` 已 import `agent`）决定了本包⛔ 能 import 真名册，于是"镜像 ⇄ 真名册"这对今天在**任何一枚用例里都没对接**（我的尺：`cmd/wisp/plan_assembled_165r2_test.go`/`plan_wiring_165r2_test.go` 断的是真名册与装配出口，`internal/agent` 这枚断的是镜像，两半零枚相等）。
- **G3（标签错位）**：`:131` 注释写着 `--- THE JUDGMENT: the sink reading must be non-zero. ---` 指在 `n == 0` 那枚门上，而 M-c-A（常量计数）绿、M-c-B（常量计数＋摘牙）红在 `:143` ⇒ 那枚门对牙不敏感；真牙在 `:143`（`rowOfTool165r3` 独立再遍历一遍行）＋`:219`/`:220`/`:199`。`:78-82` 的注释自己承认了这一点，所以记**措辞误导下一位读者**，⛔ 记恒真。

## 5. 落件前的世界（我自己量的，用来裁"派单第 9 行 ⇄ AC#6 原文"那一问）

锚＝`4366c1ca1`（＝`4f2470703^`，落件前一枚；`git ls-tree` 证两枚被裁字节当时不存在）：
- `git grep -c 'withJournal(store)' 4366c1ca1 -- internal/agent` ＝ **8 处／5 件**；夹具分布（`-B3` 现读）＝`tool-then-text` 3／`slow-tool` 2／`six-tools` 1／`disconnect-mid-toolcall` 1／`max-tokens-toolcall` 1。
- 同锚 `git grep -h -o '"name":"[^"]*"' 4366c1ca1 -- internal/agent/testdata/golden | sort | uniq -c` ＝ **echo 18／sleep 7**（`six-tools.sse` 那 6 枚也是 `sleep`）⇒ 表里当时读不出任何一枚名册真名。
- `git grep -F -l 'plan.present' 4366c1ca1 -- '*.sse' '*/testdata/*'` ＝ **0 件**。
⇒ 那两行 sink **走得过**（8 处），但**零枚用例带名册里的真工具名**打到它＝AC#6 原文 ① 与 `165r3-dispatch.txt:9` 逐字（那句带 `plan.present` 限定词）同形；裸写"零枚用例打到那两行"现量不真。

## 6. 门禁四数（只在**工作树**，逐枚落 rc）

| 门 | rc | 读数 |
|---|---|---|
| `sh scripts/d22scan.sh` | 0 | `clean - no D22 ban violations`；ban #8 射程 `design/`=39、`frontend/`=116、`internal/`=577、`cmd/`=139（工作树 design/ 尚有余件⇒射程非 0，⛔ 触发 GUARD `exit 2`） |
| `sh scripts/check-path-length-budget.sh` | 0 | VERDICT GREEN：跟踪路径 14330／超预算 60／在册 60／未在册 0 |
| `go vet ./internal/agent/` | 0 | 无输出 |
| `gofumpt -l`／`-d` 两枚被裁字节 | 0／0 | 零命中／零差异（真身 `/d/work/base/gopath/bin/gofumpt.exe`） |

## 7. 实现者名册的写面尺（逐枚 `git show --name-only`，10 枚）

`4f2470703` 夹具／`f7d6dbc98` 用例／`d498a73c` 读数件／`78fe125f` logs 15 枚／`3e4ca0e3` 证据件／`c3f55651` 票尾 bullet／`f6c76b91` 正文源 3 枚／`77553c9e` 读数＋证据／`16d4b225` logs 1 枚／`6b18b1d4` 读数件。
全部路径落在 `{internal/agent/**, .scratch/wisp/probes/165/r3/**, docs/evidence/s1/165-doorbell-instrument-r3.md, 票 165}` ⇒ **越界命中 0 枚**。
`c3f55651` numstat＝`2 0`（只追加）；票框现量＝`grep -c '^- \[x\]'` 5／`'^- \[ \]'` 1（AC#6 仍空，他一枚没翻）。logs 目录现量 **16** 枚。
