# 票 260 · AC#2 验收腿 `260-v2` 判语表（非实现者；实现者＝`260-r5`）

**本腿是什么**：票 260 AC#2「零症状那一格要有声」的**非实现者判语腿**。被裁对象＝`internal/ball/hotkey_borrow_refused_260r5_test.go`（448 行，首行 `//go:build windows`，无 `winlive`；三枚用例）＋`260-r5` 的九节证据件与两枚复跑笔。⛔ 本腿未翻票面任何框、未改 `Status:`／票名、未改任何产码；只追加票尾 1 条 bullet。
**锚**：起手取数 `2026-10-14 10:39:02`，HEAD＝`b8f883a3`（`A907`，10:33:18；派单发送时刻那枚 `1179e232` 已被它 +1）。`go version go1.27.1 windows/amd64`。
**写面**＝`docs/evidence/s1/ticket260-260-v2-accept-v1.md`＋`.scratch/wisp/probes/260/v2/**`＋票尾 1 bullet。⛔ `cmd/wisp`／`internal/tools`／`.github`／`internal/panel`／`internal/agent` 零触碰。

## §0 判语（先结论）

| # | 派单要求的格 | 判语 | 凭据 |
|---|---|---|---|
| 1 | 复跑绿（非突变态，三读数逐字） | **成立** | §1：定向 3 PASS／0 FAIL／**0 SKIP** rc=0；整包 68／1／0 rc=1（唯一红＝既有环境红） |
| 2 | 自造反证 ≥2 发（≥1 发 r5 没种过的形） | **成立** | §2：**五发**真跑（含两发 r5 从未种过的形＋r5 的 M4 复种）；每发仓外备份＋sha256 起手终值全等 |
| 3 | AC#2 判语 | **成立（无条件）**；三条具名残差**不构成条件**（§8） | §2 的红句**指名 `ball_windows.go:895` 那一笔**；尺零 `t.Skip`、零 SKIP；与 258 分开（§5） |
| 4 | 越界核（r5 名册／产码零改动／与 258 分开） | **成立** | §4：8 枚 r5 commit 名册全在 `probes/260/r5/**`＋1 枚测试尺件；`ball_windows.go` 自 r5 起手锚**零字节改动** |
| 5 | winlive 谨慎（能跑则跑、跑不了具名） | **跑了，读了** | §6：Esc 借用族 2 枚 **PASS** rc=0（10:46:03→10:46:07）；跑前跑后 `tasklist` 零残留；无占键冲突 |

## §1 复跑绿（非突变态，逐字三数）

1. 定向（派单原命令）：`go test -count=1 -v -run '260r5|BorrowRefused' ./internal/ball/`，10:39:13→10:39:15
   ⇒ **PASS 3／FAIL 0／SKIP 0，rc=0**（档 `logs/01-targeted-v.txt`，1874 B）。`=== RUN` 恰 3 行＝pattern 没匹配到别人。
2. 全名册：`go test -count=1 -v ./internal/ball/`，10:39:15→10:39:26 ⇒ **PASS 68／FAIL 1／SKIP 0，rc=1**；唯一红＝`TestC21TableColourRowsMatchTokensCSS`（`design/assets/tokens.css` 被别的会话删，早于本票的既有环境红，⛔ 不修不追；档 `logs/02-full-v.txt`，25355 B）。
3. 加固复跑：`go test -count=3 -race -v -run '260r5|BorrowRefused' ./internal/ball/`，10:43:09→10:43:18 ⇒ **9 PASS／0 FAIL／0 SKIP，rc=0**（档 `logs/03-race-count3-v.txt`）；r5 两程都没跑过 `-race`，这一发是本腿新加的。
4. 尺自身三查（本腿现跑）：首行 tag＝`//go:build windows`；`t.Skip(` 真调用＝**0**（唯一一处是 `:9` 注释里的字样）；非 winlive 真调用＝**1 枚＝`:202` 的 shipped `b.TakeEscForCancel()`**（winlive 侧 5 枚＝`hotkey_cancel_borrow_live_260_test.go:37/:70`＋`hotkey_live_test.go:373`＋`interaction_live_test.go:141`＋`live_windows_test.go:126`，档 `logs/12-ruler-scope.txt`）。

## §2 自造反证五发（本腿自己写的台件 `mutate260v2.py`，195 行）

台件规矩：起手互锁＝sha256(working)==sha256(HEAD blob)==sha256(仓外备份)，不等**拒跑**；还原源＝**仓外备份**（`C:\tmp\wisp260v2backup\*.orig`）＋HEAD blob 双份；每发记窗口起止、突变后 sha256、还原后 sha256。**最终一轮窗口 10:49:23→10:49:34**（另一轮四发版 10:42:09→10:42:23 被同一台件重跑覆盖，读数逐字相同，仅时刻不同——供他腿归因）。总账＝`logs/mutation-sha256-ledger.txt`，控制台全文＝`logs/mutation-run-console.txt`。

| 发 | 形状 | 与 r5 的关系 | rc | 红（逐字见下） |
|---|---|---|---|---|
| **M-A** | 摘掉 `ball_windows.go:895` 失败回执写入（票面判据原形） | r5 M1 同形；**突变后 md5＝`51d4024f927c55a84565dee91039cac1`，与 r5 账逐字相同** ⇒ 字节级同一发 | 1 | 两枚接线尺红＋既有环境红（整包 PASS=66 FAIL=3） |
| **M-B** | 保留写入但改成**静默**（回执＝`cancelIdleLine`） | **r5 两程都没种过** | 1 | 两枚接线尺红（写在了但没声） |
| **M-C** | 摘掉 `takeEscWithAcc` 的 `slog.Error`（另一 symptom 通道） | **r5 两程都没碰过 `hotkey_windows.go`** | 1 | 一枚红在**日志那句**断言 |
| **M-D** | 回执仍写失败、但登记集**谎称 live**（`:896` 更强读法） | **r5 没种过**（M4 的反方向） | 1 | 前提尺红在**登记集**断言 |
| **M-E** | 只摘 `:896` 登记集镜像（拿回执原样） | r5 M4 同形；**突变后 md5＝`3cfc5ee9d780e796aa20abae9f4b2791` 与 r5 账逐字相同** | **0** | 三枚全绿 ⇒ 复现 r5「本尺对 stale 镜像不敏感」的边界 |

红句逐字（档 `logs/mutation-M-*-gotest-v.txt`）：

```
M-A: hotkey_borrow_refused_260r5_test.go:283: the cancel line carries no error (...Err:<nil>...): the rig never met a refusal
M-A: hotkey_borrow_refused_260r5_test.go:325: problem lines after a borrow Win32 refused = 0 ([]), want the 1 that says so: the seed said [], and the write at ball_windows.go:895 is the only thing between the two
M-B: （同上两句逐字，:283 与 :325）
M-C: hotkey_borrow_refused_260r5_test.go:383: the refused borrow logged nothing naming "cancel takeover refused by Win32"; error-level output captured was ""
M-D: hotkey_borrow_refused_260r5_test.go:291: the registration set still claims the cancel key after a refusal: map[1:{16387 81} 2:{16387 77} 3:{16387 120} 4:{16387 80}]
M-E: （零红：--- PASS 三枚，rc=0）
```

⇒ **判据兑现**：M-A 两枚红句之一**逐字点名 `ball_windows.go:895` 那一笔**是种子与借后之间唯一的东西；M-B 证明尺读的是**内容**不是"发生过一次写入"；M-C／M-D 证明两枚 symptom 通道（日志、登记集）各自有牙；M-E 复现不敏感边界。既有语义层两枚（`TestBorrowFailureNamesTheKeyItTried260`／`TestCancelBorrowFailureIsAProblemLine`）在 M-A 下**保持 PASS**（本腿现跑，档 `mutation-M-A...txt`）＝与 `260-v1` 量到、r5 要治的形状一致。三枚用例里的反形腿 `TestRefusedBorrowIsNotAStandingAnnouncement260r5` 在 M-A/B/C/D 四发下**全绿**＝正尺不是恒定红。

## §3 字节还原链（sha256，不是 md5；md5 仅作 r5 对账）

- 起手：`ball_windows.go`＝`89c543601be926a1df26b6b0166f7c03952f62ebeca955c7c889099c975984d1`，三份（working／HEAD blob 副本 `logs/ball_windows.go.head-blob`／仓外备份）**逐字节相同**；`hotkey_windows.go`＝`13b15011c73544d16930f9a8902fea0dbeb35dd2b7e460aeea57f84c50850cdc` 同法。
- 五发逐发还原后 sha256 **全部回到起手值**（账里 `restored-ok=True` 五行）；`FINAL ... matches_start=True` 两行。
- 终态尺：`git status --porcelain -- internal/ball/`＝**空**；working hash-object＝HEAD blob＝`af9ff08a57cdf7f8535fda790f725cd15f21fc1d`。

## §4 越界核（r5 名册＋产码零改动）

- r5 全部 **8 枚** commit（`8174c325`／`629f146b`／`ed892e6e`／`6ab2efdd`／`5ad1a7b8`／`5d680d78`／`b617084e`／`8a463bc6`）逐枚 `git show --name-only`：probes 之外只有 **1 枚路径**＝`internal/ball/hotkey_borrow_refused_260r5_test.go`（测试/尺件，Oct 4 落）。派单点名的两枚复跑笔 `b617084e`／`8a463bc6` 的 probes 之外行＝**0**（档 `logs/10-r5-rosters.txt`）。
- 产码零改动：`git log -1 --format=%H -- internal/ball/ball_windows.go`＝`b1e59d63`（2026-10-04 10:01:15，**早于 r5 第一枚 14:09**）；`git diff c6cf66e6..HEAD -- internal/ball/ball_windows.go`＝**空**（c6cf66e6＝r5 复跑程起手锚）；`git diff c43a791a..HEAD -- internal/ball/`＝只有那枚测试尺件 +448 行，产码 0 行；`git log --oneline c6cf66e6..HEAD -- internal/ball/`＝**空**（档 `logs/11-boundary.txt`）。

## §5 与票 258 的 rebind×借还自钉分开（两句射程文本各引一句）

- **258 那把尺**（rebind×借还那一格的 live 形态）逐字："`TestLive258ResidentRebindsWithoutRestart is AC#2's live form: the process is UP on a planted non-default summon, then config.toml is edited underneath it; the poll (1s) must rebind the live keys inside the same process.`"（`cmd/wisp/resident_hotkey_live_258_windows_test.go:107-109`）。它的场景是**改配置时正有借用在飞**（残留白纸黑字在 `internal/ball/ball_windows.go:807` 的 `RebindHotkeys` 文档："`A rebind drops an in-flight Esc borrow`"）。
- **260-r5 那把尺**逐字："`plant "the borrow request went out but Win32 never lent the key" and make the step that says so go red - in a non-winelive case`"（`hotkey_borrow_refused_260r5_test.go:7-9`）。它种的**只有"借被 Win32 拒"**，全文件没有一处 `Rebind`。
- 交叉核（本腿现跑）：260 尺里 `Rebind` 计数＝**0**；258 那枚 winlive 尺里 `TakeEscForCancel|Borrow` 计数＝**0**；两把尺一在 `internal/ball`（非 winlive）、一在 `cmd/wisp`（winlive）。⇒ 没合并成"热键接好了"，守住了。

## §6 winlive 半（派单令：能跑则跑、跑不了具名）

- 协议：跑前 `tasklist`＝**零** wisp／go／balldebug 进程（`logs/20-winelive-prep.txt`）；`internal/ball` 干净、起手 sha256 在场。
- 跑：`go test -tags winlive -count=1 -v -run 'TestLiveCancelBorrow' ./internal/ball/`，10:46:03→10:46:07
  ⇒ `TestLiveCancelBorrowDefaultUnchanged260` **PASS**（0.35s，裸 Esc 真借走→`requireEscReturned` 归还）、`TestLiveCancelBorrowFollowsConfig260` **PASS**（0.04s，Ctrl+Alt+Y 真注册、Esc 保持空闲、1409 探测）；rc=0（档 `logs/21-winelive-cancelborrow-v.txt`）。
- 跑后 `tasklist`＝**零残留进程**（同档尾行）。**无占键冲突**（两枚没有走任何 SKIP-LOUD 支）。⇒ 这是该族**第一次真跑**（r5 两程都只 `go vet -tags winlive`），读数与 AC#2 的接线主张同向。

## §7 门禁（档 `logs/30-gates.txt`）

| 门 | 读数 |
|---|---|
| `sh scripts/d22scan.sh` | **rc=0，clean**（276 production Go files；ban #8 internal=550、cmd=125 Go 含注释与 `_test.go`） |
| `go vet ./internal/ball/` | **rc=0** |
| `go vet -tags winlive ./internal/ball/` | **rc=0**（证明 winlive 那侧编得动） |
| `gofumpt -l`（逐枚点名那把尺，GOPATH＝`D:\work\base\gopath`） | **空输出 rc=0** |
| 新路径长度 | 最长 **89** 字符（`logs/mutation-M-D-receipt-failed-set-claims-live-gotest-v.txt`），在规则 9 帽 121 之下，⛔ 未添长度债 |

## §8 具名残差（逐条：为何**不构成** AC#2 的条件）

1. **`:896` 只摘（stale 镜像）本尺不敏感**：M-E 本腿亲手复现 rc=0。AC#2 的判据句是"种'借用没成键'⇒**指名那一步**必须红"，指名的那一步＝`:895`；`:896` 是镜像，属另一格（r5 也自报"属另一格"）。本腿 M-D 把边界另一侧（谎称 live⇒`:291` 红）也量了 ⇒ 边界两侧都有读数，**归编排者判要不要单开一格**。
2. **`cmd/wisp` 侧非 winlive 真调用＝0 枚**（本腿独立复核：该目录唯一提到它的是一行注释 `resident_cancel_key_wording_260r3_windows_test.go:21`）。它是"回执→界面那句话"端到端那一跳的欠账，与 AC#1／AC#4 同一个〔未实测〕笼子，**不是本格的判据**；按票面第 33 行本腿未跨包。
3. **winlive 1409（真被占）与默认档被拒两形本尺不覆盖**：本尺的拒绝是 Win32 句柄校验层（1400），种子用组合键；这两形要么赌真注册、要么依赖句柄层，均归 winlive/机主面（r5 §5 已具名，本腿复认）。

## §9 没做的格与并发披露

- ⛔ 没翻票面任何框（只票尾 +1 bullet）；⛔ 没改 r5 任何件一字；⛔ 没改台账；⛔ 没修既有环境红；⛔ 没 push；⛔ 没用 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`。
- **并发披露（逐时段，供他腿归因）**：突变共两轮窗口——**10:42:09→10:42:23**（四发）与**10:49:23→10:49:34**（五发，含 M-E）。10:43:5x 本腿看见**别腿**的 `wisp.test.exe`（PID 25404）＋`go.exe`（PID 16484）在台面（首查 start time 失败、第二查已退出，故起始时刻**量不到**）——若那是 `315-a1` 的 `cmd/wisp` 测试且其编译窗口与上面两段重叠，它的读数可能被本腿突变洗过：**已逐段报时刻，请对时**。本腿自己的读数（§1／§3／§6）都在突变窗口之外取。
- 交件名册（`.scratch/wisp/probes/260/v2/`）：`mutate260v2.py`(7767 B)＋`logs/`22 枚＝锚 `00-anchor.txt`(581)、跑档 `01-targeted-v.txt`(1874)／`02-full-v.txt`(25355)／`03-race-count3-v.txt`(5346)、名册与核 `10-r5-rosters.txt`(3293)／`11-boundary.txt`(857)／`12-ruler-scope.txt`(5413)／`13-tags-and-doc.txt`(1142)、winlive `20-winelive-prep.txt`(301)／`21-winelive-cancelborrow-v.txt`(1419)、门禁 `30-gates.txt`(23163)、突变账与五发日志 `mutation-sha256-ledger.txt`(2959)＋`mutation-M-A/B/C/D/E-*.txt`(25237/1795/2009/2225/2037)、控制台 `mutation-run-console.txt`(2976)、还原源 blob 双份 `ball_windows.go.head-blob`(33995)／`hotkey_windows.go.head-blob`(27476)。
