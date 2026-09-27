# 票 164 · AC#2 ＋ AC#3 实现程证据件（写码位 164-r1）

- 派单＝`.scratch/wisp/dispatches/2026-09-27-185x-impl-164-r1-output-leg-with-recoverable-pointer-held-until-171-v2.md`
- 票面＝`.scratch/wisp/issues/164-background-jobs-cannot-be-read-fix-the-roster-then-add-the-output-leg.md`（含 09-27 18:5x 五条编排者定案）
- 只读设计核＝`docs/evidence/s1/164-task-output-design-core-c1.md`（先读的，未冷搜索）
- 本格性质＝写码位；**只 commit、不 push**；AC 框**不由实现程勾**。
- 证据分档：〔现跑〕＝本程同树跑过；〔第二见证〕＝引他处读数＋锚点；〔推断〕＝无读数支撑。

---

## 1. step-0 五件〔现跑〕

```
date '+%Y-%m-%d %H:%M %z'                       -> 2026-09-27 19:16 +0800
git rev-parse --abbrev-ref HEAD                 -> dev            （= 派单要求的分支，未停手）
git rev-parse HEAD                              -> f206e9f2210a9e9a41fc4270ba1796afa2e379b1   ← 本程 step-0 锚
git status --porcelain -- internal/tools/ internal/agent/ cmd/wisp/   -> 空（三枚写面干净，可开工）
```

基线读数（改前，逐包、非全仓）：

```
go test -count=1 -v ./internal/tools/ ./internal/agent/    -> EXIT=0
  ok  github.com/CarlosShao/wisp/internal/tools   18.860s
  ok  github.com/CarlosShao/wisp/internal/agent    2.369s
```

**口径（两数不同义，别混着报）**：

| 口径 | 命令 | 读数 |
|---|---|---|
| 含子测试的 RUN 行数 | `grep -c '^=== RUN' <log>` | **238** |
| **顶层** PASS 行数（子测试带缩进、不计入） | `grep -c '^--- PASS' <log>` | **172** |
| 顶层 FAIL | `grep -c '^--- FAIL' <log>` | **0** |
| 顶层 SKIP | `grep -c '^--- SKIP' <log>` | **0** |
| 顶层用例名册（差集基准） | `grep -oE '^--- (PASS\|FAIL\|SKIP): [A-Za-z0-9_/]+' <log> \| sed 's/^--- [A-Z]*: //' \| sort` | **172 枚**，存 `/tmp/164r1-baseline-names.txt` |

⇒ 基线＝两包全绿、172 枚顶层用例、零 skip。

---

## 2. 本程**没**测什么（照实记，别当查过）

1. **没跑** `probes/154/gate-clauses.sh` 与 `probes/161/r6/flip-declaration.sh`（派单硬令：那族尺的行为归票 171，且 flip 一跑就脏 tracked 日志）。关于 G3/G5/G6/G7 的一切都是**读脚本与读设计核**，不是本程读数。
2. **没跑**全仓 `go test ./...`（同树躺着别人的未提交件）；门禁只 scope 到 `./internal/tools/ ./internal/agent/`。⇒ **本程不声称全仓绿**。
3. **非 Windows 那一支**：见 AC#3 段末的 linux 容器读数；`internal/tools` 里带 `//go:build !windows` 的文件（`platform_other.go`、`staging_live_other.go`）**不参与本机 Windows 编译**，本机绿不等于它们绿——单独在容器里量过的那一发才算。
4. **AC#4 那一格今天没有被测对象**（派单末条＋票面 18:5x 定案⑤）：生产侧零后台 spawn 点 ⇒ 本程**没有**为它写判据，也没有写一发永不响的装饰。
5. **没查**注册表有没有用拼接/常量表注册 `task.*`（AC#1 裁定段 `:24` 自己留的那个洞，本程没补）。本程只新增字面量名字。

---

## 3. AC#2〔现跑〕——"今天起一个后台东西、再读它的输出"这一发能不能做到

**判据形状**：**调用级**，不是"注册表里有没有这个名字"（设计核 §5 第 2 条＋派单 §1 的提醒）。台件＝真桥 `Bridge.Execute`（`internal/tools/bridge.go:243`，即环路调用的那一枚 choke point），roster＝今天生产注册的那一枚（`cmd/wisp/run.go` 的 `tools.BuiltinFSEntries`）。注入面只有既有接缝，无 mock 代替真桥（AGENTS §1.3）。

用例＝`internal/tools/task_output_ac2_before_test.go::TestTaskOutputAC2BeforeLegIsUnreachable`。

**改前那一发的原样读数**（可复算命令在下面）：

```
go test -count=1 -v -run TestTaskOutputAC2 ./internal/tools/

--- PASS: TestTaskOutputAC2BeforeLegIsUnreachable (0.00s)
    AC#2 BEFORE verbatim: IsError=true ErrorClass="tool" Truncated=false \
      Text="未知工具 task.output，可用工具见 list_tools"
```

⇒ 落到 `bridge.go:247-253` 的 `Lookup` 失败分支：`IsError` 真、`error_class=tool`（D37 可自纠，任务继续）、**不是** Go error（拒绝不是故障）。

**结论：做不到 ＝ 本票成立，派单那句判断未被推翻。** 同一发请求在改后必须返回真 `Result.Text`（见 AC#3 段的 `TestTaskOutputAC2AfterLegIsReachable`），那一对照就是"由不可读到可读"的基准。

---

## 4. 派单前提逐条现核（都是未验证断言，核完才用）

| 派单／票面的前提 | 现量结果 | 凭据 |
|---|---|---|
| `fs.read` 参数只有 `{path, max_bytes}`、**没有偏移** | **核对为真** | `internal/tools/fs.go:113-116`（`fsReadArgs` 两枚字段）、`:119-122`（schema 只有这两个 key） |
| `fs.read` 单次上限 256 KiB | **核对为真** | `internal/tools/fs.go:55` `defaultMaxReadBytes = 256 * 1024` |
| `bridge.go:247` 那个 `Lookup` 失败分支 | **核对为真**，行为＝`未知工具 <name>，可用工具见 list_tools` ＋ `OutcomeKindNotFound` ＋ `err == nil` | `internal/tools/bridge.go:247-253`，读数见 §3 |
| `RunAsync` 生产零调用点 | **核对为真**：非测试命中只有两行**注释**（`internal/agent/loop.go:320` 自己的文档注释、`internal/tools/bridge.go:663` 的一句注释），**生产调用点 0 枚** | 〔现跑〕`grep -rn "RunAsync" --include=*.go internal cmd \| grep -v _test \| grep -v "func (l \*Loop) RunAsync"` |
| `PLAN.md:431` 写着可用 `fs.read` 再读 | **核对为真**，原文逐字含"模型可用 `fs.read` 按需再读（复用 D10 机制）"，且同一行就规定截断形状＝**头 500 token ＋ 尾 200 token ＋ 总长度 ＋ 文件路径** | 〔现跑〕`git show HEAD:docs/PLAN.md \| awk 'NR==431'` |
| `task.output` 那两行契约文本已落 | **核对为真** | `docs/PLAN.md:2564`、`docs/specs/SPEC-07-tools-and-plugins.md:71` |
| （新问）**artifacts 目录在不在 `[fs] allowed_dirs` 里**——决定"指针能不能被 `fs.read` 走通" | **默认不在**：`[fs] allowed_dirs` 全仓没有默认值可查（只有 schema 字段与"consumed"登记），`cmd/wisp/run.go:324-327` 又只用 `cfg.FS.AllowedDirs` 构造 C26 resolver | 〔现跑〕`grep -rn "allowed_dirs" internal/config/ cmd/ \| grep -v _test` ⇒ 三行：`manager.go:377`（增删方向）、`schema.go:472`（字段，无默认）、`unwired.go:129`（"consumed: cmd/wisp/run.go feeds the C26 canonicalizer"） |

⚠ 最后一行是**既有账**（票 164 之外的 C26 allowlist 语义），本程**没有**顺手修，也没有为它改判据；它决定了本程把"读回全文"这一腿判成"由 `task.output` 自带指针＋由 `fs.read` 在 allowlist 命中时走通"，见 AC#3。

**顺带结掉设计核 §3.1 那处"引文坑"**（它明说没做字节级核对、留给了后程）：〔现跑〕`git show HEAD:docs/PLAN.md | awk 'NR==431'` 打出来的是
`全文落 \`artifacts<TAB>ool-output-<id>.txt\``——**TAB 确实在打印层吃掉了一个 `\t`**，源文件写的是反斜杠路径分隔符（与 `internal/agent/spill.go:20` 的 `<artifacts>\tool-output-<encoded id>.txt` 同源）。⇒ 设计核那句怀疑**成立**；引用 `:431` 时那段路径**不能当逐字凭据**，本程引的是它的**截断形状**（头 500／尾 200／总长／路径），不是那枚路径字面量。

---

## 5. AC#2 的"改后"那一发〔现跑〕（与 §3 同一枚请求、逐字对照）

| | 请求 | 读数 |
|---|---|---|
| 改前（roster＝今天的 `BuiltinFSEntries`） | `task.output {"task_id":"bg-1"}` | `IsError=true ErrorClass="tool" Truncated=false Text="未知工具 task.output，可用工具见 list_tools"` |
| 改后（roster＝fs 六枚 ＋ `BuiltinTaskEntries`） | 同一枚 | `IsError=false ErrorClass="" Text="后台任务打印的第一行"` |

用例＝`task_output_ac2_before_test.go::TestTaskOutputAC2BeforeLegIsUnreachable` ＋ `task_output_leg_test.go::TestTaskOutputAC2AfterLegIsReachable`。
**"由不可读到可读"的比对基准就是这两行**，两边都在同一棵树上跑过，且**请求字节相同**——差别只在注册表里有没有那枚工具。

---

## 6. AC#3 —— 落点、形状与两向反向判据

### 6.1 落点（派单 §5 第 5 条要的自证）

| 文件 | 新增／改 | 内容 |
|---|---|---|
| `internal/tools/task.go` | **新建** | `task.output`（C1 Tool）＋ `TaskRoster`（进程内任务表）＋ `TaskDeps`／`TaskOutputDecl()`／`BuiltinTaskEntries()` |
| `internal/tools/task_output_leg_test.go` | **新建** | AC#3 判据 9 枚（含两枚反向判据的承接者）＋ AC#2 改后那一发 |
| `cmd/wisp/run.go` | 改既有 | `agentRuntime` 加 `tasks *tools.TaskRoster` 一枚字段；fs 注册循环之后 append `BuiltinTaskEntries`（`git diff --numstat`＝**23 增／5 删**，〔现跑〕那 5 枚删除逐行看过＝**全是 gofumpt 因注释行断开头尾对而重排 `spec/cfg/mgr/paths/store` 的对齐**，五枚字段原样在，零内容丢失） |
| `internal/agent/**` | **零字节** | 见下 |

**落点自证（派单点名的那枚命令）**：

```
git -c core.quotePath=false diff --name-only f206e9f2210a9e9a41fc4270ba1796afa2e379b1 HEAD -- internal/agent | wc -l   -> 0
git -c core.quotePath=false diff --name-only f206e9f2210a9e9a41fc4270ba1796afa2e379b1          -- internal/agent | wc -l -> 0
grep -rnE '^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID' internal/agent | wc -l                                  -> 0
```

⇒ **`internal/agent` 一字未动**，`Loop` 上**没有**新增任何"收 `taskID` 的导出方法"——`gate-clauses.sh` 的 **G3 腿**（`want G3 quiet`／`want_n 0`，为 `Q-56` 而立）射程内那枚形没有被我造出来。入口全在 tools 侧（定案④）。
⚠ 本程**没有跑** `probes/154/gate-clauses.sh`（派单硬令），上面第三行是那把尺的 pattern 本体的直接计数，不是它的读数。

### 6.2 形状（照 `PLAN.md:431`／`:2564`，不新造）

`task.output` 参数＝`{task_id}`，**没有** offset／cursor／page（定案①：不发明第二套分页读取）。三种形状，只有中间一种允许安静：

1. **查无此 id** → `IsError`，文案点名 id 并明写"这是『没有这条记录』，不是『任务没有输出』"（定案②：绝不空返回）。
2. **短输出**（≤4000 token） → 全文返回，`Truncated=false`，不需要指针。
3. **长输出** → 头 **500** token ＋ `[…输出已落文件：省略 N 字符，总长 B 字节 / 约 T token，全文见 <path>…]` ＋ 尾 **200** token，`Truncated=true`。桩文本措辞照 `internal/agent/spill.go:141` 那一句（同一枚 D15 机制，不给模型看两种方言）。
   长输出而**宿主没落文件**时（`ArtifactPath==""`）→ 同一形，但把指针位换成 `后半段不可找回`，**响**，不静默。

能力声明：`Capabilities`/`Needs` 都空 ⇒ 与 `PLAN.md:2564` 第 4 列那个 `—` 同形，**C3 一枚没碰、没加第 12 枚**（定案③）。`PathParams` 也空：模型唯一能传的字符串是名册的 key，不是路径，所以 C19 没有路径可判、C26 不需要规范化 ⇒ 那条"按 taskID 去 artifacts 目录拼路径"的地雷**从形状上就不存在**（`task.go` 全文 `import` 里没有 `path/filepath`、没有 `os`）。

### 6.3 反向判据（本格不是装饰的凭据）——两枚变异都**真跑过**

| 变异 | 摘掉哪一处 | 必须红的那一枚 | 〔现跑〕实测 |
|---|---|---|---|
| **只截不指** | `task.go` 桩文本里 `，全文见 %s` 那一味 | `TestLongOutputPointerRecoversEveryByte` | **FAIL**：`the stub must carry a recoverable pointer`（整包 rc=1）。顺带 `TestPointerPast256KiB…` 也红（它同样吃那枚指针），`TestTruncationShapeIsTheD15Triple` 仍绿——红得**各归其因**，不是一片糊。 |
| **截断形状与 D15 不符** | `refTaskSpillHeadTokens` 500 → **400** | `TestTruncationShapeIsTheD15Triple` | **FAIL，且全库只有它红**（其余 9 枚 PASS，rc=1）。⇒ 期望值写在测试里的是**字面量 500×4／200×4**，不是引用产码那枚常量——本程一开始确实引用了常量（那样改常量测不出来），发现后改成字面量，这一枚才真的抓得住。 |

两枚变异跑完**都已逐字改回**（`refTaskSpillHeadTokens = 500`、桩文本含 `全文见 %s`），改动没有留在树上。

`TestLongOutputPointerRecoversEveryByte` 断的是**读回来**而不是"提了个路径"：桩里那枚路径 (1) 直接 `os.ReadFile` 必须逐字节等于全文；(2) 再经**真 `fs.read` 调用**读一次，仍要逐字节等于全文；(3) 宣布的 `总长` 必须等于**全文**长度、不能等于桩自己的长度；(4) `省略 N` 与实际保留的头尾字节数对得上账。

### 6.4 "产物超长时后半段今天读不到"——本程的处置＝**断言它响**

`TestPointerPast256KiBIsNotFullyReadable`：落一枚 300 KiB 的 artifact，`task.output` 给出指针，再用真 `fs.read` 去走那枚指针 ⇒ 实测**只回前 256 KiB**（`len=262144`、`Truncated=true`），换 `max_bytes=9999999` 也一样（`fs.read` 的参数里没有偏移，`internal/tools/fs.go:113-116`）。

**为什么选断言、不选明写**：这条限制是**可量的行为**，写成断言它就成了一枚 canary——`Q-59` 哪天批了给 `fs.read` 加偏移，这一枚**会红并逼着回来改本票那一格**；写进证据件的散文不会。
⚠ 同时**没有为了让判据好看去扩 `fs.read`**：该文件本程一字节未动（§6.5 的名单可证）。

### 6.5 契约轴（AC#5 的形，但**不勾框**）

**本程只拥有 5 枚路径**：`internal/tools/task.go`、`internal/tools/task_output_ac2_before_test.go`、`internal/tools/task_output_leg_test.go`、`cmd/wisp/run.go`、本证据件。
〔现跑〕把这 5 枚对 `^(docs/specs/|docs/PLAN\.md|internal/risk/|internal/panel/|internal/agent/|internal/observe/thresholds\.go|tools/d22scan/|frontend/|design/|\.scratch/wisp/probes/|docs/reports/|\.github/)` 求交集 ⇒ **0 命中**；`golden|thresholds\.go|allowlist\.txt` ⇒ **0 命中**。

⚠ **口径提醒（别拿树级读数当本程读数）**：`git diff --name-only f206e9f -- <那族路径>` 在这棵**共享树**上会给出 22 枚变更，那**全是编排者与别的过程留下的未提交件**（`design/**` 的 16 枚删除、`docs/reports/HANDOVER.md`、`pending-and-issues.md`、`.gitignore`、`probes/**`…），**没有一枚是本程写的**——本程一枚没 add、没还原、没补完。派单要求的"以我 owns 的路径求交集"才是这一格的凭据。

---

## 7. 门禁读数（改前／改后各一遍，全部 scope，无全仓 `./...`）

| 门 | 命令 | 改前 | 改后 |
|---|---|---|---|
| 逐包测试 | `go test -count=1 -v ./internal/tools/ ./internal/agent/` | rc=0；`ok tools 18.860s`／`ok agent 2.369s`；RUN **238** ／顶层 PASS **172** ／FAIL **0** ／SKIP **0** | rc=0；`ok tools 16.519s`／`ok agent 1.932s`；RUN **248** ／顶层 PASS **182** ／FAIL **0** ／SKIP **0** |
| 名册两向差集 | `comm -3 <(基线名册) <(改后名册)`（名册＝`grep -oE '^--- (PASS\|FAIL\|SKIP): [A-Za-z0-9_/]+' \| sed 's/^--- [A-Z]*: //' \| sort`） | — | **左栏 0 枚**（无 vanished／无被 panic 吞掉的）、右栏 **10 枚**新增，全是本程用例 |
| D22 静态扫 | `sh scripts/d22scan.sh` | — | **rc=0**，`clean - no D22 ban violations`；本票相关两行：`scope ban #7 internal/tools/ examined 20 production Go files`、`ban #8 internal/ 425 Go files` |
| d22scan 自测 | `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | — | **rc=0**，`top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0` |
| 格式 | `gofumpt --version` ＋ `gofumpt -l <本程 4 枚 .go>` | — | **`v0.12.0 (go1.27.1)`**（先 `export PATH="$PATH:$(go env GOPATH)/bin"` 现跑）；`-l` 首轮列出 2 枚（`task_output_leg_test.go`、`run.go`）⇒ `gofumpt -w` 之后 **unformatted=0** |
| **跨平台（linux 那一支）** | 宿主交叉编译再进容器**真跑**：`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -c -o "$OUT/tools.test" ./internal/tools/` ＋ `docker run --rm -v "$OUT:/out" golang:1.27 /out/tools.test -test.run '…' -test.v` | — | **10/10 PASS、`PASS`、DOCKER_RC=0**。⚠ 挂载**先证后跑**：第一次 `cygpath`/`sed` 写坏、容器里 `ls: cannot access '/out/tools.test'` ⇒ 那次是**假绿形状**、没算读数；重来那次的容器内 `ls -l /out/tools.test` 打出 `-rwxrwxrwx … 24136577` 才算 |

⚠ **"本机绿＝跨平台绿"这句本程不说**：新用例**没有** `//go:build` 标签、不是 `*_other_test.go`，所以 linux 容器那一发是"真跑过"而不是"推出来的"。AC#4 那一族（`internal/proc/jobscope_windows.go`）本程没碰，**没有**任何 Windows-only 断言被新增。

**没跑**：`probes/154/gate-clauses.sh`、`probes/161/r6/flip-declaration.sh`（派单硬令；那族尺的行为归票 171）。
**`go test ./cmd/wisp/` 以 `0xc0000135` 收场**＝**既有环境事实**，不是本程引入：〔现跑〕`go test -count=1 -run 'ZZZNoSuchTest' ./cmd/wisp/` 一枚用例都不选也 0xc0000135 ⇒ 加载期就没过；台账 `docs/reports/pending-and-issues.md:2207` 早记过同一形状根因＝缺 `sherpa-onnx-c-api.dll`。本程的 `go build ./cmd/wisp/` = **BUILD_OK**，注册那一发是编译期过的、不是靠测试绿灯。

---

## 8. 提交记录

**AC#2 那一格**：

```
$ git log --oneline -1
bde17aca 164-r1 AC#2: 把"今天读不到后台任务输出"量成调用级读数
$ git -c core.quotePath=false show --name-only HEAD
docs/evidence/s1/164-task-output-impl-r1-ac2-ac3.md
internal/tools/task_output_ac2_before_test.go
```

**AC#3 那一格**：记在下一枚 commit 之后（本件与它同程提交；未跑出的读数不预先引用）。

---

## 9. 伪授权两栏 ＋ 没做的事

- **被权限系统拒绝的调用**：**0 次**（本程从头到尾没有一次 tool call 被拒；没有"假设成功"往下写的段落）。
- **跑过的删除命令**：**0 次**（`rm`／`git clean`／`git checkout .`／`restore` 一律没跑；两枚变异是**编辑—跑—编辑回去**，没有删文件）。
- **AC 框**：**一枚没勾**。票面五个框的勾归编排者，且还要另派非实现者做对抗验收表（`AGENTS §0.3`／`SPEC-12 §4.3`）。
- **凭据值抄录**：本程零次；测试里出现的都是自造的临时 task_id 与 `t.TempDir()` 路径。
- **AC#4**：**没有**硬造判据。今天连被测对象都没有（生产零 spawn 点，§4 已量），按定案⑤处理：`next=` 里写清可测时刻晚于前两格。

## 10. `next=`

1. **`task.output` 现在"在册＋有实现＋生产无人喂"** —— 这是 AC#1 那笔硬账的**新形状**（从"零实现"前进到"零写者"）。下一格（票 163 那批／AC#4）要做的是**后台起跑口本身**：谁 `RunAsync`、跑完把 `Text` ＋ `spill.go` 那枚 artifact 路径 `roster.Record(taskID, …)`。做完 AC#2 那句"起一个后台东西、再读它的输出"才在**生产**里成立（今天只在桥级判据里成立）。
2. **AC#4 那一格**：先把"还在写的后台尾巴"造出来（真 spawn ＋ 真写者），**再**回答"未修码上响不响"；红句要点名"取消后仍落了几字节"，形制可照 `internal/tools/fs_edit_ac4b_kill_windows_test.go:36`（"no t.Skip in this file"）。它的可测时刻**晚于本票前两格**。
3. **`Q-59`（`fs.read` 加偏移）批与不批都会咬本票**：批了 ⇒ `TestPointerPast256KiBIsNotFullyReadable` 会红，那是**设计好的 canary**，届时把 §6.4 那一格从"断言限制在"翻成"断言能续读"；不批 ⇒ 本票"后半段读不到"应升成 `PLAN.md §7` 那条五字段 DEFERRED 行（票面定案①那一支），别让它停在"在册＋无人认领"。
4. **artifacts 目录默认不在 `[fs] allowed_dirs`**（§4 末行现量）：不修，但 `task.output` 的指针能不能被模型自己走通，取决于它——属既有账，另立一枚。
