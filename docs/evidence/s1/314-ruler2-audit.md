# 票 314 · 非实现者审计（第二轮）：补洞后的源码形状尺

> 角色：**审计腿（⛔ 实现者）**。时间 2026-11-14 15:0x–15:1x +08。
> 审计对象：`cmd/wisp/panel_profile_leaf_shape_windows_test.go`，**仓 HEAD `a46352fd3` 里的那一枚**
> （blob `627c2fcdb2c5901b08aec224a2bea1050af3491b`，SHA256 `6E54A0B34D6A3CD79888F5CD33535EB0754F4549DDD3D2FE1FD8FFE3CAA723F2`），
> 即修复腿 **`c5be2d703`** 交付、`docs/evidence/s1/314-shape-ruler.md` §7 自证的那一枚。
> 树：`git archive HEAD` 导出的 **`C:\Temp\wisp-ruler2-audit-20261114\tree`**；测试二进制 `go test -c` ＋ 三枚 sherpa DLL 同置，
> cwd＝导出树的 `cmd/wisp`，只跑 `-test.run '^Test314'`。探针源与日志在
> `.scratch/wisp/probes/ruler2-audit/{probes,logs,orig}/`（只建不删）。
> 边界：仓内写面只有本件、票 314 末尾一节、`probes/ruler2-audit/**`；⛔ 未改 `cmd/wisp/**` 任何跟踪件、⛔ 未跑整包、⛔ 未 push。
> 读数根日志：`baseline.log`／`a1-a2-scan.log`／`a3-scan.log`／`a3-contrast.log`／`a4-teeth.log`／`a5-k.log`／`a5-f.log`／
> `a6-cases.txt`／`final-green.log`。

**第一行结论：a46352fd3 靶向 c5be2d703 的尺 ｜ 1 漏 ｜ 2 漏 ｜ 3 漏 ｜ 4 抓 ｜ 5 抓 ｜ 6 抓。**

---

## 0 靶怎么钉的（先说一件审计期间发生的事）

### 0.1 工作树里那枚尺**正在被别人改**，本审计只认已提交的那一枚

审计开始时仓 HEAD＝`a46352fd3`，工作树 `cmd/wisp/panel_profile_leaf_shape_windows_test.go` 与 HEAD 一致（SHA `6E54A0B3…`，1630 行）。
审计中途（约 15:05–15:12）该文件出现**未提交修改**，行数从 1630 涨到 2362，内容自述是在补 **A1b/A1c（票 442）**
（`git diff` 878 insertions / 146 deletions；期间 SHA 从 `6E54A0B3…` → `CAD3C5D5…` → `9F97653…`，还在变）。
⇒ 本件所有读数都来自 **`git archive HEAD` 导出的树**，`file:line` 全部指 HEAD 的 1630 行版本（blob `627c2fcd`）；
工作树那枚在飞的 442 实现**不**是本件的审计对象，也**没有**被本件采信或否定。这条单列，免得把「另一个腿的半成品」当成「c5be2d703 修的」。

### 0.2 基线（现跑，不是抄）

```
314 形状尺【绿】：扫了 274 个源文件，形状命中 2 处；其中已豁免 2 处、实红 0 处，豁免自检红 0 处
--- PASS: Test314FixedPanelProfileLeafShapeIsRed (0.96s)
--- PASS: Test314PanelProfileShapeRulerRingsOnTheFixedNameShape (0.01s)
FINAL_RC=0
```

`-test.v` 数出 **36 个 `=== RUN` 子例**：**35 个形状子例 ＋ 1 个汇总行子例**（`baseline.log` / `final-green.log`）。
本件用默认档构建（尺的扫描集是 `cmd/wisp/*.go` 的 glob ＋ 具名 spike，与 build tag 无关）；与派单给的 `GOFLAGS=-tags:winlive` 读数同数（274／命中 2／已豁免 2／实红 0／自检红 0）。
⚠ 一处**文件对不上**：`314-shape-ruler.md` §7.4 写「14 例 → 36 例（＋1 条汇总行控制）」、§7.7 写「36 个形状子例＋1 个汇总子例」，
实际是 **35 形状＋1 汇总＝36 个子例**。总数 36 对，**拆法多算了一条**（把汇总子例又数进了形状里）。这是计数口径，不是断言被换弱，见 §6。

---

## 1 攻击 1（D1）：折叠的边界——字面量互拼真折了，**常量名/函数拼/切片变量/展开没折**

### 1.1 造了什么（逐字，探针 `probes/zzz_audit_d1d2_windows_test.go`）

```go
const zzzD1SuffixSameFile = "-pan" + "el-profile"                     // :21 同文件 const

func zzzD1a() *PanelManager {                                          // :23
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-d1a"+zzzD1SuffixSameFile))
}

func zzzD1b() *PanelManager {                                          // :27  const 定义在另一个文件
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-d1b"+zzzD1SuffixOtherFile))
}

func zzzD1c() *PanelManager {                                          // :31
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(),
		fmt.Sprintf("%s%s", "wisp-d1c", "-panel-profile")))
}

func zzzD1d() *PanelManager {                                          // :36
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(),
		strings.Replace("wisp-d1d_PANEL_profile", "_", "-", -1)))
}

func zzzD1e() *PanelManager {                                          // :41  切片先绑变量
	parts := []string{"wisp-d1e", "-panel-profile"}
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), strings.Join(parts, "")))
}

func zzzD1f() *PanelManager {                                          // :46  []string{...}... 展开
	return NewPanelManager(nil, nil, filepath.Join(append([]string{os.TempDir()}, "wisp-d1f", "-panel-profile")...))
}

func zzzD1g() *PanelManager {                                          // :50  Join 三参
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-d1g", "panel-profile"))
}

func zzzD1h() *PanelManager {                                          // :54  filepath.Join 嵌套
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(),
		filepath.Join("sub", "wisp-d1h-panel-profile")))
}

func zzzD1i() *PanelManager {                                          // :59  Join 结果再贴一个 Join
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-d1i")+filepath.Join("", "-panel-profile"))
}
```

另一文件 `probes/zzz_audit_d1_consts_windows_test.go:8`：`const zzzD1SuffixOtherFile = "-panel-profile"`。

### 1.2 尺子红了还是静默

同一次扫描（`a1-a2-scan.log`，278 源文件＝274＋4 枚探针）：

| 形 | 探针 `:行` | 尺子 | 红句点名的行 |
|---|---|---|---|
| `"a"+suf`，`suf` 是同文件 const | `:24` | **静默** | — |
| `"a"+suf`，`suf` 是**别处** const（票面原形） | `:28` | **静默** | — |
| 叶由 `fmt.Sprintf` 拼 | `:33` | **静默** | — |
| 叶由 `strings.Replace` 拼 | `:38` | **静默** | — |
| 切片先绑变量再 `strings.Join` | `:43` | **静默** | — |
| `append([]string{…}, …)...` 展开 | `:47` | **静默** | — |
| `Join` 三参 | `:51` | **红** | `zzz_audit_d1d2_windows_test.go:51`，叶 `wisp-d1g\panel-profile` |
| `filepath.Join` 嵌套 | `:56` | **红** | `:56`，叶 `sub\wisp-d1h-panel-profile` |
| `Join(...)`＋`Join("", suffix)` | `:60` | **红** | `:60`，叶 `wisp-d1i-panel-profile` |

汇总行：`扫了 278 个源文件，形状命中 7 处；其中已豁免 2 处、实红 5 处`（5 红＝D1 的 `:51/:56/:60` ＋ D2 的 `:94/:99`），`FAIL`、rc=1。

### 1.3 它凭什么（读补洞后的折叠加）

- **Ident 不是常量求值器**：`panelProfileFold` 的 `*ast.Ident` 分支（`:637-641`）只查 `roots[node.Name]`，不在根表里就整枚 `opaque`；
  `panelProfileFixedJoinIn`（`:562-574`）在 `:565` 把 `folded.opaque` 直接判 false。⇒ 常量名（`const suf = "-panel-profile"`）
  与普通变量一视同仁，**没有类型信息、没有 `go/types` 常量求值**。同文件/另一文件都一样。
- **只折三种被点名的调用**：`*ast.CallExpr` 分支（`:655-666`）只认 `os.TempDir` 类根、`strings.Join`、`Join`；
  `fmt.Sprintf`、`strings.Replace`、`append(...)...` 落到 `:666` 的 `opaque`。
- **`strings.Join` 只吃字面量切片**：`panelProfileFoldStringsJoin`（`:694-721`）在 `:699-702` 要求第一实参是 `*ast.CompositeLit`；
  `parts` 是 `*ast.Ident` ⇒ opaque。
- **展开参数是 `*ast.Ellipsis`**：`panelProfileFoldJoin`（`:672-690`）逐参 `fold`；`…` 实参落进 `:666` 的兜底 opaque，整个 Join 连坐。

### 1.4 判读

**D1 只修好了「字面量与字面量互拼」这一半**：它点名的五形（A2a–A2e）由常驻控制的 5 条必红子例挡住、本次跑全 PASS；
这次另造的三参 Join 与嵌套 Join 也红（这两条不在上一轮审计的 8 洞里，本件只判现尺射程，不回溯旧尺）。但票面这次要求的「变量参与」——`suf` 是**全常量**、只是有名字——静默；`Sprintf`／`Replace`／切片变量／`...` 展开也静默。
头注释说「a leaf assembled from pieces is judged as the name it spells」；**一枚具名 const 拼出来的叶也是 pieces**，这句在射程外时没有写。
D1 的洞**没有全修**，只修了写法最字面的那一半。

---

## 2 攻击 2（D2）：根的边界——`Getenv(变量)`/`Dir(Executable)`/根经 `+` 重绑仍静默

### 2.1 造了什么（逐字，同一探针文件）

```go
func zzzD2a() *PanelManager {                                          // :65
	return NewPanelManager(nil, nil, filepath.Join(os.Getenv("TMPDIR"), "wisp-d2a-panel-profile"))
}

func zzzD2b() *PanelManager {                                          // :69
	const envName = "TEMP"
	return NewPanelManager(nil, nil, filepath.Join(os.Getenv(envName), "wisp-d2b-panel-profile"))
}

func zzzD2c() *PanelManager {                                          // :74
	exe, _ := os.Executable()
	return NewPanelManager(nil, nil, filepath.Join(filepath.Dir(exe), "wisp-d2c-panel-profile"))
}

type zzzD2Dirs struct{}
func (zzzD2Dirs) TempDir() string { return "" }
func zzzD2d() *PanelManager {                                          // :83
	d := zzzD2Dirs{}
	return NewPanelManager(nil, nil, filepath.Join(d.TempDir(), "wisp-d2d-panel-profile"))
}

func zzzD2e() *PanelManager {                                          // :88
	dir, _ := os.MkdirTemp("", "wisp-d2e")
	return NewPanelManager(nil, nil, filepath.Join(dir, "wisp-d2e-panel-profile"))
}

func zzzD2f() *PanelManager {                                          // :93  path.Join（不是 filepath）
	return NewPanelManager(nil, nil, path.Join(os.TempDir(), "wisp-d2f-panel-profile"))
}

func zzzD2g() *PanelManager {                                          // :97  根先绑给局部名
	root := os.Getenv("TEMP")
	return NewPanelManager(nil, nil, filepath.Join(root, "wisp-d2g-panel-profile"))
}

func zzzD2h() *PanelManager {                                          // :102 根经 `+` 重拼后再绑局部名
	root := os.Getenv("TEMP") + string(filepath.Separator) + "wisp"
	return NewPanelManager(nil, nil, filepath.Join(root, "wisp-d2h-panel-profile"))
}
```

### 2.2 尺子红了还是静默

| 形 | 探针 `:行` | 尺子 | 说明 |
|---|---|---|---|
| `os.Getenv("TMPDIR")` | `:66` | **静默** | 头注释只认 TEMP/TMP；这条是**具名的刻意射程**，不是意外 |
| `os.Getenv(envName)`，`envName` 是全常量 | `:71` | **静默** | 实参必须是 `*ast.BasicLit` |
| `filepath.Dir(os.Executable())` | `:76` | **静默** | `Dir` 不是根调用，`exe` 也没被登记成根 |
| `t.TempDir()` 形（方法调用） | `:85` | **静默** | 每发独立，**不是缺陷** |
| `os.MkdirTemp` 返回值 | `:90` | **静默** | 每发独立，**不是缺陷** |
| `path.Join`（`path` 包） | `:94` | **红** | `panelProfileCalleeName` 只看最后一段 `Join`，`path.Join` 被当 `filepath.Join`（过近似，不算错） |
| `root := os.Getenv("TEMP")` 再 Join | `:99` | **红** | `tempRoots` 的 `record` 把直接根源登记为别名 |
| `root := os.Getenv("TEMP") + sep + "wisp"` 再 Join | `:104` | **静默** | `record` 只登记「RHS 本身是根源」，BinaryExpr 不登记；随后 `root` 不在根表 |

读数出处同上（`a1-a2-scan.log` 的 7 命中里 D2 是 `:94/:99` 两条）。

### 2.3 它凭什么

- `panelProfileGetenvIsRoot`（`:874-891`）：实参必须是 `*ast.BasicLit`（`:878-880`），且字面量 ToUpper 只能是 TEMP/TMP（`:886-890`）。
  ⇒ `Getenv(变量)`、`Getenv("TMPDIR")` 都不认。
- `panelProfileIsRootSource`（`:848-869`）：只认 `os.TempDir`／`os.UserHomeDir`／`os.Getenv` 三枚 selector 调用，加已登记别名（`:850-851`）；
  `filepath.Dir(...)`、`os.Executable()` 都不在。
- `panelProfileTempRoots` 的 `record`（`:795-809`）：`panelProfileIsRootSource(rhs[0], roots)` 为真才登记。`root := os.Getenv("TEMP") + "\\wisp"` 的
  `rhs[0]` 是 BinaryExpr ⇒ 不登记；后面的 `Join(root, …)` 因而 `rooted=false`。**这是一条真缺陷**：共享 temp 下的固定夹，只因根多拼了一段就全盲。

### 2.4 判读

**D2 修好了它点名的三种调用与盘符字面量（A3a/A3b/A3c），但根这一侧仍有三条真静默**：`Getenv(常量变量)`、`filepath.Dir(os.Executable())`、
「根先经 `+` 拼一段再绑局部名」。`TMPDIR` 与 `t.TempDir()`／`os.MkdirTemp` 的静默可以接受（前一条是具名射程，后两条每发独立）。⇒ **D2 漏了**。

---

## 3 攻击 3（D3）：包级 seed 的边界——只 seed **同文件**、只 seed `ValueSpec`

### 3.1 造了什么（逐字）

文件 A `probes/zzz_audit_d3a_decl_windows_test.go`：

```go
var zzzD3OtherFileJoin = filepath.Join(os.TempDir(), "wisp-d3a-quicklook")   // :15 包级固定 Join，叶不含 profile

var zzzD3InitValue string                                                    // :17
func init() {                                                                // :19
	zzzD3InitValue = filepath.Join(os.TempDir(), "wisp-d3b-quicklook")       // :20 init 里赋值
}

var zzzD3Once sync.Once                                                      // :23
var zzzD3OncePath string                                                     // :25
func zzzD3InitOnce() {                                                       // :27
	zzzD3Once.Do(func() {                                                    // :28
		zzzD3OncePath = filepath.Join(os.TempDir(), "wisp-d3c-quicklook")    // :29 sync.Once 里赋值
	})
}

func zzzD3bUse() {                                                           // :33
	zzzD3InitOnce()
	_ = NewPanelManager(nil, nil, zzzD3InitValue)                            // :35
}
func zzzD3cUse() {                                                           // :38
	zzzD3InitOnce()
	_ = NewPanelManager(nil, nil, zzzD3OncePath)                             // :40
}
```

文件 B `probes/zzz_audit_d3b_use_windows_test.go`：

```go
func zzzD3aUse() *PanelManager {                                             // :8
	return NewPanelManager(nil, nil, zzzD3OtherFileJoin)                      // :9
}
```

对照 `probes/zzz_audit_d3c_samefile_windows_test.go`（同文件声明＋同文件使用）：

```go
var zzzD3SameFileJoin = filepath.Join(os.TempDir(), "wisp-d3x-quicklook")     // :14
func zzzD3SameFileUse() *PanelManager {                                       // :16
	return NewPanelManager(nil, nil, zzzD3SameFileJoin)                       // :17
}
```

### 3.2 尺子红了还是静默

| 形 | 在哪 | 尺子 |
|---|---|---|
| 包级固定 Join 在**另一文件**，函数里交给 `NewPanelManager` | 声明 A`:15`，使用 B`:9` | **静默** |
| 包级变量在 `init()` 里赋值，同文件函数里交给 sink | A`:20` → A`:35` | **静默** |
| 包级变量在 `sync.Once` 里赋值，同文件函数里交给 sink | A`:29` → A`:40` | **静默** |
| **同文件**包级 `var = Join(...)` + 函数里交给 sink（对照） | C`:14` → C`:17` | **红**，叶 `wisp-d3x-quicklook`，点名 C`:14` |

读数：`a3-scan.log`（只放 A/B 两件）＝ **276 源文件 / 命中 2（两枚豁免）/ 实红 0**、rc=0；
`a3-contrast.log`（再加 C 对照）＝ **277 源文件 / 命中 3 / 豁免 2 / 实红 1**、rc=1，红句只点 C 的 `:14`。

### 3.3 它凭什么

- **seed 是逐文件算的**：`panelProfileHitsInSource`（`:401`）每次只拿一个 `src`，`:405` 用 `panelProfilePackageFixedJoins(…, src.file, …)`
  取本文件的包级绑定，`:445-448` 再把这份表 seed 进本文件每个函数体。⇒ **另一文件**声明的 `zzzD3OtherFileJoin` 永远不会出现在 B 的 seed 表里。
- **只收 `ValueSpec`，不收赋值**：`panelProfilePackageFixedJoins`（`:537-556`）只遍历 `GenDecl` 里的 `*ast.ValueSpec`（`:545-553`）；
  `init()` 与 `sync.Once` 里的是 `AssignStmt`，不收。`panelProfileSinkRule` 遇到 `AssignStmt` 只把名字存进**本次遍历的** `stored`（`:497-500`），
  并不 `add(fixed)`；`stored` 是每次调用新建的（`:477-480`），出了这个函数体就丢。⇒ init/Once 的赋值对别的函数不可见。

### 3.4 判读

**D3 只修好了「同文件、包级 `ValueSpec` 初值、叶再交给 sink」这一种形状（X2 原形）**。票面这次明确点的三种变体——包级变量在**另一个文件**、
`init()` 里赋值、`sync.Once` 里赋值——**全静默**。补洞把「包级作用域」当成一个**单文件语法作用域**处理，不是包作用域。⇒ **D3 漏了**（只修了示例，没修作用域）。

---

## 4 攻击 4（D4）：三条新正控有没有牙——**逐臂对应，没有错配**

### 4.1 造了什么（逐字，仓外 `teeth` 树里改尺子；`a4-teeth.log`）

每次只改一处，`go test -c` 都 rc=0，然后只跑 `-test.run '^Test314'`：

| # | 突变（`String.Replace` 的 from → to，逐字） | build/run |
|---|---|---|
| T1 | `case "NewPanelManager", "newPanelControl":` → `case "NewPanelManager":` | rc=0 / rc=1 |
| T2 | `isSel && sel.Sel != nil && panelProfileSinkKey(sel.Sel.Name)` → `isSel && sel.Sel != nil && false` | rc=0 / rc=1 |
| T3 | `hits(expr.Value)` → `_ = expr.Value` | rc=0 / rc=1 |

### 4.2 尺子红了哪一条

```
=== T1 ... fails=1
    --- FAIL: .../the_newPanelControl_sink_arm
=== T2 ... fails=1
    --- FAIL: .../the_DataPath_field-assignment_sink_arm
=== T3 ... fails=2
    --- FAIL: .../a_package-level_DataPath_key_is_a_sink_where_it_is_declared
    --- FAIL: .../the_DataPath_key_sink_arm_inside_a_function_body
```

- **T1 只红 `newPanelControl` 臂**，没有误伤别的；
- **T2 只红字段赋值臂**（`c.DataPath = …`，源码 `:501-503`）；
- **T3 红两条**，但两条都是**同一个** `KeyValueExpr` 臂（`:523-529`）的覆盖：函数体内 `DataPath:` 键（`:1520`）与包级 `DataPath:` 键（`:1429`）——
  这不是错配，是一条臂本来就该由两条控制共同钉住。

还原后 `build_rc=0 run_rc=0`，两测全 PASS；`RESTORED_SHA == ARCHIVE_SHA == 6E54A0B3…`。

### 4.3 判读

**D4 真修好**：三条臂各自有常驻控制，打断哪条就红哪条（T3 是 1 臂/2 控制，合理）。没有出现「打断 A 却红 B」的正控错配。
**但**：这三次 `FAIL` 的**同一次运行**里，`Test314FixedPanelProfileLeafShapeIsRed` 自己的汇总行仍印
`314 形状尺【绿】…实红 0 处，豁免自检红 0 处`（`a4-teeth.log` 每段 `SUM` 行）。这不是正控错配，是 D5 的口径只覆盖了本测试的两个计数器，见 §5。

---

## 5 攻击 5（D5）：汇总行——两种被指的读数**真红了**，但「任何一处红」仍说大了

### 5.1 造了什么（逐字）

**(a) 实红 K>0**：导出树 `cmd/wisp` 加一枚探针 `zzz_audit_d5k_windows_test.go`（叶不含 `profile`，只能由 SINK 抓）：

```go
func zzzD5k() *PanelManager {
	return NewPanelManager(nil, nil, filepath.Join(os.TempDir(), "wisp-d5k-quicklook"))
}
```

**(b) 自检红 F>0**：删掉 holder 理由句 `Nothing here creates a window (the headless half never asks bringUp), `（整句从 `panel_holder_319_windows_test.go:46` 拿掉），
其余一字不动；跑完从备份原字节拷回（`SHA_AFTER_RESTORE` 与原件同串）。

### 5.2 尺子印了什么

```
# (a) K>0
zzz_audit_d5k_windows_test.go:11 命中固定名面板 profile 形状（叶 "wisp-d5k-quicklook"）
314 形状尺【红】：扫了 275 个源文件，形状命中 3 处；其中已豁免 2 处、实红 1 处，豁免自检红 0 处
--- FAIL: Test314FixedPanelProfileLeafShapeIsRed      [rc=1]

# (b) F>0
豁免 cmd/wisp/panel_holder_319_windows_test.go 的理由句在源码里找不到了："Nothing here creates a window (the headless half never asks bringUp), so"
314 形状尺【红】：扫了 274 个源文件，形状命中 2 处；其中已豁免 2 处、实红 0 处，豁免自检红 1 处
--- FAIL: Test314FixedPanelProfileLeafShapeIsRed      [rc=1]
```

日志：`a5-k.log`、`a5-f.log`。⇒ 上一轮 D5 指的两枚读数（未豁免命中、过期豁免）**都真的印红**，不是绿。

### 5.3 但它仍有半句说大了

`panelProfileSummaryLine`（`:247-256`）的 `clean := red == 0 && faults == 0` 只看**这一个测试**的两个计数器：`hits-exempted` 与 `faults`。
于是 **正控表自己红了**（T1/T2/T3 的 `Test314PanelProfileShapeRulerRingsOnTheFixedNameShape` 子例 FAIL）时，
`Test314FixedPanelProfileLeafShapeIsRed` 的汇总行照印 `【绿】`（`a4-teeth.log` 三段 `SUM` 行逐字）。
实际判词句尾是「**只有「实红 0 且自检红 0」才是绿，任何一处红都不能当绿证**」——前半句对，**后半句大于射程**：
正控表的红不在它的两枚计数里。整包/整次测试仍以 `FAIL`/rc=1 收口，所以不是旧 D5 那种「测试失败而汇总说 0 红」的误导，**是措辞没跟上射程**。
⇒ D5 的两枚读数**真修好**；「任何一处红」这句应改成「实红与自检红都为 0 才是绿」，别把另一枚测试的红也算进承诺。

---

## 6 攻击 6（复核）：既有 36 子例有没有被悄悄变钝——**没有；14 例逐字同块，断言只变强**

### 6.1 14 例逐条对照

上一轮审计读到的 14 例来自 `c5be2d703^` 的旧尺（blob `cb308abd…`，896 行，存 `probes/ruler2-audit/orig/old_ruler.go`）。
与 HEAD 的 35 例子例逐块做**空白归一化后的整块比较**（`a6-cases.txt`）：

```
old cases: 14  HEAD cases: 35  (plus 1 summary subtest = 36 RUN subtests)
 1 SAME | the audit's revert, default tier verbatim
 …
14 SAME | prose in a comment and text in a quoted string
old 14 all present and byte-identical after whitespace normalisation: True
```

⇒ 14 例的 `name`／`source`／`want` 全部**同块**（旧例没有 `lineNeedle`，新例也没有给它们补加）：不是「换了条更弱的断言」，是原样保留。

### 6.2 断言逻辑有没有被换弱

| 项 | 旧尺（`old_ruler.go:844-873`） | 新尺（HEAD `:1541-1576`） | 判 |
|---|---|---|---|
| 叶集合 | `strings.Join(got,"|") != strings.Join(want,"|")` | 同 | 一样强（精确集合） |
| 每命中行 | `hit.line < 1 || !Contains(hit.code, hit.leaf)` | `hit.line < 1 || hit.raw == "" || !Contains(hit.code, hit.raw)` | **更强**：新增 `raw` 非空要求；折叠叶用起始字面量验行，不会因为折叠值不在那一行而误判 |
| 点名行 | `hits[0].line == lineOf("\""+want[0]+"\"")` | `lineNeedle` 可显式给；未给时同样用 `"`+want[0]+`"` | 一样强（折叠叶的用例显式给 needle） |
| 静默例 | `len(want)==0` → 只比集合空 | 同 | 一样强 |

⇒ **没有任何一条旧例被换弱**。新增 21 条形状例＋1 条汇总例，是「14→35 形状＋1 汇总＝36 子例」；
§0.2 已记证据件把 36 说成「36 形状＋1 汇总」的**多算一条**，与用例强度无关。

---

## 7 复跑（逐字）

```powershell
# 导出树（HEAD a46352fd3；尺 SHA 6E54A0B3…）
git -C "D:\work\workspace\projects plans\Wisp" archive HEAD -o C:\Temp\wisp-ruler2-audit-20261114\src.zip
Expand-Archive C:\Temp\wisp-ruler2-audit-20261114\src.zip C:\Temp\wisp-ruler2-audit-20261114\tree
# 起二进制（DLL 必须挨着 exe；直接跑 go test ./cmd/wisp 的 0xc0000135 是假红）
cd C:\Temp\wisp-ruler2-audit-20261114\tree
go test -c -o ..\stage\wisp.test.exe ./cmd/wisp
Copy-Item "D:\work\workspace\projects plans\Wisp\third_party\sherpa-onnx\*.dll" ..\stage\
cd C:\Temp\wisp-ruler2-audit-20261114\tree\cmd\wisp
..\..\stage\wisp.test.exe -test.run=^Test314 -test.v -test.count=1
```

探针（只建不删）：`probes/` 下 4 件（D1/D2 一枚、D1 跨文件 const 一枚、D3 声明/使用/同文件对照三枚）。
日志：`logs/` 下 9 份。旧尺与 HEAD 尺：`orig/` 两份（blob 分别为 `cb308abd…`／`627c2fcd…`）。

---

**末尾一句：D1–D5 的真账是——D1 只把「字面量与字面量互拼」折进来（常量名、`Sprintf`、`Replace`、切片变量、`...` 展开仍静默），D2 只把四根里的三种调用与盘符字面量认进来（`Getenv(变量)`、`Dir(os.Executable())`、根经 `+` 重绑仍静默），D3 只把同文件包级 `ValueSpec` seed 进来（跨文件、`init()`、`sync.Once` 仍静默），D4 三条正控逐臂对应、真修好，D5 两种被指的读数真修好但「任何一处红」是措辞大于射程——所以 D1/D2/D3 是「修了它举的那几个例子」，D4/D5 才是「修了机制」。**
