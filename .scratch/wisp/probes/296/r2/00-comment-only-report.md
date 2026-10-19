# 票 296 `AC#5` 落地腿 `296-r2` 交件（★纯注释面·唯一产码写面＝`cmd/wisp/resident_windows.go` 的注释行）

起手锚（我自己 `git log -1` 现认，⛔ 用编排者写的号）＝`887dd356cd4eeb541df292b27c27db493657cff6`
（`dev`，Mon Oct 19 22:36:20 +0800）。交件笔＝`1d6b0cf9`（23:01:19）。

## 1. 三处改后注释原文（逐字，供编排者裁"是否只说行为、没新造承诺"）

`cmd/wisp/resident_windows.go` 里 `hotkeyReloadSource296` 文档块末段。尺＝`git diff -U0` 的 hunk 头逐字 `@@ -40,5 +40,11 @@`＝
**改前 `:40-44` 那 5 行 → 改后 `:40-50` 那 11 行**（`:39` `residentBallHotkeyChain258. The bridge does NOT keep…` 起头那行一字未动；
整块今天占 `:36-50`）。⚠ 编排者派单引的绝对行号（`AC#5` 三枚锚＝更早那颗 HEAD 上的 `:38/:40/:42`）在我这颗起手锚 `887dd356` 上
读作 `:40/:42/:44`（尺＝`git show 887dd356:cmd/wisp/resident_windows.go | sed -n '36,44p'` 逐字在场）＝内容锚对得上、行号随历史漂；
本腿一律按**内容锚**落刀。改后原文逐字：

```
// NO SLOT IS EVER ANSWERED EMPTY (ticket 296). Both branches merge through
// ball.ApplyHotkeyDefaults, which is the precondition ball.HotkeySource states
// in its own doc comment and the shape the construction half gets from
// residentBallHotkeyChain258. The bridge does NOT keep the current bindings on
// an empty answer, and Check() has two skips, neither of them an all-empty
// guard: a nil src keeps whatever is live, and a set DeepEqual to what was
// applied is left alone. Four empty strings from here would therefore differ
// from the applied set, be read at the registration end as disabled by config
// (never attempted), and leave summon, mute and panel dead until a read
// answers a DIFFERENT set - reading the same empties again only hits the
// DeepEqual skip. That is the regression this closure caused on the owner's
// machine. Whether an empty [hotkey] slot can be a user's way of turning a key
// off is Q-82, still on the owner's list; this closure takes no position on it
// and only says what it does today - it answers every empty slot with the
// compiled default before the bridge ever sees the value.
```

三处对照（①②③＝派单内容锚）：

| 锚 | 改前那句为什么是过头话 | 改后怎么落 |
|---|---|---|
| ① `Check()'s only skip is the DeepEqual diff` | 不成立：`Check()` 有**两条**不重绑出口，`r.src == nil` 那一条**恰恰就是**旧注释否认的"保持当前绑定"＝把 `AC#2` 删掉的旧承诺换形状又许一遍 | 改成 `Check() has two skips, neither of them an all-empty guard` ＋逐枚点名（nil src 保住现有的、DeepEqual 相等的不动） |
| ② `dead until the next read` | 把恢复条件说轻：再读到**同一组**空串只撞 DeepEqual 跳过，键不会自己回来 | 改成 `until a read answers a DIFFERENT set` ＋补"同一组空串只会撞上那一支跳过" |
| ③ `(an empty [hotkey] slot means "unset", not "turned off")` | 以行为语义口吻触到 `Q-82` 未拍那条线（＝把没拍的裁决写成契约） | 改成只描述今天这一支怎么处理空串，并具名 `Q-82` 仍在机主清单、这一枚闭包**不表态**；⛔ 裁决句、⛔ `none`/`off` 归一化 |

顺带把"注册端怎么读空串"写成事实（`hotkey_windows.go:494-496` 逐字 `HotkeyDisabled // disabled by config: never attempted`），
⛔ 用 `internal/ball has no all-empty guard` 那类"全仓不存在"的否定式承诺。

## 2. ★最硬那把尺：`⛔ 一行可执行代码变动`（件 `logs/00-diff-comment-only.txt`）

- 尺＝`git show -U0 1d6b0cf9 -- cmd/wisp/resident_windows.go`，射程＝**只有这一枚文件**，含注释行。
- `git diff --numstat`＝`11 5 cmd/wisp/resident_windows.go`（11 增 5 删＝**16 枚变动行**；`+++`/`---` 两枚文件头另计）。
- 逐行首三字统计＝`+//` **11** 枚、`-//` **5** 枚。
- 反尺＝`^[+-]` 开头而**不**以 `+//`/`-//` 开头的变动行＝**0 枚**。
- 函数体一字未动：`func hotkeyReloadSource296` 现在在 `:51`，体内 `:52-62`（含两枚 `return`，`:57` 与 `:61`）都不在 diff 里。
- 行数尺（⛔ 单枚字节的函数，只作对照）：文件总行 782 → 788；`^//` 开头行 212 → **218**（净增 6＝与 numstat `11−5` 逐字相等）；
  改后 CR 字节＝**0**（新写字节＝LF）。⇒ 注释净增 6 行 ⇄ 可执行字节变动 0 枚，两把尺相互咬合。

★另有一枚**函数体内**的注释（`:59-60`）逐字写着 `the three slots the owner's file leaves empty are "unset", and / the bridge would read them as "disabled by config"`：
它**不在 `AC#5` 那三处锚里**（三处锚＝`hotkeyReloadSource296` 上方文档块），本腿**一字未动**。
⚠ 具名请编排者裁：`"disabled by config"` 那半句说的是注册端今天的读法（`hotkey_windows.go:495` 逐字 `disabled by config: never attempted`）＝事实；
`are "unset"` 那半句是对**机主那份文件**的陈述（票面现量节已有他的四行逐字）＝本票立票依据，也⛔ 是我要扩的射程。若你裁它也该改口径，那是**第二格**、归你再派。


## 3. 我自己现读的 `Check()` 两条出口凭据（⛔ 动 `internal/ball/**`）

尺＝`git show HEAD:internal/ball/hotkey_reload.go | Select-String -Pattern 'r\.src == nil|DeepEqual\(r\.applied'`（HEAD＝`887dd356`，行号是快照）：

- `internal/ball/hotkey_reload.go:87` 逐字 `	if r.src == nil {`（`:88` 逐字 `		return false, r.Report()`）＝**"保持当前绑定"那一条**。
- `internal/ball/hotkey_reload.go:93` 逐字 `	if r.exists && reflect.DeepEqual(r.applied, cfg) {`（`:94` 逐字 `		return false, r.last // nothing changed: do not touch Win32`）。
- 全文件反尺＝`grep '== ""|HotkeyConfig\{\}|all-empty'` ⇒ `hotkey_reload.go` **0 命中**＝确无 all-empty 守卫（这一条我改成"两条跳过都不是 all-empty 守卫"，⛔ 再写成不存在式承诺）。
- 注册端读空串＝`internal/ball/hotkey_windows.go:494` 逐字 `		case b.Binding == "":` / `:495` `			b.Status = HotkeyDisabled // disabled by config: never attempted`。
- `Q-82` 在场尺＝`docs/reports/pending-and-issues.md:14927`（`## A788` 节里那三栏，含"不做"栏）＋代码侧唯一命中＝`cmd/wisp/resident_hotkey_296_windows_test.go:48`；⛔ 我改这两处任何一字。

## 4. `-run 'Test296|Test258'` 改前⇄改后读数（四数带尺名，件 `logs/20`–`logs/23`）

| 发 | `RUN=` | `PASS_TOP=` | `FAIL_TOP=` | `SKIP=` | rc | winsec seed（当发现跑） |
|---|---|---|---|---|---|---|
| 改前 1 | 13 | 8 | 3 | 2 | 1 | `07300bff…` |
| 改前 2 | 13 | 8 | 3 | 2 | 1 | `07300bff…` |
| 改后 1 | 13 | 8 | 3 | 2 | 1 | `07300bff…` |
| 改后 2 | 13 | 8 | 3 | 2 | 1 | `07300bff…` |

逐名名册四发**逐字一致**（13 枚顶层）＝我这枚注释**既没把包打红、也没把它打绿**；⛔ 放宽任何断言、⛔ 删用例、⛔ `t.Skip`。
四发的 `winsec.go` seed 同串＝`07300bfffc3e5b14c017029d6f84fb5a7ede3db8`（可比性凭据；`342-r2` 在写 `internal/winsec/**`，本腿**零字节**、四发里也没出现过任何点名 winsec 的红）。

★`rc=1` 那一族全部是同一枚**台面成因**，具名上报（⛔ 判自己失败、⛔ 据此动别人的字节）：
本机起手时就有一枚**活着的 `build\wisp.exe`（PID 21892，机主自己的常驻进程）**占着 `Ctrl+Alt+Q`／`Ctrl+Alt+M`／`Ctrl+Alt+P`
（尺＝`Get-CimInstance Win32_Process -Filter "Name='wisp.exe'"`，`CommandLine` 逐字 `"D:\work\workspace\projects plans\Wisp\build\wisp.exe"`）。
⇒ 三枚红逐字：`resident_hotkey_258_windows_test.go:361`（`the old binding did not survive a no-rebind host (VK=0x0): the premise moved`）／
`:448`＋`:453`（`live slots after the rebind = 1, want 2`）／`resident_hotkey_v1probe_test.go:108`（`boot verdict lost the live-count shape`）；
⇒ 两枚弃权逐字（`resident_hotkey_296_windows_test.go:172`/`:213`）：`SKIP-LOUD: boot could not register the three default combinations on this host (summon=occupied by another app …)`。
＝票面 `AC#4` 起手前提"⛔ 有活着的 wisp.exe 会把逐名名册搅红"今天不满足。**我没有停掉机主的进程**（那枚是他的常驻实例，本票一票到底⛔ 碰他的机器状态）。
两发改前串行、安静复量，⛔ 与 `./internal/winsec/` 并发、⛔ 跑 `cmd/wisp` 整包。

★台面自纠两处具名（第一发是假红，⛔ 计入判据）：
① 我把整条 PATH 过 `cygpath -wp` ⇒ `gcc` 所在目录（`/mingw64/bin`，真身 `E:\work\base\msys64\mingw64\bin`）掉出 Windows 形 PATH ⇒ `go env CGO_ENABLED` 由 **1 掉到 0** ⇒ `sherpa-onnx-go-windows` 被 build constraint 全排除 ⇒ `setup failed`＋`RUN=0`＋`rc=127` 那发连 `grep` 都找不到＝**零枚 `--- FAIL` 的整包死**。
⇒ 定式补一条给编排者：`cygpath -wp` 那条形只在**仓外导出树跑 Windows `.exe` 取 DLL** 时用；共享工作树里跑 `go test` 用公共块点名的 POSIX 前缀形（`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`）＋绝对路径 `go.exe`，CGO 保持 1、13 枚 `RUN` 到齐。
② `pwsh` 的 `cat-file blob > 仓外路径` 在目标目录不存在时报 `Could not find a part of the path`＝落盘前先 `mkdir -p`。

## 5. 门禁四把＋并排两把格式尺（逐名，件 `logs/30`–`logs/34`；⛔ 转述，每件 rc 由当发现取）

| 门禁 | rc | 读数 |
|---|---|---|
| `GOFLAGS= go build ./...` | **0** | 输出 0 字节；正控 `go version go1.27.1 windows/amd64` |
| `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`；ban#8 `cmd/` 143 枚、`internal/` 602 枚 |
| `gofmt -l cmd/wisp`（工作树侧） | 0 | 脏 **1** 枚＝`cmd\wisp\models.go`（已知 CRLF 存量：工作树 CR 334⇄HEAD blob CR 0）⇒ ⛔ 顺手洗 |
| `gofumpt -l cmd/wisp`（工作树侧，`/d/work/base/gopath/bin/gofumpt.exe`） | 0 | 同一枚 `models.go` |
| 裸 `gofumpt` | **127** | 具名 `command not found`＝⛔ 读成"跳过" |
| 判据尺＝`git archive HEAD cmd/wisp`→仓外→同两把尺 | 0／0 | **0 枚**（改前后 HEAD 各一次）；正控＝往导出树丢一枚故意脏 `.go` ⇒ 两把尺各命中 **1** 枚、摘掉后回 **0**；导出树 `.go` 计数 **139**（⛔ 0＝读数有效） |
| `go vet ./cmd/wisp/` | **0** | 输出 0 字节 |
| `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`（roster＝tree） |

跑门前后一律 `unset GOTMPDIR TMPDIR`（公共块 `A1027` 第⑤节）。本腿新写字节＝LF（该文件工作树 CR 字节 **0**、`git hash-object` 与 `HEAD` blob 同为 `37915a5a…`＝改前净）。

## 6. 禁区执行账（一律⛔）

⛔ `internal/ball/**`（两条出口是事实，我没动一字）；⛔ `internal/config/schema.go` 的 `[hotkey]` 四枚默认值；
⛔ 读或改机主 `config.toml`（本腿全程**没打开过**那个路径）；⛔ `none`/`off` 归一化（`Q-82` 未拍＝不做）；
⛔ 放宽票 258 既有断言（含 `resident_hotkey_258_windows_test.go:153` 那句"源可以返回全空"＝一字未动）；
⛔ 合并 `hotCfg258`⇄`hotReload258`；⛔ 翻任何框（`AC#5` 归非实现者裁）；⛔ `frontend/**`／`design/**`／三枚冻结件／golden／`thresholds.go`／D43 表；⛔ push。
票 296 框终量尺（我自己现跑，`grep -cE '^- \[[ x]\]'` 那把）＝总 **6** 格／已勾 **4**／未勾 **2**（`AC#2`、`AC#5`）＝我**一枚没翻**。
本腿越界尺＝`git show --name-only 1d6b0cf9` 只含 `cmd/wisp/resident_windows.go` 一枚。
