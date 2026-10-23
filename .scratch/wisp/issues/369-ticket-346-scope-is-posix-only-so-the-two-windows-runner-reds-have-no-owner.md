# 369 — 票 346 的射程只到 POSIX：**今天红在 Windows runner 上的那两枚弃权仪器无主**

- Status: **ready-for-agent**
- 立案：2026-10-23 17:4x +08 编排者立
- ★**来路（实现者签名、只登记不认领）**：`112-r1` §八（件 `docs/evidence/s1/112-r1-impl.md` **118 行**、`e2e576347`）
- 关联：票 **346**（射程缺一半）· 票 **112**（这道门今天红）· `internal/winsec/**` · `.github/workflows/ci.yml:701`

## 这是什么（人话）

CI 的 `winsec-tests.sh` 这一步（`.github/workflows/ci.yml:701`）**今天在 Windows runner 上 `failure`**。

★**红在"仪器什么都没量到"上**（★**⛔ 判定错**）：
- `…RefusesAForeignNameAndNeverDeletesIt342r1` ⇒ `no leaf leg answered (guard=0, landing's re-check=0)`
- `…PrivateFileRouteHasTheLeafGuardSpeaking322r1` ⇒ `400 raced attempts produced guard-leg refusals=0 and racer flips=411`

★★**而票 346 自称治这件事 —— 但它的射程只到两枚 POSIX 仪器**：

★票面第 4 行逐字：
> `射程＝internal/winsec/openleaf_322_other_test.go ⇄ internal/winsec/posix_seal_window_120_other_test.go 的弃权那一支`

★★**而这两枚的首行都是 `//go:build !windows`** ⇒ ★★**在 Windows runner 上它们**根本不编译****
⇒ ★★**所以"今天红在 Windows runner 上的这两枚"不在票 346 的射程里、也没有别的票认领** ⇒ **无主**。

## ★★为什么这不是"再开一张重复票"

- ★**票 346 治的是 POSIX 那两枚的**弃权支****（★**`!windows` build tag**）。
- ★★**本票治的是 Windows 上**真正在跑、真正在红**的那两枚**（★**`_348_windows_test.go`／`_322_windows_test.go` 一族**）。
- ⇒ ★**两者是**不同平台、不同文件**** ⇒ **⛔ 重复**。

## ★★已知的"为什么是胜率不是性质"（`112-r1` 现跑、我采信）

| 量 | 本机 | Windows runner |
|---|---|---|
| 322 叶子守卫命中 | `guardRefusals=8/8/15`（每发 400、`-count=3`）| ★**0**（400）|
| C26 稳态拒绝（"哪条腿先答"）| 205/400 | ★**288/400** |
| 342 `AC#1` 截断分子 | `truncated=2`（红）| **0**（绿）|

⇒ ★★**`winsec.go:167` 的 `preOpenLeafGuard` 只在 `resolveString` 之后、`pinnedParent` 之前那一段才可能轮到**
⇒ ★**runner 上 C26 把 72% 的尝试先吃掉** ⇒ ★★**红绿是**胜率**、⛔ 性质**。
★**旁证**：★**本机整包三发、三发红名各不相同**（`342r1` 两枚 ＋ `342r2` 一枚）。

## 结案判据

- [ ] **AC#0 现读＋归因（只读，⛔ 动任何产码字节）**：
      ①★**现跑** `grep -l 'go:build'` ／ `//go:build !windows` ⇒ ★**逐枚列出**票 346 射程里那两枚的 build tag**（★⛔ 抄票面）。
      ②★**现跑** Windows 上真正在跑的那两枚（`_348_windows_test.go`／`_322_windows_test.go`）⇒ ★**逐枚给 build tag 与它们在 CI 步骤里的位置**。
      ③★**量出"票 346 改完这两枚还红吗"**（★**读它的判据 ＋ 造一发 POSIX 侧的对照**）。
      ④**正控**：★**在本机跑那两枚** ⇒ ★**必须复现"红名漂移"或"guard=0"** ⇒ **证明"胜率"这件事在本机也能看见**（★**若本机全绿 ⇒ 说清为什么**）。
- [ ] **AC#1 裁形（编排者拍）**：三形各写代价＝
      〔甲〕**给 Windows 那两枚也做确定性化**（★**同票 346 的思路、但落 Windows 侧**）；
      〔乙〕**给 CI 一个"安静台面"**（★**票 346 的甲形**、**要动 `ci.yml`、须机主批 `.github/**`**）；
      〔丙〕**在 Windows runner 上明确 skip 并给出理由**（★**但这会把门挖掉一块** ⇒ ★**⛔ 推荐**）。
      ★**未裁形任何实现都⛔ 开工**。
- [ ] **AC#2 落地＋牙（写腿）**：按裁定那一形做；★**每处改动配一枚今天响、不修不响的用例**；
      ★**反向判据**＝把新加的那一跳摘掉 ⇒ 用例必须**转红并逐名报红行**；三读数 ＋ **还原核**。
      ★★**还原核的时机**：★**必须在**最后一发变异之后**量**（★`112-r1` 实测：在变异之间量到的是**变异产物**）。
- [ ] **AC#3 门禁**：`GOFLAGS= go build ./...`／`go vet ./internal/winsec/`／`go test -count=1 ./internal/winsec/`（**四数分开、⛔ 相减**）／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh`。
      ★**格式门禁甲乙两尺都报**（甲＝入库 blob、**判据**；乙＝工作树、**只作诊断**）。
- [ ] **AC#4 纪律格**：⛔ 翻任何框、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 碰票 346／342 任何原句。

## ★硬提醒（⛔ 踩）

1. ★★**⛔ 把门改绿** —— ★**能红的门才有意义**；★**本票治的是"红得没道理"、⛔ "红"**。
2. ★★**⛔ 动 `//go:build` tag 去把 Windows 那两枚排除掉** —— ★**那是〔丙〕、⛔ 推荐**。
3. ★★**票 342 **未结**（无 `-done`）** ⇒ ★**它的"修前必红"探针还在 CI 门里** ⇒ ★**本票可能也要等 342 落地**（★**`A1196④` 裁：探针搬出 CI、搬到票自己的探针目录**）。

## 排程

`369-a1`（只读＝`AC#0`）→ 编排者裁 `AC#1` → `369-r1`（写腿＝`AC#2`）→ `369-v1`（非实现者裁）。

## Progress log

- [2026-10-23 17:4x +08] agent=`编排者` did=立案（★来路＝`112-r1` §八；★它核到票 346 射程那两枚**首行都是 `//go:build !windows`** ⇒ **Windows runner 上不编译** ⇒ **今天在 Windows 上红的那两枚无主**；★它**只登记、⛔ 认领**、处置对）next=派 `369-a1` 只读腿。
