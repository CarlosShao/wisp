# 票 174 落地腿 174-r1 — the allowlist arm told the model there was ONE road back, and there are two

- 角色：写腿（实现者）。裁决者不是实现者 ⇒ 本件**不判自己过**、禁 未翻任何 AC 框、禁 未改票名/`Status:`/`-done`、禁 未碰台账 `docs/reports/**`。
- 起手锚：`b482968f`（`git rev-parse HEAD`，09:11 现跑）。终态 `git status --porcelain -- internal/ cmd/`＝`M internal/tools/task.go` ＋ `?? internal/tools/task_output_pointer_two_roads_174r1_test.go`（其余 zero）。
- 写面：`internal/tools/**` 两枚文件。零字节：`internal/tools/task.go` 以外任一 `internal/**`、`cmd/**`、`internal/risk/**`、`tools/d22scan/**`、`thresholds.go`、`[fs]` 默认值、`docs/PLAN.md`、`docs/specs/**`、`frontend/**`、`design/**`、`go.mod`/`go.sum`/`.gitignore`/`ci.yml`。

## 1. 开工前三条（禁 跳步＝空跑）

1. **`git status --porcelain -- internal/tools/ internal/agent/ cmd/`＝0 行**（09:12 现跑）⇒ `internal/tools/**`、`internal/agent/**`、`cmd/**` **无别人未提交改动**，不触发停手。另测全树 tracked-dirty 只 1 行＝`.scratch/wisp/issues/115-...-done.md`（不在我写面）。
2. **票面前提新鲜度（尺 v4 的v4 第8条）逐条现跑**：见 节2，五条全**未过期**。
3. **现量记录**：见 节2 表格（`AllowedDirs` 逐条命中行、`spill.go` 阈值与回执逐字）。

## 2. 票面 节现量 逐条复量（尺＝现跑，非转述）

| 票面断言 | 现跑尺 | 读数（09:12–09:20） | 过期？ |
|---|---|---|---|
| `run.go` 授权根只从 `cfg.FS.AllowedDirs` 拷 | `grep -n AllowedDirs cmd/wisp/run.go` | `:510` `:512` 两行（`allowed := make(...)`＋`for _, d := range`）＝**唯一来源** | 否（行号由票面 `:327-331` 漂到 `:510-514`） |
| spill 落盘目录＝`dataDir/artifacts` | `grep -n artifacts cmd/wisp/run.go` | `:1090` `ArtifactsDir:`／`:1186` `NewSpiller(filepath.Join(rt.spec.dataDir,"artifacts")...)` | 否（行号漂移） |
| 同一 `dataDir` **从未**写进 `allowed` | `grep -n allowed cmd/wisp/run.go` | 全文件只 `:510` `:512` `:514` 三处，**零** `dataDir`→`allowed` 写入 | **否** |
| `paths.go` 那句"空名单＝什么都不授权" | `grep -n 'EMPTY allowlist' internal/tools/paths.go` | `:63` 逐字 `// returns the adapter. An EMPTY allowlist authorizes nothing: every fs call` | 否（票面写 `:50-52`） |
| `SPEC-03` 的 `allowed_dirs[]=[]` | `grep -n allowed_dirs docs/specs/SPEC-03-config-secrets-envs.md` | `:35` 逐字含 `allowed_dirs[]=[]` | 否 |
| `task.go` 指针支只看 `ArtifactPath != ""` | `grep -n ArtifactPath internal/tools/task.go` | `:547` `if rec.ArtifactPath != ""`（**条件里零 stat**） | 否（票面写 `:227`） |

**结论：五条全新鲜，无过期前提 ⇒ 不停手。**

**`AllowedDirs` 全签名（逐条，禁 只报计数）**：`git grep -n 'AllowedDirs' HEAD -- internal/config internal/tools` 现跑，命中分三类——
`internal/config/allowdirs.go`（`:40` `:42` `:46` `:76` `:81` `:82` `:86` `:91` `:95` `:109` `:124` `:125` `:139` `:140` `:141` `:142` `:143`）＝**访问器与写入腿本体**；`internal/config/schema.go:486` `AllowedDirs []string \`toml:"allowed_dirs"\``（**无 `default:` tag**，`:485` 注释逐字 "Additions are loosening"）；`internal/config/manager.go:468`（`setDirection` 判宽松/收紧）；其余全在 `_test.go`（`boundary_test.go:148`／`manager_223_test.go:76/143/178`／`manager_test.go:173/194/215/229`／`migrate_test.go:92/93`／`writeguard_226_test.go:63…410`／`writeguard_handedit_227_test.go:95/96`／`internal/tools/refusal_20r2_test.go:199/213/239/240/290/310/329`）。⇒ **默认值是空**（`schema.go:486` 零 default tag ⇒ 零值切片），与我这一腿的前提一致。

**`spill.go` 截断阈值与回执逐字**（`internal/agent/spill.go`，现读）：
- 阈值：`:142` `if ApproxTokens(out.Text) <= s.b.SpillTokens { return out, nil }` ⇒ 走 `Budgets.SpillTokens`（`PLAN.md:431` 的 4000 token 经 scale）；硬帽 `:109-119` `CapRaw` 走 `s.b.RawOutputCapBytes`。
- 回执模板逐字（`:186`）＝`"%s\n[…输出已落文件：省略 %d 字符，总长 %d 字节 / 约 %d token%s，全文见 %s…]\n%s"`；notice 槽在 `约 %d token` 与 `，全文见 %s` **之间**（`:181-184` 注释具名）。
- 该腿 allowlist 支逐字（`:232-234`）＝`"；注意：这条路径现在读不到，它不在你被授权的目录范围内：" + "fs.read 会要一张 L2 卡，没人批就是直接拒；回来的路有两条——" + "要么用户把所属目录加进 [fs] allowed_dirs，要么批下那一张卡"` ⇒ **这条腿早就说了两条路**。

## 3. 改了什么（一处产码，一处新判据）

### 3.1 产码：`internal/tools/task.go` 的 `pointerNotice` allowlist 支（1 处，8 增 / 2 删 + 7 行注释）

改动前（**未修码逐字**，现跑捕获）：

```
；注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来
```

改动后（现跑捕获）：

```
；注意：这条路径现在读不到，它不在你被授权的目录范围内：fs.read 会要一张 L2 卡，没人批就是直接拒；回来的路有两条——要么用户把所属目录加进 [fs] allowed_dirs，要么批下那一张卡
```

**为什么必须改**：`[fs] allowed_dirs` 在今天是**风险判级输入、不是执行硬边界**——`FSDeps.open`（`internal/tools/fs.go:93-105`）只 `Canonicalize`、**执行腿不查授权根**，所以一张被批的 L2 卡足以把根外路径真读走。旧句"fs.read 会被拒 … **才**读得回来"里那个"才"把"没人批的情形"写成了唯一情形 ⇒ **是同一族半真话**（本票要清的就是这一族）。而 `internal/agent/spill.go:232-234` 自 `174-c2` 量出这层不对称起就一直在说两条路 ⇒ 改前**同一个判定在两条腿上给了两个答案**。改后两腿逐字同形。

**边界守住**：(i) 只读、零新增能力、零授权面变化——不改 `[fs]` 默认值、不改 C26 判定语义、`internal/risk/**` 零字节；(ii) 判"在不在根里"仍**只**经 `TaskDeps.Paths.Canonicalize`→`InAllowlist`（C26 本尊），无 `filepath.Clean|Abs`（`grep -n 'filepath\.' internal/tools/task.go` 仍只有 `:47` 那句**注释**，零实调）；(iii) 指针照样给（`全文见` 在场），模板 `%s…` 形状未动。

### 3.2 新判据：`internal/tools/task_output_pointer_two_roads_174r1_test.go`（新件，4 枚）

复用的既有无改动：`judgedTaskBridge`/`callTaskOutput`/`readThroughBridge`/`mustWriteFile`/`noticeLead`/`announceOf`（`task_output_pointer_notice_test.go`）、`gateSpy`（`bridge_test.go:71`）、`tempCanonical`/`mustCanonical`（`helpers_test.go:83/89`）、`asciiRun`（`task_output_leg_test.go:435`）、`pointerRe`（`:155`）。**既有判据文件零枚改动。**

注入面（AGENTS 1.3）：真 `Bridge`＋真 C26 `PathCanonicalizer`＋真 `fs.read`＋真 `TaskRoster`＋盘上真文件＋真 `Gate` 接缝（`GateFuncs`/`gateSpy` 走 `Options.Gate`）。零 mock 顶替路径判定。

## 4. 四枚用例逐名 + 今天为什么响

| 用例 | 承的是哪一问 | 今天为什么响（=它在**修好之后**才有意义的那一发） |
|---|---|---|
| `TestPointerOutsideAllowlistNamesEveryRoadBack174r1` | allowlist 支必须**两条路都点名**，且**不许**再留"只有一个办法" | 断言 `[fs] allowed_dirs` **和** `L2 卡` 两块标记同在，并**禁** `才读得回来` 这个"唯一办法"字面。改前：`L2 卡` 缺失 ＋ `才读得回来` 在场 ⇒ 双红（见 节5 变异 M1） |
| `TestPointerSecondRoadBackIsRealOnThisBridge174r1` | 第二路**是不是真的**（承重：句子准不准由事实兜底） | 同一桥／同一路径／同一 C26 两臂对拼：`NoGate` ⇒ `isError=true level=L2 class="user_rejected"`；`gateSpy{approveAns:Allow}` ⇒ **`isError=false`、整 20000 字节读回**。第二路一假，这句文案立刻成新谎 ⇒ 该臂转红（节5 变异 M3 实证） |
| `TestPointerSecondRoadBackIsBookedAsACard174r1` | "**会要一张 L2 卡**"是量出来的机制，不是修辞 | 同一发读走 `PendingApproval` **恰 1 次**、`PendingWindow` **0 次**，且卡片 `Reason` 逐字含 `授权目录之外`（现跑读数在 节5） |
| `TestPointerTwoRoadNoticeLeaksNoPathOrRoots174r1` | 句变长了，**不许**顺势漏：路径原文／根列表／C26 自身错误前缀 | 两块真根目录 ＋ 真路径，扫 notice 区（从首个 `注意：` 切到 `，全文见 `）禁 `filed 字面量`／`rootA`／`rootB`／`tools:`／`risk:`。改后 notice 区逐字无一项（节5 变异 M2 证明它有牙） |

**反向判据**（节5）已做三发；**恒真判据／`t.Skip`／放宽断言换绿：零枚**（`grep -c 't.Skip' 新件`＝0，`grep -c '//go:build' 新件`＝0 ⇒ 全平台真跑，不需 `portable-tests.sh` ledger 行；四个新名在 `scripts/portable-tests.sh` 里命中数各 **0**）。

## 5. 反向判据三读数（逐字，禁 只报枚数）

**先证变异落地**（每发都贴 `grep -n` 原文 ＋ `go build` rc）：

### 变异 M1 — 把新句整句摘回旧句（本腿的承重那一跳）

- 落地证明：`grep -n "才读得回来" internal/tools/task.go` → **`:859`**；`GOFLAGS= go build ./...` → **rc=0**。
- **变异红（逐名红行）**：`FAIL`，恰 **1 枚**：

```
--- FAIL: TestPointerOutsideAllowlistNamesEveryRoadBack174r1 (0.01s)
  task_output_pointer_two_roads_174r1_test.go:102: the allowlist notice must name "L2 卡", got "[…输出已落文件：省略 17200 字符，总长 20000 字节 / 约 5000 token；注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来，全文见 C:\\Users\\swq\\…\\tool-output-174r1-outside.txt…]"
  task_output_pointer_two_roads_174r1_test.go:108: the notice must not claim the allowlist is the only road back, got "[…同上…]"
```

- **还原绿**：`cp <pristine> internal/tools/task.go`，`md5sum` ＝ `432a52be634af33d089f32a269735728`（与备份同值）⇒ `-run 174r1` ＝ **`ok 0.185s`**，4 枚全 PASS。
- **原样绿**：`-run 174r1 -v` 四枚 `--- PASS`（节6 读数）。

### 变异 M2 — 把 `canon`（C26 归一化后的真路径）拼进 notice 尾巴

- 落地证明：`grep -n '" "+canon' internal/tools/task.go` → **`:860`**；`go build ./...` → **rc=0**。
- **变异红**：恰 **1 枚**：

```
--- FAIL: TestPointerTwoRoadNoticeLeaksNoPathOrRoots174r1 (0.01s)
  task_output_pointer_two_roads_174r1_test.go:259: the two-road notice must not relay "C:\\Users\\swq\\…\\003\\tool-output-174r1-sweep.txt": notice "注意：…要么批下那一张卡 C:\\Users\\swq\\…\\003\\tool-output-174r1-sweep.txt"
```

- 其余三枚 **仍 PASS** ⇒ 这条漏**只有**第 4 枚抓得住（不是与前三枚粘连）。
- **还原绿**：`md5sum` 回 `432a52be…` ⇒ 4 枚全 PASS。

### 变异 M3 — 让第二路**变假**（`bridge.go` 的 `case risk.L2:` 后插 `a = AnswerReject`）

- 落地证明：`grep -n "a = AnswerReject" internal/tools/bridge.go` → **`:467`**；`go build ./...` → **rc=0**。
- **变异红**：恰 **2 枚**（同一句承诺的文本半＋行为半各咬一处）：

```
--- FAIL: TestPointerSecondRoadBackIsRealOnThisBridge174r1 (0.01s)
  task_output_pointer_two_roads_174r1_test.go:168: road two must be REAL for the notice to name it: an approved card still got {Text:L2 审批未通过，已拒绝执行 IsError:true RiskLevel:L2 ErrorClass:user_rejected …}
--- FAIL: TestPointerSecondRoadBackIsBookedAsACard174r1 (0.01s)
  task_output_pointer_two_roads_174r1_test.go:198: the approving gate should have let this through, got {Text:L2 审批未通过，已拒绝执行 IsError:true …}
```

- `TestPointerOutsideAllowlistNamesEveryRoadBack174r1` 与 `…LeaksNoPathOrRoots174r1` **仍 PASS** ⇒ 文本半与行为半各守一块地。
- **还原绿**：`bridge.go` md5 回 `0f754371420e361111c2d30c785e72b3`；`git status --porcelain -- internal/ cmd/`＝只我两枚（`M task.go` ＋ `?? 新件`）⇒ bridge.go 逐字节回 HEAD。

## 6. 门禁（我自己现跑，逐把贴末行）

| 门 | 命令 | rc | 末行逐字 |
|---|---|---|---|
| build | `GOFLAGS= go build ./...` | **0** | （无输出） |
| vet | `go vet ./internal/tools/ ./internal/agent/` | **0** | （无输出） |
| test（四数分开、禁 相减） | `go test -count=1 -v ./internal/tools/ ./internal/agent/` | **0** | `ok github.com/CarlosShao/wisp/internal/tools 18.916s` / `ok github.com/CarlosShao/wisp/internal/agent 3.502s` |
| d22scan | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; … bans #1-5 internal/=244, bans #1-5 cmd/=42, … ban #7 internal/tools/=25, ban #8 internal/=608, ban #8 cmd/=143` |
| path-length | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree` |

**四数（禁 相减，各自独立计数）**：`^=== RUN` ＝ **443**；`^--- PASS` ＝ **336**；`^--- FAIL` ＝ **0**；`^--- SKIP` ＝ **0**。（RUN 多于 PASS 是子测试层级：`t.Run` 的 `=== RUN` 进 RUN 数、`    --- PASS:` 带缩进前缀故不落在 `^--- PASS`，两数**不是同一分母**，按派单禁 相减。）

**跑门前置**（每把都做）：`unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`；`cmd/wisp` 面本腿未跑（未碰 `cmd/**`）、`PATH` 加 `third_party/sherpa-onnx` 那一味因此未用。

## 7. 红名册逐名（禁 只报枚数）

`grep -E "^--- FAIL|^    --- FAIL"` 在整份 `-v` 输出（443 RUN / 两包）上 ＝ **0 行**。
`grep -E "^--- SKIP|^    --- SKIP"` ＝ **0 行**。
⇒ **新增红 0 枚、既有红 0 枚、skip 0 枚**（无一项需要归属说明）。日志原件：`D:\tmp\wisp\20261021\174r1\logs\gotest-v-174r1.txt`。

## 8. 提交

见交件回复（逐枚提交号＋逐枚 `git show --name-only` 的 `.go` 枚数与越界枚数）。本件与产码**分枚**提交（新规矩）。

## 9. 如实登记的边界（禁 不当作已完成）

1. **两处既有原句我没有删**：旧句"会被拒…才读得回来"在我改的**那一支**里已被替换，但 `internal/agent/spill.go` 与 `internal/tools/task.go` 的其余四支（未接线／规范化失败／不存在／非普通文件）**一字未动**——本腿只动 allowlist 这一支。
2. **`cmd/wisp/run.go:557` 的接线早已在盘上**（`TaskDeps{Roster: rt.tasks, Paths: rt.paths}`），本腿**未碰 `cmd/**`**。端到端 `wisp run` 真进程本腿**未跑**（AC#2b 那三臂归既有 `14b42887` 的 r8 判据，本腿未重跑它）。
3. **`-race` 未跑**；整包 `cmd/wisp` 未跑（不在写面，派单亦要求 `go test ./internal/tools/ ./internal/agent/`）。
4. **本腿不裁自己**：AC#2 的"文案偏严"那半句（编排者 09-27 22:3x 第 3.2 条要的"两条回来的路"）**本腿已落**，但**是否据此翻 AC#2b／AC#2 的任何条件**归非实现者与编排者。
5. **`gofumpt -l` 有两枚既有漂移**（`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`）——`git diff --stat HEAD` 对这两枚**为空** ⇒ **不是本腿产生**，本腿未洗别人的件（派单禁 `gofumpt -w` 洗别人的件）。我这两枚 `gofumpt -l` 为空。
6. **一处越界自查**：`internal/tools/task.go:648` 有一枚 `注意`（U+26A0）**在 HEAD blob 里早已存在**（`git show HEAD:internal/tools/task.go | grep -c`＝1），非本腿引入；d22scan `ban #8` 对注释豁免故 rc=0。我新增/改动的字节里被禁码位（U+00A7／U+2713／U+2264／带圈数字／emoji）现扫＝**0**。
