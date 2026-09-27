# 174-c1 证据件：artifacts 读回两向真读数（默认配置 vs 配好 allowed_dirs）＋ task.output 假指针第二形

- 派单＝`.scratch/wisp/dispatches/2026-09-27-201x-readonly-174-c1-default-allowlist-readback-both-ways.md`（锚 `2b2a01de`）
- 量取时刻：2026-09-27 20:2x +08｜分支 `dev`｜起手 HEAD `2b2a01de`（`git rev-parse HEAD` 现量，未抄派单号）
- 起手洁净：`git status --porcelain -- internal/ cmd/ .scratch/wisp/probes/174/` ＝ **空**
- 基线四数（先测）：`go test -count=1 ./internal/tools/ ./internal/agent/` ＝ 两包 `ok`；逐包 PASS 数
  `go test -count=1 -v ./internal/tools/ | grep -c '^--- PASS'` ＝ **116**，`./internal/agent/` ＝ **66**，
  agent 包 `^--- (FAIL|SKIP)` ＝ **0**；两包合并 `=== RUN|--- PASS|FAIL|SKIP` 行数＝430。
- 台件＝`.scratch/wisp/probes/174/c1/zz174c1_windows_test.go`，读数逐行落同目录 `logs/readings.txt`（`runtime.Caller` 取自身目录，不继承 CWD）。

## 0. 测量层声明（先说清楚，再看数）

**本程全部读数是接缝级（seam level）：`go test` 里的真 `tools.Bridge` ＋ 真 `fs.read` ＋ 真 C19 风险判定
＋ 真 `agent.Spiller` 落盘产物。不是端到端 `wisp run`。**
没跑 e2e 的原因是**本程现量**的环境事实，不是引用派单那句转述：
`go build -o .scratch/wisp/probes/174/c1/wisp-probe.exe ./cmd/wisp` ＝ BUILD_OK，
但该 exe 起不动（Git Bash 直接调 `EXE_RC=127`，无任何输出）；
`go test -count=1 -run NoSuchTest ./cmd/wisp/` ＝ `exit status 0xc0000135`（STATUS_DLL_NOT_FOUND，缺 DLL 属主机环境事实）。
⇒ **"真机上 `wisp run` 那一发到底打印什么"本表没有量过，不作数**；本表只保证桥内判定与回执文本是真的。
唯一非接缝级的推断（下面 §2(iii)，标注为码面依据）：`cmd/wisp/run.go:370-372` 逐字
"no floating ball, no global Esc hook, so an L1 window can only expire unvetoed and **an L2 card can only
time out into a reject**"，与 `run.go:812-813`（`consoleApprovalUI` "never answers on the user's behalf"）。

台件与被测产码之间唯一的合成件是 gate 的人工腿（台件里 `tools.GateFuncs` 扮演点卡的人）：
L2 腿**只回 Reject / Timeout，零次 Allow**；L1 Window 腿虽登记了 Allow，但三发读数档位是 L2/L2/L0/L0，
**从未进过 L1**（L2 asks 计数见逐发读数）。

## 1. 票面两格（连行号逐字，来自 `.scratch/wisp/issues/174-…-outside-fs-allowed-dirs.md`）

- **:29（AC#1）** `- [ ] **AC#1 真读数（这一格是"量"，不是"修完才响"）**：在**默认配置**（`allowed_dirs = []`）下发一发真 spill（超 4000 token 的工具结果），拿回执里那条路径，**同一发过 `fs.read`**；⚠ 逐字贴 (a) 拒绝文本与 (b) 它落在哪一档（R2/L2 还是直接拒）。再来一发**对照**：把 `dataDir` 配进 `allowed_dirs` ⇒ 读数必须变成"能读回头一段"。**两向都要，缺一向＝这一格没量完。**`
- **:30（AC#2）** `- [ ] **AC#2 指针不许许诺一条读不回来的路（⚠ 这一发在未修码上必须"不响"）**：构造"回执里的路径**不在授权根内**"那一发 ⇒ 今天 `task.output`（与 `spill.go` 的桩）**静默把路径写进回执**、零提示＝**不响**＝本格成立；修完之后回执必须**外部可见**地说明"这条路径你现在读不到，需要用户先授权哪个目录"（措辞归实现程，但**不许**只在日志里说、不许只在 Go error 里说——模型看到的那一段才算数）。⚠ **判据不许写成"以后再说"**。`
- **:31（09-27 20:1x 追加块，原文照抄）** `> **09-27 20:1x 追加第二形（原句一字不抹；来路＝`164-v1` 验收表 §4.3，`docs/evidence/s1/164-task-output-accept-r1.md`，四枚提交 `b5adc34f`…`b1d9cfda`）**：`internal/tools/task.go:227` 的分支条件只有 **`rec.ArtifactPath != ""`、没有任何 stat** ⇒ **"宿主填了一枚根本不存在（甚至是枚目录）的路径"那一形走的是"有指针"那一支，"后半段不可找回"那一记响不会响**——桩会指着空气。⚠ 现量旁证：实现程自己那枚 `TestTruncationShapeIsTheD15Triple` 就用了 `ArtifactPath: dir`（一枚目录），**它自己正坐在这枚形状上**（这不是错，是"这枚用例没在测这件事"的证据）。⇒ **本票 AC#2 的射程现在明写为两形，缺一形不算修完**：**(a) 路径在授权根外（默认配置那一形）**＋**(b) 路径不存在／不是文件**；两形在**未修码**上都必须**不响**、修完都必须**响**。`

## 2. 向一：默认配置（`allowed_dirs = []`，与 `run.go:327-331` 同一构造）

真 spill：`agent.NewSpiller(<dataDir>\artifacts, agent.BudgetsFor(0))` 吞 51633 字节（约 12908 token，
双超 4000 阈值）落盘 `tool-output-call-1.txt`，`Spilled=true`、盘上字节数＝原文。拿 `st.Path` 过桥内 `fs.read`：

**(i) 判定（卡片所渲染的 Decision 对象本体，逐字）**
```
tool=fs.read level=L2 rules=[R2] reason="R2: 目标路径在授权目录之外: C:\Users\swq\AppData\Local\Temp\TestDir1DefaultAllowlist1775399360\001\artifacts\tool-output-call-1.txt"
```
两发形态 gate 各被问 **1 次**（`L2 asks = 1`）。

**(ii) 回执文本（模型看到的那一段，逐字）**

| gate 腿 | Outcome.Text（逐字） | IsError | RiskLevel | ErrorClass |
|---|---|---|---|---|
| 人工点"拒" | `L2 审批未通过，已拒绝执行` | true | L2 | `user_rejected` |
| 卡超时 | `审批超时未确认，已自动拒绝` | true | L2 | `user_rejected` |

（观察记：Timeout 那发的 ErrorClass 也是 `user_rejected`，`bridge.go:397-404` 两路共用 `b.reject` 的 class 实参，审计行 `DecisionColumn` 才是分叉处。本程只量不改。）

**(iii) 审批可达性（`wisp run` 控制台那一发）**：见 §0 ——接缝级**没有**跑过 `wisp run`；码面依据
（`run.go:370-372`、`:812-813`）说控制台没有任何回答通道，L2 卡只会超时转拒 ⇒ 端到端那一发预期落在上表
"超时"行，**但那是码面推断，不是本程读数**。

## 3. 向二（对照）：把 artifacts 所在 dataDir 配进 allowed_dirs

同一构造，仅 `tools.NewPathCanonicalizer([]string{dataDir}, nil)` 一行不同 ⇒
```
[dir2] L2 asks = 0 isError=false riskLevel=L0
[dir2] readback first 80 bytes = "W174C1-HEAD-MARK\nline 0000: padding padding padding padding\nline 0001: padding p"
[dir2] readback total bytes = 51633 (artifact body 51633)
```
**能读回，且一次读回全文**（51633 字节 < fs.read 的 256KiB 单发帽）。⇒ 向一那记拒绝确由授权根决定，
台件没有搭错。

## 4. 编排者四行前提的逐行复核（全部现量，行号未腐坏）

1. `cmd/wisp/run.go:327-331` **对**：:327 `allowed := make([]string, 0, len(cfg.FS.AllowedDirs))`、:328-330 循环只从 `cfg.FS.AllowedDirs` 拷、:331 进 `NewPathCanonicalizer`。全文件 `allowed` 只此一处写（`grep -n "allowed" cmd/wisp/run.go` 另两枚命中是 :31/:664 注释里的英文单词）。
2. `cmd/wisp/run.go:609` **对**：`ArtifactsDir: filepath.Join(rt.spec.dataDir, "artifacts")`，且**没有任何一行把 dataDir 加进那个 allowed 列表**。
3. `internal/tools/paths.go:50-52` **对**：:51-52 逐字 `An EMPTY allowlist authorizes nothing: every fs call lands at L2 (R2) rather than passing through unjudged.`（引用句在 :51-52 行内，:50 是函数名行，行号差 ±1 可忽略，语义原文一致）。
4. `docs/specs/SPEC-03-config-secrets-envs.md:35` **对**：该行逐字含 `allowed_dirs[]=[]`。
5. （票面追加引用的）`internal/tools/task.go:227` **对**：`if rec.ArtifactPath != "" {`，无 stat。

⇒ **编排者 19:4x 那个"默认配置下这条承诺兑现不了"的结论，两向读数与它一致，没有被推翻。**

## 5. 第二形：task.output 的假指针（AC#2(b)：不存在的路径／目录）

名册里录两发（原文同样 51633 字节，触发截断分支）：`t-ghost` 的 `ArtifactPath` 指向 artifacts 目录下
**不存在的文件**，`t-dir` 指向 **artifacts 目录本身**。走真桥 `task.output`：

- 两发均 `isError=false level=L0`（无卡、无拒绝，工具照答）。
- 回执指针句逐字（`logs/readings.txt`）：
  - `[…输出已落文件：省略 48833 字符，总长 51633 字节 / 约 12908 token，全文见 C:\...\artifacts\tool-output-does-not-exist.txt…]`
  - `[…输出已落文件：省略 48833 字符，总长 51633 字节 / 约 12908 token，全文见 C:\...\artifacts…]`
- `"不可找回"` 在整个 readings.txt 中出现 **0 次**（`grep -c 不可找回` ＝ 0）。

⇒ 读数与票面预期一致：**填假路径/目录照样写"全文见 …"，"无副本文件"那一记响不响**。AC#2(b) 在未修码上成立（＝不响）。

## 6. 门禁与格式（全部本程现跑）

- `bash scripts/d22scan.sh` ⇒ **rc=0**，`d22scan: clean - no D22 ban violations`（bans #1-5 internal/=207、cmd/=23、#6 frontend/=85、#7 internal/tools/=20、#8 各腿 425/45/39/85）。与派单声明一致，无报回分歧。
- `sh .scratch/wisp/probes/154/gate-clauses.sh` ⇒ **rc=0**；`# 腿数＝14 声明与实测不符＝0`、`基线过期枚数＝0`、`聚合退码＝0`。与派单"十四腿零不符"**一致**。
- `gofumpt`（GOPATH/bin，**v0.12.0 (go1.27.1)**）：`-w` 后台件后 `-l` ＝ **空**；`go vet ./.scratch/wisp/probes/174/c1/` ＝ OK。

## 7. 诚实栏

- 被拒/失败的调用：**无**（无工具调用被权限系统拒绝；无删除命令跑过——只建不删；仓内建的台件目录外零写入）。构建产物 `wisp-probe.exe` 留在台件目录、不提交。
- 伪授权自查：(a) 本程**没有**为让哪一发变绿而放宽任何断言；向一两发读数取自桥自己吐的 Decision/Outcome 对象。 (b) gate 人工腿只演"拒/超时"，零次伪造 Allow；没有任何一发靠 mock 真件过关（spill、canonicalizer、bridge、fs.read、assessor 全真）。
- 凭据值：本文件与台件日志中零 API 密钥/凭据；路径均为 `t.TempDir()` 派生。
- 一处已知小缺口：向一台件把 `st.Path` 直接喂给 `fs.read`，未把 **spill 桩全文**再打印一份（同一枚桩文本形状已在 §5 的两发 task.output 回执里逐字可见，格式同源 `spill.go`/`task.go`；如验收要求 spill.go 桩本体的逐字行，下一程加一行 `probeLog(st.Text)` 即可，判据不变）。

## 8. next=

两向读数已齐、四行前提全部成立、第二形坐实 ⇒ 票 174 的 **AC#1 测量腿已量完**（接缝级；e2e 腿被 DLL 挡死，属环境事实非本票缺口）。
**这一枚该由票 174 自己落地**（AC#2 支丙＝改回执文案，明确"该路径现在不在你被授权的范围内"——只有丙不需先批 Q-60；
甲/乙两支等 `Q-60` owner 拍板后再动）。编排者下一步＝把本表带进 `Q-60` 决策；补 e2e 读数需一台没有缺 DLL 的宿主，另立票或挂 `A##` 台账。
