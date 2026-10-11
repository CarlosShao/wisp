# 308-v1 · 60-attribution-ruling（票 §3 留给 `308-v1` 的那一句：成因归因单裁）

## 裁的对象（逐字，来自票 §编排者裁定 §3 末行）

> ⚠ 那发当时的成因归因（腿写"Go 从 Git Bash 继承无卷标 `TMPDIR`"）我⛔ 复认（Windows 档 Go⛔ 读 `TMPDIR`，见 §1）＝**归因⛔ 立、现象⛔ 否**，留给 `308-v1` 单独裁这一句。

两件事分开裁：**归因句**（那句解释）与**现象**（那发 panic 真发生过）。

## 一、归因句＝**⛔ 立**（本腿把它当"机制断言"来跑，⛔ 引用文档）

尺 1（**读产码级的源**，⛔ 二手；件 `logs/v50-goroot-tempdir.txt`，1105 B）：`D:/work/base/go/src/os/file_windows.go:284-304` ＝ 本工具链（`go version go1.27.1 windows/amd64`，`go env GOROOT`＝`D:\work\base\go`）里 `func tempDir()` 的全部实现——它只调 `syscall.GetTempPath` / `windows.GetTempPath2`（Win32 `GetTempPathW`），**整段⛔ 出现 `TMPDIR`**；对照尺＝`grep -rn 'TMPDIR' $GOROOT/src/os/` 只命中 `file_plan9.go:546` 与 `file_unix.go:391` 两支（`os/file_windows.go` **0 命中**）。
⇒ **机制层**：Windows 档 Go 的 temp 根＝Win32 语义（`TMP` → `TEMP` → 系统目录 → `USERPROFILE`），`TMPDIR` 是 **POSIX 侧** `os.TempDir()` 的变量，Windows 档⛔ 读它。

尺 2（**当场跑**，件 `logs/v14-tempgrid.txt`）：
- `env -u TMP -u TEMP TMPDIR=/tmp` → 探针看到 `env TMPDIR=[C:\Users\swq\AppData\Local\Temp]`（**MSYS 先把 `/tmp` 改写成了 Windows 形**，那一枚 env 值根本⛔ 是 POSIX 形），而 `os.TempDir()`＝**`C:\Users\swq`**＝`USERPROFILE` 兜底那一支 ⇒ **Go ⛔ 理 TMPDIR**（若它读，答案应当是 AppData\Local\Temp）。
- 基线一发 `TMP/TEMP` 都在 ⇒ `os.TempDir()`＝`C:\Users\swq\AppData\Local\Temp`（`VolumeName="C:"`）。
⇒ **归因句两路皆反证**（机制＋现跑），且现跑还多打掉一层：就算 Go 真读 `TMPDIR`，Git Bash/MSYS 递到子进程的那枚值也**早被补成带盘符的 Windows 路径**，⛔ 存在"无卷标的 POSIX 形 `TMPDIR`"这枚东西可继承。

## 二、现象＝**⛔ 否**（本腿⛔ 否认它发生过，也⛔ 说本腿复现了它）

- 原始件在盘＝`probes/72/a2/logs/t03-broken-surface-note.txt`（本腿数过＝**24 行**；逐字 `panic: runtime error: slice bounds out of range [:1] with length 0 [recovered, repanicked]`；`rc=1`／`RUN=0 FAIL=2 SKIP=0`／`PASS=0`）。
- ★本腿**没有**在 env 台面上重放它（本机⛔ 造得出可创建的无卷标 temp，见下"三"），但本腿在**输入面**重造了**同一枚崩形、同一个后果**：把 `driveShare` 的 `len(vol)==0` 守卫拔掉 ⇒ 同一句 panic 逐字出现，且整包求值从 `RUN=201` 掉到 `RUN=35`（件 `logs/teeth-mut-fullpkg.txt`，见 20 件 §a）。
  ⇒ 现象那句"**一枚无守卫切片能让整包静默少跑**"**⛔ 必"要那台坏机器才成立**——本腿在好机器上、用⛔ 依赖机器的载具把它量出来了。**这是本件对票 §3 的实质推进**：那枚矛盾（"真崩过"⇄"本机造不出"）的**杀伤力**部分已被就地可量，⛔ 再靠历史件背书。
- 原始件里那一发的崩**位**（`TestClassifySpellingInvariance`，票面 §现量 记 `anchor:77-78`）与本腿量的位（helper `:27`）是**同一枚切片表达式的两个宿主**，⛔ 混为"同一发"。

## 三、"本机造不出无卷标 temp"这枚前提＝**成立，但理由⛔ 是 A846 写的那个**

本腿自己的网格（件 `logs/v14-tempgrid.txt` bash 侧 ＋ `logs/v15-cmd-grid.txt` cmd 侧，⛔ 依赖 a1 的任何尺）：

| 写法 | Go 看到的 `TMP` | `os.TempDir()` | `MkdirTemp` 结果 | `VolumeName` |
|---|---|---|---|---|
| bash `TMPDIR=/tmp`（`-u TMP -u TEMP`） | （MSYS 改成 Win 形） | `C:\Users\swq` | 成 | `C:` |
| bash `TMP=/` | `C:\Users\swq\.qoder-cn\bin\git` | 同 | 成 | `C:` |
| bash `TMP=//tmp` | **`\\tmp`** | **`\\tmp`** | **⛔ 成**：`mkdir \\tmp\…: The network path was not found` | **`""` ← 无卷标** |
| cmd `TMP=/` | `/` | `C:\` | 成 | `C:` |
| cmd `TMP=\` | `\` | `C:\` | 成 | `C:` |
| cmd `TMP=RISKAT~1\RUNNER~1` | 原样 | 同 | 成 | `C:` |

⇒ 读法（**这条是 A846 该收的更正**）：
1. `os.TempDir()` 这一层**可以**无卷标（`\\tmp`），所以"Windows 档 Go⛔ 交⛔ 卷标的 temp 根"这句**⛔ 成立**；
2. 真正⛔ 可构造的是"**无卷标 AND 可创建**"那一枚合取——`\\tmp` 创建即失败（UNC 无服务端），`/`、`\` 会被 `GetTempPath` **补当前盘符**；⇒ 结论与裁语一致（守卫⛔ 能挂在 env 台面），**但成立的关键是"补盘符＋UNC ⛔ 可创建"这两条 Win32 行为**，⛔ 是"`TMP=/tmp` 会被补成 `D:\tmp`"那一条（那条只对了一半：本腿 cmd 那发补的是**当前盘** `C:\`）。
3. 那发历史记录里的**形状**（`"\TestClassifyAnchorSpellingIsNotVerdict…\001\profile\…"`＝**单个前导反斜杠、无盘符**）恰好＝`filepath.Join("/", name)` 的产物，**而 `GetTempPath` 本腿量过⛔ 会交回那枚值**（它总要补盘符或落 UNC）⇒ 历史那一发的**真正成因仍⛔ 定**（候选：当时的工具链在 cmd 之外递了枚⛔ 经 `GetTempPath` 的 temp 根、或当时 cwd 在 UNC 上使"当前盘符"⛔ 存在）。⛔ 猜，登记 **`Q-308-v1-C`**（成因未定；⛔ 影响交付判语）。

## 四、单裁结论（一句话版）

**归因句⛔ 立**（两路反证：stdlib 源码⛔ 读 `TMPDIR` ＋ 现跑网格）；**现象⛔ 否**（原始件在盘，且本腿已在输入面把同一崩形、同一整包截断后果量出来，⛔ 再依赖那台机器）；**前提"本机造不出"改口径后仍成立**（可创建的无卷标 temp ⛔ 存在；`os.TempDir()` 层面的无卷标形态**存在**，A846 §1 那句要照上表第 2 条改）。

## 五、a1 那三枚被引用的在场尺，本腿自己看了

`probes/308/a1/logs/probe-T1-fwd-slash.txt`（391 B，**有真读数**：`os.TempDir()="D:\\tmp"`、`VolumeName="D:"`、`env TMP="/tmp" TEMP="/tmp" TMPDIR=""`）⇦ §1 引的就是它，**⛔ 空口**；§3 那句"五把在场尺全非空"本腿逐枚 `wc -c` 复过（T1 95–472／T2 1002／T3 372–524／default 512／bare-root 379 ⇒ **成立**）。
⚠ 但同一目录里另外三枚**看起来像**同族证据的件是**空壳**：`probe-sweep-shapes.txt`（120 B，8 行右侧**全空**）、`probe-sweep-gotmpdir.txt`（437 B，8 行里只 2 行有正文）、`probe-sweep-unset.txt`（**0 B**）。⇒ 下一位⛔ 再拿 `probe-sweep-*` 当"逐枚坐实"引（登记 `Q-308-v1-D`，⛔ 翻 a1 的账，⛔ 是本票射程）。
