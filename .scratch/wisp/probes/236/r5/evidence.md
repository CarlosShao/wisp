# 236-r5 — 票 236 AC#4 写腿交件：`ci.yml` junction 注释改成"指对对象"（非实现者写腿）

- 派单：编排者 2026-10-14 10:4x +08（存档 `.scratch/wisp/probes/orch/20261013/txt/236r5-dispatch.txt`）。
- 起手锚 `b8f883a3`（10-14 10:38:59 +08 现读；派单发送时刻写的 HEAD=`1179e232`，我起手时它已被派单提交自身推进到 `b8f883a3`）。
- `ci.yml` blob `a5075a3e…`（起手＝HEAD，无他腿脏改动）→ `80b7a393…`（改后）。票 236 blob `3786a6ad…`（起手＝HEAD）。
- 写面＝`.github/workflows/ci.yml` 注释行（6→15 行，仅注释）＋本件＋票尾 1 条 bullet。⛔ 未翻任何 AC 框／未动 `Status:`／票名。
- 判据（派单 §七 改判）：那句不是"过期"，是"**指错对象**"；改后注释写"**读数＋日期＋run 号**"，⛔ 不写"今天为真"。

## 1. 改前现读两发（gh 只读；⛔ 本腿未跑 `go test`）

取数方式：`gh api repos/CarlosShao/wisp/actions/runs?branch=dev&per_page=10` 取最近两发 `ci`（run 328／327，同 `dev @ 2c764c04`）；逐发 `gh api .../jobs/<id>` 取步级 conclusion；`gh run view --job <id> --log` 取正文（⚠ 日志里步名列全是 `UNKNOWN STEP`，按正文自报行归因——派单已预判这一坑）。

| 读数 | run 328 | run 327 |
|---|---|---|
| run id / created | 38484004521 / 2026-10-13T23:24:03Z | 38332849280 / 2026-10-13T01:35:26Z |
| job `test-windows` | id 115593146575 / conclusion=failure | id 115086275360 / conclusion=failure |
| **step 11** "PathResolver junction placeholder (real cases tickets 18/20)" | **conclusion=success**（23:34:05Z→23:34:06Z） | **conclusion=success**（01:46:08Z→01:46:09Z） |
| 同一步正文 | `--- PASS: TestPathResolverJunctionWindows (0.02s)` ＋ `runtests.sh: OK - packages=[./internal/risk/ -run TestPathResolverJunctionWindows] top-level: PASS=1 FAIL=0 SKIP=0, === RUN=1, '[no tests to run]'=0` | 同形，(0.03s) |
| step 10 "Portable windows tests (…)" | conclusion=failure，`##[error]Process completed with exit code 1.`（`portable-tests.sh: strict runner exited 1`） | 同形 |
| step 10 里的三枚 FAIL | `TestPathResolverShortNameAListDenied` (0.04s)／`TestPathResolverUNCAListDenied` (0.05s)／`TestPathResolverExtendedLengthPrefixAListDenied` (0.02s)，红点 `pathresolver_junction_windows_test.go:135/:164/:182`，正文均含 `want "C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\…"` | 同三枚（0.03s／0.02s／0.03s） |

- **rc 口径**：junction 步 rc=0——`runtests.sh` 的 `OK` 行只在 `go test` rc=0 且 SKIP=0 且 PASS>0 时打印（`tools/d22scan/runtests.sh` 尾：非 0 时 `echo "…exited $rc…" >&2; exit "$rc"`），步级 conclusion=success 互为印证；step 10 rc=1（`##[error]…exit code 1.`）。
- **附加读数（本腿自己发现，⛔ 派单没写）**：junction 用例在 step 10 里也 PASS（0.08s／0.04s）——两发里它在两条路径上都绿；红的只是同包那三枚短名用例。
- **名称更正**：派单把第二枚缩写写成 `…UNCAListDenied`；HEAD 上的真名＝**`TestPathResolverUNCAListDenied`**（无 `ShortName` 前缀，红点 `:164`）。注释里写的是真名；按 `TestPathResolverShortNameUNCAListDenied` grep 会 0 命中。
- 逐字日志：`logs/01-run328-readings.txt`／`logs/02-run327-readings.txt`（含 gh api 元数据＋原始日志行号）；两份全文日志留在 `logs/run*-test-windows.log`（见 §5 清单，未提交）。

## 2. 改了什么（只注释）

`ci.yml:799` 那个 step（`PathResolver junction placeholder`）的 805–810 六行注释 → 805–819 十五行：

- **删**：`the step FAILS on windows-latest today for its own real reason (run 35551819606 job 106188167785: …classify B, want ClassA)`——指错对象的那半句。
- **写**：junction 步实测 PASS（`dev @ 2c764c04`；run 328 job 115593146575 2026-10-13 23:34Z ＋ run 327 job 115086275360 2026-10-13 01:46Z；conclusion=success ＋ `--- PASS: …` ＋ `runtests.sh: OK … PASS=1 FAIL=0`）。
- **写**：8.3 短名之因 `did not go away - in those two runs it lands on the sibling --scope=windows step` ＋三枚具名用例（真名）＋ `want "C:\\Users\\RUNNER~1\\..."`。
- **保留** A43② 纪律句（runtests.sh 原样传播非零退出码 ⇒ 这仍是一道真门：`so a red here would still be real`）。
- 全文 diff：`logs/03-diff-comment-only.txt`。注释块 6→15 行 ⇒ 该行以下所有行号 +9（老引用 `:519-524` 早已漂；引用本件请用内容锚，不用行号锚）。

## 3. 零行为证明（两把独立尺，均 rc=0）

1. **注释行隔离**：把两版（`git show HEAD:` vs 工作树）的**所有**注释行剥掉后逐行比对——`304 ↔ 304` 行，**非注释内容逐字相同**；且 `git diff --unified=0` 的每一枚 `+`/`-` 行 `lstrip()` 后都以 `#` 开头 ⇒ **非注释改动行 = 0**。行尾未动（两版 CRLF=0，纯 LF）。
2. **YAML**：`yaml.safe_load` 两版都过；**两版解析出的结构体相等**；派单那条字面命令现跑 `yaml ok`、rc=0。仓里没有 actionlint 之类的 workflow linter（`grep -rln "actionlint" scripts tools .github` 零命中），故用"双向 safe_load＋结构相等"做零行为证明。
- 输出：`logs/04-yaml-check.txt`（yaml）／`logs/03-diff-comment-only.txt`（diff＋隔离＋逐行证明）。

## 4. 没做（具名）

- ⛔ 没跑 `go test`／编译／本机复现（CI 日志只读腿；该 step 的求值路径只在托管 runner 上）。
- ⛔ 没碰 `Status:`／票名／AC#4 的框（票尾只追加 1 条 bullet；框归验收/编排者翻）。
- ⛔ 没动其它步骤／触发条件／并发组／`if:`；没加 `continue-on-error`；没新增步骤（diff 可见）。
- ⛔ 没 push；没 amend/reset/rebase/stash；没 `add -A`／`add .`。
- ⛔ 没碰工作树里别人的脏件（起手 `git status --porcelain -- <本腿面>` 为空；本次 commit 只带显式 pathspec）。
- ⛔ 没验证 8.3 短名"机制本身怎么修"（归口票 252／115，超本格射程）；本格只证"该因落在哪一步／哪几枚用例"。
- ⛔ 没改 `docs/reports/pending-and-issues.md`（台账归编排者）。

## 5. 工作目录清单（issues/README 规则 8：只建不删，列清单）

提交集（6 件，全部 ≤200 行）：`evidence.md`、`logs/00-anchor.txt`、`logs/01-run328-readings.txt`、`logs/02-run327-readings.txt`、`logs/03-diff-comment-only.txt`、`logs/04-yaml-check.txt`。

未提交的原始件（留盘不删；01/02 引用它们）：
- `logs/run328-test-windows.log` — 1,075,089 B／6,191 行（`gh run view 38484004521 --repo CarlosShao/wisp --job 115593146575 --log`，rc=0）
- `logs/run327-test-windows.log` — 1,075,384 B／6,198 行（同命令 38332849280／115086275360，rc=0）
- `logs/run328-test-windows.err`／`logs/run327-test-windows.err` — 0 B（无 stderr）

两份全文可由 `gh api repos/CarlosShao/wisp/actions/jobs/<job id>/logs` 重取（job id 见 §1）；不提交的理由＝派单"产出 ≤6 枚 ≤200 行"上限。
