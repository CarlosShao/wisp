# 308-r2 · 00-anchor（起手读数＋射程＋尺清单）

leg=308-r2 · role=narrow write leg（只改注释行） · ticket=`.scratch/wisp/issues/308-spelling-invariance-test-panics-when-temp-has-no-volume.md`
dispatch 给的锚=`202ed4a5` · 我起手现跑 `git rev-parse --short HEAD`=**`202ed4a5`**（一致）
时间＝2026-10-11 14:2x +08 · 台面＝本机 Git Bash（POSIX sh）/ Windows / `go1.27.1 windows/amd64`

## 0. 起手盘面（尺＝`git status --porcelain -- <path>` 行数）

| 尺 | 读数 |
|---|---|
| `git rev-parse --short HEAD` | `202ed4a5` |
| `git status --porcelain -- internal/risk` 行数 | **0**（无人在飞） |
| `git status --porcelain -- internal/tools` 行数 | **0**（起手时 309-r2 未落笔；若中途变红会具名重跑） |

## 1. 射程（写面三件套，逐字照派单）

1. `internal/risk/pathresolver_anchor_spelling_windows_test.go` 注释两处（内容锚 `TestClassifySpellingInvariance` 文档段末句／ticket 308 注释块第二、三句）
2. `internal/risk/pathresolver_volumeless_spelling_308_test.go` 注释一处（`driveShare` 文档段 "single place" 那句）
3. `.scratch/wisp/probes/308/r2/**`（本件＋`10-three-sentences.md`／`20-proof-code-unchanged.md`／`30-gates.md`＋`logs/`）
4. 本票末尾 1 条 append-only bullet（1 加／0 删）

⛔ 可执行行一行；⛔ 动产码；⛔ 动两枚 `//go:build windows` 用例的逻辑；⛔ 动 helper 逻辑；⛔ 加／删用例；⛔ 引 `t.Skip`；⛔ 翻框／碰 `Status:`／改名 `-done`。

## 2. 起手逐字节基线（供 `20` 件比对）

尺＝`git rev-parse HEAD:<path>`（blob）＋`sha256sum <path>`（工作树）

| 文件 | HEAD blob | 工作树 sha256 |
|---|---|---|
| `internal/risk/pathresolver.go`（⛔ 射程） | `15349e28baf70da3bd7aa1d101ebf61cee20fe35` | `c347e1cfc5401f096153d527e57933ee61680f732b6314b2fe68f53925cb422b` |
| `internal/risk/pathresolver_windows.go`（⛔ 射程） | `7b00b022cf437638794e82eaac79e5fa92845a4a` | `1dcf0640ceb4b79c847d2fe601dd7b0c1d79987cc440b5de15aa75ecd4816771` |
| `…pathresolver_anchor_spelling_windows_test.go`（待改） | `52c8fd440d58adc3c426a86f6682600c700cbbaa` | `80c8e5143a56b04b1aef22246092200e48a089aa2094ed0200161b6db693aa12` |
| `…pathresolver_volumeless_spelling_308_test.go`（待改） | `7dc18b3af8b48d227be1f82b259650fa26fd959d` | `f2c0130702384bfd0c60749cd3163d7492fd5351f64ca32ce4529ca04897070e` |

两枚产码 sha256 ⇄ 票面 §4（编排者 14:2x 现跑值）**同一枚**；新件 `f2c01307…` ⇄ `308-r1`／`308-v1`／编排者三处记录同一枚。

## 3. 改前那发（⛔ 在共享工作树跑 ⇒ 仓外导出树）

树＝`logs/pre-tree-path.txt` 里的 `mktemp -d` 目录，`git archive --format=tar 202ed4a5 | tar -x`，`go.mod`／`go.sum` 在场，取出的 `.go` 枚数＝**995**（⛔ 空树，先证非空再信读数）。

| 台面 | 件 | rc | `=== RUN` | 顶层 `--- PASS` | 顶层名册（PASS＋FAIL＋SKIP） | `--- SKIP` | `--- FAIL` | 子层 `    --- PASS` | `panic: runtime error` | 哨兵 A | 哨兵 B |
|---|---|---|---|---|---|---|---|---|---|---|---|
| 默认 temp（`C:\Users\swq\AppData\Local\Temp`，`VolumeName="C:"`） | `logs/pre-head-full.txt` | 0 | 201 | 132 | 133 | 1 | 0 | 68 | 0 | **0** | **0** |
| 8.3 别名 temp（`C:\Users\swq\tmp\RISKAT~1\RUNNER~1`，`VolumeName="C:"`；`RISKAT~1` 字样命中 10 枚＝台面真被用上） | `logs/pre-head-full-alias.txt` | 0 | 201 | 132 | 133 | 1 | 0 | — | 0 | **0** | **0** |

哨兵 A＝`dropped the two volume-derived`（anchor 用例 `len(dropped) > 0` 那一支的 `t.Logf`）
哨兵 B＝`cannot be constructed`（junction 用例 `!haveShare` 那一支的 `t.Fatalf`）
⚠ 词面尺具名：`grep -ciE 'panic|slice bounds'`＝**4**，逐枚看过＝全部是用例名里的 `Panic` 字样（`TestR9PanicInSubAssessorFailsClosed` 两行 RUN＋两行 PASS）⇒ 承重尺＝`grep -c 'panic: runtime error'`＝**0**。

## 4. 分母测量（⛔ 仓内动交付件；探针件在仓外导出树）

探针＝`<probe tree>/internal/risk/zz_r2_denominator_probe_test.go`（`//go:build windows`，只打印不断言），两把台面各一发：
`logs/probe-denominator.txt`（默认 temp，rc=0）／`logs/probe-denominator-alias.txt`（别名 temp，rc=0）。

| 喂进去的输入 | `VolumeName` | `keep` | `dropped` | `getShortPath` | 拼法循环实际枚数 |
|---|---|---|---|---|---|
| `tmpHome(t)` 真 temp ＋ `.ssh\id_testkey`（默认台面） | `C:` | **9** `[dot-and-dotdot extended long lower not-yet-there slashes unc-admin unc-extended upper]` | `[]` | 非空（`…\TESTR2~1\001\SSH~1\ID_TES~1`） | **10**（9 ＋ `8.3-short` 条件位命中） |
| 同一枚输入，别名台面 | `C:` | **9**（名册同上，逐字同） | `[]` | 非空（`…\RISKAT~1\RUNNER~1\TESTR2~1\…`） | **10** |
| 无卷标字面量 `\tmp\wisp-308-no-volume\.ssh\id_testkey` | `""` | **7**（9 去掉 `unc-admin`/`unc-extended`） | `[unc-admin unc-extended]` | **`""`** | **7**（`8.3-short` 位打不上 ⇒ 见 `10` 件对 `(8)` 的更正） |

⇒ 建议替换句里的 **9／10／7** 我复现，**(8) 复现不了**（本机两把可达台面上 `getShortPath` 对无卷标字面量交回空串，枚数停在 7）。

## 5. 编译档位尺（`go list` 射程＝`./internal/risk/`，`GOFLAGS=`）

| GOOS | `.TestGoFiles` 枚数 | 含 `pathresolver_anchor_spelling_windows_test.go` | 含 `pathresolver_junction_windows_test.go` | 含 `pathresolver_volumeless_spelling_308_test.go` |
|---|---|---|---|---|
| windows | 22 | 是 | 是 | 是 |
| linux | 17 | **否** | **否** | 是（新件⛔ tag） |

⇒ "windows-tagged 那一支在 linux 永不求值"这把尺是**结构性**的（档位成员名册），⛔ 靠 docker 跑色推断。CI 跑机色仍＝`Q-308-v1-A`（归编排者，⛔ 本腿能取）。
