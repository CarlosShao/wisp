# 票 350 落地腿 `350-r1` —— 有界词表仪器（ban #10）落 `tools/d22scan/**`

- 程＝`worker-ticket350-bounded-wordlist`；日期＝2026-11-06 17:xx–18:xx +08
- 授权＝**票 350 AC#0 owner 批复（2026-11-06，台账 `A1571`：B 批含 350；契约类"350 门禁"按本次批复执行）**——"造门＝人工批准"这一步由机主亲自解扣，本腿凭批复落门。
- 起手锚＝派单点名的 `92d5a33d1`（现量＝**HEAD 的父**：`git merge-base --is-ancestor 92d5a33d1 HEAD` ⇒ 是；HEAD＝`d0a230547`＝A1571 批复入账笔）。
- 隔离（规矩 21 定式，同形先例＝票 140 AC#2(b) 腿）：**主仓 `git status --porcelain -- tools/d22scan` ＝ 0 行**（在飞未跟踪件＝0，未触发 STOP）⇒ **`git clone --shared` → `/d/tmp/wisp/20261106/350r1/repo`**（检出 HEAD `d0a230547`），**全部实现、测试、门禁、提交在克隆里做**；主树工作树零接触、主树共享对象库被 clone 复用（无第二份历史）。

## 0. 交付（commit `79f287ab7`，克隆 dev 分支，⛔ 未推）

| 件 | 是什么 |
|---|---|
| `tools/d22scan/riskwordlist.go`（新） | 纯件三枚：`aWordRe`（`\ba[A-Z][A-Za-z0-9_]*`）、`collectDeclaredIdents`（go/ast 全形声明名册：包级 func/type/const/var＋方法＋字段＋参数＋`:=`/range 局部＋label）、`riskCommentCitations`（`f.Comments` 出处扫描，ParseComments 构造保证"在注释里"——与 ban #9 同一保证，非行正则） |
| `tools/d22scan/main.go` | ＋19 行名册条目 **`10 phantom-symbol`**（`// Bans` 块内，自测名册动态读它）；＋`walkRiskWordlist` 走 `internal/risk`（testdata/`.git`/git-ignored 跳过、`_test.go` 不进、per-package 名册、per-(file,token) 去重）；＋`declaredScopes()` 一行 **live:true**（`examinedKey: "phantom-symbol"`，空树 fatal 归 guard 2 管）＋ scanWithIgnore 挂走 |
| `tools/d22scan/selftest.go` | `selfCoverForTag["phantom-symbol"]="riskGo"` ＋ `selfCounter` 新 case（防 silent 样本读 VACUOUS） |
| `tools/d22scan/selftestsamples.go` | ban #10 双向样本 5 枚（1 ring：`aVanished`＝blacklist.go:193 aEquiv 形状；4 silent：活引用／方法引用／`:=` 局部引用／macOS+"a " 形状边界） |
| `tools/d22scan/riskwordlist_test.go`（新） | 8 枚用例：**起手锚正控**（AC#1①）／**交件侧绿**（AC#1②）／单元（名册广度＋形状边界）／最大 token 与词边界／per-(file,token) 去重／`_test.go` 射程外／空树经台账 fatal／自测映射存在 |
| `tools/d22scan/scan_test.go` | `liveFixture` ＋种 `internal/risk/ok.go`（ban #10 变 live 后"每个 live scope 都干活"规则，同票 88 给 frontend/ 的处理）；**A207/A214 两处分母 3→4 重述**（意图一字未动、每枚新计文件逐名写进注释，见 §4） |

## 1. 现量（两把尺都可逐字重跑；枚数是锚读数 ⛔ 当通则）

**尺①（census.sh，仓外 `/d/tmp/wisp/20261106/350r1/census.sh`）＝生产件（去 `_test.go`）`sed -n 's,^[[:space:]]*//,,p'` 取注释行 → `grep -oE '\ba[A-Z][A-Za-z0-9_]*'` 去重 ＋ 每名声明形对撞：**

| 锚 | 注释侧 token | 声明命中 |
|---|---|---|
| **HEAD `d0a230547`（21 枚产码件）** | `aAnyForm`(1) `aEq`(3) `aRelates`(6) `aRuleReason`(4) `aUnder`(5) ＝ **5 枚** | 逐枚定义 ≥1 ⇒ **门红 0** |
| **起手 `29d258f9`（19 枚产码件）** | 同上 5 枚 ＋ **`aEquiv`(1)** ＝ 6 枚 | **`aEquiv` defhits=0** ⇒ 门红恰 1 |

与票面 §现量、`343-v1`、`349-v1` 三方读数逐枚同（`349-v1` 亦报 5 枚两侧全有定义）。

**尺②（工具自身，交付后 `sh scripts/d22scan.sh` 第二步真扫）＝`ban #10 internal/risk/ examined 21 production Go files`、clean。**

**`_test.go` 排除的诚实注脚（censustest.sh 现量）**：42 枚测试件注释侧 token＝`aEq`(5) `aRuleReason`(1) `aUnder`(6)，逐枚落生产声明（`aKey`/`aPath` 只在测试代码字面出现、不在注释引用侧）⇒ **排除不藏任何真阳**；读数落档，将来若测试注释点名幻影符号要另立票（泛化＝AC#2 ⛔ 本票顺手）。

## 2. AC 对表

| AC | 判据 | 读数 |
|---|---|---|
| AC#0 | 上／不上＝机主拍 | **已拍**＝上（`A1571`）。本腿即"上"路 AC#1 |
| AC#1① | 起手侧必红 1 枚＝`aEquiv` | **PASS**：`TestRiskWordlistStartAnchorControlIsExactlyAEquiv`——把 `29d258f9:internal/risk/blacklist.go` 的**提交字节**（`git show` 出 blob）种进一次性夹具再扫 ⇒ 恰 1 枚 finding、逐字点名 `"aEquiv"`＋`package risk`；同 blob 的 `:193` 行现读＝`// provably expanded fails CLOSED (aEquiv / aUnder report such a miss as a`。⛔ 未动共享树一字——控制跑在夹具里 |
| AC#1② | 交件侧 0 红 | **PASS**：`TestRiskWordlistDeliveryTreeIsGreen`（克隆真树，21 枚产码件入分母、0 finding）＋真扫 `clean` |
| AC#1③ | 误伤面＝0 且可枚举 | **PASS**：HEAD 实测 0；五类可枚举假阳（别包真符号 `hostPathBox`/`assessorFor`、平台散文 `macOS`、标记约定 `DEFERRED`、厂商名 `Nutstore`、故意删的历史名 `CloseScope`）＝`343-v1` 判语所列，本门的**形状与射程把五类结构上拒之门外**——`a[A-Z]` 形必小写头＝不可能别包真符号；`\b` 拒 macOS/DEFERRED（`D` 前是字母）；Nutstore/CloseScope 首二字非 `a[A-Z]`。⛔ 自由文本豁免 0 枚（AC#1 硬判据③） |
| AC#2 | ⛔ 泛化 | 守住：射程＝`internal/risk` 产码注释（`_test.go` 排除有 §1 注脚）；形状＝唯一的 `a[A-Z]…`；名册＝同包 go/ast；**无 allowlist 行、无私有豁免表**；抓 `fires()`/小写形＝不做 |
| AC#3 | 门有牙（两种形状都放行⛔） | **PASS 双向**：ring＝`TestRiskWordlistStartAnchorControlIsExactlyAEquiv`（真历史 blob）＋自测样本（夹具）＋**二进制级 AC#3**（`teeth.sh`：一次性夹具种注释 → `d22scan.exe` rc=1 逐字 `[phantom-symbol] … "aVanished" …`，声明该符号后复跑 rc=0——能打印≠有人管它打不打印，这次是有人管）；silent＝4 枚自测样本＋真树 0 |
| AC#4 | 零既有判据移动 | **PASS＋两处具名重述**：八枚旧禁令射程／豁免语义 0 字节改动（`git diff` 只有新增块与两处注释里的数字）；三枚冻结件／allowlist.txt／golden／thresholds.go 0 触；**任何注释文字 0 改**（含 `internal/risk/**`）。两处例外具名在 §4（分母 pin 3→4，属于"新增 live scope 的机械后果"，意图未动） |
| AC#5 | 颜色与门禁 | **PASS**：`go test ./tools/d22scan -count=1 -timeout 420s`＝**ok 117.8s**（runtests 口径 **PASS=43 FAIL=0 SKIP=0 / RUN=85**，改前 41 PASS ⇒ +2 枚顶层＋子测全绿）；`-self-test`＝**clean, 45 direction checks (21 ring / 24 silent)**（改前 37/20/17 ⇒ 名册动态读＋对，10 phantom-symbol 进名册）；`go build ./...` rc=0（**CGO_ENABLED=1**，未设 0）；gofumpt **v0.12.0 `-l`＝空**（tools/d22scan 与 tools/internal/cmd/scripts/frontend 全仓）；`sh scripts/d22scan.sh`＝正控 OK ＋ 真扫 **clean**；旁证：`internal/risk` 整包 `ok 14.4s`、`internal/winsec` 整包 `ok 45.7s`、`go vet ./cmd/wisp/` rc=0 |
| AC#6 | 纪律格 | 票框未翻一字（翻框归非实现者）；票名／Status／`-done` 未动；写面＝`tools/d22scan/**`＋`docs/evidence/s1/350-*.md`＋票尾一条 bullet；add/commit 显式 pathspec 逐枚手写；⛔ push |

## 3. registry row（票面 AC#1 的"词表件／登记行"要求）

本腿把 gate 以**ban #10** 登进三处既有登记面（全部由现役机器动态读，无新造登记文件）：

1. `tools/d22scan/main.go` `// Bans` 块第 **10** 行 `phantom-symbol`——selftest.go 的名册读取器从这块动态点名（`d22scan -self-test` 头行现读：`10 numbered ban(s) [1 bare-goroutine, …, 9 phantom-citation, 10 phantom-symbol]`）；连号校验（audit rule 0）通过。
2. `tools/d22scan/main.go` `declaredScopes()` 一行 live:true——分母自报（`ban #10 internal/risk/ examined 21`），空树 fatal（guard 2）。
3. `tools/d22scan/selftest.go` `selfCoverForTag` 一行——双向样本的 cover 计数。

**登记依据引批复**：台账 `A1571`（2026-11-06 机主批 owner 清单 B 项含 350；"契约类改动……350 门禁……按本次批复执行"）。

## 4. 既有判据的两次移动（具名，AC#4 的诚实面）

`liveFixture` 是 13 枚既有测试共用的夹具；ban #10 变 live 后"每个 live scope 必须干活"规则（票 71 AC#4／票 88 先例）要求它种一枚 risk 件。后果＝两处把 `internal/` 分母写死的 pin 从 3 变 4：

- `TestWalksSkipGitIgnoredPaths`（A207 形状）：3→4，新计文件＝`risk/ok.go`（clean 件，不被忽略规则覆盖；被排除的两枚 `build/probe.go`/`x.generated.go` 与被捞回的 `pkg/keep.go` 逐名不变）。
- `TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked`（A214/F1 形状）：3→4，同上。

两处均以"RESTATED BY TICKET 350"注释块重述（原意图逐字保留在注释里：哪些文件该被排除／该被读）。**这不是判据放松**：两枚测试断言的是忽略规则的方向，不是夹具文件数；文件数是夹具的机械事实。

## 5. 门禁读数汇总（克隆内现跑，逐条可重跑）

```
git status --porcelain -- tools/d22scan            # 主树起手 = 0 行（不 STOP）
git clone --shared . /d/tmp/wisp/20261106/350r1/repo   # 规矩 21
# 实现后（克隆内）：
go test ./tools/d22scan -count=1 -timeout 420s     # ok 117.810s（runtests 口径 PASS=43 FAIL=0 SKIP=0 RUN=85）
go build ./...                                      # rc=0（CGO_ENABLED=1）
gofumpt -l tools/d22scan                            # 空（v0.12.0）
sh scripts/d22scan.sh                               # 正控 OK + 真扫 clean，分母含 ban #10 internal/risk/=21
d22scan -self-test                                  # clean, 45 direction checks (21 ring / 24 silent)
# 二进制 AC#3 正控（teeth.sh，一次性夹具）：
种注释 → rc=1 [phantom-symbol] … "aVanished" …；声明同名符号 → rc=0
# 起手锚 blob 正控：
29d258f9:internal/risk/blacklist.go → 恰 1 枚，点名 "aEquiv"，:193 行
```

## 6. 本腿没做的（具名残余）

- 翻票 350 六格 AC 框＝非实现者（裁决者≠实现者）；本腿零翻框。
- `ci.yml` 接步（AC#4 明写那一步单独摆机主；本腿⛔ 动 `.github/**`）。
- `_test.go` 射程与"注释里其它可引用形状"（小写调用形等，票 349 名册）＝泛化禁区，另票。
- 台账 `A##` 行＝编排者收件时记；本腿不写台账（写面外）。

## 7. 临时件清单（只建不删，规矩 8）

`/d/tmp/wisp/20261106/350r1/`：`census.sh`／`censustest.sh`（两把普查尺）、`teeth.sh`（二进制正控）、`repo/`（规矩 21 共享克隆，实现与门禁全在里面，commit `79f287ab7`）、`teethfix/`／`anchorfix/`（一次性夹具）、`d22scan.exe`（克隆构建的扫描器）。
