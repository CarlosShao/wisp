# 票 35 `:75` (a) 支｜夹具两面恒真的建模——写腿 `35-r4` 的盘上遗产＋编排者代跑读数

**⛔ 本件的性质先说清**：写腿 `35-r4`（39 次调用／442 万 token／35 分 10 秒）**死于每日额度**，不是死于任务失败，也不是死于服务端掐上下文。它把**产码面（这里是夹具面）做完了并自己跑绿**，但没写成交付件、没 commit。
⇒ §0–§2 是**它的活**（我按行号复量过才写），§2 的突变读数**是我代跑的**（它那一版四发全部没跑成，根因见 §2.1），§3 起是**它的原话＋我补的口径**。
⇒ ★**判语一律不代填**：`:75` 能不能翻、这三枚牙算不算数，归非实现者 `35-v4`。

---

## 0. 起手锚（腿自己 commit 的那笔 `1eec984c`）

锚 `6c7dd5ed`（派单写的 `af68866c` 之后共享树已前进，腿按现量记）；禁区三方差 `1d70fb2b..HEAD` 现量空（`cmd/wisp/panel_host_windows.go`、`internal/panel/`、`frontend/`、`docs/specs/`、`docs/PLAN.md`）；`bridge.go` blob `bebe8e70` 未动；夹具在 `1eec984c` 时 blob `7185ab56`。⇒ **本程产码零改动**，只动测试夹具这一面（`cmd/wisp/panel_transport_35r2_test.go`，＋245／−3）。

## 1. 落了什么（我自己 `grep -n` 复量到行的四处）

| 面 | 行 | 内容 |
|---|---|---|
| 可写性表（M-B 的根） | `:128`＋`:137` | `jsObject` 新增 `unwritable map[string]bool`，注释指名"35-r4's property-writability table（ledger A684 §3, ticket 35:75 (a)）" |
| 台件入口 | `:149-160` | `markNonWritable(names ...string)`——**Go 侧台件，不是 JS 可见方法**（`Object.defineProperty` 没有实现成页面能调的东西） |
| 覆写被拒 | `:163-166` | `set` 现在读表：`if _, exists := o.props[name]; exists && o.unwritable[name] { return }`（＝静默忽略赋值，浏览器里非可写数据属性的形状） |
| receiver 严格（M-A 的根） | `:1285-1291` | 原生出口桩从 `func(_ jsValue, args …)` 改成**认 receiver**：`if recv != cw { panic("TypeError: Illegal invocation: chrome.webview.postMessage lost its receiver - called with …") }` |
| 两形世界 | `:1214-1218`＋`:1307` | `nativeExitNonWritable` 让 `openDocument` 在"非可写"世界里把 `chrome.webview.postMessage` 摆成非可写数据属性 |
| 载具 | `:1926` | `runShippedHookInOneWorld35r4(t, nonWritable)`——装**生产接线**（`installPanelTransport`：绑门＋把产码钩子字符串当值读入）→ 发一次页面自己那封信 → 交回断言材料 |
| 三枚用例 | `:1946`／`:1976`／`:2026` | `TestForwardingHookMustNotLoseTheNativeExitReceiver`（M-A）／`TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit`（M-B）／`TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent`（守卫 typeof 那半支的落点） |

腿自己跑的颜色（`logs/r4-run.txt`，22:28）：这三枚＋r2 那三枚＋r3 那三枚＝**12 枚全 PASS，rc=0**。

## 2. 四发反形读数（★编排者代跑，2026-10-08 08:17）

### 2.1 先记一枚坑：腿那一版四发**全部没有读数**

腿的 `logs/mut-summary.txt` 四节都写着 `PASS=0 FAIL=0 buildrc=1 runrc=127`。⛔ 这不是"突变没红"，是**编译就没成**：
`go: parsing overlay JSON: invalid escape sequence \w in string`。
根因在它台件的 `mk_overlay`：替换侧路径用 `cygpath -m`（正斜杠）而**被替换侧用 `cygpath -w`（反斜杠）**，Windows 反斜杠直接进 JSON 字符串就是非法转义 ⇒ Go 拒读 overlay ⇒ 没产出 exe ⇒ 后面那句 `panel-$m.exe` 找不到文件（`runrc=127`）。
⇒ 我另写一份 `scripts/orchestrator-countershape.sh`，**唯一改动＝那两条路径也走 `-m`**，突变源、`-test.run` 名册、断言全部沿用腿的原件（它的 `mut/*.go` 四枚拷贝我一字未动，`logs/orch-mut-landing.txt` 现量每枚仍是 `changed-line-count=2`＝恰一行被换）。腿原件保留在 `scripts/countershape.sh`，不覆盖、不改写。

### 2.2 代跑后的颜色（`logs/orch-mut-summary.txt`，四枚 buildrc 皆 0）

| 反形 | 改了哪一行（逐字） | 颜色 | 谁红 |
|---|---|---|---|
| `ma-hook` | 产码 `return native.call(cw, message);` → `return native(message);` | 9 FAIL／3 PASS，runrc=1 | **含具名那枚** `TestForwardingHookMustNotLoseTheNativeExitReceiver` |
| `ma-detector-gone` | 夹具 `if recv != cw {` → `if false {`（**与 ma-hook 同发**） | **0 FAIL／12 PASS，runrc=0** | —（摘掉探测器后那 9 枚全转绿） |
| `mb-set-gone` | 夹具 `exists && o.unwritable[name] {` → `exists && false {` | 1 FAIL／11 PASS | **恰那枚** `TestForwardingHookIsSilentlyUnarmedByANonWritableNativeExit` |
| `mb-typeof-gone` | 产码 `if (inside || typeof window.%[1]s !== "function")` → `if (inside)` | 1 FAIL／11 PASS | **恰那枚** `TestForwardingHookFallsBackToTheNativeExitWhenTheDoorIsAbsent` |

⚠ **两处我读出来的形状差别，不许抹平**：
- **M-A 那发是"级联红"（9 枚）不是"外科手术红"（1 枚）**：探测器自己会 `panic`，所以 receiver 一丢，凡走过原生出口的用例一起红。⇒ 归因**不能只靠"红了 9 枚"**，必须靠 `ma-detector-gone` 那一发（摘掉 receiver 检查 ⇒ 9 枚全转绿）来证明"红来自这枚新探测器，不是别处坏了"。这一对齐全，是 M-A 的牙。
- **M-B 与 typeof 那两发各只红 1 枚且名字唯一**，是更强的形状（红不可挪用）。
- ⛔ 我不据此判"`:75` (a) 成立"——那是 `35-v4` 的活。我把颜色、行数、根因交出来就停。

### 2.3 仓库没被突变洗过（`logs/orch-repo-untouched.txt`）

`cmd/wisp/panel_host_windows.go` 跑前跑后同 blob `26b5de83…`＝HEAD；夹具跑前跑后同 blob `31c96f37…`（＝本程那笔未提交编辑，overlay 不动它）。⇒ `prod-untouched=YES`、`fix-untouched=YES`、rc=0。

## 3. 射程——腿自己写的话（逐字抄，别由我转述）

> WHAT THIS BUYS AND WHAT IT DOES NOT. It buys "a hook that drops this, or whose assignment is silently refused, is no longer certified as a delivery". It does NOT buy "WebView2 really binds the receiver / really lets the page replace postMessage" - A684 §1 puts that 格 with ticket 35:52's real-window evidence, and no assertion here stands in for it.

**我补三句口径**：
1. 这枚尺**今天仍然不引真浏览器**：它跑的是手写的 JS 子集解释器＋库的 shim 拷贝，`-tags winlive` 那枚真窗凭据在 `:52`，⛔ 两格不许互相借光（第 108 条镜像形：一枚绿不许替另一枚绿作证）。
2. `markNonWritable` 是 **Go 侧台件**，页面上没有 `Object.defineProperty` 可用 ⇒ "页面自己 descriptor-read／descriptor-redefine 它的传输"这一族**在本尺外面**（下一节列全）。
3. 机主的裁度还在生效：⛔ 不合页面分支、⛔ 不推远程。`:75` 的 (c) 支（winlive 进 CI 的载体）没落地前，真窗读数**只有本机这一份**。

## 4. 未建模清单（腿指名"要写在 impl.md 里而不是悄悄跳过"的那几样）

- `Object.defineProperty` 作为 **JS 可见方法**——未实现（只有 Go 侧 `markNonWritable`）。
- `Object.getOwnPropertyDescriptor`（**读**描述符）——未建模。
- `configurable`（能否被删／再定义）——未建模；今天只有 `writable` 一维。
- `delete` 运算符——未建模。
- accessor 属性（`get`/`set` 那对）——未建模。
- ⇒ 后果要写死：把钩子写成"定义访问器属性"或"先 delete 再赋值"这类形状，**这枚尺今天看不见**。这些是 `:75` 的牙边界，不是已修好的洞。

## 5. 门禁（编排者现跑，2026-10-08 08:1x–08:2x）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| `go vet` | `go vet ./cmd/wisp/` | **0** | 空输出 |
| 卫生 | `cd tools/d22scan && go run . -root ../../` | **0** | clean；ban #8 覆盖 `cmd/` 108 枚 Go 文件（含注释与 `_test.go`）⇒ 本笔新增的英文注释里的中文"格／台件"两字符不触门 |
| 格式 | `gofmt -l` 对夹具 | 列出该文件 | ★**同一处 4 行漂移在 HEAD 版本里也在**（HEAD `:775-776`／工作树 `:807-808`，因插入而位移；CR 计数两版皆 0）⇒ 本笔**不新增**未格式化面；那 4 行是 `35-r2` 落下的，归票 275 那族格式化欠账，⛔ 不在本程顺手改（`ci.yml:188-195` 自己记着 run `37406757402` step 8 `[failure]` gofmt ⇒ step 9 `[skipped]`，这道门今天本就是红的） |
| 整包 | `cmd/wisp` 全量（原生 DLL harness：三枚 DLL 摆在测试 exe 旁＋CWD＝包目录） | 见 `logs/full-package-*.txt` | 名册与既有红族作差，逐枚在 `:5` 那节补 |

## 6. 留给 `35-v4`（非实现者）的问题，⛔ 我不代答

1. `ma-hook` 的级联红（9 枚）能不能算 M-A 的牙？`ma-detector-gone` 那发（12 全绿）够不够把它和"别处坏了"分开？
2. `mb-set-gone`／`mb-typeof-gone` 各红 1 枚——**红句的名字唯一性**要自己造一发混文案的突变来顶（票 35 `:63` 那三支的 V2 就是这么做塌的）。
3. §4 那五样未建模里，有没有哪一样能让**今天这三枚用例**在产码坏掉时照样绿？（也就是：`A684` 那两面之外，还有没有第三面。）
4. 腿的注释把台件语义写进了英文散文里且夹了中文字（`that格`、`markNonWritable台件`）——算不算该修的可读性缺陷，归它裁，⛔ 我不改。
