# 票 261 · 腿 261-r1 证据件（AC#1 走 ⓐ：让承诺兑现）

代理：`261-r1`（写码子代理）。派单裁定：两处落门，都在 `internal/llm/resolver.go`——
枚举点 `DiscoveredModels` 跳过 `Enabled==false`；选择点 `resolveEndpoint` 对目录里有但被关掉的条目具名拒绝。
**选择点是编排者据 AC#0 读数扩的射程**（AC#0 ⓐ 判死的结局是"被选中并产生花费"，只过滤枚举治不到
`text_chain`/`roles.chat` 点名那条路）。文案走 `observe.New(observe.ClassConfig, ...)`，
不复用/不近似 `llm: unknown model %q of provider %q`。

⚠ 状态：骨架（第一枚 commit 时点）。后续节为在飞草稿，终态以最后一版为准。

---

## §0 起手锚与起手名册（现跑自取）

- 起手时刻：`2026-10-03T15:55:46+08:00`（`date -Iseconds` 自取）
- 起手 HEAD：`6c7a12b656188405e7ab40d6bb7c47066ff55dd6`（`git log -1 --format=%H` 自取），分支 `dev`
- 起手 `git status --porcelain internal/llm internal/config`：**空**（rc=0，零行输出）——
  与编排者宣称"我已现跑＝空"一致，**无他人脏面**。
- 起手绿名册：`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test -count=1 -v ./internal/llm/ ./internal/config/`
  → **152 PASS / 0 FAIL / 0 SKIP**；两包 ok（llm 36.408s／config 0.943s），rc=0。
  名册逐名 `.scratch/wisp/probes/261/r1-base-roster-clean.txt`（152 行，剥时延后缀）。
- 与 `261-p1` 终态名册（`.scratch/wisp/probes/261/p1/final-roster.txt`，152 行）逐名 `comm` 比对：
  - `comm -23`（p1 有、我无）＝**0 行**；`comm -13`（我有、p1 无）＝**0 行**。**完全一致，零丢名零增名。**
  - ⚠ 注记：p1 名册是裸测试名（无时延），我最初从 `-v` log 抽出的名册带 ` (0.02s)` 时延后缀，
    首次 comm 被后缀污染（152 行全"差异"）；剥后缀（`sed 's/ (.*$//'`）后重比才干净。
    中间产物 `.scratch/wisp/probes/261/r1-base-roster.txt`（带后缀版）与 `r1-p1final-roster.txt`
    （误从带后缀 log 抽 p1 名册得 0 行的废件）**只建不删**，留作过程痕迹。
- 基线 log：`.scratch/wisp/probes/261/r1-baseline-test.log`

## §0b 起手预检读数（撞钉预检，逐枚）

1. **行号复认**（推翻清单的靶子，逐枚现读）：
   - `internal/llm/resolver.go:287` 逐字 `out := make([]string, 0, len(p.Models))`、`:288` 逐字
     `for id := range p.Models {`——**编排者写"循环体实测在 :287-288"成立**（但我注记：:287 是
     `out :=` 行不是循环行，循环体是 :288-290；语义无差）。
   - `internal/llm/resolver.go:113` `func (r *Resolver) resolveEndpoint(...)`、`:119` 逐字
     `spec, ok := p.Models[model]`——**成立**。
   - `internal/config/schema.go:369-370` 承诺句两行、`:371` 字段——**成立**（票面引文折成两行是拼接转写）。
   - `internal/config/defaults.go` 无任何 `Models` 预置（grep 零命中）⇒ `NewDefaults()` 不产模型条目，
     生产 `wisp run` 的正常路径（零手写条目）不受本门影响。
2. **撞钉预检（fixture 缺 `Enabled` 又期待可达的形状）**：
   全仓 grep `ResolveChain|ResolveRole|DiscoveredModels --include=*_test.go` 命中仅 2 枚文件：
   `internal/llm/enabled_reach_261_test.go`（本票前腿仪器，按派单改写）与 `internal/llm/catalog_test.go`。
   另枚全仓 `ModelSpec{` 测试构造点（13 处）逐枚核：
   - `internal/llm/catalog_test.go`：`:29,36,43,133,209` 全部显式 `Enabled: true`（编排者点名 5 处复认成立）；
     `:217` 逐字 `if !m2.Enabled || ...` 是 `ImportDiscovered` 的判据（不走 resolver，本门碰不到它）。
   - `internal/llm/fallback_test.go`：`:44,51,265` 全部显式 `Enabled: true`（编排者点名 3 处复认成立）。
   - `internal/llm/openaichat/mockllm_integ_test.go:418`：`if !spec.Enabled || ...` 是 discover/import
     路径的读（不走 resolver），不受影响。
   - `internal/config/{validate,loader,catalog,boundary}_test.go` 的 `ModelSpec{}`/缺 Enabled 构造：
     全部只过 `config.LoadFile`/validate（**config 包不 import llm**），resolver 门不在其路径上。
   - `internal/agent`、`cmd/wisp` 的测试：零处构造 `ModelSpec` 后过 resolver。
   - **结论：无一处"建 map 时省掉 Enabled 又期待那枚模型可达"的形状；零枚 fixture 需要补写。**
3. **`LoadFile` 第二参**：grep 全仓 `LoadFile\(.+, [^n]` → 生产调用者全部传 `nil`/变量名均为
   `st`/`res`（SecretResolver），**没有任何调用者把 `*llm.Resolver` 塞进 config.LoadFile**。
   注记：编排者预检清单里"resolver.go:157 `func LoadFile`"一句是**笔误**（resolver.go 里没有 LoadFile）；
   我的 grep 是按真实签名 `internal/config/loader.go:41 func LoadFile(path string, res SecretResolver)` 做的。
4. **`wisp run` 正常路径不受门影响**：`cmd/wisp/run.go:435` 的 cfg 出自 `config.LoadFile`；
   零手写模型条目时 `p.Models` 为空 map ⇒ 无条目可判 false，门为透明。

## §1 改了什么（逐枚 file:line）

（骨架时点：未动笔。终态补齐。）

## §2 两向读数

（终态补齐。）

## §3 变异自证

（终态补齐。）

## §4 门禁读数

（终态补齐。骨架时点基线：`go build ./...` 不跑（派单禁区）；vet/d22scan/gofumpt 待跑。）

## §5 我判不动的地方

（终态补齐。）

## §6 推翻清单

（终态补齐。§0b 已含两笔行号级注记：:287 的语义切分；"resolver.go:157 LoadFile"为派单笔误。）

## §7 量不到的格子＋复量法

（终态补齐。）
