# 369 — 票 346 的射程只到 POSIX：**今天红在 Windows runner 上的那两枚弃权仪器无主**

- Status: **DONE（2026-10-24 13:0x 编排者收口＝格全闭、文件名加 `-done`）** ~~**ready-for-agent**~~
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

- [x] **AC#0 现读＋归因（只读，⛔ 动任何产码字节）**：
      ①★**现跑** `grep -l 'go:build'` ／ `//go:build !windows` ⇒ ★**逐枚列出**票 346 射程里那两枚的 build tag**（★⛔ 抄票面）。
      ②★**现跑** Windows 上真正在跑的那两枚（`_348_windows_test.go`／`_322_windows_test.go`）⇒ ★**逐枚给 build tag 与它们在 CI 步骤里的位置**。
      ③★**量出"票 346 改完这两枚还红吗"**（★**读它的判据 ＋ 造一发 POSIX 侧的对照**）。
      ④**正控**：★**在本机跑那两枚** ⇒ ★**必须复现"红名漂移"或"guard=0"** ⇒ **证明"胜率"这件事在本机也能看见**（★**若本机全绿 ⇒ 说清为什么**）。
- [x] **AC#1 裁形（编排者拍）**：三形各写代价＝
      〔甲〕**给 Windows 那两枚也做确定性化**（★**同票 346 的思路、但落 Windows 侧**）；
      〔乙〕**给 CI 一个"安静台面"**（★**票 346 的甲形**、**要动 `ci.yml`、须机主批 `.github/**`**）；
      〔丙〕**在 Windows runner 上明确 skip 并给出理由**（★**但这会把门挖掉一块** ⇒ ★**⛔ 推荐**）。
      ★**未裁形任何实现都⛔ 开工**。
- [x] **AC#2 落地＋牙（写腿）**：按裁定那一形做；★**每处改动配一枚今天响、不修不响的用例**；
      ★**反向判据**＝把新加的那一跳摘掉 ⇒ 用例必须**转红并逐名报红行**；三读数 ＋ **还原核**。
      ★★**还原核的时机**：★**必须在**最后一发变异之后**量**（★`112-r1` 实测：在变异之间量到的是**变异产物**）。
- [x] **AC#3 门禁**：`GOFLAGS= go build ./...`／`go vet ./internal/winsec/`／`go test -count=1 ./internal/winsec/`（**四数分开、⛔ 相减**）／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh`。
      ★**格式门禁甲乙两尺都报**（甲＝入库 blob、**判据**；乙＝工作树、**只作诊断**）。
- [x] **AC#4 纪律格**：⛔ 翻任何框、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 碰票 346／342 任何原句。

## ★硬提醒（⛔ 踩）

1. ★★**⛔ 把门改绿** —— ★**能红的门才有意义**；★**本票治的是"红得没道理"、⛔ "红"**。
2. ★★**⛔ 动 `//go:build` tag 去把 Windows 那两枚排除掉** —— ★**那是〔丙〕、⛔ 推荐**。
3. ★★**票 342 **未结**（无 `-done`）** ⇒ ★**它的"修前必红"探针还在 CI 门里** ⇒ ★**本票可能也要等 342 落地**（★**`A1196④` 裁：探针搬出 CI、搬到票自己的探针目录**）。

## 排程

`369-a1`（只读＝`AC#0`）→ 编排者裁 `AC#1` → `369-r1`（写腿＝`AC#2`）→ `369-v1`（非实现者裁）。

## Progress log

- [2026-10-23 17:4x +08] agent=`编排者` did=立案（★来路＝`112-r1` §八；★它核到票 346 射程那两枚**首行都是 `//go:build !windows`** ⇒ **Windows runner 上不编译** ⇒ **今天在 Windows 上红的那两枚无主**；★它**只登记、⛔ 认领**、处置对）next=派 `369-a1` 只读腿。

## ★★`AC#2` 判语（2026-10-24 13:0x 编排者，凭写腿 `369-r1` 交件，`7d9e5b481`＋`28b15313d`）

★**`AC#2` 能翻** —— ★**它把两枚胜率判据改成了确定性判据、且牙在**。

**① ★★★判据换成了什么（5 改 1 新、**产码 0 字节**、全 `_test.go`）**
1. ★★**`openleaf_guard_322_windows_test.go`**（`168+/91-`）：★**删掉 400 发 racer 与频率门 `guardSpoke==0||flips==0`**
   ⇒ ★**改成三件确定性**：(a) ★**接线门**新件 `routeguard_369_test.go` 用 `go/parser` **读 `privateFile` 调用图**（`resolveString < preOpenLeafGuard < pinnedParent/openPinnedLeaf`，★**且 guard 那一跳必须在**函数体顶层****）
   (b) ★**直驱 guard 的链接形原样保留** (c) ★**route 枚改成**每发同步** rename 链到调用者名 ＋ `reparseAt` 证明**
   ⇒ ★**一发 `PrivateFile` 必须非 `nil`、归因必须是**两条稳态行之一**、见证不被移除、payload 不落进去** ＋ ★**正控**（链删掉后同一次调用必须落地）
2. ★**`ancestorunlink_348_windows_test.go`**（`88+/17-`）：★**叶子用例改三段** —— (A) **直驱 `openPinnedLeaf`**（bare-sibling 链）必须**逐字带 re-check 的链接形理由 ＋ 两 sentinel**
   (B) **同步 `toFlip`+`reparseAt` 证明世界能进翻态**（替代 `provedReparse==0`）(C) **400 发世界**只判 `deleted>0`**、**split 降为读数**
3. ★**`leafrace_342r2_windows_test.go`**（`26+/1-`）：★**首行 → `//go:build windows && wisp342probe`**

**② ★★★机制（量到行、这是"为什么是胜率"的答案）**
- ★★**`winsec.go:163 resolveString` 的两条稳态腿都把叶子当 prefix 走**（`resolve.go:631` **无条件重跑地板** → `placement_windows.go:59-68`，而 `winsec.go:397-421 pathPieces` **末片就是整条路径**）
  ⇒ ★**叶子一挂链在 `:167` **之前**就被拒**
- ★★**guard（`:167`）与 re-check（`pinnedleaf_342_windows.go:199-211`）只能吃到 [叶子那次 Lstat → 自己那次 CreateFile] 的**两 syscall****
  ★**而 `:179`+`:187` 的窗口**严格包含**它 ⇒ ★★★**`guard 命中 ⊆ 链命中`**
- ★**现量 400 发**：稳态 reparse `249`／链 `path escapes from parent` **`75`**／seal `23`／`resolveString refused` `9`／unreadable `1`／other `6`／nil `37`、★**guard 链接形 `0`、re-check `0`**

**③ ★★★连跑 10 次（这是"确定性"的硬证明）**
- ★**`RUN=40 PASS=40 FAIL=0 SKIP=0`**（四枚 × 10 全 PASS）
- ★★**变异前 `12:33` 与**最后一发变异之后** `12:52` 各一发同读数**
- ★★★**而读数在动**（split `0,1,0,0,1,2,0,0,2,2`）、**判据不动**
  ⇒ ★★**这才是"确定性判据"该有的样子**：★**读数可以动、判决不许动**。

**④ ★★牙在（两发变异仍红）**
- ★★`-overlay` 把 `preOpenLeafGuard` 体改 `return nil` ⇒ ★**`AC#3 RED …accepted a name carrying a reparse point`**（`:188/:190/:192`）**FAIL**
- ★★去掉 `openPinnedLeaf` 的身份 re-check ⇒ ★**`AC#2 RED: …returned a usable handle…`**（`:153`）**FAIL**

**⑤ ★正控**：★**链删掉后同一次 `PrivateFile` 落地、payload 到位**；★**「真文件/未建名必须过」两跳 `10/10` 绿**

**⑥ ★★★342 探针：移了 —— 而它用的办法比我建议的更好**
- ★**我派单说"移到票 342 自己的 probe 目录"**；★★**它用了 build tag**：★**`//go:build windows && wisp342probe`**、★**文件原地、不进默认构建**
- ★★**验证**：★`go test -tags wisp342probe -run '^TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2$' ./internal/winsec/` ⇒ **PASS**、`re-check answered 45 of 400`
  ★**默认构建 `-list 342r2` 不再列它**
- ⇒ ★★★**build tag 比搬文件好**：★**不破坏路径预算、不搬文件、不新增 roster 条目** ⇒ **我认这个偏离**。

**⑦ ★★★反向判据三读数**
1. ★**HEAD 原件 `-count=10` ⇒ `PASS=18 FAIL=2`**（★348 枚 shot 1、8 红）
2. ★★★**只把 racer pace `2ms→20ms`（判据**一字节不改**）⇒ `PASS=3 FAIL=2`**、逐发 `guardRefusals=0,0,1,1,1`、`flips=14,12,15,14,11`
   ⇒ ★★★**"颜色是时钟的函数"** —— ★**这是本会话对"胜率判据"最硬的一次证词**
3. ★**接线门对三枚已知坏件分别读成 `MISSING`／`nested`／`sameOrder=false`**
4. ★**还原核在最后一发变异之后**：★**10× 重跑 `40/40`**

**⑧ ★四数／门禁**
- ★**改前** `-count=1 ./internal/winsec/` ⇒ **`RUN=127 PASS=76 FAIL=0 SKIP=0` `rc=0`**（★**同一棵树 `369-a1` 那发是 `FAIL=1`**、红名 `…BareSiblingForm342r2` ⇒ **门级也在掷硬币**）
- ★**改后** ⇒ **`RUN=127 PASS=76 FAIL=0 SKIP=0` `rc=0`**（★**−1 探针 +1 接线门**、`ok 61.175s`）
- ★**门禁**：`go build ./... rc=0`／`vet`（含 `-tags wisp342probe`）`rc=0`／`d22scan.sh rc=0 clean`／`check-path-length-budget.sh rc=0 GREEN`／★**格式甲乙两尺都空**／★**`GOOS=linux` 编译通过**
- ★**纪律**：⛔ 动 `winsec.go`、★**开工时写面内 `376-a1` 的未提交件一枚未碰**、收尾写面已空

**⑨ ★★它两条要处置的上报 ＋ 我的裁**
1. ★★★**"CI 门仍会红，但红的是另一枚、且是真红"** ⇒ ★**它问"要不要再移这枚"**
   ⇒ ★★★**我的裁：不移、留。** 理由三条：
   - ★★**它测的是**真伤害****（`ancestoropen_342_windows_test.go:50` 的 `truncated=2/400`、★**外国对象 body 真被换掉**）⇒ **那是真缺陷**
   - ★★**`369` 的责任是"把胜率判据改确定性"** ⇒ ★**这枚**不是胜率**、是真红** ⇒ **⛔ 本票的事**
   - ★★★**所以门**应该**红** ⇒ ★**记成**票 342 的阻塞项**（★`A1229④` 同判）
   ★★**它按硬提醒 #1「⛔ 把门改绿」没动它** ⇒ **处置对**。
2. ★★**"接线门读盘、不吃 `-overlay`"** ⇒ ★**实测三发换 `winsec.go` 的 overlay 它全 PASS**
   ⇒ ★★★**它自己写清了射程**："**它抓"有人删/挪调用"、抓不到"overlay 里的产码"；其牙是对已知坏件演示的**"
   ⇒ ★★**这是诚实的自限**、★**要写进件、⛔ 假装它是全称** ⇒ **记、下一波可补**。
