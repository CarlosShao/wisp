# 票 158 非实现者验收表 r1 —— 六格一次裁，锚钉 `c7f638c`

- 我是**验收方**，没写过被验的码（`SPEC-12 §4.3` #1/#3：裁决者≠实现者）。
- 被验物＝`docs/evidence/s1/158-gate-scope-blind-spot-r1.md`（**我的被验物，不是我的权威**）：锚点上现量 **1077 行**；`git diff --stat c7f638c..HEAD -- 那枚文件` ＝ **空**（HEAD 往前走的那一枚只加了派单文件）⇒ 我从盘上读到的那一份就是被验的那一版。
- 派单＝`.scratch/wisp/dispatches/2026-09-26-203x-accept-158-six-cells-pinned-c7f638c.md`。
- 票面＝`.scratch/wisp/issues/158-the-gate-stands-outside-the-shape-…-package.md`（锚点上现量 **84 行、6 枚 `- [ ]`、0 枚 `- [x]`**；本表末尾 §8 由我给出勾法）。

## 0. step 0 四件原文（派单 §0 要求的那四件）

```
$ date
Sat Sep 26 20:28:39 CST 2026
$ git rev-parse --short HEAD
ec5ecc0
$ git status --porcelain            # 全树：23 行已跟踪脏（design/** 那一堆、probes/152/my152.py、
                                    # 152-…-accept-r1.md 那 3 行未提交自校）＋ 21 枚未跟踪，原文与
                                    # 被验件 §A 那一串同形（差＝别人这一段新落的 probes/156、probes/158/r2 等）
$ git log --oneline -3
ec5ecc0 派单 20:3x: 票 158 非实现者验收（六格一次裁，锚钉 c7f638c）
c7f638c 票 158 r2 · 交件后自纠两处（追加，§A–§H 已入库那一行一字未改：git diff numstat deleted=0）
226d050 ledger(A310): 十四枚调研全部收官、全线封批——Q-56 带回三枚独立见证；票158 的病因改档

$ git rev-parse c7f638c
c7f638c5e9e3b66073a66d680b5e8c0b85a5b7e4        # ← 本表全部判据的锚点
```

- **锚点外还核了三件事**，免得后面每一发都要解释口径：
  `git diff --name-only c7f638c..HEAD -- '*.go'` ＝ **空**（HEAD 前进的那一枚没动任何 Go 码）；
  `git diff c7f638c -- internal cmd tools scripts` ＝ **空**（工作树在这四支上就是锚点那一版）；
  `git status -sb` 首行＝`## dev...cnb/dev [ahead 285]` ⇒ 本程与实现程都是**只 commit 未 push**。
- 工具版本现读（不抄被验件）：`go version go1.27.1 windows/amd64`；
  `$(go env GOPATH)/bin/gofumpt.exe --version` ＝ **`v0.12.0 (go1.27.1)`**（与被验件 §A 末、§D.5 那格同值，这一枚对得上）。
- 争用：开测前后各一枚 `Get-Process` 都命中 `Runner.Listener`（pid 3952，2028 与 2048 两次读数同一枚 pid）
  ⇒ **本机 self-hosted runner 同机在跑，争用未排除**。本表**不含任何延迟／阈值结论**，所有 `ok …Ns` 一律不当证据。
- 本程写面只有两处：本文件 ＋ `.scratch/wisp/probes/158/accept-r1/**`。两枚半件（`probes/152/my152.py`、
  `docs/evidence/s1/152-…-accept-r1.md`）**未碰、未还原、未 commit、未补完**；`frontend/**`、`design/**` **未写一个字节**。

## 1. 三把判据仪器先各自打过正控（派单 §1）

> 规矩：**打不出响声的尺上任何"零"都不算证据。** 三把尺各自的正控读数都在下面，全部本程现跑。

### 1.1 尺一＝门本体 G5（`.scratch/wisp/probes/154/gate-clauses.sh`）——逐发复算，不转抄 rc

```
$ bash .scratch/wisp/probes/154/gate-clauses.sh c7f638c   > probes/158/accept-r1/gate-at-anchor.txt ; rc=0（脚本自身）
全文 124 行；按节数：grep -n '^## ' → 11 枚节标题（G1／G1b／G2／G3／G4／G5 主尺／G5 正控说明／G5-正控／G5 负一负说明／G5-负一负／两枚附）
```

| 腿 | 本程现量（锚点 `c7f638c`） | 被验件自报 | 对得上？ |
|---|---|---|---|
| **G5 主尺** | `UNPAIRED cmd/wisp/panel_assets.go (开方调用点=1)`、未成对枚数＝**1**、该腿 `# rc=1` ⇒ **响** | §2.2 同一行判语 | **是** |
| **G5 正控** | 名册**非空**（`bridge.go:642` 开＋`bridge.go:695` 合，两行都在）、`git grep rc=0`、未成对枚数＝**0** ⇒ **不响** | §2.3（它引 `:691`，那是注释面之前的行号；本程在锚点上读到的是 **`:695`**，位移与被验件 §B 末那条"注释净加 4 行"一致） | **是** |
| **G5 负一负** | 未成对枚数＝**5**（`panel_assets` ＋ 4 枚 `internal/risk/*_test.go`），该腿 `# rc=5` | §2.6 末"负一负＝5" | **是** |

**这把尺自己有没有牙（派单 §2 那两问，本程现量）**：

1. **成对判据剔注释那一行（`pair()` 的 `grep -vE ':[0-9]+:[[:space:]]*(//|/\*|\*)'`）承不承重？** 承。
   同一射程、同一条腿，只把那行去掉：负一负 **5 → 6**，多出的正是
   `internal/tools/bridge_scope_open_ticket158_test.go`（它唯一的 `OpenScope` 出现在**注释**里）。
   原文＝`probes/158/accept-r1/g5-with-risk.txt` 旁边那一发（本件 §1.1 表下第二条命令）。
   ⇒ 这不是"摘掉导致编不过"那种假牙：**外部可见读数（枚数）从 5 变到 6。**
2. **`internal/risk/*` 那条排除承不承重（对被验件 §2.4"其实不承重"那句）？** 确实不承重，本程独立复算：
   把 `internal/risk` 放回射程（仍排 `*_test.go`）⇒ 开方文件 3 枚（`panel_assets`／`provenance.go`／`bridge.go`），
   未成对仍＝**1**（只有 `panel_assets.go`）⇒ 判语一字未变，与 §2.4 末那条"顺带量到的口径事实"同向。
3. **主尺射程里今天有没有"字符串字面量假响"（§2.4 登记的结构性残余）？** 没有：
   本程同尺现量 roster 行数＝**3**、剔注释后调用点行数＝**3** ⇒ 射程内一枚注释/字符串混入都没有。
   那条残余是"未来会假响"的形状，不是今天的读数——被验件就是这么写的，判**如实**。

（一处**过期字**要记：门本体脚本第 96 行附近那条正控说明注释仍写"`bridge.go`（`:642` 与 `:691` 同文件成对）"，
`:691` 在锚点上已是 `:695`。它在**注释**里、不进任何读数（判语是动态算的），**不承重**，但下一位照它去 `sed -n '691p'` 会读到别的东西。）

### 1.2 尺二＝AC#4 那 28 枚禁改面的名册尺（按 commit 号集合 ＋ `--no-walk`）

本程自己跑生成器（`.scratch/wisp/probes/158/r2b/zero-byte-census.sh`，读全文核过：**只往 stdout 打，不写仓里任何东西**），
锚点仍取被验程自己的 step-0 `bb08d1f`，读数落 `probes/158/accept-r1/census-at-my-anchor.txt`：

```
## 1  枚数=5（ef34118 8225cdc ee41e63 3122518 c7f638c；反向自证 8 枚编排者提交不在集合里，逐枚看 subject 前缀都是 ledger(A3xx)/派单）
## 2  名册去重后枚数=49 ;  带 --no-walk = 52 行 ;  不带 --no-walk = 4885 行
## 4  28 支冻结面 越界枚数 全 =0，且每支"该面已跟踪文件数"非零（PLAN.md 1 / docs/specs 14 / internal/risk 37 /
     panel 20 / agent 58 / observe 25 / thresholds.go 1 / llm/golden 3 / 含 golden 的路径 6 / testdata/golden 52 /
     allowlist 1 / slo-check.ps1 1 / tools/d22scan 6 / frontend 85 / design 30 / 151·154·156 证据件与票面各 1 …）
## 5  第二口径（字节级 numstat 打同一串正则）⇒ 只剩一枚 HIT: internal/tools/bridge.go
## 8  §5 那把尺的正控（同一串正则，打在已知动过冻结面的三枚真提交上）：
     12181923 命中=2  FIRE internal/risk/assessor_test.go, internal/risk/rules_scale.go
     00bbb76f 命中=3  FIRE internal/observe/{observer_cost_test.go,sampler.go,thresholds.go}   ← thresholds.go 本尊点得出
     cbbbdf85 命中=1  FIRE internal/agent/testdata/golden/spill-tool.sse                          ← golden 点得出
```

- `--no-walk` 那一发本程自己量到的是 **52 对 4885**（被验程在 3 枚号上量到 41 对 4870）。两次的**形状**一致：
  不带旗标就沿祖先链吞整仓历史 ⇒ **被验件那条"不许拿区间 diff 当名册尺"的理由成立**，本程也照做（本表全部普查按号集合）。
- **预测核**：被验件 §C.5 末与 §I.3 留下一条可复算预测（"枚数 3→4 时名册 40→48，§5 仍只剩 bridge.go，28 支仍全 0；再加一枚＝同样形状再走一遍"）。
  本程在 **5 枚**上复算：枚数 5、名册 **49**、§5 仍只剩 `bridge.go`、28 支仍全 0 ⇒ **预测兑现**。
  多出的那 1 枚＝`probes/158/r2b/msg-selffix.txt`（自纠那一枚的提交信息件），**没有一枚落在冻结面**。
- 完整性另一问：49 枚名册逐枚分类 ＝ 47 枚在 `probes/158/r2b/` ＋ 本证据件 ＋ `bridge.go`，
  **写面之外的路径 0 枚**（`sed` 三条排除后为空集，读数在 `…/census-at-my-anchor.txt` §4 那三段 WRITEFACE）。

### 1.3 尺三＝AC#5 那三道工具（不跑 vs 跑了空输出是两件事，一律打 rc）

```
$ $(go env GOPATH)/bin/gofumpt.exe --version ; rc=0        -> v0.12.0 (go1.27.1)
$ $(go env GOPATH)/bin/gofumpt.exe -l . tools/d22scan tools/mockllm ; rc=0   -> 0 行（stderr 0 字节）
$ $(go env GOPATH)/bin/gofumpt.exe -l internal/tools/ cmd/wisp/ ; rc=0       -> 0 行
$ go vet ./internal/tools/ ./cmd/wisp/ ; rc=0                                -> 输出 0 字节
$ sh scripts/d22scan.sh ; rc=1                                               -> 唯一 finding：frontend/src/components/harness/right-rail.tsx:108 [emoji]
$ (cd tools/d22scan && go run . -root "$PWD/../..") ; rc=1                   -> 1 finding(s)，同一条；八枚作用域 examined 全非零
     bans #1-5 internal/=205  cmd/=23 · ban #6 frontend/=85 · ban #7 internal/tools/=18
     ban #8 design/=39  frontend/=85  internal/=414  cmd/=45            （原文 probes/158/accept-r1/d22scan-direct.txt）
```

⇒ 三道里两道**绿**（gofumpt／vet，且都是"跑了、rc=0、输出为空"而不是"没跑"），第三道 **rc=1**——
它是被验件 §D.5 那枚"未裁"的同一发。那一发怎么判，见 §5.2（派单 §3① 的三件，那里**分三件裁**）。
