# 308-r2 · 90-final（交回）

leg＝`308-r2`（窄落地腿，**只改注释行**）· 票＝`.scratch/wisp/issues/308-spelling-invariance-test-panics-when-temp-has-no-volume.md` §4（编排者 2026-10-11 14:2x 点名三句）
起手现跑 `git rev-parse --short HEAD`＝**`202ed4a5`**（＝派单给的值，一致）· 交回前再跑＝**`1da192e9`**（我三笔之后）· ⛔ push

## 1. 三笔

| 笔 | 内容 | 越界尺（`git show --name-only --format=` 逐笔） |
|---|---|---|
| `6f389389` | 锚件（`00-anchor.md`＋起手 logs） | 11 枚全在 `.scratch/wisp/probes/308/r2/**`，射程外 0 行 |
| `4c63d3b2` | 交付：两枚 `internal/risk` 测试件的**注释**（三句） | 2 枚＝`pathresolver_anchor_spelling_windows_test.go`／`pathresolver_volumeless_spelling_308_test.go`，⛔ 其它 |
| `1da192e9` | 证据：`10`／`20`／`30` 件＋logs（含剥注释快照、导出树读数） | 65 枚全在自家 probes 目录，射程外 0 行；新建 `.go`＝**0 枚** |
| （第 4 笔＝本腿收尾） | `90-final.md`＋票尾一条 bullet | 票＝numstat **1/0**（纯追加，⛔ 碰框、⛔ 碰 `Status:` 两行、⛔ 改名 `-done`） |

## 2. 三句落地形状（细节与尺见 `10` 件）

1. `anchor` 文档段末："every other spelling is **unconditional** / **never quietly stops measuring**" → 换成**量过的分母**（卷标齐 9、8.3 命中 10、无卷标输入 7）＋理由改成"a drop is NAMED and LOGGED"。
2. `anchor` ticket-308 注释块："seven … still **run**"／"the **case** ends with a verdict" → 收成"helper 留七枚＋**只有选择那一半**被输入面那枚用例钉住；windows-tagged 那一支⛔ 夹具（哨兵命中 0，三发台面）"。
3. 新件："the **single place** the drive letter is taken" → `Within this test surface`＋具名产码 `pathresolver.go`（守卫在前）与 `internal/tools/volume_windows.go`（读 `VolumeName`、⛔ 切片）＋普查尺留在句里。

## 3. ★我顶回 `308-v1` 建议文本的地方（交编排者复跑再判）

- **`7 (8)` 里那枚 `8` ⛔ 成立**：`getShortPath` 喂无卷标字面量＝`""`（两把台面同值，`logs/probe-denominator*.txt` 第 7 行）⇒ 无卷标那形停在 **7**。新句写 "so it is 7, not 8"。
- **"the only kind Windows hands out here"＝全类断言**：本腿只复现默认＋别名两把台面，⛔ 复跑 v1 的 env 网格 ⇒ 新句收窄成 "the two temp surfaces reachable here"。
- **"printing kept/dropped" 措辞偏松**：代码印的是保留**枚数**＋去掉**名**，且只在真去掉时印 ⇒ 新句按代码逐字。
- ⚠ 具名但⛔ 碰：`junction:157-159` 的 "on the windows CI gear"（v1 §另两处）仍在盘上——那枚文件票 §4 ⛔ 点名、本腿⛔ 动，账归 `Q-308-v1-A`。

## 4. 承重读数一览（全部本腿现跑，件在 `logs/`）

| 尺 | 读数 |
|---|---|
| 剥整行注释后 `cmp`（派发锚 blob ⇄ 交回内容，逐文件） | **IDENTICAL**／**IDENTICAL**（代码 md5 `fb351a9c…`、`c113a6ec…` 两侧同） |
| `git diff -U0` 非注释变化行 | **0 / 0**；numstat `27/6`＋`9/1` |
| 产码两枚 sha256 | `c347e1cf…`／`1dcf0640…`＝票 §1 复跑值同串，blob 与 `git status` 皆未动 |
| 门禁 | build **rc0**／vet **rc0**／长度门 **rc0**（GREEN，57⇄57、worst 228）／d22scan **rc0→rc1→rc1**（红＝`internal/tools/bridge.go` 两枚 ban#9，别腿在飞态；我这面隔离扫描 **rc0 clean**、工作树真扫描 finding 命中 `internal/risk`＝**0**） |
| 格式名册（gofmt⇄gofumpt，工作树⇄导出树，各写射程） | 窄 1⇄1／宽 4⇄5 ⇄ 导出树窄 0⇄0／宽 0⇄**1**；CR 全数过（那 4 枚 CR>0＝假枚；导出树里仍红的 `cmd/wisp/panel_transport_35r2_test.go` CR=0＝HEAD 真脏，⛔ 我账）；**我的两枚＝0 行（两把尺都是）** |
| 整包名册作差（改前仓外导出树 ⇄ 改后工作树） | `RUN 201⇄201`、顶层 `^--- PASS` **132⇄132**、含子测试 `--- PASS` **200⇄200**（子层 68⇄68）、名册枚数 133⇄133、`SKIP 1⇄1` 同名、`FAIL 0⇄0`、`panic: runtime error` 0⇄0 ⇒ **名＋色 `diff` 空**；别名台面另有一发同值（`RISKAT~1` 字样 10 枚证台面真用上） |
| `t.Skip` 代码形 | anchor `0⇄0`、新件 `0⇄0`、junction `1⇄1`（票 18 既有那枚，⛔ 我引入） |
| 票面框数（⛔ 我翻） | `- [x]`＝**6**／`- [ ]`＝**0**（起手与交回同值＝我一行框⛔ 动） |

## 5. 欠读数（⛔ 由我说哪格该翻、⛔ 我改名 `-done`）

- **`Q-308-r2-A`**＝三枚用例在 CI（`ci.yml:487` ubuntu／`:797` windows）的跑机色：⛔ push ⛔ 取得，与 `Q-308-v1-A`（＝吸收 `Q-308-r1-A/B`）同一批，归编排者、与 `A813` 同批取。
- **`Q-308-r2-B`**＝"真无卷标 temp 台面那一发"＝票 §4 已判**永久欠**（本机 `GetTempPath` 必补盘符）。本腿⛔ 复跑 v1 的 env 网格，只在两把我能量到的台面上取证，句里射程已按此收窄。
- **`Q-308-r2-C`**＝交付注释里那组**带日期的计数**（`22 ⇄ 17 test files`、`9／10／7`）会随 `internal/risk` 增减文件或拼法腐烂。census 那枚我已把尺写进句子；`go list` 那枚⛔ 写尺＝缺陷具名，重跑尺＝`GOOS=linux GOFLAGS= go list -f "{{range .TestGoFiles}}{{println .}}{{end}}" ./internal/risk/`。要不要把这两枚计数从注释里撤掉（只留机制＋尺）＝**编排者定**。
- **`Q-308-r2-D`**＝d22scan 在共享工作树交回时仍可能红（309-r2 的 `bridge.go` 在飞）。归属尺已具名（行号漂 887/966→892/971、`202ed4a5` 计数 0、`git status -- internal/tools`＝2 行）；⛔ 算我消门，也⛔ 由我判它绿。
- 沿用而⛔ 新造的账：`Q-308-v1-A..E` 五枚原样在票上，本腿只交句子，⛔ 动判语、⛔ 动框。

## 6. 越格自报（⛔ 自我盖章，等下一枚非实现者腿问）

1. 在**仓外** `mktemp -d` 导出树里新建过一枚 `.go` 探针（`zz_r2_denominator_probe_test.go`，只打印不断言）；仓内 `probes/308/r2/**` 新建 `.go`＝**0 枚**（`find` 计数 0），剥注释快照件名＝`.dispatch`／`.code`／`.final2`（⛔ `.go` 后缀）⇒ 全仓 gofumpt／d22scan 的 `.go` 分母纹⛔ 动。三枚导出树路径逐名留在 `logs/{pre,probe,mine}-tree-path.txt`，只建⛔ 删。
2. 用 `cmd.exe //c "set TMP=…"` 起 8.3 别名台面两发（整包一发＋探针一发）；沿用它人留下的 `C:\Users\swq\tmp\RISKAT~1\RUNNER~1`，⛔ 造⛔ 删。
3. `gofumpt` 不在 PATH，取自 `$(go env GOPATH)/bin/gofumpt.exe`（`v0.12.0 (go1.27.1)`）；`go list`／`go build` 属读数面⛔ 写面，具名自报。
4. 顺手登记（⛔ 删、⛔ 我账）：`git status --porcelain -- .scratch/wisp/probes | grep '\.go$'`＝**3 枚别家未跟踪 `.go`**（`probes/224/v2/posctrl.go`、`probes/33/r5/block_ac13_ac14.go`、`probes/33/r5/logs/prev-test-copy.go`）——任何把 `.scratch` 扫进来的宽域格式尺都会先撞上它们。
5. ⛔ 跑 `cmd/wisp` 测试包（sherpa DLL 前置归编排者）；⛔ 动 junction 那枚文件的任何字；⛔ `push`／⛔ `add -A`／⛔ `--amend`／⛔ `reset`／⛔ `rebase`／⛔ `stash`／⛔ `clean`；⛔ 在仓内建 worktree。

## 7. 我⛔ 做的事（口径边界）

⛔ 翻框、⛔ 碰 `Status:`、⛔ 改票名 `-done`、⛔ 宣布任何格"成立"、⛔ 动产码／helper 逻辑／用例逻辑、⛔ 加删用例、⛔ 引 `t.Skip`、⛔ 放宽任何断言、⛔ 造 `--- SKIP`。本腿的全部交付＝**三处注释**＋证据＋票尾一条。

## 8. 收尾读数（★⛔ 覆盖 §0/§1 那行的"交回＝`1da192e9`"，那枚是第三笔之时）

交回 HEAD＝**`c30b7abe`**（第四笔＝本件＋票尾一条；`git log --oneline` 里我这四笔＝`6f389389`→`4c63d3b2`→`1da192e9`→`c30b7abe`，起手＝`202ed4a5`）。收尾又现跑一遍（件＝`logs/final-readings.txt`、`logs/gate-d22scan-run4.txt`＋`.rc.txt`、`logs/closing-note.md`）：

- `git status -- internal/risk`＝**0 行**（我的两枚已入库）；`git status -- internal/tools`＝**2 行**（309-r2 仍在飞）。
- 剥整行注释后 `cmp` ⇄ `202ed4a5` blob：两枚**各 IDENTICAL**（md5 `fb351a9c…`／`c113a6ec…` 与 §4 同值）；产码两枚 sha256 仍＝`c347e1cf…`／`1dcf0640…`。
- 门禁第四发＝build **rc0**、vet **rc0**、长度门 **rc0**（`VERDICT GREEN`）、`sh scripts/d22scan.sh` **rc1**：两枚 finding 仍在 `internal/tools/bridge.go`，行号第三⇄第四发之间又从 `892/971` 漂到 **`892/975`** ⇒ "别腿在飞态"这一判读第四次复现；命中 `internal/risk` 的 finding＝**0**。⛔ 由我消、⛔ 由我盖（归编排者复跑认账＝`Q-308-r2-D`）。
- 票框现跑＝`- [x]` **6**／`- [ ]` **0**（＝起手值，⛔ 我翻）。

