# 票 318 写腿 318-r1 —— AC#1 落地：`WideningRule` 入口**整卡拒收**含标记 path（选案①）

- 腿身份＝实现者（判语归非实现者腿）；⛔ 未翻框、⛔ 未动 `Status:`/票名；⛔ 未 push。
- 锚（起手现读）＝`5893baaa`；读数窗口 2026-10-14 **16:50–17:0x +08**。
- 环境：`PATH=…/third_party/sherpa-onnx:$PATH`（MSYS 形）；`GOFLAGS=`；`GOTMPDIR=TMP=TEMP='C:\tmp\wisp318r1gotmp'`（反斜杠形；首跑未建该目录时 go 拒绝起跑 `go: creating work dir … cannot find the file specified`，先 `mkdir` 后重跑——该首发日志被 10 覆盖、未留档，具名）；判据全用 `-run` 小射程 `-count=1`；⛔ 未跑整包 `./cmd/wisp/`（314 死区）。
- 结论先行：**① 选案① 落地**（`replies.go` 产码 5 行＝入口整卡拒收守卫）；**② 修前红→修后绿逐字复到**（三形：正斜杠/反斜杠/空路径标记）；**③ 三发突变逐字点名红**（摘守卫 M1／换判词 M3／换逐 path 跳过 M2——最后一个只红"整卡拒"那条钉）；**④ 对照不回归**（无标记越界读仍推真父目录，探针真链＋单元双面）；**⑤ 门禁全 rc=0**（§7；唯一 rc=1 是 session 全包的**既有间歇红**，具名在册＋独立控制）。

## 1. 选案（先读再裁；引文逐字）

现读四处先例（全部真读）：

- `internal/session/grants.go:260-268` `recordablePattern`：`strings.Contains(p, "无法规范化")` ⇒ 拒这一行；`Record:150` 拒整笔记账。
- `internal/session/grants.go:189-197` `Covering`：**任一** path 不可记账 ⇒ 整个覆盖问题答 `(0,false)`（"puts the question back on the user"）。
- `internal/tools/bridge.go:1273` `scopeVerdictOf`：**任一** path 带标记 ⇒ 整个判定 `scopeOut`。
- `internal/tools/bridge.go:1028` `displayPaths`：标记＝原样路径＋`" (无法规范化: "+err.Error()+")"`，`Paths` 里就这一种形；C26 今天拒 traversol 的唯一文本＝`ErrReparseDenied`（`internal/risk/pathresolver.go:36`，**自身带一个 `/`**）。

**两候选裁决**：

| 候选 | 裁决 | 理由（现读） |
|---|---|---|
| ① `WideningRule` 入口拒收含标记 path | **采用** | 直落 `grants.go:264` 口径；面只有本包一枚文件；标记 ⇒ R2 fail-closed ⇒ 卡必 L2（`rules_gateway.go:36-44`；探针 READING-1 `level=L2 rules=[R2 R3]`）⇒ `WideningRule` 是**唯一不认标记的消费方**，此处拒＝链条闭合在唯一可达点。 |
| ② 让标记不进 `Paths`（改 `displayPaths`） | 否（未触、无需停报） | 动 `tools.Decision.Paths` 语义面（票面明示须先报）；且会把"C26 拒判的路径"伪装成正规 canonical 路径喂给 session/bridge/审计全部下游；另按 `dirOf` 反推它会推出**真父目录 `…/allowed/jn`** 并存成放宽行——把碎片换成另一种"基于拒判路径的放宽"，方向仍错。 |

**同一采用的第二裁决（整卡拒 vs 逐 path 跳过）**：`WideningRule` 的契约是"恰好一个目录才给一行，更宽一律拒（`replies.go:490-495` 自述）"。逐 path 跳过会让"修前拒、修后有"的卡**反向开口**（混合卡从 `ok=false` 变推干净目录），与本闸门"只收窄"方向相反；整卡拒与 `Covering:193-197`、`scopeVerdictOf:1273` 的 any-path 形同构。突变 **M2** 即此对照读数。
⇒ 未触 `cmd/wisp/**`、未触 `tools.Decision.Paths` ⇒ **不触发停手条件**（回执点具名）。

## 2. 落修（产码 5 行＋注释）

`internal/agent/approval/replies.go` `WideningRule` 入口：

```go
	for _, p := range c.Paths {
		if strings.Contains(p, "无法规范化") {
			return "", "", false
		}
	}
```

＋函数文档 16 行（记根因、记两先例、记"跳过＝反向开口"的排除）。总 diff：`replies.go +21/-0`（产码 5 行）；`zz_201_a1…_test.go +56/-72`（探针翻钉，见 §3）；新尺 `ticket318_widening_marker_rulers_test.go` 129 行。`cmd/wisp/**`、`tools/**`、`config/**`、`grants.go` **零字节改动**；`replies_201_test.go` 一字未动（既有断言零放宽，`git status` 只挂 3 枚文件）。

## 3. 修前红 / 修后绿 / 探针翻钉 逐字

- **修前红**（锚 `5893baaa` 未修树；`go test ./internal/agent/approval/ -run 'TestTicket318|TestZZ201A1' -count=1 -v`，16:50:22；`logs/10-prefix-red.txt`）⇒ **rc=1**：
  - `--- FAIL: TestTicket318AMarkerCardOffersNoWideningTarget (0.00s)`，逐字红句三形：
    `…:58: junction card as 201-a1 measured it (forward slashes): marked path still offered as a widening base: dir="C:/proj/allowed/jn/victim.txt (无法规范化: risk: path traverses a reparse point (junction" rule="[fs] allowed_dirs += \"…(junction\""`；
    `…:58: the same chain in backslash spelling: … dir="…(junction"`；
    `…:58: empty-path marker (tools: empty path), no slash in the string: … dir=" (无法规范化: tools: empty path)" rule="[fs] allowed_dirs += \" (无法规范化: tools: empty path)\""`
    ⇒ **AC#2 第一格现量**：空路径标记的"目录"＝**整条标记**（无斜杠 ⇒ `dirOf` 原样返回）。
  - `--- FAIL: TestZZ201A1JunctionMarkerCardIsRefusedForWidening (0.03s)`，真 NTFS junction 链停在翻后的 2b：`…:191: 反斜杠拼法：带标记卡面仍被当作放宽基础… dir="C:/tmp/wisp318r1gotmp/TestZZ…/001/allowed/jn/victim.txt (无法规范化: …(junction" ok=true`。
  - `FAIL github.com/CarlosShao/wisp/internal/agent/approval 1.265s`。
- **修后绿**（同命令，终树 16:57:31；`logs/23-final-scoped.txt`）⇒ **rc=0**，5 枚全 PASS（`ok … 1.353s`），真链逐字：
  - `READING-2 卡面逐字 Paths[0] = "C:/tmp/…/allowed/jn/victim.txt (无法规范化: risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions)"`（**标记仍在卡面**，本修不动卡面）；
  - `READING-2b 反斜杠拼法: … dir="" ok=false`（**两种拼法同拒**）；
  - `READING-3 ledger 卡: level=L2 tool=fs.read paths=[…同上逐字]`；`READING-4 WideningRule: ok=false`；`READING-5 dir  = ""`；`READING-6 rule = ""`；
  - `READING-11 对照（无标记的越界读）: … dir="C:/tmp/…/001/outside" rule="[fs] allowed_dirs += \"C:/tmp/…/001/outside\"" ok=true`（**真父目录**）；
  - `READING-12 结论：带标记卡面在入口整卡拒（两种拼法 + 空路径标记同拒），干净卡面仍推真父目录`。
- **探针翻钉**（`zz_201_a1_junction_marker_windows_test.go`，跟随 `315-r1` 先例：zz 探针在修落地后翻成常驻回归钉）：头注改现状说明；测试名 `TestZZ201A1JunctionMarkerBecomesAWideningFragment` → `TestZZ201A1JunctionMarkerCardIsRefusedForWidening`（`TestZZ201A1` 前缀保留，a1 复跑尺 `-run TestZZ201A1` 依旧命中）；legs 2b/3 断言翻向（修前钉 `ok=true`＋碎片尾，修后钉 `ok=false`＋`dir/rule` 清空）；**legs 4-5（碎片进 config.toml＋第二张卡）按"主题不存在"删除**——修前读数原文保留在 `probes/201/a1/logs/run3-full.txt` READING-8/9/10，不重写；对照腿加强（`EqualFold(cdir, 真父目录)`＋`cok==true`）。⛔ 未动 `probes/201/a1/00-verdict.md`（历史件一字不改）。

## 4. 突变（三发，全部 `-overlay`，工作树零突变；JSON 入库、突变 .go 留盘）

| 突变 | 读了什么 | rc | 逐字点名 |
|---|---|---|---|
| **M1 摘掉守卫**（`overlay-mut-noguard.json`） | 与修前同族全红 | 1 | `--- FAIL: TestTicket318AMarkerCardOffersNoWideningTarget`＋`…:58 … dir="…(junction" …`；`--- FAIL: TestZZ201A1JunctionMarkerCardIsRefusedForWidening`（`…:191 … ok=true`）（`logs/30`） |
| **M2 换判词＝逐 path 跳过**（`overlay-mut-perpathskip.json`） | **只有"整卡拒"那条钉红**；标记单形钉与探针全绿 | 1 | `ticket318_widening_marker_rulers_test.go:84: one marked path must refuse the whole card, got dir="C:/proj/outside" rule="[fs] allowed_dirs += \"C:/proj/outside\""` ⇒ 证明 `TestTicket318AMarkedPathRefusesTheWholeCard` 对"跳过"替代语义有牙（`logs/31`） |
| **M3 换判词＝谓词换字**（`overlay-mut-otherword.json`） | 与 M1 同两钉红 | 1 | 同 M1 两句（`logs/32`） |

- 重建配方（突变 .go 不追库，留盘；对终树现读逐字替换即可）：M1＝删 `WideningRule` 入口 5 行 for 块；M2＝删该块并把 `if strings.Contains(p, "无法规范化") { continue }` 插进派生循环首行；M3＝块内 `"无法规范化"` → `"无法规范化X"`。

## 5. 对照（不回归两面）

- **真链面**（探针，真 junction→真 bridge→真 Gate/Replies）：READING-11 见 §3——无标记越界读 `ok=true` 且 `dir` 逐字＝真父目录（终树断言 `EqualFold`）。
- **单元面**：`TestTicket318TheSamePathWithoutTheMarkerStillDerivesItsParent`——`C:/proj/allowed/jn/victim.txt`（同路径**去掉标记**）⇒ `dir="C:/proj/allowed/jn"`、rule 逐字；修前/修后双绿（它是控制，不是红钉）。
- **生产侧本闸门邻域**：`go test ./cmd/wisp/ -run 'TestAlwaysBranch' -count=1 -v` ⇒ rc=0，`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card` PASS 3.70s＋`…StoresNothingWhenTheSecondCardIsRefused` PASS 3.42s（`logs/44`；本波唯一涉窗笼跑＝控制台 harness、无系统窗口；⛔ 未跑整包）。

## 6. AC#2 边界逐条

- **空路径标记**＝**已量**：修前 `dirOf` 把**整条标记**当目录（`dir=" (无法规范化: tools: empty path)"`、`ok=true`，§3 红句逐字）；修后拒。尺：`ticket318_…:第三形`。（生产大概率到不了这形——fs 参数校验先拒——但卡面若现即拒，两边都覆盖。）
- **8.3 短名**＝**不属本修，已钉反证**：短名不是规范化失败（C26 折叠成短名/长名，252 的 `paths.go:139-155`），永远产不出标记 ⇒ 守卫与它不相干；`TestTicket318ShortNamesAreNotTheMarker` 钉住"无标记→原行为"，`internal/tools` 短名族 3 枚绿（§7）。
- **多路径 `ok=false`**＝**已量**：`[标记, 干净]` 混合卡 ⇒ 整卡拒（修后入口；修前也已因两 distinct 目录拒）＋`[标记, 标记]` 全标记 ⇒ 拒。其与"逐 path 跳过"的差别由 **M2** 逐字点名（跳过语义下混合卡变 `ok=true dir="C:/proj/outside"`）。
- 另一枚 `dirOf` 消费方 `Aggregate`（`batch.go:45`）＝**具名不属本修**：它只在 L1 走（`batch.go:36`），而标记 ⇒ R2 fail-closed ⇒ L2（现读 `rules_gateway.go:36-44`＋探针 READING-1 `rules=[R2 R3]`）⇒ 标记进不了它的射程。
- **沿既有口径的已知假阳面（登记）**：真文件名恰含四字 `无法规范化` 的卡也会被拒——`grants.go:264`/`bridge.go:1273` 同判，本修**不新造更细判据**。

## 7. 门禁 rc 表（终树；日志见 §9）

| 门 | 命令（要点） | rc | 读数 |
|---|---|---|---|
| 钉子（修后） | `go test ./internal/agent/approval/ -run 'TestTicket318\|TestZZ201A1' -count=1 -v` | 0 | 5 PASS；`ok 1.353s`（`logs/23`） |
| approval 全包 | `go test ./internal/agent/approval/ -count=1` | 0 | `ok 1.955s`（`logs/24`，63B） |
| 钉子 -race | 同 `-run` ＋`-race` | 0 | `ok 3.290s`（`logs/22`，63B） |
| session grants 先例族 | `go test ./internal/session/ -run 'TestTicket224(Mint\|Record\|Covering\|Expired\|LedgerFailsClosed\|NewLedger\|PatternCovers\|Wildcard\|TestLiterals)' -count=1 -v` | 0 | 11 PASS；`ok 1.542s`（`logs/42`） |
| session **全包** | `go test ./internal/session/ -count=1` | **1** | 唯一红＝`TestSessionListByWorkspace`（**既有间歇**）：台账 `pending-and-issues.md:17402`「新量…间歇红（安静 5 发里 2）→ 交回 166 侧」、`:17760/:17762` 同款在册；**决定性控制**＝`go list -deps ./internal/session` 数 `agent/approval` 依赖 **0**（本修物理不可达）＋本腿 `-run …ListByWorkspace -count=3` 实测 **2 红 1 绿**（`logs/40`/`logs/41`） |
| tools junction/短名族 | `go test ./internal/tools/ -run 'Junction\|ShortName' -count=1 -v` | 0 | 11 PASS；`ok 3.282s`（`logs/43`） |
| cmd/wisp always 两枚（小射程） | `go test ./cmd/wisp/ -run 'TestAlwaysBranch' -count=1 -v` | 0 | 2 PASS（3.70s/3.42s）；`ok 9.336s`（`logs/44`） |
| 构建 | `go build ./...` | 0 | 空输出（`logs/50`，0B） |
| vet | `go vet ./internal/agent/approval/` | 0 | 空输出（`logs/51`，0B） |
| 格式 | `gofumpt v0.12.0 -l` 三枚动过的件 | 0 | 空（`logs/52`，0B） |
| D22 扫描 | `sh scripts/d22scan.sh` | 0 | `clean - no D22 ban violations`（`logs/53`） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | 0 | `VERDICT GREEN`；over=60／roster=60（60=60），wall interval=0（`logs/54`） |

## 8. 具名没做

- ⛔ 未动 `cmd/wisp/**` 产码（`approval_always.go` 的 `ok=false` 拒绝支是**现成**的，本修只是让标记卡走到它）、⛔ 未动 `tools.Decision.Paths`/`displayPaths`、`config/**`、`grants.go`、`thresholds.go`、golden、`allowlist.txt`、`PLAN.md`/`specs/**`、`frontend/**`/`design/**`、`scripts/**`。
- ⛔ 未端到端真跑 `wisp run` 遇真 junction（无法在 cmd/wisp 面造卡；见 §5 生产侧替代读数）——**登记**：票面「最后半跳下游（惰性行）」自本修起**不可达**（`AddAllowedDir` 在 `ok=false` 分支之前折返）。
- ⛔ 未跑整包 `cmd/wisp`（314 死区，派单明令）；⛔ 未修 session 间歇红（在册、他面）。
- ⛔ 未删任何临时件；三枚突变 `.go` 副本留盘（未追库，重建配方在 §4）；⛔ 未 push、未 `add -A`。

## 9. 读数名册（`.scratch/wisp/probes/318/r1/`；★＝追库）

| 件 | 字节 | 是什么 |
|---|---|---|
| `logs/10-prefix-red.txt` ★ | 4,083 | 修前红：三形红句＋翻后探针 2b 红＋`FAIL 1.265s` |
| `logs/20-postfix-green.txt` | 3,095 | 修后绿（注释微调前同读数，留盘） |
| `logs/21-postfix-package-full.txt` | 63 | approval 全包 `ok 1.995s`（注释微调前） |
| `logs/22-postfix-race.txt` | 63 | `-race` 小射程 `ok 3.290s` |
| `logs/23-final-scoped.txt` ★ | 3,103 | **终树**修后绿＋对照＋READING-2…12 全逐字 |
| `logs/24-final-package-full.txt` | 63 | 终树 approval 全包 `ok 1.955s` |
| `logs/30-mut-noguard-red.txt` ★ | 4,079 | M1 摘守卫：两钉逐字红 |
| `logs/31-mut-perpathskip-red.txt` ★ | 3,275 | M2 逐 path 跳过：只"整卡拒"钉红（:84 逐字） |
| `logs/32-mut-otherword-red.txt` ★ | 4,083 | M3 换判词：两钉逐字红 |
| `logs/40-session-full.txt` | 33,515 | session 全包 rc=1＝既有间歇（`manage_test.go:78` 排序）大日志留盘 |
| `logs/41-session-scoped-rerun.txt` ★ | 6,024 | 同枚 `-count=3`：2 红 1 绿（间歇实锤） |
| `logs/42-session-grant-family.txt` ★ | 15,569 | grants 先例族 11 PASS |
| `logs/43-tools-junction-family.txt` ★ | 5,234 | tools junction/短名族 11 PASS |
| `logs/44-cmdwisp-always-branch.txt` ★ | 6,452 | cmd/wisp always 两枚 PASS（生产侧邻域） |
| `logs/50-build.txt`／`51-vet-approval.txt`／`52-gofumpt.txt` | 0／0／0 | rc=0 空输出（0B 即读数，留盘） |
| `logs/53-d22scan.txt` | 23,760 | d22scan clean（大日志留盘） |
| `logs/54-pathlen-before-commit.txt` | 1,458 | 长度门 GREEN 60=60 |
| `mutation/overlay-mut-{noguard,perpathskip,otherword}.json` ★ | 179×3 | overlay 映射（突变 .go 留盘，不追库） |

★＝追库；commit B 的 `--name-only` 即上表 ★ 名单（10/23/30/31/32/41/42/43/44＋3 JSON＋本 md）。
