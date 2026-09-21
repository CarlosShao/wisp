# 106 — 对抗验收裁决表（acceptor-93-106，独立于 `agent-ticket106`）

**验收时间：** 2026-09-21 19:12 CST（`date` 实测）
**被验对象：** `.scratch/wisp/issues/106-winsec-private-set-trips-on-la-on-the-ci-runner.md`，commits `1400e15`/`34f6959`/`b58bf2f`/`0f4c891`/`0a44598`/`c1da933`
**我的树：** 本地读数全部来自 `git archive 440dd88` 的纯净快照 `/tmp/ac936-t106`（**不含票 108 的两枚未跟踪文件** `seam_bypass_108_windows_test.go`/`ancestor_separator_108_other_test.go` ⇒ 别人的半成品不算回归）；变异在第二枚独立快照 `/tmp/ac936-mut`。**HEAD 工作树验收期间已推进到 `f4bf0fa`，我没有在带他人 WIP 的树上取过任何绿。**
**CI 读数树：** run `35591482293`（head `440dd88`），job `test-windows` id=`106306750494`。
**总判：PASS WITH CONDITIONS**（缺陷在 runner 上确实消失、反半边仍有牙、集合未被放宽；条件是 **CI 上没有任何步骤跑 winsec 自己的用例** ⇒ "runner 那一格"只被 `internal/secret` 间接证实）

档位：**① 〔独立复现〕** / **② 〔日志＋归档，我抽验〕** / **③ 〔仅自述，不背书〕**（第三档在本表里存在，且每条都写补救）。

---

## 逐格裁决（与票面 AC 1:1）

| AC | 判 | 档位 | 我亲自拿到的证据 |
|---|---|---|---|
| **AC#1** 私有集定义/使用点落到 file:line + 两侧 temp DACL 差异量出来 | **通过** | ① | 锚点我逐行打印核对（`sed -n` 原文）：`:408 privateSet()` 返回 `map[string]bool{system, admins, me}` **三员**，三者都过 `canonicalSIDString`，**表里没有任何名字**；`:362` 判定 `if !ace.grant || !set[ace.trustee]`；`:371` `if ace.trustee == me`；`:96` 通知腿 `if ace.grant && set[ace.trustee]`，否则 `out = append(out, ace.trustee+"("+ace.text+")")` ⇒ 拒绝/通知侧带表示形式、放行侧只认 SID。**本机 `icacls %TEMP%` 我自己跑了一遍**，与实现方逐字一致：`CodexSandboxUsers:(OI)(CI)(M,DC)`、`S-1-5-21-3623186960-731165060-4091685855-1717338598:(OI)(CI)(M,DC)`、`NT AUTHORITY\SYSTEM`、`BUILTIN\Administrators`、`DESKTOP-LVS7839\swq` 三条 `(OI)(CI)(F)`。runner 那份只存在于票面 CI 原文（`D:PAI(A;;FA;;;SY)(A;OICIIO;GA;;;SY)(A;;FA;;;BA)(A;OICIIO;GA;;;BA)(A;;FA;;;LA)(A;OICIIO;GA;;;LA)`）⇒ 档位②（日志档）。 |
| **AC#2** 造一台没有 `LA` 也能验的仪器，修前必须红，不许 `t.Skip` 挡 | **通过（有条件）** | ③（"修前红"那一格）＋①（"仪器有牙"由我另路证） | 我**没能**独立复现"修前红"：在 `git archive 1400e15^`（=`a64d06f` 态生产码）上放入 `private_set_sid_windows_test.go` 会**编译不过**（第三条腿点名新函数 `privateSet`），而这正是实现方自己在票面披露的手法（"那一跑里把 `:183-206` 的第三条腿摘掉了"）⇒ 它的披露**方向诚实**，但那条红我只到档位③。**补救**：`git archive 1400e15^ \| tar -x -C /tmp/ac106-pre && cp <HEAD 的该测试文件> ...` 后**手工摘第 3 腿**再 `go test -count=1 -v -run TestSealNarrows...`，两行红名（`:243`/`:290`）应由下一手重新产出。**我没有因此判 FAIL**：我用两枚变异（MUT-A/MUT-B/TEETH，见 AC#4 格与第 3 条）独立证明这三条腿**都在真实承载行为**，不是摆设。用例体内**无 `t.Skip`**（我在快照里 `grep -n "t.Skip" internal/winsec/private_set_sid_windows_test.go` 为空）。 |
| **AC#3** 该拒的仍拒 / 不该拒的不拒；放行侧只认唯一已解析形式 | **通过** | ① | **反半边我自己跑**：HEAD 生产码下 `TestGateRefusesADescriptorThatLeavesARealGrantToAnotherAccount` **PASS**（`-count=1 -v`，0.150s，`ok internal/winsec`）⇒ 给另一个真实主体（`S-1-1-0`，用 SID 种的 ALLOW 授权）**仍然被拒**。并且我把校验关掉后它**必红**（MUT-A，见下）⇒ 拒的**原因**就是这条判定，不是巧合。正半边 `TestGateJudgesThePrivateSetByResolvedSID` 三条子腿全 PASS。**集合未被放宽**：`privateSet` 仍是三员，`LA` 不在里面（我用 TEETH 变异往集合里塞名字 ⇒ 立刻红 ⇒ 结构守卫有效）。 |
| **AC#4** 变异：整块去掉 ⇒ ①那半红；什么都不查 ⇒ 必须红 | **通过** | ①（两发全部我自己重跑） | **MUT-A** `:362` → `if false && (!ace.grant || !set[ace.trustee]) {`：`go build ./internal/winsec/` **rc=0**（编译失败不算变异 ⇒ 成立），同链 `grep -n` 落地=`362:`；`-count=1 -v` 三条 ⇒ **rc=1**，`--- FAIL: TestGateRefusesADescriptorThatLeavesARealGrantToAnotherAccount` @ `private_set_sid_windows_test.go:312`「a DACL leaving S-1-1-0 a real grant on the object was accepted as private: trustees [... 0/0x0=S-1-1-0]」，另两条 PASS。**MUT-B** = MUT-A + `:96` → `if true || !set[ace.trustee] {`（什么都不查）：build rc=0 ⇒ **rc=1、2 条 FAIL**：上面那条 + `--- FAIL: TestSealNarrowsAndNamesThePrincipalItRemovedBySID` @ `:268`「the notice named the cleared principal by spelling only ... cleared=""」。⇒ **"取消检查换不来绿"，两半各红在自己的腿上，与实现方自报逐字吻合**。还原：`cmp` 快照文件 vs 纯净快照同一文件 ⇒ `RESTORED-BYTE-EXACT`；仓库 `internal/winsec/winsec_windows.go` 从未被我改过。 |
| **AC#5** 门禁（本机全绿不算过，欠 push 后的 `test-windows` 步级结论） | **通过（有条件，CI 那半由我补齐）** | ①（本地四数＋CI 步级＋日志原文） | **CI 半边**：run `35591482293`（head `440dd88`）job `test-windows` **id=`106306750494`**，**steps 逐条**：`6\|Portable windows tests (proc/secret/config)\|success`、`5\|cgo build smoke\|success`、`7\|PathResolver junction placeholder\|success`，job conclusion=success。日志（我 `gh api .../jobs/106306750494/logs`，763 行）第 249/389/449/718/719/720 行：`portable-tests.sh: platform=windows scope=[./internal/proc/ ./internal/secret/ ./internal/config/]`、`ok github.com/CarlosShao/wisp/internal/proc 11.781s`、**`ok github.com/CarlosShao/wisp/internal/secret 0.391s`**、`ok .../internal/config 0.881s`、`runtests.sh: OK ... top-level: PASS=108 FAIL=0 SKIP=0, === RUN=164`、`four numbers: === RUN=164 --- PASS=108 --- FAIL=0 --- SKIP=0`。⇒ **票面那个"`internal/secret` 整包 8+ 条在 `NewStore`/`MigratePlaten` 第一步就死"的缺陷在 runner 上消失了**，且不再是"步级 success 而内部没结论"的形状。**条件/残留见 `R-106-1`**：同一份日志里 `internal/winsec` 出现次数 = **0** ⇒ **CI 上没有任何步骤跑 winsec 自己的用例**，所以"三条新用例在 runner 上 PASS"这一半**没被证实**（`internal/secret` 变绿是**间接**证据）。**本地半边**（纯净快照 `-count=2 -v` 四包）见第 5 条。 |

---

## 编排者点名要我自己打的两点（＋三件附带）

### 1. **因果对不对：旧码错在"名字比较"，不是"集合成员"** —— 实现方的因果裁决**成立**（①，读旧码原文）

`git show 1400e15^:internal/winsec/winsec_windows.go`（=修前生产码）第 276-352 行原文：

```go
allowed, me, err := allowedSIDStrings()          // :292
for _, ace := range aceGroups(sddl) {            // :298  sddl = sd.String() —— OS 印出来的拼写
    fields := strings.Split(ace, ";")
    if len(fields) < 6 || fields[0] != "A" { foreign = append(foreign, "("+ace+")"); continue }
    if !allowed[fields[5]] { foreign = append(foreign, fields[5]); continue }   // :304  ← 用名字查表
    ...
}
func allowedSIDStrings() (allowed, me map[string]bool, err error) {
    u, err := user.Current()
    allowed = map[string]bool{"SY": true, sidSystem: true, "BA": true, sidAdmins: true, u.Uid: true}   // :343-346
    me = map[string]bool{u.Uid: true, "ME": true}                                                     // :347
}
```

⇒ 集合**里本来就有令牌用户**（`u.Uid` 是 SID 串）。失败的是 `allowed[fields[5]]` —— `fields[5]` 是 **SDDL 的 trustee 拼写**，runner 上 OS 把内置 Administrator 的 SID 印成 `LA`，`LA` 不是那张表的键 ⇒ 拒。**这是名字比较错，不是成员错。**

**"把 `BA` 组与 `LA` 账户混为一谈"这个替代解释我排除了**：旧表压根不查令牌组/组成员关系，纯字符串相等；`BA` 在表里只以 `"BA"` 与 `S-1-5-32-544` 两个键出现。

**"runner 上 `u.Uid` 解析成 `runneradmin` 而 ACE 落的是 `Administrator`"这个替代解释我也排除了，理由是构造性**：那份 DACL 只有 winsec 自己写得出来 —— 落盘器只放三个主体（`:261 sidStrings(sidSystem, sidAdmins)` + `:334 convertSID(u.Uid)`），且 `D:P` 已在 ⇒ 继承被剥。CI 那份是**六条三对、无 `ID` 位**（票面原文与实现方本机回读同形），第三对**只能是 `convertSID(u.Uid)` 那条**被 OS 用名字印回。⇒ 若真存在"`u.Uid` ≠ ACE 主体"，我们看到的会是**四个**不同主体（SY/BA/令牌用户/LA），日志里没有。**8.3 短名 `RUNNER~1` 与判定无因果**这一点同样成立：它是 `TEMP` 环境变量带进来的**路径**拼写，而修后的判定链（`:172` 二进制 ACE ⇒ `(*SID).String()`、`set` 全 SID）不经过任何路径名字解析。

**因此修法仍然成立**：把比较端换成"唯一已解析形式"正是打在真实根因上；不是靠放宽躲开症状（第 2 点独立证）。

### 2. **有没有偷偷放宽：没有**（①，三发证据）

- **集合仍是三员**：`privateSet` 的返回字面量 `{system, admins, me}`（`:432`），`LA`/`BA`/`SY` 之类**名字键一枚没有**。
- **反半边仍拒**：HEAD 下 `TestGateRefuses...` **PASS**（给 `S-1-1-0` 留了真实 ALLOW ⇒ 报错并点名 SID）。
- **变异（我重跑，不用它的数）**：`:362`→`if false` ⇒ **反半边红**（`private_set_sid_windows_test.go:312`，rc=1，build rc=0）；再连通知腿 `:96`→`if true ||` ⇒ **两条红**（上面那条 + `TestSealNarrows...` @ `:268`，`cleared=""`）。⇒ 绿是**判定换来的**，不是"取消检查"换来的。
- 加发 **TEETH**（它没做的）：`:432` 往集合里塞 `"LA","BA","SY"` 三个名字 ⇒ build rc=0、`-run TestGateJudgesThePrivateSetByResolvedSID` **rc=1**，红在 `--- FAIL: .../the_set_holds_no_name_in_the_form_it_is_compared_with` @ `:195`「the private set admits a principal by name, not by SID: "LA"」。⇒ 这条"放宽即暴露"的守卫**是真的**。

### 3. **`TestGateJudgesThePrivateSetByResolvedSID` 本机恒绿怎么判：算"钉 runner 形状的仪器"，不是死用例 —— 它有牙**（①）

我把它**人为弄红过一次**，做法与结论：把集合改回"按名字也认"（TEETH，`:432` 塞进 `"LA"`/`"BA"`/`"SY"`）⇒ 三条腿里 `the_set_holds_no_name_in_the_form_it_is_compared_with` **FAIL**（`:195`），另两条仍 PASS。⇒ **它的第二条腿承载真行为**（守卫"判定链里不许出现名字"），与"本机 LA 是不是令牌用户"无关。它的第一/三条腿（`runner_descriptor_is_judged_by_what_it_resolves_to`、`a_name_spelling_cannot_change_the_verdict`）在我这台机器上走不到 `LA == 令牌用户` 那一支 —— 本机 `LA = ...-500` 是别人，这一点实现方**主动声明**过（票面 18:52 段），我核实其声明与代码形状一致。**判：合理仪器，非残留死用例；但见 `R-106-1`** —— 它钉的那一格在今天的 CI 上**根本没有步骤会跑到**，所以"钉 runner"目前只是意图，不是台账。下一条命令：在 `test-windows` 加一步 `bash scripts/portable-tests.sh ./internal/winsec/`（票 93 已把这台仪器修好，接上即可），然后读该步的步级结论。

### 4. AC#5 由我补：**run id = `35591482293`，job id = `106306750494`，step 6 `Portable windows tests (proc/secret/config)` conclusion = `success`**，日志含 `ok .../internal/secret 0.391s`、`FAIL=0 SKIP=0`。票面那个"run id 位 = ____"现在可以填。⚠ 未拿本地绿替代。

### 5. 门禁复跑（① 全部我自己跑，纯净快照 `/tmp/ac936-t106`）

```
go test -count=2 -v ./internal/winsec/ ./internal/secret/ ./internal/config/ ./internal/proc/
rc=0
=== RUN=410   --- PASS(顶层)=270   --- FAIL=0   --- SKIP=2
ok internal/winsec 9.113s / ok internal/secret 0.354s / ok internal/config 1.758s / ok internal/proc 2.760s
```
- **口径差异如实报**：实现方写 `PASS=407 / 不同测试名=205（×2 恰等）`。我的 `^--- PASS` **只数顶层**（Go 的子测试 PASS 行是缩进的 `    --- PASS:`）⇒ 顶层 270；`410/2=205` 那个"205"是**含子测试名的 `=== RUN` 行**数，与我 `go test -list '.*'` 的**顶层用例数**不是一回事。我的顶层用例数：`winsec 27 + secret 21 + config 56 + proc 32 = 136 个不同测试名`，×2 = 272 条顶层 RUN，余下 138 条 `=== RUN` 是子测试。**两边都报，不重测到运气好。**
- **SKIP 逐条点名（且是 `-v` 读数）**：`TestHelperProcess` ×2（`internal/proc/jobscope_windows_test.go:87`，re-exec 助手本非用例）⇒ 已由票 93 的台账第 3 条显式记账；这条命令是裸 `go test`，所以它跳了没让 rc≠0 —— 正是票 93 修的病，不是本票的账。
- **FAIL=0 ⇒ `TestExternalSamplerReadsSubjectFromOutside` 我没复现到红**（它的自述是 `internal/proc` 的 `NtQuerySystemInformation: buffer never sufficient`）。我在同一次纯净跑里 proc 是 `ok 2.760s`；另加压（同时并发 `-count=2` 四包 + risk 的 budget 用例）proc 仍 rc=0。**判：负载假红，与票 103/86 同一口径；不是回归。** 多样本：加压腿 rc=0，普通腿 rc=0。
- `gofmt -l internal/winsec internal/secret internal/config internal/proc scripts` ⇒ **空**；`gofumpt.exe`（`$(go env GOPATH)/bin`，本机**有**）`-l` 同四包 ⇒ **空**（rc=0）；按包 `GOOS=linux go vet ./internal/winsec/ ./internal/secret/ ./internal/config/ ./internal/proc/` **rc=0**，同 scope 原生 `go vet` **rc=0**。
- `sh scripts/d22scan.sh`（**在纯净快照 `/tmp/ac936-t106` 里跑的，不是仓根 `go run ./tools/d22scan`**）：**rc=0 clean**，台账读数逐字：`bans #1-5 internal/=197 cmd/=20, ban #6 frontend/=37, ban #7 internal/tools/=17, ban #8 design/=16 frontend/=37 internal/=347 cmd/=26`。**与实现方自报的对照**：`internal/=347`（它登记的上一值 346 ⇒ 347，多的正是它那一枚测试文件）**逐字吻合**；`#1-5 internal/=197`、`cmd/=20`、`#7 internal/tools/=17`、`#8 design/=16`、`#8 cmd/=26` 全部吻合。**唯一不吻合的是 ban #6/#8 的 `frontend/`：它报 40、我在 `440dd88` 的纯净归档里量到 37。** 我判**不是降**、也不是本票的账：它那一跑在**仓库工作树**里，树里同时躺着票 92/77 未跟踪的 `frontend/` WIP（我今天 `git status` 仍见 `?? internal/panel/attachments.go` 等 5 枚未跟踪文件 ⇒ 同一现象），未跟踪文件会被计入但不进 `git archive`。⇒ 两份数**不可直接对照**，"不降"这条以 tracked-only 的 37 为准，登记在 `R-106-3`（谁碰 `frontend/` 谁把这条对齐）。
- **共树碰撞已按编排者要求避开**：票 108 的两枚未跟踪文件不在 `git archive` 里，我的 winsec 四数与生产码绿**不含**它们；它们的 `declared and not used: victimSidsBefore` 若把 CI 的 winsec 步弄红，**不算本票回归**。

---

## 禁改面自证 + 残留

未碰 `docs/PLAN.md`、`docs/specs/*`、`internal/risk/**`（含 `pathresolver*.go`、`assessor.go`、`rules_gateway.go`）、`tools/d22scan/**`、`allowlist.txt`（票 93 侧亦未变短/未变长）、任何阈值/golden/断言。变异只在 `/tmp/ac936-mut`；仓库工作树一枚未 add（我的 evidence 与票面追加段除外）。未 push、未建 worktree、未跑整仓门禁。

- **R-106-1（本票最重的一条）**：`test-windows` 的 portable 步 scope 是 `proc/secret/config`，**不含 `./internal/winsec/`**（我 grep 该 job 日志：`internal/winsec` 命中 0 次）⇒ AC#5 想要的"三条新用例在 runner 上 PASS"以及 `LA == 令牌用户` 那一支**至今没有 CI 台账**。补救命令（下一手）：`ci.yml` 的 `test-windows` 加 `run: bash scripts/portable-tests.sh ./internal/winsec/`，再 `gh api repos/CarlosShao/wisp/actions/runs?per_page=5` 取新 run id、读该步 steps。**在此之前，本票的"修对了"只有 `internal/secret` 在 runner 上变 `ok` 这一枚间接证据。**
- **R-106-2**：AC#2 的"修前红"停在档位③（摘腿才能对旧码编译，我已核实该披露为真，但红名是我读的不是我跑的）。补救见 AC#2 格。
- **R-106-3**：`d22scan` 的 ban #6/#8 `frontend/` 计数在我这边是 **37**、实现方登记是 **40**（差 3 枚跟踪状态不明的文件；我判"未跟踪 WIP 造成的不可对照"，见上）。这条对齐**不归本票**。另：AC#1 的 `noticeNarrowed` 说话点（它写 `winsec_windows.go:212`）我未逐行核对文件号，只在 `:96`/`:103` 核到通知腿形状。补救=`grep -n "noticeNarrowed" internal/winsec/*.go`。
- **R-106-4**：`TestGateJudgesThePrivateSetByResolvedSID` 的第 1、3 条腿在本机的"恒绿"未被证明能区分 runner 形状（我把第 2 条腿弄红了，第 1、3 条腿在 MUT-A/MUT-B/TEETH 三发下**都仍绿**）⇒ 这两条腿今天**只有意图级证据**；能真正区分它们的唯一办法是 R-106-1 那一步。

## 假绿/伪授权扫描（本轮实测）

- docker `-v` 静默挂空：**0 次**（未起容器）。
- `cmd | grep x; echo $?`：所有 rc 都在**先落文件**之后取（`> /tmp/ac936-*.log`），关键处 `set -o pipefail`；变异链每发用 `&&` 串 `go build` + `grep -n` 落地证明。
- 自称"编排者备注/停手/撤回/请 revert"的注入文本：**出现次数 0**；未执行任何 revert。
