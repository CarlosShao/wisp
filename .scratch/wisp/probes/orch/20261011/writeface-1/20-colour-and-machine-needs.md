# writeface-1 — 20 颜色面与机器需求（归口用）

"颜色"＝那一枚包的 `go test` 逐名名册；"机器需求"＝这一格在**本机、无桌面/无机主/无 CI** 时跑不跑得出来。
本腿⛔ 跑任何颜色（零 go 命令），下面全部是**票面判据文字＋静态尺**给的归属，⛔ 本腿代跑＝故意。

## 逐腿要哪一枚包的颜色

| 腿 | 必读包（票面逐格点名） | 尺名/出处 |
|---|---|---|
| `296-r2` | `./cmd/wisp/` **定向两族**：`-run 'Test296\|Test258'`；另 `GOFLAGS= go build ./...` | 票 296 `AC#5` 判据逐字（只要求 `-run` 那两族，⛔ 整包名册） |
| `295-r1` | `./cmd/wisp/` **＋** `./internal/audio/`，改前改后各 ≥2 发取交集、逐名作差 | 票 295 `AC#4` 逐字两包 |
| `294-r1` | `./cmd/wisp/` 整包，改前改后各 ≥2 发取交集 | 票 294 `AC#5` 逐字 |
| `309-r1` | `./internal/tools/` ＋ `go build ./...` ＋ `go vet ./internal/tools/` ＋ `sh scripts/d22scan.sh` | 票 309 `AC#4` 逐字 |
| `298` 残余 | ⛔ 无需颜色（落地半边零枚文件，见 `10-…` E 节；`298-v1` 已自跑三发＋overlay 发） | — |

★**同色冲突的根**：`296-r2`／`295-r1`／`294-r1` 三枚取的是**同一枚 `./cmd/wisp/` 名册**，`309-r1` 取 `./internal/tools/` ⇒ 三枚 cmd/wisp 腿的名册**互相背锅**，309 与它们⛔ 不背锅（⛔ 同包）。

## 本机跑不出来的格（＝编排者的欠账，⛔ 算腿的）

| # | 格 | 为什么本机⛔ 算交付 | 该由谁 |
|---|---|---|---|
| M-1 | `294-r1 AC#2`（真球窗现读 `IsLive(mute 镜像)`） | 要**交互桌面会话**：`startResidentBall` 起真窗，`SKIP-LOUD` 是 headless 退路。该退路逐字只在 3 枚件里（`resident_hotkey_258_windows_test.go`／`_296_`／`live_258_`，尺＝`git grep -n "SKIP-LOUD" HEAD -- cmd/wisp`＝**12 命中**，每文件 4）＝三枚件都⛔ 能裸 headless 交 | 本机可跑，⛔ 无人值守；票 294 裁定节逐字"SKIP-LOUD 在本机不算交付"⇒ **编排者或机主在场一发** |
| M-2 | `295-r1` 若把 `./cmd/wisp/` 整包名册取全 | 同 M-1，且 `resident_audio_247_live_windows_test.go`（`295-r1` **必动那 1 枚读点**）要**机主的真麦克风** | 机主在场；⛔ 记在 `295-r1` 账上 |
| M-3 | `295-r1 AC#5`（`-race` 墙钟与颜色） | ⛔ 门禁名册里根本没有 race 这一把：尺＝`git grep -nE "\-race" HEAD -- scripts .github tools` ⇒ **0 命中**（本腿自跑，换票面那把 `--include` 形复现同结论）⇒ 跑它＝一次性代价，⛔ 进门禁＝**人工批准面**（票 295 逐字"不许动 `.github/workflows/ci.yml` 一字"） | 编排者＋机主一句话 |
| M-4 | `296-r2` 的 `Test258` 那一族 | 同 M-1：`-run 'Test296\|Test258'` 命中的两族都起真窗 ⇒ ⛔ 无人值守会退成 SKIP，而 SKIP⛔ 是读数 | 本机有桌面会话即可 |
| M-5 | `298` 现量节那枚"CI 两步长期 skipped"的**根因支** | 只能由 **push 后的 CI** 给色：CI 名册现量＝`lint`/`test-core`/`lint-frontend` 跑 `ubuntu-latest`，`test-windows`/`slo-smoke` 跑 `windows-latest`，`slo-full` 跑 `[self-hosted, wisp-slo]`（尺＝`git show HEAD:.github/workflows/ci.yml \| grep -nE '^  [a-z0-9-]+:$\|runs-on:'`）；`./cmd/wisp/` 在 CI 那侧走 `scripts/portable-tests.sh --scope=cli`（由 `wisp-cli-tests.sh` 调），⛔ 本机任何绿红顶替不了那一发 | 编排者车道（push＋`gh` 取日志）——`298-v1` 的"只有真机／只有 CI 才能闭合的格"节已具名同一条，本腿复认⛔ 重述 |
| M-6 | `298` 的"改前两发同台面"严格形 | 改前内容只活在历史 blob 里 ⇒ 本机只能 `-overlay`／导出树**近似** | `298-v1` 已交近似形并具名形差；严格形＝窗口已关 ⇒ ⛔ 再派腿补 |

## 三态与口径备注

- 本件所有"票面逐字"读的是 **HEAD blob 的票面**（`git show HEAD:.scratch/wisp/issues/<票>.md` 落 `logs/`，工作树⇄HEAD 对那五枚票**逐字相同**＝`git status --porcelain -- <五枚路径>` 0 行）。
- 枚数口径：M-1 的 **12**＝整族（含注释行，射程 `cmd/wisp`，含 `_test.go`）；`-race` 的 **0**＝pathspec 形（`scripts .github tools`）。同一条尺去掉 pathspec 换成全仓 tracked ⇒ **162 枚文件**命中，本腿⛔ 逐枚判性质（绝大多数在 `.scratch/**` 与 `docs/**` 的转述里）⇒ 它⛔ 支撑"有人在跑 race"，但**射程那一维必须并写**（票 298 `AC#3` 立的规矩）。
- 层级⛔ 混着数：本件没引任何 `--- PASS` 层级数（⛔ 跑），凡引用一律带 `^` 与否的口径；`298-v1` 件里的 35/0/0 是 `runtests.sh` 的**顶层**口径，本腿只转引它"这件事存在"、⛔ 当现势。
