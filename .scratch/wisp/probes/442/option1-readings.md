# 候选①（自写单层过程间摘要）的现成读数 —— **未落地**

> **本件未落地，⛔ 当已交件读。** 票 442 尾部 `## 编排者改向裁定（2026-11-14 15:2x）` 把主判据改为候选③
> 并明写「本票不再往尺子里加规则」；本腿先照**派单**（按铁律选型后补 A1b/A1c）把①做完了、全部读数是真的，
> 随后读到该裁定 ⇒ 未提交、未 push、尺子已还原成 HEAD 那一枚（SHA `6E54…23F2`）。
> 落地物＝`option1-ruler.patch`（57796 字节；post-image SHA `9F9765322511BCE45D2F2F2170FC3E209BBC0BACB74E36E3198CA8221E98A4E2`）
> ＋本目录 `option1-panel_profile_leaf_shape_windows_test.go`（post-image 全文，84829 字节 / 2367 行 / CR=0）。
> 「442 改向复议」后一条 `git apply .scratch/wisp/probes/442/option1-ruler.patch` 即可落。
> ③ 的可行性判断见同目录 `option3-feasibility.md`。

---
# 票 442 · 面板 profile 形状尺补上「值跨函数边界」两形（A1b/A1c）

> 角色：**写腿（实现者）**，⛔ 审计者，不自审。上游＝票 442、`docs/evidence/s1/314-shape-ruler.md` §7（D1–D5）、
> `docs/evidence/s1/314-ruler-audit.md` §1（A1b/A1c 原始读数）。
> 写面：`cmd/wisp/panel_profile_leaf_shape_windows_test.go`（唯一 Go 件）、本件（新）、
> `.scratch/wisp/probes/442/**`（新）、票 442 Progress log（只追加一条 bullet）。
> ⛔ 未碰产码／其他测试／`go.mod`／`go.sum`／规格／阈值／golden／豁免表；未 push；⛔ 没跑整包 `./cmd/wisp`
> （只定向 `-test.run`，走 staged carrier）。尺子只有一枚跟踪件被改：`git diff --stat` ＝
> **1 file changed, 883 insertions(+), 146 deletions(-)**（2367 行 / 84829 字节 / CR=0）。
> 尺子 SHA256：补前 **`6E54A0B34D6A3CD79888F5CD33535EB0754F4549DDD3D2FE1FD8FFE3CAA723F2`**
> （＝ 314 补洞后那一枚；本腿导出树 `old` 现跑**同串**，证明"补前"读的就是 HEAD 的尺）、
> 补后 **`9F9765322511BCE45D2F2F2170FC3E209BBC0BACB74E36E3198CA8221E98A4E2`**。
> 仓外导出树：`C:\Temp\wisp-442-20261114`（`old`＝HEAD 的尺、`new`＝本腿的尺，两棵放同一枚探针，各 274 件 + 探针）。

**第一行结论：先选型（三候选逐种「买到什么／花什么」，选①自写有界摘要）；A1b、A1c 两个方向共三形在仓外
**可编译**探针里旧尺实红 0、新尺实红 3 并逐条点名到叶所在行；还原后尺子 SHA 同串且复绿；既有 **35 形＋1 汇总＝36 条**
子例一字未动全 PASS（票面写"36 形"是数重了，§4 给读数）；豁免表 2 条 783 字节同串；同一棵 275 件导出树
尺子 0.68s → 0.72s；四门禁 GREEN。**

---

## 1 AC#1 选型：跨过程该怎么读（判据＝最有价值＋最主流 harness）

> 铁律⓪⓪⓪：⛔ 用「风险小／动作小／门禁会红」降级；「复杂」也不加分——复杂必须买得到真价值。
> 下面三候选逐种写「买到什么／花什么」，最后是本腿拍的**一个**和**为什么它最有价值**。

### 1.1 三个候选

| 候选 | 买到什么（价值） | 花什么（代价与落点） |
|---|---|---|
| **① 自写单层（有界）过程间摘要** ✔选 | **补的是整类**：每个单返回值函数做一阶摘要——「返回值用形参怎么表达」；调用点把实参代入，两个方向都成立（根进形参／叶进形参）；helper 体内先算进局部再返回也跟（3 轮）；helper 自己的临时根也种上；同包跨文件同一条路。红句仍**点名叶所在的 file:line**（本票与审计都要的「点名」），第二句写下沉点。零新依赖，尺子仍是「一枚测试件，只读源码文本」。 | 尺子 1630 → 2367 行（+883/−146：机制约 400 行、常驻正控 8 条表内 ＋ 1 枚跨文件测试约 430 行、其余是把两条规则的作用域化）。运行 0.68s → 0.72s（同一棵 275 件树；冷文件缓存下 1.00s → 1.23s）。不跟的边界要**具名写死**（方法／字段／接口／函数值／多返回值不一致／环／深度>3）。 |
| **② 引 `golang.org/x/tools/go/ssa` ＋ callgraph** | 类型感知数据流：能穿结构体字段、切片、map、接口、方法、φ 节点；callgraph 还附送可达性。 | `go.mod`/`go.sum` 动（`x/tools` 连带 `x/mod`、`x/sync` 等），本仓依赖逐枚钉 pseudo-version、过供应链与 D22 审查——**票面要求「先报」**，本腿没报、也不该在没报的情况下动。运行时：SSA 要先把 `cmd/wisp` 在 winlive 标签下**类型检查＋构建**一遍（`go/packages` 会在测试进程里回调 `go list`／构建），而尺子今天故意只 `parser.ParseFile`；274 件＋全部依赖再走一趟，秒级预算被吃掉一大块。**语义上更贵**：折叠后的 `ssa.Const` **没有字面量位置**，红句只能退化成「这个函数」——正好丢掉本票与审计都要的「点名到行」。可达性也不是这把尺要的：它故意连死代码一起判（本件的三枚探针全写在没人调用的函数里）。 |
| **③ 判据换成运行时证据**：测试里断言 `NewPanelManager`／`DataPath` **实收**的值不是固定名（用既有 fixture 注入点） | 零假阳性：断言的是这一发真跑到的值；不读源码就不会被注释/字符串骗。 | 只覆盖**这一发走到的路径**：源码里有固定名、但这次没走到（死函数、未走分支、另一 build tag）＝看不见——**而形状尺的全部价值就在静态覆盖**（审计 A3 已证：把默认档改回固定名，`go test -c` rc=0、`Test253r6` 仍 PASS）。要断言「实收值」还得有注入面：要么改产码让 `NewPanelManager` 记下收到的 dataPath（本票写面禁止产码），要么在 fixture 里重造取值路径（那就成了「测 helper 返回」，不是测下沉点）。它也满足不了 AC#2 的「仓外探针真红并点名」（探针只编译不跑窗）。 |

### 1.2 拍①，为什么它最有价值

1. **补的是类，不是两形。** 票面两形只是「值跨了函数边界」的两种写法；① 的判据是「实参代入被调函数的返回表达式」，
   于是 helper 返回、param 两向、helper 内局部、同包跨文件**同一套机制**全落进来（本件 8 条新控制逐条钉住，
   含 4 条"必须静默"的反向控制）。② 买更大的类（内存流）但丢了「点名到行」；③ 买的是另一射程（这一发），**对静态覆盖不成立**。
2. **它保住红句的可操作性**：命中身份仍是「叶所在 file + line + 叶名」，`via` 第二句写「值在 <sink file:line>
   交给 `NewPanelManager`」——读红句的人知道**改哪一行**、**为什么这行值得改**。② 的 IR 给不了这个。
3. **成本写成上界，不赌运气**：深度 3、环即「不知道」、摘要按（函数, 实参值）记忆化，且**只有没被环/深度截断的计算才进缓存**
   （截断结果依赖起点，不能复用）。这是可复跑的数字，不是「保守一点」的托词：深度不设界且不记忆化时同一棵树读到
   **3.44s**（本腿第一次现跑）→ 设界 3 后 1.47s → 加记忆化 1.25s（冷）／0.72s（热）。

**主流 harness 对齐**：这就是主流静态检查器在没有整程序 SSA 时用的**摘要式过程间分析**——函数摘要＝「返回值用形参
怎么表达」，调用点代入，上下文深度有界，环用 guard／不动点处理（go vet 自己的若干 pass、semgrep 的 interprocedural
mode、CodeQL 的 local dataflow summary 都是这个形状）。本腿取它的最小可用子集：一层摘要＋有界深度＋记忆化，
**没有**引入整程序构建，也没有把判据从「源码形状」改成别的。

**③ 的射程外（票面要求写清）**：它对「源码里有固定名但这一发没走到」**不成立**——死函数、未走分支、另一 build tag
里的固定名，运行时断言永远看不见；而形状尺的价值恰恰是**静态**覆盖（下一发重构把它搬进活路径时，尺子仍然红）。
⇒ ③ 只配当补充证据，不配当判据。

---

## 2 改了什么（机制，可核）

1. **一层世界 `panelProfileWorld`**：扫描集解析一次，按「包目录 + 函数名」索引单返回值函数（同名多次声明＝build-tag
   对 ⇒ 标 ambiguous、不代入、不猜）；`lineOf`/`codeAt` 让任一文件的字面量都能被点名（红句点的是**叶所在文件**）。
2. **代入 `foldCall`**：`call.Fun` 是普通标识符、且同包目录唯一、实参个数与形参一致时，实参先按**调用方作用域**折叠，
   再在**被调函数自己的作用域**里折返回表达式。形参两向都成立：`mk(root, leaf)` 的 root 收 `os.TempDir()`
   （→ rooted）、leaf 收字面量（→ 文本），与内联写法同一套判据（即 A1b/A1c）。
3. **被调体的根与局部**：`panelProfileTempRoots` 先把 helper 自己绑的根种进去；局部绑定 `p := Join(...)` → `q := p`
   迭代 3 轮（`panelProfileLocalRounds`）；形参名覆盖同名包级根（遮蔽）。
4. **有界与记忆化**：深度 3（`panelProfileCrossDepth`）；环（`active`）与深度截断都把 `contextSkips` 加一 ⇒
   **被截断的计算不进缓存**；未截断的按（函数 + 每个实参的值与位置）记忆化。
5. **命中身份与第二句**：仍是（叶 display, 叶 line, 叶名）；跨边界时 `hit.via` ＝
   「值在 <sink file:line> 交给 <sink>，常量叶写在 <file:line>」。单表达式形状 `via` 为空——表里有一条控制
   专盯「没跨边界却写了跨边界说明」。
6. **没动的**：SINK 集合（`NewPanelManager`／`newPanelControl`／`DataPath` 字段／`DataPath:` 键）、四个根、
   折叠规则、名称规则（`profile` 拼写）、豁免表（2 条、783 字节同串）、汇总行口径、性能口径。

---

## 3 AC#2 仓外可编译探针：A1b / A1c 真红并点名，还原 SHA 同串再绿

探针 `.scratch/wisp/probes/442/zzz_442_probe_windows_test.go`（SHA256 `0E2B857001575FE518307100A4ECA7858D15A01A34CF013675BFFD203D457545`），
三形都写成**没人调用的函数**（尺子判源码形状、不判可达性）：

- **A1b**：`plainLeaf442Probe()` 返回 `Join(os.TempDir(), "wisp-442-a1b-quicklook")`（**避开 `profile` 拼写**，NAME 规则看不见它），调用处交给 `NewPanelManager`；
- **A1c(root)**：`paramRoot442Probe(root string)` 返回 `Join(root, "wisp-442-a1c-rootparam")`，调用处传 `os.TempDir()`（票面读法）；
- **A1c(leaf)**：`paramLeaf442Probe(leaf string)` 返回 `Join(os.TempDir(), leaf)`，调用处传字面量（审计 §1 原始读法）。

两棵导出树放同一枚探针，`go test -c ./cmd/wisp` **两棵都 rc=0**（探针能编译，不是"只被正则扫到的文本"）：

| 树 | 扫描集 | 尺子读数 | 结果 | 耗时 |
|---|---|---|---|---|
| `old`（HEAD 的尺，SHA `6E54…23F2`） | 275 件 | 命中 **2**／豁免 2／**实红 0** | **PASS**（三形全静默） | 0.68s（冷 1.00s） |
| `new`（本腿的尺，SHA `9F97…A4E2`） | 275 件 | 命中 **5**／豁免 2／**实红 3** | **FAIL** | 0.72s（冷 1.23s） |
| `new` 还原（删探针） | 274 件 | 2／2／**0** | **PASS** | 0.53s（冷 1.23s） |

`new` 那三条红句（逐字，`.scratch/wisp/probes/442/probe.log`）：

```
cmd/wisp/zzz_442_probe_windows_test.go:33 命中固定名面板 profile 形状（叶 "wisp-442-a1b-quicklook"）：
        return filepath.Join(os.TempDir(), "wisp-442-a1b-quicklook")
    这一形是跨函数边界读到的：值在 cmd/wisp/zzz_442_probe_windows_test.go:37 交给下沉点 NewPanelManager(...)，而常量叶写在 cmd/wisp/zzz_442_probe_windows_test.go:33。

cmd/wisp/zzz_442_probe_windows_test.go:42 命中固定名面板 profile 形状（叶 "wisp-442-a1c-rootparam"）：
        return filepath.Join(root, "wisp-442-a1c-rootparam")
    这一形是跨函数边界读到的：值在 cmd/wisp/zzz_442_probe_windows_test.go:46 交给下沉点 NewPanelManager(...)，而常量叶写在 cmd/wisp/zzz_442_probe_windows_test.go:42。

cmd/wisp/zzz_442_probe_windows_test.go:55 命中固定名面板 profile 形状（叶 "wisp-442-a1c-leafparam"）：
        _ = NewPanelManager(nil, nil, paramLeaf442Probe("wisp-442-a1c-leafparam"))
    这一形是跨函数边界读到的：值在 cmd/wisp/zzz_442_probe_windows_test.go:55 交给下沉点 NewPanelManager(...)，而常量叶写在 cmd/wisp/zzz_442_probe_windows_test.go:55。
```

**还原读数**：删探针后 `new` 树复跑 **274 件／2 命中／2 豁免／实红 0／PASS**，尺子 SHA **`9F97…A4E2` 三处同串**
（`REPO_RULER_SHA` ＝ `NEW_RULER_SHA` ＝ `RESTORED_RULER_SHA`；脚本第一步就从仓内拷、最后一步重新读），
探针 SHA 未变。⇒ 红是探针带来的，不是改尺子凑的；同一棵树、同一枚探针、同一条命令下，旧尺对这三形**零红**。

---

## 4 AC#3 既有子例逐条保留（含一条计数更正）

**计数更正（读数优先，不跟票面数字走）**：HEAD 的 `Test314PanelProfileShapeRulerRingsOnTheFixedNameShape`
表里是 **35 个形状子例**（`cases := []struct{…}` 段内 `name:` 计数＝35）＋**1 个汇总子例**（单独一条 `t.Run`），
即 **36 条 `t.Run`**。票面与 314 件写的「36 个形状子例＋1 汇总子例」把汇总行数了两遍；本腿按实测报：
**35 形＋1 汇总＝36 条，全部逐条保留、全 PASS**，本腿在表尾**追加** 8 条（不重排、不改旧名）⇒ 现表 43 形＋1 汇总＝44 条。

- 逐条保留：两棵导出树子例名单比对，**35 条旧名 missing ＝ 0**（`OLD_SUBTESTS=36`／`NEW_SUBTESTS=44`／`NEW_ADDED=8`）。
- 追加 8 条（表内）：A1b 必红、A1c(root) 必红、A1c(leaf) 必红、helper 局部必红、per-run 叶必静、无下沉点必静、
  实参填不进必静、两个返回值不一致必静；新增控制都带 `viaNeedle`，并反向钉「单表达式形状不得写跨边界说明」。
- 表外新测试 `Test442CrossFileHelperReturnReachesTheSink`：helper 在 `cmd/wisp/zzz_442_cross_helper.go`、下沉点在另一文件 ⇒
  断言命中点名**叶所在文件**、`via` 写出下沉点；再抽掉 helper 文件 ⇒ 同一 sink 源读 **0 命中**（证明这红是跨过程读出来的）。
- 豁免表与 HEAD **逐字节同串**（整块比对 `same=True`，783 字节、2 条）——没加一条豁免、没改一个字符。
- 仓内现跑（`.scratch/wisp/probes/442/repo-final.log`）：`^Test314|^Test442` 全 PASS，`ok … 1.00s`，
  `Test314FixedPanelProfileLeafShapeIsRed` **0.73s**，形状尺读 **274 件／2 命中／2 豁免／实红 0**（与补前同数）。

---

## 5 AC#4 性能（现跑数字）

同一棵 275 件导出树、同一条 `-test.run=^Test314FixedPanelProfileLeafShapeIsRed`，同一支脚本里背靠背跑
（同一轮：`old` → `new` → 还原；两轮＝冷/热文件缓存各一次）：

| 尺 | 扫描集 | 冷跑 | 热跑 |
|---|---|---|---|
| HEAD（补前） | 275 | **1.00s** | **0.68s** |
| 本腿（补后） | 275 | **1.23s** | **0.72s** |
| 本腿（还原到 274 件，仓内同形） | 274 | 1.23s | **0.53s** |

仓内 staged carrier 现跑：`Test314FixedPanelProfileLeafShapeIsRed` **0.73s PASS**；`…RingsOnTheFixedNameShape`
（44 条子例）**0.01s**；`Test442CrossFileHelperReturnReachesTheSink` **0.00s**。⇒ 仍是**秒级**，增量 **+0.04s（热）／+0.23s（冷）**。
两条把增量压住的是必需件：深度不设界且不记忆化 = **3.44s**（本腿第一次现跑），设界 3 = 1.47s，再加记忆化 = 1.25s。

---

## 6 AC#5 门禁四数（本件最终态现跑）

| 门禁 | 读数 |
|---|---|
| `go vet -tags winlive ./cmd/wisp` | **rc=0**（`probes/442/vet.log` 0 字节） |
| `gofumpt -l cmd/wisp` | **空**（0 行；gofumpt v0.12.0 / go1.27.1；`probes/442/gofumpt.log` 0 字节；本机 `core.autocrlf=true`，该枚工作树副本 CR=**0**，与入库 blob 同字节） |
| `sh scripts/d22scan.sh` | **rc=0 clean**：`examined 481 production Go files`；`ban #8 cmd/=280`（把这枚尺算进分母）、`internal/=1010`、`frontend/=123`、`design/=39`、`#10 internal/risk/=21` |
| `sh scripts/check-path-length-budget.sh` | **rc=0 VERDICT GREEN**：tracked=**19970**／over-budget=**69**／roster=**69**／not in roster=**0**／wall interval 内 **0**／longest=185 |

（长度门读的是**已跟踪**路径，本件三枚新件在提交前不计入；三枚路径长 62／48／35 字符，提交后复跑仍 GREEN。）

---

## 7 诚实边界（本节不吹的地方）

- **方法调用、字段、切片、map、接口、函数值**：仍 opaque（接收者/内存里的值不跟）。审计的 A1b/A1c 是普通函数，已读；
  方法那一类要另一件仪器。
- **两个返回值不一致的 helper 不做摘要**：这种 helper 传出的非 profile 叶**不跟**（每个 return 仍就地按 NAME 规则判）。
  表里有一条必静控制钉住这条边界，别把它当「已覆盖」。
- **环与深度 >3**：读成「不知道」，不是「不是缺陷」。深度 3 是取舍，不是证明。
- **跨目录不代入**：摘要按包目录隔离（`cmd/wisp` 与 `scripts/spike/webview2-latency` 都是 `package main`，但不是同一个包）。
- **只按名解析**：同包目录内同名两次（build-tag 对）⇒ ambiguous、不代入——宁可静默，不猜。
- **常量折叠的既有限制**（不改）：`Join` + 运行时后缀视为 per-run；`os.UserCacheDir`／UNC 不是根。
- **不是审计**：本件是写腿自证；按 `SPEC-12 §4.3`，补洞后的对抗验收要另派一条腿（D22 双角色）。
- 未跑整包 `./cmd/wisp`；未 push；未动阈值/golden/规格/产码；⛔ 没碰 `frontend/**`、`design/**`。

---

## 8 复跑命令（逐字）

```sh
# 仓内：只定向跑这两族（staged carrier；⛔ 裸 go test，会 DLL 假红 0xc0000135）
GOFLAGS='-tags=winlive' scripts/wisp-cli-staged-tests.sh -test.run '^Test314|^Test442' -test.v

# 仓外：两棵导出树 + 同一枚探针（old / new / 还原 三段读数；约 3 分钟）
pwsh -NoProfile -File .scratch/wisp/probes/442/probe.ps1

# 四条门禁
go vet -tags winlive ./cmd/wisp
gofumpt -l cmd/wisp
sh scripts/d22scan.sh
sh scripts/check-path-length-budget.sh
```
