# 245-wording 证据：`wisp -h` 正文里 cancel 槽那句改成说清两半（裸键借、组合常驻）

- **票**：`.scratch/wisp/issues/245-the-bare-esc-cancel-default-is-a-global-hotkey-so-the-resident-process-steals-esc-from-every-other-app.md`（无 `-done`，`Status` 未动）
- **腿**：`245-wording`（文案写腿；**不是裁决腿**，AC 框一枚未翻）
- **起手 HEAD**：`e1cc304e0`（`dev`）。落盘时 HEAD 已是 `b062d52a7`（另一条腿推进）；`cmd/wisp/main.go` 在 `e1cc304e0` 之后无新改动（`git diff --numstat -- cmd/wisp/main.go` = `7 4`，只含本腿这段）。
- **写面（逐名）**：`cmd/wisp/main.go`（**仅 `usage` 那个原始字符串里的 cancel 槽句子**）／`docs/evidence/s1/245-usage-wording.md`（本件，新）／`.scratch/wisp/probes/245/wording/**`（台件）／票 245 Progress log（只追加一行）
- **产品行为零改动**：`internal/**` 全树零字节；没有新增/修改任何测试、判定、阈值、golden、规格文字。本腿只把一个字符串字面量改成与既有行为一致的说法。

## §1 为什么这句必须改（现读代码，不转抄自述）

| 事实 | 现读位置 | 读数 |
|---|---|---|
| 稳态注册集的分叉点 | `internal/ball/hotkey_windows.go:577` | `if p.id == hkCancel && !cancelHeldWhileIdle(bindings[i]) {` ⇒ 裸键走 `cancelIdleLine`（`HotkeyStandby`，从不注册）并 `continue`；带修饰键落到 `:587 registerHeldBinding`（**常绑**） |
| 「有没有真修饰键」的判据 | 同文件 `:474 cancelHeldWhileIdle` | `return acc.Mods&realModifiers != 0`（`MOD_NOREPEAT` 不算修饰键） |
| 换算成用户能看见的键 | `internal/ball` 的 `DefaultHotkeys()` | cancel 默认值仍是裸 `Esc`（本次未改） |
| 归还不是「删掉就完」 | 同文件 `:564 cancelIdleLineWith`，由 `ball_windows.go:1032` 的归还路径调用 | 带修饰键的 cancel 在会话结束后**重新绑回**；裸键支仍只 drop、回 standby |

⇒ 改前那句 `bound only while a confirmation card is waiting and handed back when it leaves` 对**裸键**成立、对**组合键**不成立（`Ctrl+Alt+X` 配好后是常驻）。改后两半都写明。

## §2 改前 / 改后（逐字，宽度守住）

改前（`git show 9e84ae42e:cmd/wisp/main.go`，即本腿起手时 HEAD 里的同一段）：

```
  wisp             GUI resident process: boots the runtime skeleton, hosts the
                   floating ball window, its tray icon and its three always-on
                   global hot keys in this same process (summon, mute, panel;
                   the cancel slot is the fourth key, bound only while a
                   confirmation card is waiting and handed back when it leaves,
                   so an idle wisp holds no bare key against the rest of the
                   desktop), shows confirmation cards on that ball, and takes
                   tasks from the operator's console when this process really
```

改后（工作树 `cmd/wisp/main.go:28-40`，逐字）：

```
  wisp             GUI resident process: boots the runtime skeleton, hosts the
                   floating ball window, its tray icon and its three always-on
                   global hot keys in this same process (summon, mute, panel;
                   the cancel slot is the fourth key: a bare binding such as the
                   default Esc is borrowed only while a confirmation card is
                   waiting and handed back when it leaves, so an idle wisp holds
                   no bare key against the rest of the desktop, while a binding
                   carrying a modifier such as Ctrl+Alt+X stays bound all the
                   time because a combination takes no bare key from any other
                   program), shows confirmation cards on that ball, and takes
                   tasks from the operator's console when this process really
                   owns one (ticket 246 AC#7: no interactive console means the
```

要点（逐条对派单）：

1. **裸键半边**（默认 `Esc`）：`is borrowed only while a confirmation card is waiting and handed back when it leaves, so an idle wisp holds no bare key against the rest of the desktop` — 保留原句事实，并显式点名 `bare binding / default Esc`。
2. **组合键半边**（如 `Ctrl+Alt+X`）：`stays bound all the time because a combination takes no bare key from any other program` — 常驻，并说清它不吞别人的裸键。
3. **数没动**：`its three always-on global hot keys` 逐字保留，仍是 **three**；没有把 always-on 的枚数写成 four。`the cancel slot is the fourth key` 是改前就有的「槽位是第四枚」说法，**本次未改这半句**（被改的只是它后面那个假条件）。
4. **宽度**：这 7 行长度依次 **80 / 76 / 80 / 79 / 77 / 78 / 77**，续行缩进 19 空格，与同段其余行（72–80）一致；`usageBare133` 钉的行首 `^  wisp {2,}\S` 那一行未碰。
5. `git diff --numstat -- cmd/wisp/main.go` = `7	4	cmd/wisp/main.go`（+7/−4，改的只有这 4 行、新增的下文 3 行属同一句的换行）。
6. `cmd/wisp/main.go` SHA256：改前（HEAD blob，`git cat-file` 逐字节）`89F2FE0557D6CEE7AD9C2D9414DA25DD9080CE4C39932CA27BFBD418F450B6C7` → 改后 `351F4DFE08D049FCB4F8B271CADF1870CE7E5FE8DAC33F4AD8D2DC530A86F784`。

## §3 真跑 `wisp -h`（暂存 exe，不是系统里的旧 wisp）

构建（`GOFLAGS=""`，cgo＋sherpa 走仓里既有方式：mingw gcc 与 pinned DLL 目录上 PATH，链接库由 Go 绑定模块自带）：

```
$env:GOFLAGS=""
$env:PATH="E:\work\base\msys64\mingw64\bin;<repo>\third_party\sherpa-onnx;"+$env:PATH
go build -o .scratch/wisp/probes/245/wording/bin/wisp.exe ./cmd/wisp
=> BUILD-RC=0（15s）

E:\work\base\msys64\usr\bin\bash.exe scripts/stage-sherpa-dlls.sh .scratch/wisp/probes/245/wording/bin
=> STAGE-RC=0（onnxruntime.dll / sherpa-onnx-c-api.dll / sherpa-onnx-cxx-api.dll 三枚并置）
```

新鲜度（防「拿旧 exe 读数」）：`cmd/wisp/main.go` mtime **14:40:27** < `wisp.exe` mtime **14:41:10** < 实跑时刻 **14:41:43**；`wisp.exe` 51,503,946 B，SHA256 `DF01A2277FBD819983B67033FC2C0C93A12AC53809A17A72769E48D77AB84B4F`。

运行（cwd = 该 bin 目录，stdout/stderr 分开落件）：

```
.\wisp.exe -h > ..\help-after.out.txt 2> ..\help-after.err.txt
=> rc=0，679 ms
   help-after.out.txt 4802 B，SHA256 839AF70270C8C165040624C1B620DE74FC95497C33AFBC6629AC7DDF051EE917
   help-after.err.txt  135 B
```

stderr 只一行（既有 winsec 安装 INFO，不属 usage，原样抄）：

```
time=2026-11-14T14:41:43.349+08:00 level=INFO msg="winsec: sealing path resolver installed" resolver=risk.c26Pipeline probes_passed=2
```

stdout 逐字（全文 68 行，与 `.scratch/wisp/probes/245/wording/help-after.out.txt` 同文）：

```
wisp - personal voice agent for Windows

Usage:
  wisp             GUI resident process: boots the runtime skeleton, hosts the
                   floating ball window, its tray icon and its three always-on
                   global hot keys in this same process (summon, mute, panel;
                   the cancel slot is the fourth key: a bare binding such as the
                   default Esc is borrowed only while a confirmation card is
                   waiting and handed back when it leaves, so an idle wisp holds
                   no bare key against the rest of the desktop, while a binding
                   carrying a modifier such as Ctrl+Alt+X stays bound all the
                   time because a combination takes no bare key from any other
                   program), shows confirmation cards on that ball, and takes
                   tasks from the operator's console when this process really
                   owns one (ticket 246 AC#7: no interactive console means the
                   task entry says so out loud and the ball keeps running;
                   Ctrl+C stops the process, and so does the tray's Exit item,
                   which asks the resident loop for the orderly D38(e) shutdown
                   instead of exiting in place (ticket 228 AC#11))
  wisp run "task"  run one task end to end through the agent loop (ticket 12:
                   the S1 text path - streamed reply + notification + exit code
                   that reflects the error class)
  wisp sessions    cross-restart conversation list (ticket 166 outlet #1,
                   wired by 364-r1a): sessions list prints the newest sessions
                   through the existing session.Manager.List - id, updated_at,
                   title
  wisp commands    the command roster (ticket 213 AC#4): prints every command
                   this machine can really execute, one row each - the name, what
                   it does, its argument shape, whether a human approval is
                   involved, and which existing door runs it. The rows come from
                   internal/command's production roster, so a capability with no
                   door cannot appear here
  wisp providers   provider catalog path (ticket 12, ruling A11): discover
                   <provider> lists /v1/models; probe <provider>/<model> runs
                   the capability probe suite into provider_health
  wisp secret      credential entry (ticket 63): set/get/list/unset a DPAPI
                   blob without the key ever entering argv, a log line or chat
  wisp config      config.toml keys from the console (ticket 65): get/set/list
                   over the keys this CLI owns - today ball.style (glass|
                   jarvis, the appearance switch 默认玻璃球), written through
                   the same guarded setter the settings page uses; the receipt
                   names the tier honestly (a CLI write lands in the file, a
                   running wisp picks the new style up at next boot)
  wisp models      signed model store (C29, ticket 121 AC#2): list the verified
                   manifest, re-check an installed model, or make one available
                   through the Downloading hand-off and its re-verification
  wisp doctor      toolchain and native-DLL self-check, prints PASS/FAIL
  wisp panel-assets  embedded panel bundle (ticket 77): -manifest lists what the
                   binary carries, -render <path> writes the bytes the WebView2
                   host would serve (proof the UI needs no node and no network)
  wisp panel-inbound  inbound panel dispatch leg (ticket 33 AC#9): reads one raw
                   composer envelope per stdin line and hands it to the inbound
                   router; -data names the data root (this leg never resolves one)
  wisp slo         SLO sampling driver (ticket 08): one state per run, JSON
                   verdict; driven by scripts/slo-check.ps1
  wisp diagnostics  diagnostics bundle (ticket 45 engine, wired by the
                    deferred-disposition 20261110 #5): the total switch is
                    [privacy] diagnostics_opt_in in config.toml - false refuses
                    the whole command (fail-closed, no preview, no log read),
                    true runs preview (what the zip WOULD carry, nothing
                    written) and export (same preview, then the bundle zip;
                    -o names the output, -slo-snapshot attaches a sampler JSON).
                    Transcripts and watchdog snapshots stay default-excluded
                    with per-item consent recorded in the manifest
  wisp version     print version information
  wisp help        show this help

Environment: WISP_ENV in {prod|dev|test}, default from build (SPEC-03 §5).
```

## §4 门禁（本腿现跑，逐条读数）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `GOFLAGS="" go build ./cmd/wisp`（PATH 同 §3） | **rc=0**（输出 0 字节） |
| 暂存 exe 构建 | `go build -o .scratch/wisp/probes/245/wording/bin/wisp.exe ./cmd/wisp` | **rc=0**，51,503,946 B，SHA256 `DF01A227…84B4F` |
| gofumpt | `D:\work\base\gopath\bin\gofumpt.exe -l cmd/wisp` | **空输出，rc=0** |
| d22scan | `sh scripts/d22scan.sh`（msys bash） | **rc=0，148 s**；`runtests.sh: OK - packages=[./...] top-level: PASS=45 FAIL=0 SKIP=0, === RUN=93`；末行 `d22scan: clean - no D22 ban violations`；`ban #8 cmd/ examined 280 Go files`（emoji 扫描含改动后的 `cmd/wisp/main.go`） |
| 真跑 `wisp -h` | §3 | **rc=0**，stdout 逐字含改后那句 |

d22scan 原始日志：`.scratch/wisp/probes/245/wording/d22scan.log`（27,532 B）。

## §5 边界自查（只给名册与读数）

```
$ git status --porcelain cmd/wisp
 M cmd/wisp/main.go
 M cmd/wisp/panel_profile_leaf_shape_windows_test.go
```

- 本腿自己的 `git status --porcelain` 相关项：` M cmd/wisp/main.go` ＋ 未跟踪的 `docs/evidence/s1/245-usage-wording.md`、`.scratch/wisp/probes/245/wording/`、票 245 Progress log（追加）。
- **禁区逐名 0**：`internal/**` **0**（包括 `internal/ball/hotkey_windows.go`、`internal/ball/ball_windows.go`；本腿只读不写）、`scripts/**` **0**、`frontend/**` **0**、`design/**` **0**、`docs/PLAN.md`／`docs/specs/**` **0**、`internal/observe/thresholds.go` **0**、golden **0**、`tools/d22scan/allowlist.txt` **0**。
- 未新增 `DEFERRED(D-xx)` 标记；未改 `Status`；未 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`；未 push；未跑整包测试。
- `git status --porcelain cmd/wisp` 里的 `panel_profile_leaf_shape_windows_test.go` 是**别人未完成的改动**（本腿进场前已在，未 staged、未提交）。

## §6 没证明什么 / 残差（明写）

1. **本腿不是裁决**：AC#6① 框的翻与不翻归编排者；本件只交「文案现在与行为一致」＋「真跑 `-h` 的原文」。
2. **没跑 winlive**：`Ctrl+Alt+X` 常驻行为本身已由 `d44e89216` 的确定性用例钉住，本腿复跑的是同一句代码事实的**文案**，没有在真桌面上再注册一次键（写腿禁真机）。
3. **没新增测试**：派单要求只改文案；没有新增词面型仪器（这本身是票 245 的既有约束）。
4. **示例键不是默认值**：`Ctrl+Alt+X` 只是「带修饰键」的例子；产品默认 cancel 仍是裸 `Esc`（`DefaultHotkeys()` 未动），所以第一半句才是默认行为。
5. **起手 HEAD 与会话开始时的读数**：本腿开始读仓时 HEAD 是 `e1cc304e0`；收尾时 `b062d52a7`，两者之间 `cmd/wisp/main.go` 无改动（§2 的 diff 只含本段 4 行）。
