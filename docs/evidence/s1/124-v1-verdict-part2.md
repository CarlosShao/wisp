# 票 124 终裁腿 124-v1 裁决 part2（sec 二、三、四）

> 前接 `docs/evidence/s1/124-v1-verdict.md`（sec 0 与 sec 一）。

## 二、逐格结论（我只判，不翻）

| 格 | 能不能翻 | 缺什么 |
|---|---|---|
| **AC#2b** | **NOT 能**（票级） | 本包 33 => 0 已成立，但票级终判据"软链形红名＝0"余额仍在：`internal/tools` **38**（扣票 123 三枚豁免后 35）＋`cmd/wisp` **37**。批 2/3 未做。 |
| **AC#2** | **NOT 能** | 票面 `:170` 已定"2b-1..2b-4 全绿再按复算翻"。今天 2b-2/3/4 全未做。 |
| **AC#3** | **能翻（就批 1 这一发而言）** | 四条尺全过（sec 一.3）。若编排者要票级口径，等批 3 合并态再翻亦可；**批 1 层面我不认为还缺东西**。 |
| **AC#4** | **能翻（就批 1 这一发而言）** | 变异红/还原/原样三态我独立复现，42 与 6 名额外红逐数对上。唯一缺的是**它没交 sha256 还原证词**（我已用干净工作树替代证明）。 |
| **AC#1 / AC#2a** | 早已翻（非本轮） | — |
| **AC#5** | **NOT 能** | 票面写"受影响包"，2b-2/3/4 的受影响包还没进门禁账（与 2b-1/2/3 同处理）。 |
| **AC#6** | **NOT 能** | 与本轮零交集（POSIX 侧 `sameAbsoluteness` 那格未动）。 |

## 三、那 2 枚残余红：独立复量与"要不要登台账"

**读数是摇摆，NOT "改前就红"**。我在四台各 `-count=20` 只跑这两枚：

| 台 | `TestNothingWrittenPromise...` P/F | `TestHandDeletedSchemaVersion...` P/F |
|---|---|---|
| PRE 纯净件 · 普通形 | 7 / 13 | 11 / 9 |
| POST 纯净件 · 普通形 | 7 / 13 | 14 / 6 |
| PRE 纯净件 · 容器原生 overlayfs 副本 | 7 / 13 | 12 / 8 |
| 宿主 Windows/NTFS 原生 | **20 / 0** | **20 / 0** |

- **两枚都摇**（NOT 只有第二枚）。实现腿说"第二枚 `-count=5` 得 2 PASS/3 FAIL"——那是**一发样本**；
  我那发 POST `-count=5` 得 3/2，方向相反。`-count=5` 太小，任何一侧结论都过不了。
- **失败断言只有一条**：`mtime did not move ...`（两处行号 `writeguard_unnailed_341_test.go:239`、
  `writeguard_handedit_227_test.go:136`）。**字节断言（`was NOT re-laid` / `left byte-identical`）0 命中**
  => 字节确实变了，只有 mtime 没动。
- **机制**：`newTestManager` 的 mutate 闭包是 `time.Sleep(10ms)` 后 `os.WriteFile`，两枚用例比对 `statMod` 前后是否相等。
  我另打一发 mtime 粒度探针（10ms 间隔写 20 次 × 三台文件系统）：**distinct-mtime 20/20，三台全部分得开**
  => 不是粒度问题，是这条 10ms 竞态在 Linux 上会输给某些写路径的 mtime 戳。
- **"既有"这个定性成立，但成立的方式与它说的不同**：两枚用例的引入 commit 是 `3b8c114e`（票 227）与 `6948eb08`（票 341），
  都在票 124 之前；`git diff 4ac4549c HEAD --` 那两枚定义文件 **空**；PRE 锚上 20 发里照摇（13/9 次红）
  => **与本腿零因果**。但它说的"改前普通形就红（PRE-PLAIN FAIL=2 正是这两名）"我没复现到——我那一发 PRE-PLAIN **0 红**。
- **整包普通形重复 8 发**（这是出口判据的真正分母）：PRE 得 TOPFAIL 取 0/1/2，`0×2 / 2×3 / 1×3`；
  POST 得 `0×1 / 2×3 / 1×4`。**两侧都横跨 0..2，且从未冒出第三名**
  => 普通形红数在本包**本来就不确定**，单发"枚数相同"不能当判据；能当判据的是**名册集合**（两侧恒为那 2 名）。
- **要不要登台账：要，而且两枚一起登**（我的判）。理由三条：
  (1) 它是一条**只发生在 POSIX** 的摇摆（Windows 20/20 从不红），任何"普通形 0 红"的门禁在 POSIX 上都会非确定性地假红；
  (2) 现有台账对第一枚只有 **1 处**命中（`:19381`，是票 227 收件时记它 `--- PASS` 的一句，**NOT 摇摆登记**），
      第二枚 **0 命中**（`:20351` 那一处是编排者本轮的话，NOT 登记）=> 两枚实质上都未登记；
  (3) 不登的后果正是本轮踩到的：同一形状两发能得出相反的结论，下一位裁腿会把它读成"本腿引入"或"本腿修好"。
  注意：台账 `docs/reports/pending-and-issues.md` 不在我的写面 => **我只判"要登"，登由编排者落**。

## 四、门禁 rc 与末行逐字（我自己现跑，`unset GOTMPDIR TMPDIR TEMP TMP` + `GOFLAGS=`）

- `GOFLAGS= go build ./...` => **rc=0**，输出 **0 行**（无末行）。
- `go vet ./internal/config/` => **rc=0**，输出 **0 行**。
- `go test -count=1 ./internal/config/` => **rc=0**，末行逐字：`ok` + 包路径 + `2.914s`。
- 同发 `-v` 四数分开量：顶层 `PASS=114 FAIL=0 SKIP=0`、`RUN` 行 **195**、子测试 `PASS=81 FAIL=0 SKIP=0`（末行同格式 `2.175s`）。
- `sh scripts/d22scan.sh` => **rc=0**，末行逐字节录（首尾）：
  `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=244, bans #1-5 cmd/=42, ban #6 frontend/=116, ban #7 internal/tools/=25, ban #8 design/=39, ban #8 frontend/=116, ban #8 internal/=608, ban #8 cmd/=143`
  正向对照：`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`
- `sh scripts/check-path-length-budget.sh` => **rc=0**，末两行逐字节录：
  `check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`
  `check-path-length-budget.sh: unit cross-check: the byte view and the character view select the SAME 60 paths on this tree, recomputed just now in both units`
  （`tracked paths=17493  over-budget=60  covered by roster=60  not in roster=0`）
- 未跑项与理由：`CGO_ENABLED=0 go build ./...` 按派单明示在本仓必然失败（sherpa 全文件被 build constraint 排除），
  **我未跑**，故不报读数、也不当缺陷；`gofumpt` 我未跑（实现腿报 v0.12.0，属它那一发的事实，我不替它背书）。
