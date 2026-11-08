# 裁决表：整包 runner 把 `internal/speech` 路由过 staging carrier（票 409 / wave99 speech-runner 腿）

- 腿：`speech-runner`（写腿；写面＝`scripts/portable-tests.sh` ＋ `scripts/portable-tests-selftest.sh` ＋ `scripts/speech-tests.sh` ＋ 本件；零 Go 产码、零断言、零 thresholds/golden、零 CI workflow、零 System32 触碰）
- 起手锚：`66a03a320`（派单时 HEAD）；落笔时工作树 HEAD 漂到 `eccae1fd7`（他腿并行落笔，本腿 diff 与两枚脚本自持字节为准）
- 任务来源：台账 `A1645` §三「整包 runner 该路由」裁定（该做）＋ `A1644` §一「语音族 39」归因。

## 1. 三核（派单要求，全部现跑）

1. **两份成品脚本全文现读**（⛔ 重写）：`scripts/stage-sherpa-dlls.sh`（85 行）＝读 `deps.toml` `[sherpa-onnx.dll."…"]` 钉值 → GUARD → 拷贝；`scripts/speech-tests.sh`（87 行）＝`go test -c` → stage 三枚 DLL 到测试二进制旁 → 包目录 cwd 真跑 → runtests 形四数。本腿**只做最小增量**（见 §3）。
2. **grep 现有 runner 对 speech 的处理**：`scripts/portable-tests.sh:268-281`＝「`./internal/speech/` IS DELIBERATELY NOT CLAIMED…census row is reported to the owner」（彼时 carrier 不存在，认领即永久红）；CI（`ci.yml:764`）只挂 `runtests.sh ./internal/speech/cer/`（纯 Go 子包）；`ci.yml:1057` census 步是 GUARD D 唯一 CI 调用点。
3. **整包 runner 真名**＝`scripts/portable-tests.sh`（ticket 93/111/250/251/254 系列的载体）。39 语音红的真实来源＝编排者自己的全仓复跑（`A1644`：`./cmd/wisp ./internal/...` 形，四数 RUN=707/PASS=497/FAIL=51/SKIP=4，51 红里语音族 39）——census/scope 族并不把 speech 计入分母，治法是把 speech **变成本地整包跑法的一档**，与派单里「最简做法」一致。

## 2. 结构判定（派单「如实」条款）

`portable-tests.sh` **按包路由结构上成立**：tier＝scope 列表＋pin＋`case $mode in` 分支，census 自检（ticket 251 AC#4）只要求 roster 与 branches 同集合。故「把 `internal/speech` 摘出裸 `go test` 名单、路由过 carrier、四数并入」在既有骨架内就是一条合法 tier——不需要 `-exec` 包装（派单已证无效）、不需要把 DLL 预先 stage 到 `$PWD/build`（对 `go test -c` 产出的二进制无效，该二进制落在 stage 目录而非 build/）。

## 3. 落刀（三文件，全部最小增量）

- `scripts/portable-tests.sh`（+98/−3）：
  - `speech_pin='github.com/CarlosShao/wisp/internal/speech'`（:374）＋ `tiers=` 名册加 `speech`（:386）＋ `case` 加 `speech)` 分支（scope=`./internal/speech/`，:440）＋ census pin 循环加 `speech`（:548）——四件同一 commit，GUARD C/D 与 census 自检互锁。
  - 2026-11-08 dated update 段：原「deliberately not claimed」段**原句不抹**，其后追加 supercession 说明（carrier 已移除其前提）。
  - run 步：`mode=speech` 时**不走 strict runner**，改调 `scripts/speech-tests.sh`，stdout 经 tee 进同一 capture（:961-977 附近）；strict runner 路径对其余 tier 逐字不变。
  - GUARD B speech 形（:1019-1035 附近）：carrier 真跑的是裸测试二进制，不存在逐包 `ok <path>` 结果行 ⇒ speech 档的「tested」判据＝`-v` 分母非零（PASS=0 且 FAIL=0 即拒，exit 1）；残留 SKIP **不当 fatal**（理由：carrier 自身 exit code 已携带全部断言判决；6 枚残留 skip 点名标点模型缺件＝fixture 类环境缺口，非放宽断言；strict runner 的 skip 规则对其他 tier 一字未动）。
- `scripts/speech-tests.sh`（+12/−4）：`-v` 流由 `>"$log"` 改为 `2>&1 | tee "$log"`＋`rc=${PIPESTATUS[0]}`——stdout 保留同一 `-v` 字节流（整包四数来自它），log 文件与 skip 读数不变。
- `scripts/portable-tests-selftest.sh`（+1/−1）：`unknown-scope-still-hard` case 的 `known:` 钉补上 `speech`（roster 合法扩张所致，非仪器弱化——该 case 断言的是 unknown scope 仍 exit 2 且**点名全部现存 tier**）。

## 4. 验收读数（全部本腿现跑，台面＝`git clone --shared` 导出树 `/d/tmp/wisp/20261108/speech-runner-tree`，untracked 资产 `third_party/spike-models`＋`third_party/sherpa-onnx` 逐件拷入）

| 发 | 命令 | 读数 |
|---|---|---|
| 主验收 | `bash scripts/portable-tests.sh --scope=speech`（导出树，PATH 带 build+third_party/sherpa-onnx） | **RUN=106 PASS=83 FAIL=0 SKIP=6，FINAL_RC=0**；日志含 `mode=speech routes [./internal/speech/] through the staging carrier`；guard 链（310 root → resolve → GUARD C 钉比对 → GUARD A → ledger staleness → 路由 → 四数）全过 |
| 作差（治前） | 台账 `A1644`：语音族在整包＝**39 枚假红**（遮蔽→47 skip→strict runner 记红） | **39 假红 → 0**；同一口径（carrier 单跑）修前 PASS=42/SKIP=47，修后 PASS=83/SKIP=6（+41 转真跑绿，`A1645` §二） |
| 残留 skip | 主验收日志逐名 | 6 枚全部点名「punctuation model not available: set WISP_PUNC_MODEL …」（标点模型缺件），**遮蔽文案 1.17.1 零命中** |
| census | `bash scripts/portable-tests.sh --scope=census`（导出树） | `internal/speech 16/0 speech`（**UNCLAIMED-HAS-TESTS → CLAIMED BY speech**）；census 自检（roster⇄branches）rc0；GUARD D 剩余 4 枚（`cmd/balldebug`／`internal/agent/scheduler`／`internal/command`／`internal/plugin/manifestslot`）＝**与 speech 无关的存量洞**，本腿⛔ 越权代领，归后续各票 |
| 仪器反噬 | `bash scripts/portable-tests-selftest.sh`（导出树） | 首跑 32 case/1 断言失败＝`unknown-scope-still-hard` 的 roster 钉过期（§3 第三条）；按钉面修正后 **32 case/0 失败 GREEN**——GUARD C/D/A/B 与 census 全部仍咬合 |

## 5. 门禁（导出树，`unset GOTMPDIR TMPDIR TEMP TMP`；gofumpt 绝对路径）

| 门 | 读数 |
|---|---|
| `gofumpt -l .`（`/d/work/base/gopath/bin/gofumpt.exe`） | rc=0，28 行全为 `.scratch`/`forks` 存量，`scripts/` 命中 **0** |
| `go vet ./...` | rc=0 |
| `go build ./...` | rc=0 |
| `sh scripts/d22scan.sh` | rc=0，`clean - no D22 ban violations` |
| 三枚脚本 LF | `portable-tests.sh`/`speech-tests.sh`/`portable-tests-selftest.sh` CRLF 行 **0/0/0** |

## 6. 如实

- census 的 GUARD D 在 speech 闭合后仍红 4 枚——那是**本来就存在的洞**（census 修前也会报，只是被 speech 行挤在旁边），本腿修它们＝越权；逐名已列 §4。
- 「整包 579s 到终态」的全仓复跑属编排者，本腿按派单⛔ 不重跑；能覆盖 speech 路由的那一档（`--scope=speech`）已现跑到终态。
- 本腿首跑 selftest 有 1 断言红——那是本腿 roster 扩张的**直接后果**、按钉面修正（不是把仪器调松：case 的主断言 unknown scope exit 2 原样保留）。
- 撤销口令：`git revert <本腿提交>` 即回退三文件；语音治本体＝`git revert 656a2ea25`（台账已记）。
