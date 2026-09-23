# 票 124 AC#2b 批次2（`worker-ticket124-ac2b-2`）— tools 20 + llm 17 + perm 5 + approval 1 接上票 119 那条纪律

**日期**：2026-09-23 · **性质**：转换方交件（改 `.go`，但只改 `_test.go`）
**本批范围**：清点账 `docs/evidence/s1/124-ac2a-leg-classification.md` §5.8（tools 行 107-124、126-127，除行 125 归票 123）、§5.5（llm 行 52-68）、§5.6（perm 行 69-73）、§5.3（approval 行 44），共 **43 枚**，账上全标「可转」。
**票面判据**：`.scratch/wisp/issues/124-*.md`「AC#2b 的放行与分批」（09-23 16:33 编排者）那张表 + 五条共用结案判据；批次 1 的 `next=` 两条本批提醒照办。

## 0. 锚点 / 快照 / 跑法（可复核）

- 开工首读 `git rev-parse --short HEAD` = **`5417a3c`**（与派单给的锚点一致；下文记 **pre 锚**，全部改前读数在它上面量）。
  ⚠ 共树漂移：测量途中 HEAD 先后经 `4d266f9`、`38c09df` 到本批改件 `dfa3dc4`。已核 `git diff --stat 5417a3c dfa3dc4^ -- internal/ go.mod go.sum` **零 hunk**；`git diff --name-only 5417a3c dfa3dc4^` = 9 枚文件，全部在 `.scratch/wisp/issues/`、`docs/`、`cmd/wisp/**`（兄弟代理的 docs 与 2b-4 地界的 `cmd/wisp` 码）⇒ 与本批四包无交集，pre 读数量的是被改前同一版本。
- 改件 commit：**`dfa3dc4`**（`test(124,AC#2b-2)`，18 枚文件，+161/-30；其中删除侧 30 行全是递根点，见 §3）。
- 全部测量都在 `git archive <sha> | tar -x` 的纯净快照里做；工作树里别人的半成品（`.scratch/.../134-*.md`、`docs/reports/*.md`、`cmd/wisp/leg_dispatch_gate_133_test.go`、`cmd/wisp/part133b.go.txt`）一个字都不读。
- 容器 `golang:1.27`（`go1.27.1 linux/amd64`，容器内 `CGO_ENABLED=1`），复用票 119 的命名卷 `ac119-gomodcache` / `ac119-gocache` ⇒ 离线可编。
- 挂载一律 `/d/...` + `MSYS_NO_PATHCONV=1`；每枚样本进容器第一件事打 `ls -l /src/go.mod` + `md5sum /src/go.mod` = `f6ef661732b1851e5c3db348113cb605`、`md5sum resolve.go` = `b6876a5efe759f6e17434d1b50a129c3`（与 AC#1/AC#2a/AC#2b-1 逐字同字）⇒ 非空挂自证。
- 形状硬断言沿用批次 1：软链形 `exit 97`（`/varlink` 不是 symlink）/ `exit 98`（`readlink -f /varlink/w124tmp != /realpriv/w124tmp`）；普通形 `exit 99`（`/varlink` 必须根本不存在、`/plainroot` 不许是链接）。
- 跑法脚本（新建，未改 2b-1 那枚）：`/d/tmp/wisp124-2b2-h.sh`（`bash /h.sh <link|plain> <tag> [pkgs]`，`COUNT=2` 切门禁那一发；`-timeout 25m`——tools 软链形含票 123 那两枚 ~300.0x s 的已知腿，12m 会掐）。变异台：`/d/tmp/wisp124-2b2-m.sh`（判据⑤：先 grep 证落地 + `go build` rc=0 再读红名）。转换件：`/d/tmp/wisp124-2b2-convert.py`（30 处行号锚定替换，锚不中即 `ANCHOR MISS` 停手）。日志全在 `/d/tmp/wisp124-2b2-logs/`。
- 开测前查在飞 run（本机就是 self-hosted runner）：18:13 本地（`10:13z`）`gh run list --limit 5` 最近 5 枚全部 `completed` ⇒ 无在飞争用。
- 快照表：

| 快照目录 | 内容 | 文件数 |
|---|---|---|
| `/d/tmp/wisp124-2b2-pre` | `git archive 5417a3c`（改前基线） | 1039 |
| `/d/tmp/wisp124-2b2-post` | `git archive dfa3dc4`（改后） | 1046 |
| `/d/tmp/wisp124-2b2-probe` | 待量：改后 + 判据④取字串探针 | - |
| `/d/tmp/wisp124-2b2-mut` | 待量：改后 + 判据⑤变异 | - |

## 1. 对派单数字的复核（派单给的每个数字都是断言）

- 清点账点名（读 `124-ac2a-leg-classification.md` 逐行）：tools §5.8 = 行 107-127 共 21 枚，其中行 125（`TestL1WriteGoesThroughTheRealBlockWindow`）标「归因待票 123」不计本批 ⇒ **20**；llm §5.5 = 行 52-68 ⇒ **17**；perm §5.6 = 行 69-73 ⇒ **5**；approval §5.3 = 行 44 ⇒ **1**。合计 **43**。⇒ 与派单「20+17+5+1=43」**一致，登记差 0**。
- 实测复算（pre 锚软链形逐名 `--- FAIL`，两形各一枚）：待量（PRE-L / PRE-P 在跑，读数落 §2 表）。
- ⚠ 票 123 那两枚已知红（`TestL1Write` / `TestLateVeto`，同包 `wiring_test.go`，300.0x s 那族）**不在本批 43 枚**：`internal/tools/wiring_test.go` 本批一字未动（§3 的 `git show dfa3dc4` 文件清单可核），它们在改后软链形照常红，逐名登记为已知遗留。

## 2. 判据①：逐枚转绿且点名（`-v` 才有 PASS 名）

待量（PRE-L / PRE-P / POST-L / POST-P 四台对照表 + 逐包四数）。

## 3. 判据②：普通形一枚都不许多红 + `git diff` 证判定分支一字未动

待量。

## 4. 判据③：两形各取一枚 + RUN/SKIP 差逐包解释

待量。

## 5. 判据④：本批「另一种拒」复算 — 逐枚实拿字符串

待量。

## 6. 判据⑤：变异自证（票面 AC#4）— 三态原文

待量。

## 7. 门禁读数（票面 AC#5）

待量。

## 8. 改动面与本批用的 helper

- 用的既有 helper = **`proc.SealableRoot`**（`internal/proc/envfork.go`），即票 119 生产路 `TestDataDir` / `DefaultLayout` / `cmd/wisp` 的 `resolveDataDir` 走的那一枚；批次 1 同法。
- 新增五枚**单行委托**（`tools`、`tools_test`、`llm_test`、`perm`、`approval_test` 各一枚；tools 因内/外两枚测试包作用域不同而拆两枚，**不是**第二份解析实现）+ 一枚 `sealableTempCanonical124`（= `mustCanonical(sealableTempDir124(t))`，被测试的规范化器照旧跑，变的只是喂进去的根拼写）。零新解析路数：`R-125-2` 那本副本账不增行。
- 递根点 30 处（tools 22 行 / llm 1 / perm 6 / approval 1），对应 43 枚（llm 的 1 处 fixture 覆盖 17 枚；tools 的 fs_write 3 处覆盖 4 枚含两子测试）。
- 未动：`internal/tools/wiring_test.go`（票 123 两枚）、`internal/tools/paths_ticket107b_probes_test.go` 的 ProbeA/ProbeB 与 `ticket107_portable` 非红枚、`fs_write_test.go`/`paths_workspace_test.go`/`pathshape_portable_test.go` 里服务非红用例的 `tempRaw`/`t.TempDir()` 站点、`internal/winsec/**`、`internal/risk/**`、`internal/memory/**`、`internal/config/**`、`cmd/wisp/**`。

## 9. 未验证项 + `next=`

待量后写。

## 10. 临时件清单（只建不删）

`/d/tmp/wisp124-2b2-{pre,post}/`、`/d/tmp/wisp124-2b2-logs/`、`/d/tmp/wisp124-2b2-h.sh`、`/d/tmp/wisp124-2b2-convert.py`（探针/变异台目录待量后登记）。
