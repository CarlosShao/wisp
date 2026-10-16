# `00-readings.md` — 票 317 验收腿 `317-v1`（非实现者；只读产码，读数全落本目录）

起手锚＝现读 `git rev-parse --short HEAD` = **55e11b2f**（分支 `dev`；本腿读数台＝`git archive 8235e683` 的仓外导出树）。
写面＝`.scratch/wisp/probes/317/v1/**` ＋ `docs/evidence/s1/317-gotmp-separator.md` ＋ 票 317 末尾一条 bullet。
⛔ 未动任何产码／脚本／测试件／`.github/**`；⛔ 未勾票面框、⛔ 未动 `Status:`／票名／`-done`；⛔ 未 push。

## 0. 台面与尺（先声明射程，后面每个数都指回这里）

| 尺 | 射程 | 量/算 | 备注 |
|---|---|---|---|
| 导出树 `D:\tmp\wisp317v1\tree`＝`git archive 8235e683 \| tar -x` | 仓外、无 `.git` | — | 六件套逐一 `[ -e ]` 验在：`go.mod go.sum internal cmd tools scripts`；`find … -name '*.go' \| wc -l` ＝ **1087 枚 .go 文件**（非空尺） |
| 另两枚突变树 `tree-m1`／`tree-m2` | 仓外，同一发 `git archive` | — | 突变只打在 `scripts/portable-tests.sh`，`diff` 逐字列出（见 §3），⛔ 未用 `-overlay` |
| 用例名尺 `go test -run <名> -list <名>` | 导出树 `./internal/risk/` | 量 | 命中判据＝该枚名字**自己被打出来**（防"不存在的名数 0"） |
| 四数尺 `grep -c '^--- PASS' / '^--- FAIL' / '^--- SKIP' / '^=== RUN'` | 各自日志 | 算（从 `-v` 输出数**顶层用例**） | 单位一律写"顶层用例"，⛔ 与子测试混计 |
| `gofmt -l` | 两台面分别跑（见 §5） | 量 | 本机 `gofumpt` ⛔ 装＝`command not found`／rc=127 ⇒ ⛔ 拿 gofmt 的读数冒充门禁那把 |
| 折法名册尺 `grep -rFn '${LOCALAPPDATA//\\//}'` | 共享工作树 `scripts tools` | 量 | ★第一次用 `grep -rn` 正则形报 **0 命中**＝假 0（正则里 `\\` 吃掉一层）；改 `-F` 固定串＋正控（同串在 `*.sh` 命中 **16 行**）后才是读数 |

## 1. 派单要我自查的三条"更正"（逐条复核，⛔ 直接采信）

- **(a) 元凶行号已移位＝成立。** 票面写 `scripts/portable-tests.sh:144`；现读该组成式在
  `:155`（基座 `gotmp_base=${WISP_TEST_GOTMP_ROOT:-"$LOCALAPPDATA\\wisp-test-gotmp"}`）＋`:156`（尾接 `gotmp_run="$gotmp_base\\portable-…"`）；
  另两枚入口各一行：`scripts/winsec-tests.sh:69`、`scripts/wisp-cli-tests.sh:82`。`scripts/portable-tests.sh:144` 今天是一行注释。
- **(b) `wisp_311_ref_path` 现量＝4 处赋值（⛔ 三处）。** `scripts/portable-tests.sh:818`／`scripts/winsec-tests.sh:174`／
  `scripts/wisp-cli-tests.sh:186`／`tools/d22scan/runtests.sh:155`（尺＝`grep -rn wisp_311_ref_path scripts tools`，全部同形 `${LOCALAPPDATA//\\//}`）。
- **(c) `tools/d22scan/runtests.sh:92` 今天＝混形（⛔ 纯 posix 形）。** 逐字节（`sed -n '92p' \| od -c`）：
  `gotmp_run="${WISP_TEST_GOTMP_ROOT:-"$LOCALAPPDATA\\wisp-test-gotmp"}/d22scan-$(date …)-$$"`＝基座反斜杠＋尾接正斜杠。
  出身＝编排者两笔微修只动了基座：`581dd252`（`-  …${LOCALAPPDATA//\\//}/wisp-test-gotmp…` → `…"$LOCALAPPDATA\wisp-test-gotmp"…`，尾接 `/` 原样留着）
  ＋`76b9e214`（单反斜杠 → 双）。⇒ 实现腿那句"仍折形"按现量要重述成"基座已原生、尾接仍正斜杠"。

## 2. AC#1（复现）——我判"那些读数算不算凭据"

- **① 红在哪一行（判"是不是建目录失败"）**：实现腿 `logs/10` 与我自己的 R2b 落在**同四行**——
  `internal/risk/syncdirs_ancestor_actable_leg_116_test.go:182`（`precondition: deepestExistingAncestor(…) = 原生形, want 正斜杠基座形`）、
  同件 `:286`（`precondition: the ~ spelling must resolve under the relocated home "C:/…"`）、
  `internal/risk/syncdirs_ancestor_reparse_ticket105_windows_test.go:59` 与 `:116`（**真 `mklink /J` 的目标参数**被 cmd 当开关）。
  ⛔ 任何一发红是 `mkdir`/建根失败；harness 的 `DEGRADED` 行在两发里都不出现。**读数成立。**
- **② 成对两发是否同一台面只差分隔符＝原发有一处欠洁，我用另一发补上了。** 实现腿那对是
  `C:/tmp/wisp317r1gotmp/posix` ⇄ `C:\tmp\wisp317r1gotmp\nat`：同一台面（同树同命令）、但**叶名也不同**（`posix` ⇄ `nat`），
  严格说不是"只差分隔符"。我这一发把两形钉在**同一个目录**上（`C:\tmp\wisp317v1pair` ⇄ `C:/tmp/wisp317v1pair`，`od -c` 逐字进日志）：
  原生形 ⇒ 顶层用例 `PASS=4 FAIL=0`、rc=0；正斜杠形 ⇒ `PASS=0 FAIL=4`、rc=1。同四枚名、同四行号 ⇒ 混淆项排除，原读数随之加固。
- **③ 四名今天真存在＝成立。** R1 逐枚 `-list` 自名命中：四枚 116/105 ＋ `TestPathResolverJunctionWindows`（5/5）。
- **"先读 310／311 的原意"这一半**：实现腿 §1.1 引 `probes/310/r1/10-impl.md` §2 的两行（`:23` 组成式、`:24` 本机现量形 `C:/Users/…`）我现读到；
  `scripts/portable-tests.sh:113-151` 四条长注释里 `grep -n 'fold|slash|separator|posix'` 只命中 **ticket 317 自己那两行**（`:31`/`:38` 相对位），
  旧注释无一句把折法写成需求 ⇒ "拼接风格"这一判读有凭据。"无测试读 GOTMPDIR 拼法"我也独立复核：`grep -rn GOTMPDIR --include=*.go internal cmd tools`
  只命中**别腿在飞件** `internal/winsec/openleaf_guard_322_windows_test.go:95` 的一行**注释**（0 处真读取）。
- **★AC#1 判语＝成立。**

## 3. AC#2（修）——四支逐枚独立复核（物理突变在仓外导出树，⛔ overlay）

| 发 | 台件 | 结果（四数＝顶层用例） | 读数件 |
|---|---|---|---|
| R3 | 出厂形 harness：`env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT bash scripts/portable-tests.sh ./internal/risk/`（导出树 `tree`） | GOTMPDIR 逐字 `C:\Users\swq\AppData\Local\wisp-test-gotmp\portable-20261016-090055-1512`；`RUN=200 PASS=132 FAIL=0 SKIP=0` rc=0 | `logs/r3-harness-asshipped.txt` |
| R4＝**M1 退回原状** | `tree-m1` 里 `:155-156` 逐字改回 `78805129^` 那两枚（`diff`＝只这两行） | GOTMPDIR `C:/Users/swq/AppData/Local/wisp-test-gotmp/portable-20261016-090221-552`；`PASS=128 FAIL=4` rc=1 | `logs/r4-harness-m1-prefix-fold.txt` |
| R5＝**M2 反形**（基座原生、只把尾接 `\\` 写成 `/`＝runtests.sh:92 today 的形状） | `tree-m2`，`diff`＝只 `:156` 一行 | GOTMPDIR `C:\Users\swq\AppData\Local\wisp-test-gotmp/portable-20261016-090256-813`；`PASS=128 FAIL=4` rc=1；红句 `mklink … Invalid switch - "portable-20261016-090256-813\TestSy…"` | `logs/r5-harness-m2-mixed-tail.txt` |
| R6＝override 仍原样 | 出厂形 `tree` ＋ `WISP_TEST_GOTMP_ROOT='C:/tmp/wisp317v1/ovr'`（`od -c` 进日志） | GOTMPDIR `C:/tmp/wisp317v1/ovr\portable-20261016-090416-462`（操作员拼法未被改写）；`PASS=128 FAIL=4` rc=1 | `logs/r6-posix-override.txt` |

- **① 三处 GOTMPDIR 赋值**：现读三枚入口都是**两处位置都原生**（基座 `"$LOCALAPPDATA\\wisp-test-gotmp"`＋尾接 `\\<entry>-`），
  且 R3 打印出的运行时值是**单反斜杠**（`${X:-"…\\…"}` 吃一层这条今天仍成立；`bash -n` ×4 rc=0）。
- **② "override 仍原样采用"防的是什么＝复现通道本身**：R6 实测 posix override 原样进 GOTMPDIR 并把四枚打回红（`FAIL=4` rc=1）
  ⇒ 这一支今天仍成立（四枚入口的 `${WISP_TEST_GOTMP_ROOT:-…}` 写法逐枚 grep 现读，无别腿改写过）。
- **③ "不修的那几处"判对没有＝判对，但其中一条理由无仪器凭据**：4 处 `wisp_311_ref_path` 只喂两样——卷读取器
  （读取器第一行自己再折一次：`wisp_311_arg=${1//\\//}`，4 枚入口各有，现读 `portable:811/winsec:167/cli:179/runtests:148`）与
  INCONCLUSIVE／GUARD 的显示文本；**没有任何 Go 测试或 CI 步骤断言那行文案**（尺＝`grep -rn wisp_311_ref_path --include=*.go --include=*.yml --include=*.ps1 .` 排除 `.scratch` ⇒ 0 命中，
  正控＝同串在 `*.sh` 命中 16 行）。⇒ "不修"这个**处置**对（折形永不流向 GOTMPDIR，且读取器两形都收）；
  实现腿那句"改它会动 311 已钉的显示文案"作为**仪器事实**没有凭据（仓内无人钉），只能算措辞偏强。
- **④ ★恒真检查＝两个方向都红，尺不是空的**：M1（退回原状）4 红、**M2（反形）也 4 红**、且红的是**同四枚名同四行**；
  出厂形 R3 全绿。⇒ 这把尺对"分隔符"在**基座与尾接两个位置**都敏感，317-r1 那两处改动都是承重的，
  ⛔ 落进"绿了就说明不敏感 ⇒ 证据不足"那一支。
- **★AC#2 判语＝成立。**

## 4. 第四处同形（`tools/d22scan/runtests.sh:92`，混形）——会不会爆

| 发 | 命令（出厂形 `tree`，`env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT`） | 结果 |
|---|---|---|
| S1 | `runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows`（＝`ci.yml:831` 真用的那一格） | GOTMPDIR 逐字 `C:\Users\swq\AppData\Local\wisp-test-gotmp/d22scan-20261016-090656-522`（混形）；`top-level: PASS=1 FAIL=0 SKIP=0, === RUN=1, '[no tests to run]'=0` rc=0 |
| S2 | 同一枚 runner、同一枚混形根，`-run` 换成四枚 116/105 名 | 四枚**全红** rc=1；红句 `mklink … Parameter format not correct - "d22scan-…\TestSy…"`（cmd 报错措辞与 M2 的 `Invalid switch` 不同，同因） |

读数件＝`logs/r8-runtests-mixed.txt`。⇒ **"同形、今天未爆"这句要重述**：**形状会爆**（实测在 runtests.sh 自己的根上把四枚打红），
**今天不爆只因为射程**——CI 里直呼 runtests.sh 的三格是 `tools/d22scan ./...`（`ci.yml:82`）、`internal/proc TestLayoutForTestEnv`（`:462`）、
`internal/risk TestPathResolverJunctionWindows`（`:831`），后者的 link/target 走 `filepath.Join`（＝Clean 归一化）所以吃得下混形；四枚 116/105 一名都不在射程里。
⇒ 建议**新立一票**（一行：`/d22scan-` → `\\d22scan-`），并把这条从"登记为未爆"改成"射程遮蔽中的活坑"。本腿写面⛔ 含 `tools/**`，⛔ 未动。

## 5. 门禁五件（共享工作树 `D:\work\workspace\projects plans\Wisp`＝含别腿在飞态；读数件 `logs/r7-gates-sharedtree.txt`）

- `bash -n` 四枚脚本（`portable`／`winsec`／`cli`／`runtests`）＝各 **rc=0**。
- `go build ./...`＝**rc=0**（`322-r1` 的 `internal/winsec` 在飞态此刻可编译 ⇒ 本腿无别腿红要另归）。
- `go vet ./internal/risk/`＝**rc=0**。
- `sh scripts/d22scan.sh`＝**rc=0** clean（`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77`；扫描行 `examined 280 production Go files`／`ban #8 internal= 571 Go files`＝分母含别腿新件，⛔ 归本票）。
- `sh scripts/check-path-length-budget.sh --with-self-test`＝**rc=0** `VERDICT GREEN`（`tracked paths=13497 over-budget=60 covered=60 not-in-roster=0`）。
- `gofmt -l internal/risk scripts tools`（窄于门禁那把；本机 `gofumpt` rc=127 ⛔ 装）＝列出 **1 枚文件**：`internal\risk\provenance.go`。
  ★另节具名＝**台面条件、⛔ 计入本票**：`git status --porcelain -- internal/risk` 空、`git ls-files --eol` 读 `i/lf w/crlf attr/text eol=lf`、
  字节级 `head -c 120 \| od -c` 显示工作树副本为 `\r\n`；同一枚文件在 8235e683 导出树里 `gofmt -l` **0 枚列出**。
  ⇒ 共享工作树那 1 枚是 CRLF 工作树副本（内容＝blob），与本票无关。
- 长度门在导出树⛔ 跑（无 `.git` ⇒ `not inside a git work tree`＝rc=2 射程拒绝，⛔ 当红）；已在仓根跑过＝上面那发 rc=0。
- ⛔ 跑：完整 `--scope=windows` 八包与 `--scope=cli`（票 314 死亡区，本票不需要）；`go test ./cmd/wisp` 整包⛔ 跑；`go mod tidy`／`go get`⛔ 跑。

## 6. 盘上只建不删（本腿新增根）

`C:\tmp\wisp317v1pair`（两形共用同一目录）、`C:\tmp\wisp317v1\ovr\portable-20261016-090416-462`、
`C:\Users\swq\AppData\Local\wisp-test-gotmp\{portable-20261016-090055-1512, portable-20261016-090221-552, portable-20261016-090256-813, d22scan-20261016-090656-522, d22scan-20261016-090659-565}`，
仓外导出／突变树 `D:\tmp\wisp317v1\{tree, tree-m1, tree-m2}`（突变只在树内，⛔ 回写仓内）。
