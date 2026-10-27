# 票 347 终裁腿 `347-v1` — POSIX 侧 `resolveHandle` 真落地的非实现者裁决

- 锚：判读树面＝`0c730cac2`（＋`3849e5c95` 补）；门禁/d22scan 实跑树＝**现 HEAD `ad309e0ac` 的 `git archive` 导出树**（打包时 26-r3 已落；其 4 枚 `cmd/wisp` 件不属 347 射程，见判 9）。起手核：`git diff --cached` 空、index 无前置阻塞。
- 我自己跑的：Linux `golang:1.27` 容器（`unset GOTMPDIR TMPDIR`、GOCACHE 指仓外）＋ Windows Git Bash 档，全部在导出树；五发突变全部落在仓外导出树 `D:\tmp\wisp\20261026\347v1\export-mut`，工作树零突变；日志在 `D:\tmp\wisp\20261026\347v1\logs\`（只建不删）。⛔ 未落任何 `.go`（仓内）、⛔ 未翻框、⛔ push。
- 落笔前五问已过（`ruler-scope-checklist.md` 202 行通读）；两次踩尺的坑当场纠正并记进判语（导出树归属核 hash 遇管道重编码→改 `git cat-file` 重定向字节级比；计数遇 `--- SKIP` 行在日志里出现两次→按名去重）。

## 判 1：票面八格现状（逐格现跑数）

| 格 | 现状 | 判 |
|---|---|---|
| AC#0 前置 | 未勾。编排者按铁律第 4 条自裁「一」（`A1275`），与票面「摆机主」字面冲突；台账留「待机主追认（不认可⇒`revert 0c730cac2`）」 | **⛔ 翻**，保持未勾＋残余具名（判 10） |
| AC#1 代价表→甲落地 | 表在（`347-a1`）；甲真落：`pathresolver_other.go` 两桩已删、`resolveHandle`＝`EvalSymlinks`＋空/错守卫、`reparseComponents`＝自 `/` 逐组件 `Lstat` 收 `ModeSymlink`（现读 `:39-51`/`:67-91`） | **能翻** |
| AC#2 单独摘守卫 | **被落地退役**（判 3：M2 我复刻＝0 红） | **能翻**（判语＝退役） |
| AC#3 搬牙不删牙 | 牙全在且承重（判 4） | **能翻** |
| AC#4 两平台成对 | 我自跑四数见判 5 | **能翻** |
| AC#5 零新导出名 | diff 新顶层函数全小写；`pathresolver_windows.go` 零字节；无新 JSON 键/`panel.*` | **能翻** |
| AC#6 纪律格 | 写面＝`internal/risk/**`＋`scripts/portable-tests.sh`＋证据件（派单逐字给的三处）；票面未动一字节（派单禁改票面压住票面那条 bullet——编排者已收件，冲突已由编排者裁定） | **能翻** |
| 禁区 | 票 55 未动、票 342（`internal/winsec`）零字节、票 339 票面未动（现核 `0c730cac2` 文件单） | **能翻** |

## 判 2：AC#1 —— **成立，能翻**

- `resolveHandle` 不再恒 false：现文件 `:39-51`（`EvalSymlinks`＋`p==""`/err 守卫）；`reparseComponents` 已落：`:67-91`（`!IsAbs` 早退、首枚 Lstat 失败即停、`ModeSymlink` 收拼写照）。
- **A 档 POSIX 可达的正向证明我复刻过**：M4（停收 symlink）红恰 2 枚，mutant 逐字 `Resolved=true … Canonical=…/outside/secret.txt`（跟随到界外）⇒ 载荷在 `reparseComponents`；还原后同筛 rc=0。界外 symlink 拒、界内默认拒、例外后 `Classify=A` 的三段读数以 `347-r1` §6 ＋ `pathresolver_symlink_other_test.go`（150 行，untagged）为证。
- 两条射程补充核：`normalizeLocalUNC` 编译期早退在 `pathresolver.go:249-261`；partial 制造器＝dangling symlink（`anchor_fallback_teeth_336_other_test.go`），带"Stat 必败＋Lstat 必成"自断言。

## 判 3：AC#2 —— **「摘守卫＝全 A」被落地退役；判语：能翻，翻成"退役"**

- 票面要防的形：只摘 `blacklist.go:363` 首行守卫、不落解析器 ⇒ `isUnder("","")` 恒真 ⇒ POSIX 反向全量 A。
- **我复刻 M2**（在 `0c730cac2` 树面上只删那两行，联合筛 336＋339＋347 六枚）：**rc=0，`ok`，0 FAIL**。机理＝落地后 `formsOf` 的 root 最短到 `/`，两侧 root 恒非空且不构成祖先关系 ⇒ 该守卫在 POSIX 结构上死掉。守卫的注释（`blacklist.go:369-374`）也已把"ticket 347 landed realpath on POSIX, so 'no resolver here' is no longer such a case"写明——票面要的形已不存在。
- **退役不等于失防**：它防的"缺世界退化成全 A"由 M1 形承重——我复刻 M1（`aRelates` 摘 fail-closed 支）红恰 3 枚（336＋339 两枚），mutant 逐字 `class=B … reason="…id_*"`（A 降 B 且可覆盖）＝牙抓的正是那格。
- ⇒ 判语：**AC#2 按票面判据"去牙那一发点名这一形"已不可能再红（形不存在），终裁结论＝"已退役"（`347-r1` 具名了它而未掩盖，诚实）**。这一格翻成"退役"不是放水：M1/M3/M4/M5 四发定向红证明失败模式全部仍被别处的牙盯着。

## 判 4：AC#3 搬牙不删牙 —— **成立，能翻**

- 牙在：`TestAnchorFallbackTeethPerAnchor336`（`anchor_fallback_teeth_336_roster_test.go:103`，8 枚 roster untagged）＋`TestAnchorMissReasonNamesCandidateAndMissingAnchorPosix339`/`TestAnchorMissKeepsNonOverridableClassPosix339`（`anchor_miss_reason_339_other_test.go:64/:125`，两枚真牙）＋`candidateOffAnchorChain334` 保名搬 `anchor_world_334_test.go:113`（216 消费者 `name_claim_two_arm_216_test.go:89` 未断）。
- **牙承重我验过**：M3（`resolveHandle` 回桩）红恰 4 枚＝336＋339 两枚＋`TestResolveFollowsExemptedSymlinkOntoRealPath347`；还原后同筛 rc=0。⛔ 未删牙、⛔ 未静默绿——搬之前先有 shot B（5 红逐名）过渡读数（我核过 `347r1\linux-shotB.log` 的 FAIL 名册与账面 5 枚逐名同）。

## 判 5：AC#4 两平台门禁 —— **我自跑的四数**（导出树 `ad309e0ac`，⛔ 相减、逐枚点名）

| 档 | 四数（顶层） | rc | 红集（逐名） |
|---|---|---|---|
| Linux `--scope=core`（容器） | **PASS=1394 FAIL=9 SKIP=1**（RUN=2052） | 1 | `NothingWritten…`＋`Dsn196r3…`＋`C26RewriteAccount…`（存量 A380）＋panel 4＋`ReadGit…`；SKIP＝`TestAppendSessionTurnCrashChild` |
| Windows `--scope=windows`（导出树） | **PASS=526 FAIL=1 SKIP=0** | 1 | `TestC26RewriteAccount…`（存量 A380，唯一） |

- **与 347-r1 账面对账**：其门禁级差集＝`+1 枚 config`（已被 cfg125-r1 消掉，判 8）；我这两发里 config 枚不在红集 ⇒ 消红在我的独立读数里成立。`NothingWritten`/`HandDeleted` 是容器内已知 writeguard 轮换对（cfg125-r1 §5 具名）——我那发轮到 `NothingWritten`、cfg125 自己那发轮到 `HandDeleted`，红集其余逐名相同。
- **±1 的尺差（具名，⛔ 当矛盾）**：`347-r1` 账面 Linux 本腿 `1392/11/1`；我按 `^--- (PASS|FAIL|SKIP)` 重数其日志得 `1393/10/1`——差的一枚是 SKIP 行在日志出现两次被两处不同计法吸收（我两发同法计数，自洽）。**这不动任何判定**：红集逐名（不是枚数）才是 AC#4 的尺。
- 包级交叉核（我跑，`0c730cac2` 树面）：Linux risk 包 `127/1/0`（唯一红＝存量 A380）与账面同；Windows risk 包 `147/1/1` 同（日志 `C-risk-mut-premut.log` 前身）。`portable-tests.sh` 行号现读：`core_pin :270`/`win_pin :300`（账面 `:291/:292` 为普查/票面校正值，与本判无关）。

## 判 6：AC#5 零新导出名 —— **成立，能翻**

- `git diff 0c730cac2~1 0c730cac2` 五件产码：新增顶层 `func` 仅 `resolveHandle`/`reparseComponents`（全小写、`!windows` 件内）；无新 `type/var/const` 导出；无 `panel.*`/JSON 键；`pathresolver_windows.go` 零字节；`Result` 结构零改动。20 件清单与账面逐一相符（`+1064/−608` 我重加总＝同值）。

## 判 7：AC#6 边界审计 —— **成立，能翻**

- `git show --name-only 0c730cac2`：`cmd/` **0**、`internal/{tools,config,secret,memory,session,speech}` **0**、`testdata|frontend|design|.github` **0**、`.scratch|docs/specs|docs/PLAN|docs/audit|docs/reports` 0（唯一 docs 面＝证据件自身）；`3849e5c95` 仅证据件。⛔ 未 push（远端核对过）。票面 `Status:`/`-done` 未动。

## 判 8：cfg125-r1 消红复核 —— **成立（我独立复跑）**

- 我现跑 `TestAC1POSIXSeamHoldsC26Pipeline125|TestAC2POSIXSeamInstallSurvivesASymlinkSpelledTemp125`（导出树、容器）：**rc=0 双 PASS**，双子（control_plain_temp/measured_symlink_spelled_temp）全绿、两子进程逐字打 `winsec: sealing path resolver installed`。
- 主题仍在量的证明我核过原文：被测子进程 TMPDIR 仍穿链（`…/varlink125/tmproot125`），只把 AC#1 探针 base 用 `EvalSymlinks` 锚定；三条断言零删。
- 我的 Linux core 门禁 `1394/9/1` 与 cfg125 本腿发 `07-gate-core-postfix.log`（我重数 `1394/9/1`）**同数同红集**（除 writeguard 轮换对互换）⇒ 门禁级红集差 −1 枚（config 那枚）在我的独立读数里坐实。

## 判 9：d22scan 导出树 ＋ 正控（第 21 处定式）

- 导出树（`ad309e0ac`，git 不可咨询＝大声方向）：**rc=1，1 枚 finding**＝`cmd/wisp/resident_voice_windows.go:311 [phantom-citation] "cmd/wisp-side"`。**归属核**：字节级比对＋`git grep` 证明该行是 **26-r3 提交（`ad309e0ac`）带入**；`0c730cac2`/`3849e5c95` 对 `cmd/` 命中 0 ⇒ **347 树面自身零违禁**（347-r1 当时导出树 clean 的账面与此相容）。该枚归 26-r3，具名上报编排者，⛔ 本票责任。
- 正控：同工具扫同树拷贝、塞一枚 `U+2264` 进 `internal/risk/zz_d22_positive_control_347v1.go` 字符串 ⇒ **rc=1 且逐字点名该文件**（`ban #8 … string literals included`）⇒ 尺是活的。

## 判 10：AC#0 追认项（本终裁的特殊格）

- **自裁在"结果对不对"上站得住吗——站得住**。用 `347-a1` 代价表对照：「一」是唯一让 A 档 fail-closed 在 POSIX **真正可达**的形（二＝只修两处不实注释、敞口照旧；三＝零动作），且零契约字面改动（`C26`/`SPEC-06 §4`/`D7` 三处现读核过）、可逆、不花钱——不属铁律第 4 条"必须摆机主"的三类（改契约文字/花钱/不可逆）；普查 §6 独立判"选一成立"。落地后果（判 2-5）全部闭环 ⇒ **追认建议：认可编排者自裁「一」，无需 `revert`**。
- **AC#0 格的处置建议**：**保持未勾＋残余具名**。它问的是"程序合规"（摆机主后才许开工），机主追认那一格只有机主能勾——"结果对"不能替程序格翻框。残余句保持台账原文（「编排者按铁律第 4 条自裁『一』、与票面 AC#0 字面冲突、留待机主追认」），本判语把"结果正确性"的证明补齐。
- **"可单独回退"前提核实——成立**：①产码与牙同在 `0c730cac2`（M3 证明只回退产码⇒牙全红、shot B 证明只回退牙⇒绊线红，两个方向都红⇒只能同退）；②`0c730cac2..HEAD` 对那 20 枚文件**零触碰**（现跑 `git diff --name-only … -- internal/risk/ scripts/portable-tests.sh`＝空）⇒ `revert 0c730cac2` 无冲突；③回退后 cfg125 的修法仍自洽（无链 base 上 `EvalSymlinks` 恒等）。

## 残余清单（具名，⛔ 藏）

1. **A380 存量红** `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`：我两发门禁都在（两平台），根因 `internal/panel` 三处直读 `Result.Canonical`＝禁写面，非 347 可修。
2. **`pathHandleVerify` 仍 `false`**（`syncdirs_other.go`）：347 后只剩"Lstat 过而 realpath 不过"残留；收紧＝独立决定＋独立分母。
3. **`syncdirs.go` 的 `s.home`/`roots[].canon` 折叠形喂 `formsOf`**：同类折叠缺陷、精度非颜色（347-r1 只修了 blacklist 侧）。
4. **macOS 未实测**：`t.TempDir()` 走 `/var→/private/var` 会被默认拒命中 ⇒ **票 55 落地前必须先解决例外注入**（347-r1 在 Linux 容器复现同形＝cfg125 那枚红的同源）。
5. **容器 writeguard 轮换对**（`NothingWritten`⇄`HandDeleted`）：旁证链归"容器既有"，非 347 责任。
6. **26-r3 带进的 d22scan phantom-citation**（`cmd/wisp/resident_voice_windows.go:311`）：现 HEAD 门禁树上的 1 枚红，归 26-r3，具名上报。
7. **AC#0 程序冲突待机主追认**（判 10；告知非请示）。

## 能不能带残余结案 —— **判：能**

AC#1–AC#6 六格判据全部闭环（多数经我独立复跑/复刻）；AC#0 保持未勾＋残余具名＝程序格的正确终态，且"结果正确性"已由代价表＋落地读数补证。残余 1/3 本就不是 347 的判据面（存量/精度）、4 是票 55 的前置、2 是具名独立决定、5/6 是容器与他腿噪声。**建议编排者：六格翻勾＋AC#0 留勾＋残余 1-7 随票具名结案；若机主不追认 ⇒ `revert 0c730cac2` 单独可行（前提已核）。**

## 提交与越界自查

- 本腿写面＝`docs/evidence/s1/347-v1-verdict.md`（本件）＋`D:\tmp\wisp\20261026\347v1\**`（仓外）。⛔ 未改 `internal/**`、`cmd/**`、`testdata/**`、`.scratch/wisp/**`、`docs/reports/**`、`docs/specs/**`、`docs/PLAN.md`、`docs/audit/**`；⛔ 未落任何仓内 `.go`；⛔ 未翻框/改 `Status:`/加 `-done`/改台账；`add`/`commit` 显式 pathspec、⛔ push。
- 日志（只建不删）：`logs/A-cfg125-target.log`（rc=0）、`B-gate-core.log`（1394/9/1）、`E-gate-windows.log`（526/1/0）、`M1..M5-union.log`（红 3/0/4/2/1）、`Mut-restored-union.log`/`Pristine-union.log`（rc=0/rc=0）、`C-risk-mut-premut.log`（127/1/0）、`D-d22scan.log`/`D-d22-poscontrol.log`（rc=1 正控逐字点名）。SHA256 还原核：五枚产码件逐枚 SAME＝True（时机＝最后一发变异之后）。
- 我没核（⛔ 当已核）：macOS 实机；`0c730cac2` 之前的基线门禁我未重跑（用 347-r1/cfg125-r1 的基线日志交叉对账）；`HandDeleted`/`NothingWritten` 轮换只归因到旁证链。
