# ci-r1 实现件（阶段 0-1 ｜ CI 超时与工具链腿 ｜ 写腿）

取数窗口 2026-10-26 10:45–11:40 +08；仓库 HEAD `f559f14f`；本机 go1.27.1 / Windows / Git Bash。
写面实碰：`.github/workflows/ci.yml`（件一、件二）＋本件。件三**停手上报**（§3），`internal/ball` 一字节未动。

## 0 三件的落点（一句话版）

| 件 | 结果 | 落在哪 |
|---|---|---|
| 一 六 job 加 `timeout-minutes` | **做了** | `ci.yml` 六处：20 / 20 / 30 / 20 / **55** / **10** |
| 二 升 `staticcheck`（`lint` 真因） | **做了** | `2026.2.1` ⇒ **`2026.2.2`**（一行 ＋ 证据注释） |
| 三 `internal/ball` 那条 `tokens.css` 路径 | **停手上报** | 唯一能过断言的文件**根本没进版本库**，且 `design/**` 零写面 |

## 1 件一：六个 job 的真实耗时 ⇒ 定档

尺：`gh api repos/CarlosShao/wisp/actions/runs/<id>/jobs`，逐 job 取 `completed_at - started_at`（**执行分钟，排队不算**）。
样本：**近 40 次 `ci` run**（`workflow 362354689`，2026-10-14 → 10-25，GitHub 上能翻到的全部 40 次）；四舍五入到 0.01。

| job | n | min | p50 | p90 | **max** | 40 次里跑完几次 | **我定** | 建议 | 差 |
|---|---|---|---|---|---|---|---|---|---|
| `lint` | 40 | 1.13 | 1.67 | 1.82 | **1.90** | **0 / 40** | **20** | 20 | 0 |
| `test-core` | 40 | 1.65 | 2.00 | 2.20 | **3.03** | **0 / 40** | **20** | 20 | 0 |
| `test-windows` | 40 | 5.73 | 8.08 | 11.90 | **14.78** | **0 / 40** | **30** | 30 | 0 |
| `slo-smoke` | 40 | 2.15 | 2.58 | 2.85 | **3.35** | 40 / 40 | **20** | 20 | 0 |
| `slo-full` | 40 | 1.53 | 2.65 | 5.68 | **35.52** | 31 / 40 | **55** | 90 | **−39%** |
| `lint-frontend` | 40 | 0.33 | 0.48 | 0.57 | **0.65** | 40 / 40 | **10** | 15 | **−33%** |

★★**两条必须让机主看到的读数**：

1. ★★**3h23m 那次 run 的成因不是"job 跑太久"，是排队**：run `40238092102` 的 `slo-full` `created_at=22:00:04Z`、`started_at=01:23:02Z` ⇒ **队等 202.97 min（3h23m）**，然后**执行 2.65 min** 就完了（整 run 墙钟 3h26m）。
   ⇒ **`timeout-minutes` 从 job `started_at` 起算，治不了排队**；本腿只能把"真挂住"的 job 从 360 min 默认值上拽下来。全样本里 `slo-full` 队等最长 **694.43 min（11h34m）**（run `38575787259`，执行 1.85 min）——那个形状要的是"runner 别睡"，不是 yaml 里的数字。
2. ★★**三个 job 40/40 从未跑到底** ⇒ 它们的 `max` 是**下界**，不是成本：`lint` 死在 `staticcheck`（导入期就崩，下面的步全没付钱）、`test-core` 死在 portable core 第一枚红、`test-windows` 死在那两枚既有红。**"max × 1.5" 对"从没跑完过的 job"不是合法上限** ⇒ 这三个数按**冷路径**定：两次 `go install` 现编（gofumpt ＋ staticcheck）、`docker compose --build mock-llm` 冷镜像、`third_party` 冷缓存（`actions/cache` 键是 `deps.toml`，键一换就全量取 sherpa-onnx）。

**两处我与建议值不同，按实测改（>50% 的差价已解释）**：
- `slo-full` **55**（建议 90）：实测 max **35.52 min**（run `38575155791`，那天冷依赖；同日一发 25.60 min 是 success）⇒ 55 = 1.5 × 该 max；90 会让"真挂住"白烧 1.5 小时。另注：p90 只有 5.68 min，所以 90 也不是从 p90 来的。
- `lint-frontend` **10**（建议 15）：40/40 全绿、max 0.65 min（全文件最紧的一条），10 已是 ~15×；再高只是把"挂住"多留 5 分钟。

## 2 件二：`staticcheck` 2026.2.1 ⇒ 2026.2.2

**逐字复现（本机，不是引用）**：`staticcheck@2026.2.1` 在 **go1.27.2 / go1.27.3** 下跑 `./internal/ball/`：

```
-: internal error in importing "cmp" (cannot decode "cmp", export data version 5 is greater than maximum supported version 4); please report an issue (compile)
```

＝**与 `docs/audit/2026-10-26-audit-verification.md` §3 的红句同一条**（审计引到 "maximum supported version" 截断，真实尾巴是 `4`）。

**为什么这个版本兼容（两层证据，都是实测）**：

1. **源码层**：`2026.2.1 = v0.8.1` 内嵌 `golang.org/x/tools v0.44.1`，其 `internal/pkgbits/version.go` 的枚举**止步 `V4`**（`V0…V4`，`numVersions-1 = 4`），崩溃点在 `decoder.go:86`。`2026.2.2 = v0.8.2` 内嵌 `x/tools v0.51.0`，同文件多出
   `V5: encodes the index of methods to preserve relative order of nongeneric and generic methods (go.dev/issue/81188)` ⇒ `numVersions-1 = 5`。
2. **矩阵层**（三工具链 × 两二进制，`staticcheck ./internal/ball/`，热构建缓存，GOOS=windows）：

| Go | 2026.2.1（v0.8.1） | 2026.2.2（v0.8.2） |
|---|---|---|
| go1.27.1 | 38 findings / **0** crash | 38 / 0 |
| go1.27.2 | 0 findings / **5** crash | **38 / 0** |
| go1.27.3 | 0 findings / **5** crash | **38 / 0** |

＋**按本步自己的循环形状**（`for mod in . tools/d22scan tools/mockllm`，即 `ci.yml` 那段逐字）在 go1.27.3 下：
`2026.2.1` = **0 findings / 15 crash lines**（3 模块 37 包）；`2026.2.2` = **99 findings / 0 crash lines**（同 3 模块 37 包）。
⇒ **兼容不是"试了最新版不报错"，是"同一把尺把 V5 读出来了"**，且发现数从 0 变 99（gate 恢复有牙）。

**CI 用的 Go 是哪一版（推断链，明写）**：`actions/setup-go@v5` + `go-version-file: go.mod`，而 `go.mod` 是 `go 1.27` / `toolchain go1.27.1`。
- 本机 **go1.27.1 写的是 V4**（上表 0 crash），所以 **CI 的 Go ⩾ go1.27.2**；代理上 1.27 的可用档只到 **go1.27.3**（`goproxy.cn/golang.org/toolchain/@v/list`），故 CI 必在 **1.27.2 / 1.27.3**。
- ⚠ 我没能从 CI 日志里读 `go version`：`gh run view --log` 走 `results-receiver.actions.githubusercontent.com`，本机到那一跳**连不上**（如实写，不假装读了）。这条推断的立足点是"V5 崩溃本身"。
- **没有降 Go**（brief 禁）：病在 consumer（`x/tools` 的 pkgbits 上限），不在 toolchain；且降 Go 会动另外五个 job 的前提。

## 3 件三：`tokens.css` —— 查清"该读什么"，然后**停手**

### 断言到底想看什么
`internal/ball/tokens_table_test.go:1465` 是**三方对账**（票 74 AC#4）：表（`docs/evidence/s1/c21-native-tokens.md`）的 `tokens.css 值` 列 = CSS 文件里 `:root` / `[data-theme="light"]` 的声明 = `internal/ball/tokens.go` 的 `Palette.*`；78 条配色行 ＋ `--font-sans` 必须含 `FontFamily`。
读不到文件是 `t.Fatalf`，注释逐字 `the CSS leg of this check must never skip` ⇒ **这条腿的语义就是"源 → 原生副本"的第三条边**，不是格式检查。

### 四个候选面逐一称重（用 `go test -overlay`，**不写仓**）
| 候选 | 实测 | 判 |
|---|---|---|
| `design/assets/tokens.css`（现值） | `rc=1`：`open ...\design\assets\tokens.css: The system cannot find the path specified`；`git cat-file -e HEAD:design/assets/tokens.css` ⇒ **fatal: 不在 HEAD**；历史删除在 `72b65c614` | **已死**（不是"CI 缺件"，是 HEAD 里就没有） |
| **`design/old/assets/tokens.css`** | **`rc=0` `ok`（PASS）** —— 内容**正好**满足 78 行 ＋ `--font-sans`，文件 287 行 / 156 条声明 | ★**内容对，位置死**：`git ls-files design/old` = **0 枚**（`design/` 跟踪 14 枚、盘上 75 枚）⇒ **CI 的 checkout 永远没有它**。把常量改到这条只会让**本机**变绿、CI 仍红（换一句报错），那是**假修**，比现在更坏（它会掩护"这个文件没入库"这件事） |
| `frontend/src/styles/tokens.generated.css` | `rc=1`：`parsed to 43 :root and 0 [data-theme] declarations: ... the CSS leg would pass vacuously` | **不是同一套词表**（无 `--tint-*`、无 light 块）；且它是 `frontend/scripts/gen-tokens.mjs` 从脚本内一张手写表生成的 ⇒ 读它等于拿**派生件核对派生件**，正是该测试开头拒绝的"fake the very provenance"。另 `frontend/**` 零写面 |
| `design/doubao/demo/styles.css` | 51 条声明、`--background: 210 20% 98%` 这类 HSL 名 | 另一套设计系统，**无 `--tint-*`**，满足不了表 |

**全仓核对**（`git grep -l -- '--tint-success-hi'`）：只有 `docs/evidence/s1/c21-native-tokens.md`、`internal/ball/tokens.go`（注释引用）、测试本身，以及 `.scratch` 里的历史日志 ⇒ **HEAD 里没有任何一枚跟踪文件携带 C21 的 CSS 词表**。

### 所以为什么停手
判据挂在**①已降级为参考的面（`design/old/`）②根本没进版本库的文件③归机主外部 agent 所有的目录**上。三条里任何一条都不是本腿能自己解的：
- 改常量 ⇒ 只治本机，CI 照红（上表第 2 行实测）⇒ 假修，且会掩盖"未入库"；
- 让它读派生件 ⇒ **改的是这条腿的语义**（把"源 vs 副本"换成"派生件 vs 派生件"），属于**改判据**，且 `frontend/**` 零写面；
- 缺件即 SKIP ⇒ 四条硬边界第 4 条明令禁止。

⇒ **摆机主的三条路（我不挑，因为每条都要他的判断）**：
① 他把 `design/old/assets/tokens.css`（或一份副本）**提交到某个跟踪路径** ⇒ 本测试只需改一行常量就能真绿（内容已被我实测证明是对的）；
② 他裁"C21 的 CSS 腿作废或换基" ⇒ 动的是 `A24-D6` 那条权威链的语义；
③ 什么都不做、把它登记为**已知红**（本腿零代码改动，CI 照旧红这一枚）。
⚠ 同一根因还有一枚：`internal/panel/tokens_fourway_test.go`（四way）读的也是这枚文件 —— 一并归上面三条路，本腿同样没动。

## 4 语法与静态检查（改 `ci.yml` 后现跑）

| 检查 | 命令 | rc | 读数 |
|---|---|---|---|
| YAML | `python -c "yaml.safe_load(...)"` | **0** | 六 job 全在，`timeout-minutes` = 20/20/30/20/55/10 |
| `bash -n` | 从 yaml 抽全部**多行 `run: |` 块**（3 块：`lint` step7 gofumpt、`lint` step11 staticcheck、`test-core` step4 probe）逐块 `bash.exe -n` | **0 / 0 / 0** | 含本腿改过的那块（step11） |
| `actionlint` | `actionlint -oneline .github/workflows/ci.yml`（v1.7.12，本机构建） | **1** | 只剩**既有**一条 `label "wisp-slo" is unknown`（自定义 self-hosted 标签）：**改前改后同一条**，只有行号 1014 → 1103 随插入位移；**零新增** |
| 空白 | `git diff --check -- .github/workflows/ci.yml` | **0** | — |
| 编码 | 新增 108 行逐行 `[^\x00-\x7F]` | — | **0 行非 ASCII**（`ci.yml` 无 emoji、无箭头、无圈码写入） |

`gh workflow` dry-run：**没跑**（`gh workflow run` 会真触发一次 run ＝ 真推送测量，本腿只 commit 不 push，且那等于开一次 6-job 全套）；上面的 `actionlint` ＋ YAML ＋ `bash -n` 是本机能给的最接近的替代。

## 5 四条硬边界：没碰的证明

| 边界 | 尺 | 读数 |
|---|---|---|
| ⛔ `continue-on-error` / `if: false` / 删步骤 / 放宽断言 | `git diff -U0` 的 **108 行新增**逐行 grep | `continue-on-error` = **0**、`if: false` = **0**、`if: true` = **0**、`skip` = **0**；**删除行仅 1 行** = 旧的 `staticcheck@2026.2.1` 那行（`2026.1 → 2.2` 的替换），无任何步骤被删 |
| ⛔ `design/**` / `frontend/**` | 本腿 commit 的 pathspec（逐个名字面点名）只有 `ci.yml` ＋本件；对这两个目录**零写操作** | 件三的四个候选**全部用 `-overlay` 读**（`go test -overlay`，仓内零字节改动）；`design/old/assets/tokens.css` 只读被称重 |
| ⛔ `pull_request` 触发器 | `git diff` 里含 `pull_request` 的 `+/-` 行 | **0 行**（`trigger` 块与 `concurrency` 表达式未动） |
| ⛔ "环境不满足即 SKIP" | 新增行 grep ＋ 步骤表 | **0**：没有新 `if:`、没有条件 SKIP、没有 `-skip`；三条既有红照旧红 |

## 6 哪几件做了、哪几件停手

- **做了**：件一（六 job 超时，含两处按实测改档并写明差）、件二（`2026.2.2` ＋ 兼容性证据）。
- **停手上报**：件三。**停手原因不是"不会改"，是"改了只会让本机变绿、CI 照红"，且真相源那一侧不在本编队写面** —— 见 §3 的三条路。
- **没做**（照 brief）：没动 `pull_request`、没加 SKIP、没动 SLO 阈值/golden/断言、没碰 `design/**`/`frontend/**`/`docs/specs/**`/台账/票面。

## 7 门禁 rc（本腿现跑）

| 命令 | rc | 说明 |
|---|---|---|
| YAML `safe_load` | 0 | 见 §4 |
| `bash -n` × 3 块 | 0 | 见 §4 |
| `actionlint` | 1 | **既有** `wisp-slo` 标签告警，改前改后同一条 |
| `git diff --check` | 0 | — |
| `go test -count=1 -run TestC21TableColourRowsMatchTokensCSS ./internal/ball/` | **1** | **未改**：`read design/assets/tokens.css: ... cannot find the path`（件三停手，见 §3；overlay 到 `design/old/` 时 rc=0，仅作为"该读什么"的读数） |
| `staticcheck` 三模块循环（go1.27.3） | 1 | 旧版 **0 findings / 15 crash**；新版 **99 findings / 0 crash** —— **仍是 rc=1，且这是对的**：`lint` 绿不是本腿的判据（票 85a／122 明文），本腿只挪掉"工具链读不了格式"这一层 |
