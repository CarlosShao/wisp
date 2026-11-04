# 189 AC 裁决（非实现者：编排者，2026-11-04 15:0x +08）

> 来路：写腿 `189-r1`（dud，未提交、build 失败、类型未定义）→ 收尾腿 `189-r2` 补齐并提交 `e1345b0a0`。编排者＝非实现者，收件只认盘上三把尺＋自己复跑。

## 一、落地（`e1345b0a0`，`git show --name-status` 恰好三枚，无夹带）
- `M internal/panel/git.go`（+55/−8）：`GitView` 加 `Uncommitted`(:178)／`Diff`(:181) 两**嵌套子格**（`composer.git` 下，**无顶层新键** ⇒ 不触 Q-51 键集）；`newGitView` 默认第三态；`ReadGit` 非仓/不可读早退同第三态；working tree 定后按 `CurrentWorktree` 调 `readGitWorkingState`(:242)；头注 "NO EXTERNAL PROCESS" 改诚实版。
- `A internal/panel/git_readonly_exec.go`（833 行）：五类型 `GitChangedFile`(:89)/`GitUncommittedCounts`(:111)/`GitUncommittedView`(:135)/`GitDiffFile`(:163)/`GitDiffView`(:200)；`Counts *T omitempty`（未读＝无键、渲染不出 0）、`BytesTotal *int64 omitempty`、`Files` 恒 `[]T{}`；`GitReadNoGit` 第三态＋`gitNoBinaryReason`；解析逻辑承接 r1 一字未改（仅修两处编译/误判缺陷）。
- `A internal/panel/git_189_test.go`（827 行，六枚判据）：①合成真仓双尺对账（`git status --porcelain`／`git diff HEAD --numstat`）②not-a-repo/not-probed：`counts` 键缺席、`files:"[]"` ③PATH 置空→`GitReadNoGit` 绝不报 0 ④`.git/index`＋`HEAD`＋工作树 sha256 跑前后不变 ⑤AST 结构尺（唯一 `exec.CommandContext` 在具名文件、无 os 写调用、禁行动词字面量）⑥正控种 `os.WriteFile`＋`git commit` 必须响。

## 二、门禁读数（编排者现跑）
| 尺 | 读数 |
|---|---|
| `go build ./internal/panel ./internal/tools` | rc=0 |
| `go test -count=1 ./internal/panel -run '189\|Uncommitted\|Diff'` | **6/6 PASS**（1.799s） |
| `go test -count=1 ./internal/panel`（整包） | **PASS=178 FAIL=1**，唯一红＝`TestC21DesignTokensFourWayAgree` |
| `gofumpt -l` 三件 | 空 |
| `sh scripts/d22scan.sh` | clean（rc=0） |
| `go vet ./internal/panel` | rc=0 |

## 三、两处非实现者自证（编排者亲跑）
1. **正控（第 21 处纪律：红是不是别人造成的）**：`git clone --shared` 导出父提交 `66fdb2275` 跑同一测试 ⇒ `TestC21DesignTokensFourWayAgree` **同样 FAIL**（`tokens_fourway_test.go:638/682` 逐字读数）⇒ **该红非 189 引入**，属 tokens 四向域（`internal/ball/tokens.go` vs `docs/evidence/s1/c21-native-tokens.md` vs `frontend/`），本票无关、不修。
2. **变异自证（第三态牙）**：`-overlay` 仓外 mutant 把 `gitUncommittedUnavailable` 在 no-git/unreadable 时改成返回 `Kind:GitReadCounted`＋`Counts:&GitUncommittedCounts{}`（即"没读却报 0"），build 过；两枚第三态测试**双双 RED**，红线逐字：
   - `git_189_test.go:322: not-a-repo: uncommitted kind="counted" counts=&{0 0 0 0 0 0}, want an unread state with nil counts`
   - `git_189_test.go:392: the no-git state carries a number: counts=&{0 0 0 0 0 0} diffTotal=0`
   ⇒ 「0＝干净／未读≠0」这条本票命门**有真牙**。

## 四、判
- **AC#2 只读落地＝成立**（值来自真仓，双尺对账）。
- **AC#3 只读边界＝成立**（AST 结构尺＋正控，无写落点；不注册为 Tool）。
- **AC#4 大 diff 代价＝成立**（`ByteCap`/`Truncated`/`BytesTotal`/`RecoverPointer` 四元组在，测试钉 `byteCap==gitDiffByteCap` 与 `RecoverPointer` 含 `diff HEAD --`）。
- **AC#5 契约轴＝成立**（commit 零碰 PLAN/specs/thresholds/golden/allowlist；无 C17 新方法）。
- **AC#6 门禁＝成立**（上表；`internal/panel` 已知红 1 枚且经正控证非本票引入）。
- **AC#1 只读设计核**：由 `189-a1` 出件（`docs/evidence/s1/189-uncommitted-count-design-a1.md`＋`probes/189/a1/`）承接，本裁决不代勾（归该件之非实现者判语，已存）。
- ⚠ 遗留（非本票引入，189-a1 §9 未裁）：`readGitWorkingState` 随每次 `ReadGit` 同步起 `status`/`numstat`/最多 12 枚 `diff`，泵重算节奏与缓存键未定 —— 归 145/后续腿。
