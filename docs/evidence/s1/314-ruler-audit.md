# 票 314 · 非实现者审计：源码形状尺（`panel_profile_leaf_shape_windows_test.go`）

> 角色：**审计腿（⛔ 实现者）**。时间 2026-11-14 13:5x–14:2x +08。
> 审计对象：`cmd/wisp/panel_profile_leaf_shape_windows_test.go`，SHA256 `7D0CB395303D09212DE920A9E24F3658F06E93FD48371B373A5A86C03F89A6E7`，
> blob `cb308abde6e2ebf5844cc3fc34cb839d28d28117`；`git log -- <该件>` 只有一枚提交 `0b28c7d98`，此后零改动
> （审计期间仓 HEAD 从 `50285dff1` 前进到 `069f97cac`，是别的腿的票 245，未动这把尺）。
> 树：`git archive HEAD`（`50285dff1`）导出的 `C:\Temp\wisp-ruler-audit-20261114` —— 攻击样本全在这里，⛔ 仓内零污染。
> 读数根：`.scratch/wisp/probes/ruler-audit/`（13 份日志、2 枚探针源 + 2 份探针 C 变体、2 份原件备份）。
> 边界：仓内写面只有本件、票 314 末尾一节、`probes/ruler-audit/**`；⛔ 未跑整包（只 `-test.run '^Test314'`）；
> ⛔ 未 push；⛔ 未改 `cmd/wisp/**` 任何跟踪件（改只发生在导出树，还原后 SHA256 逐位同串，见 §4.4）。

**第一行结论：攻击 1 抓住／攻击 2 漏了／攻击 3 漏了／豁免成立（且过期即红）／14 例计数对得上（9 必红＋5 必静）。**

---

## 0 现跑基线（先钉读数）

导出树里按仓自己的方式起测试二进制（`go test -c` + 把 `third_party/sherpa-onnx` 三枚 DLL 放在二进制旁边，
cwd＝`cmd/wisp`；`go test ./cmd/wisp` 直跑会 `0xc0000135`，与本尺无关）：

```
314 形状尺：扫了 274 个源文件，2 处固定名命中（2 处具名豁免、0 处红）
--- PASS: Test314FixedPanelProfileLeafShapeIsRed
--- PASS: Test314PanelProfileShapeRulerRingsOnTheFixedNameShape
```

| 项 | 自述 | 现跑 | 判 |
|---|---|---|---|
| 扫描集 | 273 件 `cmd/wisp/*.go` ＋ 具名 spike 1 件 ＝ 274 | `(Get-ChildItem cmd/wisp/*.go).Count`＝**273**；尺子自印 **274** | 对得上 |
| 命中/豁免/红 | 2 / 2 / 0 | **2 / 2 / 0** | 对得上 |
| 耗时 | 约 0.82s | 本机四读 **1.24 / 0.99 / 0.71 / 1.07s** | 量级一致（另说，不算缺陷） |
| 常驻正控 | 14 例（9 必红＋5 必静） | `-test.v` 打出 **14 个 `=== RUN` 子例、14 个 `--- PASS`** | 见 §5 |

日志：`logs/baseline-ruler.log`、`logs/clean-recheck.log`、`logs/final-green-after-probes.log`。

---

## 1 攻击 1：固定名挪进 helper 形参／helper 返回值

### 1.1 造了什么（逐字，导出树 `cmd/wisp/zzz_ruler_audit_probe_a_windows_test.go`）

```go
// A1a - 票面攻击 1 的原样：固定叶在 helper 里拼，返回值交给 DataPath 下沉点
func auditMkProfile() string {
	return filepath.Join(os.TempDir(), "wisp-audit1a-panel-profile")
}

func auditAttack1a() {
	_ = NewPanelManager(nil, nil, auditMkProfile())
}

// A1b - 同一形状，但叶避开 "profile" 字样：只剩 SINK 规则可能看见
func auditMkPlain() string {
	return filepath.Join(os.TempDir(), "wisp-audit1b-quicklook")
}

func auditAttack1b() {
	_ = NewPanelManager(nil, nil, auditMkPlain())
}

// A1c - 头注释自己交底的射程外 #1：根在 helper 里，字面量在调用点
func auditMkParam(leaf string) string {
	return filepath.Join(os.TempDir(), leaf)
}

func auditAttack1c() {
	_ = NewPanelManager(nil, nil, auditMkParam("wisp-audit1c-panel-profile"))
}

// X1 - 头注释三条之外：包级结构体字面量的 DataPath 键 + 非 profile 叶
var auditPackageSink = auditSink{DataPath: filepath.Join(os.TempDir(), "wisp-auditpkg-sink")}
```

（另有 A2a–A2c、A3a–A3c 同文件，见 §2/§3。）

### 1.2 尺子红了还是静默

**A1a：红。** 加了这枚文件后扫描集 276 件（273 ＋ spike ＋ 2 枚探针），尺子交出的**完整**命中表是 2 条豁免 ＋ 1 条红：

```
panel_profile_leaf_shape_windows_test.go:187: cmd/wisp/zzz_ruler_audit_probe_a_windows_test.go:19 命中固定名面板 profile 形状（叶 "wisp-audit1a-panel-profile"）：
        return filepath.Join(os.TempDir(), "wisp-audit1a-panel-profile")
314 形状尺：扫了 276 个源文件，3 处固定名命中（2 处具名豁免、1 处红）
--- FAIL: Test314FixedPanelProfileLeafShapeIsRed (0.63s)     [rc=1]
```

日志：`logs/probe-a-scan-set.log`。

**A1b / A1c / X1：静默。** 尺子把每条命中都打印，`3 处命中` 里只有 A1a 与两条具名豁免——所以同一次运行就证明其余 9 形
（A1b、A1c、A2a–c、A3a–c、X1）**零命中**。逐形隔离读数在 §1.4 的表（探针 B，同一匹配器、单发 16 例）。

### 1.3 它凭什么（读匹配代码）

- **抓 A1a 的不是数据流，是「字面量自证」。** `panelProfileNameRule`（`:371-380`）对每个 `Join` 调
  `panelProfileFixedJoinIn`，命中条件在 `:375`：
  `if ok && strings.Contains(strings.ToLower(fixed.leaf), "profile")`。
  叶里只要有 `profile`，这个 Join 在**任何**上下文都算固定名——helper 返回值、包级常量、没人调用的函数，全一样。
  尺子从来不看 `NewPanelManager` 那个实参：`auditMkProfile()` 这条调用与命中毫无关系。
- **A1b 为什么静默：SINK 规则跨不过函数边界。** `panelProfileSinkRule`（`:386-440`）只在下沉点实参**本身**是固定 Join 时记命中
  （`:429` `hits(arg)` → `:390-392` → `panelProfileFixedJoinIn`），或者经一次**局部**绑定（`:405-407` 存 `stored[name]`，`:414-416` 存 `ValueSpec`，
  `:423-427` 再用）。`auditMkPlain()` 是 `*ast.CallExpr`：`:449` 判 callee 名字不是 `Join`，`panelProfileFixedJoinIn` 立刻 false；
  `:419` 判 `auditMkPlain` 不是下沉点。⇒ 返回值这条路，两条规则都不走。
- **A1c 为什么静默：根与字面量从不同居一个表达式。** helper 体内 `filepath.Join(os.TempDir(), leaf)` 的 `leaf` 是 `*ast.Ident`
  （`:455` 只认 `*ast.BasicLit` + `token.STRING`），`:466` `if rooted && fixed.leaf != ""` 因此不成立；调用点的字面量又不是 Join 的实参。
- **X1 为什么静默：SINK 只走函数体。** `panelProfileHitsInSource` 对包级 `GenDecl` 只调 NAME（`:344-348`），
  NAME＋SINK 一起走的是 `fn.Body`（`:353-365`）。所以「包级 `var` 里 `DataPath:` 键 + 非 profile 叶」两边都不看。
  X1 不在头注释的三条射程外里；它是**第四条**（同因的还有「包级固定 Join 变量，函数里再交给下沉点」，见 §1.4 的 X2）。

### 1.4 隔离读数（探针 B，`panelProfileHitsInText` 同一匹配器，单发 16 例）

`logs/probe-b-table-final.log`，逐行原文：

```
[HIT]    A1a helper return, profile leaf, value to sink        leaves=["wisp-audit1a-panel-profile"]
[SILENT] A1b helper return, plain leaf, value to sink          leaves=[]
[SILENT] A1c param helper (header escape 1)                    leaves=[]
[SILENT] A2a two half literals inside one Join                 leaves=[]
[SILENT] A2b whole path by + concatenation                     leaves=[]
[SILENT] A2c leaf by strings.Join                              leaves=[]
[SILENT] A2d join then suffix, outside the Join                leaves=[]
[SILENT] A2e concatenated value bound to a local, then to sink  leaves=[]
[SILENT] A3a root os.Getenv("TEMP")                            leaves=[]
[SILENT] A3b root home directory                               leaves=[]
[SILENT] A3c absolute literal root                             leaves=[]
[SILENT] X1 package-level DataPath key, plain leaf             leaves=[]
[HIT]    C1 control: inline join, plain leaf, to sink          leaves=["wisp-plain"]
[HIT]    C2 control: NAME alone, no sink                       leaves=["wisp-name-only-profile"]
[HIT]    C3 control: DataPath field assignment                 leaves=["wisp-field"]
[HIT]    C4 control: DataPath key inside a function body       leaves=["wisp-key"]
```

C1 是关键的对照：**同一个非 profile 叶，Join 直接写在实参里就红（C1），隔一层 helper 返回就静默（A1b）**。

**X2（包级固定 Join 变量 + 函数里交给下沉点）**单跑：探针 C 与探针 A 同法放进扫描集（该文件只需能 parse，不重编译），
`275 件 / 2 命中 / 2 豁免 / 0 红`、`rc=0`（`logs/probe-c-pkgjoin-plain.log`）；
只把叶改成含 `profile`（一个字符改动），同一棵树立刻 `275 / 3 / 2 / 1`、红句点名 `zzz_ruler_audit_probe_c_windows_test.go:18`
（`logs/probe-c-pkgjoin-profile.log`）。⇒ 包级作用域不是护栏，**叶的拼写才是**。

### 1.5 判读

攻击 1 按票面逐字（helper **返回值** → `DataPath`）**抓得住**，但抓它的是 NAME 的字面量自证，不是数据流；
文件名里没有 `profile` 时同一形状全线静默（A1b），头注释自己交底的「根藏进形参」（A1c）也确实是静默的。
⇒ 「NAME 做数据流吗」＝不做；「SINK 只看就近文本吗」＝比就近略强（认一次**局部**绑定），但**不认返回值、不认包级绑定**。

---

## 2 攻击 2：两个半截字面量拼叶

### 2.1 造了什么（逐字）

```go
// A2a - 同一个 Join 里，叶由两截字面量拼
func auditAttack2a() {
	_ = NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-33r1"+"-panel-profile"))
}

// A2b - 整条路径三段拼
func auditAttack2b() {
	_ = NewPanelManager(nil, nil, os.TempDir()+"/wisp-33r1"+"-panel-profile")
}

// A2c - strings.Join 拼叶
func auditAttack2c() {
	_ = NewPanelManager(nil, nil, filepath.Join(os.TempDir(),
		strings.Join([]string{"wisp-33r1", "-panel-profile"}, "")))
}
```

探针 B 另加两形（逐字在 `probes/zzz_ruler_audit_probe_b_windows_test.go`）：

```go
// A2d - Join 之后再拼尾巴（拼接在 Join 外面）
_ = NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-33r1")+"-panel-profile")

// A2e - 拼出来的值先绑局部，再交给下沉点
p := filepath.Join(os.TempDir(), "wisp-33r1") + "-panel-profile"
_ = NewPanelManager(nil, nil, p)
```

### 2.2 尺子红了还是静默

**五形全静默。** §1.2 那次跑（探针 A 在扫描集里）A2a/A2b/A2c 零命中：`3 处命中` 只有 A1a ＋ 两条豁免。
A2d/A2e 在探针 B 表里 `[SILENT]`。同一次运行的对照 C1 证明 SINK 规则确实在工作（它不是瞎了）。

### 2.3 它凭什么

- **叶必须是单枚 `*ast.BasicLit`**：`:455` `if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING`。
  `"wisp-33r1"+"-panel-profile"` 是 `*ast.BinaryExpr`，不满足 ⇒ `fixed.leaf` 留空 ⇒ `:466` 不成立。A2a 静默。
- **`+` 串接分支要求「一边是字面量、另一边是 temp root」**：`:469-478`，两处 `panelProfileIsTempRoot` 判断在 `:473`/`:476`。
  A2b 的 AST 是 `(os.TempDir()+"/wisp-33r1") + "-panel-profile"`：外层右边是字面量但左边是 BinaryExpr（不是 root）；
  内层虽然 root＋字面量成立，叶是 `"/wisp-33r1"`（不含 `profile`，NAME 放行），且不在下沉点实参位置。
- **A2d/A2e 同一根因**：`hits(arg)`（`:429`）只把**那一个实参表达式**交给 `panelProfileFixedJoinIn`，不递归子表达式；
  A2e 的绑定只在下沉点实参**是**固定 Join 时才存（`:401`/`:405`），而这里 RHS 是 BinaryExpr。
- **`strings.Join` 不特殊**：`panelProfileCalleeName`（`:578-588`）把 `strings.Join` 也叫 `Join`，但它的实参是
  `CompositeLit` ＋ 空字面量 `""`（`:456` 判空后 `:466` 仍不成立）。A2c 静默。

### 2.4 判读

头注释射程外 #2 **成立**：半截字面量（在 Join 里拼、在 Join 外加尾巴、先绑局部、`strings.Join`）全部静默。
这把尺只认**单枚字面量叶**。

---

## 3 攻击 3：非 `os.TempDir()` 的根

### 3.1 造了什么（逐字）

```go
// A3a - %TEMP% 经 os.Getenv
func auditAttack3a() {
	_ = NewPanelManager(nil, nil, filepath.Join(os.Getenv("TEMP"), "wisp-audit3a-panel-profile"))
}

// A3b - 家目录
func auditAttack3b() {
	home, _ := os.UserHomeDir()
	_ = NewPanelManager(nil, nil, filepath.Join(home, "wisp-audit3b-panel-profile"))
}

// A3c - 绝对路径字面量，根本没有根表达式
func auditAttack3c() {
	_ = NewPanelManager(nil, nil, `C:\Temp\wisp-audit3c-panel-profile`)
}
```

### 3.2 尺子红了还是静默

**三形全静默**：同样由 §1.2 的单次运行证明（命中表只有 A1a ＋ 2 条豁免），探针 B 的 A3a/A3b/A3c 也都是 `[SILENT]`。
注意 A3a/A3b 的叶**含 `profile`**——NAME 想抓也抓不动，因为缺的是根，不是叶。

### 3.3 它凭什么

`panelProfileIsTempRoot`（`:559-574`）只认三种拼法：`os.TempDir()` 调用（`:562-567`，selector 必须是 `os`.`TempDir`）、
已登记的别名 `*ast.Ident`（`:568-569`）、括号（`:570-571`）。`os.Getenv("TEMP")` 是 Callee `Getenv`、`home` 未登记、
裸字面量连 `Join` 都没有（`:445-481` 只处理 CallExpr/BinaryExpr）⇒ `rooted` 永远 false。
别名登记也只在「恰好赋值 `os.TempDir()`」时发生（`:525-555` 的 `record`，`:537` 用同一个 `IsTempRoot` 判 RHS）。

### 3.4 判读

头注释射程外 #3 **成立**，且是**最宽**的一条：换根一行代码，NAME 与 SINK 一起失效，连「叶里写着 profile」都救不回来。

---

## 4 复核：豁免表有没有被写宽

### 4.1 两条豁免逐条核

| # | 文件 | 叶 | 理由句逐字在源码里？ | 仍覆盖真命中？ |
|---|---|---|---|---|
| 1 | `cmd/wisp/panel_holder_319_windows_test.go` | `wisp-319r1-holder-profile` | **在**：`:46` `// at. Nothing here creates a window (the headless half never asks bringUp), so`、`:48` `// where it matters, and shared here on purpose: two fixtures in one run would` | **是**：命中 `:53`（尺子日志逐字 `豁免命中：cmd/wisp/panel_holder_319_windows_test.go:53`） |
| 2 | `scripts/spike/webview2-latency/main.go` | `wisp-spike-webview2-profile` | **在**：`:307` 同一行同时含 `flag.String("datapath"` 与 `"WebView2 user data folder"` | **是**：命中 `:307` |

现跑基线（`logs/baseline-ruler.log`）的两行 `豁免命中` 就是这两处，读数 `2 命中 = 2 豁免`：**表没有第三行、也没有一条空转**。

### 4.2 表为什么不可能「悄悄写宽」

- 身份是 `(file, leaf)` **精确**相等（`panelProfileLeafExemptionFor` `:258-265`），没有前缀/正则/glob 形态（表注释 `:126-129` 明说）；
  改一个字的叶名就不搭旧豁免（§4.3 的 M2 实测）。
- 每次跑三条自检（`:216-254`）：仍覆盖 ≥1 处（`:222-232`）、`quotes` 逐字仍在源文件里（`:238-244`，读 `src.src` 原文）、
  `noWindow` 的文件不得出现三个建窗调用名（`:245-252` → `panelProfileWindowCallIn` `:593-615`）。
- 顺带一条**边界**（不是缺陷，是形状）：`noWindow` 是**三个名字的直接调用**检查（`:104`）。holder 件现在调 `newPanelFixture319`
  （`:59`）——如果哪天那个 fixture 内部真建窗，这条自述仍然「真」，尺子看不见。这一条尺子没声称能看，记在这里免得被当成保证。

### 4.3 过期即红：四发突变（全在导出树，跑完还原）

| 突变（逐字） | 尺子反应 | 日志 |
|---|---|---|
| **M1** 删掉 holder `:46` 那句 `Nothing here creates a window (the headless half never asks bringUp), so` | **红**：`:190 豁免 cmd/wisp/panel_holder_319_windows_test.go 的理由句在源码里找不到了：\"Nothing here creates a window…\"`；`--- FAIL`、rc=1 | `logs/mutant4a-quote-deleted.log` |
| **M2** holder `:53` 叶改 `"wisp-319r1-holder-profile-renamed"` | **红两处**：新叶本身被当固定名命中（`:187`，点名 `:53`）＋`:190 豁免表过期：… 的固定名叶 \"wisp-319r1-holder-profile\" 在这次扫描里覆盖 0 处`；读数 `2 命中/1 豁免/1 红` | `logs/mutant4b-leaf-renamed.log` |
| **M3** spike `:307` 的 `flag.String("datapath"` 改名 `flag.String("data-path"` | **红**：`:190 豁免 scripts/spike/… 的理由句在源码里找不到了：\"flag.String(\\\"datapath\\\"\"` | `logs/mutant4c-spike-flag-renamed.log` |
| **M4**（加测）holder 末尾加 `func auditNoWindowProbe() { bringUp(0) }` | **红**：`:190 豁免 … 自述「这里不建窗」，但 cmd/wisp/panel_holder_319_windows_test.go:389 调用了 bringUp` | `logs/mutant4d-nowindow-false.log` |

一处读数口径要说清：M1/M3/M4 的**汇总行仍印 `0 处红`**（`logs/mutant4a…` 最后一行），因为 `:192` 那个计数器只数
「未被豁免的形状命中」；豁免自检的红是独立的 `t.Errorf`。⇒ 别拿汇总行的 `0 红` 当绿证，看 `--- FAIL`/rc。

### 4.4 还原与残留

两枚被改文件在突变前整文件备份到 `probes/ruler-audit/orig/`，还原后 SHA256 与备份**逐位同串**：
holder `D79647435FC74C73E7DD3484F9734199480616DCEC0895DC29A7F4CE74DF692A`、
spike `C85E952DBA92DD3D6C5B71EE4FD46717F0EE031517080A138DCD6972565C2648`；最后再跑 `274 / 2 / 2 / 0`、rc=0
（`logs/final-green-recheck.log`、`logs/final-green-after-probes.log`）。

⇒ **豁免表没有被写宽：2 条都与源码逐字对得上、各覆盖 1 处真命中、身份精确、过期即红（四种删改都实测）。**

---

## 5 复核：14 例合成正控是真的吗

### 5.1 尺子自己不打印正控计数

`Test314PanelProfileShapeRulerRingsOnTheFixedNameShape` 只在 `t.Run` 里断言，**没有**「正控 N 例」这类汇总行；
`-test.v` 能数出来的是 `=== RUN` 14 条、`--- PASS` 14 条（`logs/baseline-ruler.log`）。
所以按派单要求读代码逐条数（`:637-874`，用例表 `:643-841`）：

| # | 用例（源码 `:行`） | 期望 | 例 |
|---|---|---|---|
| 1 | 审计那行原样 `:643-657` | 1 命中 `wisp-253r6-panel-profile` | 必红 |
| 2 | 字面量在前、根在后 `:658-672` | 1 | 必红 |
| 3 | 跨行 `:673-690` | 1 | 必红 |
| 4 | 局部别名 `root := os.TempDir()` `:691-706` | 1 | 必红 |
| 5 | 包级别名 `:707-723` | 1 | 必红 |
| 6 | `os.TempDir() + "/lit"` 串接 `:724-735` | 1 | 必红 |
| 7 | helper 返回的 profile 叶（319-holder 形）`:736-750` | 1 | 必红 |
| 8 | 非 profile 叶直传 `NewPanelManager` `:751-765` | 1 | 必红 |
| 9 | 非 profile 叶经局部绑定进下沉点 `:766-781` | 1 | 必红 |
| 10 | 非 profile 固定临时文件（320 那形）`:782-795` | 0 | 必静 |
| 11 | `fmt.Sprintf` 每发叶 `:796-810` | 0 | 必静 |
| 12 | `panelProfileLeaf(...)` helper `:811-819` | 0 | 必静 |
| 13 | 测试局部 `d.TempDir()` `:820-834` | 0 | 必静 |
| 14 | 注释散文＋引号文本 `:835-841` | 0 | 必静 |

**9＋5＝14，与自述对上**；用例表的 9/5 分界正好落在 `:781`/`:782`。

### 5.2 一处覆盖不均（计数对，但 14 例没覆盖四条 SINK 臂的全部）

SINK 规则点名四种下沉形态（`:491-506`：`NewPanelManager`、`newPanelControl`、`DataPath` 字段、`DataPath:` 键），
14 例只覆盖了 **`NewPanelManager` 的两种传法**（#8 直传、#9 局部绑定）。剩下三臂**没有常驻正控**：
`newPanelControl` 调用、`DataPath` 字段赋值、函数体内 `DataPath:` 键。
我用探针 B 的 C3/C4 证明后两臂**当前是好的**：

```
[HIT] C3 control: DataPath field assignment            leaves=["wisp-field"]
[HIT] C4 control: DataPath key inside a function body  leaves=["wisp-key"]
```

但它们只是我这一发的临时探针，不在常驻 14 例里 ⇒ 若哪天有人改坏 `:409-411`（字段赋值）或 `:431-436`（键），
14 例会全绿而尺子已瞎一半。这是「正控计数对得上」与「正控覆盖够不够」两件事，前者是，后者差三臂。

### 5.3 判读

**计数对得上**：14 = 9 必红 ＋ 5 必静，现跑全 PASS，且每条都真的在断言期望叶集合（`:854-856`）与行号（`:857-862`）。
唯一要记的是 §5.2 的覆盖缺口。

---

## 6 这枚尺能翻 314 的 AC#4 吗

**不能单独翻。** 它买到的是「A3 那种回退会红」这一条子条件，不是 AC#4 要的处置闭环：

1. **射程是真的窄，而且是它自己交底的窄。** 攻击 2、3 五形全静默（§2、§3），攻击 1 的「根进形参」静默（A1c），
   另有两条它没交底的静默（§1.4：helper 返回非 profile 叶、包级绑定）。换根一行 = 全盲；拆字面量一行 = 全盲。
   ⇒ 它保证的是「这几种写法不再偷偷回来」，不是「仓里没有固定名 profile 夹」。
2. **残余的两处固定名是「登记」不是「清零」。** holder `:53` 与 spike `:307` 都在豁免表里活着；AC#4 的「修（谁）／登记（多长）」
   里，「登记」这一半现在是机器可判的（§4 四发突变），「清零」这一半没有发生，票 314 里也没有授权去改 holder 的无窗共享
   （那是审计 `today-four-fixes-audit` 拍的保留）。
3. **AC#2（静桌面读数）仍缺，与这把尺无关**：尺子不碰桌面、也不产出那条读数。
4. **它确实进了 CI**：`ci.yml:1165` 跑 `bash scripts/wisp-cli-tests.sh` → `portable-tests.sh --scope=cli` →
   `wisp-cli-staged-tests.sh`；那条链只用 `-test.skip` 跳语音用例（`wisp-cli-staged-tests.sh:455`/`:471`），
   **没有 `-test.run` 过滤**，`Test314*` 不匹配任何 skip 模式 ⇒ 它在 windows 腿的分母里。这一条 AC#4 的门禁前提成立。

**结论句：这枚尺能翻 314 的 AC#4 吗——不能，它只把「A3 回退无牙」这一条从纸面变成默认套件里的确定性红；
AC#4 要的「残余固定名处置」与「AC#2 静桌面读数」仍在尺子射程之外。**

---

## 7 复跑（逐字）与读数清单

```sh
# 1) 导出树（审计时 HEAD 50285dff1）
git archive HEAD -o C:/Temp/wisp-ruler-audit-20261114/src.zip && \
  powershell -c "Expand-Archive -Path C:/Temp/wisp-ruler-audit-20261114/src.zip -DestinationPath C:/Temp/wisp-ruler-audit-20261114 -Force"

# 2) 起二进制（本机 go test ./cmd/wisp 直跑会 0xc0000135：DLL 必须挨着 exe）
cd C:/Temp/wisp-ruler-audit-20261114 && go test -c -o stage/wisp.test.exe ./cmd/wisp
cp "D:/work/workspace/projects plans/Wisp/third_party/sherpa-onnx/"*.dll stage/   # 导出树没有 third_party（git-ignored），从仓里拷三枚
cd C:/Temp/wisp-ruler-audit-20261114/cmd/wisp && ../stage/wisp.test.exe -test.run '^Test314' -test.v
```

探针源（只建不删，逐字）：`probes/ruler-audit/probes/` 下 4 件（A、B、C 两变体）。
日志 13 份：`probes/ruler-audit/logs/`。原件备份 2 份：`probes/ruler-audit/orig/`。
