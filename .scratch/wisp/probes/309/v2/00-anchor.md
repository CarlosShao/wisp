# `309-v2` 起手件（锚点／写面／台面／四道门禁）

派单＝`.scratch/wisp/probes/orch/20261011/txt/309v2-dispatch.txt`（逐字读过，判格只按票面原文与 `probes/309/v1/**` 的条件，⛔ 按转述）。

## 1. 起手两把尺（现跑）

- `git rev-parse --short HEAD`＝**`f6067e5d`**＝派单 §1-1 要求的锚 ⇒ 开工。
- `git -c core.quotePath=false diff --cached --name-only | wc -l`＝**0**（与派单"派单时刻＝0"一致 ⇒ 我没看到 staged 别家路径，也没清空过任何东西）。
- 起手框数尺（§5）＝`grep -cE '^[[:space:]]*- \[x\]' .scratch/wisp/issues/309-*.md`＝**3**／`- [ ]`＝**3**。
- 起手 `git status --porcelain -- internal cmd tools scripts`＝**0 行**（⇒ 我这发的"工作树面"读数＝HEAD blob 面，两处口径可互换；别家别腿的在飞态在 `.gitignore`／`go.mod`／`design/**`，⛔ 我碰过）。
- 工具链＝`go version go1.27.1 windows/amd64`；`gofumpt` 在 `$(go env GOPATH)/bin` 在场（本轮⛔ 用）。

## 2. 台面与 provenance（三类，逐读数归类）

| 面 | 怎么来的 | 我这发用在哪儿 |
|---|---|---|
| 工作树 | 母仓 `D:\work\workspace\projects plans\Wisp`，`internal/cmd/tools/scripts` 零脏 | 四道门禁、词面 grep 尺、`find`/`go list` 根尺、`Assessor:` 普查、锚点在盘尺 |
| HEAD blob | `git show f6067e5d:<path>` | 守卫行逐字计数＝**1**；剥注释两版 `cmp`（⇄ `bd124b2a`）＝IDENTICAL |
| 仓外导出树 A／B | `git archive f6067e5d \| tar -x -C /tmp/wisp309v2tree{.Vvd7,B.lfdT}`，**⛔ 混面**（无 `cp` 别的锚 ⇒ `Q-309-r2-tree-rigs-mixed-face` 对我这些读数已结）；路径记在 `logs/tree-path.txt`／`logs/treeB-path.txt` | 全部 20 发突变 Rig（A 组＝仪器四形＋两肢；B/C/E 组＝枚数钉；D 组＝调用点） |

仓内⛔ 建 worktree、⛔ 动一枚交付码／测试码字节、⛔ 删任何件（Rig 加的文件只落在两棵 `/tmp` 导出树里，且每枚先 `cp` 进 `logs/*.code` 留证再删；`probes/**` 下我新建的件后缀只有 `.md`／`.log`／`.txt`／`.code`＝⛔ `.go`、⛔ `.out`）。

## 3. 四道门禁（各一发，`rc=` 逐行显式追加，件⛔ 0 字节）

| 门禁 | 件 | `rc=` | 件字节 | 旁证 |
|---|---|---|---|---|
| `GOFLAGS= go build ./...` | `logs/gate-build.log` | **`rc_build=0`** | 11 B（＝只有 `rc_build=0` 那一行，输出 0 字节＝无诊断） | — |
| `GOFLAGS= go vet ./internal/tools/` | `logs/gate-vet.log` | **`rc_vet=0`** | 9 B（同上） | — |
| `sh scripts/d22scan.sh` | `logs/gate-d22scan.log` | **`rc_d22scan=0`** | 23095 B | 末行逐字含 live 门禁那句 `d22scan: clean - no D22 ban violations`（⛔ 只装自测输出），并带 `bans #1-5 internal/=229`、`ban #8 internal/=531`、`ban #8 cmd/=119` 分母 |
| `bash scripts/check-path-length-budget.sh --with-self-test` | `logs/gate-pathlen.log` | **`rc_pathlen=0`** | 2547 B | 末行 unit cross-check 那句（byte 视图⇄char 视图选中同一 57 枚路径） |

⚠ `set -eu` 结构那两把（d22scan／pathlen）绿时 stdout 非 0 字节，仍按派单 §3 各显式追加了 `rc=` 行；build／vet 那两发 stdout 本就 0 字节 ⇒ 那两枚件的字节数＝`rc` 行自己（⛔ 拿 0 字节件当"本格已交"）。

## 4. 我这发写的件（写面只有这些）

`.scratch/wisp/probes/309/v2/**`（本件＋`10-rigs-and-readings.md`＋`90-final.md`＋`logs/`）
`docs/evidence/s1/309-v2-verdicts.md`（裁决表，AGENTS §1.5 规定的落点）
票 309 末尾 **1 条** append-only Progress log bullet（`git diff --numstat`＝**＋1／−0**）
⛔ 翻框／⛔ 动 `Status:`／⛔ 改名 `-done`／⛔ push／⛔ `git add -A`／`git commit` 不带 pathspec（提交动作见 `90-final.md` §4）。
