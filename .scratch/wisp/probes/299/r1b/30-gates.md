# 票 299 接续腿 `299-r1b` · 第 3 笔：门禁四把 ＋ vet ＋ 见证批

时刻 `2026-10-13 09:1x +08`。锚＝HEAD `91c1d8ab`；最终盘面 `cmd/wisp/panel_control_windows.go` sha256＝`ae714aeeafe2f24c584f2cb4eb14a36c3a75f252fc3340cbc969d88087de1e46`（下面每把尺都在这套字节上跑；跑法逐字在 `logs/` 各件头尾）。

| # | 尺（逐字见件内 command_line） | 读数 | rc | 件 |
|---|---|---|---|---|
| g1 | `GOFLAGS= go build ./...` | 无声（0 行输出） | `rc=0` | `logs/gate-build.txt` |
| g2 | `sh scripts/d22scan.sh` | **红→绿**（见 §g2）：终态逐字 `d22scan: clean - no D22 ban violations` | `rc=0` | `logs/gate-d22scan.txt`（终）+ `logs/gate-d22scan-before-fix-red.txt`（红发） |
| g3 | `sh scripts/check-path-length-budget.sh --with-self-test` | `VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`；分母＝`tracked paths=12084 over-budget=57 covered by roster=57 not in roster=0`；自证 3/3 | `rc=0` | `logs/gate-length.txt` |
| g4 | 格式两把尺×两棵树＝8 宫格（gofumpt／gofmt × 工作树／仓外 blob） | 我的三枚件在 8 名册里 **0 命中 ⇒ 净增 0 枚**（见 §g4） | 各 `rc` 见件 | `logs/gate-format-summary.txt` ＋ `logs/gate-format-{blob,work,gofmt-blob,gofmt-work,attrib}.txt` |
| g5 | `tasklist` 跑前／跑后 | 三枚进程（`wisp.exe`／`wisp.test.exe`／`balldebug.exe`）**两把读数全 0** | — | `logs/tasklist-before-after.txt` |
| g6 | `go vet -tags winlive ./cmd/wisp/ ./internal/ball/`（CI 同形） | 干净树**无声**`rc=0`；M3 变异体（把库的 uintptr→Pointer 写回）**命中并逐字点名**`possible misuse of unsafe.Pointer` `rc=1` | 两把各见件 | `logs/vet-clean.txt`、`logs/vet-unsafepointer.txt` |
| g7 | 六钉见证批 `go test ./cmd/wisp/ -count=1 -v -run '^(AC13×2｜AC2｜255ruler3｜33r10｜L1)$'` | **6 PASS / 0 FAIL**，`ok cmd/wisp 2.468s`（第 3 发，跑在提交字节上） | `rc=0` | `logs/run-final-green.txt` |

## §g2 d22scan：红→绿两发（这是本腿抓到的第二枚 r1 残钉）

- **红发**（临时把修复前的拼法放回，跑完立即复原）：`logs/gate-d22scan-before-fix-red.txt` 逐字——

  ```
  scan_test.go:269: repo HEAD violates: cmd/wisp/panel_control_windows.go:100: [phantom-citation] comment cites repo path "internal/w32/w32.go" which does not exist on disk (ticket 212 ban #9): fix the citation or create the file; shorthand forms are a separate prescription, not a violation
  --- FAIL: TestScannerSelfScanOfRealRepoIsGreen (1.59s)
  --- FAIL: TestRealRepoLedgerIsHonest (1.75s)
  runtests.sh: go test exited 1 - packages=[./...] top-level: PASS=33 FAIL=2 SKIP=0
  rc=1
  ```

  成因：`3b806876` 新件的注释把**依赖模块内部**的文件写成 `internal/w32/w32.go` —— 恰好呈本仓 `internal/` 前缀形状（`tools/d22scan/main.go:890` 的 `repoPathRe` 前缀名册），而该路径在本仓盘上不存在 ⇒ ban #9 判 phantom。**r1 的代提信息只写了 `go build rc=0`，这枚 d22scan 红是第一次被量到。**
- **修法（纯注释改拼，⛔ 不动任何语义）**：`:100` 改写成不呈仓内路径 token 的形状——"the library's own internal w32 package carries (that module's `w32/w32.go:55-120`)"（`internal ` 带空格不成 token；`w32/w32.go` 前缀不在名册）。
- **绿发**：`logs/gate-d22scan.txt` 尾逐字 `d22scan: clean - no D22 ban violations`、`rc=0`，分母逐枚在场（`bans #1-5 internal/=231 cmd/=40；ban #8 internal/=533 cmd/=121` 等）。
- **本笔一处自抓（如实记）**：红发复现时用了 python 默认行尾写回，把该件写成 CRLF ⇒ 我自己的 sha256 哨兵当场 `RESTORE-MISMATCH`；随即用 `mut/ctl-pristine-r1b.go.bak`（`ae714aee`）逐字节复原、`grep -c $'\r'`＝0、`git diff` 仍纯注释、gofmt/gofumpt 复净。⛔ 未流入任何提交。

## §g4 格式两把尺（工作树＋仓外 blob，8 宫格）

台面：仓外 blob 树＝`git checkout-index` 把 **1003 枚 tracked .go 的 HEAD blob** 铺到 `%LOCALAPPDATA%/Temp/wisp-299r1b-blobs`（17MB，LF）——⛔ 仓内建树。

| 尺 | blob 树 | 工作树 | 判定 |
|---|---|---|---|
| gofumpt `-l . tools/d22scan tools/mockllm` | 20（其中 cmd/wisp＝**1**：`panel_transport_35r2_test.go`） | 48（cmd/wisp＝2：前枚＋`models.go`） | 工作树独有 28 行＝别腿**未跟踪**样本＋CRLF 假枚（`models.go`、`internal/agent/*`、`internal/risk/provenance.go`、`257/r2/mut/*`、`263/v1`），逐名在 `gate-format-summary.txt` |
| gofumpt over tracked set | 20 | 33 | 同上归因 |
| gofmt 同两形 | 15（cmd/wisp＝**0**） | 41／27（cmd/wisp＝1：`models.go`） | 同族 |

**我的三枚件（`panel_control_windows.go`／`panel_assets_299_windows_test.go`／`panel_host_windows.go`）在 8 个名册里逐枚 0 命中 ⇒"新增 0 枚"**（尺＝集合交，脚本逐名跑，见件尾 "MY three edited files in ANY of the 8 offender sets: NONE"）；三枚件另各过 `gofumpt -l`／`gofmt -l` 单价＝空。
`models.go` 的成因＝295-r1 具名的 CRLF 假枚（工作树独有、blob 无、`git status` 干净）；`panel_transport_35r2_test.go`＝HEAD 自带债（Q-295-r1-gofumpt-two，两树皆有，⛔ 本票写面）。
CI 的 tracked-set 仪器 `attrib.sh --tracked-only` 另跑一发：**`attrib_rc=2`（无判语）**，因 tracked 集里别腿的故意坏样本（`fs_broken.go`、`block_ac13_ac14.go` 等）让 gofumpt 以 123 退出——仪器拒绝在不可读集上给判语；可用的两棵树对拉就是上表。

## §g6 悬空引用（派单点名的两处＋同族一处）

`panel_control_windows.go` 头注释与 `Window()` 注释原文点名 `probes/299/r1/logs/vet-unsafepointer.txt`——**该件不在盘上**。处置：真跑并发件——
- **干净树**：`go vet -tags winlive ./cmd/wisp/ ./internal/ball/`＝无声、`rc=0`（`logs/vet-clean.txt`）。
- **M3 反控**（1 处 `unsafe.Pointer(lparam)` 写回，`-overlay` 驱动）：逐字命中 `.scratch\wisp\probes\299\r1b\mut\m3-panel_control_windows.go.txt:271:22: possible misuse of unsafe.Pointer`、`rc=1`（`logs/vet-unsafepointer.txt`）。
- 两处注释改指 **r1b 的真实件**并按实测改口径（"WHEN THE PORT CARRIES THE WRITE"——库自己的拷贝在依赖体里、本包 vet 不扫依赖体；本仓的显式偏离理由成立）。
- 同族第三处：`panel_assets_299_windows_test.go:34` 点名 `probes/299/r1/30-teeth.md`（亦不在盘上）⇒ 改指本腿 `20-ac2-reverse.md`＋其原始日志（均已落盘）。`grep -rn -- "299/r1/" cmd/wisp/` 现＝0 命中。

## §g7 见证批附件

第三发（提交字节上）逐字件头带 `workdir`／`command_line`／`rc=0`；六枚钉子＝`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`（页面自报 `"1"`）／`TestAC13BringUpSurvivesAReusedThreadQuit`／`TestPanelPageFetchesItsOwnAssets299`（页面自报全量）／`TestTicket255PanelHostBuildsItsWindowOptions`／`TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe`／`TestPanelHostOpensNoListeningSocketL1`。

## 纪律自查

- 全程 ⛔ 起 HTTP listener／⛔ 开端口（反控走的都是 `-overlay` 或临时改件；g5 两把读数全 0）；⛔ 动 `go.mod`／`go.sum`／`.gitignore`／`design/**`／`frontend/**`；⛔ 翻 `AC` 框／碰 `Status:` 行；⛔ 提交 `.png`。
- 临时改件共 **4 处**（AC2 物理发＋M4/M5/M6），每处都当场复原且 **sha256 逐字节对拉**（脚注逐字挂在各 `logs/run-m{2,4,5,6}-*.txt` 尾；其中 AC2 物理发与 M4/M5 复原到锚哈希，M6 复原到该件当刻快照）；一处滑脱（CRLF）被哨兵抓下并复原，见 §g2 末。
