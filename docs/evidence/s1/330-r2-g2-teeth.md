# 票 330 · 落地腿 330-r2 · 缺口 G2 的牙与 G3 的注释（测试面，⛔ 产码字节）

- 起手锚（本腿第一条命令现取）＝`feab45e3abee588ae30245c820a762d3b695b91f`；交件时 HEAD＝本腿第四枚提交。
- 写面＝`internal/agent/plan_present_crowded_world_330r2_test.go`（新建 272 行）＋
  `internal/agent/plan_present_doorbell_165r3_test.go`（G3 注释：1 删 3 加，221→223 行）＋
  `.scratch/wisp/probes/330/r2/**`＋本件＋票 330 末尾一条 bullet。⛔ 碰 `cmd/wisp/**`、⛔ 产码、⛔ golden 内容。
- 台面＝仓外导出树 `D:\tmp\w330pre`（锚 feab45e3）与 `D:\tmp\w330pst`（锚 560945c2），
  配方＝`git archive <锚> go.mod go.sum internal cmd tools scripts frontend models docs`＋手工三件
  （`build/*.dll`、`third_party/sherpa-onnx/*.dll`、`frontend/dist`）＋`GOTMPDIR=D:\tmp\w330tmp`（先 mkdir）；
  目录名 8 字符（≤8.3 段）；取出后 `.go` 计数＝747（pre）／749（post）证非空。脚本＝`export-and-run.sh`。

## 1. 起手三把尺（本腿自己现跑，⛔ 抄派单行号）

- `sha256sum` 基线：`testdata/golden/plan-present.sse`＝`7d26fa889b59fc3a…`（28 行，本腿⛔ 动它一字）、
  `plan_present_doorbell_165r3_test.go`＝`00c32049994c441c…`（221 行）。
- `grep -rn 'func bm25Select' internal/agent`＝`internal/agent/inject.go:172`；`:219` 那枚
  `if score > 0 { hits = append(...) }` **今天在场** ⇒ `G2` 的机理成立，本腿按派单继续（未触发"停下上报"分支）。
- 镜像条目 `Resident: true` 现位＝`plan_present_doorbell_165r3_test.go:69`（本腿⛔ 改它一字）。

## 2. 造出来的世界（为什么档位这时才说话）

`crowdedDoorbellWorld330r2()` ＝ 真镜像条目（`planPresentDirectory165r3()` 组合而来，⛔ 第二份抄本）
＋**一枚**不驻留的长尾候选 `ledger.sweep330r2`（夹具检索专用，⛔ 声称是真名册条目——那是 G1／330-r1 的活）。

- 尺＝`bm25Select` 的名额削法（`inject.go:232-243`）：第一名无条件入选，其后 `used+cost > budget` 即断；
  `budget = SectionToolIndex = 300`（`BudgetsFor(128000)`，`scaleInt` 在 factor>=1 时原值返回）。
- 现量：长尾那条索引行 cost＝**291**、镜像那条＝**20** ⇒ 单独都装得下、两条合计 311 装不下。
- 于是档位＝唯一切换：镜像条目在驻留支 ⇒ 无条件上线；被拨到检索支 ⇒ 与长尾抢同一枚名额、抢输 ⇒ **整枚不上线**。
- 用例 1（`TestCrowdedWorldStillOffersTheBrake330r2`）走真 Loop／真 openai-chat 适配／真 store，
  断言里**⛔ 提名长尾候选** ⇒ 把新增那枚候选删掉，它就逐字退回 165 的单候选世界（＝第 5 节必答题的根据）。
- 用例 2（`TestResidencyTierDecidesTheInjectionArm330r2`）首枚在本包直接读 `Inject` 两支臂的用例
  （本锚点 `grep -rn 'Inject(' --include=*_test.go .`＝0 命中；产码调用点只有 `loop.go:573`）。
  门序＝前提(枚数)→前提(预算)→镜像驻留臂→长尾不入驻留臂→长尾入检索臂，故意让三发突变各打不同一行。

## 3. 颜色四数（`go test ./internal/agent/ -count=1 -timeout 30m -v`，仓外树，本腿现跑）

| 发 | rc | RUN | PASS_TOP | FAIL_TOP | SKIP | PANIC |
|---|---|---|---|---|---|---|
| 改前＝锚 feab45e3（起手分母，自跑复现） | 0 | 110 | 86 | 0 | 0 | 0 |
| 改后＝锚 560945c2（含本腿两枚文件） | 0 | 112 | 88 | 0 | 0 | 0 |

- 作差：RUN +2／PASS_TOP +2／**终态红名册与起手作差＝0**（两枚红名单皆空）。同包别枚用例⛔ 一枚被打红 ⇒ 无需上报裁定。
- 正控（三枚点名用例在未突变树上）：`rc_forward=0`，三枚全 PASS ⇒ 件 `logs/m0-forward-unmutated.log`。

## 4. ★成对定向突变（全部只打在仓外导出树字节上；逐发当场还原，`finally`＋哈希核回）

原始读数＝`logs/m1-g2-replay.log`／`m2-defang.log`／`m3-control.log`／`m4-deleted-only.log`／`m5-deleted-plus-g2-replay.log`，
驱动器＝`mutations-330r2.py`（还原失败即 exit 4；五发跑完 `logs/tree-final-hashes.log` 证树内字节＝提交内字节）。

- **M1＝G2 那一发的重放**（`plan_present_doorbell_165r3_test.go:69` `Resident: true → false`）⇒ `rc=1`：
  - `plan_present_crowded_world_330r2_test.go:188: 330r2 crowded world: round 1 body must offer the brake: missing "\"name\":\"plan.present\"" in "{…`（请求体截断，全文在 .log）
  - `plan_present_crowded_world_330r2_test.go:240: residency switch ignored: resident arm = [], want "plan.present" in it - the mirror flags it Resident: true and Inject put it somewhere else`
  - 同发里 `--- PASS: TestPlanPresentCallBooksTheToolCallRow165r3 (0.05s)` ＝ 165 那枚用例本身**依旧看不见**档位
    （G1/G2 原判语的形状未被本腿推翻，本腿的牙在新用例上）。
- **M2＝去牙发**（新夹具那枚不驻留候选 `Resident: false → true`）⇒ `rc=1`，唯一红：
  - `plan_present_crowded_world_330r2_test.go:247: tier leak: resident arm = [plan.present ledger.sweep330r2], want "ledger.sweep330r2" kept out of it - the fixture declares it Resident: false and only the tier decides this arm`
  - `--- PASS: TestCrowdedWorldStillOffersTheBrake330r2`（用例 1 不靠档位之外的东西红）。
- **M3＝反控发**（`crowdedWorldQuery330r2` 换成与长尾**无共同字**的 `"frobnicate widget quietly"`）⇒ `rc=1`，唯一红：
  - `plan_present_crowded_world_330r2_test.go:254: retrieval arm = [], want "ledger.sweep330r2" in it: it scores above zero against "frobnicate widget quietly", so the capped slot is the only thing that could have dropped it`
  - 红句与行号（254）与 M2（247）**不同** ⇒ 这枚用例响的是档位与检索各自的门、⛔ 恒响。
- 还原哈希（五发各自 `finally` 现算，逐枚与提交内一致）：
  `plan_present_doorbell_165r3_test.go` → `6542773d3d2edcc6d554a47f091ce2dffe35586a92c1f948f79c361454683fca`、
  `plan_present_crowded_world_330r2_test.go` → `891fbf5ccc9bdcd23f5075538c2a590e441dc87a100c2634a45d29cbbdd2b59d`。
  基线件⛔ 动一字：`testdata/golden/plan-present.sse` 仍 `7d26fa889b59fc3a…`；产码 `inject.go` 树内外同串 `937bcc113b6018b7…`。

## 5. 必答题（编排者指定那一发）

把**新增的那枚候选删掉**（M4＝只删夹具、⛔ 档位突变；M5＝删掉之后再重放 G2 的 `Resident` 拨 `false`）：

- **M5 里 `TestCrowdedWorldStillOffersTheBrake330r2` ＝ `--- PASS`** ⇒ 对"世界那一半的牙"答**是**：
  档位突变退回单候选世界后再次不可见（`inject.go` 那支检索把唯一候选无条件捞进来），牙确实挂在多候选世界上。
- M5 里 `TestResidencyTierDecidesTheInjectionArm330r2` 仍红，但红在**枚数前提门**：
  `plan_present_crowded_world_330r2_test.go:221: 330r2 premise: the crowded world must be the frozen mirror plus exactly one retrieval decoy, got 1 entries [plan.present]`
  —— 接住它的是这扇门（用例 2 的断言按名字提名了那枚被删的候选），⛔ 是档位门 `:240`。M4（未突变）同红同一行 ⇒
  这枚红与档位无关，属"夹具枚数完整性"门，具名交此，不当作 G2 的牙。

## 6. ★零产码字节那把自证尺（`logs/self-rulers.log`，逐字跑派单给的两条）

- `git diff --numstat feab45e3 HEAD -- internal/agent` ⇒ 只有两行：`272 0 …crowded_world_330r2_test.go`、
  `3 1 …plan_present_doorbell_165r3_test.go`。
- `git diff --name-only feab45e3 HEAD -- internal/agent | grep -vE '(_test\.go|^internal/agent/testdata/)'`
  ⇒ **空、`filtered_rc=1`** ＝ 零产码字节。
- ★正控（⛔ 共享工作树，跑在 `mktemp -d /d/tmp/w330pos.XXXXXX` 的 `git clone --local` 副本里）：
  把副本 `internal/agent/inject.go` 的一处非注释字面量 `MaxRetrievedTools` 改成 `MaxRetrievedTools+1` 后跑同尺
  ⇒ 输出 `internal/agent/inject.go`、`control_filtered_rc=0` ⇒ 这把尺两种形状都分辨得出、⛔ 失明。

## 7. 门禁四数（只在共享工作树；`logs/gates.log`）

| 尺 | rc | 备注（本腿现跑） |
|---|---|---|
| `gofumpt.exe -l internal/agent/` | 0 | 列出 `approval/pending_read.go`、`tools.go` 两枚＝**本锚点既存漂移**（`git status --porcelain` 对这两枚为空、本腿⛔ 碰过）；本腿两枚文件不在名单内。 |
| `sh scripts/d22scan.sh` | 1 | 唯一 finding＝`cmd/wisp/panel_snapshot_deliver_329.go:14 [phantom-citation]`（并行腿 329-r1 已提交字节）；`mentions_of_this_legs_new_file=0`。起手第一次跑还多一条 `panel_pump.go:19`，随后由 329-r1 自己修掉 ⇒ 共享面在飞，本腿⛔ 判它红。 |
| `go run . -root .`（d22scan 第 2 步单跑） | 1 | 同上唯一一条，逐字见 `logs/gates.log`。 |
| `sh scripts/check-path-length-budget.sh` | 0 | `VERDICT GREEN`，14425 条 tracked 路径、超帽 60 条全部在册。 |
| `go vet ./internal/agent/` | 0 | 输出 0 字节。 |

## 8. 本腿⛔ 关掉的（归口）

- `G1`（镜像⇄真名册零对接）＝票 330 `AC#1/#2`，形甲落 `cmd/wisp/**`，归 `330-r1`；本腿⛔ 碰那面一字。
- `Declared` 那一支＝`165-v2` 已量得**有牙**（M-b2 唯一红 `:157`、M-g 五枚红），本腿⛔ 重复装钉。
- `G3`＝只把 `:131` 那句注释说实（"冗余非零门，牙在 :145 的 `rowOfTool165r3` 独立遍历"），
  ⛔ 删门、⛔ 重放"把 sink 计数喂成常量"那一发（派单只要求改注释；`165-v2` 已量，本腿未复测，不写成本腿读数）。
