# 361-a1 普查腿（只读）part 2／2：正控＋两侧对照＋自检＋门禁＋停手上报

- 续件：`docs/evidence/s1/361-a1-census.md`（sec 1-3：前提新鲜度、35 枚逐名归类、归属规则逐字）
- 时刻：2026-10-22 16:4x +08；HEAD：`21f05dd4`；本机 `core.autocrlf=true`；gofumpt `v0.12.0 (go1.27.1)`

## sec 4 正控（硬闸）：真脏 tracked `.go` 必须红并点名

**未在仓内造**（那是落 `.go` 字节）。在仓外 `git archive` 导出树 `D:\tmp\wisp\20261022\361a1\archive` 内造：

- 新件：`internal/census361a1/positive_control.go`（非台件、纯 LF，`tr -dc '\r'` = **0**）
- 已 `git add` 并 commit（该树自己的 `.git`，`git ls-files` 可见）
- 内容故意非 gofumpt-clean（`x := 1` 后多余缩进的 `return x`）

读数（`archive-with-poscontrol.txt`，rc=**1**）：

```
   TRACKED-DIRTY	internal/census361a1/positive_control.go	(last commit touching it: bf6f6e92 10-22 16:39 census)	<== CI IS RED ON THIS
attrib.sh: (A) tracked-tree  lines=21 files=21  rule: MUST be empty
attrib.sh: RED (tracked-only) - tracked-dirty=21 files_dirty=21. ...
```

＝> **尺是活的**：20 -> 21，且**逐名点名**了正控那一枚。`NOT` 恒红、`NOT` 恒绿。

## sec 5 两侧对照：仓外 `git archive` 导出树（与 CI 同形）

导出：`git archive --format=tar HEAD | tar -x -C <仓外目录>`，再 `git init && git add -A && commit`（该树自足索引）。
该树性质：**无未跟踪台件状态**（只有 tracked 内容）、`*.go` 全 **LF**（抽查 4 枚非台件 `CR=0`，20 枚 B 类 `archive_CR=0`）。

同一发读数（`archive-tracked-only.txt`，rc=**1**）：

```
== (A) tracked tree: ... (tracked .go files handed to it: 1155)
   TRACKED-DIRTY	.scratch/wisp/probes/163/a1/main.go ...
   ... 共 20 行 ...
attrib.sh: (A) tracked-tree  lines=20 files=20  rule: MUST be empty
attrib.sh: RED (tracked-only) - tracked-dirty=20 files_dirty=20. ...
```

**35 -> 20。** 掉的 15 枚正是 sec 2 的 C 类（CRLF 假红），留下的 20 枚全是 B 类。
这 20 枚在该导出树里与 HEAD blob **逐枚 sha256 相同**（20/20 IDENTICAL），即**入库字节本身就脏**。
再用版本稳定的 `gofmt -l` 复核：20 枚中 **15 枚连 gofmt 都判脏**（另 5 枚是 gofumpt-only 规则，如 `0644` -> `0o644`、复合字面量缩进）。
把 `ci.yml:172` 那步的形状（`gofumpt -l . tools/d22scan tools/mockllm`）在同一导出树上跑，得到**同样这 20 行**。
另核：这 20 枚的 blob 在最后一次已推送提交 `626fcc7b` 上**逐枚相同**（20/20 identical_at_pushed）。

＝> **仓外导出树上它仍然红，且不是台件＋CRLF 能解释的。触发停手条件 (2)。**
＝> 按 `ci.yml:191-193` 自己的话，checkout 上这两步读同一批字节、必须一致；本读数是 20 枚同形红。

## sec 6 `--self-test` 的 rc 与四数

`sh .scratch/wisp/probes/161/r5/attrib.sh --self-test` ＝> rc=**0**（本机，未被 CRLF 影响）

四数：`cases=10`、`failures=0`、`ring=5`、`silent=3`（`named=3` 全在）
逐字尾行：`attrib.sh: --self-test cases=10 failures=0` / `SELF-TEST GREEN - ...`
其中 275 名册控制那一格：`PASS  want rc=2+roster-names-the-sample  got rc=2 rows=1`，临时树留在 `/tmp/wisp-275-r2-rosterctl.C7GT1n`（脚本自身 create-only 行为）。

## sec 7 门禁

| 门 | rc |
|---|---|
| `sh scripts/d22scan.sh` | **0**（`clean - no D22 ban violations`） |
| `sh scripts/check-path-length-budget.sh` | **0**（`VERDICT GREEN`；denominator 17564、over-budget 61、roster 61、not in roster 0） |

## sec 8 交付件与纪律

- 本腿**新增 `.go` 文件数 = 0**（正控造在仓外导出树，不在本仓）
- 本腿写面：`docs/evidence/s1/361-a1-census.md` ＋ 本件 `docs/evidence/s1/361-a1-census-part2.md`（各 160 行内）＋ 仓外 `D:\tmp\wisp\20261022\361a1\**`
- 未改 `.github/**`、未改 `.scratch/wisp/**`、未改 `Status:`、未翻任何框、未 push
- 正控尺（可复跑）：在仓外导出树放一枚非台件、纯 LF、非 gofumpt-clean 的 tracked `.go` ＝> 该发必 rc=1 且逐名点名它

## sec 9 停手上报（未定义即停，两条都触发）

**停手条件 (1)：`attrib.sh` 的归属规则不认 `.scratch` 台件在 `--tracked-only` 下的红。**
票面前提「31 枚台件、归属规则允许它们红」是错的。逐字读数见 sec 3：归属规则（`ticket_of` :185-193 ＋ `ticket_known` :196-203）只在 `classify_line()` 里合流，而 `classify_line()` 只被 `--classify-line`、`--self-test`、`ruler_b` 调用；`--tracked-only` 只跑 `ruler_a`，那条路径上**没有任何归属判定**，tracked 且不 clean 即红。
实测反证：`sh attrib.sh --classify-line .scratch/wisp/probes/163/a1/main.go` ＝> `TRACKED<CI-RED> ... rc=1`；同一发 `--tracked-only` 里它同样红。

**停手条件 (2)：仓外导出树上它也红（20 枚），所以本机红不只是台件＋CRLF。**
20 枚 B 类在纯 LF、无未跟踪台件状态的导出树上仍红，且与 HEAD blob 逐枚 sha256 相同、与已推送 `626fcc7b` 的 blob 逐枚相同；其中 15 枚连版本稳定的 `gofmt` 都判脏。
＝> 「本机 35 枚红不代表 CI 会红」这句只在 15 枚 CRLF 假红上成立；对 20 枚 B 类，本机与 CI **同形**，且 CI 的 `ci.yml:172` 那步会在同一批字节上先红。

**未取到的读数（按 AC#3 口径登记为待推送取数，不用本机冒充）：**
本机存的 CI 日志（`.scratch/ci-logs/run-37166458550-failed.log`、`run-37158259050.log`）里 `lint` 作业只留了 `D22 scanner positive control` 与 `staticcheck` 两步的日志，**没有** `attrib.sh` 那一步的行。
故「CI 上它是绿」这一条**目前没有任何 CI 侧原件支撑**；本腿不认它、也不反认它，只报「导出树同形读数是红」。

**需编排者裁的两件事（本腿不自行假设）：**

1. 票面「本机红 / CI 绿」这个前提要不要按本读数重写（改成「20 枚真红、15 枚 CRLF 假红」）——这是改票面文字，属编排者。
2. AC#1 的三形（甲/乙/丙）是否还成立。若 20 枚真红属实，则「让 `attrib.sh` 本机按入库 blob 读」只消掉 15 枚假红；剩下的 20 枚要靠**修字节**（或另立票），不在 `attrib.sh` 射程内。
