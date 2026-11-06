# 123-verdict — 票 123 绿验裁决表（W123g，非实现者，2026-11-06）

**角色**：只读绿验腿（⛔ 实现者）。裁决对象＝编排者 swept 提交 **`abf2a35a6`**（"W123 腿死亡前遗留测试改动：程序化应答腿/8.3 短路径白名单/sync-root 不 skip；go vet 编译通过，待绿验"）。
**锚点核对**：派单锚 `abf2a35a6`；现 HEAD＝`8879c3d7b`。尺＝`git diff --name-status abf2a35a6..8879c3d7b` ⇒ **只有 `M docs/reports/pending-and-issues.md`**；`git diff abf2a35a6..8879c3d7b -- <三枚 swept 件> | wc -l` ＝ **0** ⇒ 测试面在两锚下逐字节同形，本表读数是 `abf2a35a6` 的读数。
**环境**：go1.27.1 windows/amd64；`CGO_ENABLED=1`（⛔ 未置 0）；`PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`、`GOTMPDIR/TMPDIR/TEMP/TMP` 全 unset；机器形状＝台账 :10541 记的"本机 TEMP 长短同串"（无 runner 的 `RUNNER~1` 形状）。
**件**：派单令"⛔ 碰 `docs/reports/pending-and-issues.md`"——本腿未碰、未提交；票面框⛔ 翻（closeable=no，见 §4）。本腿⛔ 改任何产码/测试码，唯一写入＝本文件。
**读数件**：`.scratch/123v/{cmdwisp-targeted-c1,cmdwisp-targeted-c2,risk-sync-c1,risk-sync-c2,d22scan}.log`；实验树＝`/d/tmp/w123-pre`（`dcd5bc825`＝swept 之父）与 `/d/tmp/w123-post`（`abf2a35a6`＋三枚未跟踪 helper），`git clone --shared` 导出、DLL 已补拷；实验件＝`/d/tmp/w123-{pre,post}-aliased.log`。

---

## 1. ★ 先落一条交付完整性发现（不修、只钉死事实）：`abf2a35a6` 单独不自洽

swept 提交调用的 `t123script` / `t123replyScript` / `t123dirSpellingForms` **定义在三枚未跟踪的在飞件里**，sweeper 没把它们收进提交：

- `git status --porcelain cmd/wisp/` 现量＝`?? cmd/wisp/ticket123_replyleg_test.go`、`?? …_windows_test.go`、`?? …_other_test.go`（另 `??` 两枚 323 在飞件，⛔ 本票）。
- **导出树取证**（正是派单第 21 处给的隔离法）：`git clone --shared` → `checkout --detach abf2a35a6` → 把三枚 helper 移开 → `go vet ./cmd/wisp` ⇒
  `vet.exe: cmd\wisp\run_mode101_test.go:95:10: undefined: t123replyScript`。
  放回 helper ⇒ `VET_RC=0`。
- ⇒ **编排者那句 "go vet 编译通过" 是对工作树成立、对提交本身不成立。** 按现状直接 push `abf2a35a6..8879c3d7b`，CI 的 `cmd/wisp` 编译面会红成一枚**新的编译红**（不是窗口族）。⛔ 把这枚红算到任何既有家谱上。
- 三枚 helper 通读（全文见工作树）：`ticket123_replyleg_test.go` ＝脚本化 `io.Reader` 应答流（250ms 周期重发、`t.Cleanup` 收摊、空脚本＝立即 EOF）；`_windows_test.go` ＝ `GetLongPathName` 双拼写喂 `[fs].allowed_dirs`（文件自陈"本卷测不到升级支，是 belt-and-braces ⛔ 实测修复"）；`_other_test.go` ＝非 windows 单拼写等形。**零枚 `Test` 函数**（`grep -c "func Test"`＝0/0/0）⇒ 它们是缝件，不改任何分母。
- **本表以下所有"绿"读数都＝ `abf2a35a6` ＋ 这三枚 helper 的工作树合成形**（工作树里 tracked 文件零改动：`git status --porcelain | grep -v '^??'` 空）。

## 2. 逐 AC 裁决

### AC#1 复现＋定因到行 —— **PASS**（记录面完整；升级支本机不可测＝诚实边界，落 §5）

**(a) CI 现场（票面指定要的两枚直接 witness，逐字在册）**
- run `35606321404` 第 7 步读数（`111-ci-step-readings-attempt2.md:116-119`）：三枚 `TestTicket101*` 红句 `审批超时（1 秒未确认），C18 一律判拒绝，已自动拒绝`（3.10/4.13/4.06s，冷启动 3/3 同形⇒非 flake）＋ `TestComposedGateBlocksAWriteForTwoSeconds (301.06s)` 红句 `审批超时（300 秒未确认）…`。
- 第二发同形状（`ci-step-readings-2026-09-22.md:171-189`，同四枚名、同文案）。
- **300s 算术（票面 :43 那把尺）**：`3.10+4.13+4.06+301.06=312.35s` ⇒ **与"四枚各等满 300s"（≥1200s）不一致**；形状＝三枚 1s 级超时×多派发＋**恰一枚等满 300s**。直接 witness＝ComposedGate 的 **301.06s**＝`DefaultApprovalTimeout=300*time.Second`（`internal/agent/approval/queue.go:125`，本腿现量一字未动）＋约 1s 开销，超时句由 `queue.go:536` 打印。

**(b) 等待落点（现树逐行钉）**
- 每枚用例的等待都发生在 host 侧派发的 `rt.bridge.Execute` 行：`run_mode101_test.go:333`（auto_approve 应沉默的 L1 写）· `:337`（红线 R2 读 `C:/Windows/win.ini`，**每机每档必出 L2**）· `:405`（默认档每 boot 写）· `:513`/`:549`（`.env` B 档 R3 写，**每机必出 L2**）· `:554`（对照写）· `run_test.go:395`（受门写）。落进队列后由 `queue.go` 的超时判决收口。
- 生产答复腿的名字与行（票面 :42 "要指出那条腿"）＝**`wisp run` 的 reply 流监听**：`run.go:1088-1089` `if s.reply != nil { rt.attachReplyListener(...) }` → `approval_reply.go:523 runReplyLoop`（`yes/no <corr>` 动词表 `:557-562`；拒收行打印后续跑 `:539-544`；L1 窗口答不了＝`ErrRouteHasNoAllow` 支 `:213-216`）；生产侧流来源＝`main.go:180 interactiveStdin()`。**hosted runner 上没人对这条流打字**＝四枚红的结构因。
- 本机为什么不需要等：①在-allowlist 的写在本机不升级（拼写同形 ⇒ `pathresolver.go:304-305` 记的 `RUNNER~1`↔handle-长形分叉不发生）⇒ 走 L1 窗口/静默放行，无人在场也不需要人；②**每机必出的那两枚 L2 卡（`:337` 红线读、`:513`/`:549` `.env`）本机 sweep 前是真的在等**——等满夹具自己写的 `confirm_timeout_sec = 40`（`run_mode101_test.go:119`），因断言接受"拒绝"文案而绿着（本腿在 pre 树实测：**`dcd5bc825` 形 ManualSwitch `--- PASS (42.76s)`**＝40s 等满＋开销；与 371 证据的 41.48/41.67s、`orphan-reds-cmdwisp-180.md:114-116` 的 40.75/80.77s 同形）。⇒ 票面 :30 的"本机有人/有人造的确认腿"猜测**证伪**：本机没有应答腿，绿的因是"断言容得下超时拒绝"，红的因是"runner 上连在-allowlist 的调用也被升级后等满再拒"。

### AC#2 修法只许测试内可编程应答腿；拒绝形必须同样通过 —— **PASS**

- **只走许可缝**：注入面＝`runSpec.reply`（AGENTS.md §1.3 具名的 `wisp run` CLI 缝）。swept 件把四枚用例的静默 console 流换成脚本流（`run_mode101_test.go:167` `reply = h.script`；调用点 `:329/:397/:467/:541`、`run_test.go:367`），消费方仍是**生产** `runReplyLoop`——⛔ 新 mock 闸、⛔ 碰 `rules_gateway.go`/`assessor.go`/`approval/**`（commit 足迹＝3 枚 `_test.go`，见 §2-AC#3）。
- **拒绝形实测**（AC#2 的判据原句）：`ManualSwitch` boot2 用 `no t101-corr-2` 答红线 R2 卡，断言仍逐字要 `拒绝` 句（`:357-359` 一字未动）⇒ 本腿现跑 **PASS（2.24s/2.26s/2.68s）**；`SessionGrant` 用 `no t101-corr-1` 答 `.env` 卡，`firstErr`/`refused` 断言（`:522`/`:566`）未放宽 ⇒ **PASS（2.90s/2.68s/2.83s）**。**同机前后对照：pre-sweep 42.76s（绿在等满 40s）→ post-sweep 2.24s（绿在 250ms 应答）**——应答腿打的是真队列条目，⛔ 恒真 stub。
- **yes 线不是放宽**：L1 路由上 yes/no 被生产循环按 unknown-corr 拒收（`:539-544`）；窗口本身"只有否决、没有允许"（`:213-216`）⇒ `yes` 永远短路不了 L1 窗口——ComposedGate 的 `blocked >= 2s` 尺照旧成立（本腿现跑 4.33s/3.92s/3.89s PASS，含 ≥2s 断言）。yes 只在"该机形确实升级"的支上兜底，且兜不住时报的是真 tripwire（`cards != 0`），⛔ 超时文案。
- ⛔ 阈值/断言/golden 放宽：`DefaultApprovalTimeout`、`l1_window_sec`、SLO、golden、`thresholds.go` 一字未动（commit 足迹为证）。

### AC#3 反半边（生产真确认腿还在；R20/M4 不被绕过）—— **PASS**

- 足迹尺：`git show --stat abf2a35a6` ＝ **只有 3 枚测试件（+113/−9）**，产码 0 字节 ⇒ `main.go:180`（真 console 输入接 reply 流）、`run.go:1088` 监听接线、`modeConfirm` 的 `confirmModeSwitch` 生产默认适配器全部在位。
- R20/M4 的钉**没被绕也没被换**：swept 未碰的 `TestTicket101ModeSwitchUsesTheRealL2Gate` 本腿同场现跑 **PASS（1.84s/2.16s/2.66s）**——"每次切进 auto_approve 恰一张 L2 卡"（`cards==1`，`:621-622`）、"无人答则档不动"（`:624-625`）、"console run 无 native 通道必拒"（`:617-619`）三条逐字在跑。审计面：`MODE-SWITCH … result=applied`（`:291-294`）、`MODE-READ origin=startup`（`:362`、`:421`）断言全保留。

### AC#4 runner 真读数 —— **未取＝FAIL 态，E 尾交编排者（票面 :49 本就写明子代理不许 push）**

`next=`（照票面格式写全 run/step/期望，读数人回填 id）：
1. **先决**：编排者把 §1 的三枚 helper（`cmd/wisp/ticket123_replyleg_{test,windows_test,other_test}.go`）随显式 pathspec 补落一枚提交，⛔ 不带它们 push（CI 必红在编译面）。
2. push `origin/dev` 后的**首个 run** → job **`test-windows`**（`.github/workflows/ci.yml:810`）→ step **`cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)`**（`ci.yml:886`；票面旧称"第 7 步"，该 job 今日已有 25 枚 `- name:` 步骤 ⇒ **按步名取⛔ 按序号**）。
3. 期望读数：四枚名逐枚 `--- PASS`；该步 `--- FAIL` 计数 **0**；⚠ 票面 "PASS=33" 是 2026-09 的分母，包已长 ⇒ 操作判据＝**名册四枚全绿＋步 rc=0**；取数人须同时落 **run id＋job id＋step 号**（票面 :51：说不出＝这道门不存在；**skipped ≠ 通过**）。
4. **形状 witness 顺手取**：四枚时长应为秒级（`ManualSwitch≈2-3s`、`SessionGrant≈3-4s`、`ComposedGate≈4-6s`）；**只要再出现 300.0x/40.0x 等满时长或 `审批超时` 文案，就是升级支没被拼写尺封住**——这一发同时是 §5 两条"本机测不到"支（8.3 升级支、AC#5 替身支）的唯一真读数。
5. 顺带同 job 取 `TestSyncRedTeamRealOneDrive` 在 runner（无活体根）上的 **stand-in 支 `t.Logf("ticket 123 AC#5 stand-in: …")`** 行（见 AC#5）。

### AC#5 windows 腿未记账 skip —— **PASS（选形＝"永不 skip 的替身"，票面要的两选一已明写落在件里）**

- 选择登记：`syncdirs_redteam_windows_test.go:202-216` 头注逐字写明"skip 已去、无活体根时用注入的 registry-grade 合成根跑同一套拼写电池，并在日志里响亮声明这是替身⛔ 机器现实"＋"选替身⛔ 豁免清单＝本票 AC#5 的记录裁"。
- 代码核（现树 `:217-252`）：两处 `t.Skip` 已删（原 :207/:220）；无 profile home ⇒ `t.TempDir()` 顶＋响亮 Logf（`:223-224`）；无活体根 ⇒ 合成根注入（`:242-244`）＋ `SyncDetectionComplete()` 自锁（`:245-247`，注入根必须是确认检测、⛔ suspect-fallback 蒙绿）＋电池断言与根来源无关（`:270-290` 同一套拼写/落字节尺两条支共用）。
- **fresh 读数**：`go test ./internal/risk -run Sync -count=1 -v` ⇒ **rc=0；顶层 RUN=33／PASS=32／FAIL=0／SKIP=1；`--- SKIP: TestSyncRedTeamRealOneDrive` 计数＝0，该枚 `--- PASS (0.01s)`**。本机有活体根（日志逐字 `P12 evidence: live OneDrive root via env: C:\Users\swq\OneDrive`）⇒ **本机走的是活体支；替身支未在本机点亮**（诚实登记，CI 补读＝AC#4 第 5 条）。
- 边界具名（⛔ 本票格、给编排者入账）：①同一枚文件 `:50-57` 的 `shortPathOf` 仍有 `t.Skipf("no 8.3 alias …")`——那是 **307 拼写电池的前置尺**，⛔ 票面 :53 具名的 `TestSyncRedTeamRealOneDrive`；②同场 SKIP 那 1 枚＝`TestSyncRegistryProbeLive`（自带文案"HKCU Accounts without UserFolder is the documented reality here"，既有的环境前提 skip，⛔ 未记账新债）。"windows 腿归零 skip"这件事按票面口径（具名那枚）成立，按整包口径仍差这两枚具名既有件。

### AC#6 门禁 —— **PASS（本腿可取部分全绿；两枚具名例外）**

| 尺 | 命令（原样可复跑） | 读数 |
|---|---|---|
| cmd/wisp 定向 count=1 | `go test ./cmd/wisp -run "TestTicket101\|TestComposedGateBlocksAWriteForTwoSeconds" -count=1 -v -timeout 420s` | **rc=0**；顶层 **RUN=6／PASS=6／FAIL=0／SKIP=0**，`ok 21.029s`；四枚票 123 名逐名 PASS（2.24/8.04/2.90/4.33s） |
| cmd/wisp 定向 count=2 | 同上 `-count=2` | **rc=0**；**RUN=12／PASS=12／FAIL=0／SKIP=0**，`ok 43.902s`；两回路时长同形（2.26/2.68、7.81/8.03、2.68/2.83、3.92/3.89…） |
| internal/risk Sync count=1 | `go test ./internal/risk -run Sync -count=1 -v` | **rc=0**；顶层 PASS=32/FAIL=0/SKIP=1（见 AC#5 两枚具名）；含子测试 47 PASS 行 |
| internal/risk Sync count=2 | 同上 `-count=2` | **rc=1**——**归因见 §3 name-diff，⛔ 本票** |
| gofmt/gofumpt | `gofmt -l`／`/d/work/base/gopath/bin/gofumpt.exe -l` 于全部六件（3 swept＋3 helper） | **空列表、rc=0**（repo 级 `gofmt -l .` 报的全是 `.scratch/**`、`build/**` 未跟踪探针件，⛔ tracked 面） |
| go vet | `go vet ./cmd/wisp ./internal/risk` | **rc=0**（含三枚未跟踪 helper 的工作树合成形） |
| d22scan | `sh scripts/d22scan.sh` | **rc=0 clean**；自检正控 `runtests.sh: OK PASS=35 FAIL=0 SKIP=0`；scope 现量：bans#1-5 internal/=**363**、cmd/=**58**，ban#6 frontend/=116，ban#7 internal/tools/=**48**，ban#8 design/=**39**、frontend/=116、internal/=**907**、cmd/=**216** |
| 台账各 scope 不降 | 对 `pending-and-issues.md` 末记（`ban #8 internal/=623/:15359 段 525/531`、`cmd/=119`、`tools/=25`、`frontend/=85`、`design/=39`；`bans#1-5 internal/=229、cmd/=39`） | **每枚 ≥ 末记**（907≥623、216≥119、48≥25、116≥85、39=39、363≥229、58≥39）⇒ 不降成立 |

## 3. name-diff：本腿采到的两枚红都做了归属取证（第 21 处口径）

1. **`risk -run Sync -count=2` 二回路红** ＝ `TestSync307AliasHomeUnderProfileStaysSuspect`／`TestSync307AliasShortRootMembership` 各第 2 回路红（第 1 回路绿），红句 `ticket 307 fail-open … home="C:\\Users\\swq\\TESTSY~1\\001\\THIRTY~1"`。归属尺＝**pre-sweep 导出树（`dcd5bc825`，⛔ swept 改动、⛔ 三枚 helper）同尺同形复现**：`--- PASS / --- PASS / --- FAIL / --- FAIL`。两枚件＋`syncdirs.go` 均⛔ 在 `abf2a35a6` 足迹里 ⇒ **既有的 307 仪器 count=2 隔离红（t.TempDir 清理后 8.3 别名 `TESTSY~1/THIRTY~1` 复用面），归 307 家族账，⛔ 票 123**。（机理线索：回路间 `tryShortPath` 拿到与回路 1 同名的复用短形，判定链答 not-sync。）
2. **`abf2a35a6` 单独编译红**（§1）＝ swept 足迹不完整，非工作树在飞件拖累本腿——反向钉死：本腿没碰任何未跟踪件，是本票交付**需要**它们。
3. 整包 `cmd/wisp` rc=1 窗口族：按派单⛔ 判、⛔ 跑；本表全部读数 `-run` 定向＋逐名。

## 4. closeable? —— **no**（两点，缺一即不结）

- **AC#4 未取**：这道门只有 runner 能答（§2-AC#4 已把 run/step/期望写成 next=）。
- **交付不完整**：三枚 helper 未落库（§1）。⛔ 在二者落定前 push 后按本表取 CI 读数。

**safe-to-flip？**
- 现在：**⛔ 翻任何一格**（合成形才刚绿，交付面还缺 §1 那一步）。
- helper 补落＋push＋CI 读回（四枚逐名 PASS、步 FAIL=0、时长秒级、stand-in 行在案）之后：**AC#1/AC#2/AC#3/AC#5/AC#6 五格可翻**（AC#5 翻格时在票面注一行"替身支首读＝CI 那发"）；**AC#4 随读数翻**；四格齐前票⛔ `-done`。
- 撤销口令（给编排者一句话复原本表）：回退 `docs/evidence/s1/123-verdict.md` 一笔即可，本表⛔ 附带任何代码/台账/票面动作。

## 5. 边界与限制（宁缺毋造）

1. **8.3 升级支本机未测到**： planted 实验（`TEMP/TMP=C:\Users\swq\tmp\123ALI~1` 真 8.3 祖先形，两树同尺）——pre-sweep ManualSwitch **42.76s PASS**、post-sweep **2.94s/5.13s PASS**：本卷的判定链没把在-allowlist 调用升级出去（`pathresolver.go:304-305` 的分叉是 runner 的 USERPROFILE↔handle 形状，本机长短同串）。⇒ 拼写尺（`t123dirSpellingForms`）在本机是**惰性支**，与 `_windows_test.go:12-14` 自陈"belt and braces ⛔ measured repair"一致；它的真读数排在 AC#4 那发（判据＝时长里再无 300.0x/40.0x 等满形）。
2. **AC#5 替身支本机未点亮**（本机有活体 OneDrive 根＝走的活体支），同排 AC#4 第 5 条。
3. 实验树里 `third_party/`、`build/` 属 gitignore（clone 后⛔ DLL），本腿从工作树补拷三枚 `.dll` 才起跑——首两发 `0xc0000135` 即此因，⛔ 任何测试面结论。
4. `count=2` 的 307 红＝既有账（§3）；本票格⛔ 因此翻，编排者入账时⛔ 记到 123 头上。
5. 本表绿＝`abf2a35a6`＋三枚未跟踪 helper 的合成形；helper 的内容裁（零 Test 函数、纯缝件）见 §1，其格式/扫描面已并入 §2-AC#6 各尺。
