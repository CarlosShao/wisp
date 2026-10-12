# 311-a1 `30-cleanup-failure` — 那 4 行 `Access is denied` 锁在哪一步（`AC#0` 问③）

答案先行：**锁在"封／种之后的自我锁死"那一支，⛔ 在"根本身拒 `WRITE_DAC`"那一支**；
而且"自我锁死"本身在两枚根上**一模一样**，换卷之所以让它现形，是因为**父目录给⛔ 给 `FILE_DELETE_CHILD`（0x40）这第二道门**（根 A 给、根 B ⛔ 给）。四句成因各自标的见每节末。

## 1. 那 4 行逐字（根 B 那发，`logs/21-gotest-rootB-dtmp.txt`，行号同件）

```
130:    testing.go:1617: TempDir RemoveAll cleanup: openfdat D:/tmp/wisp-test-gotmp/311a1-B-20261012-133040-1594\TestGateJudgesThePrivateSetByResolvedSIDrunner_descriptor_is_ju1386323932\001\runner: Access is denied.
133:    testing.go:1617: TempDir RemoveAll cleanup: openfdat D:/tmp/wisp-test-gotmp/311a1-B-20261012-133040-1594\TestGateJudgesThePrivateSetByResolvedSIDa_name_spelling_cannot_2711238538\001\named: Access is denied.
135:    testing.go:1617: TempDir RemoveAll cleanup: openfdat D:/tmp/wisp-test-gotmp/311a1-B-20261012-133040-1594\TestGateJudgesThePrivateSetByResolvedSID3972594964\001\la: Access is denied.
142:    testing.go:1617: TempDir RemoveAll cleanup: openfdat D:/tmp/wisp-test-gotmp/311a1-B-20261012-133040-1594\TestSealNarrowsAndNamesThePrincipalItRemovedBySID1250499198\002\la: Access is denied.
```
根 A 那发同行尺＝`grep -c 'TempDir RemoveAll cleanup' = 0`（件 `logs/20-…`）。

## 2. 这行句子是谁打的（GOROOT `D:\work\base\go` 逐字行）

```
testing/testing.go:1613:  c.tempDir, c.tempDirErr = os.MkdirTemp(os.Getenv("GOTMPDIR"), pattern)
testing/testing.go:1617:          c.Errorf("TempDir RemoveAll cleanup: %v", err)
testing/testing.go:1633:  dir := fmt.Sprintf("%s%c%03d", c.tempDir, os.PathSeparator, seq)   ← 所以路径里那枚 \001 \002 \003 = 第几次 t.TempDir()
testing/testing.go:1650:  const arbitraryTimeout = 2 * time.Second                            ← removeAll 的重试窗（对上 1.55s/1.93s/1.57s 那几枚时长）
testing/testing_windows.go:25:  if err == syscall.ERROR_ACCESS_DENIED { return true }          ← 正因为"可重试"，每一枚锁死叶都要白跑满窗
```
⇒ **这条错误⛔ 出自任何一条判据**：它是 `t.TempDir()` 注册的 `Cleanup` 在**测试正常跑完之后**打的一声 `Errorf`（`go test` 收尾阶段），用例的断言全部走完了〔读码推的（源码逐字在上）＋真机量的（§1 那 4 行的宿主是 `testing.go:1617`；两发全文里判据句一枚都没响，尺见 `20-red-sentence.md` §1 末那张 `grep -c` 表）〕。

## 3. 锁死的是哪一步：`os.RemoveAll` 的两段尝试都被拒（GOROOT 逐字）

```
os/removeall_at.go:66-68   func removeAllFrom(parentFd, base) { err := removefileat(parentFd, base)   ← 第一段：按名字删
os/removeall_at.go:79      if err != EISDIR && err != EPERM && err != EACCES { … return unlinkat … }  ← 被拒才往下走
os/removeall_at.go:91      file, err := openDirAt(parentFd, base)                                     ← 第二段：打开它以递归清空
os/removeall_at.go:104     recurseErr = &PathError{Op: "openfdat", Path: base, Err: err}              ← ★我们看到的 "openfdat … Access is denied" 就是这一行
os/root_windows.go:356     func removefileat(dirfd, name) { return windows.Deleteat(dirfd, name, FILE_NON_DIRECTORY_FILE) }
os/root_windows.go:360     func removedirat (dirfd, name) { return windows.Deleteat(dirfd, name, FILE_DIRECTORY_FILE) }
internal/syscall/windows/at_windows.go:246-249  NtOpenFile(&h, FILE_READ_ATTRIBUTES|DELETE, …)        ← 第一次：要对象自己的 DELETE
internal/syscall/windows/at_windows.go:260-262  // Access denied, try opening with DELETE only.
                                                 // This may succeed if the file has restrictive permissions
                                                 // but the caller has delete child permission on the parent directory.
internal/syscall/windows/at_windows.go:263-270  NtOpenFile(&h, DELETE, …)                             ← ★第二次：靠的就是父目录的 FILE_DELETE_CHILD
```
把盘上读到的 DACL 套进这两段（叶与父的逐字 DACL＝`10-acl-before-after.md` §3／§1）：

| 步 | 需要的那一枚权利 | 根 B 上给⛔ 给 | 根 A 上给不给 |
|---|---|---|---|
| 第一次 `NtOpenFile(FILE_READ_ATTRIBUTES\|DELETE)` 打叶 | **叶自己的 DELETE** | ⛔ 给（叶的 DACL 逐字＝`D:PAI(A;;FA;;;LA)`／`D:PAI(…SY)(…BA)(…LA)`／`D:PAI(A;;FA;;;BG)`，`…-500`／`S-1-5-32-546`／`S-1-5-32-544`(deny-only)／`S-1-5-18` 一枚都⛔ 在本 token 里） | 同样⛔ 给（同一枚码、同一枚 DACL） |
| 第二次 `NtOpenFile(DELETE)` 打叶 | **父目录的 FILE_DELETE_CHILD(0x40)** | ⛔ 给：父（`…\001`／`…\002`／`…\003`）逐字只有 `S-1-5-11: mask=0x1301BF` 与 `S-1-5-32-545: mask=0x1200A9` 落到本 token，`0x1301BF_has_0x40=False` | **给**：父链每一层都带 `S-1-5-21-…-1001`(本 token 自己) `mask=0x1F01FF`，`0x1F01FF_has_0x40=True` |
| 第三段 `openDirAt` 打叶（第二段失败后的兜底） | 叶的 LIST_DIRECTORY | ⛔ 给 ⇒ 报 `openfdat … Access is denied`（`removeall_at.go:104` 那一行） | 走⛔ 到这一步（第二步已经删掉了） |

⇒ 问③的具名答复：**"封之后自我锁死"是必要的一半，但单它⛔ 成立**；同一枚锁死叶在根 A 上被收尾删得掉，靠的是第二次 `NtOpenFile` 用父目录的 `FILE_DELETE_CHILD` 兜了过来。
根 B 那枚卷的继承集合里**没有一枚给本 token 0x40**，于是两段全拒、兜底也拒 ⇒ 4 行 ＋ 目录留在盘上。〔"哪一段需要哪一枚权利"＝读码推的（源码逐字在上）；两枚根各自的 DACL／掩码／token 组成员＝真机量的（`logs/10`、`logs/11`、`logs/16`、`logs/18`、`logs/19c`、`logs/05`＋`whoami //groups` 那发）〕

## 4. ⛔ 把"根本身拒 `WRITE_DAC`"当成成因——两笔反证（都是真机量的）

1. **写 DACL 在根 B 上确实成功了**：那 6 枚残留叶的 DACL 就是被写进去的（`logs/16`／`logs/18`／`logs/19c` 逐字 `dacl_protected_flag=True`＋`sddl=…D:PAI…`），写它用的正是 `SetNamedSecurityInfo`（`plantDescriptor` 与 `applyDescriptorWindows` 走的那一枚）；同一发的根 B 还打出 58 条封通告（`grep -c 'cleared_inherited=' = 58`，其中 53 条有内容），通告只在"封写成功之后"才发。⇒ 那一发里 `WRITE_DAC` 一路够用。
2. **根 B 的继承集合本来就⛔ 给我 `WRITE_DAC`，而那正好是"仍写得动"的那一支**：`maskbit ChangePermissions=262144`（0x40000）⇄ 落在我 token 上的最大掩码 `0x1301BF` 里⛔ 含它 ⇒ 写成功靠的是**属主的隐式 `WRITE_DAC`**（四枚属主读数逐字 `owner=DESKTOP-LVS7839\swq`）。
   （"属主隐式持 WRITE_DAC"这一句＝读码推的（Windows 的既有语义，⛔ 是本仓的码）；属主读数与掩码读数＝真机量的。）
3. 顺带一条排除：失败的**⛔ 是"根"那枚目录**。根 B 那枚 `…\311a1-B-20261012-133040-1594` 本身在 §1 里逐字带 8 枚继承 ACE、可列举可写；4 行里的路径首段全都落在 `TestGateJudges…\001\runner`、`…\002\la` 这种**叶子**上（`removeall_at.go:131-134` 那句 `pathErr.Path = base + PathSeparator + pathErr.Path` 拼出来的就是"父\子"）。

## 5. 正控与反形（同一条结论的两个方向，都是真机量的）

- **反形正控**＝同一发（根 B）里那枚**被封成"私有集合含本 token"**的目录：`TestSealNarrows…1250499198\001`（`store` 那枚，`PrivateDirAll` 的产物）在残留名册里**没有**——`logs/14-residue-deep-list-both-roots.txt` 与 `logs/17-residue-census-two-rulers.txt` 两把尺里，那枚 tempDir 下只剩 `\002`、`\003`（及其 `la`），`\001` 已被正常收尾删掉。⇒ **"父不给 0x40 也能删"的另一半条件＝叶自己的 DACL 里有本 token**：`store` 已被删掉、盘上取不到它的"封后"读数，它的 DACL 含 `S-1-5-21-…-1001` 这一条⛔ 是逐字量的、是从断言推的——`privateACL` 给本 token 一枚 grant（`winsec_windows.go` 的 `privateEntries`），且 `TestSealNarrows…` 走到 `verifyPrivate(root)` 那一支而两发全文里 `grep -c 'grants the current user nothing' = 0`／`'names principals outside the private set' = 0`／`'DACL is missing or unprotected' = 0`〔"删掉的因＝它给了本 token"＝读码推的；"它被删了、那三句判据⛔ 响"＝真机量的〕。两枚条件缺一即锁。
- **正控**＝同一枚探针 SDDL 在两枚根上落出**逐字相同**的 6-ACE protected DACL（`logs/13-acl-after-plant2-both-roots.txt`：A 与 B 的 `plant2-runner` 都是 `D:PAI(A;;FA;;;SY)(A;OICIIO;GA;;;SY)(A;;FA;;;BA)(A;OICIIO;GA;;;BA)(A;;FA;;;LA)(A;OICIIO;GA;;;LA)`、`ace_count=6`、`dacl_protected_flag=True`），且两枚都**进不去**（`logs/14` 的 `CANNOT_ENUMERATE=…plant2-runner … is denied.` 一行在 A、B 各一发）⇒ 锁死状态⛔ 随卷变；随卷变的只有父链那枚 0x40。
- **枚数与行数的关系**（⛔ 把 4 当"4 枚"）：根 B 锁死叶共 **6** 枚（`logs/14` 逐枚具名），`TempDir RemoveAll cleanup` 只有 **4** 行——因为 `removeAllFrom` 只把**第一枚**错带回来（`removeall_at.go:129` 那句 `if recurseErr == nil {`），而 `testing.go` 的清理是**每枚 tempDir 一行**（`testing.go:1613` 那枚 `c.tempDir` 才是清理单位，`testing.go:1633` 的 `%03d` 只是它下面的序号）。4 行 ⇄ 4 枚 tempDir（`3972594964`／`runner_descriptor_…1386323932`／`a_name_spelling_cannot_…2711238538`／`1250499198`）逐一对得上；"父\子"那个拼法＝`removeall_at.go:126` 那句 `pathErr.Path = base + string(PathSeparator) + pathErr.Path`。
