# 274 · AC#0 两枚数量现量清账（只读腿；2026-11-06）

**腿性质**：只读/分析。零产码、零 CI yaml 改动、零工单改动、零推送。
**锚点**：仓外共享 clone `D:\tmp\wisp274-ac0`（`git clone --shared`）检出 **`86fb145eed4ee8ce0912c5468d3208fc9fe3bf8d`**（= dispatch 所指 `bca0990f8` 之后 3 笔的 dev 尖；**关键路径现量等价**：
`git diff --stat bca0990f8 86fb145ee -- .github/workflows/ci.yml scripts/build.ps1 frontend docs/specs/SPEC-11` ＝ **空输出**，故本文件全部读数在两个锚点上同值）。
**CI 侧佐证 run**：**`42205868481`**（2026-11-06T10:51Z，sha `ab808c35`；该 sha 与 `86fb145ee` 在同一组构建胶水路径上的 diff 同样为空 ⇒ 跑的就是本锚点的 ci.yml＋build.ps1）。本日后续 run `42209230080`（`86fb145ee` 本尊）尚在跑，故取最近完成发。

> 判据定义（AC#0 两枚数的口径）：(a) CI 的 `lint-frontend` 一发真跑成功后，`npm run build` 产出几枚文件、落在哪；(b) 同一发 run 里 `build.ps1` 造出的 exe，`wisp.exe panel-assets` 回 `built=true` 还是 `false`、枚数差几枚。下文全部数字为本腿现量，不引自任何既有票面/台账。

---

## §1 数 (a)：CI `lint-frontend` 的 `npm run build` 产出

**读数：3 枚文件，落在 `frontend/dist/`（lint-frontend 作业自己的工作目录里，job 结束即丢）。**

| # | 产物 | 大小（vite 报 / 字节） |
|---|---|---|
| 1 | `dist/index.html` | 1.04 kB ＝ 1,068 B |
| 2 | `dist/assets/index-Dxt2y5oW.css` | 65.08 kB ＝ 65,088 B |
| 3 | `dist/assets/index-g8i3FA7_.js` | 625.00 kB |

出处＝run `42205868481` 的 `lint-frontend` 步 `build (vite build -> frontend/dist, the bytes go:embed carries)` 日志（2026-11-06T10:52:10Z，vite v8.3.0，2452 modules）。

**结构尺（HEAD `86fb145ee` 现量，ci.yml 全文 1,378 行）**：

```bash
cd /d/tmp/wisp274-ac0
grep -n 'npm run build' .github/workflows/ci.yml        # 命中 1 枚 = :1321，在 lint-frontend（job 起 :1278）内
grep -n 'setup-node'     .github/workflows/ci.yml        # 命中 1 枚 = :1292，同在 lint-frontend
grep -c '^  needs:'      .github/workflows/ci.yml        # 0   （跨 job 依赖边为零；另一把尺 grep -c 'needs' = 9，全是注释/散文）
grep -n 'uses: actions/upload-artifact' .github/workflows/ci.yml   # 仅 :1158 / :1233 两枚，均为 SLO JSON 报告，无 dist、无 exe
```

⇒ 数 (a) 的完整答案：**产出 3 枚，落 `frontend/dist`，且没有任何 upload/needs 边把它送出去——job 结束即弃**。它对出货 exe 的唯一历史作用是"证明这台 Linux 机器上有 node"；现在连这个作用都被 (b) 取代（见 §2）。

## §2 数 (b)：`build.ps1` 造出的 exe 带不带页面

**读数（CI，run `42205868481` 三枚产 exe job 全部现量）：`built=true` 形；dist 名册 4 枚 vs 干净检出形 1 枚 ⇒ 差 3 枚页面产物。**

CI 侧 `scripts/build.ps1` 调用＝**3 处**：

```bash
grep -n 'scripts/build.ps1' .github/workflows/ci.yml
# :886  test-windows   "cgo build smoke (build.ps1 fetch-deps + mingw link + doctor)"
# :1151 slo-smoke     "Build wisp.exe"
# :1227 slo-full      "Build wisp.exe (deps cached on the runner)"
sed -n '812,1122p;1123,1191p;1192,1277p' .github/workflows/ci.yml | grep -cE 'node|npm'   # 三枚 job 自身步内 = 0 枚 node/npm 步
```

三枚 job 自身零 node/npm 步，但 `build.ps1` 第 2 步（`scripts/build.ps1:86` 起，全文 370 行）**自己真跑 npm**。run `42205868481` 实测（test-windows 与 slo-smoke 两发日志逐行在案；slo-full 以 job conclusion=success + 同一脚本同一闸门结构佐证，其日志本会话 gh 拉取反复 EOF，未逐行取回——如实登记）：

```
build.ps1: frontend: node = C:\Program Files\nodejs\node.exe (v22.23.4); npm = C:\Program Files\nodejs\npm.cmd (10.9.9)
build.ps1: frontend: node_modules ABSENT -> running npm ci
build.ps1: frontend: pre-build dist snapshot = 1 file(s) [.gitkeep]
build.ps1: frontend ok: 4 file(s) in frontend/dist (entry index.html is 1068 bytes):
           .gitkeep=0 assets\index-BPq-izHv.js=625045 assets\index-Dxt2y5oW.css=65088 index.html=1068
           || provenance: this run's npm added new artifact name(s) [assets\index-BPq-izHv.js assets\index-Dxt2y5oW.css index.html]
build.ps1: go build ok (cgo linked against sherpa-onnx C API)
```

⇒ CI 托管 Windows runner 的作业 PATH **现量有 node v22.23.4 / npm 10.9.9**（不再是"取不到"）；出处闸门（pre-snapshot 对拉）绿；exe 出炉。

**读数（仓外共享 clone 本机复现，同锚 `86fb145ee`，可逐字重跑）**：

```bash
git clone --shared "D:\work\workspace\projects plans\Wisp" "D:\tmp\wisp274-ac0"
cd /d/tmp/wisp274-ac0 && git checkout -q 86fb145ee
# 原生依赖：copy 主树 third_party/sherpa-onnx/{*.dll,.cache-manifest.json}（三 DLL sha256 前缀
#   7f66f939…/300e0c88…/02ca4606… 与 deps.toml 钉值逐枚相等 ⇒ fetch-deps 走 cache-hit 快路）
cp "/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx/"*.dll third_party/sherpa-onnx/
cp "/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx/.cache-manifest.json" third_party/sherpa-onnx/
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/build.ps1 -Env dev     # rc=0（本机 node v24.21.0 / npm 11.19.0）
./build/wisp.exe panel-assets
```

本机终读数：

```
panel assets embedded: 4 files, entry=index.html built=true      rc=0
panel assets check: entry=index.html built=true 2 asset refs resolve [./assets/index-BPq-izHv.js ./assets/index-Dxt2y5oW.css]
```

**exe 内嵌物 ↔ 本趟 dist 逐枚哈希对拉（sha256 前 16 位，本机现量）**：

| 文件 | exe `-manifest` 报 | 盘上 dist 实测 | 相等 |
|---|---|---|---|
| `.gitkeep` | `e3b0c44298fc1c14` (0 B) | `e3b0c44298fc1c14` (0 B) | ✔ |
| `assets/index-BPq-izHv.js` | `5ac41dcbdfac6f03` (625,045 B) | `5ac41dcbdfac6f03` (625,045 B) | ✔ |
| `assets/index-Dxt2y5oW.css` | `37ad4b9194c89e2f` (65,088 B) | `37ad4b9194c89e2f` (65,088 B) | ✔ |
| `index.html` | `c9d04324af553e8b` (1,068 B) | `c9d04324af553e8b` (1,068 B) | ✔ |

⇒ 数 (b) 的完整答案：**`built=true`、内嵌 4 枚（含跟踪锚 `.gitkeep`），对干净检出形（1 枚）差 3 枚页面产物；内嵌字节＝本趟 dist 字节（同上下文逐枚哈希相等）**。本机 Windows fresh build 与 CI Windows build.ps1 的名册**逐枚同名同大小**（`index-BPq-izHv.js`/`index-Dxt2y5oW.css`/`index.html 1068`）。

## §3 摆放计数（placement，HEAD `86fb145ee` 现量）

| 尺 | 读数 | 位置 |
|---|---|---|
| `npm run build`（ci.yml） | **1** | `:1321`（lint-frontend 内） |
| `scripts/build.ps1`（ci.yml 真 `run:`） | **3** | `:886` test-windows · `:1151` slo-smoke · `:1227` slo-full |
| 三枚产 exe job 内自带 node/npm 步 | **0** | ——（npm 由 build.ps1 第 2 步自带） |
| `setup-node` | **1** | `:1292`（仅 lint-frontend） |
| `needs:` 依赖边 | **0**（另一把尺 `needs` 子串 9，均为注释） | —— |
| `upload-artifact` 真步 | **2** | `:1158`/:1233，均 SLO JSON |
| `panel-assets` 在 ci.yml | **1 处且是注释**（`:1330`，权限串示例；CI 从不对 exe 跑它） | —— |
| `cache: npm` | **1** | `:1297`（仅 lint-frontend；三枚 exe job 每次冷 `npm ci`） |
| build.ps1 里的旧谎话（`no frontend yet` / `lands in S5` / `skipped (no frontend`） | **0 命中** | 第 2 步头 `:86-88` 现为「This step BUILDS the page bundle the exe is supposed to carry」 |
| build.ps1 第 2 步闸门 | 10 枚具名 `Fail` 支 + 唯一 `exit 1`（`:49`）+ 出处快照对拉（pre-snapshot → `:212` 红因） | `scripts/build.ps1:86-215` |

## §4 时长读数（run `42205868481` 步级现量）

| 步 | 读数 |
|---|---|
| lint-frontend 整 job | 31 s（10:51:45→10:52:16）；其中 `npm run build` 步 ≈6 s |
| slo-smoke `Build wisp.exe` | frontend 段 ≈19 s（npm ci ≈12 s + vite ≈7 s）＋ go build ≈21 s |
| test-windows cgo build smoke | npm ci ≈10 s（10:53:30→10:53:40）＋ vite ≈5 s，绿 |

## §5 裁决/建议表（其余部分要动什么；只改注释还是真接线）

| # | 项 | 现量事实（本腿） | 判 | 要动什么 | 类 |
|---|---|---|---|---|---|
| 1 | build.ps1 第 2 步是否真产页面 | `:86-215` 真跑 npm＋5 道具名闸门＋出处快照对拉；CI 3 job、本机 clone 双侧实测 rc=0、exe `built=true`/差 3 枚 | **已闭（实测，非推断）** | 无 | 真接线——**已在 HEAD 落地并测绿** |
| 2 | 那句两半都假的过期注释 | 旧字符串全仓 `build.ps1` 内 **0 命中**；`.DESCRIPTION`/步头现在说实话 | 已闭 | 无 | 注释——**已落地** |
| 3 | CI 产 exe 三 job 摸不摸得到页面 | 三 job 经 build.ps1 第 2 步必然产页面且被出处闸门咬住；本发 run 三 job 全绿 | 已闭 | 无 | 真接线——已闭 |
| 4 | lint-frontend 那 3 枚产物被弃 | `needs:0`＋无 dist artifact ⇒ job 结束即丢 | **良性冗余**：exe 侧自产自校（闸门咬出处），传 artifact 不买来任何 exe 正确性 | 不建议加跨 job artifact 传递（徒增依赖边与时长） | 结构备注 |
| 5 | CI 从不对出炉 exe 跑 `panel-assets` | ci.yml 仅 `:1330` 一处注释提及 | 弱残留：上游闸门＋`go:embed` 确定性已使 `built=true` 是"构造保证"；可选加固＝slo-smoke build 后补一发 `./build/wisp.exe panel-assets -check`（秒级） | **可选**，非必需 | 保险丝（锦上添花） |
| 6 | `SPEC-11 §2.2` 文字 | 仍写 `--with-frontend`/`assets/web`/"存在则跳过"，与已落地的无条件 npm 步**反向矛盾**（同树重跑被出处闸门判红是设计行为） | **唯一剩余的文字级欠账**；且"存在则跳过"那半句如今与实现相反 | 改 `docs/specs/**` 文字＝**人工批准**（契约面，本腿零触碰） | 契约文字——待人批 |
| 7 | 托管 runner PATH 上的 node | **现量有**（v22.23.4/npm 10.9.9）；但仅 lint-frontend 有 `setup-node` 钉子，三 exe job 依赖镜像自带 | 观察项：未来 runner 镜像若撤预装 node，三 job 会**按设计红**（具名红因 `:119`/`:123`，不是静默跳过） | 届时再在三 job 补 `setup-node`；现在不动 | 观察项 |
| 8 | exe job 冷 `npm ci`（无 npm cache） | ≈10-12 s/发 | 可接受 | 可选加 cache，非必需 | 优化项（低利害） |
| 9 | 同一树第二次 `build.ps1` 必红（同名册被出处闸门判红） | 设计行为（红因文本自带处置：清 dist 重跑）；CI checkout 首发不受影响 | 接受 | 无 | 已接受代价 |

## §6 总判

**既不是"只改注释"，也不再是开着的接线缺口**：出货 exe 的页面接线在 HEAD 已真实落地（build.ps1 第 2 步＋具名闸门），且本腿用两条独立证据链（CI run `42205868481` 三 job + 仓外干净共享 clone 逐哈希复现）现量坐实 `built=true`/差 3 枚。剩余待办只有一枚**契约文字项**（表 #6，人工批准）与三枚低利害可选项（#5/#7/#8，可不做）。

## §7 纪律自证

- 零产码：本腿写面仅本文件；`cmd/ internal/ scripts/ tools/ .github/ frontend/ design/` 零写入；工单文件零触碰；零 push。
- 环境事实（如实）：本机 `pwsh 5.1` 因 PSModulePath 被 PS7 污染缺 `Get-FileHash`，故本机 build 用 `pwsh 7` 跑（CI 用 `powershell` 5.1 跑，run 实测绿——差别是本会话环境，不是脚本缺陷）；clone 内原生依赖来自主树已验文件（deps.toml 钉值逐枚对上，fetch-deps 报 `cache hit ... matches deps.toml (sherpa-onnx 1.13.8)`）；`CGO_ENABLED` 未置 0（build.ps1 自置 1）。
- 未逐行取回项：slo-full 的 build.ps1 日志行（gh API 本会话反复 EOF）；以 job conclusion=success＋同一脚本同一闸门结构佐证，已如实登记。
- clone 为一次性台件（`D:\tmp\wisp274-ac0`），构建产物只落在 clone 内，主树 `frontend/dist`、`build/` 未被本腿碰过。
