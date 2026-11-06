# 票 85（wave34 / W85 腿）：lint 两把尺的门禁核验 ＋ AC#4 步骤改名

- 腿：`W85`（**实现腿**；本件是本腿读数，不是裁决表——裁决者≠实现者）
- 时间：2026-11-06 16:1x–16:4x +08
- 基线：开工 HEAD `59ac99f92`（A1570 收 wave33 的那枚；派单原文写 `6ae44b1fb`，系派单起草后 HEAD 又前进一枚）
- 写面（实际触碰）：`.github/workflows/ci.yml`（**只有 AC#4 改名**，4 增/2 删）· 本件
- 本机：go1.27.1 windows/amd64；`/d/work/base/gopath/bin/gofumpt.exe` = `v0.12.0 (go1.27.1)`；
  `/d/work/base/gopath/bin/staticcheck.exe` = `2026.2.2 (0.8.2)`
- 环境纪律：`PATH` 前置 `build/`＋`third_party/sherpa-onnx`；`GOTMPDIR/TMPDIR/TEMP/TMP` 全 unset；
  全程 **CGO 默认开**（未设 `CGO_ENABLED=0`）；未碰 `docs/reports/pending-and-issues.md` 一字。

---

## 0. 结论速览

| 格 | 判定 | 判据 |
|---|---|---|
| `gofumpt` 已钉版本（派单第 1 件事） | **早已在 HEAD，本腿复核 PASS** | `ci.yml` 安装行 `go install mvdan.cc/gofumpt@v0.12.0`（commit `57e764f82`，85-pin 腿 2026-11-05）；本腿现跑 `@latest` 命中 **0** 行 |
| `staticcheck` 能跑、能产出真实判据（派单第 2 件事） | **早已在 HEAD，本腿实证 PASS** | pin `@2026.2.2`（commit `702c13f51`，ci-r1 腿）；本腿现跑 `./internal/ball/`＝**37 条 finding、0 条 crash line**、rc=1（红在真 finding 上） |
| AC#4 步骤改名（上一腿 STOP-3 留给"下一腿"的活） | **本腿完成（本地一半）** | 名字 `PathResolver junction placeholder (real cases tickets 18/20)` → `PathResolver junction real case (tickets 18/20; wiring first-green tickets 72/75)`；**只改名字与同名注释，零行为改动**；真 run 结论＝编排者 push 后取 |
| `gofumpt -l internal cmd tools scripts frontend` | **PASS（空）** | rc=0，0 行 |
| `go build ./cmd/... ./internal/... ./tools/...` | **PASS** | rc=0 |
| `sh scripts/d22scan.sh` | **PASS** | rc=0，`clean - no D22 ban violations`，runtests `PASS=35 FAIL=0 SKIP=0` |
| `go vet ./...` | **rc=1，红因＝他腿未跟踪遗留**（非本腿、非 CI 可见） | `build/mut_rvw.go:255:6: undefined: voiceTtsFiles`；`git ls-files build \| wc -l`＝**0** ⇒ 该文件不跟踪，CI 纯净 checkout 结构上看不到它 |
| `go test -count=2` | **未跑（如实记）** | 本腿触碰的 Go 包集合＝**空**（零 `.go` 改动），无"碰过的包"可点名 |

---

## 1. 进场核对（现跑，不抄票面/前任证据）

1. 派单说"pin gofumpt + make staticcheck runnable"——**盘上事实**：这两件事在开工 HEAD 上**已经落地**：
   - gofumpt：`57e764f82`（"85a2/gofumpt: 安装串钉到 @v0.12.0"），证据 `docs/evidence/s1/85-gofumpt-pin.md`；
   - staticcheck：`702c13f51`（ci-r1，2026-10-26），`2025.1.1→2026.2.2`（两次升 pin 的因与实测都写在步骤注释里）。
   ⇒ 本腿把这两件事**当复核对象而不是重做对象**：重做=在同一行上写第二枚 pin，只会制造 diff 噪声。
2. 票 85 还开着、且写面落在 ci.yml 的剩余项＝**AC#4**（上一腿 `85-gofumpt-pin.md` §6-STOP-3 明写"留给下一腿"）。
   ⇒ 本腿的 ci.yml 落刀就落在它上面，一刀 4 增/2 删。

## 2. 落刀（AC#4，只改名字）

- 步名：`"PathResolver junction placeholder (real cases tickets 18/20)"`
  → `"PathResolver junction real case (tickets 18/20; wiring first-green tickets 72/75)"`。
  票面 AC#4 建议含 18/20/72/75：18/20＝真用例归属，72/75＝接线与首次真绿的腿（票面 :26 的翻转注释）。
- 同步改了两处**引用旧名字的注释**（`test-windows` job 头部那句 + 步名下方注释开头），**run: 那一行一个字没动**：
  `bash tools/d22scan/runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows`。
- 零行为改动：没有新增/删除步骤、没有动任何 `if:`、没有 `continue-on-error`、没有放宽任何断言。
  `git diff --numstat` ＝ **4 / 2**（两个 hunk：job 头注释 + 步名行）。

## 3. 门禁读数（本腿现跑，2026-11-06）

| 门 | 命令（逐字） | 读数 |
|---|---|---|
| YAML | `python -c yaml.safe_load(...)` | ok；jobs=`['lint','test-core','test-windows','slo-smoke','slo-full','lint-frontend']`，lint 15 步；test-windows 的 PathResolver 步名读出＝新名 |
| gofumpt | `/d/.../gofumpt.exe -l internal cmd tools scripts frontend` | **空**，rc=0（上一腿脚注③的 CRLF 假阳性文件今日已不出现） |
| staticcheck | `/d/.../staticcheck.exe ./internal/ball/` | **37 条 finding、0 条 `export data version` 行**，rc=1（红在真 finding；样例 `win32_windows.go:257:2: field guidItem is unused (U1000)`） |
| staticcheck | 同物 `./internal/risk/` | **5 条**（见 §4 漂移登记），0 crash line |
| build | `go build ./cmd/... ./internal/... ./tools/...` | **rc=0** |
| vet | `go vet ./...` | rc=1，**唯一红点**＝未跟踪 `build/mut_rvw.go`（tracked-in-build=0，他腿 11-04 遗留，前任证据脚注①同一枚）；**不删**（临时件只建不删） |
| d22scan | `sh scripts/d22scan.sh` | **rc=0**；runtests `PASS=35 FAIL=0 SKIP=0, === RUN=77`；八格 scope 分母 `363/58/116/48/39/116/907/217`，verdict=`clean` |

## 4. 新漂移登记（归票 122 / 85b 清账，本腿不修）

票面与预做都写「`internal/risk` 0 条」——**今天不是了**（这是票面 :43"开工先重测"那句话兑现的一个实例）：
`staticcheck ./internal/risk/` 现报 **5 条**，全部在 `_test.go`、全部来自票 334/394 的新测试文件：

1. `anchor_world_334_test.go:55` const `hostilePOSIX347` unused (U1000)
2. `anchor_world_334_test.go:60` const `danglingTarget347` unused (U1000)
3. `markseam_394_test.go:24:20` string literal contains U+200B (ST1018)
4. `markseam_394_test.go:162:48` 同上 (ST1018)
5. `proccpu_windows_test.go:20:15` `golang.org/x/sys/windows.GetCurrentProcess` deprecated (SA1019)

处置（按票面 :41-49 的规矩）：**逐条判**，未删、未加 `//lint:ignore`、未排除任何包——留给票 122 清账。

## 5. STOP / 本腿没做的格

- **STOP-1 · AC#1/AC#4 的"真 run"一半**：要 push 之后的托管 run 日志行（gofumpt 下载行、改名步骤的 conclusion）。
  本腿不 push（铁律）⇒ 编排者 push 后取。预期：`gofmt (gofumpt)` 步日志首行含 `@v0.12.0`；
  `test-windows` 里改名步骤应继续 `success`（它跑的用例在 run 328/327 已两次真绿，本腿没动 run: 行）。
- **STOP-2 · AC#2（35/34/78 条清算）与 AC#3（变异回验）**：归**票 122 / 85b**（票面 :95、:219 已写死）。
  本腿零 findings 处置。本腿新增的实据：risk 5 条当日值（§4），票 122 清账时应一并纳入分母。
- **STOP-3 · lint job 整体仍红**：裁定①批准的形状（红在读得懂的 finding 上），不是失败。
- **STOP-4 · 票面（`.scratch/wisp/issues/85-*.md`）未动**：不在本腿写面；收件一步按 §3 的读数回填勾框。

## 6. closeable? —— **否（附最短路径）**

票 85 不能在本腿结案：AC#1 后半句（真 run 日志）、AC#2/AC#3（清账＋变异）仍开。
但**本腿这半张已无 ci.yml 遗留**：push 一枚 run 之后 AC#1/AC#4 即可翻；AC#2/AC#3 走票 122。
**结案判据清单（给收件）**：① run 日志里 `go: downloading mvdan.cc/gofumpt v0.12.0`（或安装行原文）；
② 同 run `staticcheck` 步自报行 `findings=<M> toolchain-crash-lines=0`（M≥30，linux 形状）；
③ 同 run `test-windows` 的 `PathResolver junction real case ...` 步 conclusion=success；
④ 票 122 清完账后 staticcheck 步 rc=0（或按裁定留红的具名清单）。

## 7. 撤销口令

一句话复原：把 test-windows 的步名改回
`"PathResolver junction placeholder (real cases tickets 18/20)"`（job 头注释同步），即回到落刀前；
按仓规应 `git revert <本腿 commit>`，禁 `--amend`/`reset`。
