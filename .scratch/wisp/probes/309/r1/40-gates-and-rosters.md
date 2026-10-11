# 309-r1 · `AC#4` 门禁与越界（每枚门禁件自落 `rc=N`）

## 三枚编译/门禁 rc=0（原始＝`logs/gates.txt`；HEAD 现量 `1ad48333`，非我基线——见 90-final 台面节）
- `GOFLAGS= go build ./...` → **rc=0**
- `GOFLAGS= go vet ./internal/tools/` → **rc=0**
- `sh scripts/d22scan.sh` → **rc=0**（clean；`internal/=229`/`cmd/=39`/`internal/tools/=23` production Go、`ban#8 internal/=531` 含注释与 `_test.go`；⛔ 放宽任何断言、⛔ 改 allowlist.txt）
- `go test ./internal/tools/`（全量）→ **rc=0**（`logs/after-fullquiet.log`）

## d22scan 曾一度 rc=1（本轮抓到的真违规，非仪器噪声）
`toolrequest_nail_309_test.go:103` 的 `t.Fatalf` **字符串字面量**里放了 `≈`(U+2248，落 `U+2200–U+22FF` 禁带) ⇒ `TestScannerSelfScanOfRealRepoIsGreen`＋`TestRealRepoLedgerIsHonest` 双双红（两枚都现读真仓）。注释里的 `⛔`(U+26D4)/`⇒`(U+21D2) 按 Q-46(c) **注释豁免**、未被扫到；`①`/`②`(U+2460 圈数字) 是刻意留的空隙、⛔ 扫。**修法＝把字符串里的 `≈零` 改 `接近零`**（⛔ 改任何断言逻辑）⇒ 重跑 rc=0。⛔ 动 `tools/d22scan/allowlist.txt` 一字节。

## 格式名册（两把尺并排，⛔ 顺手修、⛔ 因既有债判红）
| ruler | `internal/tools`（我射程） | `cmd/wisp`（既有债，⛔ 我动） |
|---|---|---|
| `gofmt -l` 工作树 | **0** 枚 | 0 枚 |
| `gofmt -l` HEAD-blob 导出树 | **0** 枚（基线） | 0 枚 |
| `gofumpt -l` 工作树 | **0** 枚 | **2**＝`models.go`＋`panel_transport_35r2_test.go`（编排者面的另立与否，⛔ 我修） |
| `gofumpt -l` HEAD 导出树 | **0** 枚 | 35r2 残留（本票无关） |
⇒ 我的改动在 `internal/tools` 新增 **0** 枚格式问题（射程＝`internal/tools/**` 的 `gofmt`/`gofumpt`，两树皆 0↔0）。

### 本轮踩到的 eol 坑（具名，非既有债）
`bridge.go` 经我编辑后工作树变成 **CRLF**（`CR=1313`），而 `task.go`/`subagent_197.go` 等未触文件与 HEAD blob 均 LF（`.gitattributes: *.go text eol=lf`）⇒ `gofmt -l` 唯一红＝`bridge.go`。这是我引入的真回归、⛔ 当既有债放着 ⇒ `perl -pi -e 's/\r\n/\n/g'` 归一到 LF（纯 eol、⛔ 改内容）⇒ 工作树 `gofmt -l internal/tools`=0、`gofmt -d`=空、`git diff --numstat` 仍 `32 3`。两枚新 `_test.go` 本就是 LF、未犯。

## 逐名作差（改前⇄改后，同层级尺，⛔ 混着数）
`go test -v` 名册（`logs/baseline-2.log` ⇄ `logs/after-1.log`；名集合差＝`baseline-2-names.txt` ⇄ `after-1-names.txt`）：
- 顶层 `^--- PASS`：**209 → 211**（+2＝我两枚顶层用例，无子测试）。
- 子测试 `^    --- PASS`：**75 → 75**（未动）。
- `--- FAIL` 顶层/子：**0 → 0**；`--- SKIP`：**0 → 0**（⛔ 造 SKIP 换绿）。
- 名集合差＝**纯新增**（只两条 `> Test309…`），⛔ 删除线 ⇒ **无覆盖搬家式消失**。
- **两枚真喂空 id 进桥的 236 载体照绿**（Ⓑ 行为不变必须成立）：
  `Test236R2TaskSpawnRefusesWhenHostGaveNoTaskID`（含 `:288 TaskID: "", CorrelationID: ""` 直喂 `Execute`）＝**PASS**；
  `Test236R2TaskCancelRefusesWhenHostGaveNoCallerID`（经 `:118 x.cancel(t, "", target)` → `task_cancel_221_legs_test.go:116-119 TaskID: callerID`）＝**PASS**。

## Git 纪律（每笔显式 pathspec）
- commit1 `1cd57ec1`＝锚件（`-- .scratch/wisp/probes/309/r1/00-anchor.md`）。
- commit2（`feat(309-r1)`）＝`bridge.go`＋两枚 `_test.go`；`git diff --cached --name-only` 数＝**3**、越写面＝**0**（`logs/staged-code.txt`）。
- commit3＝本证据件＋`logs/`（显式逐个 pathspec，⛔ `git add .scratch/wisp/probes/309/r1` 整目录吞入临时导出树；该导出树已 `rm`，读数已落文件）。
- ⛔ `push`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`／`add -A`／`add .`；`.gitignore`/`go.mod` 未提交改动非我名下、⛔ 动。`AC#4` 闭合。
