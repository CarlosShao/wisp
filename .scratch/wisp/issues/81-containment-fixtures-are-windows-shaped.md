# 81 — 两份 containment 测试的**阳性对照是 Windows 形状**，ubuntu 上必红（票 76/79 的判据仪器在 Linux 上不自证）

**Status:** code 侧完成、四框已勾（两侧本地实测）· **未闭**：ubuntu CI 侧 `test-core` 要一次 push 才能复验（本代理不 push）、裁决表 `docs/evidence/s1/81-*.md` 归对抗验收代理
**Type:** 测试夹具的平台可移植性缺陷（不是生产洞，但**它让"防逃逸的那条判据"在 Linux 上等于没跑**）
**Blocks:** 票 70 的 AC#2/AC#6（`test-core` 现在只剩 10 条红，其中 2 条就是本票）· **Blocked by:** nothing
**Packages:** `internal/agent/spill_path_invariant_test.go`、`internal/memory/artifacts_path_invariant_test.go`
（必要时它们的兄弟测试文件）。**禁改**：`internal/agent/spill.go`、`internal/memory/artifacts.go` 的**生产行为**、
`internal/risk/**`（票 70/82 的地界）、`tools/d22scan/**`、`allowlist.txt`、任何 `docs/PLAN.md`/`docs/specs/*.md`。

## 实测证据（CI run **35558750456**，commit `17efc2c`，job `test-core` 步骤 `Portable package tests`，ubuntu）

**① `TestArtifactsContainmentByDirectoryListing`（`artifacts_path_invariant_test.go:464`）**
```
purge removed [data/artifacts/nested\canary-separator.txt data/artifacts/real-1.txt data/artifacts/real-2.txt],
want exactly [data/artifacts/nested data/artifacts/nested/canary-separator.txt
              data/artifacts/real-1.txt data/artifacts/real-2.txt]
```
根因：夹具用**字面反斜杠**造"子目录里一个文件"。Windows 上 `\` 是分隔符 ⇒ 真的是 `nested/` 目录；
Linux 上 `\` 是合法文件名字符 ⇒ 落成一个**顶层文件**，目录集合自然对不上。

**② `TestSpillContainmentByDirectoryListing`（`spill_path_invariant_test.go:280`、`:286`）**
```
control: escape id "..\\..\\..\\..\\CONTROL-escape" resolves to ".../data/artifacts/tool-output-..\\..\\..\\..\\CONTROL-escape.txt"
the listing never saw the deliberate escape outside the data dir (…); AC#2 would be vacuous
```
根因：**阳性对照**要证明"如果真有文件写到 data dir 之外，列举一定看得见"。Linux 上那串 `..\` 不构成逃逸，
文件老老实实待在 `artifacts/` 里 ⇒ 对照本身失效。**注意这不代表生产有洞**：票 79 的编码器已经把所有
`/ \ : .` 百分号转义，Linux 那个长文件名恰恰是**正确**行为。

## 判据（1:1，裁决表 `docs/evidence/s1/81-*.md`）

- [x] **AC#1** 两条用例都改成**平台无关地表达同一个意思**：逃逸/嵌套的"物理路径"用 `filepath.Join` 或
  `filepath.Separator` 构造，同时**保留**字面反斜杠拼写作为额外一例（它在 Linux 上是合法文件名、
  在 Windows 上是分隔符，这个不对称本身就是要钉的东西）。**不许**给这两份文件加 `//go:build windows`
  把它们变成"Linux 上静默不跑"——那是 A49② 里票 78 代理拒绝过的那类修法（整包被排除 = 静默跳过）。
- [x] **AC#2** 两个方向各自**制造一次红**：(i) 在 Linux（docker `golang:1.27`）证明改后这两条**真跑且绿**，
  且把生产侧的转义临时拆掉时它们**变红**（说明判据仍咬得住）；(ii) 在 Windows 本机同样跑一遍。
  ⚠ 变异锚点选**真正承载行为的那一行**，同一条 `&&` 链里先 grep 证明落地，还原后 `git diff --quiet` 证干净。
- [x] **AC#3** **全仓扫同族**：`grep -rn` 出所有在没有 build tag 的测试文件里用字面 `\` 当目录分隔符的地方
  （含 `"a\\b"`、`filepath.ToSlash` 反用、`\r?\n` 之类合法的除外），逐条列进本票 log：
  要么本票一并修掉，要么写清"它在两侧语义相同、不需要修"的理由。**不许只报"扫了没问题"**——要给出命中清单。
- [x] **AC#4** 门禁（只跑自己碰的包）：`gofmt -l` 空、`go vet` rc=0、`go test -count=2` rc=0，
  并 `GOOS=linux go vet` rc=0；逐跑点名 `--- SKIP`/`--- FAIL` 行数与名字。

## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；每次 commit 同步票面 Status + 勾框 + 末行 `next=`；
`git commit -q -F - -- <显式路径>` + **带引号** heredoc `<<'MSGEOF'`；禁 `git add -A`/`.`；commit 前核对
`git diff --cached --name-only`；禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`（共树 A34）；**不 push**；
不在仓内建 worktree（A38④），纯净树用 `git archive HEAD | tar -x -C /tmp/...`；
docker 挂载**快照**而不是工作树（别人的在飞改动会混进你的结论）；
票面 Progress log append-only：**要改的那一行先读再替换**。

## Progress log（append-only）

- **C1（AC#1 的代码侧）** 两份夹具改成平台无关表达，字面 `\` 一例保留并**显式钉住不对称**：
  - `internal/memory/artifacts_path_invariant_test.go`：`inv76Shapes` 的 `separator`/`dotdot` 改用
    `inv76Sep()`（`filepath.Separator`）拼，target 一律 `filepath.Join(被指的目录, 调用者给的名字)`
    = "这台机器把它解析成哪儿"；新增 `separator_backslash_literal`、`dotdot_backslash_literal` 两例保留字面 `\`。
    (e) 段原本硬编码 4 条 `data/artifacts/...` 键（Linux 上必然对不上），改为由夹具真实种下的路径经
    `inv76ReclaimKeys` 派生（自己 + 被 OS 认可的父目录，artifacts 根本身不算）；purge 计数也不再硬编码 2，
    改由 `before` 快照里 artifacts 树下的**文件**条数 -1 推得（`PurgeArtifacts` 的 n 只数文件，实测 Windows 6 键/3 文件）。
    新增 `TestArtifactsLiteralBackslashKeepsItsAsymmetry`：同一串字节在 Windows 是两层、在 Linux 是一个扁平文件名，
    两侧各自断言，并注明"加 build tag 就等于删掉这条判据"。
  - `internal/agent/spill_path_invariant_test.go`：`inv76aEscapeID` 改用 `inv76aEscapeIDWithSep(..., string(filepath.Separator))`
    （Linux 上 `..\` 不构成逃逸 ⇒ 阳性对照失效，这就是票面 ② 的根因）；字面 `\` 拼写保留为**第二条对照** (d)，
    Windows 断言它落在 data dir 之外、Linux 断言它落回 artifacts 里成为一个扁平名；端到端用例的 `nested\canary.txt` /
    `..\canary.txt` 扩成"平台分隔符 + 字面反斜杠"双拼写。
  - Windows 本机：`gofmt -l` 空、`go vet ./internal/agent ./internal/memory` 无输出、
    `-run 'Containment|LiteralBackslash|HostileShapes|StaysUnderDataDir|RejectsTheFour'` 两包 `ok`（0 SKIP / 0 FAIL）。
  - AC#1 勾框**留到 Linux 侧绿了再打**——票面要的是两侧都跑，不在只测了一侧时就记完成。
  next= docker `golang:1.27` 跑 `git archive` 快照，证明这两条在 Linux 上真跑且绿（AC#2(i)），然后做生产侧转义/守卫的变异检验
- **C2（Linux 首轮实测：改对了方向，但暴露我第一版仍带 Windows 预设）** 快照 `git archive HEAD` → `D:\tmp\wisp81snap`，
  `docker run golang:1.27 -count=2 -v` 两包：50 条 `=== RUN`、**4 条 FAIL**（同一对，×2）、0 SKIP；agent 包 `ok`，memory 包红。
  失败原文（`artifacts_path_invariant_test.go:502`）：
  `purge removed [data/artifacts/..\canary-dotdot-literal.txt data/artifacts/nested data/artifacts/nested-backslash\canary-literal.txt
  data/artifacts/nested/canary-separator.txt data/artifacts/real-1.txt data/artifacts/real-2.txt] want exactly [去掉第一条]`。
  根因是**我自己的 C1 版本**：`wantRemoved` 仍按名字点了两例字面反斜杠 canary，而 Linux 上 `..\canary-dotdot-literal.txt`
  这个字面名**不越界、老老实实落在 artifacts 里**，于是它也归 purge 管——这正是票面要的不对称的第二面，
  第一版把它漏掉了（同一原因让 `LiteralBackslash` 那条断言误报"意外反斜杠条目"）。
  修法：期望集合改为**按前缀从 target 派生**（`inv76Rel(root, sh.target)` 落在 `data/artifacts/` 里就算），
  不对称测试 `TestArtifactsLiteralBackslashKeepsItsAsymmetry` 重写成两侧通用的"含 `\` 而本平台不当它是分隔符 ⇒ 必是扁平一条"，
  并加 `len(planted) < 4` / `slices.Equal(flats, want)` 两条反-vacuous 卫兵。
  Windows 本机重跑：`rc=0`、50 条 `=== RUN`、0 FAIL、0 SKIP。
  next= 重新 archive 快照跑 Linux（AC#2(i) 要绿），再做 AC#2 的变异检验
- **C3（AC#2(i) Linux 侧绿）** 快照 `git archive HEAD`（`4683c34`）→ `D:\tmp\wisp81snap2`，
  `docker run golang:1.27`，`go test -count=2 -v -run 'Containment|LiteralBackslash|HostileShapes|StaysUnderDataDir|RejectsTheFour' ./internal/agent ./internal/memory`
  → **rc=0，RUN=50 PASS=50 FAIL=0 SKIP=0**（两包各 `ok`；`-count=2` 与 `-count=1` 的 RUN 比 50:25 成立 ⇒ 真的跑了两遍，不是空匹配）。
  Windows 侧同一条 `-run` 过滤器：工作树（内容 = `4683c34` 的这两份文件）`-count=2` → `rc=0，RUN=50 PASS=50 FAIL=0 SKIP=0`；
  同一快照 `-count=1` 基线 → `RUN=25 PASS=25 FAIL=0 SKIP=0`，快照还原后再跑一遍仍 25/25。
  平台日志（两侧都留了证据行）：
  - Linux：`GOOS=linux separator='/': literal-backslash names living as ONE flat file name in the artifacts dir:
    [..\canary-dotdot-literal.txt nested-backslash\canary-literal.txt]`；
    `control: escape id "../../../../CONTROL-escape" resolves to "/tmp/.../001/CONTROL-escape.txt"`（对照真的越界）；
    `literal-backslash control: id "..\..\..\..\CONTROL-literal" resolves to ".../data/artifacts/tool-output-..\..\..\..\CONTROL-literal.txt"`（Linux 上它**不**越界，断言的就是这一条）。
  - Windows：`GOOS=windows separator='\': ... : []`（零扁平）；两条对照都落在 root。
  ⇒ **AC#1/AC#2(i) 打勾**：两侧同一批 25 条用例、同一套断言，无 build tag、无 SKIP。
- **C4（AC#2(ii) 变异检验：两侧各制造一次真红）** 全部在快照 `D:\tmp\wisp81snap2` 里做，**工作树一次都没脏过**
  （`git diff --quiet -- internal/agent internal/memory` 当场证明）。锚点先 `grep -n` 证明改动落地、同一条 `&&` 链里再跑：
  - **MUT-A `internal/agent/spill.go:210`**（真正承载转义的那一行 `if c == '_' || ... {`，加 `|| c == '/'` ⇒ 分隔符不再转义）：
    Linux `rc=1，RUN=25 FAIL=5`：`TestSpillContainmentByDirectoryListing`、`TestSpillIntoRealStoreThenDeleteStaysUnderDataDir`、
    `TestSpillCallIDHostileShapesSanitizedToBareNames{,/separator,/dotdot}`；Windows 同一变异 `rc=1，RUN=25 FAIL=5`，同样三条顶层用例红。
    （第一次尝试因把 `\` 也塞进 sed 替换而写成非法字符字面量 `'\'`，**编译失败不算行为变异**，已作废重跑；`/` 在两侧都是 `filepath.Join` 的分隔符，足以证明判据咬得住。）
  - **MUT-B `internal/memory/artifacts.go:297`**（`validArtifactName` 头部插 `if name != "" { return nil }` ⇒ 名字守卫整体失效）：
    Linux `rc=1，RUN=25 FAIL=16`（含 `TestArtifactsContainmentByDirectoryListing` + 两条 `RejectsTheFourHostileShapes` 的 6 个子测试各两侧），
    实测红字：`DeleteArtifact("nested/canary-separator.txt") accepted a caller-named path`、
    `hostile names changed the tree: + [] / - [data/artifacts/..\canary-dotdot-literal.txt data/artifacts/nested-backslash\canary-literal.txt
    data/artifacts/nested/canary-separator.txt data/canary-dotdot.txt]`（Linux 版）——Windows 版同一条断言报的是
    `[data/artifacts/nested-backslash/canary-literal.txt ... data/canary-dotdot-literal.txt data/canary-dotdot.txt]`，
    即**同一判据在两侧各自咬住各自的物理事实**，这正是 AC#1 要的形状。
  - 两个变异都 `git archive HEAD | tar -x` 还原（`grep -c` 命中归零），还原后 Windows 重跑 25 RUN / 25 PASS / 0 FAIL / 0 SKIP。
    ⇒ **AC#2 打勾**。
- **C5（AC#3 全仓同族扫描：命中清单）** 仪器三条，都在 `dev`@`4683c34` 上跑，全量输出不截断：
  1. **广谱**（任一测试文件里含转义 `\\` 或含 `\` 的反引号串）：**389 行 / 70 个文件**；
     其中被 windows 门（文件名 `_windows_test.go` 或 `//go:build windows`）挡住的 **13 文件 / 46 行**，
     未被挡的 **57 文件 / 343 行**。这 343 行的大头是**正则转义**（`\b`、`\.`、`[0-9]` 类，如 `internal/ball/tokens_test.go:95`、
     `internal/observe/nobarego_test.go:26`）与 `"\n"`/`"\r\n"` 类合法转义，以及**只当数据用的 Windows 路径字面量**
     （`internal/config/{unwired,manager,migrate,boundary}_test.go` 的 TOML 串、`internal/memory/dao_test.go:385` 的授权 Pattern、
     `internal/observe/logging_test.go:107` 的被脱敏文本、`internal/ball/position_test.go` 的 `\\.\DISPLAY1` 设备名、
     `internal/risk/{provenance,rules,taintmatch,syncdirs}_test.go` 的评估输入字符串）——
     两侧都是同一串字节、不构造物理路径，**不需要修**（`internal/risk/**` 与 `internal/config/**` 亦为票 70/80/82 地界）。
  2. **窄谱（决定性）**：`\` 与路径构造调用同行的 → **11 行**，去掉 5 条 `"\n"` 假命中（`internal/tools/{wiring,fs_write}_test.go`、
     `internal/observe/logging_test.go:191,199`）与 1 条写入文件**内容**的 `C:\Users\x\Dropbox\.dropbox`（`internal/risk/syncdirs_test.go:48`），
     **真正用字面 `\` 构造物理路径的只有 5 行**：
     - `internal/memory/artifacts_path_invariant_test.go:94,:104` — **本票自己保留的两例字面反斜杠夹具**，
       由 `TestArtifactsLiteralBackslashKeepsItsAsymmetry` 显式钉住两侧语义（Windows 嵌套 / Linux 扁平），且 target 定义为
       "本平台把它解析成哪儿"，不再是硬编码形状。✅ 本票已修。
     - `internal/tools/pathshape_portable_test.go:110`、`internal/risk/pathshape_portable_test.go:75,:207` —
       票 75 的"可移植路径形状"族：三处**都在 `filepath.Separator` 分支里**（如 `if filepath.Separator != '\\' { ... }`、
       `foreignSepT()`、`home + "\\id_k/x.pem"` 双侧折叠断言），断言的是"这台机器上反斜杠是/不是分隔符"本身 ⇒
       **两侧语义各自明确、已按平台分流，不需要修**（且 `internal/risk/**` 是禁改地界）。
     - （同族旁证：`internal/memory/artifacts_stray_test.go` 票 79 的套件全用 `filepath.Join(arts, "stray", "deep", "hidden.bin")`
       造嵌套，反斜杠只出现在报错文案里 ⇒ 不需要修；本票的修法就是照它对齐的。）
  3. **`filepath.ToSlash` 反用**：把 ToSlash 结果直接喂回 OS 的调用 **0 处**；最接近的一例是
     `internal/tools/wiring_test.go:162` 把 `filepath.ToSlash(target)` 当 `fs.write` 的 JSON 参数——
     那是"另一种分隔符拼写也必须能用"的**正向**探针，两侧同一断言，不需要修。
  另：`internal/agent/spill_path_invariant_test.go` 与 `internal/memory/artifacts_path_invariant_test.go` 里剩下的 `\` 字面量
  （`p\q`、`\\fileserver\share\payload`、`C:\Windows\System32\drop`、`{"/", "\\", ":", ".."}`）都是**喂给编码器的输入串**
  或**对编码结果的负断言**，不构造路径 ⇒ 两侧语义相同，保留。
  next= AC#4 门禁（两包 -count=2 全量 + GOOS=linux go vet，两侧各一次），见 C6
- **C6（AC#4 门禁，只跑本票碰的两个包）**
  - **Windows 工作树**（我的两份文件 = 已提交状态，`git diff --quiet -- internal/agent internal/memory` 当场证干净）：
    `gofmt -l internal/agent internal/memory` → 0 行；`go vet ./internal/agent ./internal/memory` → rc=0；
    `GOOS=linux go vet ./internal/agent ./internal/memory` → rc=0；
    `go test -count=2 -v ./internal/agent ./internal/memory` → **rc=0，RUN=284 PASS=282 FAIL=0 SKIP=2**，两包各 `ok`（agent 4.52s / memory 29.63s）。
  - **Linux docker `golang:1.27` 快照**（`git archive HEAD` 重导，挂快照不挂工作树）：`gofmt -l` 0 行、`go vet` rc=0、
    `go test -count=2 -v` → **rc=0，RUN=284 PASS=282 FAIL=0 SKIP=2**，两包各 `ok`（agent 7.69s / memory 27.33s）。
    两侧用例清单逐条同数：284 = 2×142，过滤集 50 = 2×25。
  - 四种假绿逐跑点名：① `--- SKIP` 2 条 = `TestSubprocessCrashWriter` ×2（`internal/memory/concurrent_test.go:186`
    的子进程角色占位；父用例 `TestCrashRecoveryKillMidWrite` 真跑且 PASS），两侧同名同数、票前既有，不是本票换来的绿；
    ② `-run` 空匹配：五份日志（两侧 gate + 两侧过滤集 + 快照基线）里 `no tests to run` 命中 **0**；
    ③ `-count=N` 已按上面的 2 倍关系核对；④ 静默跳步：docker 侧 `vet && test` 同链（vet 不绿就没有测试输出，而输出在），
    `gofmt -l` 单独核了**输出行数为 0**（日志首行即 `=== RUN`，没有文件名）。
  - 附带（同一枚 commit）：`artifacts_path_invariant_test.go` 里 "four hostile shapes" 的旧文案对齐成 "every hostile shape / six of them"，
    Layer 2 标题处写明**测试函数名里的 "Four" 故意不改**——`docs/evidence/s1/76-adversarial-acceptance.md` 逐字引用它，
    改名字等于替别人重写验收账。
  - **本票没闭的两件事**：(a) 票面 ①② 两条红是在 **ubuntu CI**（run 35558750456）上观测的，本代理不 push，
    所以"CI 侧真的转绿"要编排者推一次才能复验——我的 Linux 证据是本地 docker 快照；(b) 裁决表 `docs/evidence/s1/81-*.md`
    是对抗验收代理的产物，不属写码代理地界。四框都勾了，但 Status 不写 `-done`。
  next= 交回编排者：push 一次让 `test-core` 在 ubuntu 复验这两条；随后派对抗验收出 `docs/evidence/s1/81-*.md`
