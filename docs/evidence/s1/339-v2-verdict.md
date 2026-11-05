# 票 339 第 7 格（`AC#2`）裁决 · `339-v2`（非实现者裁决腿·只读产码·只写本件）

裁决者＝`339-v2`（⛔ 实现者、⛔ 产码、⛔ 翻框、⛔ 改票、⛔ push）。起手锚＝`a2b3b80c97`（`git rev-parse --short=10 HEAD`，2026-11-05 15:1x）。射程＝**票 339 唯一未勾那一格 `AC#2`**（"去牙反控＝回落支整支摘掉 ⇒ 修好后的仪器必须在两平台都点名红"），另附逐 AC 复核表；票面其余六格维持 `339-v1`／编排者的既有判语，本腿只复核其是否能被今天的读数推翻。
★被裁对象＝落地腿 `339-r2` 的交件 `docs/evidence/s1/339-r2-impl.md`（提交 `dfd6f624b`）。

## 0. 起手现量（⛔ 抄腿，逐条现跑）

| 事实 | 读数 | 我这把尺 |
|---|---|---|
| 票面框 | 未勾 **1**／已勾 **6** | `grep -cE "^- \[ \] \*\*AC#"`／`"^- \[x\] \*\*AC#"`（两条各命中 1／6） |
| 被裁字节在 HEAD 上原样 | `git diff --stat dfd6f624b HEAD -- internal/risk/` = **0 行**（internal/risk 自落地笔起逐字未动） | 同命令 |
| 共享工作树里 internal/risk 干净 | `git status --porcelain -- internal/risk` = **0 行** | 同命令 |
| 落地笔的 `.go` 面恰三件 | numstat＝`334_test.go 29/31`、`334_windows_test.go 38/0`、`339_other_test.go 11/0` | `git show dfd6f624b --numstat -- "internal/risk/*.go"` |
| 落地笔未碰票 334 的 issues 件 | 该笔 name-only 里 `334-` 命中 **1**＝只那枚 `_windows_test.go`，无 issues 件 | `git show dfd6f624b --name-only` |

台件＝**仓外导出树**（`git archive HEAD | tar -x`，`/d/tmp/wisp/20261105/adj339/tree`，⛔ 复用腿的树、⛔ `-overlay`）；Windows 档＝本机 `go1.27.1 windows/amd64`（跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）；POSIX 档＝`golang:1.27` 容器、`CGO_ENABLED=0`、`GOPROXY=off`、宿主 module 缓存只读挂载、`GOCACHE=/tmp/gocache`、`TMPDIR` 留容器内 `/tmp`（⛔ 绑 Windows 卷＝`342-v1` 那枚定式）。四数只在 `^=== RUN`／`^--- PASS`／`^--- FAIL`／`^--- SKIP` 取，⛔ 相减。日志在 `/d/tmp/wisp/20261105/adj339/logs/**`（仓外，逐枚非 0 字节；清单见末节）。

## 1. `AC#2` 决定读数：同一枚突变，两平台各自点名红

突变＝`internal/risk/blacklist.go:298` 的 `aRelates` 尾行回落支整支摘掉（`return uncertainAnchorMiss(cand, anchor)` ⇒ `return false`，锚命中 **1**；与 `334-v1` 的 M-V1 同形）。突变只打在仓外树的物理字节上，逐发用 `git show HEAD:...` 还原并 `git hash-object` 自证。

| 台面 | 树 | RUN | PASS | FAIL | SKIP | rc | 红名册（逐名） |
|---|---|---|---|---|---|---|---|
| Windows | HEAD 未突变（正控） | 22 | 3 | 0 | 0 | **0** | — |
| Windows | HEAD＋去牙 | 22 | 0 | **3** | 0 | **1** | `TestAnchorFallbackTeethPerAnchor336`、`TestAnchorMissReasonNamesCandidateAndMissingAnchor`、`TestAnchorMissKeepsNonOverridableClassAndNeedL2` |
| POSIX | HEAD 未突变（正控） | 24 | 5 | 0 | 0 | **0** | — |
| POSIX | HEAD＋去牙 | 24 | 0 | **5** | 0 | **1** | 上三枚＋`TestAnchorMissReasonNamesCandidateAndMissingAnchorPosix339`、`TestAnchorMissKeepsNonOverridableClassPosix339` |
| POSIX | **未修形**＋去牙（对照） | 19 | 0 | 3 | 0 | **1** | 336 牙＋两枚 `...Posix339`；**票 334 仪器两枚名字不在名册里** |
| POSIX | HEAD 还原后再跑（正控 2） | 24 | 5 | 0 | 0 | **0** | — |

- 突变哈希：`P0(git hash-object)=44f6eb9b508eeaa0943cdeb9123bd19ad993bc78` ⇄ `MUT=305b6418e31d4801618c796ffbbfbd7ab850339c` ⇄ `RESTORED=44f6eb9b508eeaa0943cdeb9123bd19ad993bc78 match=YES`（Windows 发）；POSIX 发以 sha256 证同一对字节：`c37b0af9...` ⇄ `550b9863...` ⇄ 还原 `c37b0af9...`。`P0` 与腿报的 `restored=44f6eb9b...` **逐字同号**＝我量的是它交的那棵树的同一枚字节。
- POSIX 去牙后的红句逐字（`logs/linux-deteeth-after.log:86/88/95/96/97`，与腿报 `:81/:106/:143/:147/:151` 五句同形，行号差来自 339-r2 自己的注释重写）：
  `fixture guard: the fail-closed shape must still be A here, got class=B allow=false needL2=true reason="B-list sensitive path; single-file override requires one L2 confirm: id_*"`／`class: got B, want A (the fail-closed tier is unchanged by AC#1)`／`Gate: got class=B ... want class=A needL2=false allow=false (non-overridable)`／`an A deny from the fail-closed clause was unlocked by an override key: class=B allow=true ...`
- 作差归位：POSIX 未修形 19 枚名字 **全部**在修后 24 枚里（`comm -23` 空），多出的 5 枚＝2 顶层＋3 子用例＝**票 334 仪器本体**；未修形去牙的红名册里没有它＝"摘牙对它零变化"。⇒ 修后那两枚红**只能来自 339-r2 补的那半个世界**，⛔ 来自"另养一枚更弱的 POSIX 孪生"。

**⇒ 票面原句"修好后的仪器必须在两平台都点名红"今天成立**：同一枚仪器、同一批用例名（`TestAnchorMissReasonNamesCandidateAndMissingAnchor`／`TestAnchorMissKeepsNonOverridableClassAndNeedL2`）在 Windows 与 POSIX 上都被同一枚突变打成红，且两平台正控各自全绿。

## 2. "没有放宽断言换绿"的独立尺（⛔ 采腿自报）

| 尺 | 读数 |
|---|---|
| 被删行里有没有承重断言 | `git show dfd6f624b -- internal/risk/ \| grep '^-' \| grep -E 'want A\|ClassA\|t\.Skip'` = **0 命中** |
| 新增 `t.Skip`／`testing.Short` | `+.*t\.Skip` = **0**；三件里 `t.Skip\|testing.Short` grep = **无此两串**（`anchor_world_334_test.go` 同尺也 0） |
| 用例数有没有被删 | Windows 名册本腿量得 `22`（腿报改前改后同为 22，我未单独重建 Windows 改前树）；POSIX `19 ⊆ 24`、无名字消失 |
| 三句承重断言逐字在场 | `Classify(cand) != ClassA`（`:141`，红句在 `:143`）、`d.Class != ClassA \|\| d.NeedL2 \|\| d.Allow`（`:146`，红句在 `:147`）、override 键 `u.Allow \|\| u.Class != ClassA`（`:150`，红句在 `:151`）——现读 `sed -n '138,152p'` |

## 3. 逐 AC 表（`AC#2` 本腿裁；其余＝维持既有判语并复核"今天有没有被推翻"）

| AC | 判 | 依据（本腿现跑的那一把） |
|---|---|---|
| `AC#0` 起手现跑 | **通过**（维持） | 两平台基线本腿自打：Windows `22/3/0/0 rc=0`、POSIX `24/5/0/0 rc=0`；整包 Windows 共享树 `281/170/1/1 rc=1` 唯一红＝C26 存量（见第 4 节） |
| `AC#1` 三形代价表 | **通过**（维持） | 甲形是当年唯一合法形（乙按字面＝动产码）；今天 347 把解析器落地后，甲形的"少跑 5 枚名字"被乙形收拢回去＝`19 ⊆ 24`、多回 5 枚，与本腿读数同源 |
| `AC#2` 去牙反控两平台 | **★通过（本腿＝决定性格）** | 第 1 节六发读数；同一枚突变、两平台各点名红、正控全绿、还原哈希对得上 |
| `AC#3` 不许偷偷改写票 334 的账 | **通过**（维持） | 落地笔 name-only 里票 334 的 issues 件 **0 命中**；被动的只有它自己写面内的 `.go`（纯注释＋取候选那几行＋一枚 tag） |
| `AC#4` 颜色与门禁 | **通过（按 `dfd6f624b` 锚）**，另附两枚 HEAD 级红（见第 4 节） | 本腿复跑：`gofumpt.exe -l` 三件 = **0 行 rc=0**、`go vet ./internal/risk/` **rc=0**、整包 Windows `281/170/1/1`（只 C26）、POSIX 名册两发同号 |
| `AC#5` 禁区 | **通过**（维持） | `blacklist.go`／`pathresolver*.go` 非注释字节零动；`git diff dfd6f624b HEAD -- internal/risk/` 空＝HEAD 上也没有后来的腿动过它 |
| `AC#6` 纪律格 | **通过**（维持） | 票面框现量 6 勾／1 未勾＝腿一枚没翻；票名／`Status:`／`-done` 未动（本腿只读，未改一字） |

## 4. 两枚"不是 339"的红，具名归位（⛔ 记到本票头上）

1. **`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（整包唯一红）＝存量、⛔ 本票**〔量的〕。红句逐字点名 `cmd/wisp/panel_pump.go:283/:291` 与 `internal/panel/git_readonly_exec.go:409/:412`（"reads Result.Canonical without reading the rewrite account"）。★我另打一枚**独立性对照**：把 339-r2 只动的那三件按父提交 `5fd76dc7f` 逐枚还原（并删掉当时不存在的 `_windows_test.go`）、其余字节不动，同一发再跑 C26 ⇒ 仍 **rc=1**、四条红句逐字同。⇒ 这枚红与 339-r2 的字节无关，属台账 `A380` 那族（panel 禁写面）。
2. **`sh scripts/d22scan.sh` 今天在 HEAD 上 rc=1**〔量的〕。红因不涉 `internal/risk`：`TestScannerSelfScanOfRealRepoIsGreen` 命中 `internal/panel/composer_handlers.go:33` 的 phantom-citation（引 `docs/evidence/s1/R9-impl.md`，该件不在盘上），`TestRealRepoLedgerIsHonest` 报 ban #8 的 `frontend/` 射程扫到 0 枚文本件。两个被点名面最后被动在 `8e0702fdb`／`b4ee24cbe`（都在 `dfd6f624b` 之后、都不是 339 的腿），而 339-r2 的写面在 `internal/risk/**` 三枚测试件。⇒ **HEAD 级后笔红，⛔ 339**；票 339 `AC#4` 的门禁读数仍以 `dfd6f624b` 锚上那组（腿报 d22scan rc=0 clean）为准，本腿不在 HEAD 上重判它。
   另：`sh scripts/check-path-length-budget.sh` 今天 **rc=1**＝既有 roster 漂移（`179-...-sentence.md` 已改名 `-done`，roster 没跟上），与本票零关系。

## 5. 一句话终态

票 339 的 `AC#2` 今天**成立**：同一枚去牙突变（`aRelates` 尾行回落支整支摘掉）在 **Windows `22/0/3/0 rc=1`** 与 **POSIX `24/0/5/0 rc=1`** 上都把票 334 那两枚仪器点名打红，未突变正控两平台各自全绿（`22/3/0/0`／`24/5/0/0`），未修形对照组在 POSIX 上那两枚名字根本不在名册（`19/0/3/0`）＝补的是世界、⛔ 藏红。⇒ **七格齐（六维持＋`AC#2` 通过）＝票 339 可以结案**；翻框与 `-done` 归编排者，本腿一枚没翻。

## 6. 具名（⛔ 阻结案，给下一位）

- 本腿的台件与日志在仓外 `/d/tmp/wisp/20261105/adj339/**`（`tree/` 导出树、`logs/` 10 枚日志＋2 枚 pristine 拷贝、`sh/` 三枚脚本）。⚠ 按派单"只写两枚判语件"的约束，**本腿未把 probes 落进 `.scratch/wisp/probes/339/`**；要长期留证需由编排者决定是否搬进库。
- POSIX 侧本腿只跑名册 `-run 'TestAnchorMiss|TestAnchorFallbackTeeth'` 与整包 Windows 一发，⛔ 跑 `./cmd/wisp` 整包、⛔ `go test ./...`。
- 本腿 ⛔ 动票 339 一字、⛔ 动任何 `.go`、⛔ 任何 git 写操作（`git show`／`git archive`／`git hash-object` 只读）。

## 日志清单（仓外，逐枚非 0 字节）

`win-roster-after.log 6166`／`win-deteeth.log 14664`／`win-fullpkg.log 174246`／`win-roster-recheck.log 53`／`linux-after.log 7784`／`linux-deteeth-after.log 15946`／`linux-deteeth-before.log 13763`／`linux-again.log 7788`／`c26-head.log 1175`／`c26-parent-shape.log 1175`（字节数，`stat -c%s`）。
