# 197-subagent-entity-r1b — 票 197 实体层 r1b 的证据件（正控由 `197-r1c` 现跑）

**本文件是谁写的、写的是什么**：代码由 `197-r1b`（commits `7ea14ce3` ＋ `64c1eea0`）交完，它交完码就被模型连接打断，
**证据件一枚没落盘**。本文件由后继腿 `197-r1c` 写，内容是：派单 `2026-09-28-192x-impl-197r1b-pool-honesty-and-format.md`
第 3 节要求的三样里欠的那一样，**并且两发变异正控是本腿亲手跑的**（不是转述上一腿）。
本腿**没有改动任何产码或测试文件**：盘上动作只有「临时改 → 跑 → 逐字改回 → `git diff --exit-code` 自证」，
唯一入库的新件就是这份 `.md`。

---

## ① 起手锚点（本腿现取）

```
$ date "+%Y-%m-%d %H:%M %z"
2026-09-28 20:23:57 +0800          # 起手；写完这份的复核时刻 20:32

$ git log --oneline -6
a90677e9 ticket 200 (200-r2): consume the rewrite account before naming a tree, and give the panel carrier a producer
64c1eea0 197-r1b 票 197 §2(c) 丙：把那枚注释声称的钉真写出来（cmd/wisp/subagent_stream_key_197_test.go）
7ea14ce3 197-r1b 票 197/211 甲：池的诚实数 8->4（= 桥的 D38d 天花板）＋常驻判据＋两枚桥上路＋gofumpt
d1c440ad 派单第三轮只读调研：上一轮自己点名的三处未覆盖面（openchamber mobile/vscode/extensions＋DSH 52 枚 ui-* 分包＋MiniMax agent-modules 12 枚 147 个 ts），不再交目录级抽样
a928b83d ledger(A412 收 200-r1＋……)＋派 200-r2／197-r1b＋票 200 记界面侧那一跳
84feec44 ticket 200: read the project's own instructions for the AI (all five harnesses do)

$ git status --short | wc -l
102                                 # 全程与别人的脏文件共存；本腿一枚没碰、一枚没提交
```

起手时**我这一路的写面是干净的**（这是后面所有「零残留」读数的参照点）：

```
$ git diff --exit-code -- internal/tools/ cmd/wisp/ ; echo $?
0
```

别人的在飞件（`.gitignore`、`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/r6/logs/**`、
`design/**` 那 16 枚 ` D`、以及 200-r2 刚入库的那批 probes）**一律未碰**。

编排者 20:2x 的现量（口径照引，不重复浪费轮次，下面标注「编排者现量」处即此）：
`go test -count=1 -run 'Test197|SubagentStreamKey' ./internal/tools/ ./cmd/wisp/` 两包**都 ok**；
新判据 `Test197SubagentPoolNeverExceedsBridgeCeiling` 确实在 `internal/tools/subagent_197_test.go:368`。

**起手时本证据件不存在**（这一条决定了下面 §⑤/§⑥ 的两笔账）：

```
$ ls docs/evidence/s1/ | grep -c 197-subagent-entity-r1b
0
$ git ls-files -- docs/evidence/s1/197-subagent-entity-r1b.md | wc -l
0
```

---

## ② (a)(b)(c) 三处的改前读数 → 改后读数

「改前」一律取自 `7ea14ce3` 的父commit `d1c440ad`（＝上一腿动手之前的树），命令是
`git cat-file blob d1c440ad:<path>`，**不依赖任何人的转述**；「改后」取自当前工作树（起手时＝`a90677e9`）。

### (a) 池的诚实数：`MaxConcurrentSubagents` 与桥的 D38d 天花板对齐

| 尺 | 改前（`d1c440ad`） | 改后（现树） |
|---|---|---|
| 池常量 | `internal/tools/subagent_197.go:59` `MaxConcurrentSubagents = 8` | `internal/tools/subagent_197.go:73` **`= 4`** |
| 桥天花板（契约，未动） | `internal/tools/bridge.go:24` `const MaxToolConcurrency = 4` | **同行同值 4**（本腿与上一腿都没开它） |
| 池 − 天花板 | **+4**（名册多承认 4 枚机器同时跑不动的） | **0** |
| 常驻判据 | **不存在**：`git cat-file blob d1c440ad:internal/tools/subagent_197_test.go \| grep -c MaxToolConcurrency` ⇒ **0** | `subagent_197_test.go:368` `Test197SubagentPoolNeverExceedsBridgeCeiling`，判的是 `MaxConcurrentSubagents > MaxToolConcurrency`（两枚**独立读数**相比，不是常量自比） |
| 旧断言的形状 | `:363 Test197SubagentPoolCapsAtEight` ＋ `:694` 「池/深度常量漂了：%d / %d, **want 8 / 1**」 | 用例改名 `:406 Test197SubagentPoolCapsAtBridgeCeiling`（名字与行为一致），手打的 8 换成 `MaxConcurrentSubagents`；深度 1 不动 |
| 给模型读的文本 | 说明由常量生成 | `subagent_197.go:156-162` 仍由常量生成，并有 `:383 Test197SpawnDescriptionNamesTheRealPoolCap` 钉住「文本与常量不分家」 |

AC#1 并排打印（现跑）：

```
$ grep -n "MaxToolConcurrency = \|MaxConcurrentSubagents = " internal/tools/bridge.go internal/tools/subagent_197.go
internal/tools/bridge.go:24:const MaxToolConcurrency = 4
internal/tools/subagent_197.go:73:	MaxConcurrentSubagents = 4
```

AC#4（测量改走真桥）：`spawnDirect` 作为函数**已不在盘上**——
`grep -rn "func spawnDirect" internal/tools/` ⇒ 无匹配（rc=1）；全仓只剩 `subagent_197_test.go:400-401` 两行注释在说它已被去掉、
以及**为什么**去掉（池钉到天花板之后，绕行量不出任何这条真桥路量不到的东西；旧绕行正是「8」能一路绿灯过整套的原因）。
上一腿选的是「去掉」那一支，不是「保留并写清它量的是池语义」那一支。

AC#3（池满硬拒＋可读理由＋不占行不占位）由 `:491 Test197FullPoolRefusesNextSpawnWithReadableReason` 承着，
本腿复跑为绿（§③/§④），它现在也带同一枚前置守卫（`:493` 池>天花板即 Fatalf）。

### (b) 格式门：跟踪集里本写面的源件跑成 gofumpt 干净

尺：`git ls-files -z '*.go' | xargs -0 "$GOPATH/bin/gofumpt.exe" -l`，仪器版本
`gofumpt v0.12.0 (go1.27.1)`（现跑 `-version` 读数）；跟踪的 `.go` 共 **647** 枚。

| | 读数 |
|---|---|
| 改前（**本腿亲手量**：把 `d1c440ad` 的两枚 blob 抽到 `.scratch/wisp/probes/197/r1c/pre/` 再 `gofumpt -l`） | 点出 **2 枚**：`subagent_197.go`、`subagent_197_test.go` —— 即上一腿动手时这两枚确实未净 |
| 编排者 19:2x 的口径（派单 §2(b)） | 跟踪集 7 枚：`cmd/wisp/run.go` ＋ 上面两枚 ＋ `.scratch/wisp/probes/**` 4 枚 |
| 编排者 20:2x 的口径 | 「已不含任何 `internal/`、`cmd/` 下的文件，只剩 `.scratch/wisp/probes/**` 那 4 枚」 |
| **改后（本腿现量）** | `internal/`、`cmd/` 下 **0 枚**（与编排者现量一致）；`.scratch/wisp/probes/**` 下 **5 枚**（逐名见下，**比派单/现量口径多 1 枚**，见 §⑥ 第 1 条） |

`.scratch/wisp/probes/**` 那 5 枚的逐名归属（都是别人刻意做坏的变异样本，**不归本票、一枚没动**）：

| 文件 | 最后触碰它的 commit |
|---|---|
| `.scratch/wisp/probes/163/a1/main.go` | `e13196bd` 2026-09-27 21:27 |
| `.scratch/wisp/probes/174/c2/zz174c2_wiring_pair_windows_test.go` | `6479d540` 09-28 14:55 |
| `.scratch/wisp/probes/183/r2/mut/task-boxset-off.go` | `78ea1d19` 09-28 14:07 |
| `.scratch/wisp/probes/185/c1/mut/fs_broken.go` | `4813567e` 09-28 14:42（这枚 gofumpt 直接吐语法错误：`4:1: imports must appear before other declarations`） |
| `.scratch/wisp/probes/185/r1/mut-m1/hostpath_185.go` | `ee1e2118` 09-28 16:43 |

本写面三枚的定点复量（现跑）：
`gofumpt -l internal/tools/subagent_197.go internal/tools/subagent_197_test.go cmd/wisp/subagent_stream_key_197_test.go` ⇒ **输出为空**（rc=0）。

### (c) 那枚假注释：`SubagentStreamKeyPrefix` 声称的「两包相等」的钉

| | 读数 |
|---|---|
| 改前（`d1c440ad:internal/tools/subagent_197.go:47`） | 注释逐字：`// other, so cmd/wisp/subagent_stream_key_197_test.go pins them equal.` —— 而这枚文件**盘上不存在**：`git ls-tree -r --name-only 7ea14ce3 -- cmd/wisp/ \| grep -i subagent_stream_key` ⇒ 无匹配（rc=1）；`git grep -l SubagentStreamKeyPrefix 7ea14ce3 -- cmd/wisp/` ⇒ rc=1 |
| 改后（现树） | `cmd/wisp/subagent_stream_key_197_test.go` **已入库**（`git ls-files -- <path>` 回显该路径，`git show --numstat 64c1eea0` ＝ **63 增 0 删**），内含两枚判据：`TestSubagentStreamKeyPrefixAgreesAcrossBothPackages`（`:30`，比两包常量相等 ＋ 把 `"subagent:"` 这枚字面量本身也钉住）与 `TestSubagentStreamKeyBuildersAgreeAcrossBothPackages`（`:52`，比两侧 `SubagentStreamKey()` 对生产会用的 id 拼出的键相等） |
| 选的是哪一支 | **丙**（把注释说的事真做出来）。理由：`internal/tools` 与 `internal/panel` 不许互相 import，只有同时 import 两者的组合根比得了；丁（改注释）会让「键形状分家」这一形在盘上仍然无人测。**没有**顺手合并成单一真源——那是票 197 载体层 r3 的活（裁在 `A406`） |

---

## ③ 两发正控（本腿亲手跑）＋「哪一枚 commit 才算未修码」

仪器：所有 `go test` 都带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（不带即 `exit status 0xc0000135`，一枚用例都不跑）。
原始 stdout 逐发落盘在 `.scratch/wisp/probes/197/r1c/`（`base-*`、`mut1-*`、`mut2-*`、`restored-*`、`gate-*`、`gofumpt-tracked-final.txt`、`d22scan-final.txt`、`pre/`），只建不删。

### 正控一（池）：把 `MaxConcurrentSubagents` 临时写成 8

改前的绿基线（本腿现跑，`-run Test197 ./internal/tools/`）：**11 枚全 PASS**，`ok 0.059s`——
`Test197SpawnPublishesIdentityRow`／`Test197RowExistsBeforeFirstChildModelCall`／`Test197SubagentPoolNeverExceedsBridgeCeiling`／
`Test197SpawnDescriptionNamesTheRealPoolCap`／`Test197SubagentPoolCapsAtBridgeCeiling`／`Test197FullPoolRefusesNextSpawnWithReadableReason`／
`Test197ChildCannotDeriveSubagent`／`Test197ConclusionCarriesTaskOutputTaint`／`Test197CancelIsPerRowAndNeverCascades`／
`Test197SubagentHasNoSelfApprovalOutlet`／`Test197StreamKeyShapeIsLiteral`。

**改坏 → 跑**（唯一改动：那一枚常量 4→8，注释与其余一切原样）：

```
$ go test -count=1 -v -run Test197 ./internal/tools/        # 池 = 8
RC=1
--- FAIL: Test197SubagentPoolNeverExceedsBridgeCeiling (0.00s)
    subagent_197_test.go:370: 池 8 大于桥的 D38d 天花板 4：多出来的 4 枚会被桥排成队，名册却说它们在跑；诚实数 = 天花板（票 211 甲），要更多并发得先动契约（票 211 乙），不是改这枚常量
--- FAIL: Test197SubagentPoolCapsAtBridgeCeiling (0.00s)
    subagent_197_test.go:409: 池 8 大于桥的天花板 4：这一发要同时占 8 枚桥位，多出来的会排在 sem 后面起跑不了；先修 Test197SubagentPoolNeverExceedsBridgeCeiling 那条账
--- FAIL: Test197FullPoolRefusesNextSpawnWithReadableReason (0.00s)
    subagent_197_test.go:494: 池 8 大于桥的天花板 4：先修 Test197SubagentPoolNeverExceedsBridgeCeiling
FAIL	github.com/CarlosShao/wisp/internal/tools	0.070s
```

⇒ 红 **3 枚**（其中 §AC#2 要求的那枚 `NeverExceedsBridgeCeiling` **确实真红**，原文如上）；其余 **8 枚仍 PASS**，
包括 `Test197SpawnDescriptionNamesTheRealPoolCap`——它**结构上抓不到这一形**（文案由常量生成，池怎么改它都自洽），
别把它算进防线。

**逐字改回 4 → 复跑 → 自证零残留**：

```
$ go test -count=1 -v -run Test197 ./internal/tools/        # 池 = 4（改回后）
RC=0    --- PASS=11 / --- FAIL=0    ok github.com/CarlosShao/wisp/internal/tools 0.076s
$ git diff --exit-code -- internal/tools/ ; echo $?
0
```

**哪一枚 commit 才算未修码（这条不许含糊）**：

- 未修码 = **`d1c440ad` 及其之前的整条 197 r1 线（`b9fa815b`＋`a4b75096`）**：那棵树里池 = 8、
  测试件里 `MaxToolConcurrency` 出现 **0 次**（上面 §②(a) 现量），另有 `:363 Test197SubagentPoolCapsAtEight` 与
  `:694 want 8 / 1` 两枚**把 8 钉成期望值**的断言。⇒ 关键结论：**那棵树的套件是绿的**，「池 8」在当时的尺上根本不是失败。
- 因此「会红」**不是旧代码的性质，是 `7ea14ce3` 新钉的性质**：我这一发红的状态 = `7ea14ce3`（含至 HEAD）的测试件 ＋ `subagent_197.go` 那一枚常量回成 8。
  换句话说，**只有「保留判据、把常量抬回去」这一形可见**；
  而「整枚 `7ea14ce3` 连测试件一起回退」（＝真回到未修码）**今天没有任何常驻判据会红**——本腿没有实测这一形（不能 checkout，见 §⑤ 第 3 条），
  它是从上面两枚静态读数（`grep -c MaxToolConcurrency` ⇒ 0、旧断言写死 want 8）推出来的，**推论就该记成推论**。
- 修吗？判据落在 `7ea14ce3`；语义形状（超一枚＝硬拒、不排队、不留名册行）没动，深度 1 没动。

### 正控二（流键的两包钉）：把 `SubagentStreamKeyPrefix` 临时改成别的字面量

改前基线：`-run SubagentStreamKey ./cmd/wisp/` **2 枚全 PASS**，`ok 0.051s`（本腿现跑）。

**改坏 → 跑**（唯一改动：`internal/tools/subagent_197.go` 那枚字面量 `"subagent:"` → `"subagent-"`）：

```
$ go test -count=1 -v -run SubagentStreamKey ./cmd/wisp/       RC=1
--- FAIL: TestSubagentStreamKeyPrefixAgreesAcrossBothPackages (0.00s)
    subagent_stream_key_197_test.go:32: 两包的前缀分家了：tools = "subagent-", panel = "subagent:"（写侧键与读侧键从此对不上，名册里的子代理会有流却没人读得到）
    subagent_stream_key_197_test.go:39: 前缀 = "subagent-", want "subagent:"（§0 冻的键形状：全小写、冒号分隔）
--- FAIL: TestSubagentStreamKeyBuildersAgreeAcrossBothPackages (0.00s)
    subagent_stream_key_197_test.go:60: 同一个 id "task-197-a" 两侧拼出不同的键：tools = "subagent-task-197-a", panel = "subagent:task-197-a"
    subagent_stream_key_197_test.go:60: 同一个 id "00000000-0000-4000-8000-000000000000" 两侧拼出不同的键：...
    subagent_stream_key_197_test.go:60: 同一个 id "带中文的id" 两侧拼出不同的键：...
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.082s
```

顺带把两包一起跑（同一枚变异，读「还有谁抓得到」）：`internal/tools` 那枚自家
`Test197StreamKeyShapeIsLiteral` 也红（`:845` 键形状三行、`:849` 前缀 = "subagent-", want subagent:），
`--- FAIL github.com/CarlosShao/wisp/internal/tools 0.061s`。

**逐字改回 → 复跑 → 自证零残留**：

```
$ go test -count=1 -run 'Test197StreamKeyShapeIsLiteral|SubagentStreamKey' ./internal/tools/ ./cmd/wisp/
ok  github.com/CarlosShao/wisp/internal/tools  0.057s
ok  github.com/CarlosShao/wisp/cmd/wisp        0.058s
$ git diff --exit-code -- internal/tools/ cmd/wisp/ ; echo $?
0
```

（这也是本腿最后一次全量还原自证：**两发变异加起来改过 3 次、还原 3 次，`git diff --exit-code` 一律 0**。）

**哪一枚 commit 才算未修码**：

- 未修码 = **`7ea14ce3` 以及更早**：注释（`d1c440ad:internal/tools/subagent_197.go:47` 起就在说）声称那枚文件钉住了两包，
  而 `cmd/wisp/` 下**没有**任何文件引用 `SubagentStreamKeyPrefix`（两条 git 读数见 §②(c)）。
  那一棵树上「tools 侧前缀被改」这一形**只有 tools 自家那枚 `Test197StreamKeyShapeIsLiteral` 会红**——
  它在写侧包内，**改常量和改它是一笔 diff**，锁不住「两侧一起漂」；「panel 侧单独漂移」那一形当时**全仓无人能看见**。
- 修吗？钉落在 `64c1eea0`（＋63 行，纯新文件，零删除）。
  本发正控实测到的红来自 **`cmd/wisp` 的两枚新钉 ＋ `internal/tools` 的旧钉**，三者各异的射程：
  新钉的 `:32` 抓「两包分家」、`:39` 抓「两包一起漂离 §0 冻的字面量」、`:60` 抓「键拼出来就不一样了」。

---

## ④ 门禁终态逐包读数 ＋ 在册红名册逐名比对

终态＝**两发变异全部还原之后**、写面为 `a90677e9` 的现树。

| 门禁 | 读数 |
|---|---|
| `go build ./...` | **rc=0**，输出空 |
| `go vet ./internal/tools/ ./cmd/wisp/ ./internal/panel/` | **rc=0**，输出空 |
| `go test -count=1 ./internal/tools/`（整包，`-v`） | **`ok 15.191s`**；`--- PASS` **158** 枚、`--- FAIL` **0** 枚 |
| `go test -count=1 ./cmd/wisp/`（整包） | **`ok 87.571s`**；FAIL **0**（⚠ 长跑约 1.5 分钟，与 r1 证据件那发 101.997s 同量级） |
| `go test -count=1 ./internal/agent/`（整包，池的契约拷贝在这包） | **`ok 2.060s`** |
| `go test -count=1 ./internal/panel/`（整包，**只读、一枚没改**） | **`FAIL 1.191s`**，红 **4** 枚 ⇒ 见下表 |
| `sh scripts/d22scan.sh` | **rc=0 clean**：`no D22 ban violations`；ban#1-5 `internal/`=214、`cmd/`=24 枚产码，ban#7 `internal/tools/`=22，ban#8 `internal/`=450、`cmd/`=49（注释与 `_test.go` 在射程内）、`design/`=39、`frontend/`=85；`tools/d22scan` 自身 `runtests.sh: OK packages=[./...] PASS=34 FAIL=0 SKIP=0` |
| gofumpt（跟踪集 647 枚 `.go`） | `internal/`＋`cmd/` 下 **0 枚**；只剩 `.scratch/wisp/probes/**` **5 枚**（逐名见 §②(b)） |

**在册红名册逐名比对**（起手 4 枚 vs 终态 4 枚，`internal/panel`）：

| 起手在册红（编排者给名） | 终态现量 | 结论 |
|---|---|---|
| `TestApprovalCardViewJSONKeysMatchFrontendTypes` | `--- FAIL` 在 | **名册未变** |
| `TestComposerContractTypesMatchFrontend` | `--- FAIL` 在 | **名册未变** |
| `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | `--- FAIL` 在 | **名册未变** |
| `TestC21DesignTokensFourwayAgree` | `--- FAIL` 在 | **名册未变** |

`grep '^--- FAIL' gate-panel.txt | sort` 的点名的就是这 4 枚、**没有第 5 枚、也没有任何一枚消失**。
四处**一未修、二未 Skip、三未放宽**（本腿连写面都没开）。红因按 200-r2 的读数仍是
`Snapshot emits [instructions] that interface PanelSnapshot does not declare`（两枚）、
`frontend/src/components/harness/right-rail.tsx:89` 的第二处样式源（一枚）、
`read design/assets/tokens.css: open ...: The system cannot find the path specified`（一枚，
对应那 16 枚 ` D design/**` 还在别人的暂存区里）——**都归票 200／前端线与 F9，不属本票**。

---

## ⑤ 我没测什么（具名，不留空）

1. **没测「池 8 ⇒ 7/8 红」这个生产症状本身。** 正控一红的是**常量比较**与两枚前置守卫（`:409`/`:494` 直接 `Fatalf`），
   没有任何一枚常驻判据在池被抬回 8 时真的去桥里排队、真的撞 C22 的 30s。
   「8 枚并发被桥排成队 ⇒ 7/8 红」是 **`197-r1` 的实测**，本腿只是复算它与常量的因果形状（票 211 现量表第 5 行）。
   把这条当本腿的读数用就是越账。
2. **没测 panel 侧单独漂移。** 写面禁区（不碰 `internal/panel/**`），所以正控二只从 tools 侧改字面量。
   `:32` 那枚判据是对称的，逻辑上一样红，但**这是推理不是读数**。
   同理，`cmd/wisp` 测试件注释里写的「panel 的 `SubagentStreamKey` 对空/纯空白 id 回 `""`、tools 侧无此分支」这一形，
   本腿**没有实测**（未跑空 id），它由 `64c1eea0` 注释自陈，未入库任何判据。
3. **没测「整枚 `7ea14ce3` 连测试件一起回退」会不会红。** 断言依据是两枚静态读数（旧测试件 `MaxToolConcurrency` 出现 0 次、旧断言写死 want 8），
   **不是一次实跑**——本腿不许 `checkout`/`restore`，也不许删件，所以那一形的「绿」是推理。要实测得用 `git cat-file` 把旧测试件抽到树外再叠，本腿未做（成本一轮长跑，且非派单所求）。
4. **没跑全仓 `go test ./...`**（34 枚 top-level 包，`tools/d22scan` 的 runtests 只覆盖它自己那一片），
   只跑了 `internal/tools`／`cmd/wisp`／`internal/agent`／`internal/panel` 四包 ＋ build/vet/d22scan/gofumpt 四门。
   `internal/risk` 那两枚起手在册的别腿红（`TestResolvePerCallBudget`、`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`）**本腿未复跑**，
   它们的收口归 200-r2（`a90677e9`）。
5. **没跑 `-race`、没跑 `wisp slo`**、没量任何 SLO/延迟预算；`thresholds.go`、golden、`allowlist.txt`、`bridge.go`、`budgets.go` 一字未动
   （AC#5 的证：`git diff --stat d1c440ad HEAD -- internal/tools/bridge.go internal/agent/budgets.go` ⇒ **空**；工作树 `git diff --exit-code HEAD --` 同两枚 ⇒ **0**）。
6. **没测 `task.spawn` 的真端到端**（面板上点进去看子代理流）：本票三层里的载体层 r3 不在射程，
   键的两包钉只保证「拼出来的字符串一致」，不保证「泵真的订阅了那枚键」。
7. **没自证这份证据件的读数能被别人复算**——除了 §③ 的两发（还原后 `git diff` 为 0 这一枚是可核事实）。原始 stdout 留在
   `.scratch/wisp/probes/197/r1c/`（未跟踪、只建不删），编排者要核哪一发直接点名。

---

## ⑥ 对派单的不服（含对本腿任务书的不服）

派单与任务书里**核过为真**的：`7ea14ce3` 的 numstat（`subagent_197.go` +29/−9、`subagent_197_test.go` +213/−50）、
`64c1eea0` 的 +63/0；判据在 `subagent_197_test.go:368` 且判的是 `>` 关系；
「跟踪集 gofumpt 未净里 `internal/`、`cmd/` 为 0」；起手在册红 4 枚的名字与终态逐名一致；
`d1c440ad` 那一版池 = 8、桥 = 4、`MaxToolConcurrency` 在旧测试件出现 0 次。**没有异议。**

要不服的、要更正的：

1. **枚数错了 1**：任务书 20:2x 说 gofumpt 未净「只剩 `.scratch/wisp/probes/**` 那 **4** 枚」，
   派单 §2(b) 也写 4 枚。本腿现量是 **5** 枚（逐名＋归属见 §②(b)）。
   差别不影响判据（`internal/`＋`cmd/` 为 0 才是交件判据），但**账上的数应改成 5**。
   一个可能的成因（只摆形状，不裁定）：`185/c1/mut/fs_broken.go` 让 gofumpt 吐的是**错误行**
   （`4:1: imports must appear before other declarations`）而不是普通的文件名，xargs 传上来的 rc 是 **123**，
   看着像「仪器坏了」而被整行漏数。
2. **派单 §2(a) 说「人为把池写成 8 ⇒ 那枚判据必须真的红」——成立，但它没告诉你红几枚。** 现量：红 3 枚，
   其中**只有 1 枚是判据**（`:368`），另 2 枚（`:409`、`:494`）是**前置守卫 `Fatalf`**，它们说的是「先去修那枚判据」，
   不是「桥级症状复现了」。下一程若把「红 3 枚」读成「三枚测量都反对」就是误账，本文件 §③ 已把它拆开。
3. **任务书给的旧行号已漂，但漂得无害**：派单 §2(a) 现量写 `subagent_197.go:59`/`:205`，那是 `d1c440ad` 的行号；
   现树是 `:73`／`:223`。本腿按派单第 0 节「别引用行号当现量」处理，全部重取。
4. **一笔真不服，形状与 §2(c) 一模一样**：`7ea14ce3` 与 `64c1eea0` 各自交了一枚**指向不存在的文档的注释**——
   `internal/tools/subagent_197_test.go:361` 与 `cmd/wisp/subagent_stream_key_197_test.go:16-18` 都写着
   「正控跑过了，读数在 `docs/evidence/s1/197-subagent-entity-r1b.md`」。
   本腿起手的现量是：那枚文件**既不在盘上也不在索引里**（`ls | grep -c` ⇒ 0、`git ls-files` ⇒ 0，见 §①）。
   ⇒ 「**注释声称的东西盘上不存在**」这一形，在修它的那一腿自己身上**又复发了两次**（这次复发的是证据件，不是钉）。
   派单第 3 节把证据件列为「缺一不可」，所以这不是新洞，是**上一腿没交完的件**；本文件把它做真——
   但要说清：**两发正控的实际执行时刻是 20:25–20:32 +08、执行者是 `197-r1c`，不是 `197-r1b` 的「at delivery」**。
   建议编排者在账上（`A411` 那一带）把这句口径改准；**本腿无权改产码/测试件注释，所以那一行注释本腿没动**（一字未动，可 `git diff` 复）。
   另请裁：要不要给 d22scan 加一形「注释/文档引用仓内路径 ⇒ 该路径必须存在」的常驻检查（**本腿未判它属不属于 ban 的形状，也没做可行性普查**）。
5. **仪器提醒（派单第 1 节写了，任务书也写了，仍然记一笔以免下一腿摔）**：`cmd/wisp` 整包是 **87 秒**级，
   r1 证据件记的是 102 秒。别按「快跑」估这一票的核实成本；不带那两枚 PATH 目录时拿到的是 `0xc0000135`，
   本腿两轮共 8 发全带 PATH，**没有复演那一发**。

---

## 附：本腿的入库件

`docs/evidence/s1/197-subagent-entity-r1b.md`（**就是这份，枚名按派单第 3 节，未改名**），一枚路径、一次 commit、带显式 pathspec、只 commit 不 push。
`git diff --cached --numstat` 的删除列 ＝ **0**（本腿没删过任何东西，包括自己写的临时件）。
