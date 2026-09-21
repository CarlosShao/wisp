# 89 — `0o600` 在 Windows 上是装饰品：artifact / 密钥 / 数据库**从来没真的"只有我能读"过**（A51①②）

**Status:** **accepted-done（复验通过）**（2026-09-21 17:3x 第二轮独立对抗复验判 **退回单四条全部 PASS**，
复验表 `docs/evidence/s1/89b-reacceptance.md`；第一轮表 `89-adversarial-acceptance.md`。
独立复现要点：删 `propagatePrivate` ⇒ rc=1、`=== RUN` 25、**14 绿 1 红、红名唯一**，
红因原文 `NOT PRIVATE operator-widened-file.txt: foreign SID(s) S-1-1-0`（在 HEAD 上同锚点红两条，更强）；
迁移备份同一份 `.bak-plaintext` BEFORE `Everyone:(I)(RX)` → AFTER 只剩 SY/BA/我；
符号链接那一格**真测掉了**（未提权 + 开发者模式=1 ⇒ `os.Symlink` 成功、`os.Remove` 能拆、目标存活）；
带外授权被清除时**不再是静默的**，且白名单**按 SID 不按名字**（九发对照）。
四包回归 402 RUN=2×201、0 FAIL。⚠ 复验同时**更正了修复者一句「全量 0 SKIP」**：实测有 **2 行 SKIP**
（`TestSubprocessCrashWriter` ×2，票前既有），而**非 `-v` 的输出根本不印 SKIP** ⇒ 「0 SKIP」这种说法必须配 `-v` 才成立。
⚠ 三笔随改名落的账：**R-89b-5 立案票 104**（只对单个孩子 `SealFile` 时，它那份**继承来的**授权被静默清掉、0 条 WARN）、
`internal/config` 的同类明文站点与它的 `.bak-<schema>` **仍敞**（归票 95）、
以及 **`MigratePlaintext` 至今无生产调用方**（⇒ 这条密封今天是「能力就绪、没人叫它」，票 95 接线时一并处理）。
**票 94 的牙没被卸**：未动 `winsec.go`/`resolve.go`，铸造口未被绕开，攻击性输入仍被拒。）
`agent-ticket89b` 已在 `c8d5c94`（第 1、2 条判据+接线）/`01e7007`（第 3、4 条 + 前三次变异读数）/本枚（第 4 条变异读数 + 全量门禁）
**把四条全部补完，各带实测红名** ⇒ **等重验收，仍不改名 `-done`**。裁决表 `docs/evidence/s1/89-adversarial-acceptance.md`。
票面那两句错话（"普通权限建不出目录符号链接"、"第三个主体会在写入点直接失败"）**已按实测改掉**，见 Progress log 17:5x/18:1x。
~~blockers：`internal/winsec/` 此刻有 `agent-ticket94` 在写~~ ⇒ 已解：94 的 `ResolvedPath` 铸造口与 `internal/risk/winsec_c26.go`
形状**原样保留**，本枚未新增任何绕开它的入口。）
**Type:** 安全（本机数据泄露面：同机其他用户/进程可读我们的私有数据）
**Blocks:** nothing · **Blocked by:** nothing（`internal/memory/artifacts.go` 与 `internal/agent/spill.go` 此刻无人写）
**Packages:** 新建一个 Windows ACL 的小工具文件（建议 `internal/acl/` 或 `internal/winsec/`，**由你定，但要在票面写理由**）
+ `internal/memory/`、`internal/agent/` 的落盘路径。**禁改**：`docs/PLAN.md`、`docs/specs/*.md`、
`internal/risk/**`、`rules_gateway.go`、`tools/d22scan/**`、`allowlist.txt`、`frontend/**`（票 77 在写）、
`internal/agent/approval/**`（票 87 刚交回，别再动）、`internal/panel/**`、`cmd/wisp/**`（票 77 在写）。

## 实测事实（票 79 的代理发现、编排者登记为 A51①②）

- `os.OpenFile(path, …, 0o600)` 在 Windows 上**不落地**：artifact 落盘实测是 **`-rw-rw-rw-`**，
  权限由**父目录的 ACL 继承**。⇒ **"工件只有当前用户可读"这句话从来没有成立过**（旧代码同样如此，不是回归）。
- 同族：`0o600` 出现在 `internal/memory`（artifacts / blobs）与 `internal/agent/spill.go`；
  **`secrets\`（DPAPI blob 目录）、`staging\`、`wisp.db`（SQLite + `-wal`/`-shm`）是不是同病，要一起查**——
  那三处比 artifact 更机密。
- A51②：`os.Remove` **删不掉"指向目录的符号链接"** ⇒ 票 79 的游离子树回收路径碰到这种条目会清不掉。
  （它的"只用 `os.Remove`、绝不 `RemoveAll`"的选择是对的——递归删一个可能是链接的东西等于把删除半径交给别人。
  **不要为了删除成功而改成 `RemoveAll`**，那会开一个更大的洞。）

## 判据要先于实现（这票最容易"写了 API 但没测到东西"）

- [x] **AC#1** 先给**基线证据**：真机上把这四类各建一个文件（artifact / secret blob / `wisp.db` / staging 临时件），
      用 `icacls <path>` 打印**实际 ACL**，逐条贴原文。⚠ **不许用 Go 的 `FileInfo.Mode()` 当证据**——
      它给的就是那套被 Windows 忽略的假位。
- [x] **AC#2** 定**方向**（写进票面并说明为什么）：给**文件**逐个设 DACL，还是给 **data 根目录**设一次、
      子项继承？后者覆盖面大但要求目录已存在且早于任何写入 ⇒ **建目录与建文件的先后**必须有用例钉住。
      两条路都要保证：**继承不能把权限放大**（拿 `icacls /inheritance` 之类的显式断言，别靠默认行为）。
- [x] **AC#3** 实现 + **真红过的判据**：在**当前用户之外还要有第二个账户或可读探针**才算证明。
      拿不到第二账户时的**可接受替代**：断言 ACL 里只有当前 SID 有读写、且 `SYSTEM`/`Administrators` 之外的
      任何 SID 不出现；并跑一次 `icacls` 原文比对。**"我把 `0o600` 换成了调用 SetSecurityDescriptor"不算证明**——
      那只是写了代码，没测到东西。
- [x] **AC#4** A51② 那条单独钉：造一个"artifact 位置是指向目录的符号链接"，
      断言 ①**不被递归删除**（链接目标目录里的文件必须还在），②给出**具名错误**而不是静默失败。
      ⚠ 若 Windows 上普通权限无法建符号链接，**就如实写"本机无法构造"并给出 CI 侧（Linux）等价构造**
      + 说明两个平台上这条判据各自的强度。**不许**因为"造不出来"就把这条判据删掉。
- [x] **AC#5** 失败方向必须是**收紧**：ACL 设不上时报错并**拒绝落盘**（或退到"不落敏感数据"），
      **绝不允许**"设不上就算了、继续以宽权限写"。判据用例：模拟一次失败注入 ⇒ 必须是错误。
- [x] **AC#6** 门禁（只跑你碰的包）：`gofmt -l` 空、`gofumpt -l <files>` 空、`go vet <pkgs>` rc=0、
      `GOOS=linux go vet <pkgs>` **按包作用域** rc=0（⚠ 别用 `GOOS=linux go vet ./...`，那条在 Windows 主机上
      因 CGO=0 排除 sherpa 而永远 rc=1，A54③）、`go test -count=2 <pkgs>` rc=0 且逐条点名 `--- SKIP`/`--- FAIL`。
      ⚠ **非 Windows 平台上这条能力应当"显式不适用"而不是"静默通过"**：POSIX 侧要么测真实的 `0600`，
      要么打一条明说"此平台由 ACL 保证"的 skip——**skip 要逐条点名，不许当 ok**（票 79 交回时就因为
      "两条 SKIP 其实是同一个测试名跑两遍"这种细节差点被误读）。

## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；**每完成一组就把结论追加进票面**（别攒——本仓已有代理死在轮数上限）；
`git commit -q -F - -- <显式路径>` + **带引号** heredoc；禁 `git add -A`/`.`；commit 前核对
`git diff --cached --name-only`，**别把活着代理的票面 add 进去**；
禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`；**不 push**；不在仓内建 worktree（A38④）；
纯净树 `git archive HEAD | tar -x -C /tmp/<带你会话后缀的目录>`；
四种假绿逐跑点名；变异：锚点=承载行为的那一行、同链 grep 自证、还原后 `git diff --quiet`，**编译失败不算变异**；
数红/绿用全量输出仪器；票面 append-only，**要改的那行先读再替换**。

## Progress log（append-only）

- 2026-09-21 18:1x（`agent-ticket89b`，**收尾枚：第 4 条的变异读数 + 全量门禁 + d22scan**）：
  **变异读数 #4（证明第 4 条那条通知是承重的，不是装饰）**：`git archive 01e7007 | tar -x -C /tmp/wisp89b-nt`，
  把 `applyDescriptorWindows` 里的 `noticeNarrowed(narrowNotice{…})` 换成 `_ = before // MUTATION-89SILENT`
  （同链 `grep -n` 命中 144 行、`go build ./internal/winsec/` rc=0 ⇒ 可编译且行为改变）⇒
  `go test -count=1 -v -run 'TestSeal' ./internal/winsec/` **rc=1，`=== RUN` 2，两条全红**：
  `--- FAIL: TestSealReportsThePrincipalsItCleared (0.23s)`、`--- FAIL: TestSealNoticeIsRecordedByDefault (0.07s)`
  ⇒ "重新变回静默清除"这件事只有这两格会响，红名与判据一一对上（`-run` 仅用于变异读数，条数按 `=== RUN`=2 点名）。
  **对票 94 形状的尊重**：本枚**没有新增任何绕开 `ResolvedPath` 铸造口的入口**——第 2 条走的是既有公开 API
  （`winsec.PrivateFile` / `winsec.SealFile`，两者都在内部过 `resolveString`），第 4 条只动 `applyDescriptorWindows`
  这一个已存在的平台内部接缝；`RemoveUnlinked` 刻意不解析路径那一条**未碰**（票 94 自认弱处，验收方在判，且要动 `internal/memory`）。
  **最终门禁（工作树 = 本枚）**：`gofmt -l internal/winsec internal/secret internal/memory internal/agent` 空；
  `gofumpt -l .`（**全仓**）空；`go vet ./internal/winsec/ ./internal/secret/ ./internal/memory/ ./internal/agent/` rc=0；
  `GOOS=linux go vet` 同四包（按包作用域，A54③）rc=0；
  **`go test -count=2` 四包全 `ok` rc=0**：winsec 25.132s / secret 1.533s / memory 34.235s / agent 6.958s，
  全量输出里 `--- SKIP` **0 条**、`--- FAIL` **0 条**（用 `grep -E "SKIP|FAIL"` 扫全量输出计数，不是 `head`）。
  **`sh scripts/d22scan.sh` 在纯净树上跑**（`git archive HEAD | tar -x -C /tmp/wisp89b-d22`）：
  `d22scan: clean - no D22 ban violations`，`tools/d22scan` 的 positive control `ok 6.891s`，
  生产文件计数 bans#1-5 internal/=197 / cmd/=20、ban#6 frontend/=37、ban#7 internal/tools/=17、ban#8 internal/=337、cmd/=25；
  **`allowlist.txt` 未改**（仍是那 5 行非注释，一条没加）。
  **对 owner 的口径可以升格到哪一步（诚实版）**：退回单第 1、2 条补完之后，
  "只有当前用户可读"现在点名覆盖 **artifacts / DPAPI blob / db+wal+shm / staging / `config.toml.bak-plaintext`（含既有备份的 repair 腿）
  / 迁移临时件（rename 后落在 `config.toml` 上）**，且"既存带显式 ACE 的老树"这一格现在有用例咬（含"删掉 walk 必红"）；
  **仍然不覆盖**：`internal/config` 与 `internal/models` 的写路径（票 95/90 的账）、日志（票 95，裁定只能走 `SealFile`）。
  本机 `CodexSandboxUsers` + 一个解析不出名字的别名 SID 对未密封的落点仍持有 `(M,DC)`——这条基线没变，
  变的是"密封过的对象上它们不存在"这一判据现在覆盖到了迁移备份。
  next= **交编排者重验收**（四条各带读数与红名，见 17:2x / 17:5x / 本枚三条）。请顺带处理两件我不做的事：
  ① `docs/reports/pending-and-issues.md` 两条登记更正——A51② 在 **junction 与目录符号链接两类对象上都不复现**
  （本机实测），以及新语义"**winsec 拥有 data 树的唯一授权权：带外授权会在下一次密封时被清除，并且打一条
  `level=WARN … cleared=<trustee>(<ace>)` 的可 grep 记录**"（验收 §8 要求登记的那条，registry 是编排者的面）；
  ② 若验收要"第二账户真的读不到"的强证据，仍需一台能 `runas` 的机器（本箱拿不到那两个沙箱账户口令）。


- 2026-09-21 17:5x（`agent-ticket89b`，**退回单第 3、4 条落地 + 上一枚 `c8d5c94` 欠的三次变异读数**）：
  **变异读数 #1（第 1 条要求的"删掉 `propagatePrivate` 必红"）**——在仓外快照
  `git archive c8d5c94 | tar -x -C /tmp/wisp89b-walk` 里把 `winsec_windows.go` 的
  `return propagatePrivate(path)` 改成 `return nil // MUTATION-89WALK`（`applyDescriptor` 保留）：
  先证**可编译**（`go build ./...` rc=0）再在同一条链里 `grep -n MUTATION-89WALK` 命中 293 行，然后
  `go test -count=1 -v ./internal/winsec/` ⇒ **rc=1，`=== RUN` 25、顶层 14 PASS / 1 FAIL、0 SKIP**，
  **红名唯一**：`--- FAIL: TestAC2SealDirNarrowsChildrenThatCarryTheirOwnExplicitACEs (1.06s)`，
  红因是 OS 自己的原文（三个带显式 ACE 的子项全都还是宽的）：
  `icacls text of operator-widened-file.txt still names the foreign principal (Everyone)` /
  `NOT PRIVATE operator-widened-file.txt: foreign SID(s) S-1-1-0 (principals: [Everyone NT AUTHORITY\SYSTEM BUILTIN\Administrators DESKTOP-LVS7839\swq])`，
  `operator-widened-dir` 与 `inside.txt` 同形 ⇒ **变异只打到该打的那一格**，其余 14 条（含两条 AC#1 反向钉子、AC#4/AC#5、
  生产端到端）保持绿，正是验收报告 §6(b) 预测的形状。快照在仓外，工作树未做变异。
  **变异读数 #2（第 2 条）**：`/tmp/wisp89b-mig`（同一 SHA）把 `migrate.go` 两处退回
  `os.WriteFile(..., 0o600) /*MUTATION-89MIG*/`（`grep -c` = 2、`go build ./internal/secret/` rc=0）⇒
  `--- FAIL: TestAC2MigrationBackupIsPrivate (2.58s)`（`=== RUN` 1），红因就在备份文件自己身上：
  `NOT PRIVATE config.toml.bak-plaintext: foreign SID(s) S-1-1-0 (principals: [Everyone BUILTIN\Administrators NT AUTHORITY\SYSTEM DESKTOP-LVS7839\swq])`
  （`-run` 过滤只用于这次变异读数，条数已按 `=== RUN`=1 点名）。
  **变异读数 #3（我加的 repair 腿）**：`/tmp/wisp89b-rep` 把 `} else if serr := winsec.SealFile(backupPath); serr != nil {`
  换成永不触发的等价式（`grep -c`=1、build rc=0）⇒ `--- FAIL: TestAC2PreExistingMigrationBackupIsRepaired (2.24s)`，
  **同一条路径前后两次 icacls 原文对照**：
  BEFORE `config.toml.bak-plaintext carries foreign SID(s) S-1-1-0, S-1-1-0` →
  AFTER `icacls text … still names the foreign principal (Everyone)` +
  `NOT PRIVATE … (principals: [Everyone Everyone BUILTIN\Administrators NT AUTHORITY\SYSTEM DESKTOP-LVS7839\swq])`
  ⇒ "备份永不覆盖"这条既有语义在没接线的 build 上确实是**永久漏**，接线的确是它被修掉的原因；
  同时断言了备份**字节不变**（first-seen original must win 没被 repair 腿破坏）。
  ⚠ 头两次 sed 因把标记写成行尾 `//` 注释而**编译失败**（`if err := …; // 注释` 吃掉后半行 / `undefined: serr`）——
  **编译失败不算变异**，两次都重做后才取数，记下来是因为这就是本仓规矩存在的原因。
  **第 3 条（票面那句事实错）——已按实测改掉，不是"仍未测"**：见上面被替换的 AC#4 构造可行性框与 AC#4 条目。
  读数：未提权（`IsInRole(Administrator)=False`；`swq` 在 Administrators 组但令牌未提升）+
  `AllowDevelopmentWithoutDevLicense = 1` ⇒ `os.Symlink` 到**非空目录**返回 `<nil>`；
  测试里 `os.Symlink ok: artifact -> …\someone-elses-tree (Lstat mode Lrw-rw-rw-, FILE_ATTRIBUTE_REPARSE_POINT set)`。
  两条新用例（`TestAC4DirectorySymlinkAtArtifactPositionIsNotRecursed`、`…InsideSealedTreeIsNotWalked`）
  = 目标内容存活 + 只拆链接 + 注入 `deleteLink` 失败后仍拿 `ErrIsReparsePoint` 且带路径 + 密封游走**不穿**符号链接
  （`verifyPrivate(keep)` 仍失败 ⇒ 目标的 DACL 没被我们改写）。**A51② 的更正登记再扩一格**：`os.Remove`
  对目录符号链接也**能**直接拆（`os.Remove cleared the directory symlink directly; A51② does not reproduce for it on this box`）。
  ⚠ 换一台关掉开发者模式的机器这条构造会失败 ⇒ `mkDirSymlink` 打**具名 `t.Skipf`**（带 `os.Symlink` 原始错误），
  本次跑 **0 SKIP**（不拿"条件构造"当挡箭牌）。真实数据目录未碰，构造全在 `t.TempDir()` 里、测完自清。
  **第 4 条（严格度描述反了）——选了 ②"让清除既存主体这件事报出来"，白名单不开口**：
  **实测语义（把票面/registry 那句错话换成这个）**：给某目录**有意**加第三个主体 ⇒ **写入点不会失败**
  （PROTECTED DACL 让子项根本不继承父上的第三主体，`PrivateFile(<root>/artifact.txt)` 会成功），
  真正发生的是**下一次 `SealDir`/`SealFile` 把它从该对象自己的 DACL 上清掉** ⇒ 风险比"拒写"更大（授权无声消失）。
  为什么不选 ①（白名单扩展点）：验收 §8 的三条理由在本票证据里是可复现的——m4 变异去掉 `PROTECTED` 之后
  `CodexSandboxUsers` 带着 `0x1301ff` 长进"已密封"的文件，`verifyPrivate` 是唯一当场拦住它的东西；
  白名单一开口这条拦截同时失效，而这个包唯一的承诺就是"只有我"。
  **实现**：`winsec_windows.go` 的 `applyDescriptorWindows` 在 set **之前**先读一次 DACL，
  `explicitForeignPrincipals` 只挑"**挂在该对象自己身上**（SDDL flags 不含 `ID`=INHERITED_ACE）且不在 {我,SY,BA} 白名单"的 ACE
  （deny/audit 等非 `A` 条目一并算），密封成功且 `verifyPrivate` 通过之后把每一条交给 `noticeNarrowed`；
  默认实现是 `slog.Warn`，字段 `path` / `cleared` / `policy`，走应用既有的日志管道
  （winsec **不能** import `internal/observe`——图是 observe→secret→winsec，那条边会成环）。
  **为什么只报显式、不报继承来的**：本包存在的理由就是清继承噪声，全报等于没有信号；这一格有用例钉住
  （`TestSealReportsThePrincipalsItCleared` 要求"只继承 root 那条宽 ACE 的孩子**不出现**在通知里"，同时
  root 与自带显式 ACE 的孩子**必须**出现，且出现之后 `verifyPrivate` 真的过——通知不是在对空气喊）。
  **读数（HEAD，`-count=1 -v ./internal/winsec/`）**：`=== RUN` 29、顶层 **19 PASS / 0 FAIL / 0 SKIP**（另含 10 条子测试），
  新增具名 4 条：`TestAC4DirectorySymlinkAtArtifactPositionIsNotRecursed` /
  `TestAC4DirectorySymlinkInsideSealedTreeIsNotWalked` / `TestSealReportsThePrincipalsItCleared` /
  `TestSealNoticeIsRecordedByDefault` 全 PASS；默认通知那条的实测日志原文：
  `level=WARN msg="winsec: seal cleared principals that were placed on this object explicitly" path=…\data cleared=SU(A;OICI;0x1200a9;;;SU) policy="winsec owns the grants on this tree; out-of-band ACEs are removed at the next seal"`。
  门禁（本枚）：`gofmt -l`/`gofumpt -l` 空、`go vet ./internal/winsec/ ./internal/secret/` rc=0、
  `GOOS=linux go vet` 同两包 rc=0（按包作用域，A54③）、`go test -count=2` 两包 rc=0（winsec 12.261s / secret 0.704s）。
  ⚠ **POSIX 侧不适用**：`sealDir` 在 `_other.go` 里没有 walk（chmod 不区分"显式/继承"主体），所以第 4 条的语义是 Windows 专属；
  Linux 真机复跑仍归票 89 已记录的 AC#6 那批（本枚未改 `_other.go`，不重测）。
  next= ① 第 1/2 条的判据+接线+变异读数已齐，第 3/4 条已落 ⇒ 剩 `docs/reports/pending-and-issues.md` 两条更正登记
  （A51②：junction **与**目录符号链接两类上 `os.Remove` 都能拆；新语义："winsec 拥有 data 树的唯一授权权，
  带外授权会在下一次密封时被清除并且**会打一条 WARN**"）——那是编排者的面，我不写；
  ② 收尾枚跑 `sh scripts/d22scan.sh` + 全量四包门禁 + Status 同步。


- 2026-09-21 17:2x（`agent-ticket89b`，**退回单第 1、2 条：判据 + 接线落地**）：
  **第 1 条（AC#2 的覆盖面主张没有用例咬）——判据已进包内**：
  `internal/winsec/acl_windows_test.go::TestAC2SealDirNarrowsChildrenThatCarryTheirOwnExplicitACEs`。
  形状就是验收报告 §6(b) 那台探针（搬进来，不是引用）：`PrivateDirAll(root)` 之后先
  `icacls root /grant *S-1-1-0:(OI)(CI)(RX)`（模拟运维带外授权），再造三个**自带显式 ACE** 的子项
  （`icacls <child> /grant *S-1-1-0:(RX)` 写的是**该对象自己的** DACL——这正是"只有子项自带显式 ACE 时 walk 才承重"那一格）
  与一个**纯继承**的对照孩子，然后只调公开入口 `SealDir(root)`。两条断言腿：
  `assertNoForeignPrincipalText`（**icacls 原文**里 `S-1-1-0` 与 `Everyone` 两种拼法都不许出现）+
  `assertPrivateACL`（名字→SID 解析后按白名单判，两把尺子分开，防止一个解析 bug 让两条一起说谎）。
  **先断言宽再断言窄**：`rawNamesEveryone` 不成立就 `t.Fatalf` 拒绝跑完（票 79 那次"判据被编辑掉"的红就是缺这条）。
  诚实边界也钉死在文件里：那个纯继承的孩子**在删掉 walk 的 build 上仍然绿**（OS 自己重算），
  注释与 `t.Logf` 都明写它"不区分任何事"——所以承重格是"子项自带显式 ACE"，不是"随便建个子文件"。
  **读数（HEAD，`-count=1 -v`，具名 4 条全 PASS、0 SKIP、0 FAIL）**：
  `TestAC2SealedDirCoversFilesItNeverTouched` 1.77s / `TestAC2SealDirNarrowsChildrenThatCarryTheirOwnExplicitACEs` 0.59s /
  `TestAC2MigrationBackupIsPrivate` 0.87s / `TestAC2PreExistingMigrationBackupIsRepaired` 0.25s。
  ⚠ "删掉 `propagatePrivate` 必红"的变异读数在下一枚 commit（本枚只落判据，不许先把结论写上）。
  **第 2 条（同包最坏的产物没封）——一行级接线 ×2 已落，另加 repair 腿 ×1**：
  `internal/secret/migrate.go` 的 `os.WriteFile(backupPath, raw, 0o600)` → `winsec.PrivateFile`、
  `os.WriteFile(tmpPath, out, 0o600)` → `winsec.PrivateFile`；**并**把"备份已存在就不写"那条分支补上
  `winsec.SealFile(backupPath)`——理由：`config.toml.bak-plaintext` 是**永不覆盖**的，
  一台在本修复之前跑过迁移的机器上，那份明文副本会**永远**宽下去，只接新写的路径等于只修未来不修现场。
  **密封前后各取一次 icacls 原文（`internal/winsec/migrate_windows_test.go`，同一对象、同一条路径）**：
  备份文件先 `icacls /grant *S-1-1-0:(RX)` ⇒ BEFORE 读数 `NOT PRIVATE … foreign SID(s) S-1-1-0`（测试自己把它 `t.Logf` 出来），
  跑 `secret.MigratePlaintext` ⇒ AFTER 读数只剩 `NT AUTHORITY\SYSTEM / BUILTIN\Administrators / DESKTOP-LVS7839\swq`，
  且断言备份**字节没被改写**（first-seen original must win 这条既有语义不能被 repair 腿破坏）。
  另一条腿 `TestAC2MigrationBackupIsPrivate` 在带外宽父目录里跑真迁移：迁移前的 `config.toml` 仍是宽的
  （那是 `internal/config` 的写路径，票 95 的账，本票不越界），**它备份出来的 `.bak-plaintext` 已经私有**，
  `secrets\` 全树 sweep 私有，且迁移后的 `config.toml` 也私有（tmp 封了 + rename 保留 descriptor）。
  门禁（本枚）：`gofmt -l internal/winsec internal/secret` 空、`gofumpt -l` 空、
  `go vet ./internal/winsec/ ./internal/secret/` rc=0、`go test -count=1` 两包 `ok`（winsec 7.380s / secret 0.378s）。
  next= 本枚 SHA 上做两次变异（删 `propagatePrivate`、把 migrate 两行退回 `os.WriteFile`）取红名 → 交第 3 条（符号链接真测）
  与第 4 条（静默清除改成会报出来）。

- 2026-09-21 16:1x（编排者，**验收退回单，四条待补；代理已收工 ⇒ 修复另派一人**）：
  `acceptor-ticket89` 八项全部自己测过（未抄本票一个数字），基线/红转绿/"差点假绿"/PROTECTED 变异/
  fail-closed 变异/junction 复现/AC#6 真 Linux 都成立。挡在结案前面的只有四条：
  1. **AC#2 的覆盖面主张没有用例钉住**：验收代理实测**删掉 `sealDir` 里的 `propagatePrivate(path)`
     （保留 `applyDescriptor`、可编译）⇒ 包内 **34 条全绿**，一条都不红。
     它同时给出真实边界：**只有"子项自带显式 ACE"时那次 walk 才承重**（它自造的探针里
     HEAD 全收、去掉 walk 后三条仍宽）。⇒ **修法**：把它现成的探针形状搬进包内，
     补一条"已存在的、带外来显式 ACE 的子项，在父目录被 `SealDir` 之后 icacls 原文里**不再出现 `S-1-1-0`**"的用例，
     并要求它在"删掉 `propagatePrivate`"的 build 上**必红**（验收代理已经验过这个方向）。
  2. **同包一处更重的账（验收代理新挖，本票自己没报）**：`internal/secret/migrate.go:154`
     写的是**迁移前的原始 config 备份**（`config.toml.bak`），里面**含明文密钥**，而它用的是
     **装饰性的 `0o600`**（本票的 AC#1 已经证明那串数字在 Windows 上不落地）
     ⇒ **同一条已结案的缺陷，最坏的一份产物被漏掉了**：主 config 封了、迁移前的明文没封。
     另：`:168` 的 tmp 文件同形状。**修法是一行级 ×2**（走 `winsec.SealFile`/`PrivateFile`）。
  3. **本票面有一句事实错**：它写"普通权限建不出目录符号链接"。验收代理实测：
     **未提权 + 开发者模式=1 ⇒ `os.Symlink` 直接成功** ⇒ 那一格属于"**可测而未测**"，不是"测不了"。
     ⇒ 把这句改成"本机开发者模式下可构造，已构造并测过"或"仍未测"，**二选一**，不许留着。
  4. **本票面把严格度描述错位**：它说"将来有意给某目录加第三个主体 ⇒ 会在**写入点直接失败**"。
     验收代理实测是**下次 `SealDir` 静默清除该主体**（不是当场报错）。
     ⇒ **真实风险更大而不是更小**：静默清除意味着"某人给服务账户加的授权会无声消失"。
     修法二选一并写进票面：要么加一条**显式允许的白名单扩展点**，要么让"清除既存主体"这件事**报出来**。
  **不在本票的**：`winsec.go:126` 的 `filepath.Abs`（票 94，验收代理也明写"不挡本票"）。
  **对 owner 的口径**：验收代理的原话是——"只对我可读"**今天不能无条件说**，只能点名四条主路
  （artifacts / DPAPI blob / db+wal+shm / staging）；`config.toml.bak`、日志、既存带显式 ACE 的老树**仍然宽**。
  等第 1、2 条补完才把这句话升格。
  next= 等票 94 交件让出 `internal/winsec/` ⇒ 派修复代理（一次 commit 能补完 1/3/4，第 2 条是两行接线）。

（空）
- [x] **AC#2 方向：目录链先封、文件兜底（覆盖面靠目录，两条都要）** — 理由：
  (1) 我们**控制不了**所有落盘点：`wisp.db-wal` / `-shm` 是 SQLite 自己建的、`spill.go` 的 `.retry` 是 rename 前的临时件，
  逐个设 DACL 一定会漏掉一个；给父目录设**显式、不再继承自外部**的 DACL（`/inheritance:r` + 只授当前 SID），
  之后创建的子项**继承到的就是这份**——这才真覆盖"别人替我建的文件"。
  (2) 但"先建目录、后设权限"没用：子项继承的是**创建那一刻**的 ACL，所以必须**边建边封**——
  `PrivateDirAll` 逐层建 + 逐层封，任何一层封不上就**整个拒绝**（AC#5）。这条先后关系有用例钉。
  (3) 继承不能靠"应该会收窄"：判据是把父目录**故意**设成带外来 ACE（`Everyone:(OI)(CI)(RX)`），
  再断言子项 `icacls` 原文里**不出现**该 SID；只测默认继承就是测了个空。
  落点包名 `internal/winsec`（票面建议之一）：API 中性 + `_windows.go`/`_other.go` 分平台实现，
  POSIX 侧走真 0600/0700，不静默跳过。
- [ ] AC#1 基线 **icacls 原文**：下一枚 commit 给（本枚只交骨架）。基线环境：Windows 11 Pro，
  当前 SID `S-1-5-21-1228170099-895614386-1166154857-1001`（账户 `swq`）；本机无既存 `wisp.db`（data 根尚未创建）。
  ⚠ `internal/winsec/winsec_windows.go` 此刻是**故意的占位**（只有 `os.Chmod`，也就是仓库今天的行为），
  所以本包判据测试**预期先红**——红完才换 SetSecurityDescriptorInfo。
- [x] AC#4 构造可行性（**2026-09-21 `agent-ticket89b` 实测更正，原句"普通权限建不出符号链接"是错的**）：
  **这台机上未提权就能建目录符号链接**——`IsInRole(Administrator)=False`（账户 `swq` 在 Administrators 组里，
  但进程令牌没提升）+ `HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\AppModelUnlock\AllowDevelopmentWithoutDevLicense = 1`
  ⇒ `os.Symlink(dir, link)` 返回 `<nil>`，`Lstat` 给 `Lrw-rw-rw-`、`os.ModeSymlink` 置位、
  `FILE_ATTRIBUTE_REPARSE_POINT` 也置位（⇒ 实现的 `isReparsePoint` 用属性而不是 `ModeSymlink` 是**必须**的，
  否则密封游走会漏掉所有 junction）。⇒ 那一格属于"**可测而未测**"，现已测：
  `internal/winsec/reparse_windows_test.go::TestAC4DirectorySymlinkAtArtifactPositionIsNotRecursed`
  + `…InsideSealedTreeIsNotWalked`（构造全在 `t.TempDir()` 内，测完自清，未碰任何真实数据目录）。
  实测读数：`os.Remove` 对**目录符号链接**也**能**直接拆掉、目标内容毫发无损
  （`os.Remove cleared the directory symlink directly; A51② does not reproduce for it on this box`）
  ⇒ A51② 在 junction 与目录符号链接**两类对象上都不复现**，本票的判据因此双向断言"目标仍在"。
  **特权前提没被抹掉**：换成关了开发者模式、又没给 `SeCreateSymbolicLinkPrivilege` 的机器，`os.Symlink` 就会失败
  ⇒ `mkDirSymlink` 在这种情况下打**具名 `t.Skipf`**（含 `os.Symlink` 的原始错误），不是静默通过。
- [x] **AC#1 基线：四类私有数据的 icacls 原文**（判据测试 `internal/winsec/acl_windows_test.go`，四类各走**生产入口**：
  artifact = `spill.go:245` 的 `O_CREATE|O_EXCL,0o600`；secret blob = `secret.NewStore`+`Store`（真 DPAPI 落盘）；
  `wisp.db`/`-wal`/`-shm` = `memory.Open` 真 SQLite，**在 store 打开状态下**探测侧文件；staging = `downloader.go:224/:389` 同形 idiom。
  当前用户 SID `S-1-5-21-1228170099-895614386-1166154857-1001`（`DESKTOP-LVS7839\swq`）。**四类每一类的 DACL 首条 ACE 都是同一个外来主体**：

  ```
  acl_windows_test.go:197: icacls "C:\\Users\\swq\\AppData\\Local\\Temp\\TestAC1BaselineProductionPaths240781008\\001\\data\\artifacts\\artifact-89.txt"
  C:\Users\swq\AppData\Local\Temp\TestAC1BaselineProductionPaths240781008\001\data\artifacts\artifact-89.txt DESKTOP-LVS7839\CodexSandboxUsers:(I)(M,DC)
  S-1-5-21-3623186960-731165060-4091685855-1717338598:(I)(M,DC)
  NT AUTHORITY\SYSTEM:(I)(F)
  BUILTIN\Administrators:(I)(F)
  DESKTOP-LVS7839\swq:(I)(F)
  ```
  `secrets\baseline`（DPAPI blob）、`data\artifacts`、`data\secrets`、`data\staging`、`wisp.db`、`wisp.db-wal`、`wisp.db-shm`、
  `staging\model.bin` 全部是**同一份五条 ACE 的形状**（目录多 `(OI)(CI)`），逐条原文见交回报告 / `-v` 输出。⇒
  **`DESKTOP-LVS7839\CodexSandboxUsers` 与一个解析不出名字的别名 SID `S-1-5-21-3623186960-…-1717338598`
  对工件、DPAPI blob、数据库与 WAL 全都持有 `(M,DC)` = MODIFY（含读+写）+ DELETE CHILD**，
  全部由父目录**继承**而来。"只有我能读我的数据"在本机不是"没验证"，是**当场为假**。
  本机主体清单（这条判据到底能不能算"有第二账户"）：本地用户 = `Administrator, CarlosShao, CodexSandboxOffline,
  CodexSandboxOnline, DefaultAccount, Guest, shaowq, swq, WDAGUtilityAccount, WsiAccount`；
  `CodexSandboxUsers` 组成员 = `DESKTOP-LVS7839\CodexSandboxOffline`、`DESKTOP-LVS7839\CodexSandboxOnline`
  ⇒ 持有我们私有数据 MODIFY 的是**两个真实存在、与 `swq` 不同的本机账户**，不是一个空泛的组名。
  但仍**拿不到**"以那两个账户之一真的去 open() 成功/失败"的证据（`runas` 需要口令，非交互拿不到），
  所以 AC#3 走的仍是票面允许的**替代判据**（SID 白名单 + icacls 原文比对），只是它的红因落在真账户上。
- [x] **AC#3 判据先红**（`-count=1 -v`，`go test ./internal/winsec` rc=1）：
  `--- FAIL: TestAC3SealedWritesCarryNoForeignSID/artifact-exclusive`、`…/secret-blob`、`…/staging-temp`
  （3/3 子测试红）+ `--- FAIL: TestAC2SealedDirCoversFilesItNeverTouched`，红因全是被测代码只有 `os.Chmod`：
  `NOT PRIVATE artifact-exclusive: foreign SID(s) S-1-1-0 (principals: [Everyone BUILTIN\Administrators NT AUTHORITY\SYSTEM DESKTOP-LVS7839\swq])`
  ——即"父目录给了 Everyone 读，`0o600` 一声不响地让它长在文件上"。
  ⚠ **本文件差点假绿过一次**：第一版 `aclSIDs` 跳过 `icacls` 首行，而 icacls 把**第一条 ACE 和路径印在同一行**——
  于是"外来主体持有 MODIFY"恰好被丢掉，AC#3 三个子测试在**没有任何修复**的情况下全绿。
  现已改为剥掉路径前缀 + 数每一条 `:(` 行（注释里记了这段），且 `TestAC1BaselineForeignACEPropagatesIntoModeOnlyWrites`
  是**反向钉子**：它断言"宽父目录的 ACE 一定会传到子文件"，一旦这条不成立就红，防止判据测空气。
- [ ] 待办：`winsec_windows.go` 换真 SetNamedSecurityInfo（PROTECTED DACL + 逐层封 + 传播），AC#4 链接、AC#5 失败注入。
- [x] **AC#3 实现落地 + 红→绿**（同一判据、同一台机、同一份 icacls 仪器）：
  `go test ./internal/winsec -count=2` **rc=0**；`-count=1 -v` 全名逐条：
  `--- PASS: TestAC1BaselineProductionPaths` / `--- PASS: TestAC1BaselineForeignACEPropagatesIntoModeOnlyWrites` /
  `--- PASS: TestAC3SealedWritesCarryNoForeignSID`（3/3 子测试）/ `--- PASS: TestAC2SealedDirCoversFilesItNeverTouched`。
  **0 个 SKIP、0 个 FAIL**（本包没有"非 Windows 静默通过"：POSIX 那侧走真 `0600/0700` + 读回校验，见 AC#6）。
  绿后的 icacls 原文（密封过的工件）：
  ```
  C:\Users\swq\AppData\Local\Temp\TestAC3SealedWritesCarryNoForeignSID2468636177\001\data\artifact-exclusive NT AUTHORITY\SYSTEM:(F)
  BUILTIN\Administrators:(F)
  DESKTOP-LVS7839\swq:(F)
  
  ```
  以及 AC#2 覆盖面的那条（**winsec 从没碰过**的文件，只靠父目录继承）：
  ```
  C:\Users\swq\AppData\Local\Temp\TestAC2SealedDirCoversFilesItNeverTouched1351714451\001\data\sqlite-like-sidecar.tmp NT AUTHORITY\SYSTEM:(I)(F)
  BUILTIN\Administrators:(I)(F)
  DESKTOP-LVS7839\swq:(I)(F)
  
  ```
  ⇒ 外来 `Everyone:(I)(RX)` 没了，`CodexSandboxUsers` 那种外来主体也没了，只剩 `SY/BA/我`。
  **`verifyPrivate` 把这条判据搬进了运行时**：每次密封后用 OS 自己的 SDDL 读回，
  出现任何非 `{我, SY, BA}` 主体、或 DACL 不是 `D:P`、或我自己反而没有非 inherit-only 的权 ⇒ `ErrNotSealable`。
  ⇒ "ACE 条数"不能当判据：容器上 OS 会**materialize** `OICIIO` inherit-only 伴生 ACE（实测 3 条变 6 条），
  第一版据此断言 `3 of 3 expected` 直接把 AC#2 打死；换成"无外来主体 + 我自己有实权"才对。
- [x] **实现路上三个真实坑（点名，别下次再踩）**：
  (1) `BuildSecurityDescriptor` 返回 self-relative SD，`sd.DACL()` 读回 present=false ⇒
      报"descriptor carries no DACL"——**假失败但方向是紧的**；改用 `ACLFromEntries`(SetEntriesInAcl)。
  (2) `SetSecurityInfo` **按 handle** 设 DACL 需要句柄带 `WRITE_DAC`，而 `os.OpenFile(O_WRONLY)` 只给了
      GENERIC_WRITE ⇒ `Access is denied`；owner 对该对象**隐式**持有 WRITE_DAC，所以改成按**名字**设
      （`sealHandle` 现在就是 `applyDescriptor(f.Name())`，且在文件还是空的时候调）。
  (3) **我自己的测试重写把 `PrivateDirAll(nested)` 那次调用删掉了**，于是 AC#2 在"什么都没做"的位置上红 ——
      那条红**不是实现缺陷**，是判据被编辑掉了；重加调用后才真绿。记这条是因为：红/绿名如果不带原文，
      这种"测试自己坏了"的红会被当成实现的功劳或罪状。
- [ ] 待办：AC#4（链接）、AC#5（失败注入）、把 memory/agent/secret 三条落盘路径接上 `winsec`、AC#6 的 POSIX 侧点名。
- [x] **AC#5 失败方向只能收紧（红→绿有名字）**：`internal/winsec/private_fail_test.go`（无 build tag，两平台都跑）。
  注入点 = `applyDescriptor` 这个包内接缝（不是"假装失败"，是真的把平台实现换掉）：
  `TestAC5FailedSealRefusesTheWrite` 4/4 子测试绿，且要求**同时**满足 ①错误里 `errors.Is(err, ErrNotSealable)`
  ②`assertNoBytesOnDisk` —— 磁盘上**一个字节都不许留**。四条腿：exclusive 工件 / 覆盖式写 / 目录链 / 直接 SealFile+SealDir。
  ⚠ 这条一开始**红了 4/4**，红因是真的：调用点原样把注入错误往外抛，只有平台实现自己 wrap 时才有
  `ErrNotSealable` ⇒ 契约不能依赖"每个实现都记得 wrap" ⇒ 在公共边界加 `sealError()` 归一化。
  **变异检验（AC#5 的承载行）**：把 `winsec.go:91` 的 `if err := sealHandle(f); err != nil {` 改成
  `if err := sealHandle(f); false && err != nil { // MUTATION-89`（字面就是"设不上就算了、继续以宽权限写"）⇒
  `--- FAIL: TestAC5FailedSealRefusesTheWrite/exclusive_artifact` + `…/replacement_write`（2 红，其余 13 条保持绿 ⇒
  变异只打到该打的判据，不是全炸），红因 `private_fail_test.go:48: PrivateFile accepted a refused seal`；
  还原后 `grep -rn MUTATION-89` 为空、全套重跑 rc=0。编译失败不算变异：本次变异**可编译且行为改变**。
- [x] **AC#4 A51②：本机可构造，用的是 junction（`mklink /J`，普通权限即可）**。
  `internal/winsec/reparse_windows_test.go` 三条全绿：
  ① `TestAC4JunctionAtArtifactPositionIsNotRecursed` —— 链接指向"非空目录"，`RemoveUnlinked` 拆链接本身，
  目标里的 `sub/keep-me.txt` **必须还在**（红一次就是越界删除）；
  ② `TestAC4UnlinkableLinkGivesNamedError` —— 注入 `deleteLink` 失败 ⇒ 必须拿到包着 `ErrIsReparsePoint`
  **且带路径**的错误（这条也先红过一次：错误里路径用 `%q` 打印成 `C:\Users\…`，判据按原文匹配路径 ⇒ 红；改成 `%s` 才绿 —— 记下来是因为
  "错误里有没有路径"这种事只有断言才能发现）；③ `TestAC4SealedWalkSkipsLinks` —— `SealDir` 的传播**不穿链接**
  （用 `verifyPrivate(keep)` 仍失败来证明目标的 DACL 没被我们改写）。
  **诚实结论（A51② 本身）**：本机 Go 1.27 + Win11 上 `os.Remove` **能**删掉指向非空目录的 junction
  （测试里那条 `os.Remove cleared this junction directly; A51② does not reproduce for it` 就是实测日志），
  但无论成功与否都**断言了目标内容仍在**；而**目录符号链接**（**此句原写"这台机普通权限建不出来"，
  2026-09-21 `agent-ticket89b` 实测推翻：开发者模式=1 + 未提权 ⇒ `os.Symlink` 直接成功**）
  ⇒ 那条构造**Windows 侧已补**（`TestAC4DirectorySymlinkAtArtifactPositionIsNotRecursed` /
  `…InsideSealedTreeIsNotWalked`），Linux 侧那条继续作为免特权平台的**第二套**等价构造而不是替代品：
  `internal/winsec/private_other_test.go` 的 `TestPOSIXSymlinkAtArtifactPositionIsNotRecursed`（`os.Symlink` 免特权），
  两侧强度：Windows 侧证"实现不递归、失败具名"，POSIX/Linux 侧证"同一 API 在 unlink 语义下也不越界"。
  **没有**因为 os.Remove 这次成功就把判据删掉：`RemoveUnlinked` 走的是
  `CreateFile(FILE_FLAG_OPEN_REPARSE_POINT|BACKUP_SEMANTICS)+SetFileInformationByHandle(FileDispositionInfoEx)`，
  这是"只可能删到链接本身"的那条形。生产侧接点：`memory/artifacts.go` 的 `removeStray` 两处 `os.Remove(full)` 已换成它。
- [x] **落盘路径接线（票面点名的三个家）**：`agent/spill.go`（artifacts 目录 → `PrivateDirAll`；`.retry` 临时件 → `PrivateFile`；
  `writeFileExclusive` → `PrivateFileExclusive`）、`secret/store.go`（secrets 目录 + DPAPI blob）、
  `memory/open.go`（data 根 / artifacts / backup）、`memory/artifacts.go`（reclaim 用 `RemoveUnlinked`）。
  **生产路径端到端 icacls 证据**：`TestAC3ProductionDataRootIsPrivateEndToEnd`（宽父目录下 `memory.Open`+`secret.NewStore`，
  打开态探 `wisp.db`/`-wal`/`-shm`，再全树 sweep：`sweep checked 8 entries … all private`）；
  `internal/agent/spill_acl_windows_test.go::TestAC3SpillArtifactLandsPrivate` 走真 `Spiller.Prepare`：
  `icacls tool-output-call_acl.txt -> [NT AUTHORITY\SYSTEM BUILTIN\Administrators DESKTOP-LVS7839\swq]`。
  ⚠ **没接的同族**（本票落点之外，要编排者拍板）：`models/downloader.go:224/:389`（staging，AC#1 实测同样带外来 ACE）、
  `config/migrate.go:83`+`config/parse.go:213`（0o600 配置备份）、`observe/logging.go:244`（0o644 日志）、`ball/position.go:77`。

- [x] **AC#6 门禁（只跑本票碰的包）**：
  `gofmt -l internal/winsec internal/agent internal/memory internal/secret` 空；
  `gofumpt -l . tools/d22scan tools/mockllm`（**全仓**，每次 commit 前重跑）空；
  `go vet ./internal/winsec/ ./internal/agent/ ./internal/memory/ ./internal/secret/` rc=0；
  `GOOS=linux go vet` 同四包（**按包作用域**，A54③）rc=0；
  `go test -count=2` 同四包 rc=0（winsec 15.8s / agent 9.0s / memory 32.2s / secret 1.0s）。
  **四种假绿逐条点名**：`--- SKIP` **0 条**；没有用 `-run` 过滤来交数（`-run` 只用于变异检验那两次，
  且当时的红/绿名逐条列出）；`-count=2` 用 `-count=1 -v` 另跑一遍核对名字与条数
  （**winsec 包 17 条具名结果**：AC#1 2、AC#2 1、AC#3 1+3 子、AC#4 3、AC#5 1+4 子、生产端到端 1、AC#5 happy-path 1……
  全 PASS；另有 `internal/agent` 的 `TestAC3SpillArtifactLandsPrivate` 1 条走真 `Spiller.Prepare`），0 FAIL 0 SKIP；
  没有"步骤被静默跳过"——POSIX 侧在 **Linux 真机**（docker alpine:3.20，交叉编译 `go test -c` 的二进制）跑出来：
  `TestPOSIXPrivateFileIsReally0600` / `TestPOSIXPrivateDirIsReally0700` /
  `TestPOSIXSymlinkAtArtifactPositionIsNotRecursed` / `TestPOSIXMissingFileIsNotAnError` **全 PASS + 真断言真 0600**，
  `TestAC5*`（1+4 子）在 Linux 上同样 PASS ⇒ Linux 侧共 **10 条具名结果、0 SKIP**：非 Windows 平台不是"跳过"，是"测另一套机制"。
  ⚠ 一个**格式化工具互相打架**的坑记下来：`gofumpt -w` 生成的多返回值 map 字面量缩进 `gofmt -l` 不认（两把尺子互斥），
  本仓两把都要空 ⇒ 改成先赋值再 return 的写法，两边都满意。
- [x] **变异检验 #2（AC#2 承载行）**：锚点 = `winsec_windows.go:53` 的
  `windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION`；
  去掉 `PROTECTED` 后（可编译、行为改变，不是编译失败）：
  `--- FAIL: TestAC3SealedWritesCarryNoForeignSID/{artifact-exclusive,secret-blob,staging-temp}` +
  `--- FAIL: TestAC2SealedDirCoversFilesItNeverTouched`，红因是 OS 读回的原文自己招了：
  `D:AI(A;;FA;;;SY)(A;;FA;;;BA)(A;;FA;;;S-1-5-21-…-1001)(A;ID;0x1200a9;;;WD)(A;ID;FA;;;BA)…`
  ⇒ **`WD`(Everyone) 带着 `0x1200a9`（读+执行）从父目录长了进来**，"继承会放大权限"被实测出来，
  而且运行时 `verifyPrivate` 当场拒绝落盘（失败方向=收紧）。
  还原后 `git diff --quiet -- internal/winsec/winsec_windows.go` 干净，`grep -rn MUTATION-89 --include=*.go .` = 0 命中。

**交回前还剩什么（诚实）**：本票落点之外的同族 `0o600/0o755` 站点**没有全接**：
`models/downloader.go:224/:389`（staging，AC#1 实测带外来 ACE）、`config/migrate.go:83` +
`config/parse.go:213`、`observe/logging.go:72/:244`（0o644 日志）、`ball/position.go:73/:77`、
`secret/migrate.go:154/:168`、`cmd/wisp/doctor.go:248`、`cmd/wisp/slo_windows.go:559`（禁改区）。
接线都是一行级的替换，但要有人裁决"日志/模型缓存算不算私有数据"。

  **2026-09-21 `agent-ticket89b` 更正（退回单第 2 条之后）**：上面这串里的 **`secret/migrate.go:154/:168` 已接**
  （两处 → `winsec.PrivateFile`，并给"备份已存在就不写"那条分支补了 `winsec.SealFile` 的 repair 腿，
  判据 `internal/winsec/migrate_windows_test.go` 两条 + 删掉接线必红的变异读数）。
  **其余站点不变**：`models/downloader.go`、`config/*`、`observe/logging.go`、`ball/position.go`、`cmd/wisp/*`
  仍等票 95/90 的裁决与地界。

  next= 编排者：(1) 三枚 commit（`b994a2c`→`de15a6b`→`57bdbb2`+本次这枚）已过全部署门禁，可派验收；
  (2) 上面那串同族站点要不要一起收（建议开一张小票，或并到票 90 的"权限模式"里）；
  (3) A51② 需要更正登记：**本机 Go 1.27/Win11 实测 `os.Remove` 能删指向非空目录的 junction**
  （票 79 的观察在这台机上不复现，但目标内容存活被双向断言了），目录符号链接仍受特权限制、构造放在 Linux 侧；
  (4) 验收若要"第二账户真的读不到"的更强证据，需要一台能 `runas` 的机器（本箱拿不到那两个沙箱账户的口令）。
