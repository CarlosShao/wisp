# 票 111 — 111-r1b 证据件（交件版）

代号：`111-r1b`（编排者 2026-10-05 派，作为"重派腿"替代 `111-r1`）
交件时刻：`2026-10-05 17:0x +08`
写面（硬边界，照做）：只许改 `scripts/` 下的测试范围表 `scope=(...)`；⛔ 不碰 `.github/workflows/ci.yml`；⛔ 零 `go test/build/vet/env`（仅 `go list` 允许，且不跑专项脚本的 go test 分支）。

---

## §0 起手锚与现读

- HEAD 锚：`9761d630`（`2026-10-05 15:32:16 +08`）。票面 AC#0 锚 = `4e66817`（`2026-09-24`，漂 11 天）。
- `go list ./...`：**本机 GOOS=windows = 35**（`| wc -l`=35）。票面写 **33 ⇒ 漂**（多出者含 `cmd/balldebug` 改名史/`frontend`/`internal/agent/approval`/`internal/llm/adaptertest`/`internal/projctx`/`internal/session` 等；不逐一坐实，只报总数差）。
- GOOS=linux 的 `go list ./...` 在 `cmd/wisp` 拉 `sherpa-onnx-go` 时 `go list` 即退出非零（CGO=0 无解）——票面 AC#4 的 ubuntu 障碍**复现**，非我这枚可修面。
- 5 枚目标包 `_test.go` 枚数（`git ls-files <pkg> | grep -c '_test\.go$'`）：
  `internal/ball`=16（票 11）· `cmd/wisp`=65（票 5）· `internal/perm`=3（票 2）· `internal/plugin`=1（票 1 复现）· `cmd/llmrecord`=1（票 1 复现）。
- ★★ **前提推翻（HEAD 提交真身，用 `git show HEAD:scripts/portable-tests.sh` 权威读，不执行被污染的工作树）**：
  HEAD 的 `core_pin`/`win_pin`/`cli_pin` **已认领全部 5 枚**：
  - `cmd/llmrecord` @ `:141`（core）/`:168`（win）；`internal/ball` @ `:145`/`:169`；`internal/perm` @ `:158`/`:171`；`internal/plugin` @ `:159`/`:172`；`cmd/wisp` @ `cli_pin :178`。
  ⇒ 票面"另有 5 枚带测试的包**零覆盖**"在 HEAD **不成立**。这 5 枚早在 `8fe5c7c`/`7699ec3`（`agent-ticket111`，2026-09-21，两枚 commit 已 `git log -1` 核实真在）就被接入范围表。

## §1 CI 名册真身

- `.github/workflows/ci.yml` 触发点（⛔ 我不改，只引行号）：
  - `:373` `bash scripts/portable-tests.sh --scope=core`（test-core／ubuntu）
  - `:471` `bash scripts/winsec-tests.sh`（windows 腿 winsec 门）
  - `:507` `bash scripts/wisp-cli-tests.sh`（windows 腿 cmd/wisp 门）
  - `:543` `bash scripts/portable-tests.sh --scope=windows`
- `scripts/portable-tests.sh` scope 数组真身（HEAD）：`core` @ `:207-216`（含 `./internal/ball/ ./internal/perm/` `:214`、`./internal/plugin/ ./cmd/llmrecord/` `:215`）；`windows` @ `:220-223`（含同 4 枚 `:222`）；`cli` @ `:230`（`./cmd/wisp/`）；`winsec` @ `:233-239`。
- scope glob 解析集 == pin 复算（仅 `go list`，不跑测试）：core `go list` → **25 枚 = core_pin 25 枚**（diff 空）；windows → **8 = win_pin 8**（diff 空）；cli → `cmd/wisp` = cli_pin；winsec glob `./internal/winsec/...` → `internal/winsec` 单枚 = winsec_pin。⇒ **GUARD C 咬合，无潜在红**。

## §2 五枚逐枚可纳入性

（判据：build tag 原文／靠哪种 runner 才出结论／纳入要不要新依赖／第一发最坏颜色）

1. **internal/ball** — build tag：`//go:build windows`（`hotkey_*`／`sta_release_windows_test.go` 等 6 枚，逐字见 `internal/ball/hotkey_test.go:1` 等）；`//go:build windows && winlive`（`hotkey_live_test.go:1`／`interaction_live_test.go:1`／`live_windows_test.go:1`／`live_guard_windows_test.go:1`／`hotkey_cancel_borrow_live_260_test.go:1` 共 5 枚）；另有 5 枚无 tag（dock/liquid/position/tokens/tokens_table）。runner：core(ubuntu)+windows 腿已认领；winlive 半边**永不跑**（见下第 6 项）。新依赖：无（不需 DLL/APPDATA/凭据）。最坏颜色：windows 腿已 `PASS` 实测过（票 111 AC#2"两腿"）；ubuntu 半边若有平台 bug＝红即发现。**已在册。**
2. **cmd/wisp** — build tag：22 枚 `//go:build windows`、6 枚 `//go:build windows && winlive`（如 `resident_hotkey_live_258_windows_test.go:1`）、1 枚 `//go:build !windows`（`secret_dataroot_119b_test.go:1`）、其余无 tag。runner：仅 `cli` scope，由 `scripts/wisp-cli-tests.sh` 在 **windows 腿**跑（`:60` 处 `GOOS!=windows` 直接 `exit 2` 拒跑）。新依赖：**需 `third_party/sherpa-onnx` 的 DLL**（wisp-cli-tests.sh 从 `deps.toml` 的 `[sherpa-onnx.dll.*]` 现推名单、缺任一即 `exit 1`）。最坏颜色＝票面坑①：**缺 DLL 时加载期 `0xc0000135`、stdout 0 字节、无 `--- FAIL`**＝看着像绿；现有防线＝wisp-cli-tests 的 DLL 预检 + `portable-tests.sh` GUARD A/B + `runtests.sh`（零 PASS+零 FAIL fatal）。**已在册（cli）。**
3. **internal/perm** — build tag：3 枚 `_test.go` **全部无 `//go:build`**（grep 命中 0）。runner：core+windows 均已认领。新依赖：无。最坏颜色：普通红＝发现。**已在册。**
4. **internal/plugin** — build tag：`disposal_test.go` **无 tag**。runner：core+windows。新依赖：无。最坏颜色：普通红。**已在册。**
5. **cmd/llmrecord** — build tag：`main_test.go` **无 tag**。runner：core+windows。新依赖：无。最坏颜色：普通红。**已在册。**
6. **winlive 半边（横切 ball/wisp）** — `ci.yml` `winlive` 命中 **0**（我复跑：`grep -c winlive .github/workflows/ci.yml` = 0；`grep -n` 零行；且 ci.yml 无 `-tags`/`GOFLAGS`/`WISP_LIVE`；`portable-tests.sh`+`runtests.sh` 均不传 `-tags`）。⇒ `windows && winlive` 用例（ball 5 + wisp 6）**默认零编译、CI 永不跑**。纳入它们**必须改 ci.yml 加 `-tags winlive`**（⛔ 我这枚禁面）＋多半要 self-hosted 真机（这些是 live 硬件用例）⇒ **判"要动 ci.yml，本轮不纳入"**。

## §3 实际改动与正控

- **本轮我对 `scripts/` 的改动＝零。** 原因（具名）：
  1. 我被分的 5 枚在 HEAD 已全部认领（§0/§2），无新增可写；写入将是重复行，不改变任何分母。
  2. **写面被占**：`scripts/portable-tests.sh` 工作树里有**另一枚在飞腿 `111-r1` 的未提交编辑**（`git diff --stat` = +106/−11；新增 `GUARD D` 5 处 + 把 `internal/session`/`internal/projctx` 塞进 core/windows 的 scope 数组与 pin）。我对该文件跑 `git commit -- scripts/portable-tests.sh` 会把 **r1 的未提交活以我的名义吞进提交**＝事故 `3f0c4fff` 那一类。**故我绝不动、不提交此文件**（"⛔ 不覆盖／不复用它的读数"）。
- 真实"下一步洞"不是我票面的 5 枚，而是 **session+projctx**（HEAD `git show` 证其未被任何 pin 认领，`go list -f` 证各有编译测试：session TestGo=2、projctx XTest=1）——**而这两枚 `111-r1` 正在补**。我若去补＝与 r1 在同一文件对撞。归口见 §5。
- **正控（读真身后复述，非新造）**：零 PASS 即 fatal 的真身 = `tools/d22scan/runtests.sh`（`:104` "go test exited 0 but reported NO top-level result at all" ⇒ exit 1；`:98` 任何 `--- SKIP` ⇒ exit 1；`:94` 非零 rc 透传；`:75` 强制 `-v -count=1`）。`portable-tests.sh` 在其上加 GUARD A（`:477` 声明在 scope 却 0 编译测试文件 ⇒ 红）／GUARD B（`:643` 每包须打出**自己的** `^(ok|FAIL)\s<转义包名>\s<时长>` 结果行，缺行 ⇒ 红）／GUARD C（`:455` 解析集≠pin ⇒ 红，双向）。cmd/wisp 的"产物没摆出来 vs 包有 bug"由 `wisp-cli-tests.sh:74-99` DLL 预检承担（缺 → `exit 1` 点名缺哪个）。**"这一步哪天一个用例都没跑谁会响亮地报"＝这三级 GUARD + runtests.sh 的零-PASS-fatal。** 本轮无需复制新形状（现有正控已覆盖这 5 枚），亦无新增 scope 行去配新尺。

## §4 壳尺读数（带时刻）

- `sh scripts/d22scan.sh` @ **2026-10-05 17:03:28 → 17:03:53 +08**，**rc=0 clean**。
  正控段（`runtests.sh -C tools/d22scan ./...`）：`PASS=35 FAIL=0 SKIP=0 === RUN=77`。
  实扫（`go run . -root`）：bans #1-5 internal/=228、cmd/=38；ban#6 frontend/=85；ban#7 internal/tools/=23；ban#8 design/=39、internal/=512 Go、cmd/=104 Go；"clean - no D22 ban violations"。
- `sh scripts/check-path-length-budget.sh --with-self-test` @ **17:04:07 → 17:04:11 +08**，**rc=0 VERDICT GREEN**。
  正控 3/3 ok（基准绿→种入超长路径被点名红→移除复绿）；tracked paths=5974、over-budget=57、covered by roster=57、not-in-roster=**0**；longest=180（一枚 `.scratch/wisp/issues/252-...` 名）。
- ⚠ **窗口归属（不自判绿红）**：这两把读数落在 **`111-r1` 正在改 `scripts/portable-tests.sh`（shell 文件）** 的窗口里；两把壳尺都**不执行** `portable-tests.sh`（d22scan 只跑 d22 模块测试+扫 `.go`/文本；path-length 只看 tracked 路径名），r1 的 shell 编辑也不改路径名、不触 emoji（emoji 尺只扫 `.go`），故 rc=0 对 r1 那一处编辑**稳健**。但编排者点名 **`evidence-close-1` 也在跑这两把同尺**——**若与安静窗口不一致，以编排者复跑为准，我不据此判绿**。我只报数与时刻。

## §5 判不动／量不到（具名归口）

- **需动 ci.yml 才能纳入者（本轮硬不纳入）**：winlive 半边（ball 5+wisp 6 枚 `windows&&winlive`）。归 `.github/workflows/ci.yml` 那面（编排者说归 `111-r1`）。
- **写面被占、我不叠加者**：`internal/session` + `internal/projctx` 的真实零覆盖洞——`111-r1` 正以未提交改动补（GUARD D + 纳入 scope/pin）。**我不写同一文件**，避免吞其活。此二者归 `111-r1`；若编排者判应由我这枚收，请先让 r1 落定/判死并释放 `scripts/portable-tests.sh`。
- **步级结论（"那一步真给过结论的 run id + step 号"）**：本机产不出，需 push 后读 CI——归编排者（票 111 `next=` 1–4 全是"要 push"）。⛔ 不拿本地绿冒充 CI 绿。
- **cmd/wisp 的 ubuntu 半边（19 枚红）**：需逐条读码分"产物布局 vs 跨平台真 bug"，且要改 ci.yml 才能在 ubuntu 腿跑 ⇒ 归票 98/新票，非我这枚可写面。
- **幻影/前提核对**：编排者"111-r1 盘上零足迹、无任何 commit"**两处不实**——(a) `f2d069d7`（`2026-10-05 15:48:54`，晚于 15:14 判死）已 commit r1 骨架；(b) r1 此刻在 `scripts/portable-tests.sh` 有未提交编辑＝**仍活着**。我未进 r1 目录、未读其读数、不覆盖不写该路径。

## §6 交件判语与 commit 链

- 交件判语（逐枚）：
  - `internal/ball`／`internal/perm`／`internal/plugin`／`cmd/llmrecord`：**已覆盖（core+windows 两腿，HEAD 提交真身坐实）**，本轮无需我写；winlive 半边默认不跑（归 ci.yml 面）。
  - `cmd/wisp`：**已覆盖（cli scope，windows 腿，带 DLL 预检正控）**；ubuntu 半边不纳入（需 ci.yml + 逐条读码，归票 98）。
  - 结论：票面"5 枚零覆盖"＝过期断言；我这枚**零 `scripts/` 改动**（非未做，而是"无可安全新增 + 写面被占"）。
- commit 链：
  - `b9306d69`（16:50:26）骨架（`-- .scratch/wisp/probes/111/r1b/evidence.md`，暂存清单恰 1 行）。
  - 交件 commit（本枚第二笔，见下条 git log）。
