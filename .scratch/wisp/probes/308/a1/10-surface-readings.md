# 308-a1 · 10-surface-readings（票 308 `AC#0`：四把台面各自的颜色）

时刻 2026-10-11；锚笔 HEAD `02053534`；`internal cmd` 起手脏行数 0（工作树 Go 面干净，已 commit-first 复核）。

## 0. 名册尺口径（两把尺各报一次，⛔ 相加当一枚数）

命令逐字（台面 env 见各发）：

    go test ./internal/risk/ -count=1 -timeout 600s -v -run "TestPathResolverUNCAListDenied|TestClassifySpellingInvariance"

- 顶层尺：`grep -cE '^=== RUN   Test[^/]+$'` ／ `grep -c '^--- PASS'` ／ `'^--- FAIL'` ／ `grep -cE '^[[:space:]]*--- SKIP'`
- 子项尺：`grep -c '^    --- PASS'` ／ `grep -c '^    --- FAIL'`
- panic 尺：`grep -cE 'panic|fatal error|test executable failed'`
- 包级终态行：`grep -E '^(ok|FAIL) '`

## 1. ★头条结论（先说结果，再逐台面）

**本机（win32 + Go 1.27.1 + 这枚 Git Bash/MSYS 工具）造不出"无卷标 temp"台面。**
`os.TempDir()` 走 Windows `GetTempPath`，对任何能建目录成功的 `TMP/TEMP` 形状都会**补上当前盘的盘符**，
`filepath.VolumeName(os.MkdirTemp("",…))` 恒为非空（在场尺逐发坐实，见下）。
⇒ **两枚崩点用例在 T1／T2／T3 三把台面全部 `PASS`、`panic=0`、包级终态行 `ok github.com/CarlosShao/wisp/internal/risk` 在**；
⇒ 本腿**没有一发 live panic**（票面预设的"无卷标⇒崩"在本机不可触发）。历史那发 live panic 只在 `probes/72/a2/logs/t03-broken-surface-note.txt` 里有原始件（见 §4）。

## 2. 台面在场尺（`os.MkdirTemp("", "W308PROBE")`，件 `logs/probe-*.txt`、`logs/T3-alias-probe.txt`）

| 台面 | 喂给 test 子进程的 env（cmd.exe 里 set，绕过 MSYS 改写） | `os.MkdirTemp` 返回长这样 | `VolumeName` 空串？ |
|---|---|---|---|
| T1 posix `/tmp` | `TMP=/tmp TEMP=/tmp` | `D:\tmp\W308PROBE…`（见 `probe-T1-fwd-slash.txt`） | 否（`D:`） |
| T1 bare-root `\` | `TMP=\ TEMP=\` | `D:\W308PROBE…`（见 `probe-T1-bare-root.txt`） | 否（`D:`） |
| T1 default（本工具继承） | 不设（shell 内 `TMP=/tmp`，MSYS 给子进程改成 Windows 形） | `C:\Users\swq\AppData\Local\Temp\W308PROBE…`（`probe-default-D.txt`） | 否（`C:`） |
| T2 有卷标（本机默认） | `TMP=TEMP=C:\Users\swq\AppData\Local\Temp` | `C:\Users\swq\AppData\Local\Temp\W308PROBE…` | 否（`C:`） |
| T3 别名（8.3 短名） | `TMP=TEMP=C:\Users\swq\tmp\RISKAT~1\RUNNER~1` | `C:\Users\swq\tmp\RISKAT~1\RUNNER~1\W308PROBE…`（`T3-alias-probe.txt`） | 否（`C:`） |

**在场尺结论**：能建的台面全带盘符；带不了盘符的形状（相对 `x`、UNC `//tmp`）`os.MkdirTemp` **建失败**（`probe-sweep-gotmpdir.txt`：`GetFileAttributesEx …\x … not found` / `mkdir \\tmp\… network path not found`），且相对形会落到 `internal/risk` 内＝本腿写面禁止。⇒ 无卷标且可建的 temp 在本机⛔ 存在。

## 3. 四发名册读数（逐字见 `logs/`）

| 发 | 台面 | rc | 顶层 RUN | 顶层 PASS | 顶层 FAIL | SKIP | 子项 PASS | 子项 FAIL | panic | 包级终态行 |
|---|---|---|---|---|---|---|---|---|---|---|
| T1 posix | `/tmp`（drive-qualified） | 0 | 2 | 2 | 0 | 0 | 0 | 0 | 0 | `ok github.com/CarlosShao/wisp/internal/risk 0.138s` |
| T2 vol | 本机默认 temp | 0 | 2 | 2 | 0 | 0 | 0 | 0 | 0 | `ok github.com/CarlosShao/wisp/internal/risk 0.081s` |
| T3 alias | 8.3 短名（有卷标） | 0 | 2 | 2 | 0 | 0 | 0 | 0 | 0 | `ok github.com/CarlosShao/wisp/internal/risk 0.060s` |
| T4 full(posix台面) | `/tmp` 整包名册 | 0 | 132 | 131 | 0 | 1 | 68 | 0 | 0 | `ok github.com/CarlosShao/wisp/internal/risk 2.711s` |
| 对照 full(vol台面) | 本机默认整包 | 0 | 132 | 131 | 0 | 1 | 68 | 0 | 0 | `ok github.com/CarlosShao/wisp/internal/risk 4.425s` |

- 两枚崩点在 T1/T2/T3 都各自有 `=== RUN` 行（件 `logs/T1-posix-roster.txt` 等），**都跑到了、都 PASS** ⇒ 本腿对两枚都交的是实测颜色，⛔ 一发挡另一发。
- T4 整包：顶层 `RUN=132` 与 vol 对照**逐枚同数**，131 顶层 + 68 子项 = **混合层级 199**（引此数必标"含子项"，⛔ 当纯顶层）；唯一 SKIP ＝ `TestSyncRegistryProbeLive`（`syncdirs_windows_test.go:133`，无卷标台面下与 vol 台面同跳过）。
- T4 与对照同数 ⇒ **在 posix 台面这一发里，整包其余用例全拿到了结果**（终态行 `ok` 在、panic=0），这是与"崩那族形状＝二进制早死、其余零求值"的**正向对照**。

## 4. 历史 live panic（非本腿实测，逐字自 `probes/72/a2/logs/t03-broken-surface-note.txt`）

那发（72-a2，2026-10-11 10:5x，⛔ 设任何 temp env）四数：`rc=1`／名册只 2 枚／`PASS=0`／`FAIL=2`；
名册止于第 2 枚后崩：

    panic: runtime error: slice bounds out of range [:1] with length 0 [recovered, repanicked]

崩的那枚名册位＝ `TestClassifySpellingInvariance`（原 `anchor_spelling` 体内 `vol[:1]`）；
崩时**同包其余用例零求值**（终态行虽写了 `FAIL … 0.085s`，但那是崩触发的收尾、名册只 2 枚 `=== RUN`，其余⛔ 跑）。
⚠ 本腿⛔ 能在当前工具台面复现这一发（见 §1／§2）；本节仅供编排者对照"真崩"长什么样。
