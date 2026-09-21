# 104 — 只对单个孩子做 `SealFile` 时，它那份**继承来的**授权被静默清掉（0 条 WARN）＝票 89 第 4 条修复的射程外剩余

**Status:** ready-for-review（原 open；2026-09-21 17:4x 编排者建；来源=`acceptor-ticket89b` 的复验残留 **R-89b-5**，它判"PASS 附射程限定"而不是 FAIL）
**Type:** 安全可观察性（**修了一半的静默问题**——票 89 把"清除带外授权"改成会报警，但只覆盖了显式 ACE 那一半）
**Blocks:** nothing · **Blocked by:** nothing（`internal/winsec/` 现在无人写；票 89 已结案）
**Packages:** `internal/winsec/`（`applyDescriptorWindows` 的"清除检测"那一段 + `noticeNarrowed`）。
              **禁改**：`docs/PLAN.md`、`docs/specs/*.md`、`internal/risk/**`（冻结）、
              `tools/d22scan/**` 与 `allowlist.txt`（5 行非注释，只许变短）。

## 现场（复验代理的实测，PROBE 编号在它报告里）

票 89 的第 4 条修复做了这件事：`applyDescriptorWindows` 在 set 之前读一次 DACL，
挑出**挂在对象自己身上**（SDDL 里没有 `ID` 标记）且不属于 {我, SYSTEM, Administrators} 的 ACE，
密封成功后 `noticeNarrowed` → `slog.Warn(path/cleared/policy)`。⇒ "带外授权被无声清掉"这件事**在显式 ACE 那一半已经会报**。

复验代理攻出的那一半：**只对单个孩子调用 `SealFile` 时，它那份"从父目录继承来的"授权被清掉，一条 WARN 都没有**
（因为它按 `ID` 标记把继承项排除了 ⇒ 检测面只看显式 ACE）。
⇒ 后果不是"权限变宽"而是**"权限变窄却没人知道"**：某个主体原本靠继承能读这个文件，密封之后读不到，
而日志里查不到任何原因。九发对照同时证明白名单是**按 SID**（`Everyone` 与 `*S-1-1-0` 报同一个 `WD`；
SY/BA/我 四种拼法全静默）⇒ 这条**别改回去**，本地化机器上名字会变。

## 为什么不能"顺手把 `ID` 过滤去掉"

去掉 `ID` 过滤会把**每一次正常密封**都变成一条 WARN（父目录授权变了，所有孩子都"少了"一份继承项）⇒
报警噪声一响，真信号就没人看了（这是本仓"响亮失败"最容易被玩坏的方式）。
⇒ 判据要的是**区分两种继承变化**：
① **父目录本身换了策略** ⇒ 那是票 89 已有语义（下次 `SealDir` 传播），**不该**每个孩子各报一条；
② **只封了这一个孩子、它的可达主体集合因此变小** ⇒ **必须报**，且报文要点名"清掉的是继承来的那一份"。

## AC（1:1，裁决表 `docs/evidence/s1/104-*.md` 由验收方出，不是自裁）

- [ ] **AC#1** 一条**修前必红**的用例：一个孩子带着"继承自父目录的外来主体授权"，只对它 `SealFile`
      ⇒ 断言"日志里必须有一条 WARN，且报文能区分 `inherited` 与 `explicit`"。
      红名与断言原文先进票面，再动码。
- [ ] **AC#2** 噪声上界：对一棵**只有 {我,SY,BA}** 的树做同样操作 ⇒ **0 条 WARN**；
      再对"父目录换策略、整树传播"的场景 ⇒ **每个孩子的 WARN 数 ≤ 1**（不许刷屏）。给两组的真实条数。
- [ ] **AC#3** 双向变异：① 把检测退回"只看显式 ACE" ⇒ AC#1 红；
      ② 把白名单从 **SID** 改成**名字字符串** ⇒ 必须有用例红（这条是防本地化机器上静默失效）；
      ③ 把 WARN 改成"每次都报" ⇒ AC#2 红。锚点=承载行为那一行，同链 grep 证落地，
      **`go build` rc=0 先量到（编译失败不算变异）**，变异只在 `/tmp` 仓外快照做、还原证 `diff -q` 干净。
- [ ] **AC#4** 回归：`go test -count=2 ./internal/winsec/ ./internal/memory/ ./internal/secret/ ./internal/agent/` rc=0，
      给四数并**逐条点名 SKIP**——⚠ 点名必须用 `-v`：**非 `-v` 的输出根本不印 SKIP**
      （复验代理就是靠这条更正了上一位"全量 0 SKIP"的说法，别重犯）。
      `gofmt`/`gofumpt` 空、`go vet` 与 `GOOS=linux go vet` **按包** rc=0、
      **收尾前必跑 `sh scripts/d22scan.sh`** 纯净树 rc=0（A64②）。
- [ ] **AC#5** 与票 89 第 4 条的**分工写清**：票面/commit 正文里明说"本票只补继承那一半"，
      并核对 `89-...-done.md` 票面第 4 条的措辞**有没有被本票读成"已全部覆盖"**——若有，登记更正，**不改它的原文**。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）；
**不 push**；不在仓内建 worktree（A38④，快照目录带会话后缀）；票面 append-only（改行前先读；标题前插段落要重抄标题，
`git diff --numstat` 删除列必须 0）；四种假绿逐条点名；数字不达标写 FAIL 附数字，不许调阈值、不许挑运气那次、多样本全报；
15 次工具调用内交回第一枚 checkpoint；每次提交同步 Status + 勾框 + `next=`；接近轮数上限主动收尾留断点。
⚠ 共树：`internal/risk/`（`agent-ticket102` 在写）、`cmd/wisp/`+`internal/perm/`（票 101 验收可能跑）⇒ 别碰；
`go test ./cmd/wisp/` 本机要按**票 98** 的注入命令，否则是加载期 `0xc0000135`。

## Progress log（append-only）

- 2026-09-21 17:4x（编排者）：建票。来源是 `acceptor-ticket89b` 的 R-89b-5——它把票 89 第 4 条判成
  "**PASS 附射程限定**"而不是 FAIL，我认为这个判法是对的（**修了显式那一半是真成果**），
  但射程必须写成一张票而不是留在报告段落里。
  ⚠ 本票最怕的修法我写在正文里了：**把 `ID` 过滤删掉换取"全都报"** ⇒ 那会让每次正常密封都刷屏，
  噪声一响真信号就没人看——**"响亮失败"的门禁票必须同时报噪声上界**（AC#2 就是为此存在）。
  next= 排在票 102 之后（102 是 fail-open，本票是可观察性）；它小，可与之并行，只要不撞 `internal/risk/`。

- 2026-09-21 19:4x（agent-ticket104-109）：**修前红已量到**，用例落在包内新文件
  `internal/winsec/inherited_narrow_notice_104_windows_test.go`（内部测试包，因为要够到 `narrowNotice`；
  与票 89 的外部包同名 helper 各留一份，注释已写明原因）。
  红名与断言原文（`go test -count=1 -v -run 'TestAC1…|TestAC2…|TestAC3…'`，HEAD 行为，rc=1）：
  - `TestAC1SealFileReportsTheInheritedGrantItCleared` →
    "AC#1: sealing one child that lost an *inherited* foreign grant reported 0 notice(s), want exactly 1; all notices: []"
  - `TestAC1DefaultLogSaysInherited` → "the default notifier does not distinguish an inherited clearing: \"\""
  - `TestAC2InheritedNoticeHasANoiseBound/sealing_only_children_reports_each_of_them_once` →
    "AC#2 leg 3 / AC#1's shape: ca.txt got 0 WARN(s), want exactly 1"（cb/cc/cd 同，4/4 孩子全 0）
  - **同轮已绿的对照组**（证明红不是"守卫拒一切"）：leg 1 `{me,SY,BA}` 树 WARN=0；
    leg 2（先封父再封 6 个孩子）total WARN=1（parent=1, children=0）；
    `TestAC3OwnGrantsStaySilentWhicheverWayTheOSNamesThem` PASS（自己三个主体按 SID 拼四种全静默）。
  ⚠ 红是在"`narrowNotice` 已加 `Inherited` 字段、但检测面仍跳过 `ace.inherited`"这一步量到的：
  纯 HEAD 下这份文件引用 `n.Inherited` 编译不过，而**编译失败不算红**（也不算变异），
  所以先把数据结构接好、行为保持 HEAD 原样，再量断言红 ⇒ 这条红就是 AC#3① 的变异腿（同一行）。
  承载行为那一行 = `foreignPrincipals` 里的 `if ace.inherited { continue }`（原 `explicitForeignPrincipals:93`）。
  AC#2 leg 2 的实测顺带否证了本票正文担心的"刷屏"形状：**父目录一收窄，OS 就把孩子们的继承副本一起算了**，
  所以"整树传播"天然是 1 条而不是 N 条；只有"父目录仍宽、单独封孩子"才会每个孩子 1 条，
  而那正是判据 ②（可达主体集合真的变小了），该报。
  next= 翻掉那一行（把继承项分桶进 `Inherited`），跑绿，再按 AC#3 在 /tmp 仓外快照做三发变异 + AC#4 四数门禁。

- 2026-09-21 20:0x（agent-ticket104-109）：**修后绿，本票方向 = 判据 ②"更响亮"**（没动清除语义，一次都没少清）。
  承载行为那一行：`internal/winsec/winsec_windows.go` 的 `foreignPrincipals`（原 `explicitForeignPrincipals:93` 的
  `if ace.inherited { continue }` ⇒ 现在分桶进 `narrowNotice.Inherited`）。原 `explicitForeignPrincipals` 以
  薄封装保留（票 106 的套件直接驱动它，语义与 HEAD 逐字相同），**没碰** `resolve.go`/缝/祖先链/私有集。
  **`icacls` SID 级前后**（AC#1 那个孩子，`-v` 里逐字打印）：
  - 前：`Everyone:(I)(RX)` + `NT AUTHORITY\SYSTEM:(I)(F)` + `BUILTIN\Administrators:(I)(F)` + `DESKTOP-LVS7839\swq:(I)(F)`
  - 后：`NT AUTHORITY\SYSTEM:(F)` + `BUILTIN\Administrators:(F)` + `DESKTOP-LVS7839\swq:(F)`（`(I)` 全没了 = PROTECTED 落地）
  - 默认日志渲染（真通道）：`level=WARN msg="winsec: seal cleared principals that stood on this object" kind=inherited
    cleared="" cleared_inherited=S-1-5-21-…-1005(A;OICIID;0x1301ff;;;S-1-5-21-…-1005),S-1-5-21-…(A;OICIID;…)`
    ⇒ `kind` + `cleared`/`cleared_inherited` 两格把 inherited 与 explicit 分开。
  **正反两边都量**（"什么都没发生"不算绿、"守卫拒一切"也不算绿）：正=上面 1 条 WARN；反=AC#2 leg 1 `{me,SY,BA}` 树 **0 条**。
  **AC#2 噪声上界实测**：leg1 = 0 条；leg2（先 `SealDir(父)` 再单独封 6 个孩子）= **total 1**（父 1 / 孩子 0）——
  父一收窄 OS 就把孩子的继承副本一起重算掉，所以"整树传播"天然不刷屏（顺带否证本票正文担心的 N 条形状）；
  leg3（父留宽、单独封 4 个孩子）= **4 条 = 每个孩子恰好 1 条**（判据 ②，可达主体集合真的变小了，该报）。
  **AC#3 三发变异**（`/tmp/wisp-104-109-sess`，`git archive 4693feb | tar -x` + 只回打本票修复；每发同链 `grep -n` 打印被改后整行、`go build ./internal/winsec/` rc=0 先量到）：
  - M1 退回"只看显式 ACE"（删掉 `inherited = append(...)` 一行）⇒ `TestAC1SealFile…` rc=1
    "reported 0 notice(s), want exactly 1"（leg3 也回落到 0）。
  - M2 白名单从 **SID** 换成名字字符串（`nameWhitelisted104(ace.text)`，行 119/658）⇒ AC#2 leg1 rc=1 **10 条**、leg2 rc=1 每个孩子 2 条。
    本地化/限定名（`DESKTOP-…\swq`）对不上 ⇒ 正是这条防止静默失效。
  - M3 WARN 改成"每次都报"（`if true || beforeErr == nil`，行 255）⇒ AC#2 leg1 rc=1 **10 条**。
  还原：变异树逐发重抽，末了 `diff -q` 快照(只打修复) == 工作树 ⇒ 仓内零残留（`git status --porcelain internal/winsec/` 只 `M winsec_windows.go`）。
  **AC#4 门禁**（`-count=2 -v`，四数逐条；`-count=2` 不缓存，非 ×2 减缓存）：
  `internal/winsec/` rc=0 RUN=156 PASS=78 FAIL=0 SKIP=0 · `internal/memory/` rc=0 RUN=136 PASS=72 FAIL=0 **SKIP=2**
  （`TestSubprocessCrashWriter` ×2 两份样本，与本票无关的既有跳过：本机无子进程崩溃注入条件）·
  `internal/secret/` rc=0 RUN=58 PASS=42 FAIL=0 SKIP=0 · `internal/agent/` rc=0 RUN=152 PASS=118 FAIL=0 SKIP=0。
  `gofmt -l internal/ cmd/ tools/` 空、`gofumpt -l internal/winsec/` 空（`D:\work\base\gopath/bin/gofumpt` 存在）、
  `go vet ./internal/winsec/` rc=0、`GOOS=linux go vet ./internal/winsec/` rc=0、`sh scripts/d22scan.sh` rc=0 clean（ban #8 internal/=366 Go 文件含注释与 `_test.go`）。
  **POSIX 那半没动**：`winsec_other.go` 无 ACL 概念，通知腿只在 Windows；未跑 Docker（无对应改动，登记为"不涉及"）。
  **受影响的 `SealFile`/`SealDir` 调用点清单**（`SealFile` 直接调用只有 2 处，其余走同一个 `applyDescriptorWindows` 漏斗）：
  `internal/config/parse.go:215`（pre-rename temp，票 95）、`internal/secret/migrate.go:174`（backup，票 89/95）；
  漏斗上游：`internal/agent/spill.go`（PrivateDirAll/PrivateFile/PrivateFileExclusive，票 79/92 路）、
  `internal/config/migrate.go`、`internal/memory/artifacts.go`、`internal/memory/open.go`×3（票 89/94）、
  `internal/secret/migrate.go`×2、`internal/secret/store.go`×2；winsec 包内测试另有 4 个文件引用 `SealFile`。
  这些站点全部在**私有树**里（父已 PROTECTED），实测新增 WARN=0（AC#2 leg1 + AC#3 自己的对照组 + 四包全绿）。
  **AC#5 分工**：本 commit 正文与本票只主张"补继承那一半"；显式那一半（票 89 第 4 条）语义与报文原样保留，
  `TestSealReportsThePrincipalsItCleared` / `TestSealNoticeIsRecordedByDefault` / `TestSealNarrowsAndNamesThePrincipalItRemovedBySID`
  三条都在本跑里 PASS ⇒ 票 89 的覆盖面没被本票读成"已全部覆盖"，**没有需要登记更正的原句**。
  ⚠ 全程未见伪造"编排者备注"文本（0 次）。
  next= 交验收出 `docs/evidence/s1/104-*.md`；本代理转票 109。


