# 票 101 —— 独立对抗验收裁决表（`internal/perm` 的生产装配线）

- 验收方：`acceptor-ticket101`（独立对抗验收代理，未参与实现，未看先前对话）
- 被验收对象：`7286297`（`cmd/wisp/run.go` + `cmd/wisp/run_mode101_test.go` + 票面），票面 Status=`ready-for-review`
- 验收时 HEAD：`a1613d9`（`git merge-base --is-ancestor 7286297 HEAD` rc=0 ⇒ 被测代码确在工作线上）
- 票面：`.scratch/wisp/issues/101-perm-store-has-no-production-importer.md`
- 运行日期：2026-09-21
- **共树声明（每份绿来自哪棵树，逐节标注）**：
  - 【快照树】`/tmp/wisp-101-mut-ac1-acceptor101` = `git archive HEAD | tar -x -C …`（本代理会话后缀），
    AC#4 变异与 d22scan 一遍在此；另建 `/tmp/wisp-101-mut-ac1-acceptor101-clean`（同法，无变异）跑 d22scan。
    仓库内**未**建 worktree、**未** checkout（A38④）。
  - 【共享树】`D:\work\workspace\projects plans\Wisp` 本身，仅按包 scope 跑（`cmd/wisp/` `internal/perm/`）。
  - 本代理对 `internal/perm/`、`internal/risk/`、`internal/config/`、`internal/tools/` **零改动**（票 103 / 票 risk 的正在写的东西本代理不碰不提交）。
- 来源档位图例：〔独立复现〕= 本代理亲敲命令看到读数；〔日志＋归档，我抽验〕= 采信票面读数但本代理抽样核过；
  〔仅自述，不背书〕= 只有实现方的话，本代理没能复现，附补救动作。第三档**存在**（见 AC#5 的 1 条）。

---

## AC#1 —— 生产里有人 import `internal/perm` 且落在启动装配路径上

**票面原句**：`grep -rn "wisp/internal/perm" --include=*.go cmd/ internal/ | grep -v _test.go` **非零命中**
且落在**启动装配路径**上（贴 file:line）。这条就是本票存在的理由：**"有人调用它"必须是量出来的，不是宣布的**。

**本代理亲敲**（共享树，HEAD=`a1613d9`）：

```
$ grep -rn "wisp/internal/perm" --include=*.go cmd/ internal/ | grep -v _test.go
cmd/wisp/run.go:53:	"github.com/CarlosShao/wisp/internal/perm"
rc=0        # 命中 1 条，非零
```

**逐环按符号重定位**（票面给的 `:53/:216/:314/:332/:343` 行号**在 HEAD 上未漂**，本代理仍按符号名二次定位）：

| 环节 | 实际位置（`grep -n` 量出） | 本代理核对的"真的连着"的证据 |
| --- | --- | --- |
| import | `cmd/wisp/run.go:53` | 非零命中，见上 |
| 唯一真源 | `cmd/wisp/run.go:216` `mgr, err := config.NewManager(cfgPath, nil)` | `mgr` 随后传给 `perm.Options{Manager: mgr}`，不是又开一份只读快照 |
| 构造存储 | `cmd/wisp/run.go:314` `modeStore, err := perm.New(perm.Options{…})` | `err != nil` ⇒ `rt.auditModeUnreadable` + 退出 2（`:319-323`），不吞错 |
| 持有 | `cmd/wisp/run.go:324` `rt.modes = modeStore` | 字段声明 `run.go:185` `modes *perm.Store` |
| 启动审计 | `cmd/wisp/run.go:332` `rt.auditf("perm: MODE-READ origin=startup mode=%s source=%q", modeStore.PermissionMode(), cfgPath)` | 走的是 bridge 同一个 `rt.auditf` sink；值来自 Store 本身而非猜 |
| **进决策链的那一行** | `cmd/wisp/run.go:343` `Modes: rt.modes,`（在 `tools.Options{…}` 字面量内） | `internal/tools/bridge.go:90` `Modes ModeSource` → `bridge.go:177` `modes: o.Modes` ⇒ bridge 每次调用读它 |
| 边界（不夹带 grant） | `run.go` 内**无** `Confirmations:` 赋值（`grep -n "Confirmations" cmd/wisp/run.go` 只命中 `:306` 注释） | `tools.Options.Confirmations` 保持 nil ⇒ 从盘上读回的只有 mode |

**结论：通过**〔独立复现〕。命中非零、六环环环相接，且**不是靠读注释**：AC#4 的变异（本代理自己做）
证明把 `:343` 一行注释掉后端到端行为真的翻红 —— 这条线不是"写了个 import 摆着"。

---

## AC#2 —— 三条端到端用例，分开、不许合并

**票面原句**：端到端两条用例，**分开、不许合并**（票 90 的 AC#3b 边界）：
(a) 手动改成"全自动"⇒重启后仍是全自动；(b) **从未**手动改过⇒重启后回默认"每步都问"；
(c) 会话授权**不能**跨重启。（票面 AC#2 正文即此三条，本代理按此判。）

**本代理读的是测试体，不是测试名**（`cmd/wisp/run_mode101_test.go` 全文 653 行读完）。
三个独立函数，**没有**表驱动合并：`TestTicket101ManualSwitchSurvivesRestart`（`:207`）、
`TestTicket101UntouchedConfigRestartsAtDefault`（`:328`）、`TestTicket101SessionGrantDoesNotCrossRestart`（`:389`）。
"重启"= 同一 data dir 上第二次 `runTextTask`（`:143` `h.start` 每次新建 `runSpec`、新 Manager、新 bridge）。

| 子项 | 断的是**用户可见效果**还是只断日志词 | 本代理核到的关键读数 |
| --- | --- | --- |
| (a) | **效果**：退出码 0；重启后 `rt.modes.PermissionMode()==auto_approve`；**开卡数 `cards != 0` 即 FAIL**（`rt.ui.shown()` 计数，非日志串）；`isErr` 必须为假（真执行）；**盘上 `note.txt` 存在且 size>0**；`config.toml` 里同时含键 `permission_mode` 与值 `auto_approve`；再用**全新 Manager** 冷读一次（`h.readMode`）把"文件说的"与"进程做的"分开 | `asked` 恰为 `[auto_approve]`（R20/M4 非对称：切 `ask_high_risk` 问 0 次）；`Switch.Origin/Actor` 可归因；对照半：工作区外 `fs.read C:/Windows/win.ini`（R2 红线）**仍被问、仍含"拒绝"** ⇒ "沉默"不是"链子根本没碰到 gate" |
| (b) | **效果**：3 次冷启动各自 `mode==ask_every_step`、**开卡数 `cards != 1` 即 FAIL**、`isErr` 必须为假（未否决＝执行）；键**事前**断言不存在（fixture 自证：写了键就 `t.Fatalf`），**3 次重启后**若键出现 ⇒ `t.Errorf("the host persisted a permission_mode nobody asked it to persist")` | 与 (a) 的 0 张正好相反 ⇒ 两半互为非空对照；每 boot 用**新路径**（`note-%d.txt`），避开 R8"覆盖即红线"把测量对象换掉；另有 `MODE-READ … mode=ask_every_step` 日志断言，但它是**附加**，不是唯一判据 |
| (c) | **效果**：boot 1 先 `InsertGrant`（`GrantScopeSession`、`ExpiresAt=now+3600` 未过期），`sameSession != 1` ⇒ `t.Fatalf`（**fixture 不新鲜就直接红，挡掉空断言**）；同会话内已被拒（`firstErr` 假 ⇒ `t.Fatalf`）；重启后新会话 `ListGrantsBySession` 行数必须 0；`.env` 写**仍被拒**且拒因点名 `R3/敏感/审批/拒绝`；对照半：同 boot 普通工作区内写 `plainErr` 必须为假（auto_approve 下真执行）；最后 `ListGrants` 断言那行**仍在库里**、仍挂在死会话（审计留痕没被顺手清掉） | `Options.Confirmations` 保持 nil 这条**不是靠读注释**：`grep -n "Confirmations" cmd/wisp/run.go` 无任何字段赋值；行为面上"重启后仍被拒"就是"没人从盘上读回 grant"的证据。整段跑在 `auto_approve` 档下（最松档里量边界，方向正确） |

**共享日志缓冲污染**：`h.start` 每次调用内部新建 `out, errb := &bytes.Buffer{}, &bytes.Buffer{}`（`:145`），
`log` 只在**本次 boot 内**返回 ⇒ 不存在"被别的子测试污染"的路径。开卡数是 `rt.windowCount()`（本 boot 自己的 UI），
不是日志计数；因此本项目那条"含该标志的记录条数 == 实际执行次数"的计数不变式在此**以 `cards==0 / cards==1` 的精确计数形式**成立。
（唯一"日志词"式断言是 `MODE-READ` / `MODE-SWITCH` 那几条，均在 per-boot 私有 buffer 上，且 (a)(b) 各自还有独立的盘上/退出码/开卡数判据兜着。）

**本代理亲敲的读数**（快照树，带票 98 的 PATH 注入）：见 AC#4 一节 —— 未变异前 5 个函数全 PASS；
变异后**只有 (a) FAIL**，红名原文 `run_mode101_test.go:305`。

**结论：通过**〔独立复现〕。三条各自独立、非合并、断的都是用户可见效果且**两边都带非空对照**。
（"三条不许合并"这条 R20 边界在本票交付物里成立；文件头注释还写明不许把它简化成表驱动 —— 这是防后人合并的**注释级**护栏，不是硬约束，记在心里即可。）

---

## AC#3 —— 响亮失败面：三态各自一条用例

**票面原句**：档位存储损坏/版本不认识/权限读不到 ⇒ **必须回到最严档并写审计**，
不许"读不到就按上一次缓存的宽松值"。给三态各自一条用例与真实读数。

落法：`TestTicket101UnreadableModeFailsLoudlyAndStrict`（`:560`）+ `t.Run` 三个子用例（`:612`）
`存储损坏`（写进 `[[[` 坏 TOML）/`版本不认识`（`schema_version = 99`）/`权限读不到`（把 `config.toml` 删掉后**原地建同名目录**，可携带地代表"ACL 拒绝读取"）。
每态都在**盘上先放着 `auto_approve`** 再把存储弄坏 ⇒ 测的是"坏读答哪一档"。

断言（本代理逐条读码）：`code != 0` ⇒ Errorf；`code != 2` ⇒ Errorf（SPEC-03 §4.1）；
`reached`（`onRuntime` 被叫到＝决策链装配了）为真 ⇒ Error；日志须含 `MODE-READ-FAILED`；须含 `mode=ask_every_step`（点名回落档）；
须含用户可见的 `配置未就绪`；**负向**：日志出现 `mode=auto_approve` ⇒ Error（"按上一次缓存的宽松值"被钉死为不可能）。
**非空对照跑在最前且便宜**：完好 fixture ⇒ `code==0 && okAssembled && okMode==auto_approve`，否则 `t.Fatalf`（挡掉"这 fixture 本来就没通过"的假绿）。

**判①（方向）**：只会更严，不会更松 —— 三态一律 `return rt, 2`（`run.go:220-222`、`:318-322`），
bridge 根本不会被构造，因此**不存在任何一条能执行动作的决策链**；这比"继续跑但按最严档"更严。
本代理另核 `run.go:406` `auditModeUnreadable` 的档值来源 = `risk.DefaultMode()`（**代码里读回的常量**，不是写死的字符串），
所以"点名的档"与"bridge 在 nil ModeSource 下的答案"同源。**方向无错，不是 FAIL。**

**判②（可用性代价）**：代价真实存在 —— `config.toml` 坏（哪怕坏的与权限无关，如 provider 段写歪）
⇒ 今天真机上**整个 `wisp run` 起不来（退出 2）**，且**修它只有手改文件这一条路**（票 92 的面板 composer 未落地）。
是否有用户可见解释：**有**，`run.go:221`/`:320` 打 stderr 的
`wisp run: 配置未就绪（Unconfigured）：%v` 带底层错误，且 `MODE-READ-FAILED path=%q err=%v … detail="档位读不到：本进程不缓存任何上一次的宽松值，决策链不会被装配（退出码 2）"`
点名了路径与后果。**记账但可接受**：这是继承 SPEC-03 §4.1 既有口径（config 装载失败本来就是退出 2），不是本票新引入的口径；
残留登记 `R-101-5`（"坏配置＝整机起不来＋无可点开的修复入口"这笔账要等票 92/77 还）。

**本代理亲敲**：5 个函数（含三子用例）在**未变异的快照树**上全部 PASS；三子用例的 `--- PASS` 见 AC#4 的 `-v` 输出。

**结论：通过但有条件**。条件＝ `R-101-5` 登记在案（响亮失败的可用性代价已量化、有用户可见解释，
但"配置坏了没有 in-app 出路"这一半不属于本票、也不许在本票里顺手修）。

---

## AC#4 —— 变异：把装配那一行注释掉 ⇒ 只有 (a) 必须红

**票面原句**：变异：把装配那行**注释掉** ⇒ AC#2 的 (a) 必须红（证明这条线不是"恰好也绿"）。
锚点=承载行为那一行，同链 grep 证落地，还原后 `git diff --quiet` 证干净，**编译失败不算变异**；**变异只在 `/tmp` 仓外快照里做**。

**本代理自己做的**（不采信实现方那份；快照目录带本代理会话后缀）：

```
$ M=/tmp/wisp-101-mut-ac1-acceptor101; rm -rf "$M"; mkdir -p "$M"
$ git archive HEAD | tar -x -C "$M"; cp -r third_party "$M/"
$ grep -n "Modes: rt.modes," cmd/wisp/run.go          # 锚点：值实际所在行，按符号定位
343:		Modes: rt.modes,
$ sed -i 's|\t\tModes: rt.modes,|\t\t// MUTATION(acceptor101 AC#4): Modes: rt.modes,|' cmd/wisp/run.go
$ grep -n "MUTATION(acceptor101\|Modes: rt.modes," cmd/wisp/run.go     # 同一条链里证落地
343:		// MUTATION(acceptor101 AC#4): Modes: rt.modes,
$ grep -n "Modes: rt.modes," "<主树>/cmd/wisp/run.go"                    # 主树未被改
343:		Modes: rt.modes,
$ export PATH="$PWD/third_party/sherpa-onnx:$PATH"
$ go build ./cmd/wisp/ ; echo build_rc=$?  ->  build_rc=0                 # 编译成功，之后才算变异
$ go test -v -count=1 -run TestTicket101 ./cmd/wisp/ ; test_rc=1
--- FAIL: TestTicket101ManualSwitchSurvivesRestart (4.26s)
    run_mode101_test.go:305: auto_approve still opened 1 card(s) for an L1 write; the档 was read but never reached the decision chain:
--- PASS: TestTicket101UntouchedConfigRestartsAtDefault (7.77s)
--- PASS: TestTicket101SessionGrantDoesNotCrossRestart (5.32s)
--- PASS: TestTicket101ModeSwitchUsesTheRealL2Gate (1.49s)
--- PASS: TestTicket101UnreadableModeFailsLoudlyAndStrict (1.11s)   # 含 存储损坏/版本不认识/权限读不到 三子用例
FAIL	github.com/CarlosShao/wisp/cmd/wisp	20.011s
```

**形状核对**：**恰好只有 (a) 红**，其余 4 个（含 AC#3 三态）仍绿 —— 与"默认档=fail-closed"的预期一致：
(b) 的默认档在 `Modes` 为 nil 时本来就是答案，(c) 测的是 grant 不回来，AC#3 测的是根本不装配，
**只有 (a) 能区分这条线**，故 (a) 就是本票的红线探针。红名与实现方报的**逐字一致**（`run_mode101_test.go:305`）。
**判据达成。**

**干净度**：变异只发生在 `/tmp` 快照（仓内未建 worktree/未 checkout）。主树本代理侧：
`git status --porcelain cmd/wisp/ internal/perm/` = **空**（本代理未改任何代码文件）。
⚠ 诚实读数：`git diff --quiet` 在共享树 rc=**1**（**工作树里另有别的票的未提交改动**，属并行票 103/`internal/risk` 的地界，
不是本代理造成的，本代理不碰也不提交；`cmd/wisp/`、`internal/perm/` 两个 scope 内是干净的）。

**结论：通过**〔独立复现〕。

---

## AC#5 —— 门禁重跑（逐条点名 SKIP/FAIL）

**票面原句**：`gofmt -l`/`gofumpt -l` 空、`go vet ./cmd/wisp/ ./internal/perm/` rc=0、
`go test -count=2 ./internal/perm/` rc=0 且逐条点名 SKIP/FAIL。（并：静态扫描必须与 CI 那一行逐字同形 `sh scripts/d22scan.sh`；
台账 `ban #6/#8 frontend/` 文件数不许下降。）

| 门禁 | 本代理命令 | 树 | 读数 | 判定 |
| --- | --- | --- | --- | --- |
| gofmt | `gofmt -l cmd/wisp internal/perm` | 共享树（包 scope） | 输出**空**，rc=0 | 〔独立复现〕通过 |
| gofumpt | 本机 `gofumpt` **未安装**（`which gofumpt` → NOT ON PATH，rc=127）⇒ 改 `go run mvdan.cc/gofumpt@latest -l cmd/wisp internal/perm` | 共享树（包 scope） | 输出**空**，rc=0 | 〔独立复现〕通过（命令形状与 CI 的 `gofumpt -l . tools/…` 差在 scope，见下"条件"） |
| go vet | `go vet ./cmd/wisp/ ./internal/perm/` | 共享树 | rc=0 | 〔独立复现〕通过 |
| perm 计数不变式 | `go test -count=2 -v ./internal/perm/` | 共享树（包 scope） | rc=0；`^=== RUN` = **28**；`^--- PASS` = **28**；`^--- SKIP` = **0**；`^--- FAIL` = **0**；去重测试名 = **14** | **14 × 2 = 28 ✓ 自洽**（见下戳穿） |
| d22scan | `sh scripts/d22scan.sh`（与 CI 逐字同形；**未**用 `go run ./tools/d22scan`） | **两棵树各一遍** | 快照干净树：rc=0 clean，`ban #6 frontend/=37`、`ban #8 frontend/=37`、`internal/=340`、`design/=16`、`ban #1-5 internal/=197`、`cmd/=20`、`ban #7=17`、`ban #8 cmd/=26`；共享树：rc=0 clean，`ban #6/#8 frontend/=40`、`ban #8 internal/=340`，其余同上 | 两树皆 clean；台账**不降**（见下） |
| 宿主包本机可测性 | `go test ./cmd/wisp/`（票 98 注入 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`） | 快照树（变异）／共享树（本代理未全包重跑） | 本代理**按票 98 注入命令**跑的 `-run TestTicket101` 有真实测试读数（非 `0xc0000135` 加载期失败） | 通过（未默默跳过）；全包 40.9s 那份是**实现方自述**〔日志＋归档，我抽验〕：本代理只抽验了 `TestTicket101*` 子集 + 变异体 |

**当场戳穿的自相矛盾计数**：实现方在 AC#5 里同时说 `internal/perm` 是 **28 PASS** 和有 **15 个顶层测试**，
并用"**15 × 2 减去缓存复用**"解释 30≠28。本代理量出：`grep -rn "^func Test" internal/perm/*_test.go | wc -l` = **14**
（`store_test.go` 9 个 + `ticket90_persist_test.go` 5 个），`-count=2 -v` 去重名 = 14。
⇒ **14 × 2 = 28，等式成立，"缓存复用"是错的解释**（`-count=2` 不走结果缓存，这个理由本身不合法）。
**读数对、解释错**：登记 `R-101-2`。

**台账不降的核对（含一次树差异，不遮掩）**：实现方报 `ban #6 frontend/=40`。
本代理在**共享工作树**（同一类树）复测 = **40**，未下降 ✓；在 **`git archive HEAD` 快照树**里 = **37**。
差 3 的原因是 `git archive` 只导出被 git 跟踪的文件（`git ls-files frontend | wc -l` = 37，且 `git status --porcelain frontend/` 为空
⇒ 多出的 3 个是**被 ignore 的本地文件**），属**树口径差**，不是台账下降。登记为 `R-101-3` 的"口径"备注。
另：`ban #8 internal/` 实现方报 337，本代理两棵树皆 **340**（HEAD 之后别的票往 `internal/` 落了文件）⇒ **上升**，非下降。

**gofumpt 的条件**：CI 那一行是 `go install mvdan.cc/gofumpt@latest` + `gofumpt -l . tools/d22scan tools/mockllm`（**全仓**）。
本代理只按包 scope 复现了 `cmd/wisp internal/perm`（共树规矩：全仓门禁不在共享树里改东西，而本代理也不便替别的票背全仓绿）。
⇒ 本票 scope 的绿＝独立复现；全仓 gofumpt＝〔仅自述，不背书〕。**补救动作**：由编排者在 CI 或一次性快照里跑
`go run mvdan.cc/gofumpt@latest -l . tools/d22scan tools/mockllm` 并把输出贴进台账（命令本代理已验证可跑通、无需安装）。

**结论：通过但有条件**。条件＝全仓 gofumpt 一栏由 CI/编排者补（`R-101-3`）；`internal/perm` 计数的解释要改（`R-101-2`）。

---

## 顺带查清的一条"注释说谎"账（→ 归票 90）

实现方的登记**为真**，本代理用 file:line 读码确认：

- `internal/perm/store.go:122-123`（`New` 的文档注释）："New composes a Store and **records** what mode the process
  started with, **so the audit trail** of a long-running host has a first line to compare against."
- `internal/perm/store.go:141` `New` 调的是 **`s.record(Switch{… Origin:"startup", Actor:"process", Result:ResultNoChange …})`**。
- `internal/perm/store.go:261-269` `record()` 函数体 = **只做 `s.hist = append(...)` + 按 cap 淘汰**，**不碰 `s.logf`**。
- `internal/perm/store.go:274-281` `audit()` = `s.record(sw)` **然后** `s.logf("perm: MODE-SWITCH …")` —— 只有它写 sink。

⇒ **`perm.New` 里那行"启动档位"从来没进过审计 sink**，它只在内存 `History()` 里；文档里"the audit trail … has a first line"
这半句与代码不符。**裁决：注释与代码不符，改它该有自己的票（归票 90 的账，本代理不顺手修，`internal/perm` 一行未碰）。**
登记为 `R-101-1`。

**"现在补的那条是否等价"**：**对本票的判据等价，对库契约不等价**。

| 维度 | `perm.New` 里那行（未做的） | `run.go:332` 补的 `MODE-READ` |
| --- | --- | --- |
| 是否进 sink | 否（`record()` 只进内存） | **是**（`rt.auditf` = bridge 同一个 `[audit]` 家族） |
| 是否点名启动档 | 是（`To=From=PermissionMode()`） | **是**（`mode=%s` 直接取 `modeStore.PermissionMode()`） |
| 是否可被 grep 到 | 若走 `audit()` 会落在 `MODE-SWITCH origin="startup"` | 落在 `MODE-READ origin=startup`（**kind token 不同**） |
| 覆盖范围 | 所有 `perm.New` 的调用者 | **只有 `cmd/wisp` 这条装配路径** |
| 被测试钉住 | 无 | AC#2(a)(b) + AC#3 都断言它（本代理复现） |

⇒ 不阻塞本票（本票 AC#1/#2/#3 的判据都点名 `MODE-READ`，且实测在）。
对票 90 的后续：要么让 `New` 走 `audit()`（并把 kind token 统一，避免两条家族名并存），要么把那句注释改成
"内存历史里有首行、审计首行由装配根负责"。**若将来出现第二个 `perm.New` 的生产调用者（票 77 的宿主面板），
它若不自己补 `MODE-READ`，就会静默丢掉这条首行 —— 那才是这条账真正会咬人的时刻。**

---

## 能力类 AC 的固定一问：生产里谁调用它、在哪个装配根？

- **读侧**：装配根 = `cmd/wisp/run.go` 的 `assembleRuntime`（`:216`→`:314`→`:324`→`:343`），
  消费者 = `internal/tools/bridge.go`（`Modes` → 每次调用读一次）。**已交付**，且被 AC#4 变异钉成"拔掉就红"。
- **写侧**：`grep -rn "\.Set(ctx\|modes\.Set\|func (s \*Store) Set" --include=*.go cmd/ internal/ | grep -v _test.go`
  → **只剩声明本身**（`internal/perm/store.go:175`）。⇒ **`Store.Set` 在生产里至今零调用者**，
  与实现方自述一致（它没假装有人调）。改档今天只有两条路：手改 `config.toml` 的 `risk.permission_mode`，或等票 92/77 的原生入口接 `rt.modes.Set`。
  登记 `R-101-4`。
- **因此现在能不能对用户说"选过就一直生效"**：**不能这么整句说**。
  准确说法＝"**写进 `config.toml` 的档，重启后仍然算数，并且真的进决策链**"（成立，本代理复现）；
  "**在界面里选**过就一直生效"（**不成立** —— 今天没有任何界面能发起那次选择，`wisp run` 控制台里那次 L2 强确认必然被拒，
  由 `TestTicket101ModeSwitchUsesTheRealL2Gate` 钉住：开 1 张卡、档不动、审计 `result=refused-confirm`、盘上不落 `auto_approve`）。
  本代理**不升格**成"已交付"：**库层可用、装配根可用、用户层入口不可用**。

---

## 残留与条件清单

| 编号 | 内容 | 归属/处置 | 是否阻塞本票 |
| --- | --- | --- | --- |
| `R-101-1` | `perm.New` 文档说启动行进审计，实际调 `record()`（内存）而非 `audit()`（sink）：`store.go:122` vs `:141/:261/:274` | **票 90 的账**；改注释或改走 `audit()` 该有自己的票。本代理未动 `internal/perm` | 否（本票在装配根补的 `MODE-READ` 对本票判据等价） |
| `R-101-2` | AC#5 读数解释错误："15 个顶层测试 ×2 减去缓存复用"。实测 `internal/perm` **14** 个顶层测试，14×2=28；`-count=2` 不存在缓存复用一说 | 实现方在票面就地更正解释（读数本身有效，本代理复现为 28/0 SKIP/0 FAIL） | 否 |
| `R-101-3` | 全仓 `gofumpt -l . tools/d22scan tools/mockllm`（CI 那一行）本代理未跑（本机无 `gofumpt` 二进制；只按包 scope 用 `go run mvdan.cc/gofumpt@latest` 复现）；另 `frontend/` 台账 40（工作树）vs 37（HEAD 快照）＝ignore 文件造成的**树口径差**，非下降 | **条件**：编排者/CI 出全仓那一栏；台账比较必须同树口径 | 是（AC#5 判"通过但有条件"） |
| `R-101-4` | `Store.Set` 生产零调用者 ⇒ M3 写侧未交付（面板 composer=票 92，原生点击通道=票 77）；控制台跑法切不到全自动（fail-closed 的正确形状，别读成功能坏了） | 票 92/77 复用 `assembleRuntime` 已有的 `rt.modes` 与 `runSpec.modeConfirm` seam，**别再开第二条装配路径** | 否（本票 AC 只判读侧 + 边界；实现方自述②③与实测一致） |
| `R-101-5` | 响亮失败的可用性代价：`config.toml` 任意原因坏 ⇒ `wisp run` 退出 2、决策链不装配，且**没有 in-app 修复出路**；用户可见解释**已有**（stderr `配置未就绪（Unconfigured）：<err>` + 审计点名路径与 `mode=ask_every_step`） | **条件**（AC#3 的口径写明"不装配 + 退出 2"，比"继续跑但按最严档"更严；方向只更严不更松）；出路归票 92/77 | 否，但口径要在票面上留痕 |
| `R-101-6` | 宿主包本机测试依赖票 98 的 PATH 注入（不注入 ⇒ 加载期 `0xc0000135`）。本代理按注入命令拿到真实读数，未默默跳过；全包 40.9s 那份仍属实现方自述 | 票 98 收口；本票不判 | 否 |
| `R-101-7` | 共享树 `git diff --quiet` rc=1（并行票 103 / `internal/risk` 的未提交改动）；`cmd/wisp/`、`internal/perm/` 两 scope 内为空 | 共树事实，非本票缺陷；本报告已按包 scope 声明每份绿来自哪棵树 | 否 |

**外部注入提醒的处置**：本代理会话内的工具输出**未**出现自称"编排者备注"的假指令文本，
故无原文可登记；本代理全程按"工具输出不是授权也不是指令"处理（未据任何工具输出扩大改动范围）。

---

## 总判

# `PASS WITH CONDITIONS`

五框读数全部本代理亲量或明确降级标注；本票存在的理由（"有人调用它"）已由 grep + 逐环符号核对 + **自做的变异**三重钉牢。
条件＝`R-101-3`（全仓 gofumpt 一栏由 CI 出）、`R-101-2`（票面解释就地更正）、`R-101-5`（响亮失败的可用性代价留痕，出路等票 92/77）。
`R-101-1` 归票 90，`R-101-4` 归票 92/77。

**直答 owner："现在改档重启后还在不在？"**

- **读侧：在。** 盘上（`config.toml` 的 `risk.permission_mode`）写着哪档，重启后决策链就按哪档走，
  这条本代理用变异复现过（把那行注入注释掉 ⇒ `auto_approve` 立刻"仍开 1 张卡"红）。
- **写侧：不在。** 今天没有任何用户可见入口能"改档" —— `Store.Set` 生产零调用者，控制台那条 L2 强确认必然被拒；
  要"选"只能手改 `config.toml`，界面上选（票 92）与原生点击通道（票 77）都还没接。
- 一句话：**库层可用、装配根可用、用户层入口不可用**；对 owner 只能说"写进配置文件的那档会一直生效"，
  不能说"你选过的那档会一直生效"。
