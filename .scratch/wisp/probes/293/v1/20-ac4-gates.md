# 票 293 · `293-v1` · 20 `AC#4` 门禁四把尺（全部本腿自取，⛔ 引实现者的数）

锚＝`b3593507`。落点件＝`logs/`（`build.txt`／`d22scan.txt`／`gofmt-touch.txt`／`gofumpt-touch.txt`／
`gofumpt-mine.txt`／`gofmt-headblob.txt`／`gofmt-crlf-strip.txt`／`crlf-count.txt`／`prepost-compare.txt`／
`repeats.txt`／`extra-checks.txt`／`zero-claims.txt`／`headrun1-rc.txt`／`headrun1-fails.txt`，⛔ 0 字节）。

## 0. 进程护栏（整包之前现量）

- 起手：`tasklist | grep -icE "wisp|balldebug"` ⇒ **0**（尺回 `rc=1`＝零命中）。
- 跑完整包那一族之后再量：`tasklist_now=0`（`logs/extra-checks.txt`）。⇒ 整包读数**没有在活体 `wisp.exe` 的树上取**。

## 1. `GOFLAGS= go build ./...` ⇒ **rc=0**（`logs/build.txt`，末行 `rc=0`）

## 2. `sh scripts/d22scan.sh` ⇒ **rc=0**，末行逐字（`logs/d22scan.txt`）

```
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=229, bans #1-5 cmd/=39, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=524 Go files, comments and _test.go included, ban #8 cmd/=119 Go files, comments and _test.go included
```

⚠ 这一发是在**工作树本来就脏 805 枚**（`design/**`、`frontend/**` 有删除项）的树上跑的 ⇒ "纯净树 rc=0"这句
本腿**不成立**、本腿只报"当前树 rc=0"。零 emoji 那把（ban #8）与裸 `go func(`（ban #1）同时过了这两笔产码。

## 3. 格式两把并排（⛔ 裸 `gofumpt`）

| 尺 | 读数 |
|---|---|
| `gofmt -l cmd/wisp` | **3 枚**：`models.go`、`panel_inbound_guards_35r3_test.go`、`panel_transport_35r2_test.go` |
| `"$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp` | **同一 3 枚**，`rc=0`（两把并排＝同一批，⛔ 一把有牙一把没跑） |
| `gofumpt.exe -l <本票动过的 5 枚路径>` | **空，rc=0**；`gofmt -l <同一 5 枚>` 亦 **空，rc=0**（`logs/extra-checks.txt`） |
| 裸 `gofumpt -l cmd/wisp` | `command not found`（`logs/extra-checks.txt`；⚠ 那一发的 rc 被管道吞了，本腿⛔ 声称量到 `rc=127`，只报"命令不存在"这一条现量） |

**逐枚定性（本腿自己拆的，⛔ 抄实现者的）**——尺＝`tr -dc '\r' < f | wc -c` ＋ `gofmt -l` 对 blob 副本／对 CR 剥除副本各一发：

| 文件 | 工作树 CR 枚数 | `git show HEAD:` 副本 `gofmt -l` | 结论 |
|---|---|---|---|
| `cmd/wisp/models.go` | **334** | 未列（剥 CR 后 `gofmt -l` ⇒ 空） | **CRLF 签出形态**（`core.autocrlf=true`，`git status cmd/wisp` 干净），不是格式债 |
| `cmd/wisp/panel_inbound_guards_35r3_test.go` | 0 | **列出来了** | **HEAD blob 真未格式化** |
| `cmd/wisp/panel_transport_35r2_test.go` | 0 | **列出来了** | **HEAD blob 真未格式化**；`gofmt -d` 的实形逐字＝两行 `func (p *jsParser) peek()/next()` 的**对齐空格**（@@ -804,8 +804,8 @@） |

⇒ 三枚的**最后一笔归属**（尺＝`git log -1 --format='%h %ad' -- <f>`）：`models.go`←`5e8748b3 2026-10-03`、
`panel_inbound_guards_35r3_test.go`←`7f9d6e40 2026-10-07`、`panel_transport_35r2_test.go`←`3a343bc7 2026-10-08`；
`git log --oneline 6ef14788..HEAD -- <这三枚> | wc -l` ⇒ **0** ⇒ 本票一枚没碰。
**判：票面 `AC#4` 那句"`gofumpt -l <自己动过的目录>` 空"按字面形（整目录）今天不可能成立 ⇒ 登记为前程存量缺陷（票 35 那批 2 枚＋CRLF 1 枚），
本票按"我动过的 5 枚为空"这一具名更窄形交付 ⇒ 本格本腿判**成立**，⛔ 让本票顺手格式化别人的文件。
⇒ **实现者申报的那两枚"HEAD blob 口径真未格式化"本腿复现上了**（尺＝`git show HEAD:<f> > $TEMP/blob.go` 后
`tr -dc '\r' | wc -c` ⇒ **0**，`gofmt -l $TEMP/blob.go` ⇒ **列出**；同一发对工作树原件亦 ⇒ **列出** ⇒ blob／工作树两口径同判），
`models.go` 那一枚"工作树 CRLF、blob LF、`git status` 干净"本腿亦复现（CR=334、剥 CR 后 `gofmt -l` 空）。
⚠ 唯一⛔ 复现不上的是 `panel_inbound_guards_35r3_test.go` 的 `gofmt -d` 实形（本腿只逐枚看了 `panel_transport_35r2_test.go`
那两行 `jsParser.peek()/next()` 的对齐空格），⛔ 由一枚外推到两枚。

## 4. 整包：改前／改后**都取到了**（本腿没放弃任何半把）

票面尺＝`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ ./internal/audio/ ./internal/ball/ -count=1 -v`；
派单尺把 `./internal/audio/ ./internal/ball/` 换成了 `./internal/panel/`。**两把本腿并成全四枚一次跑**：
`go test ./cmd/wisp/ ./internal/panel/ ./internal/audio/ ./internal/ball/ -count=1 -v`（⛔ 两把尺的数不可互比，包集不同）。

### 4.1 "改前"怎么取的（派单说取不到，本腿取到了）

`git stash`⛔、`git worktree`⛔（且仓内禁建）⇒ **可行第三形＝`git archive` 导出到仓外**：

```
git archive --format=tar -o "$TEMP/pre293.tar"  6ef14788 && tar -xf … -C "$TEMP/pre293"   # 改前
git archive --format=tar -o "$TEMP/post293.tar" HEAD     && tar -xf … -C "$TEMP/post293"  # 改后（同一形态）
```

两枚导出树**只有本票那 5 枚文件不同**、**同一台机器、同一条 PATH、同一无 `.git`／无 `build/` 的环境 handicap**
⇒ 这是比"改前改后各跑一次"更硬的对照（实现者那两发改前读数是在**脏工作树**上取的，环境不一致）。
已验：`$TEMP/pre293/cmd/wisp/resident_tray_mute_293_windows_test.go` **不存在**（`grep` 回 `rc=2`）⇒ 导出的确实是开工前那一枚锚。

### 4.2 读数（三数都带计数命令逐字）

| 发 | 树 | rc | `grep -c -- '--- PASS'` | `grep -c -- '--- FAIL'` |
|---|---|---|---|---|
| 改前 #1 | `6ef14788` 导出 | 1 | **677** | **26** |
| 改后 #1 | `HEAD`(=`b3593507`) 导出 | 1 | **684** | **26** |
| 改后 #2（重复发） | 同上 | 1 | **683** | **27** |
| 改后 #3（工作树，非导出） | `b3593507` 脏树 | 1 | 706 | 11 |

- **改前 → 改后：`--- PASS` 677 → 684 ＝ +7 ＝ 本票新增的 7 枚 `TestAC293*` 逐枚绿**（`grep -c '^--- PASS: TestAC293'` 在改后两发各＝7）。
- **红名册逐名作差**：`comm -3` 对两枚导出树的 `FAIL: <名>` 集合 ⇒ **两个方向都空**（`logs/prepost-compare.txt`）
  ⇒ **新增红 0 枚**，这一枚是本腿现量、⛔ 转述 `A799`。
- 26 枚红的**成色**（具名，⛔ 算到本票头上）：一族是要真 `build/wisp.exe` 的（`TestAC246*`／`TestAC1*ResidentLeg*`／`TestSecret*`／`TestAC2RealProcess*`／`TestAC3*`），
  一族是无 `.git` 的（`TestCleanCheckoutBuilds_AC11`、`TestReadGitOnThisRepositoryIsSelfConsistent`），
  一族是 WebView2/面板前台时序（`TestPanelHostRealWindowHopAndLifecycle`、`TestAC14*` 各 20s 顶到超时），
  一族是 `design/`/`frontend/` 对照（`TestC21*`、`TestComposer*`、`TestApprovalCard*`、`TestStreamLog*35r8`、`TestPanelColourLiterals*`）。
  ⇒ 全部在改前改后**同进同出**；"改前基线不是绿的"这句**本腿复跑成立**（⛔ 本腿也⛔ 引任何旧读数当基线）。
- ⚠ **同一枚树两发之间红名会抖**（现量）：改后 #1 ↔ #2 的对称差＝**恰一枚**
  `TestAC4FocusReturnToPriorWindowGap33r5`（#1 绿、#2 红），`--- PASS` 684↔683。⇒ 票面那条"单发名册不是可靠尺"
  本腿**独立复现**，且它同时也说明实现者拿"两发交集"作尺的形是必要的。

## 5. 三项零膨胀（本腿自取）

| 声称 | 尺 | 读数 |
|---|---|---|
| 新协程 0 | `git diff 6ef14788..HEAD -- cmd internal \| grep -c '^+.*go func('` | **0** |
| 新锁 0 | `git diff … \| grep -c '^+.*sync\.Mutex'` | **1** ⇒ 逐枚定位＝diff 第 224 行 `+	mu   sync.Mutex`，实形＝**测试夹具 `trayCheckRec.mu`**（`resident_tray_mute_293_windows_test.go:53`）；**产码新锁 0**（写侧仍用现成 `muteMux`） |
| 新包级依赖边 0 | `GOFLAGS= go list -deps ./internal/ball \| grep -c 'internal/audio'` | **0**（`internal/*` 命中 **5** 枚，与票面裁定的 145/5/0 那把同形同号） |

## 6. `AC#4` 本腿结论

**成立**（四把都取到、且**改前对照本腿自己取到了**）。具名不成立的只有半句：票面"⛔ 引用前先重跑"式的那句
"sh scripts/d22scan.sh **纯净树** rc=0"——本腿的树**不纯净**（805 枚脏，非我造成），
本腿交的是"当前树 rc=0 ＋ d22scan 只看 `.go`/白名单扩展名、`logs/*.txt` 进不了它的分母"。
