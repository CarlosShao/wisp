# 303-a1 `00-anchor.md` — 起手锚（票 303 `AC#0`＋`AC#1`：复现钉死＋仓外 bisect 到一笔）

腿＝`303-a1`。射程＝**只有 `AC#0`＋`AC#1`**（⛔ 修复、⛔ `AC#2`..`AC#5`）。⛔ 任何产码／测试码字节。
写面＝只新建 `.scratch/wisp/probes/303/a1/**` 里的 `.md`／`.txt`（⛔ `.go`／⛔ `.out`）。⛔ push。

---

## 0. 起手现量（同一条命令里取，date 在第一段）

| 尺 | 逐字 | 读数 |
|---|---|---|
| 时刻 | `date` | `Sat Oct 10 16:54:59 CST 2026` |
| HEAD | `git log -1 --format='%H %ad %s'`；`git rev-parse --abbrev-ref HEAD` | `dce0f133b5123967e527ffbbcf9f3fb3ed4ce14c` `Sat Oct 10 16:51:56 2026 +0800`／分支 `dev` |
| 工作树射程 | `git status --porcelain -- cmd internal scripts .github docs \| wc -l` | **0 行** ⇒ 工作树≡HEAD（G4 同口径） |
| 进程 | `tasklist //FI "IMAGENAME eq wisp.exe"` / `balldebug.exe` 计数 | `wisp.exe=0`／`balldebug.exe=0` |
| Go 台面 | `go version` | `go1.27.1 windows/amd64` |
| sherpa DLL | `ls third_party/sherpa-onnx/*.dll`；`ls build/*.dll` | 两处各 3 枚（`onnxruntime`／`sherpa-onnx-c-api`／`sherpa-onnx-cxx-api`） |
| dist 跟踪面 | `git ls-files frontend/dist` | **只有 `frontend/dist/.gitkeep`** ⇒ 干净 clone 无产物字节（票面 `:17` 的陷阱成立） |

⚠ 我开工时的 HEAD＝`dce0f133`（现量，⛔ 抄派单里任何号）。派单没给号，只说"自己现量"——已量。

---

## 1. 我要用的尺清单（逐字，含 PATH 两枚目录）

`PATH` 铺的两枚目录＝**母仓绝对路径**（clone 里没有 `build/` 与 `third_party/` 的 DLL，被 `.gitignore` 挡着；
缺 sherpa DLL 会在 LOAD 时死掉＝`exit status 0xc0000135` 且零 `--- FAIL`＝用例根本没跑＝**无效步**）：

1. `D:/work/workspace/projects plans/Wisp/third_party/sherpa-onnx`
2. `D:/work/workspace/projects plans/Wisp/build`

| # | 尺 | 用途 |
|---|---|---|
| R1 | `go test -count=1 -timeout 420s -v -run 'TestAC14AwaitedBindingReplyReachesThePage' ./cmd/wisp/` | **最小复现集**（票面 `AC#0` ②，⛔ `TestAC13…` 因它依赖 `frontend/dist` 产物字节） |
| R2 | 同 R1 但 `-count=3` | 复现率三色（`AC#0` ③） |
| R3 | `grep -nE '^--- (PASS\|FAIL\|SKIP): TestAC14AwaitedBindingReplyReachesThePage'` | 每步色 |
| R4 | `grep -F 'no report "ac14r-0" from the page within 15s'` | **bisect 判红的唯一判据**；其余红（`[build failed]`／`no tests to run`／`0xc0000135`／别的句子）＝该步**无效**，具名记件并 `git bisect skip` |
| R5 | `git show --name-only --format= <sha>` | 每笔"动的是产码还是测试码"（必分的一刀） |
| R6 | `git log --oneline cc315261..bcd0a543 -- <4 files> \| wc -l` | 嫌疑面枚数复核（两把尺） |

bisect 区间＝`cc315261`（good，10-06 CI 台面拿到回执）→ 我先量到的 bad（终点在 clone 里现量，⛔ 抄号）。
bisect 一律在**仓外** clone（`$HOME/wisp-303-bisect`）；⛔ 仓内 worktree／checkout／stash／clean（AGENTS §1.4）。

---

## 2. 盘上事实的复核（本腿现跑，⛔ 引用行号当权威）

- 基线那发（`cc315261`，10-06 CI）两枚绿句，尺＝`grep -n -e REPLIED -e PUSHED-33R5-OK`
  于 `probes/301/orch/logs/ci-baseline-cli-block.txt`：件 `tmp-baseline-grep.txt`（rc=0，**2 命中**）
  - `:1293` `panel_resident_windows_test.go:825: AC#14 nail 1 (reply hop), page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)`
  - `:1300` `panel_resident_windows_test.go:867: AC#14 nail 2 (Eval push hop), page's own words: title="PUSHED-33R5-OK"`
- 改后那发（`bcd0a543`）三枚红句，尺＝`grep -n -e 'ac14r-0' -e 'no report'` 于 `ci-after-cli-block.txt`：
  `:1335` `"ac13-probe"`／`:1360` `"ac14r-0"`／`:1367` `"ac14-push"`，三句同尾
  `(what DID arrive at the door: nothing at all)`（rc=0）
- 编排者 16:4x 本机那发（`probes/302/orch/g4-local-three-cases.txt`，`G4_rc=1`）：三枚全 `--- FAIL (20.01s)`，
  逐字含 `wisp: panel window is up (test-harness, cold -1.0 ms, hot path 0.0 ms)`、
  `panel thread ending ... window_opened=true`、`panel thread exited cleanly shows=1` ⇒ 窗那一侧是好的。

---

## 3. 与票面／派单的差异（具名，先报后改）

- 无（起手阶段）。派单转述的嫌疑面枚数（**742** 枚全程／**7** 枚那把尺）＝`AC#0` 未要求的读数，
  我在 §1 R6 现量复核后才写数；若与票面 `:14`-`:15` 不符，以票面原文＋我的现量为准并具名报回。
