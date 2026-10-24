# 369-r1 落地腿（写腿）— 两枚 Windows 用例：**胜率判据 → 确定性判据**

- 腿：`369-r1`＝**写腿**。⛔ 翻框／⛔ 改票名·`Status:`·`-done`／⛔ 改台账／⛔ push／⛔ 动 `internal/winsec/winsec.go`。
- 锚：开工 `HEAD=e485af5a`（2026-10-24 12:08:42 +08）。读数期间别人又落了 `4997c59d`（票 376 普查件，**未触** `internal/winsec/**`）。
- 台面：本机 `go1.27.1 windows/amd64`；`GOFLAGS=`；**只** `unset GOTMPDIR`（`TMPDIR/TEMP/TMP` 一字未动）。
- 写面：`internal/winsec/{openleaf_guard_322_windows_test.go, ancestorunlink_348_windows_test.go, leafrace_342r2_windows_test.go, routeguard_369_test.go}`＋本件；原始件逐枚在 `D:\tmp\wisp\20261024\369r1\{logs,mut,probe}\**`（⛔ `.out`）。

## 0. 一句话

★两枚的"守卫腿/re-check 命中过 ≥1 次"**不是 landing 的性质、是机器时钟的性质**（同一 head 两发 `0` ⇄ `1`）⇒ 判据换成**三件确定性东西**：①route 的**接线**（读 `winsec.go` 的调用图）②leg 的**行为**（直接驱动，不经 race）③route 的**安全**（叶子挂链时 `PrivateFile` 永不返回 nil，且回答归到稳态那一行）⇒ `-count=10` **10/10 一致**、坏腿变异**仍红**、正控绿。

## 1. 开盘读数（自己现跑，⛔ 抄）

| 量 | 读数 |
|---|---|
| 两发 head | `gh run view … --json headSha` ⇒ **两发同为 `23b663c2dc70…`** |
| 发 `39884431726` step 4 | **failure**；342 枚 `guard=0, re-check=0, provedReparse=63`；322 枚 `guardRefusals=0 flips=411 otherRefusals=249 returnedNil=151`；四数 `RUN=122 PASS=69 FAIL=2 SKIP=0`；`gate is RED (rc=1)` |
| 发 `40034438521` step 4 | **success**；`guard=1, re-check=0, provedReparse=65`；`guardRefusals=1 flips=521 otherRefusals=195 returnedNil=204`；`RUN=122 PASS=71 FAIL=0`；`ok … 11.912s` |
| 本机 `-count=5`（改前） | **`RUN=10 PASS=7 FAIL=3 SKIP=0`**；342 枚逐发 split `(0,1)P (0,0)F (1,0)P (0,0)F (0,0)F`；322 枚 `guardRefusals=12/11/10/8/15` 全 PASS |
| build tag | `ancestorunlink_348_windows_test.go:1`＝`//go:build windows`；`openleaf_guard_322_windows_test.go:1`＝`//go:build windows` |

⇒ ★**0 命中红、1 命中绿**：本机也出 0（`369-a1` 的 `-count=5` 里就有 0 与 1），本腿复跑亦然。

## 2. 机制：为什么"守卫腿被撞上"是胜率（精确到行）

`privateFile` 的顺序（`winsec.go`）：`:163 resolveString` → `:165` 返回 → `:167 preOpenLeafGuard` → `:179 pinnedParent` → `:187 openPinnedLeaf`（内含 re-check）→ `:196 sealHandle`。

★★**`:163` 的两条稳态腿都把"叶子"当 prefix 走**：
1. `resolve.go:605 resolveAccounted` → 装好的 C26（本包测试二进制里是 `risk.c26Pipeline`，启动行 `winsec: sealing path resolver installed`）⇒ 叶子挂链答 `path traverses a reparse point`；
2. 紧接着 `resolve.go:631` **无条件**再跑地板：`builtinVerifier.Resolve`（`:664`）→ `:674 platformVerifyPlacement` → `placement_windows.go:59-68` 逐 prefix `isReparsePoint`，而 **`winsec.go:397-421 pathPieces` 的最后一片就是整条路径**（`:414` append `path[:i]`）⇒ 叶子同样被走。

⇒ ★**叶子一挂链，`:163` 就返回了**，`:167` 的 guard 与 `openPinnedLeaf` 的 re-check（`pinnedleaf_342_windows.go:199-211`）**只能吃到"翻在 [叶子那次 Lstat → 自己那次 CreateFile] 之间"的两 syscall 那一格**。而 `:179`+`:187` 的窗口**严格包含**这一格 ⇒ **guard 命中 ⊆ 链命中**。

★现量（12:15 那发，400 attempts，`ancestorunlink_348_windows_test.go:107` 的世界，七个行合计＝400）：稳态 reparse 行 **249**、`path escapes from parent`（链）**75**、seal **23**、`resolveString refused the spelling` **9**、`could not read the name` **1**、`some other failure` **6**、`returnedNil` **37**；guard 链接形 **0**、re-check **0**；`provedReparse` 144、`truncated` 2、`deleted` 0 ⇒ **75 ⇄ 0** 就是"链吃掉 guard 那一格"的形状。
★★所以 `guardRefusals ∈ {0,1}`（runner）不是样本量问题：**加样本量不会把它变成确定**，因为它的期望值本身就随机器而变（本机 0~15）。

## 3. 选形：为什么是这一形（引普查哪几条）

- 引 `369-a1-census.md:31`（同 head 0⇄1）、`:38`（本机 `11/2/8/0/1`）、`:39`（稳态 174/400 ⇄ runner 288/400）、`:85`（"guard=0 在本机就能看见、⛔ 环境差"）⇒ 判据本身是胜率。
- ★**先量后选**：三个〔乙〕形候选（"计数器/钩子"、"绕开 C26 前置拦截"、"定点安装守卫再触发"）**本腿都在仓外量过、全部落空**，故不选：
  - **deny-all share（把手）**：`FILE_READ_ATTRIBUTES` **不受 share 门控** ⇒ guard 返回 nil（探针 `zz369r1_probe_test.go:TestZZProbeDenyAllShare`）；同一发里 `PrivateFile` 反而把调用者的文件截成 `""` 再报 sharing violation。
  - **DACL 拒 `FILE_READ_ATTRIBUTES`**：ACE 真写进去了（回读 `ace 0: type=1 mask=0x80`）但 `CreateFile(FILE_READ_ATTRIBUTES)` **仍然成功**（`zz369r1_probe3_test.go`）⇒ Windows 不按 DACL 判这一位。
  - **尾分隔符拼写**：guard 直调**确实**拒（"could not read the name"），但 `PrivateFile` 走的是**解析器的答案**（尾分隔符被清掉）⇒ 返回 nil，路由根本看不见。
  ⇒ ★**结论：`preOpenLeafGuard` 的"链接形"在 `PrivateFile` 上确定性不可达**（被两条稳态腿挡在前面），凡"它有没有答"必是胜率。
- 故取〔甲+乙〕混合形（也是本仓既有先例：`internal/observe/sealcalls_132_test.go` 的接线门、`ancestorwindow_342r2_windows_test.go:3-27` 的"重装世界一次驱动"）：**接线用调用图量、leg 用直驱量、route 用同步进入的状态量**。

## 4. 改了什么（逐行）

**`internal/winsec/openleaf_guard_322_windows_test.go`**（`168+/91-`）
- 删 400 发 racer 与其频率门（原 `:240 guardSpoke == 0 || flips == 0`）＋`sync/atomic`/`fmt` 等随之不用。
- 头注释新增"票 369 为什么换判据"段（`:33-56`，机制引到行）。
- 新常量 `:75-80`（`winLeafShots322r1 = 3` 三发静默重复）、`:85-92`（五行 roster，逐字抄产码文本）、`leakRow322r1` `:97`（**白名单**归因，未知名不算"某腿答了"）。
- `TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1` `:234` **改写**：每发**同步**把 bare-sibling 链 rename 到调用者名上并 `reparseAt` 证明（`:252-260`）⇒ 一发 `PrivateFile`，判 (a) **必须非 nil**（`:262`）、(b) 归因必须是**两条稳态行之一**（`:300`）、(c) 见证名未被移除、payload 未落进见证；(d) **正控**（`:315-322`）：把链删掉、写回调用者自己的对象 ⇒ 同一次调用**必须落地**。split/截断残留只报不判。
- `TestAC3WindowsLeafGuardRefusesLinkNamesAndAcceptsRealOnes322r1` `:152` **一字未动**（直驱：file-symlink／dir-symlink／junction 必须拒，真文件/未建名必须过）。

**`internal/winsec/routeguard_369_test.go`**（**新件**，无 build tag ⇒ linux 腿也扫）
- `routeGate369` `:31`＝`[resolveString preOpenLeafGuard pinnedParent openPinnedLeaf]`；`routeCalls369` `:46` 用 `go/parser` 取 `privateFile` 体里的调用与位置；`atTopLevel369` `:88` 要求 guard 那一跳**在函数体顶层**（夹在循环/闭包/第二个 `if` 里 ⇒ 红）；`routeGateOrder369` `:120` 判顺序；门 `:163`。

**`internal/winsec/ancestorunlink_348_windows_test.go`**（`88+/17-`）
- 叶子用例 `:117` 重写为三段：**(A)** `:118-171` 直驱 `openPinnedLeaf`（bare-sibling 链，312 那形），要求拒且**逐字**带 re-check 的链接形理由 `:151`＋两 sentinel；**(B)** `:173-188` **同步** `toFlip()`+`reparseAt` 证明世界能进翻态（替代 `provedReparse==0`）；**(C)** `:190` 保留 400 发世界，**只判 `deleted>0`** `:195`，split 降为读数 `:199-209`。
- 头注释 CASE 2 段改写（`// 票 369 换掉了 (b)`）。

**`internal/winsec/leafrace_342r2_windows_test.go`**（`26+/1-`）
- 首行 `//go:build windows` → **`//go:build windows && wisp342probe`**＋26 行头注说明"为什么出 CI 门、怎么跑、谁接替它的claim"。

## 5. 读数一：同一发测试连跑 ≥10 次（逐次红绿）

命令：`go test -count=10 -v -run 'TestAC2WindowsLeafIdentityReCheck…|TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking…|TestAC3WindowsTheRouteCallsTheLeafGuardBeforeThePin…|TestAC3WindowsLeafGuardRefusesLinkNames…' ./internal/winsec/`
⇒ **`RUN=40 PASS=40 FAIL=0 SKIP=0`**，四枚 × 10 发逐发全 PASS。**两发都贴**：`20-count10-after.txt`（12:33，变异前）与 `95-restore-check-count10.txt`（12:52，**最后一发变异之后**）——两发同读数，后者兼作 §9 的还原核。
逐次 348 枚 split 读数 `0,1,0,0,1,2,0,0,2,2`、`truncated` 每发 2（残留）——**读数在动、判据不动**，这正是要的形状。
route 枚逐发 `shot 0/1/2 = "path traverses a reparse point"`（三次重复同一行）；门逐发 `route order in winsec.go: [resolveString preOpenLeafGuard pinnedParent openPinnedLeaf]`。

## 6. 读数二：坏腿变异 ⇒ 改成确定性之后仍然必须红

变异全在仓外（`D:\tmp\…\mut\**`）＋`-overlay`，**无工作树突变窗口**：
| 变异 | 内容 | 结果 |
|---|---|---|
| `mut-guard` | `winsec_windows.go` 的 `preOpenLeafGuard` 体改成 `return nil` | **红**：`AC#3 RED file-symlink/directory-symlink/junction: preOpenLeafGuard accepted a name carrying a reparse point…`（`:188/:190/:192`）⇒ `--- FAIL` |
| `mut-recheck` | `pinnedleaf_342_windows.go` 的 `openPinnedLeaf` 去掉身份 re-check | **红**：`AC#2 RED: openPinnedLeaf returned a usable handle for a name that carries a link…`（`:153`）⇒ `--- FAIL` |
| 接线门三枚已知坏件 | `winsec.go` 删调用／塞进循环／移到 pin 后 | 检查器分别读成 `preOpenLeafGuard(MISSING)`／`(nested)`／`sameOrder=false`（`:31-door-teeth.txt`），真文件读 `sameOrder=true` |

★**诚实边界（量到就写）**：接线门**运行时读盘上 `winsec.go`**，所以用 `-overlay` 换 `winsec.go` 它**看不见**（三发 MUT-C/N/O 全 PASS，已实测）；它的牙是对"已知坏件"演示的（上表第三行），不是对 overlay 演示的。

## 7. 读数三：正控（守卫腿好 ⇒ 绿）

- 三发静默重复的 `PrivateFile` 在链删掉后**必须落地**（`:315`）⇒ 10 发全绿；`TestAC3WindowsLeafGuardRefusesLinkNames…` 里"真文件/未建名必须过"两跳 10 发全绿。
- 世界侧：`(B)` 的 `reparseAt(...)=true` 每发都印（10/10）。

## 8. `342` 那条探针：移了没有／移去哪／移出后门能不能绿

- **移了。** 移的是 `TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2`（`leafrace_342r2_windows_test.go:323`，红句 `:369-375` 的 `recheck == 0`）——★**它也是胜率门**（本机 41/400、45/400、0/400 三种读数都出现过）。
- **移到哪**：`.scratch/wisp/probes/342/**` 在**本腿写面之外**，故用**写面内的机制**＝build tag：`//go:build windows && wisp342probe`。文件仍在原处、不再进默认构建，手跑：
  `go test -tags wisp342probe -v -count=1 -run '^TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2$' ./internal/winsec/` ⇒ 本发 **PASS**，`re-check answered 45 of 400 (guard=0, proved flips=65)`。
- 机械读数：默认构建 `go test -list 342r2` **不再列出**它，同文件族的另三枚（`ancestorwindow_342r2`）仍在；带 tag 则列出。
- ★**移出后 CI 门能不能绿：不能，而且红的是另一枚、是真红。** `bash scripts/winsec-tests.sh` ⇒ **rc=1**、`RUN=127 PASS=75 FAIL=1`，红名 `TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1/symlink-dir`，红句 `ancestoropen_342_windows_test.go:50` 的 `truncated>0`（**2 of 400 真把外国对象的 body 换了**）。**这不是仪器假红**（假红＝"什么都没量到"），是票 342 未落地的**伤害分子** ⇒ 按硬提醒 #1「⛔ 把门改绿」**本腿没动它**。⇒ 口径：**"342 落地前门不可能稳定绿"仍成立，但原因已从"仪器假红"收窄成"真伤害读数"**。

## 9. 反向判据（三读数）＋ 还原核

1. **HEAD 原件、同一台机、`-count=10`** ⇒ `RUN=20 PASS=18 FAIL=2`：342 枚在 shot 1、8 **红**，其余绿 ⇒ **回到"有时红有时绿"**。
2. **只把 322 那枚 racer 的 pace 2ms→20ms（判据一个字节不改，`-overlay` 换测试件）** ⇒ `RUN=5 PASS=3 FAIL=2`：逐发 `guardRefusals=0,0,1,1,1`、`flips=14,12,15,14,11`（世界**真的**动了）⇒ ★**颜色是 pace 的函数，不是 leg 的函数** —— 这一发是"它真的在治胜率"的证据。
3. **接线门三枚已知坏件**（§6 第三行）⇒ 三种坏法各自红、真文件绿。
- **还原核（时机：最后一发变异之后）**：`git status --porcelain -- internal/winsec` 只有本腿 3 改 1 新、`git diff --numstat` 逐枚对得上（`88/17`、`26/1`、`168/91`、新件 `190` 行）；**12:52（全部变异之后）**在还原的树上重跑 §5 那发 10× ⇒ **`RUN=40 PASS=40 FAIL=0 SKIP=0`**（`95-restore-check-count10.txt`）⇒ 读数不是变异产物。

## 10. 四数（`go test -count=1 ./internal/winsec/`，⛔ 相减）

| 时机 | 四数 | rc | 备注 |
|---|---|---|---|
| 改前（本腿 12:12） | `RUN=127 PASS=76 FAIL=0 SKIP=0` | 0 | 绿发；**同一棵树在 `369-a1` 的 12:xx 发是 `PASS=75 FAIL=1`**（红名 `…BareSiblingForm342r2`）⇒ 门级也在掷硬币 |
| 改后（12:44） | `RUN=127 PASS=76 FAIL=0 SKIP=0` | 0 | `ok … 61.175s`；数目不变＝**−1**（342r2 探针出门）**+1**（新接线门） |
| 改后（12:47，CI 步原样跑） | `RUN=127 PASS=75 FAIL=1` | **1** | 红＝`ancestoropen_342_windows_test.go:50` 的 `truncated=2/400`（**真红**，非本腿造） |

## 11. 门禁 rc（逐条现跑）

- `go build ./...` ⇒ **rc=0**；`go vet ./internal/winsec/` ⇒ **rc=0**；`go vet -tags wisp342probe ./internal/winsec/` ⇒ **rc=0**。
- `sh scripts/d22scan.sh` ⇒ **rc=0**，逐字 `d22scan: clean - no D22 ban violations`（`ban #8 internal/ examined 643`）。
- `sh scripts/check-path-length-budget.sh` ⇒ **rc=0**，`VERDICT GREEN … roster equals the tree`（`61 paths`，两单位一致）。
- **格式甲乙两尺**：甲（入库 blob）＝`git show :<path> | gofumpt -d` 逐枚空、`rc=0`；乙（工作树）＝`gofumpt -l` 只列本腿 4 枚 ⇒ **空**（`gofumpt -l internal/winsec/` 输出 **0 行**；⛔ 用 `-l .`，它扫 `.scratch` 的故意坏件）。
- `git diff --cached --name-status`（提交前核索引）⇒ 只有本腿 4 条 `A/M`（⛔ 出现本笔不该有的 `D`/`A`）。

## 12. 未定义即停／具名上报

1. **写面里有别人的未提交件**：`docs/evidence/s1/376-a1-census.md`（在飞腿 `376-a1`，mtime 12:42:23）⇒ **一枚未碰、未 add**；本腿 `add`／`commit` 两处逐枚字面点名。
2. **"移出后 CI 门能绿"不成立**，红的换了人（§8）⇒ 这是本腿**量出来的**结论，⛔ 假设；处置权在编排者（票 342 落地／或另裁）。
3. **"两枚红其实是仪器错"部分成立**：两枚的**红因**（"什么都没量到"）确是仪器错——仪器用一个**时钟量**当判据；**守卫腿与 re-check 本身没坏**（坏腿变异下它们照样被抓）。故本票**改性质＝改判据**，⛔ 改性质成"产码洞"。
4. **接线门读盘、不吃 `-overlay`**（§6 诚实边界）⇒ 它抓"有人删/挪调用"，抓不到"overlay 里的产码"。
