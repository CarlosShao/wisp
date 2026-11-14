# 票 314 · 形状尺：de-share 319 真建窗那处 + 一枚让 A3 回退真红的源码形状尺

> 角色：写腿（实现者；审计者是另一条腿，本件不自审）。时间 2026-11-14 13:35–13:5x +08。
> 被写树：`dev` 工作树，HEAD `307becc69`。读数根 `.scratch/wisp/probes/314-ruler/`（日志逐字在此，入仓另存）。
> 起因：非实现者审计 `docs/evidence/s1/today-four-fixes-audit.md` 的两条硬读数——A3 那一形造不出红（`:305-331`），
> D③ 固定名未清零（`:373-379`、`:421-435`）。
> 边界：未碰产码、`scripts/**`、`internal/**`、`frontend/**`、`.github/**`、阈值/golden/契约；未跑整包 `./cmd/wisp`；
> 未 push；未建 worktree/checkout。写面 4 件＝两枚 `cmd/wisp` 测试件、`probes/314-ruler/**`、本件、票 314 Progress log 一行。

**第一行结论：两件都落盘。**

| # | 件 | 现跑读数 |
|---|---|---|
| 1 | `panel_death_replay_319_winlive_test.go` 的固定名夹改成 `panelProfileLeaf(...)` | 该文件固定名命中 **0**；`GOFLAGS=-tags=winlive go test -c` **rc=0**（63,072,718 B）；未跑窗（见 §1.3） |
| 2 | 新源码形状尺 `cmd/wisp/panel_profile_leaf_shape_windows_test.go` | 绿：**274 源文件 / 2 命中 / 2 具名豁免 / 0 红**，0.82 s；把默认档改回固定名：**rc=1**，红句逐字点名 `panel_unbound_253r6_windows_test.go:73` |

---

## 1 第一件：319 那处**真建窗**的固定名夹

### 1.1 改动（一处赋值 + 一处 import + 六行注释）

`cmd/wisp/panel_death_replay_319_winlive_test.go`（`//go:build windows && winlive`）：

```diff
-	dataPath := filepath.Join(os.TempDir(), "wisp-319r2-live-profile")
+	dataPath := panelProfileLeaf("wisp-319r2-live-profile")
```

`path/filepath` 在该文件里的**唯一**使用点就是这一行，所以 import 同步收窄（`os` 仍在 `:223-224`/`:357-358` 的 MutexName 用着，未动）。
helper 是既有件 `cmd/wisp/panel_host_windows_test.go` 的 `panelProfileLeaf`（`:70`），叶＝`prefix-<pid>-<nanotime>`，**零新 helper**。
`git diff --stat` = `1 file changed, 8 insertions(+), 2 deletions(-)`。

### 1.2 为什么这处必须改（审计已量，本件不重测）

审计 §4 D 逐字：这处是 `windows && winlive`，`startPanelHolder319`（改后 `:79`）在 `:128/:222/:341` 三处被用，**真建窗**。
同件 `:45-50` 的 exempt 只适用于「不建窗」的那一枚（`panel_holder_319_windows_test.go:53`）。真建窗而共用固定夹＝
审计 §4 测量到的「同一 UDF 同 tick 两个顶层 browser host，约 1.2 s 共存」那一形。

### 1.3 验证做到哪一步（诚实边界）

- **编译**：`GOFLAGS='-tags=winlive' go test -c -o .scratch/wisp/probes/314-ruler/wisp-winlive.test.exe ./cmd/wisp` **rc=0**，
  二进制 63,072,718 B（`wisp-winlive.test.exe` 的 LastWriteTime 2026-11-14 13:46:09）。⇒ 改后的 winlive 件能编译、能进构建。
- **没有跑这枚用例**：它真建窗，而本刻**另有腿在用桌面**（派单明令「只许 `-test.run` 定向，另有腿在用桌面」）。
  我不拿别人的桌面窗口去换一条读数。语义上这行只换了夹具夹名，`panelProfileLeaf` 已被同包 8 处 + 本件 `:81` 采用，
  其行为（唯一叶＋清 >1h 同族旧夹）不是本件引入的。
- **没有重测**审计的 A/B、并发 coexist 读数：派单明令「⛔ 不要重测」。

---

## 2 第二件：源码形状尺（让 A3 那种回退真红）

### 2.1 为什么必须是**源码形状**尺

审计 A3 逐字：把默认档 `panel_unbound_253r6_windows_test.go:71` 改回固定名后，`go test -c` rc=0（62,272,456 B），
直跑 `Test253r6TheSilentDropBranchIsLoud` **PASS**；仓里没有任何断言读「叶名唯一性」，唯一控制是
`.scratch/wisp/probes/314-deshare/concurrent.ps1`，而 `scripts/**`、`.github/**` 对它 0 引用。
⇒ 行为尺看不见这一形（共享夹只在两发同时落地时才成缺陷，A/B 两侧都是 rc=0），**「叶是不是固定名」是源码事实**，
所以尺子读源码。这也正是票面裁定选的路（票 314 `:123` 逐字：「不走改 `ci.yml` 那条路，Go 测试天然进 CI」）。

### 2.2 尺子读什么（274 个源文件、两条规则）

`cmd/wisp/panel_profile_leaf_shape_windows_test.go`，`//go:build windows`，**不建窗、纯 go/parser、0.82 s**：

- 扫描集＝`cmd/wisp/*.go`（**273** 件，产码与测试都扫：字面量可以只挪一个文件）+ 具名的
  `../../scripts/spike/webview2-latency/main.go`（审计 §5.3 的 spike 默认名那 1 处）＝ **274**；
  空扫描集 **Fatal**（盲尺不许印绿），读不到/解析不了 **Fatal**（跳过不是豁免）。
- 规则 **NAME**：`Join` 里同时出现 `os.TempDir()`（或它的局部/包级别名）与字符串字面量叶，且叶名含 `profile`
  ⇒ 红（参数顺序、跨行、`os.TempDir() + "/lit"` 串接都算）。这就是 A3 那一形，**不看是否流向窗**，
  所以「helper 返回固定名」的 319-holder 那形也在射程内（它就是靠这条进豁免表的）。
- 规则 **SINK**：同一个 Join（**叶名任意**）的值一旦进了 WebView2 data-path 下沉点
  （`NewPanelManager`、`newPanelControl`、`DataPath` 字段赋值、结构体字面量的 `DataPath:` 键），
  直接传或经一次局部变量绑定 ⇒ 红。没有这条，一个避开 `profile` 字样的固定名会整条漏掉。
- 为什么用 AST 不用字符串搜：本包注释里 `os.TempDir()` 出现极多（`panel_host_windows_test.go:53-56`、
  `panel_models_320_test.go:510/565` 等），文本尺会拿散文定罪；而且 AST 给得出**行号**、能把叶判成「字面量实参」
  而不是「任何同名出现」。

### 2.3 豁免表（2 条，全部写在测试件内部，逐字理由）

| 文件 | 叶 | 理由（读那段注释抄下来的） | 机器复核 |
|---|---|---|---|
| `cmd/wisp/panel_holder_319_windows_test.go` | `wisp-319r1-holder-profile` | 本文件一次窗都不建（headless tier 从不向 bringUp 要窗口），所以那个夹只是一个名字、从来不是 WebView2 user-data 夹；本文件内两个夹具共用它是故意的。 | `quotes` 逐字含 `Nothing here creates a window (the headless half never asks bringUp), so`（`:46`）与 `shared here on purpose: two fixtures in one run would`（`:48`）；另加 `noWindow:true`：文件里一旦出现 `bringUp`/`newPanelControl`/`startPanelHolder319` 调用即红——审计担心的「将来给它一个 window」变成机器可判 |
| `scripts/spike/webview2-latency/main.go` | `wisp-spike-webview2-profile` | 手工逐发跑的延迟测量工具（P11/D32 的 cold/hot 基数）的 `-datapath` 默认值：共用同一个夹正是它要测的 warm profile，每次调用由人显式串行发起、且 `-datapath` 可覆盖；测试套件从不启动它。 | `quotes` 逐字含 `flag.String("datapath"` 与 `"WebView2 user data folder"`（`:307`）——覆盖开关或用途说明被删，豁免跟着死 |

三条自检（每次跑）：①豁免必须仍覆盖 ≥1 处活命中（过期即红，防「空白支票」）；②`quotes` 每条必须仍逐字在源文件里
（理由被删改，豁免跟着死，不是留一张没出处的通行证）；③`noWindow` 的声明按 AST 复核。
豁免身份＝「文件＋叶名」，不是行号（行会漂），所以改叶名不会搭旧豁免的便车。

### 2.4 正控（派单必须做）：把任一处 `panelProfileLeaf(...)` 改回固定名 ⇒ 尺子**真红**

突变：`cmd/wisp/panel_unbound_253r6_windows_test.go`（**审计 A3 同一处**）
`dataPath := panelProfileLeaf("wisp-253r6-panel-profile")` → `dataPath := filepath.Join(os.TempDir(), "wisp-253r6-panel-profile")`，
import 加回 `os` + `path/filepath`（两行 ⇒ 该行从 `:71` 漂到 `:73`）。**突变体先被 carrier 编译通过**（`go test -c` 在跑之前，
日志第一段是 `building test binary`；rc 是测试的，不是编译的）。

逐字红（`mutant-fixedname-red.log`，`bin_rc=1`，`--- FAIL: Test314FixedPanelProfileLeafShapeIsRed (0.58s)`）：

```
    panel_profile_leaf_shape_windows_test.go:187: cmd/wisp/panel_unbound_253r6_windows_test.go:73 命中固定名面板 profile 形状（叶 "wisp-253r6-panel-profile"）：
            dataPath := filepath.Join(os.TempDir(), "wisp-253r6-panel-profile")
        这个形状会让两发并发互踩同一个 WebView2 user-data 夹：一个夹就是一个浏览器 profile，两个测试进程在同一条夹上各起一套顶层 browser host —— 先到的 CreateEnvironment 争同一夹会回 S_OK+NULL environment，后到的 environment 会 JOIN 一个正在退场的浏览器进程（internal/ball/orb_windows.go:432-458 的两条守卫）；票 314 的测量把同一夹同 tick 两个顶层 host 抓成过约 1.2 秒的共存读数。
        ticket 314 的非实现者审计量过这形的牙是零：把默认档 panel_unbound_253r6_windows_test.go:71 改回固定名，go test -c rc=0、Test253r6 PASS —— 仓里没有断言读叶名唯一性，所以只有这把源码形状尺能让这种回退红。
        改法＝panelProfileLeaf("wisp-253r6-panel-profile")（cmd/wisp/panel_host_windows_test.go 的既有 helper，叶＝prefix-pid-纳秒）；确有理由共用同一个夹的，在 panelProfileLeafExemptions 里按「文件＋叶名」具名登记，并把源码里那句理由逐字抄进 quotes。
    panel_profile_leaf_shape_windows_test.go:192: 314 形状尺：扫了 274 个源文件，3 处固定名命中（2 处具名豁免、1 处红）
```

还原与同串读数：

| 项 | 读数 |
|---|---|
| 文件 | `cmd/wisp/panel_unbound_253r6_windows_test.go` |
| SHA256 突变前 | `C9A14528D6E64CE98E13EDEF5F1064585B975234DF22BD384E78ADCAFB240551`（与审计 §0.1 的 `C9A14528D6E64CE9…` 同串） |
| SHA256 还原后 | `C9A14528D6E64CE98E13EDEF5F1064585B975234DF22BD384E78ADCAFB240551`（**逐位同串**） |
| 还原方式 | 突变前整文件备份到 `probes/314-ruler/orig/…bak`，还原＝原字节拷回（不是手工反编辑） |
| `git status --porcelain -- <该文件>` | **空**；`git diff -- <该文件>` **空** ⇒ 树上零残留 |
| 还原后再跑 | `bin_rc=0`，`--- PASS` ×2，`ok … 1.00s`，274 源文件 / **2** 命中 / 2 豁免 / **0 红** |

### 2.5 尺子自带的常驻正控（14 例，不依赖任何一次手工突变）

`Test314PanelProfileShapeRulerRingsOnTheFixedNameShape` 把形状喂给同一个匹配器，机制本身每次都被验：

- 必红 9 例：审计那行原样；字面量在前根在后；跨行；局部别名 `root := os.TempDir()`；包级别名；`os.TempDir() + "/lit"` 串接；
  helper 返回的 profile 形（319-holder 那形）；非 profile 叶但直接进 `NewPanelManager`；非 profile 叶经局部绑定进下沉点。
- 必静 5 例：`fmt.Sprintf` 每发叶；`panelProfileLeaf(...)`；`t.TempDir()` 局部夹具；注释散文＋引号文本（同一例）；
  320 那形「非 profile 的固定临时文件」。共 9＋5＝14 例，全部 PASS（`0.00s`）。

### 2.6 射程外（诚实边界，写在测试件头注释里）

1. 把根藏进 helper 形参（`fixedProfile("lit")` 里写 `os.TempDir()`）：根与字面量从不同现一个表达式，NAME 与 SINK 都跨不过函数边界；
2. 两个半截字面量拼出的叶（`"wisp-" + "319r1-..."`），或用非 pid/纳秒的假占位；
3. 根不是 `os.TempDir()` 的（`os.UserCacheDir`、`%TEMP%` 经 `os.Getenv`、绝对路径字面量）——票面裁定的形状是 `os.TempDir()`。
把固定名挪进这三类是新决定，要在票上论证，不是这把尺默许的编辑。

---

## 3 门禁（本件最终态现跑）

| 门禁 | 读数 |
|---|---|
| `go vet -tags winlive ./cmd/wisp` | **rc=0** |
| `gofumpt -l cmd/wisp` | **空**（0 行；gofumpt v0.12.0 / go1.27.1） |
| `sh scripts/d22scan.sh` | **clean**（见 §3.1） |
| `sh scripts/check-path-length-budget.sh` | **GREEN**（tracked=19,856；over-budget=69 全在名册；not in roster=0；最长 185；wall interval 内 0） |

### 3.1 d22scan 读数

`sh scripts/d22scan.sh` **rc=0**（日志 `.scratch/wisp/probes/314-ruler/d22scan.log`；正控 `runtests.sh` 先跑、再真扫描）：

```
d22scan: examined 481 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: scope ban #8 cmd/              examined 280 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=404, bans #1-5 cmd/=77, ban #6 frontend/=123, ban #7 internal/tools/=51, ban #10 internal/risk/=21, ban #8 design/=39, ban #8 frontend/=123, ban #8 internal/=1010, ban #8 cmd/=280
```

（`ban #8 cmd/=280` 把新增的这枚形状尺算进分母；它零 emoji，注释里也没有被 `U+2200–U+22FF` / `U+2600–U+27BF` 这些带扫的范围命中。）

其他读数：`go test -c`（默认档，carrier 内）rc=0；`GOFLAGS=-tags=winlive go test -c` rc=0（63,072,718 B）。

---

## 4 复跑命令（逐字）

```sh
# 形状尺（默认档，不建窗；走 cgo+DLL 的 staged carrier，且只定向跑）
sh .scratch/wisp/probes/314-ruler/run-ruler.sh -test.run '^Test314'

# winlive 件编译（只编译，不跑窗）
GOFLAGS='-tags=winlive' go test -c -o .scratch/wisp/probes/314-ruler/wisp-winlive.test.exe ./cmd/wisp
```

`run-ruler.sh` 只做一件事：`cd` 到仓根后 `sh scripts/wisp-cli-staged-tests.sh "$@"`（carrier 在 `cmd/wisp` 下直跑二进制，
`-test.run` 之后不会碰任何别的用例）。**没有一条命令跑整包。**

---

## 5 两处设计取舍（记下来，免得被当成放宽）

1. **初版规则更宽，被一个假阳性当场否掉**：第一版把「`os.TempDir()` ＋ 任何字面量」全判红，第一次跑就命中
   `cmd/wisp/panel_models_320_test.go:580` 的 `planted-key-320-temp.txt`——那是票 320 泄漏检测的**阳性控制文件**
   （`tempRoot := os.TempDir()` 在 `:517`，控制文件在 `:580` 附近种/读/删），不是 WebView2 profile 夹。
   ⇒ 处置不是「加一条豁免把表写宽」，而是把规则收紧到**票面裁定的主体**：叶名含 `profile`（NAME）或值真的流向
   data-path 下沉点（SINK）。320 那处因此**天然静默**，一条豁免都不欠。
2. **spike 为什么在扫描集里**：它是审计 §5.3 数出来的第 3 处固定名，且真的会开真 WebView2；把它具名进扫描集
   （而不是让它留在射程外），那处就成了豁免表里**唯一一条带 `-datapath` 覆盖证据的非测试豁免**；同时缺失即 Fatal，
   逼着「spike 搬家」与「豁免表」同一次提交。

## 6 一处数字更正（不重写历史）

第一枚提交 `0b28c7d98` 的 message 与本件初稿把 `cmd/wisp/*.go` 的件数写成 **272**，实测是 **273**
（`(Get-ChildItem cmd/wisp/*.go).Count` ＝ 273；加具名的 spike 1 件 ＝ 尺子日志的 **274**）。
274 是尺子自己打印的读数（`扫了 274 个源文件`），272 只是我手抄的加数，抄错一件。按仓规不改写已提交历史，
本件与票 314 Progress log 已更正为 273，此段留档。

---

## 7 审计后补洞（写腿 `w314-ruler2`，2026-11-14 14:1x–14:5x +08）

> 角色：**写腿（实现者）**，⛔ 审计者，不自审。上游＝`docs/evidence/s1/314-ruler-audit.md`（非实现者审计，D1–D5）。
> 写面 4 处：`cmd/wisp/panel_profile_leaf_shape_windows_test.go`、本件（只追加本节）、`.scratch/wisp/probes/314-ruler2/**`、
> 票 314 Progress log（只追加一行）。⛔ 未碰产码/其他测试/`scripts/**`/`internal/**`/`.github/**`/规格/阈值/golden；⛔ 未 push；
> ⛔ 未跑整包 `./cmd/wisp`（只 `-test.run '^Test314'`）；⛔ 未加宽豁免表（原样 2 条）。
> 仓外读数树：`C:\Temp\wisp-ruler2-20261114`（`old`＝HEAD 的尺、`new`＝补洞后的尺＋同一枚探针、`teeth`＝只放补洞后的尺）。
> 尺子 SHA256：补洞前 `7D0CB395303D09212DE920A9E24F3658F06E93FD48371B373A5A86C03F89A6E7`（＝审计件钉的那一枚；导出树 `old` 现跑**同串**）、
> 补洞后 `6E54A0B34D6A3CD79888F5CD33535EB0754F4549DDD3D2FE1FD8FFE3CAA723F2`。跟踪件里只有这一枚被改（`git diff --stat` ＝
> `1 file changed, 849 insertions(+), 115 deletions(-)`）。

**第一行结论：D1 修好／D2 修好／D3 修好／D4 修好／D5 修好。审计那 8 条共 10 形：8 形修好，2 形（A1b/A1c）仍漏，已在测试件头注释具名并写明为什么。**

### 7.1 五条洞：机制与守它的东西

| 洞 | 机制（补洞后） | 守它的常驻正控（必红/必静） | 仓外可编译突变体 |
|---|---|---|---|
| **D1** 多段拼接全静默 | `panelProfileFold`（`:623`）：`+` 链、`strings.Join` 字面量切片、`Join(...)` 再贴字面尾巴都折叠成常量，叶的判定在**折叠后的值**上；只有**整段**常量才算固定名（`opaque` 不算），所以「常量前缀＋pid」仍静默 | 5 条必红＋2 条必静 | 探针 `:36/40/44/48/52` 五形**全红** |
| **D2** 别的根全静默 | `panelProfileIsTempRoot`（`:830`）＝`os.TempDir()`／`os.Getenv("TEMP"\|"TMP")`／`os.UserHomeDir()`／别名／括号／**盘符绝对字面量**（`panelProfileDriveAbsolute` `:896`）；`home, err := os.UserHomeDir()` 这种二值赋值也把首名登记成根（`panelProfileTempRoots` `:790` 的 `record`） | 5 条必红＋3 条必静 | 探针 `:57/61/69/73/77` 五形**全红** |
| **D3** 包级作用域漏 | 包级声明**同时**走 NAME 与 SINK（`panelProfileHitsInSource` `:401`）；包级固定 `Join` 变量收集成表（`panelProfilePackageFixedJoins` `:537`）再 seed 进每个函数体，故「包级变量→函数里交给下沉点」也红 | 2 条必红＋1 条必静 | 探针 `:84/86` 两形**全红**；「包级 join 从不到下沉点」保持静默 |
| **D4** 三条臂零常驻正控 | 臂本身没改，缺的是**控制**：`newPanelControl`、`DataPath` 字段赋值、函数体内 `DataPath:` 键各加 1 条常驻必红控制（见 7.4） | 3 条必红 | 反证在 7.5：把臂打断，对应控制**当场红** |
| **D5** 汇总行口径误导 | 汇总口径抽成纯函数 `panelProfileSummaryLine`（`:247`）：印「形状命中 N／其中已豁免 M／实红 K」＋「豁免自检红 F」，并在 `K>0 或 F>0` 时写 **【红】**、绿时写 **【绿】**；`panelProfileCheckExemptions`（`:287`）改为返回自检红计数，不再只 `Errorf` | 1 条常驻控制（三条读数：clean／1 条未豁免命中／只有豁免自检红） | 反证在 7.5 T4 |

### 7.2 审计那 8 条，逐条结论

| 审计条目 | 形 | 结论 | 证据 |
|---|---|---|---|
| A2a | `Join(root, "wisp-33r1"+"-panel-profile")` | **修好** | 常驻 `two_literal_halves_inside_one_Join`；探针 `:36` 红 |
| A2b | 整条 `os.TempDir()+"/wisp-33r1"+"-panel-profile"` | **修好** | 常驻 `the_whole_path_by_concatenation`；探针 `:40` 红 |
| A2c | `strings.Join([]string{...},"")` | **修好** | 常驻 `the_leaf_by_strings.Join`；探针 `:44` 红 |
| A2d | `Join(...)` 外再贴尾巴 | **修好** | 常驻 `a_literal_suffix_outside_the_Join`；探针 `:48` 红 |
| A2e | 拼完先绑局部再传 | **修好** | 常驻 `a_concatenated_leaf_bound_to_a_local,_then_to_the_sink`；探针 `:52` 红 |
| A3a | 根 `os.Getenv("TEMP")` | **修好**（TEMP 与 TMP 都算） | 常驻 `root_os.Getenv(TEMP)`／`root_os.Getenv(TMP)`；探针 `:57/61` 红 |
| A3b | 根 `os.UserHomeDir()` | **修好**（含 `home, err :=` 二值赋值） | 常驻 `root_os.UserHomeDir_through_its_two_results`；探针 `:69` 红 |
| A3c | 根＝盘符绝对字面量 | **修好** | 常驻 `root_is_a_drive-letter_absolute_literal`／`a_drive-letter_absolute_literal_joined_with_a_leaf`；探针 `:73/77` 红 |
| X1 | 包级 `DataPath:` 键 | **修好** | 常驻 `a_package-level_DataPath_key_is_a_sink_where_it_is_declared`；探针 `:84` 红 |
| X2 | 包级固定 `Join` 变量交给下沉点 | **修好** | 常驻 `a_package-level_fixed_join_handed_to_a_sink_in_a_function`；探针 `:86` 红 |
| A1b | helper 返回非 profile 叶、值在下沉点 | **仍漏** | 函数边界：root 与字面量在扫描文件里**从不在同一表达式**；要跨过程分析，是另一件仪器 |
| A1c | 根在 helper 形参里（`mk(leaf)`） | **仍漏** | 同上：`Join(os.TempDir(), leaf)` 的 leaf 是形参，`mk("lit")` 在调用点，两者不在一处 |

仍漏的两形**具名写在测试件头注释**（`A1b … STILL OUT`／`A1c … STILL OUT`），并写明「要跨过程分析」。
另有两类**刻意不追**（也写在头注释）：四根之外的根（`os.UserCacheDir`、UNC 字面量、`%TEMP%` 经别的变量重拼）；
一半常量一半运行时的叶（`Join(root,"wisp-")+pid`：这**不是**共用夹缺陷，判红会冤枉 per-run 叶）。

### 7.3 仓外可编译突变体：同一枚探针，旧尺 0 红 / 新尺 12 红

探针源（只建不删）：`.scratch/wisp/probes/314-ruler2/zzz_ruler2_probe_windows_test.go`，`//go:build windows`，
**能编译**（`go test -c` 在两棵树里都 `rc=0`）。它把 A2a–A2e、A3a/A3b/A3c、X1/X2 写成 12 个可编译函数/声明，
外加 1 个反向控制（包级固定 join 从不到下沉点）。

| 树 | 尺 | 扫了 | 命中 | 豁免 | 实红 | `go test -c` | 跑 |
|---|---|---|---|---|---|---|---|
| `old` | HEAD 那一枚（SHA `7D0CB395…`） | 275 | 2 | 2 | **0** | rc=0 | **PASS**（全部 12 形静默） |
| `new` | 补洞后（SHA `6E54A0B3…`） | 275 | 14 | 2 | **12** | rc=0 | **FAIL**（`--- FAIL: Test314FixedPanelProfileLeafShapeIsRed`，0.70s） |

新尺的 12 条红句逐字点名到探针的 `:36/:40/:44/:48/:52/:57/:61/:69/:73/:77/:84/:86`，叶分别是
`wisp-audit2a-panel-profile`、`/wisp-audit2b-panel-profile`、`wisp-audit2c-panel-profile`、`wisp-audit2d-panel-profile`、
`wisp-audit2e-panel-profile`、`wisp-audit3a-panel-profile`、`wisp-audit3a2-panel-profile`、`wisp-audit3b-panel-profile`、
`C:\Temp\wisp-audit3c-panel-profile`、`C:\Temp\wisp-audit3c2-panel-profile`、`wisp-audit-pkg-sink`、`wisp-audit-pkg-var`
—— 折叠后的值，不是任何单个 token。汇总行同时印
`形状命中 14 处；其中已豁免 2 处、实红 12 处，豁免自检红 0 处`。
日志：`.scratch/wisp/probes/314-ruler2/mutants2.log`（旧尺在同文件前半段）。两次跑都走仓外导出树，⛔ 仓内零污染。

**还原**：仓内跟踪件从头到尾没被这两棵树碰过；`teeth` 树里改完再还原后 SHA 与仓内**同串**（见 7.5）。

### 7.4 常驻正控：14 例 → 36 例（＋1 条汇总行控制）

**既有 14 例逐条保留、逐条 PASS**（`test.v` 打出 14 个旧名 RUN/PASS），新增 **22 例**：

| 组 | 必红 | 必静 |
|---|---|---|
| D1 折叠 | `two_literal_halves_inside_one_Join`、`the_whole_path_by_concatenation`、`the_leaf_by_strings.Join`、`a_literal_suffix_outside_the_Join`、`a_concatenated_leaf_bound_to_a_local,_then_to_the_sink` | `a_folded_leaf_with_no_temp_root_stays_silent`、`a_constant_prefix_plus_a_runtime_suffix_is_not_a_fixed_name` |
| D2 根 | `root_os.Getenv(TEMP)`、`root_os.Getenv(TMP)`、`root_os.UserHomeDir_through_its_two_results`、`root_is_a_drive-letter_absolute_literal`、`a_drive-letter_absolute_literal_joined_with_a_leaf` | `a_Getenv_that_is_not_TEMP_or_TMP_is_not_a_root`、`a_cache_directory_is_not_one_of_the_roots_this_ruler_claims`、`a_relative_literal_root_is_not_a_root` |
| D3 包级 | `a_package-level_DataPath_key_is_a_sink_where_it_is_declared`、`a_package-level_fixed_join_handed_to_a_sink_in_a_function` | `a_package-level_fixed_join_that_never_reaches_a_sink_stays_silent` |
| D4 三条臂 | `the_newPanelControl_sink_arm`、`the_DataPath_field-assignment_sink_arm`、`the_DataPath_key_sink_arm_inside_a_function_body` | — |
| D5 汇总 | `the_summary_line_cannot_read_green_while_anything_is_red`（三条读数：clean／1 条未豁免命中／只有豁免自检红） | — |

每条必红控制都钉两队数字：叶名表（折叠后的值）＋**红句点名的行必须是该字面量所在行**（`lineNeedle`；折叠叶不是一行能写下的，
所以新增用例显式给 needle，而不是拿折叠值去搜源码）。每条必静控制都断言「零命中」。

### 7.5 新常驻控制有牙吗：四处「尺子自残」当场红

审计 §5.2 的原话是「把 `:409-411`（字段）或 `:431-436`（键）改坏了，14 例照样全绿」。补洞后，
在仓外 `teeth` 树里**逐条打断臂**（每次 `go test -c rc=0`），看对应控制是否当场红：

| 突变 | 打断的是 | 现跑 |
|---|---|---|
| T1 | 从 `panelProfileSinkCall` 里删掉 `newPanelControl` | **1 条 FAIL**：`.../the_newPanelControl_sink_arm` |
| T2 | 让「selector-LHS 是 DataPath」那条分支永不触发 | **1 条 FAIL**：`.../the_DataPath_field-assignment_sink_arm` |
| T3 | 让 `KeyValueExpr`（`DataPath:` 键）那条分支不做任何事 | **2 条 FAIL**：`.../the_DataPath_key_sink_arm_inside_a_function_body`＋`.../a_package-level_DataPath_key_is_a_sink_where_it_is_declared` |
| T4 | 汇总口径退回「只数未豁免形状命中」 | **1 条 FAIL**：`.../the_summary_line_cannot_read_green_while_anything_is_red` |

四次都 `build_rc=0`、`run_rc=1`；每次改完立刻从仓内拷回原件，**还原后** `build_rc=0 run_rc=0`、
`314 形状尺【绿】… 实红 0 处`、`RESTORED_SHA == REPO_SHA == 6E54A0B3…`。日志：`probes/314-ruler2/teeth.log`。

### 7.6 两处「不靠加豁免」的处置（记下来，免得被当放宽）

1. **盘符绝对字面量先被判红过一次，红的是尺子自己的控制表。** 第一版把「单独的绝对字面量」也交给 NAME 规则，
   第一次现跑就命中本件 `:1374/:1387` 两行 `want: []string{`C:\Temp\…-panel-profile`}` —— 那是**控制数据**，不是数据通路。
   处置：`panelProfileBareAbsolute`（`:583`）把「孤零零一个绝对字面量」从 NAME 规则里摘出来、只留给 SINK 规则
   （A3c 那一形正是交给 `NewPanelManager` 的），并在头注释写明「测试数据里一条只是拼着 `profile` 的 Windows 路径不算被告」。
   **没有加任何豁免**；这一形仍被 2 条必红控制钉住（裸字面量＝SINK 红、`Join(abs, leaf)`＝NAME 红）。
2. **一半常量的叶不算固定名。** `Join(root, "wisp-panel-"+strconv.Itoa(os.Getpid()))` 判静默（`opaque`），
   因为「常量前缀＋pid」正是 per-run 叶的写法，判红会制造假红。这是**收紧**（比旧版窄），已写进头注释「刻意不追」第 3 条。

### 7.7 门禁（逐条现跑）

| 门禁 | 读数 |
|---|---|
| `go vet -tags winlive ./cmd/wisp` | **rc=0**（19.7s） |
| `gofumpt -l cmd/wisp` | **空**（rc=0） |
| `sh scripts/d22scan.sh` | **clean**，`cmd/` 280 Go 文件全扫，ban #1-5/#8 零违规（rc=0，132.7s） |
| `sh scripts/check-path-length-budget.sh` | **VERDICT GREEN**（19930 条 tracked、over-budget 69 全在 roster、wall interval 0） |
| `Test314FixedPanelProfileLeafShapeIsRed` | **PASS**：274 源文件／命中 2／豁免 2／**实红 0**／自检红 0，0.72s |
| `Test314PanelProfileShapeRulerRingsOnTheFixedNameShape` | **PASS**：36 个形状子例＋1 个汇总子例全 PASS，0.01s |

日志：`probes/314-ruler2/{baseline2,final-green,d22scan,path-budget}.log`。

**复核（2026-11-14 15:1x +08，免拿旧绿当现绿）**：本节读数跑完后，另两条腿先后落了 `a7f0e8d05`／`eecf7afe3`，
所以两条树级门禁在 `HEAD=eecf7afe3` 上**又跑了一遍**：`d22scan` **clean／rc=0**、path-length **VERDICT GREEN／rc=0**
（`probes/314-ruler2/{d22scan,path-budget}-postcommit.log`）；本腿的 `Test314` 在 `25acce221` 上复跑仍是
**274 文件／2 命中／2 豁免／实红 0**、两测全 PASS（`probes/314-ruler2/post-commit-green.log`）。

### 7.8 诚实边界（本节不吹的地方）

- **A1b/A1c 仍漏**，理由是函数边界（7.2 末两行）；要修得引入跨过程分析，那是另一件仪器，不在本票射程。
- **扫描集仍是 274 件**：`cmd/wisp/*.go` 273＋具名 spike 1；探针文件只在仓外树里，**没有进仓**（避免把探针函数永久留在包内）。
- **耗时**：本机四读 0.54/0.66/0.70/0.72s（形状尺本体），量级与补洞前一致。
- **未跑整包 `./cmd/wisp`**；未 push；未改任何阈值/golden/规格。
- 本件是写腿的自证，**不是**非实现者审计：按 `SPEC-12 §4.3`，补洞后的对抗验收要另派一条腿。

