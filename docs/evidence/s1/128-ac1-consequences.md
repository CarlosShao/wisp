# 票 128 AC#1 — `%APPDATA%` 缺失时数据根回落到 CWD：把**后果**量成读数

**测量会话：** `q7` · **锚定 sha：** `2620836`（dev）· **日期：** 2026-09-23（首条读数 09:37:40 +08:00）
**范围：** 仅 AC#1（四样落点在两个 CWD 下的实际路径 + 第二遍读到什么 + `icacls` 归属 + GUI 腿拒绝独立复现）。
**AC#2/AC#3/AC#4 不在本票段内**（AC#2 是语义裁定，见票面 Progress log 的 `next=`）。

## 0. 仪器与安全前提

- **仓外临时根：** `/tmp/wisp128-q7/` ⇒ Windows 实路径 `C:\Users\swq\AppData\Local\Temp\wisp128-q7\`。
  两个 CWD 是它的两个子目录（`cwd-A`、`cwd-B`），全部写点被限制在这棵树里。
- **owner 真实数据目录：未写入。** 开工前定位到两枚，且**开工前两者均为空目录**（`find -type f` 计数 = 0、`du -sb` = 0）：
  - `C:\Users\swq\AppData\Roaming\wisp`（目录 mtime 2026-09-19T14:49:20）
  - `C:\Users\swq\AppData\Roaming\wisp-dev`（目录 mtime 2026-09-20T07:06:25）
  ⇒ 基线是"零文件"，所以任何落在里面的文件都是**可判定的写入**；收尾复核见 §6。
- **不用 docker**（本机 Git Bash 下 `docker -v "C:\…"` 会静默挂空且 rc=0 = 假绿）。
- **不取时序/内存读数**：本机 `slo-full` 就跑在 self-hosted runner 上（A103），本段无需排队。

### 0.1 一处必须写明的偏离（编排者指令 vs 被量代码的形状）

派单要求"测量全程 `WISP_ENV=test` + 显式临时根"。**照做就量不到 AC#1**：
`resolveDataDir` 在 `env == "test"` 时提前 `return proc.TestDataDir()`，
根本走不到 `base, err := os.UserConfigDir(); if err != nil { base = "." }` 那三行。
所以本段主腿用 **`WISP_ENV=dev`**（落点 `wisp-dev`，与票面"已量"一节记的 `wisp-dev\logs` 同形），
`test` 只作对照腿（§5）。
隔离手段由"env=test"换成"**CWD 本身就在仓外临时根里**"——回落出来的 `wisp-dev\` 于是落在
`C:\Users\swq\AppData\Local\Temp\wisp128-q7\cwd-A\wisp-dev\`，与 owner 的
`C:\Users\swq\AppData\Roaming\wisp-dev\` 同名不同树；§6 的复核就是证这一点。

### 0.2 源码事实（读 `2620836`，非推断）

- `cmd/wisp/doctor.go:233` `resolveDataDir(env string) string`：portable 覆盖 → `env=="test"` → `proc.TestDataDir()` →
  否则 `base, err := os.UserConfigDir() // %APPDATA%`；`if err != nil { base = "." }`；
  再 `base = proc.SealableRoot(base)`；`env=="dev"` ⇒ `filepath.Join(base,"wisp-dev")`。
- `resolveDataDir` 的**四个生产调用点**：`cmd/wisp/doctor.go:101`、`cmd/wisp/doctor.go:276`（`dataDirForDisplay`）、
  `cmd/wisp/models.go:110`、`cmd/wisp/providers.go:81`、`cmd/wisp/run.go:155`。
  ⇒ run 腿的日志/`config.toml`/`memory.db` 都跟着它。
- GUI/boot 那条腿走 `internal/proc/envfork.go:236` `DefaultLayout`：
  `if err != nil { return Layout{}, fmt.Errorf("proc: user config dir: %w", err) }` ⇒ **拒绝，不回落**。
- `cmd/wisp/secret.go:123` `resolveSecretLayout` 同样**直接返回错误**（`wisp secret: user config dir: %w`），不回落。
- `SealableRoot(".")` 把 `.` 解析成**绝对路径**（`filepath.EvalSymlinks` 那半），所以日志上看到的不是 `.` 而是 CWD 的全名。

## 1. run 腿 · CWD-A（第一遍）

**仪器补正：** `go build ./cmd/wisp` 出来的 exe 在**加载期**就要 `sherpa-onnx-c-api.dll`
（首跑 rc=127 `error while loading shared libraries`，零输出＝假绿形状），
故把 `third_party/sherpa-onnx/*.dll` 三枚复制到临时根与 exe 同目录后才可执行。
**这条是本票的一个仪器坑：加载期失败会让"落点为空"被误读成"没有回落"。**

调用形（两腿一致，只有 CWD 不同）：
`cd <CWD> && env -u APPDATA -u WISP_TEST_DATA_DIR WISP_ENV=dev <根>/wisp.exe run "AC1 探针"`

### 1.1 `wisp doctor`（同形，第一发）

- rc=1，`[PASS] data dir writable (dev)    wisp-dev` —— **打印的是相对拼写**，不是绝对树。
  实测该发在磁盘上创建了 `C:\Users\swq\AppData\Local\Temp\wisp128-q7\cwd-A\wisp-dev\`（空目录）。
- 另两枚 `[FAIL]`（sherpa/onnxruntime 版本 pin）与 `commit=unknown` 是本机源码构建的正常形状，与本票无关，如实登记不掩盖。

### 1.2 `wisp run`（A，无 config.toml）：rc=2

```
time=2026-09-23T09:39:37+08:00 level=INFO msg="wisp: persistent log sink installed" dir=wisp-dev\logs min_level=info
[audit] perm: MODE-READ-FAILED path="wisp-dev\\config.toml" err=config: config.toml read: open wisp-dev\config.toml: The system cannot find the file specified. mode=ask_every_step origin=startup result=fail-closed
wisp run: 配置未就绪（Unconfigured）：config: config.toml read: open wisp-dev\config.toml: The system cannot find the file specified.
```

**A 树实际落点（`find` 全量，size 字节）：**

| # | 落点 | 实际绝对路径 | 状态 |
|---|---|---|---|
| 1 | 日志 | `C:\Users\swq\AppData\Local\Temp\wisp128-q7\cwd-A\wisp-dev\logs\wisp-20260923-001.jsonl` | 写入了，1138 字节（两发之后） |
| 2 | `config.toml` | `C:\Users\swq\AppData\Local\Temp\wisp128-q7\cwd-A\wisp-dev\config.toml` | 读点在此；初始不存在，探针文件见 §1.3 |
| 3 | DPAPI 私钥存储 | `C:\Users\swq\AppData\Local\Temp\wisp128-q7\cwd-A\wisp-dev\secrets\` | **目录被创建并被 winsec 封条过**（`sealing path resolver installed resolver=risk.c26Pipeline probes_passed=2`）；blob 未写，见 §3.2 |
| 4 | `memory.db`（实名 `wisp.db`） | 应为 `…\cwd-A\wisp-dev\wisp.db` | **本发未创建**——rc=2 在 `memory.Open` 之前就 return，见 §3.3 |

### 1.3 给 A 树放一枚**可指名**的 `config.toml`

内容 `schema_version = 987654`（故意大于本机理解的 2，让**读到的文件名+值**出现在错误里）：

```
[audit] perm: MODE-READ-FAILED path="wisp-dev\\config.toml" err=config: config.toml: schema_version 987654 was written by a newer build (this build understands 2); upgrade Wisp or restore a backup
wisp run: 配置未就绪（Unconfigured）：config: config.toml: schema_version 987654 was written by a newer build (this build understands 2)
```

⇒ **A 确实读的是 `cwd-A\wisp-dev\config.toml`**（读数里带着我写进去的那个整数，不是巧合可撞出来的）。

## 2. run 腿 · CWD-B（第二遍，另一棵树）

同一条命令、同一个 exe、同一份环境（只 `cd` 到 `cwd-B`）：rc=2

```
time=2026-09-23T09:40:22+08:00 level=INFO msg="wisp: persistent log sink installed" dir=wisp-dev\logs min_level=info
[audit] perm: MODE-READ-FAILED path="wisp-dev\\config.toml" err=config: config.toml read: open wisp-dev\config.toml: The system cannot find the file specified. mode=ask_every_step origin=startup result=fail-closed
wisp run: 配置未就绪（Unconfigured）：config: config.toml read: open wisp-dev\config.toml: The system cannot find the file specified.
```

**B 树落点（`find` 全量）：**

```
cwd-B/wisp-dev                       d
cwd-B/wisp-dev/logs                  d
cwd-B/wisp-dev/logs/wisp-20260923-001.jsonl   f 551
cwd-B/wisp-dev/secrets               d
```

⇒ B 独立长出一棵 `wisp-dev\`：自己的日志（551 字节，与 A 的 1138 不同文件、各自从 `-001` 编号）、
自己的 `secrets\`，**没有** A 的那份 `config.toml`。

## 3. 四样落点汇总 + "第二遍读到的是不是同一棵树"

### 3.1 结论：**不是同一棵树。** 第二遍读到的是**空配置**，且是"读不到文件"而不是"读到默认值"

判据是 §1.3 与 §2 的**差分**：同一个整数 `987654` 在 A 树被点名、在 B 树变成
`The system cannot find the file specified`。这不是"看起来会搬家"，是**同形下两份状态互相看不见**。

| 落点 | CWD=A 的实际路径 | CWD=B 的实际路径 | 是否同一棵树 |
|---|---|---|---|
| 日志 | `…\wisp128-q7\cwd-A\wisp-dev\logs\wisp-20260923-001.jsonl`（1138 B） | `…\wisp128-q7\cwd-B\wisp-dev\logs\wisp-20260923-001.jsonl`（551 B） | **否**，两份，文件名还**同名**（`-001`） |
| `config.toml` | `…\cwd-A\wisp-dev\config.toml`（存在，被读，值 987654 被点名） | `…\cwd-B\wisp-dev\config.toml`（**不存在**，报错 exit 2） | **否**，B 读到的是"无配置" |
| DPAPI 私钥存储 | `…\cwd-A\wisp-dev\secrets\`（已创建+封条） | `…\cwd-B\wisp-dev\secrets\`（已创建+封条） | **否**，两枚互不相干的 DPAPI 目录 |
| `memory.db`/`wisp.db` | 未创建（rc=2 早退，见 §3.3） | 未创建 | **未验证**，只验证了"读点拼写"随 CWD 走 |

派生的三条后果（都可从上面表格直接读出来，不需要再猜）：
1. **配置丢失式失败**：换目录 ⇒ 凭据引用、档位、模型链全部"看不见"，用户看到的是一个刚装好的空 app。
2. **日志分裂**：同一台机器上按 CWD 长出 N 份 `wisp-dev\logs`，且**文件名相同**（日期+序号），
   事后取证无法靠文件名判断哪一份是哪棵树写的，只能看绝对路径。
3. **凭据面分裂 + 仍然私有**：`secrets\` 跟着搬家，票面"已量"那一句（"同迁、照样被封"）在两棵树下都成立。

### 3.2 DPAPI 腿的**读点**与 run 腿不同形（本票没量到的那一半，如实标）

- run 腿（`cmd/wisp/run.go:222` `secret.NewStore(s.dataDir)`）拿的是 `resolveDataDir` 的结果 ⇒ 会落到 `<CWD>\wisp-dev\secrets`。
- `wisp secret` CLI 子命令（`cmd/wisp/secret.go:123` `resolveSecretLayout`）在 `APPDATA` 未设时**直接返回错误、不回落**。
  实测（A 树、`WISP_ENV=dev`、`-u APPDATA`）：（读数见 §5）
- ⇒ 因此**同一次 CWD 回落里，`secrets\` 目录被创建了两枚（A/B 各一）但都是空的**：本段**未能**把一枚真 DPAPI blob 写进
  `<CWD>\wisp-dev\secrets\`。**这条记为未验证**，不当成结论的地基。

### 3.3 `memory.db` 未创建的原因（读源码，非猜）

`runTextTask` → `assembleRuntime` 的顺序里，配置未就绪即 `return code`（rc=2），
而 `wisp.db` 由 `internal/memory/open.go:39 dbFileName = "wisp.db"` 在更晚的存储打开步骤创建。
⇒ 本段对 `memory.db` 只量到"**读点拼写随 CWD 走**"（日志与 `secrets\` 已证同一 `dataDir` 派生），
**没量到**它真的落在 CWD 树里。**未验证。**

## 4. `icacls` 归属读数（对照票 95 的私有目录纪律）

同一棵树（`cwd-A\wisp-dev\`）上四枚落点的 DACL 原文，逐条对读：

| 对象 | `icacls` 读数（本机 `swq`） | 性质 |
|---|---|---|
| `…\cwd-A\wisp-dev\secrets` | `SYSTEM:(F)` + `SYSTEM:(OI)(CI)(IO)(F)` + `Administrators:(F)` + `Administrators:(OI)(CI)(IO)(F)` + `swq:(F)` + `swq:(OI)(CI)(IO)(F)` | **每一条都没有 `(I)` 标记** ⇒ 继承被切断、DACL 是 winsec 自己写上去的**显式私有** ACL |
| `…\cwd-A\wisp-dev` | `SYSTEM:(I)(OI)(CI)(F)`、`Administrators:(I)(OI)(CI)(F)`、`swq:(I)(OI)(CI)(F)` | 全部 `(I)` ⇒ **纯继承自启动目录** |
| `…\cwd-A\wisp-dev\logs` | 同上，三条 `(I)(OI)(CI)(F)` | 纯继承 |
| `…\cwd-A\wisp-dev\config.toml` | 三条 `(I)(F)` | 纯继承 |
| `…\cwd-A\wisp-dev\logs\wisp-20260923-001.jsonl` | 三条 `(I)(F)` | 纯继承 |

**对照控制组（同机、只读）：**

- owner 真实数据根 `C:\Users\swq\AppData\Roaming\wisp-dev`：`SYSTEM/Administrators/swq` 三条 `(I)(OI)(CI)(F)` ——
  **与回落树的 `wisp-dev\` 一模一样**。⇒ 本机这次测量的 CWD 恰好落在**同一个用户自己的 Temp 下**，
  继承源与 `%APPDATA%` 同形，所以**归属没有被这次测量放大**。这是"这次没出事"，不是"不会出事"。
- 继承源可以是多坏的样子（**只读 `icacls C:\Users\Public`，未在该树下写入**）：
  `BUILTIN\Administrators:(OI)(CI)(F)`、`CREATOR OWNER:(OI)(CI)(IO)(F)`、`SYSTEM:(OI)(CI)(F)`、
  `NT AUTHORITY\INTERACTIVE:(OI)(CI)(IO)(M,DC)`、`NT AUTHORITY\INTERACTIVE:(RX,WD,AD)`、
  `NT AUTHORITY\SERVICE:(RX,WD,AD)`、`NT AUTHORITY\BATCH:(RX,WD,AD)`。
  ⇒ 若启动目录在这样的树下，回落出来的 `logs`/`config.toml` 会**继承**任何交互式登录的 RX+写子项，
  而 `wisp-dev\` 本身仍会带上 `INTERACTIVE` 的 `(WD,AD)`。

**一句话结论：** 票 95 的私有目录纪律在 CWD 那种树上**只剩四分之一**——
`secrets\` 因 winsec 自带 DACL 而**照样私有**（与票面"已量：同迁、照样被封"一致，两棵树各证一次）；
`wisp-dev\`、`logs\`、`config.toml`、日志文件这四枚**一律跟随启动目录的继承**，
Wisp 不写自己的 DACL，也就没有任何一条断言说"这棵树归谁"（memory 里那条 winsec 裁定
「winsec 不管这棵树归谁」在这里正是它的代价面）。
**未实测的一格**：没有把 CWD 换到 `C:\Users\Public` 之下去**实际写一次**再取 DACL（那要在共享目录里落文件），
所以上面第二条是"继承源的读数 + 源码里无人改写 DACL"的合取，不是落点实测 ⇒ 标**未验证**。

## 5. GUI 腿在同一形下的拒绝（独立复现）

同一枚 exe、同一个 `cwd-A`、同一份 `-u APPDATA WISP_ENV=dev`：

```
$ env -u APPDATA WISP_ENV=dev wisp.exe        # 无子命令 = 常驻/GUI 腿，runResident
rc=1
wisp: boot failed: proc: user config dir: %AppData% is not defined
```

第二枚独立复现（`wisp secret list`，走 `resolveSecretLayout`）：

```
$ env -u APPDATA WISP_ENV=dev wisp.exe secret list
rc=2
wisp secret: wisp secret: user config dir: %AppData% is not defined
```

第三枚对照（run 腿，同形同环境同目录）：**不拒绝**，rc=2 但已经把
`cwd-A\wisp-dev\logs\wisp-20260923-001.jsonl` 与 `cwd-A\wisp-dev\secrets\` 写到磁盘上（§1.2）。

⇒ **三条腿三种行为**：`run` 静默搬到 CWD、常驻/GUI 腿 rc=1 拒绝、`secret` 子命令 rc=2 拒绝。
票面"已量"记的是两条腿（run 搬家 / GUI 拒绝），本段补上第三条：`wisp secret` 与 GUI 同形拒绝，
但它拒绝的是**自己那条读点**，**管不到 run 腿已经在 CWD 里创建出来的那枚 `secrets\`**。
⇒ AC#2 要裁的不一致面比票面写的宽：不是"两腿"，是"三腿"。

**`test` 环境对照腿（证明派单原本要的 `WISP_ENV=test` 量不到这一格）：**
`WISP_ENV=test` 时 `resolveDataDir` 提前 `return proc.TestDataDir()`，
`base = "."` 那一支**根本不执行**，`APPDATA` 是否设置对落点无影响 ⇒ 主腿必须用 `dev`（理由见 §0.1）。

## 6. owner 真实数据目录复核（收尾）

测量全程结束后（读数时刻 2026-09-23 09:42 +08:00）：

| 路径 | 开工前 | 收尾 | 判定 |
|---|---|---|---|
| `C:\Users\swq\AppData\Roaming\wisp` | 0 文件 / 0 子项，目录 mtime `2026-09-19 14:49:20` | 0 文件 / 0 子项，mtime `2026-09-19 14:49:20` | **未写入** |
| `C:\Users\swq\AppData\Roaming\wisp-dev` | 0 文件 / 0 子项，目录 mtime `2026-09-20 07:06:25` | 0 文件 / 0 子项，mtime `2026-09-20 07:06:25` | **未写入** |

本次全部写点（`find /tmp/wisp128-q7 -type f`，仓外）：
`cwd-A/wisp-dev/config.toml`、`cwd-A/wisp-dev/logs/wisp-20260923-001.jsonl`、
`cwd-B/wisp-dev/logs/wisp-20260923-001.jsonl`、临时根本身的 `wisp.exe` + 三枚 sherpa/onnxruntime DLL + 测量输出文件。
`C:\Users\Public` **只做了 `icacls` 读**，未创建任何文件。

## 7. 本段没量到的（明写未验证，不当任何结论的地基）

1. **`memory.db`/`wisp.db` 的实际落点未创建**：`runTextTask` 在 `config/角色/provider` 三步任一步失败即 return，
   `memory.Open` 在更晚处（`internal/memory/open.go:39 dbFileName="wisp.db"`，根是同一个 `dataDir`）。
   要让 run 腿走到存储打开需要一个能应答的 `/v1` provider（compose mock-llm），本段按派单**不用 docker** ⇒
   **只量到"它的读点拼写随 CWD 走"（与 `logs`/`secrets` 同源），未量到它在 CWD 树里真的落盘。**
2. **DPAPI blob 未真实写入 CWD 树**：`wisp secret` 那条 CLI 腿在此形下直接拒绝，写不进东西；
   run 腿创建 `secrets\` 但不写条目 ⇒ "CWD 树里的 DPAPI blob 能不能被另一个 CWD 的实例读到"**未验证**。
3. **共享目录作为 CWD 时的真实落点 ACL 未实测**（见 §4 最后一条）。
4. **prod 环境（`WISP_ENV=prod`，落点 `wisp`）未跑**：源码同一条 `base = "."`，`filepath.Join(base,"wisp")`，
   形状与 dev 只差目录名 ⇒ 未实测，登记为推断。
5. **`portable.txt` 那一支与 CWD 回落的优先级未测**（§0.2 读源码可知 portable 覆盖在前，命中即不回落）。
