# 票 88 独立对抗验收 —— `ban #6` 这道门现在到底有没有牙

**验收代理：** `acceptor-ticket88`（独立对抗验收，不写生产码、不 push）
**被验收的票：** `.scratch/wisp/issues/88-arm-ban6-now-that-frontend-exists.md`（`Status: ready-for-review`，六框自勾）
**代码 SHA：** `84e4161`（`git merge-base --is-ancestor` ⇒ 确为 `HEAD` 的祖先）
**接续账：** docs commit `aa5b90c` / `557b458`
**我的仓外纯净快照：** `/tmp/wisp88acc-88`（`git archive 84e4161 | tar -x`，**后缀 `wisp88acc-88` 是我自己的会话名**，
与前任的 `/tmp/wisp88-base`、88b 的 `/tmp/wisp88b-flip`、`/tmp/wisp88b-head` 不同名）
**工作树状态：** 我只读不写；工作树里 `frontend/**`、`internal/{config,tools,risk}`、`cmd/wisp/panel_assets.go`
等**在飞改动全部未被我触碰**（`git status` 起手即录，收工再录，见 §9）。
**所有变异只在 `/tmp` 快照内做**，仓内**没建 worktree、没 checkout、没 stash**。

---

## 0. 结论速览（逐条独立复现，每条我自己的原文读数）

| # | 检查项 | 裁决 |
|---|---|---|
| 1 | 台账真数（纯净快照 `sh scripts/d22scan.sh` rc=0 + `ban #6 examined N>0`） | **PASS**（我的 N=**35**，见 §1） |
| 2 | 牙口正向（种 `approval.decide` 进快照 `frontend/` ⇒ rc=1 并点名） | 见 §2 |
| 3 | 牙口反向（删 `frontend/` 但 `live:true` ⇒ 致命，不许变绿） | 见 §3 |
| 4 | "空仪器不算判据"未被削弱（ban 文本/正则/规则原文/allowlist 5 行） | **PASS**，见 §4 |
| 5 | AC#3 五条台账用例（`-count=1`，逐条读断言本体 + 变异） | 见 §5 |
| 6 | AC#1 的"无后缀过滤器"是否真成立 + 钉子变异 | 见 §6 |

---

## 1. 台账真数：我自己的纯净快照读数（AC#1 / AC#5 的门禁本体）

命令（原文）：
```
rm -rf /tmp/wisp88acc-88 && mkdir -p /tmp/wisp88acc-88
git archive 84e4161 | tar -x -C /tmp/wisp88acc-88        # rc=0
cd /tmp/wisp88acc-88 && sh scripts/d22scan.sh            # 逐字同 CI
```
**我的 rc = 0**。逐作用域**原文**（一字未删，全量输出）：
```
d22scan.sh: positive control - go test ./... (tools/d22scan)
ok  	github.com/CarlosShao/wisp/tools/d22scan	4.050s
d22scan.sh: scan of /tmp/wisp88acc-88
d22scan: examined 207 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/wisp88acc-88
d22scan: scope bans #1-5 internal/      examined 187 production Go files
d22scan: scope bans #1-5 cmd/           examined  20 production Go files
d22scan: scope ban #6 frontend/         examined  35 text files
d22scan: scope ban #7 internal/tools/   examined  16 production Go files
d22scan: scope ban #8 design/           examined  16 text files
d22scan: scope ban #8 internal/         examined 314 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined  25 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=187, bans #1-5 cmd/=20, ban #6 frontend/=35, ban #7 internal/tools/=16, ban #8 design/=16, ban #8 internal/=314, ban #8 cmd/=25; ban #8 emoji coverage: design/ 16 text files; internal/ 314 Go files, comments and _test.go included; cmd/ 25 Go files, comments and _test.go included
```
**我的 N = 35（> 0）**，独立交叉核对：`find frontend -type f -not -path '*/node_modules/*' | wc -l` = **35**（同一快照）。
⇒ 报告器自报的文件数**与真实文件数一致**，不是"计数器凭空 +1"。

⚠ **我的 35 为什么与编排者/88b 的 35 相同**：**因为我 archive 的是同一枚 SHA `84e4161`，树没前移到我量的地方**——
相同数字在这里**不是巧合也不是抄来的**，而是"同一棵树的同一把尺"必然结果。
真正的差异发生在**树前移之后**：88b 自己记录工作树里是 38（票 77 在飞加文件），
所以 §7 我另量了 **`HEAD` 纯净快照**那一档（多样本全报），看的是**我的数会不会随树漂**。

---

## 4. 这三处是本票最可能偷偷让步的地方 —— 逐行证明**没有**

### 4.1 没有任何 ban 的文本/正则被改窄

```
git diff 84e4161^ 84e4161 -- tools/d22scan/main.go | grep -E '^[-+]' \
  | grep -E 'regexp\.MustCompile|banned \(|is banned|panelCheck|approvalPanelRe|MatchString'
⇒ 零命中（grep rc=1）
```
再上一道**字节级**保险（把 9 条 ban 正则所在的声明块整块取出来 sha256）：
```
git show 84e4161^:tools/d22scan/main.go | sed -n '95,112p' | sha256sum   = 65dbe9bd813ab357990f79ffae6055d2c5e433d5af5392bb8259131ded94caa5
git show 84e4161  :tools/d22scan/main.go | sed -n '97,114p' | sha256sum   = 65dbe9bd813ab357990f79ffae6055d2c5e433d5af5392bb8259131ded94caa5
⇒ 完全相同（含 `approvalPanelRe = regexp.MustCompile(\`approval\.decide\`)` @ 84e4161:main.go:105，
   `panelCheck` 的命中语 @ :711）
```
⇒ `approval.decide` 这条**判定文本与正则一个字都没动**，D22 的所有权没被代理碰。

### 4.2 "live 作用域 examine 0 个文件 ⇒ 致命"这条规则原文未动

```
规则文案（`an empty instrument is not a verdict (ticket 71 AC#4)`）在 84e4161^:main.go:913
与 84e4161:main.go:962 都存在，行号漂移=纯插入造成的位移。
判据函数本体哈希：
  git show 84e4161^:tools/d22scan/main.go | sed -n '/^func emptyLiveScope/,/^}/p' | sha256sum
    = 9a0cadd297b51a8e413c096efd06d267449ef0e83fcbe65aa2593daa7fece80c
  git show 84e4161  ... 同上
    = 9a0cadd297b51a8e413c096efd06d267449ef0e83fcbe65aa2593daa7fece80c   ⇒ 逐字节相同
配套的反向守卫 `driftedAbsentScope` 函数体同样逐字节相同
  （两枚均为 bd44c25f2208ecb1e837395655574ea30ec33e7db4f9e865397a870eb78f503e）
```
`git diff 84e4161^ 84e4161 -- tools/d22scan/main.go` 全量 **132 行**我**整块读过**：
除注释、`declaredScopes()` 里 ban #6 那一条（`live: true`、去 `absentOK`、改 `note`）、
`scanScope.absentOK` 的**字段注释**、以及把 `verdict()` 拆成
`verdict() → verdictWithScopes(..., declaredScopes(root))` + 测试专用的 `fixtureVerdict()` 之外，**没有任何逻辑改动**。
⇒ 关键点：`main()` 走的仍是 `verdict()`，它**只**通过一次 `declaredScopes()` 取账；
`fixtureVerdict`（可注入豁免作用域的那个）**在 `_test.go` 之外无人调用**（我用 grep 复核，见 §5.4）。
⇒ **裁决：PASS**（这一条门没有被削弱，ban #6 现在反过来**归它管**了）。

### 4.3 `allowlist.txt` 仍是 5 行非注释

```
grep -v '^[[:space:]]*#' tools/d22scan/allowlist.txt | grep -c .
  @84e4161 = 5      @HEAD = 5
git diff 84e4161^ HEAD -- tools/d22scan/allowlist.txt   ⇒ 空输出（本票与其后所有 commit 都没碰它）
```
