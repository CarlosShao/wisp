# 票 350 非实现者终裁腿 `350-verdict` —— ban #10 有界词表门（七格读数＋closeable 判语）

- 腿名 `350-verdict`；日期 2026-11-06 18:0x–18:2x +08；角色＝**非实现者**（实现者程＝`worker-ticket350-bounded-wordlist`，其提交 `747a38edc`＋`58d86064a` 均已核为 HEAD 祖先且已在 `origin/dev`）。
- 终裁起手 HEAD＝`218ed2ee9`（A1577 收 wave35 那笔；`git merge-base --is-ancestor 747a38edc HEAD` ⇒ YES，`58d86064a` 同 ⇒ YES；`git branch -r --contains` 两枚均在 `origin/dev`）。
- 本腿写面＝**仅本件**（`docs/evidence/s1/350-verdict.md`）；⛔ 票面、⛔ 产码、⛔ `tools/d22scan/**`、⛔ 台账、⛔ push。cmd/wisp 那 2 枚未跟踪票-323 测试件全程未碰未跑。
- 环境卫生（每条 shell 现场带）：`unset GOTMPDIR TMPDIR TEMP TMP`；`CGO_ENABLED` 现读＝**1**（`go env CGO_ENABLED`，全程未设 0）；go1.27.1 windows/amd64；gofumpt **v0.12.0**。

## 0. 独立 fresh 复跑读数（本腿现跑，非转抄；每条附可重跑尺）

| # | 门禁 | 尺（逐字） | 读数 |
|---|---|---|---|
| 1 | d22scan 模块测试 | `cd tools/d22scan && go test . -count=1 -timeout 420s` | **ok 215.304s，rc=0** |
| 2 | 根 `go build ./...` | `go build ./...` | **rc=1，红全部且仅来自 `build\mut_rvw.go`**（`git check-ignore -v` ⇒ `.gitignore:14:build/`；`git status --porcelain -- build`＝0 行＝gitignored 未跟踪残渣，W21 突变台件族，mtime 2026-11-04）——**按派单排除**；对照尺 `go build ./cmd/... ./internal/... ./tools/...` ⇒ **rc=0** |
| 3 | gofumpt | `gofumpt -l tools/d22scan` | **空**。旁证（全仓 `-l .`）：非空行全部落 `.scratch/**` 探针夹具与 gitignored `build/probe61r2/`，**tracked 件 0 枚** |
| 4 | 正控＋真扫 | `sh scripts/d22scan.sh` | 正控 `runtests.sh: OK - packages=[./...] top-level: PASS=43 FAIL=0 SKIP=0, === RUN=85`；真扫 **clean**；分母现读 **`ban #10 internal/risk/ examined 21 production Go files, comment symbol citations`**；rc=0 |
| 5 | 自测 | `d22scan.exe -self-test`（在包目录内跑） | **clean - all 45 direction checks passed (21 expect-ring, 24 expect-silent)**；名册头行 `10 numbered ban(s) […, 10 phantom-symbol]`；ban #10 五行全 OK（1 ring＋4 silent） |

⚠ 5 的运行契约：`-self-test` 从 `main.go` 所在包目录读名册源文件——在别的 CWD 直接跑二进制会 `FATAL cannot read the roster source "main.go"`（本腿实测撞过一次）。这是调用约定不是缺陷；CI 的步（ci.yml:144-145 `go run . -self-test`＋`working-directory: tools/d22scan`）满足它。

## 1. AC#3 二进制级双向突变（本腿在规矩 21 克隆独立重演，共享树零编辑）

- 克隆：`git clone --shared . /d/tmp/wisp/20261106/350verdict/repo` ⇒ 检出 `218ed2ee9`，`git status --porcelain` 干净起手。
- 从克隆构建 `d22scan.exe`（rc=0）。
- **RED**：向克隆种 `internal/risk/verdict350probe.go`（注释点名 `aPhantomVerdict350`、无任何声明）⇒ 扫描 **rc=1**，逐字：
  `internal/risk/verdict350probe.go:3: [phantom-symbol] comment names "aPhantomVerdict350", which no declaration in package risk defines (ticket 350 ban #10): …`
- **GREEN**：同文件补 `func aPhantomVerdict350() {}` ⇒ 复跑 **rc=0 clean**。
- 退出码是真门：1＝有 finding 拒绝，0＝修复后放行——「能打印≠有人管它打不打印」这一关，这次是**有人管**（exit code 消费者＝CI）。

### 1.1 AC#1① 起手锚正控（本腿独立重演，不经实现者的测试函数）

- 夹具（仓外 `/d/tmp/wisp/20261106/350verdict/anchor/`，零共享树编辑）：`git show 29d258f9:internal/risk/blacklist.go` 落盘为唯一 risk 件；`:193` 行现读＝`// provably expanded fails CLOSED (aEquiv / aUnder report such a miss as a`；blob 内 `aEquiv` 出现次数＝**1**；再加 gateFixture 最小面＋8 枚 filler 过 `checkRoot` 的 ≥10 产码件防空转闸（首次跑 rc=2 被这枚闸拦下＝**防空转闸本身有牙**，顺带验证）。
- 扫描 ⇒ **恰 1 枚 finding，逐字点名**：
  `internal/risk/blacklist.go:193: [phantom-symbol] comment names "aEquiv", which no declaration in package risk defines (ticket 350 ban #10)…`，rc=1。
- 与票面「起手侧红 1 枚＝aEquiv＝真阳、假阳 0」逐字对上。交件侧 0 红另有三路独立读数：① go test 里的 `TestRiskWordlistDeliveryTreeIsGreen`；② d22scan.sh 真扫 clean；③ 本腿克隆真扫 clean（examined 21/22，视克隆侧是否含本腿探针件而定）。

## 2. 逐 AC 判语

| AC | 判语 | 凭据（本腿现读） |
|---|---|---|
| **AC#0**（上/不上＝机主拍） | **PASS** | 台账 `A1571`（pending-and-issues.md:31347-31352）现读：B 批含 **350** 全批；且「契约类改动（… **350 门禁** …）按本次批复执行」逐字在账。「造门＝人工批准」这步确由机主解扣 |
| **AC#1①**（起手侧恰红 1＝aEquiv） | **PASS** | §1.1 独立重演；`TestRiskWordlistStartAnchorControlIsExactlyAEquiv` 在 1 的套件里同绿 |
| **AC#1②**（交件侧 0 红） | **PASS** | §0 表 4/5＋§1.1③ 三路独立 clean；分母 21 枚产码件非零（无 VACUOUS 读法） |
| **AC#1③**（误伤面＝0 且可枚举） | **PASS** | 交件侧实测 0；五类历史假阳（别包真符号/平台散文/标记约定/厂商名/历史名）被形状**结构性**拒之门外：`a[A-Z]` 必小写头＝不可能他包符号、`\b` 拒 macOS/DEFERRED、首二字拒 Nutstore/CloseScope。**allowlist 0 行**：`tools/d22scan/allowlist.txt` 现读 5 行全是旧禁令的，无一 `phantom-symbol`，且 `747a38edc` 触及清单（7 件）不含它（最后一次被碰＝`38b37153b` 票 70）；自由文本豁免 0 枚（`riskwordlist.go` 全文无豁免表） |
| **AC#2**（⛔ 泛化） | **PASS** | 唯一调用点 `main.go:290`＝`walkRiskWordlist(filepath.Join(root,"internal","risk"))`，射程写死单树；形状唯一（`aWordRe`）；名册＝同包 go/ast 全形声明；`_test.go` 排除有现量注脚（evidence §1，42 枚测试件注释 token 逐枚落生产声明，排除不藏真阳）；`fires()`/小写调用形**没做**（归票 349） |
| **AC#3**（门有牙） | **PASS** | §1 二进制双向（rc=1→rc=0）＋ring 正控两枚（自测样本＋起手 blob）＋4 枚 silent 自测样本全 OK（§0 表 5） |
| **AC#4**（零既有判据移动） | **PASS（核读）＋CI 接步 owner 面（见 §3）** | `git show 747a38edc` 逐行核：八条旧禁令射程/豁免 0 字节（diff 全为新增块）；既有行仅两处分母 pin 3→4，各带「RESTATED BY TICKET 350」注释块、原意图逐字保留——性质＝ban #10 变 live 后「每个 live scope 必须干活」既有规则（票 88 先例）的**机械后果**，不是判据放松；三枚冻结件/golden/thresholds.go/`internal/risk/**` 任何注释 0 改（触及清单仅 tools/d22scan 6 件＋evidence 件） |
| **AC#5**（颜色与门禁） | **PASS** | §0 四读数全绿；车道纪律：本腿未跑根 `go test ./...`、未跑 `./cmd/wisp` 整包；并现证根模块模式 `go test ./tools/d22scan` 不可解（`main module does not contain package`，rc=1）——派单「必须 cd 进模块跑」这条是承重的，不是仪式 |
| **AC#6**（纪律格，观察项） | **合规（观察）** | 票面除允许的 progress-log 一条 bullet（`58d86064a`，2 行）外 0 触；无 `-done`；未 push（两枚提交由编排者收入 origin/dev，符合「腿只 commit 不 push」）；翻框归收件编排者，本腿与实现腿都零翻框 |

## 3. AC#4 的 CI 接步——owner 面，本腿不判

- 票面 AC#4 明写：`.github/**`/`ci.yml` 若需接进步骤＝**单独摆机主拍**。本腿按派单只记不裁。
- 现状现读：`ci.yml:111-174` 已有三步——正控（runtests.sh）、`-self-test`、`sh scripts/d22scan.sh`。这三步都是**数据驱动**（self-test 名册从 `main.go` 现读、runtests 跑全包测试），ban #10 上车**不需要任何 yml 改动**，随既有步进 CI。这是从代码结构推出的**推断**，不是证据：Linux runner 上的实测读数本腿给不了（本地 Windows verdict ≠ runner verdict）。摆机主的那一问因此缩成一句：**是否需要在 ci.yml 里为 ban #10 加显式断言行；若不加，默认随既有三步 ride-along**。

## 4. 实现者证据件（350-bounded-wordlist-gate.md）的两处笔误——具名，不动 AC

1. **「改前 41 ⇒ +2」错，实为 35 ⇒ +8**。本腿基线：`git archive 747a38edc^ | tar -x` 到仓外目录现跑 `runtests.sh` ⇒ **PASS=35 FAIL=0 SKIP=0, === RUN=77**；HEAD 侧 PASS=43 / RUN=85（§0 表 4）。计数器对撞：pre＝35 枚顶层 `func Test`（scan_test 31＋selftest_test 4）、HEAD＝43 枚（＋riskwordlist_test 8）——**+8 恰是 riskwordlist_test.go 的 8 枚新测试**，与「PASS 计顶层 --- PASS 行」的 runtests 口径（runtests.sh:222）自洽。证据件的 41/43 对不上任何口径；正确读数＝**35/43**。纯文档 erratum，方向没变（新增测试全绿、无既有红）。
2. **§5 命令块的 `go test ./tools/d22scan` 形式写法不成立**：该 pattern 从主（或克隆）仓根解析失败（本腿现证 rc=1 `does not contain package`）；真实车道＝`cd tools/d22scan && go test .`（＝派单给的形式）。读数本身（ok 117.8s 那条）与事实不矛盾，只是命令按字面复跑不通；本腿同门禁 fresh 读数＝ok 215.304s（§0 表 1）。
3. 旁注（非笔误）：证据件 §5 的 `go build ./... rc=0` 在其克隆语境成立（克隆不含主树 gitignored 残渣）；主树上根 `./...` 红仅来自 `build/mut_rvw.go`（§0 表 2），与派单预注一致，已排除。

## 5. closeable 判语

- **票 350＝PASS，closeable＝YES**：AC#0/#1(①②③)/#2/#3/#5 六读数全 PASS；AC#4 的「零判据移动」半核读 PASS，CI 接步半按票面归 owner（§3，一问已缩好）；AC#6 观察合规。
- 收口动作（翻七格、`-done`、台账 A 行）归编排者收件时做——本腿⛔ 票面。
- **撤销口令**：「撤 350-verdict」⇒ 删本裁决件、终裁退回未裁；实现提交 `747a38edc`/`58d86064a` 与 origin/dev 状态不受影响。

## 6. 本腿临时件清单（只建不删，规矩 8；全部在仓外）

`/d/tmp/wisp/20261106/350verdict/`：`repo/`（规矩 21 共享克隆 @218ed2ee9，内含本腿一次性探针件 `internal/risk/verdict350probe.go`，未 commit）、`anchor/`（起手锚夹具＋8 filler）、`d22scan.exe`（克隆构建）、`pre350/`（747a38edc^ 基线树，§4.1 基线读数用）。
