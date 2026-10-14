# 313-r3 锚＋门禁＋自报（票 313 挂账②：多元素逗号链补钉·写腿）

判语（供编排者复核）：正钉一枚落 `cmd/wisp/panel_config_313_test.go`（**append-only，+69/-0**）；
两发突变（按 `;` 切／只切第一枚）各证牙＋红句逐字；门禁 6 枚 rc 全 0；⛔ 未改产码、未 push。

## 0. 锚

- 起手 `git rev-parse HEAD`＝`b8cb0b8f`（派单写 `b8f883a3`——共享树在呼吸）；落笔期间他人提交
  （260-v2 勘误笔）把 HEAD 推到 `351214de`。本腿 diff 只含下节两处，与二者无关（见 §4）。
- 触件 sha256：`cmd/wisp/panel_config_store.go`＝`43acff43…`（**全程未变**，两发突变窗口后逐字节还原）；
  `cmd/wisp/panel_config_313_test.go`＝`db7d4050…`（本腿新增 69 行）。逐枚见 `logs/00-head.txt`。
- ⛔ 零 push。

## 1. 落点与正钉（现读后选最小形）

- **选 `cmd/wisp/panel_config_313_test.go`（端到端），不选 `internal/config/settings_313_test.go`**：
  逗号语义只活在 `splitChainValue`（`cmd/wisp/panel_config_store.go:328`，door 侧）；config 层 setter
  直接吃 `[]string`，在那里加钉对本次突变**无牙**（不可能红）⇒ 不满足派单 §2「新钉必须红」。
- **新用例名**：`Test313MultiElementCommaChainKeepsEveryElement`（具名 `-run` 现跑绿）。
- **值**＝`"deepseek/deepseek-chat, deepseek/deepseek-reasoner"`（逗号后带一枚空格：一并钉逐元素
  `TrimSpace`）；两枚模型名**现读自仓内既有 fixture**（`settingsSeed` 的 `deepseek-chat`；
  `deepseek-reasoner`＝313/257 种子里那枚第二模型，`internal/config/settings_313_test.go:244` 一带、
  `cmd/wisp/firstrun_257_test.go:41`）⇒ ⛔ 未编名。
- **判据三面**（要求 ≥2，本腿三面全给，观察法抄既有 313 用例）：回执（`apply313`——键路径
  `voice.cloud_asr_chain`＋重启文案）／落盘（config.toml 数组行**逐字**锚定
  `cloud_asr_chain = ['deepseek/deepseek-chat', 'deepseek/deepseek-reasoner']`）／读面载荷
  （`config.get` 的 `voice.cloudASRChain` 两枚按序）。落盘检查锚**数组整行**而不查全文字符串片段——
  种子自己的 model 段也含 `deepseek-reasoner`，片段查法会对「只落一半」假绿。
- **行为层现读本就工作**（预期如此）：baseline 与新钉均绿 ⇒ 本钉是覆盖补强，非缺陷修复。

## 2. 门禁（rc 逐枚；logs/）

| 门 | 命令 | rc |
|---|---|---|
| 编译 pre/post | `GOFLAGS= go build ./...` | 0/0 |
| 格式 | `gofumpt -l`（我触过两件） | 空 0 |
| Test313 三包 pre/post | `go test -count=1 -run 'Test313' ./internal/config/ ./internal/panel/ ./cmd/wisp/` | 0/0 |
| 新钉具名 | `-run 'Test313MultiElementCommaChainKeepsEveryElement'` | 0 |
| d22scan | `sh scripts/d22scan.sh` | 0 clean |
| 长度门 | `sh scripts/check-path-length-budget.sh` | 0 GREEN 59=59 |

- ⛔ 未跑整包 `./cmd/wisp/`（314 死区）。315-a1 在飞探针件 `cmd/wisp/zz315a1probe_test.go` 在树内，
  本腿小射程 `-run` 不触发其用例；⛔ 未动它（其 test 名 `TestZZ315A1…` 与本腿 `Test313…` 零重叠）。

## 3. 突变窗口（仓外备份 `C:\tmp\wisp313r3backup\panel_config_store.go.orig`＋sha256 逐字节证）

原/还原 sha256＝`43acff433ae72f2677b4efb6e3fdab9a35a6a393a89ed84107eaa647a9119a9d`（前后相等）。

- **窗 M1（按 `;` 切，照 313-v1 的 B 形）**：起 `10:52:14` → 还原 `10:52:44`；mutated sha256
  `9e8f63c1…`。新钉红，红句逐字（`logs/20-mutM1-red.txt`）：
  `panel_config_313_test.go:306: config.set was refused: config: config: refused to write voice.cloud_asr_chain。 第 2 种拒因：行不存在：链元素点名的 provider/model 必须已经在目录里（先用建行那一步，或手改文件）。: config: config.toml: voice.cloud_asr_chain element "deepseek/deepseek-chat, deepseek/deepseek-reasoner" references unknown model "deepseek-chat, deepseek/deepseek-reasoner" of provider "deepseek"`
- **窗 M2（只切第一枚 `SplitN(v, ",", 2)[:1]`）**：起 `10:53:05` → 还原 `10:53:47`；mutated sha256
  `7263362d…`。新钉红两发（`logs/21-mutM2-red.txt`，全文 body 原文留 `22-mutM2-red.raw.txt` 不入库）：
  `panel_config_313_test.go:318: config.toml does not carry the two comma-separated elements of "deepseek/deepseek-chat, deepseek/deepseek-reasoner" as "cloud_asr_chain = ['deepseek/deepseek-chat', 'deepseek/deepseek-reasoner']":`
  `panel_config_313_test.go:330: voice.cloudASRChain read back as [deepseek/deepseek-chat], want both elements of "deepseek/deepseek-chat, deepseek/deepseek-reasoner" in order`
- **窗 M1R（M1 重跑留 raw 档）**：起 `10:55:47` → 还原 `10:55:56`（sha256 同 M1）。
- 每窗还原后具名 `-run` 复跑**绿**（rc=0）；还原方式是 `cp` 仓外备份（⛔ 未用任何 git 状态命令）。

## 4. 「只加不改」怎么成立

- `git diff --stat -- cmd/wisp/panel_config_313_test.go`＝`1 file changed, 69 insertions(+)`；
  以 `^-` 行计数＝**0** ⇒ 存量用例一字未动（连文件头注释也未碰——本钉全部上下文在新增函数的注释里）。
- 存量 313 件 diff 现读为空：`internal/config/settings_313_test.go`、`internal/panel/config_313_test.go`、
  `internal/config/settings_313.go`、`cmd/wisp/panel_config_store.go`（含 `splitChainValue` 本体）。
- ⛔ 未动三枚冻结件／`thresholds.go`／golden／`allowlist.txt`／`PLAN.md`／`specs/**`／`.github/**`；
  `frontend/**`／`design/**` 零读零写；⛔ 未改任何产码语义（若必须改产码才能过＝停手上报，本次未发生）。

## 5. ⛔ 没做的格（自报）

1. config 层链级（元素级）用例未加：对本次突变无牙、非最小形（§1 理由）——若编排者认为要补，另派。
2. 整包 `./cmd/wisp/` 与 `-race` 未跑（314 死区／未要求）。
3. 未 push；未翻票框、未动 `Status:`／票名／AC 框；票尾只追加 1 条 bullet。
