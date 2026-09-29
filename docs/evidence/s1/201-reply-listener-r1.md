# 211-r1 — 答复的听众：`wisp run` 接上原生侧与面板路线两条答复通路

写码子代理 `211-r1`。工单＝`.scratch/wisp/issues/201-nobody-can-answer-an-approval-l1-runs-by-itself-l2-dies-on-timeout.md`
（票面编号 211 那一族；文件名用 201，以内容为准）。
本腿只做派单里的三件事：**接上"答"这一侧的听众**、**证明超时不落执行**、**答复进审计且拒绝理由能回给模型**。
`frontend/**`／`design/**` 零写、零读；未新增任何 C17 方法名；`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／
`allowlist.txt` 一字未动；`internal/panel/**` 一字未动。

---

## ① 起手锚点 ＋ 今天的绿用例名册（逐包点数）

```
date                      2026-09-29 10:22:57 +0800
git rev-parse HEAD        24eef597a11e8aa4ff6c3637aef6ce7cea7c4ac3   (dev)
git status --short        起手脏项（本腿一律未碰）：M .gitignore、
                          M .scratch/wisp/issues/219-…md、M .scratch/wisp/probes/152/my152.py、
                          M .scratch/wisp/probes/161/r6/logs/flip-*.txt、
                          16 枚 D design/assets/** 与 design/screens/** 删除、
                          M design/doubao/**（4 枚）、M docs/evidence/s1/152-…-accept-r1.md、
                          ?? docs/evidence/s1/219-approval-reply-surface-c1.md（只读腿的在文件，本腿未碰未提交）、
                          ?? docs/reports/missing-features-2026-09-28-v3.md、?? .scratch/wisp/probes/**（多家）、
                          ?? part1-*.txt … part3-*.txt（8 枚散件）、?? .zcodeignore、?? design/old/ 等
```

起手复算（**不是转述派单**，全部本腿现跑）：

```
grep -rn "DecideFromPanel|DecideFromNative|\.Veto(" --include=*.go internal/ cmd/ tools/ | grep -v _test.go
  → 只剩 gate.go:610 / gate.go:622 两枚定义行（含各自注释行），生产调用者 = 0
  ⇒ 派单/普查那句"答这一侧是零调用者的死代码"在本腿起手仍然成立，本腿就是来收这一格的。
grep -rn "approval\.New(" --include=*.go cmd/ internal/ tools/ | grep -v _test.go
  → 生产构造点只有一枚：cmd/wisp/run.go:406（另 gate.go:108 是包内 Options 默认路径的注释行）
  ⇒ 今天唯一装配了审批 gate 的生产进程入口是 `wisp run`；球（internal/ball）只被 cmd/balldebug 构造，
    cmd/wisp/resident_windows.go 对 Gate/ball 零引用（现读 grep 零命中）。
```

### 今天的绿用例名册（整包 `-v` 点数，非 `-run` 单跑）

命令（本腿全程带 PATH，否则 cmd/wisp 与 internal/audio 等在装载期就死 0xc0000135）：

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" WISP_ENV=test go test -v -count=1 ./...
    → .scratch/wisp/probes/211/r1/logs/baseline-v-full.txt（789996 字节）
    → 逐包点数：.scratch/wisp/probes/211/r1/logs/baseline-per-package.txt
    → 全仓合计：=== RUN 1693 ／ --- PASS 1121 ／ --- FAIL 6 ／ --- SKIP 7 ／ 包 ok 24 ／ 包 FAIL 3
      （RUN 含子测试；PASS/FAIL/SKIP 只数顶层行 —— 这就是"点数口径"，别把它当用例总数）
```

本腿要动的包，起手读数逐名：

| 包 | run | pass | fail | skip | 终态 |
|---|---|---|---|---|---|
| `cmd/wisp` | 157 | 97 | 0 | 0 | ok（单包 95.352s） |
| `internal/agent/approval` | 54 | 34 | 0 | 1 | ok |
| `internal/tools` | 207 | 158 | 0 | 0 | ok |
| `internal/panel` | 152 | 95 | **4** | 0 | FAIL |
| `internal/ball` | 55 | 48 | **1** | 0 | FAIL |
| `internal/risk` | 198 | 128 | **1** | 1 | FAIL |
| `internal/perm` | 14 | 14 | 0 | 0 | ok |
| `internal/agent` | 94 | 73 | 0 | 0 | ok |

起手在册红＝**6 枚，不是 4 枚**（派单只给了 `internal/panel` 那四枚）。逐名与出处：

1. `internal/panel` — `TestApprovalCardViewJSONKeysMatchFrontendTypes`
2. `internal/panel` — `TestComposerContractTypesMatchFrontend`
3. `internal/panel` — `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`
4. `internal/panel` — `TestC21DesignTokensFourWayAgree`
   （以上四枚＝派单名册，逐名对上；派单写的是 `TestC21DesignTokensFourwayAgree`，
   盘上真名是 `FourWay`（大写 W），这是名册里的一处字面漂移，不影响"同一枚"判定）
5. `internal/ball` — `TestC21TableColourRowsMatchTokensCSS`
   **起手新增红，不在派单名册里**。红因是别人的脏工作树，不是代码：
   `read design/assets/tokens.css: … The system cannot find the path specified`
   ＝起手 `git status` 里那 16 枚 `D design/assets/**`／`D design/screens/**` 删除（另一条腿在做的），
   该用例自己的句子写着"the CSS leg of this check must never skip"。**本腿不碰 `design/**`，不修。**
6. `internal/risk` — `TestResolvePerCallBudget`
   **起手新增红，不在派单名册里**，且是**性能预算**判据：
   `C26 Resolve: 5665468 ns/op = 5.665 ms/op (budget 1.000 ms, 490 samples)`
   ＝共享机器上并发跑测试把这条尺压爆了（本腿起手时另有腿在跑整包）。
   阈值与 `thresholds.go` 一字节都不许动 ⇒ **本腿不修、不 Skip、不放宽，只登记。**

⇒ 这两枚新增红对 §④ 的比对口径有直接影响：本腿交件的判据是
**"终态名册 ⊇ 起手名册且没有新名字"**，不是"必须为空"（派单给的名册本身起手就不全）。
