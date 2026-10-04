# 212-r3 修码证据件（票 212 第 3 轮小写码腿；三格欠账＋一枚前缀同步钉）

- 本程：写码腿 **212-r3**。派单来源＝非实现者验收腿 **212-v2** 件
  `docs/evidence/s1/212-comments-phantom-citation-v2.md`（33578 字节，本程现读）§8 表里点给 212-r3 的三笔具名欠账
  ＋编排者派单追加的一枚"两处前缀集合相等"钉。
- 起手锚：`dev` HEAD＝`b1e59d63`（260-r1 的措辞枚）。本程写面只有四枚文件：
  `tools/d22scan/main.go`·`tools/d22scan/selftestsamples.go`·`tools/d22scan/scan_test.go`（新钉）·
  `internal/agent/approval/pending_read.go`，加证据件目录 `.scratch/wisp/probes/212/r3/**`。
- 权威文本：票面 `.scratch/wisp/issues/212-comments-cite-evidence-files-that-do-not-exist.md`（只读，本程零写；
  含 2026-10-03 的「编排者终裁节」与 2026-10-04 的「编排者翻勾节」）。
- 本程一句话交付：**行为零变化**——(a) 与 (c) 是纯措辞、(b) 是加样本与措辞、钉是加测试；
  摘豁免／摘前缀这两道突变的成对读数在 §2／§5，终态门禁对比在 §7。

## §0 起手锚与本程自跑的基线尺（先落档再动手）

起手 md5（本程现跑，档 `logs/md5-baseline.txt`）：

| 文件 | 起手 md5 |
|---|---|
| `tools/d22scan/main.go` | `e6d1745996e8e02c82ab691ac2cbadf8` |
| `tools/d22scan/selftestsamples.go` | `c61abbeba9b55a6ea452cec8478c2c93` |
| `internal/agent/approval/pending_read.go` | `e0c9dc3d5798514f7acba07984605680` |

基线读数（本程现跑，档 `logs/selftest-pristine.txt`／`logs/gate-pristine.txt`）：

| 尺 | 起手逐字末行／关键行 | rc |
|---|---|---|
| `cd tools/d22scan && go run . -self-test` | `d22scan -self-test: clean - all 37 direction checks passed (20 expect-ring, 17 expect-silent)` | 0 |
| 同上·名册行 | `d22scan -self-test: roster read from main.go = 9 numbered ban(s) [1 bare-goroutine, 2 pathresolver-bypass, 3 plaintext-key, 4 wallclock-timeout, 5 mirror-hash, 6 panel-approval, 7 internal-artifact-tool, 8 emoji, 9 phantom-citation] + 1 finding type(s); 37 cases, 10 tag(s) covered, both directions required per tag` | 0 |
| `sh scripts/d22scan.sh` | `runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0` | 0 |

八枚 scope 行的起手逐字（档 `logs/gate-pristine.txt:241-248`）与 `clean` 行逐字在 §7 现跑对比；
本程以**自己这一发**为 AC#5 的基线，不搬 `212-v2` 的数（并发腿会把 `ban #8 internal/` 的分母推走——
起手这发已经是 500，v2 那发是 498，这一枚差 2 的文件数与本程无关，具名登记在 §7）。

## §1 (a) `pending_read.go` 的出处指错——改成实话

本节记：三处现读逐字（`internal/tools/gate.go`／`cmd/wisp/panel_pump.go`／`internal/agent/approval/pending_read.go`）、
改前改后逐字对照、以及"哪些字段属哪一句"的归属裁决依据。

## §2 (b) 三形 silent 样本＋成对突变

本节记：三枚新 expect-silent 样本的逐字 src、每形对应 `212-v2` §3 的哪一枚探针（P3/P4/P5）、
摘豁免突变的红句逐字、还原后的 md5 与 `git diff -- tools/d22scan`。

## §3 (b) 名册行与 `main.go` 自述改成真实射程

本节记：`// Bans` 第 9 枚名册行、`scanGoFile` 里 ban #9 的自述段、`shorthandPathStarts` 的自述段三处改前改后逐字，
以及"为什么这是把自述追平行为、不是把行为放宽"的判据。

## §4 (c) 仪器自己注释里的那枚幻影

本节记：`selftestsamples.go` 引的 `docs/evidence/s1/212-citation-ruler.md` 盘上不存在（`ls` 现跑）、
改指到哪一枚真实件、ring 向为什么改完仍响（夹具不 seed 那一枚）、以及 `tools/**` 是否在射程内（不在，本程不扩）。

## §5 前缀同步钉（`repoPathRe` ↔ `shorthandRegionRe`）

本节记：钉的落点 `file:line`、正控＝在 `repoPathRe` 单独加一枚前缀 ⇒ 指名那枚测试必须红（逐字）、
还原⇒绿、以及为什么"集合相等"而不是"字符串相等"是这条钉的正确强度。

## §6 撞钉预检：哪些断言里硬写了数字

本节记：起手那发 grep 的逐枚命中（`tools/d22scan/**.go`＋`scripts/**`＋`.github/workflows/ci.yml`＋台账），
逐枚判"硬写的数字"是断言还是散文，以及本程同步了哪一枚、为什么那是同步不是放宽。

## §7 门禁读数终态（AC#5：八枚 scope 行逐名对比）

本节记：终值 `-self-test` 行逐字（含新分母）、`go test` 的 PASS/FAIL/SKIP、`sh scripts/d22scan.sh` 八枚 scope 行
与起手的 `diff` 结果与 rc。

## §8 判不动的地方（具名，不用"应该没问题"填空）

本节记：本程判不动／归编排者／归后续程的每一枚，逐枚写清"为什么本程动不了"。
