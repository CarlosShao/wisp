# 300-r1 — `00` 起手锚（落地写码腿）

腿＝`300-r1`（实现者）。非实现者审腿＝`300-v2`（后飞）。
本件只记**起手三读**＋**仪器前置尺**，⛔ 任何判断、⛔ 任何读数结论。

## 1. 起手三读（命令逐字，顺序⛔ 换）

```
$ date '+%Y-%m-%d %H:%M:%S %z'
2026-10-10 13:04:07 +0800

$ git rev-parse HEAD
ec5cdc51ce924ff9e1c7eb9bf6c33f46c635073d

$ git status --porcelain -- internal cmd | wc -l
0
```

- 尺名＝`wc -l`，射程＝`git status --porcelain` 对 `internal`＋`cmd` 两枚路径的输出。
- 判读＝**0 行** ⇒ 树上我射程内（`internal/**`）无别人的未提交改动 ⇒ 起手门**通过**，⛔ 停手。
  （⚠ 这一枚读数是**那一刻**的快照；每笔 commit 前我另跑 `git diff --cached --name-only` 名册尺，见各门件。）

## 2. 环境（同一条命令带出的）

```
$ go version; git branch --show-current; git log --oneline -3
go version go1.27.1 windows/amd64
dev
ec5cdc51 A804＋§4.0be：★我自己把票 300 的 AC#5 那三跳预跑了一遍 ⇒ 抓到一枚通用事实＝**本仓"有没有人跑这枚用例"那把门只到包级、⛔ 到用例级**（任何包里 windows-tagged 的判据天生不在 CI 执行面，而⛔ 任何一枚门会为此变红）；＋⚠⚠**记我自己一枚新形＝编排者自己是一枚隐形写腿**（我把票面改成未提交状态，被在飞的 300-v1 那笔 3da1a1a1 连带 15 枚插入一起收进它的 commit，而它的 message 自称"只加一行"）
3da1a1a1 probe(300-v1): append one reviewer progress line to ticket 300 (verdicts 1-4, AC#2 = go); no checkbox flipped, no other byte of the ticket touched
890cd466 probe(300-v1): harness verdicts - export-tree cmd/wisp dies at load 0xc0000135 without the untracked sherpa dlls (proved by my own two-shape probe), 293-v1's 677/684 need no discount because their red roster names are 8/8 all in cmd/wisp/
```

分支＝`dev`。锚＝`ec5cdc51`。

## 3. 进程前置尺（跑任何读数之前）

```
$ tasklist | grep -iE 'wisp\.exe|balldebug\.exe|mockllm\.exe' ; echo "rc=$?"
mockllm.exe                  20904 Console                    1      6,552 K
mockllm.exe                  25660 Console                    1      7,188 K
rc=0
```

- `wisp.exe` 命中枚数＝尺＝上面 stdout 里 `wisp.exe` 出现次数 ⇒ **0**。
- `balldebug.exe` 命中枚数＝同一段 stdout 里 `balldebug.exe` 出现次数 ⇒ **0**。
- ⚠ **具名登记，⛔ 我杀、⛔ 我因它判任何红绿**：盘上有两枚 `mockllm.exe`（PID `20904`／`25660`），
  编排者已写死＝**别人的**、10-08 起就在。`rc=0` 是 `grep` 命中那两枚 mockllm 的结果，⛔ 是"有 wisp 进程"。

## 4. 本件射程声明

- 本笔 commit 名册＝**一枚文件**（`git show --stat` 尺，逐字进 `logs/`）。
- 本程写面＝`internal/audio/**` ＋ `.scratch/wisp/probes/300/r1/**`（派单写死的射程）。
- ⛔ 动 `cmd/wisp`、`internal/observe/thresholds.go`、`internal/audio/level.go` 的 `SineLevelTolerance`、
  `internal/ball/liquid.go:30` 的 `SilenceLevelGate`、`tools/d22scan/allowlist.txt`、D43 转移表、`PLAN.md`、
  `docs/specs/**`、`frontend/**`、`design/**`、三枚冻结件、golden／`testdata/golden`。⛔ 删文件、⛔ 建 worktree、⛔ push。
- 票面：⛔ 改任何原句、⛔ 翻任何 `- [ ]` 框；只在 `## Progress log` 追加一行（钟点用 `date` 的 stdout）。
