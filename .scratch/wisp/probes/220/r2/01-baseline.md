# 票 220 r2 落点甲接线腿（腿名 220-r2）— 第 0／1 步读数：改前基线

- 落盘时刻：`2026-10-16 10:25 +08`（`date '+%F %T %z'` 现跑）
- 起手 HEAD：`e396301d4905fee3a0da25ee3dceb50e36432f3b`（`git rev-parse HEAD` 现跑）
- 派单：`.scratch/wisp/probes/orch/20261016/txt/220r2-dispatch.txt`
- 本件只记"改前"的数，改后读数见同目录 `02-wiring.md`。

## 1 已落的那一半（读回，不重写）

`git merge-base --is-ancestor 480b970d2 HEAD` ⇒ rc=0（是祖先，本腿不重写它）。
`git show --stat 480b970d2` 的产码三枚：`internal/panel/pump.go`（+39）、
`internal/panel/subagent_blocked_220_test.go`（+216，新建）、
`internal/panel/subagent_roster_197.go`（+81）。

在 HEAD 上逐枚现读（尺＝`git grep -n 'L1Windows' HEAD -- internal/panel/pump.go`，射程＝该一枚文件，**量的**）：

| 行 | 内容 |
|---|---|
| `:117` | 注释点名 `LiveL1Windows`（approval 侧的枚举口） |
| `:191` | `PumpSources.L1Windows` 的文档块 |
| `:200` | `L1Windows func() []L1WindowWait`（字段本体） |
| `:324` | `if p.src.L1Windows != nil {`（读者） |
| `:325` | `waits = p.src.L1Windows()` |
| `:327` | `snap.Tasks = taskRosterSectionFrom(p.src.Tasks(), cards, waits)` |

枚数＝5 行命中、消费点 1 处（`:323-327` 那三行是一段）。
`git cat-file -e HEAD:internal/agent/approval/window_read.go` ⇒ rc=0，
枚数＝该文件里 `LiveL1Windows` 命中 2 行（尺＝`git grep -c`，射程＝该一枚文件）：
`:49` 注释、`:63` `func (g *Gate) LiveL1Windows() []L1Window`。

## 2 本腿要补的最后一格：cmd/wisp 的"没人填"那一格（改前基线）

编排者 10-13 量过 `grep -c L1Windows cmd/wisp/run.go`＝0；**本腿自己在 HEAD 与工作树各重跑一遍**，
射程逐枚写清：

### 2.1 尺一：HEAD（提交面，不含任何在飞字节）

```
git grep -c 'L1Windows' HEAD -- cmd/wisp
rc=1，且**零行输出**
```

⇒ **改前生产调用点枚数＝0 枚**（0 个文件命中、0 行命中）。
射程＝`cmd/wisp` 目录全树；HEAD 上 `cmd/wisp/**` 的 `.go` 枚数＝**134 枚**
（尺＝`git ls-tree -r --name-only HEAD cmd/wisp | grep -c '\.go$'`，量的）。
⚠ 工作树 `ls cmd/wisp/*.go | wc -l`＝**133 枚**：尺一按 HEAD 计，工作树那一枚差值与本腿无关
（差在 `cmd/wisp` 之下的子目录件，`git ls-tree -r` 会算、`ls *.go` 不会），本腿不解释它、只把两把尺的号都写下来。

### 2.2 尺二：工作树（落笔前）

```
for f in cmd/wisp/*.go; do grep -c 'L1Windows' "$f"; done | grep -v '^0 '
（零行输出）
grep -rl 'L1Windows' cmd/wisp/
rc=1
```

⇒ 工作树同样 **0 枚命中**，与 HEAD 一致 ⇒ 这一格没有别人的在飞字节，
基线可用作"改前"。旁证（尺＝`git status --porcelain -- cmd/wisp internal/panel internal/agent/approval docs/evidence/s1 .scratch/wisp/issues`，**零行输出**＝这五处工作树干净）。

### 2.3 正控（报"0 命中"之前先跑，两发都在同一刻）

```
git grep -c 'L1Windows' HEAD -- internal/panel/pump.go          ⇒ HEAD:internal/panel/pump.go:5   rc=0
git grep -c 'LiveL1Windows' HEAD -- internal/agent/approval/window_read.go ⇒ ...:2  rc=0
```

⇒ 尺子本身照得见这两处存在的命中，所以 2.1／2.2 的"0"不是仪器瞎了。

## 3 票框尺（起手第 2 步）

尺＝`grep -c '^[[:space:]]*- \[ \]' .scratch/wisp/issues/220-...-drops.md` ⇒ **5**
尺＝`grep -c '^[[:space:]]*- \[x\]'` 同一枚文件 ⇒ **0**（grep 回 rc=1、无匹配行）
⇒ 与派单预期的"5 未勾／0 已勾"一致。本腿一枚都不翻（翻框归非实现者）。

## 4 台面准备与两枚已知雷（本腿先记账，读数在 05-gates.md）

- 跑 go 命令前置 `PATH="$PWD/third_party/sherpa-onnx:$PATH"`；缺它＝`0xc0000135` 台面故障，⛔ 计入判据。
- `GOTMPDIR` 用原生反斜杠形且先 `mkdir -p`；d22scan 把 `GOTMPDIR` 指到非 `%LOCALAPPDATA%` 卷会回 rc=2 的
  311 GUARD 射程拒绝（⛔ 当红）。
- `go test ./cmd/wisp` **整包＝票 314 死亡区**（174–212 s、假分母）⇒ 本腿只跑靶向 `-run`，整包那一发归编排者。
- `gofumpt` ⛔ 在 PATH（rc=127），盘上有 `/d/work/base/gopath/bin/gofumpt.exe` ⇒ 按绝对路径跑，⛔ 不冒充 `gofmt -l`。
