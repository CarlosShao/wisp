# 证据件：票 195 AC#5 刷新（票 411 AC#6，2026-11-09）

本件由票 411 AC#6 写腿现量出具（非 195 实现者）。所有读数带可复跑命令；引文逐字。

## 1. AC#5 原文（票 195 line 28 逐字，一字未动）

> - [ ] **AC#5 门禁四数＋`gofumpt`**：`sh scripts/d22scan.sh` clean；`bash .scratch/wisp/probes/154/gate-clauses.sh` 的 **BAD 腿名册只 `G6neg`**（在册常红）；`go test -count=1 ./internal/config/` 全绿；`"$(go env GOPATH)/bin/gofumpt" -l` 空。⚠ **别拿 `ban #8 internal/` 那格的文件枚数当不变量**（`A384` 已更正：那是分母，不是违规数）。

本腿只在该框下追加具名更新一行（`SPEC-12 §4.1`：AC 原文是契约，只能追加不能改写）。

## 2. 真读数对照（四支判据逐支）

| 判据支 | 原文（立票时口径） | 修后真读数（锚 `cb2287a5`，2026-11-09 09:4x） | 判 |
|---|---|---|---|
| `d22scan.sh` clean | 要求 clean | **rc=0 clean**（本腿 11:0x 复跑：451 产文件，bans #1-5 internal=387＋cmd=64、ban#7 tools=48、ban#8 internal=966＋cmd=244 全过） | 绿 |
| gate-clauses BAD 名册 | 「只 `G6neg`（在册常红）」 | **rc=9，BAD＝G1/G1b/G2/G5/G5pos/G5neg/G6neg/G7/G7neg 九腿**（票 411「修后真读数」表；证据 `docs/evidence/s1/gate-clauses-argv-overflow.md` §3） | 原句作废（见下） |
| `go test -count=1 ./internal/config/` | 全绿 | ok 2.419s（收件r1 导出树 `cda38fd52` 复跑，十枚 AC 命名用例逐名 PASS） | 绿 |
| `gofumpt -l` | 空 | 空（收件r1，四交件文件） | 绿 |

原句为何过期：gate-clauses.sh 因 `run()`/`pair()` 内 `-- $@` 未引号 argv 溢出（1084 词／39,836 字符 ＞ Windows 32767 上限）⇒ 14 腿里 10 腿假 BAD（`git-grep-rc-126`）、聚合恒 rc=10。「只 G6neg」是溢出期假读数。票 411 修法＝`-- "$@"`（判据一字未动），修后射程从 bash 两层展开变 git 递归 glob，9 枚盲区文件第一次被看见 ⇒ rc=9 是仪器史上第一次全树真读数，九腿全部具名、无一是假红。

## 3. 九腿各自归谁（具名命中 → 责任族）

| 腿 | 具名命中（票 411 修后真读数表） | 归属族 | 涉 195 射程？ |
|---|---|---|---|
| G1 | `internal/plugin/manifestslot/tier2.go:46` | manifestslot/plugin | 否 |
| G1b | `tier2.go:46`、`internal/plugin/tier2_gate.go:210`、`tier2_loader.go:195`、`internal/tools/subagent_197.go:616` | manifestslot/plugin/tools | 否 |
| G2 | `cmd/wisp/run.go:1514`（OpenTask）＋`:1520`（CloseTask） | resident（cmd/wisp） | 否 |
| G5 | 在册 `cmd/wisp/panel_assets.go`＋新 `cmd/wisp/resident_present_windows.go` | cmd/wisp resident | 否 |
| G5pos | `cmd/wisp/resident_present_windows.go` | cmd/wisp resident | 否 |
| G5neg | resident_present＋394/382/403 家族测试＋provenance_syncdirs_other_test（负控枚数漂移） | cmd/wisp | 否 |
| G6neg | `internal/tools/web_bridge_test.go`（只开不关） | tools 测试 | 否 |
| G7 | `resident_voice_bargein_windows.go`、`resident_voice_windows.go`、`internal/plugin/goja_runtime.go`、`manifestslot/command.go`、`manifestslot/lifecycle.go` | resident/plugin/manifestslot | 否 |
| G7neg | 同 G7 家族漂移 | resident/plugin/manifestslot | 否 |

195 射程＝本票交件 `e6791e6bb` 四文件（`internal/config/manager.go`、`sectionwrite.go`、`sectionwrite_195_test.go`、`settings.go`）。九腿具名命中与该射程**零交集**。

本腿现跑 scoped 复量（HEAD 树，射程只 `internal/config/**/*.go`，六把词表尺）：

```
git grep -nE "ToolRequest\{" HEAD -- "internal/config/**/*.go"                          # 0 命中
git grep -nEw "\.Execute\(" HEAD -- "internal/config/**/*.go"                           # 0 命中
git grep -nEw "OpenTask|CloseTask" HEAD -- "internal/config/**/*.go"                    # 0 命中
git grep -nE "OpenScope\(|CloseScope\(|\.Close\(\)" HEAD -- "internal/config/**/*.go"   # 0 命中
git grep -nEw "Defer(Named)?\?|(New)?DisposalScope" HEAD -- "internal/config/**/*.go"   # 0 命中
git grep -nEw "OpenTask|OpenScope" HEAD -- "internal/config/**/*.go"                    # 0 命中
```

六尺全 0 ⇒ **本票射程内 0 红**（2026-11-09 11:0x 本腿亲跑）。

## 4. AC#5 能不能翻（判语）

1. 四支判据里三支绿（d22scan／config 包测试／gofumpt）；第四支 gate-clauses 的判据本体是「BAD 腿名册」读数——原文写死的名册（只 G6neg）已被票 411 证明是溢出期假读数，且修后射程语义变化（基线重登）归裁决方按 171 r3 规矩处理，不归本票。
2. 九腿真红逐枚归别票族（manifestslot/plugin/tools/resident），零涉 `internal/config/**` 与本票交件四文件 ⇒ **凭「本票射程内 0 红」，AC#5 具备翻框条件**。
3. 但本腿纪律 ⛔ 翻框（本编队只有 collect 腿受权翻框），且 411 AC#4/基线重登仍待裁决方 ⇒ **勾选留裁决者**：裁决者若采「射程内 0 红」口径即可翻 AC#5（追加行本腿已落，翻时原句保留）。
4. 全票状态：AC#0 永缺（10-02 裁定，判「起手对得上」的现量表两格已不符，不再复跑）；AC#1/2/3/4/6 已翻〔成立〕（收件r1 2026-11-09 09:4x）；**AC#5 待裁决者** ⇒ 若翻，票 195 全票收口（AC#0 按 10-02 裁定作「永缺不算红」口径）。

## 5. 本腿门禁（零产码 ⇒ gofumpt/vet/build 无对象）

- `d22scan` clean 核一发：**rc=0**（本腿 11:0x，读数见 §2 表首行）。
- 改动件：票 195（AC#5 追加行＋Progress log 一行）、票 411（Progress log 一行）、本件（新建）；全 LF。
