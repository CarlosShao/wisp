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
