# 票 329 收口复核 — AC#0 交件三把尺 ＋ AC#7 cmd/wisp 整包终态四数（329-closure，只跑尺＋记账）

- 腿名＝`329-closure`（收口/复核腿，⛔ 产码零改动：`frontend/**`／`cmd/wisp/**` 产码／`internal/**` 产码／`PLAN.md`/specs 全程未碰；写面只有本件＋票 329 票尾一条 bullet）。
- 台面＝共享工作树 `D:/work/workspace/projects plans/Wisp` 起手锚 `6e2eb5845`（HEAD，2026-11-08 11:31 现跑）；整包/门禁全在 `git clone --shared` 导出树 `/d/tmp/wisp/20261108/329-closure-tree`（同一锚）跑；环境按规矩 `unset GOTMPDIR TMPDIR TEMP TMP`、PATH 前插 `third_party/sherpa-onnx:build`、⛔ `CGO_ENABLED=0` 未设、显式 `-timeout 30m`。
- 九格现量交件时复量＝AC#1–#6、#8 已勾未动，本腿⛔ 翻任何 AC 框、⛔ push。

## 1. AC#0 尺一：提交真在（逐枚 `git show --name-only`，2026-11-08 本腿现跑）

**Go 半边名册（落地腿 329-r1 七笔）**——逐枚核对：每笔只落 `cmd/wisp/**` ＋ `.scratch/wisp/probes/329/r1/**`，issues/evidence 逐枚 0，**全部合格**：

| 提交 | 内容 | 文件面 | 判 |
|---|---|---|---|
| `60e82f11c` | 落地 1/5 交付腿本体＋snapHost 字段 | `cmd/wisp/panel_snapshot_deliver_329.go`、`cmd/wisp/run.go` | ✓ |
| `bed1c348e` | 落地 2/5 publish 挂推跳 | `cmd/wisp/panel_pump.go`、`cmd/wisp/run.go` | ✓ |
| `a0a2f4409` | 落地 3/5 evalPanelJS 宿主半边 | `cmd/wisp/panel_host_windows.go` | ✓ |
| `8e1ea71d5` | 落地 4/5 装配根挂推跳 | `cmd/wisp/resident_windows.go` | ✓ |
| `2c4f829a7` | 落地 5/5 AC#3 逐位相等尺 | `cmd/wisp/panel_snapshot_deliver_329_test.go` | ✓ |
| `1d1261e04` | 追加更正（头注路径截断） | `cmd/wisp/panel_snapshot_deliver_329.go` | ✓ |
| `a0d6ec85a` | 交件探针面一（page-contract＋三把尺读数） | `probes/329/r1/logs/**`（13 枚）＋`probes/329/r1/page-contract.md`＋`probes/329/r1/txt/**`（4 枚） | ✓ |

**同票非落地腿四笔（另账登记，非 329-r1 交件名册，本腿如实列明）**：

| 提交 | 内容 | 文件面 | 判 |
|---|---|---|---|
| `7535d83b4` | 编排者死腿遗产代提（r1 票尾 bullet＋r1 证据件代提） | `probes/329/r1/**`（12 枚）＋`docs/evidence/s1/329-snapshot-delivery-r1.md`（一枚短名件）＋票尾 bullet＋`probes/orch/20261017/txt/**`（3 枚，代提人自己的派单存档射程） | cmd/wisp 0、产码 0；`probes/orch/**` 不在 AC#0 四槽名册内＝代提腿自有射程，如实记 |
| `c27d5909e`／`4b4243686`／`249786fd5` | 裁腿 329-v1 探针面一/二/三 | 仅 `.scratch/wisp/probes/329/v1/**` | ✓ 零产码；⚠ 落 `probes/329/v1/**` 不在 AC#0 字面四槽（那格为 r1 写）＝裁决腿自有派单射程，如实记 |
| `7ccbdae44` | 裁腿 329-v1 判语件 | `probes/329/v1/**`（3 枚）＋`docs/evidence/s1/329-v1-verdict.md`（一枚短名件） | ✓ |
| `bf9ef6329` | 裁腿 329-v1 票尾一条只追加 bullet | 票 329（0 删 1 增）＋`probes/329/v1/txt/msg-5-bullet.txt` | ✓ |

**页面半边 `33b264d70`（按派单指令纳入核对、另账登记）**：文件面＝`frontend/package.json`、`frontend/scripts/render-snapshot-receiver.tsx`、`frontend/src/lib/panel.ts`、`frontend/src/main.tsx`、`frontend/src/panel-root.tsx`＋`docs/evidence/s1/329-ac5-page-receiver.md`＋票 329 票尾一条 bullet。**它不进 Go 半边 AC#0 名册**：AC#0 票面四槽只许 `cmd/wisp/**`／`probes/329/r1/**`／evidence 一枚短名件／票尾一条 bullet，而票面 `:42` 对 329-r1 明写"⛔ 动 `frontend/**` 一字"——页面写面是后来经台账 `A1622` 一线授权、以独立页面腿（`329-ac5-page-receiver`）落地的，AC#5 页面那半的凭据在它自己的证据件与票尾 11-08 bullet 里，与 Go 半边名册分两本账。

## 2. AC#0 尺二：件真在（wc -c／行数／CR／占位符／0 字节，本腿现跑）

| 件 | 字节 | 行 | CR | 占位符 |
|---|---|---|---|---|
| `docs/evidence/s1/329-snapshot-delivery-r1.md` | 11184 | 100 | 0 | **`待填` 3 处（:68/:69/:83）** |
| `docs/evidence/s1/329-v1-verdict.md` | 17869 | 119 | 0 | 0 |
| `docs/evidence/s1/329-ac5-page-receiver.md` | 7092 | 64 | 0 | 0 |
| `.scratch/wisp/probes/329/r1/page-contract.md` | 4434 | 65 | 0 | 0 |

- 0 字节件：`find .scratch/wisp/probes/329 -type f -empty`＝**0 枚**（导出树同读数）。
- ★如实：r1 证据件那 3 处 `待填` 就是 329-v1 判语记的两处缺陷（§表 :68/:69 改前/改后 skip3 四数＋:83 桌面卫生尺），至今无人回填＝缺陷在件不在字节，本腿写面不含该件，⛔ 代填；:92 那行是编排者代提注里的登记文字，不算占位符本体。第 4 节把这两行四数从盘上日志复量出来备查。
- 桌面卫生尺（:83 那格）：v1 腿已在自己的判语件与票尾 bullet 现跑补过（PID 21892 活着）；本腿未复跑该 PID（进程属 10 月时点，如实注明）。

## 3. AC#0 尺三：自己复跑（起手三核，共享工作树现跑）

- `grep -rn lastPanelSnapshot --include='*.go' internal cmd | grep -v _test.go`＝5 行：定义 `panel_pump.go:489`＋注释 3 行＋**真读者 1 枚 `panel_snapshot_deliver_329.go:95`**（`_, data, seen := rt.lastPanelSnapshot()`，AC#2 尺继续成立；行号较票面 `:417` 已漂＝panel_pump.go 后续腿所致，如实记）。
- 入向名册钉：`grep -cE 'Method[A-Za-z]+ +=' internal/panel/bridge.go`＝**7**（未撞）。
- 页面接收器：`grep -rn wispOnPanelSnapshot frontend/src frontend/scripts`＝**21 命中**（ac5 腿交件时报 11，其后前端腿增量所致，非零即合同已有人接）。
- `git merge-base --is-ancestor 33b264d70 HEAD`＝真（页面半边在现 HEAD 历史里）。

## 4. r1 证据件 :68/:69 两行四数的盘上日志复量（备查，⛔ 不回填原件）

对 `.scratch/wisp/probes/329/r1/logs/` 六枚整包日志逐枚 `grep -c`（尺＝`^=== RUN`／`^--- PASS`／`^--- FAIL`／`^--- SKIP`）：

| 日志 | RUN | PASS | FAIL | SKIP | ok | 尾停 |
|---|---|---|---|---|---|---|
| `full-before-anchor-feab45e3.log` | 166 | 115 | 2 | 0 | 0 | `TestAC4FocusReturnToPriorWindowGap33r5` |
| `full-before-anchor-feab45e3-run2.log` | 203 | 148 | 0 | 0 | 0 | `TestAC13BringUpSurvivesAReusedThreadQuit` |
| `full-before-anchor-feab45e3-skipAC13.log` | 203 | 148 | 0 | 0 | 0 | `TestAC13BringUpRefusesAThreadWithAQueuedClose` |
| `full-before-anchor-feab45e3-skipAC13pair.log` | 205 | 149 | 1 | 0 | 0 | `TestAC14GoSideEvalPushReachesThePage`（红＝`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card`） |
| `full-before-anchor-feab45e3-skipset2.log` | 203 | 147 | 1 | 0 | 0 | `TestAC13NamedRefusalEndsThePanelThreadInsteadOfRefusingForever`（红＝`TestRunPacketLightsTheRosterRowOfItsLiveL1Window`） |
| `full-after-tip-1d1261e0-skipset2.log` | 206 | 151 | 0 | 0 | 0 | `TestBallPanelGesturesReachThePanelThread` |

skip3 那两行（:68 A／:69 B）与上表哪两枚一一对应**盘上判不出**（没有具名 B 侧 skipAC13 日志）⇒ 本腿只登读数、⛔ 替原件断言对应关系。六发共同形状＝无一跑到 `ok`，全部死在 AC13/AC14 真窗区（与 329-v1 判语一致）。

## 5. AC#7：cmd/wisp 整包终态四数（导出树，DLL 补齐后）

- **首发生态教训（如实）**：第一发 `exit status 0xc0000135`（STATUS_DLL_NOT_FOUND，0.070s 即死）——`build/*.dll`／`third_party/sherpa-onnx/*.dll` 不入库，clone 导出树必须照票面 AC#7 自己那句**手工补**；本腿从主树拷 6 枚 DLL 进导出树后重放。
- **整包一发（终态）**：`PATH=$PWD/third_party/sherpa-onnx:$PWD/build:$PATH go test ./cmd/wisp/ -count=1 -timeout 30m -v`，起 11:40:53，**`panic: test timed out after 30m0s`，1802.442s，rc=1，无 `ok` 行**。
  - **四数＝RUN 6 / PASS 3 / FAIL 2 / SKIP 0**（`^=== RUN` 含子测试 ⇄ `^--- PASS` 只数顶层；本腿逐枚点名，未跑完全集＝**FAIL，如实报**）。
  - 终态红名册（逐名＋红因）：
    1. `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`（24.71s）——`always_write_no_clobber_226_test.go:119`："the guarded write applied the hand-edited theme … reading a hand edit back is ticket 223"；盘上历史底片＝**FAIL 17 枚 log／PASS 43 枚 log**，属机器态两见老噪声族。
    2. `Test219AC2WidenAppliedAuditLineCarriesTheFourFields`（84.85s）——`approval_always_219ac2_test.go:187`："the L2 card for t219-ac2-audit-corr never reached the ledger"；盘上历史底片＝该名 **FAIL 0／PASS 2**（其同族兄弟 `Test219AC2RevokeRoundTripByteForByte` 历史 FAIL 1／PASS 2）——今日红无历史 FAIL 先例，如实记。
  - **死点（卡死非崩）**：`Test219AC2RevokeRoundTripByteForByte` 挂满 26m55s 被超时帽打穿；goroutine 310 阻在 `approval_always_219ac2_test.go:243` 的 `<-w.done`（`chan receive, 26 minutes`）＝测试件自己的子进程驱动器没等到回执，**不是产码 panic、不是 329 新产码打红**；与 10-17 六发尾停 MARK299 真窗区死亡同属"真窗/子进程区在本机饿死"族（今日本机三腿同飞＋ac5 腿已记的面板窗环境异常，见票尾 11-08 bullet）。
  - ⛔ 本腿不把上两枚红记到 329 头上、⛔ 不重测到好运。
- **另路取数（329-v1 已立、台账 A1014 已采的同一枚具名 13 枚 `-run` 名册，同树同环境）**：**RUN 13 / PASS 13 / FAIL 0 / SKIP 0，rc=0，`ok … 72.997s`**。红名册两侧皆空＝对 v1 交件树基线（13/13/0/0）**作差 0**；对 v1 起手树（12/12/0/0）唯一差枚＝`TestBookedSnapshotReachesHostBitForBit329`，今天 **PASS 1.48s**（`--- PASS: TestBookedSnapshotReachesHostBitForBit329 (1.48s)`，现读在 `329-closure-roster13.log`）。
- 底片：`/d/tmp/wisp/20261108/329-closure-full.log`（整包 228 行）＋`/d/tmp/wisp/20261108/329-closure-roster13.log`（另路 13 枚）。

## 6. 门禁四数（导出树 `6e2eb5845`，本腿现跑）

| 门禁 | 命令 | 读数 |
|---|---|---|
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l ./cmd/wisp/`（真身绝对路径） | **空，rc=0** |
| gofumpt | 同尺 `./internal/` | **空，rc=0** |
| vet | `go vet ./cmd/wisp/` | **rc=0** |
| build | `go build ./...` | **rc=0** |
| d22scan | `sh scripts/d22scan.sh` | **clean，rc=0**（ban #8 internal/=963、cmd/=241；ban #6 frontend/=119） |

## 7. 本腿记账（读数归读数，翻框归编排者）

- AC#0：三把尺本腿全部现跑＝**名册字节合格**（七笔落地只落 cmd/wisp/**＋probes/329/r1/**）；缺陷仍在件不在字节——r1 证据件 3 处 `待填` 未回填，与 329-v1 "AC#0 不成立" 判语时点同形，本腿⛔ 代填、⛔ 翻框，是否凭本件第 2/4 节读数翻格归编排者。
- AC#7：整包终态四数不达标（未跑完全集）＝**FAIL 如实报**；另路 13 枚名册 13/13/0/0 rc0、红名册作差 0。两维读数并列，采哪维判格归编排者。
- 本腿只 commit 不 push；改动件＝本件＋票 329 票尾一条 bullet（纯追加）。
