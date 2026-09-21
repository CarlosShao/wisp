# 89 — `0o600` 在 Windows 上是装饰品：artifact / 密钥 / 数据库**从来没真的"只有我能读"过**（A51①②）

**Status:** open
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

- [ ] **AC#1** 先给**基线证据**：真机上把这四类各建一个文件（artifact / secret blob / `wisp.db` / staging 临时件），
      用 `icacls <path>` 打印**实际 ACL**，逐条贴原文。⚠ **不许用 Go 的 `FileInfo.Mode()` 当证据**——
      它给的就是那套被 Windows 忽略的假位。
- [ ] **AC#2** 定**方向**（写进票面并说明为什么）：给**文件**逐个设 DACL，还是给 **data 根目录**设一次、
      子项继承？后者覆盖面大但要求目录已存在且早于任何写入 ⇒ **建目录与建文件的先后**必须有用例钉住。
      两条路都要保证：**继承不能把权限放大**（拿 `icacls /inheritance` 之类的显式断言，别靠默认行为）。
- [ ] **AC#3** 实现 + **真红过的判据**：在**当前用户之外还要有第二个账户或可读探针**才算证明。
      拿不到第二账户时的**可接受替代**：断言 ACL 里只有当前 SID 有读写、且 `SYSTEM`/`Administrators` 之外的
      任何 SID 不出现；并跑一次 `icacls` 原文比对。**"我把 `0o600` 换成了调用 SetSecurityDescriptor"不算证明**——
      那只是写了代码，没测到东西。
- [ ] **AC#4** A51② 那条单独钉：造一个"artifact 位置是指向目录的符号链接"，
      断言 ①**不被递归删除**（链接目标目录里的文件必须还在），②给出**具名错误**而不是静默失败。
      ⚠ 若 Windows 上普通权限无法建符号链接，**就如实写"本机无法构造"并给出 CI 侧（Linux）等价构造**
      + 说明两个平台上这条判据各自的强度。**不许**因为"造不出来"就把这条判据删掉。
- [ ] **AC#5** 失败方向必须是**收紧**：ACL 设不上时报错并**拒绝落盘**（或退到"不落敏感数据"），
      **绝不允许**"设不上就算了、继续以宽权限写"。判据用例：模拟一次失败注入 ⇒ 必须是错误。
- [ ] **AC#6** 门禁（只跑你碰的包）：`gofmt -l` 空、`gofumpt -l <files>` 空、`go vet <pkgs>` rc=0、
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
- [ ] AC#4 构造可行性：普通权限建不出符号链接（需 `SeCreateSymbolicLinkPrivilege`/开发者模式），
  但 **junction 不需要特权**（`mklink /J`），而它正是"指向非空目录的链接"这一类，`os.Remove` 对它
  报 `ERROR_DIR_NOT_EMPTY` ⇒ **本机可构造**，不必只靠 CI/Linux 侧等价构造。
