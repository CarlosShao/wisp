# `00-anchor-and-gates.md` — 票 317 落地腿 `317-r1`（实现者；判语归后续非实现者腿）

起手锚＝现读 `git rev-parse --short HEAD` = **a5cee1f5**（分支 `dev`）；本机形＝`LOCALAPPDATA=C:\Users\swq\AppData\Local`（反斜杠，与 CI runner 同形）＋`GOOS=windows` ⇒ 折法在本机与 CI 同形触发。写面＝`scripts/{portable-tests,winsec-tests,wisp-cli-tests}.sh` ＋本目录。⛔ 零 push、⛔ 未动 `Status:`/票名/任何框。

## 1. 原意（为什么当时折成正斜杠）— 逐条具名

1. **310-r1（`6a8f537a`）＝拼接风格**：折法只出现在 `gotmp_base=${WISP_TEST_GOTMP_ROOT:-"${LOCALAPPDATA//\\//}/wisp-test-gotmp"}`，其唯一作用是让基值与紧随其后的 `"$gotmp_base/portable-$(date …)-$$"` **共用一套分隔符**（把环境给的 `C:\…\Local` 改写成该腿全部读数都在用的 `C:/…` 拼法）。310-r1 自己的 `probes/310/r1/10-impl.md` §2 把该式与"本机现量形 `C:/Users/swq/AppData/Local/wisp-test-gotmp/portable-20261012-115401-2013`"并列登记；块内四条长注释（WHY GOTMPDIR／WHY NOT UNDER %TEMP%／WHY %LOCALAPPDATA%／WHY WINDOWS ONLY）**没有一句**把折法写成需求 ⇒ 它是转录风格，不是判据，也从未被任何读数验证过。
2. **311（死腿 `3b857e53` → `f8b1f446`／`10734fbd`／`0f888a4f`）＝对照位**：`wisp_311_ref_path=${LOCALAPPDATA//\\//}` 只喂两处——卷读取器（而读取器自己第一行就再折一次 `wisp_311_arg=${1//\\//}`，故该折对判语冗余）与 INCONCLUSIVE/GUARD 行的显示文本（块注释逐字写明 volume "read in both the Windows form (C:/x and C:\x)"）。它**从不流向 GOTMPDIR**。
3. **票面三假设逐条排除**：① **不是 linux 腿形状**——整块在 `if [ "$(go env GOOS)" = windows ]` 内，非 windows 支只打一行 `GOTMPDIR left as-is on GOOS=…`；② **不是"只许 C: 卷"**——折法（10-12 12:15）早于 311 的换卷守卫（10-12 15:04），卷的取舍是量出来的（同码 `58/0` ⇄ `56/2`），且两稿都是正斜杠拼法 ⇒ 折法从未编码卷规则；③ **无测试依赖**——310 落地时 `GOTMPDIR` 全仓 0 命中（`310-a1` 普查＋编排者复跑，见票 310 日志），Windows 上 `os.TempDir()` 不读 GOTMPDIR，没有测试读它的拼法。**反证**：真正存在的依赖是反向的——`internal/risk` 的 116/105 族把 `t.TempDir()` 的字面拼法直接用作字符串比较（`tree + sepStr + …`）与真 `mklink /J` 目标（`target = t.TempDir()` 不经 `filepath.Join`），它们要的是**原生**根拼法。

## 2. 复现（两形；直测 ＋ harness 前后）

- 直测（236-v2 形状，`go test -run <四名> -count=1 ./internal/risk/`）：
  - posix `GOTMPDIR='C:/tmp/wisp317r1gotmp/posix'` ⇒ **rc=1、4 FAIL / 0 PASS**（`logs/10-red-posix-gotmpdir.txt`）。逐字红句：`precondition: deepestExistingAncestor("C:\tmp\…\001\%WISP116ANC\B%\new.md") = "C:\tmp\…\B%", want "C:/tmp/wisp317r1gotmp/posix\TestSync…\001\%WISP116ANC\B%"`；`mklink /J C:\tmp\…\002\proj C:/tmp/wisp317r1gotmp/posix\…\001 failed: exit status 1: Invalid switch - "tmp".`
  - 反斜杠 `GOTMPDIR='C:\tmp\wisp317r1gotmp\nat'` ⇒ **rc=0、4 PASS / 0 FAIL**（`logs/11-green-backslash-gotmpdir.txt`）。
- harness（修的正是这条；`GOTMPDIR`／`WISP_TEST_GOTMP_ROOT` 抽干后 `bash scripts/portable-tests.sh ./internal/risk/`）：
  - 改前（`logs/12-harness-before-internal-risk.txt`，rc=1）：`portable-tests.sh: ticket 310 AC#1: GOTMPDIR=C:/Users/swq/AppData/Local/wisp-test-gotmp/portable-20261014-133428-14242 - …`；四数 `=== RUN=200  --- PASS=128  --- FAIL=4  --- SKIP=0`；红句逐字 `mklink /J C:\Users\swq\…\002\proj C:/Users/swq/AppData/Local/wisp-test-gotmp/…\001 failed: exit status 1: Invalid switch - "Users".` ＝**与 CI 同句尾**。
  - 改后（`logs/13-harness-after-internal-risk.txt`，rc=0）：`…GOTMPDIR=C:\Users\swq\AppData\Local\wisp-test-gotmp\portable-20261014-133651-14508 - …`（原生）；四数 `=== RUN=200  --- PASS=132  --- FAIL=0  --- SKIP=0`；`ticket 311` 命中 **0**（同卷，守卫应静且静）。
- 边界（不许把尺拆掉）：显式 posix override `WISP_TEST_GOTMP_ROOT='C:/tmp/wisp317r1gotmp/posix-override'`（`logs/17-posix-override-risk.txt`）⇒ 值原样进 GOTMPDIR（`C:/…\portable-…`）、四数仍 `FAIL=4` rc=1 ⇒ 用例的敏感度与复现通道未动，改的只是 **harness 自产默认值**的拼法。

## 3. 修（逐处判）

- **ⓐ 三处 GOTMPDIR 赋值＝改为原生分隔符合成**：`portable-tests.sh`（`gotmp_base`＋`gotmp_run`）、`winsec-tests.sh`、`wisp-cli-tests.sh`——`"${LOCALAPPDATA//\\//}/wisp-test-gotmp"` → `"$LOCALAPPDATA\\wisp-test-gotmp"`，尾接 `\\<entry>-…`；`WISP_TEST_GOTMP_ROOT` **仍原样采用**（不重写操作员拼法，见 §2 边界发）。三处各加 ticket 317 注释（portable 全释，另两枚指向 portable 的注）。`git diff --stat`＝`15+/2-`、`8+/1-`、`8+/1-`（`logs/24-no-fold-left.txt`）。
- **ⓑ `wisp_311_ref_path=${LOCALAPPDATA//\\//}` 三处＝判不修**（`portable-tests.sh`、`winsec-tests.sh`、`wisp-cli-tests.sh` 各一）：不流向 GOTMPDIR（见 §1.2），改它会动 311 已钉的显示文案；读取器两边形状都收。
- **ⓒ 第四处同形（射程外，具名上报）**：`tools/d22scan/runtests.sh:92` 的 `"${LOCALAPPDATA//\\//}/wisp-test-gotmp"` 同形仍在。本腿写面＝`scripts/**`（派单 §3），⛔ 未动；它今天不被 CI 咬（直呼 runtests.sh 的 windows 步只有 junction placeholder，其 mklink 目标经 `filepath.Join`＝Clean 归一化；run `38508311376` 的 test-windows 5 增里无它）。本腿 `sh scripts/d22scan.sh` 那一发的 GOTMPDIR 行仍是折形（`logs/22-d22scan.txt:2`）而门照旧绿 ⇒ 登记为"同形、未爆"。建议下一枚同形腿按同法修（一行）。
- ⛔ 未动 174-r5 已入册 ledger 行（`portable-tests.sh` 的 `174r3` 条目文本逐字未变；diff 只有上面三处 hunk）；⛔ 未放开/删任何测试或 skip；⛔ 未动 `.github/**`、⛔ 未动三枚 ref_path。

## 4. 修后读数与门禁

- winsec 入口（`bash scripts/winsec-tests.sh ./internal/winsec/`，`logs/14-winsec-entry-after.txt`，rc=0）：GOTMPDIR＝`C:\Users\swq\AppData\Local\wisp-test-gotmp\winsec-20261014-133730-14659`（原生）；`=== RUN=101  --- PASS=58  --- FAIL=0  --- SKIP=0`；`ticket 311` 行 0。
- cli 入口（`logs/15-cli-entry-after.txt`）：GOTMPDIR＝`…\cli-20261014-133811-14870`（原生）；该档本机 `rc=1` 而四数 `PASS=129 FAIL=0`＝**已知票 314 形状（cmd/wisp 测试二进制整进程死亡）**，与票 311 日志同判 ⇒ ⛔ 本修（该发期间 `316-r1` 正在写 `cmd/wisp`）。
- 311 换卷守卫正控（`logs/16-winsec-crossvolume-control.txt`）：`WISP_TEST_GOTMP_ROOT='D:\tmp\wisp317r1gotmp-xvol'` ⇒ 310 行先建根（原生 `D:\…\winsec-…`）、随后 311 `GUARD` 具名两 offender、`rc=2`、`=== RUN` 0 ⇒ 守卫未被本修松动。
- 门禁：`bash -n` ×3 **rc=0**（`logs/20-bash-n.txt`）；`--scope=census` **rc=0**，totals 行逐字＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`（`logs/21-census-after.txt`）；`sh scripts/d22scan.sh` **rc=0** clean（`logs/22-d22scan.txt`）；长度门 `--with-self-test` **rc=0** `VERDICT GREEN`（`logs/23-pathlen.txt`；本批件入索引后再跑一发＝`logs/25-pathlen-after-staging.txt` rc=0 同判）。
- **下一发 CI windows 腿判法**：`--scope=windows` 步的四数应从 `FAIL=4` 转 `FAIL=0`，且日志 `GOTMPDIR=` 行应是 `C:\Users\runneradmin\…` 原生形；若仍见 `C:/…` 或 `Invalid switch` ⇒ 先读那一行（本修未生效 ⇄ 另一折法）。

## 5. 具名没做 / 只建不删

- 没跑：完整 `--scope=windows` 八包（只跑其中的 `./internal/risk/`）与全量 `--scope=cli` 判语（票 314 域）；末发 CI 色由编排者取。
- 没动：`tools/d22scan/runtests.sh`（§3ⓒ）、`.github/**`、三枚 `wisp_311_ref_path`、ledger 行、任何测试/断言/skip。
- 本腿只建不删的盘上新根：`C:\Users\swq\AppData\Local\wisp-test-gotmp\{portable-20261014-133428-14242, portable-20261014-133651-14508, winsec-20261014-133730-14659, cli-20261014-133811-14870, portable-20261014-134311-15445, d22scan-20261014-134346-15842}`、`C:\tmp\wisp317r1gotmp\{posix, nat, posix-override}`、`D:\tmp\wisp317r1gotmp-xvol\winsec-20261014-134148-15253`（跨卷正控那发 310 先建、311 随后拒绝）。
