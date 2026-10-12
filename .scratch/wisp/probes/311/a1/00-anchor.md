# 311-a1 `00-anchor` — 锚、本程两问、台面与尺

## 1. 锚

- 仓库＝`D:\work\workspace\projects plans\Wisp`，分支＝`dev`，起手锚＝HEAD `8e3cb3701ccbb541964e4acc4e74d8b56a8c99c7`
  （逐字读数：`git log -1 --format='%H %ad %s'` ⇒ `8e3cb3701ccbb541964e4acc4e74d8b56a8c99c7 2026-10-12 13:22:08 +0800 派 310-r2：…`）。
- 本程＝票 311 `AC#0`（只读归因）＋`AC#3` 那一问。**⛔ 产码、⛔ 动判据、⛔ 翻框、⛔ 碰 `Status:`、⛔ 改名 `-done`、⛔ 删任何东西。**
- 派单正文＝`.scratch/wisp/probes/orch/20261012/txt/311a1-dispatch.txt`（已整份读）；票面＝`.scratch/wisp/issues/311-winsec-gate-colour-depends-on-which-volume-the-test-temp-root-sits-on.md`（已整份读）；
  编队正典＝`.scratch/wisp/probes/orch/gate-roster-304.md`（已整份读）。

## 2. 本程只答的两问（射程⛔ 扩大）

1. `AC#0` 三问：① 两枚根的**继承集合**差哪几枚 trustee（逐枚 SID＋掩码）；② `TestGateJudgesThePrivateSetByResolvedSID` 的红句里那枚"解析出来却⛔ 在集合内"的 SID；③ 那 4 行 `Access is denied` 锁在哪一步（封之后自我锁死 ⇄ 根本身拒 `WRITE_DAC`）。
2. `AC#3`：`TempDir RemoveAll cleanup: … Access is denied` 算不算票 310"死亡三形"之外的第四形、落点在哪张票、触发条件写成一句能判的话。

答复正文＝`10-acl-before-after.md`（①）／`20-red-sentence.md`（②）／`30-cleanup-failure.md`（③）／`90-final.md`（`AC#3`＋残留名册＋自报）。

## 3. 台面（两枚根都是我本程新建的，⛔ 碰任何既有目录的 ACL）

| 代号 | 全路径 | 形状来路 |
|---|---|---|
| 根 A（出厂默认形） | `C:\Users\swq\AppData\Local\wisp-test-gotmp\311a1-A-20261012-133040-1594` | `${LOCALAPPDATA}/wisp-test-gotmp/<本次一枚>`＝`scripts/winsec-tests.sh:55` 那同一枚基（`winsec-<ts>-$$` 换成本腿的 `311a1-A-<ts>-<pid>`），⛔ 走那三枚入口件 |
| 根 B（换卷形） | `D:\tmp\wisp-test-gotmp\311a1-B-20261012-133040-1594` | 与编排者 12:21 那发同基（`D:/tmp/wisp-test-gotmp/…`），只差末段名 |

逐字自落读数（`logs/00-runids.txt`）：
`ROOTA=C:\Users\swq\AppData\Local/wisp-test-gotmp/311a1-A-20261012-133040-1594`／`ROOTB=D:/tmp/wisp-test-gotmp/311a1-B-20261012-133040-1594`／`TS=20261012-133040-1594`。

两发对照的跑法（⛔ `scripts/*.sh`，直接 `go test`，逐字在 `logs/20-…`／`logs/21-…` 的第一行）：

```
cmd=GOFLAGS= GOTMPDIR=C:/Users/swq/AppData/Local/wisp-test-gotmp/311a1-A-20261012-133040-1594 go test -v -count=1 ./internal/winsec/
cmd=GOFLAGS= GOTMPDIR=D:/tmp/wisp-test-gotmp/311a1-B-20261012-133040-1594 go test -v -count=1 ./internal/winsec/
```

台面同码：`go env GOOS=windows`、`GOROOT=D:\work\base\go`、`GOFLAGS` 空；两发都 `-count=1`（⛔ 靠缓存）。

## 4. 颜色复现（本腿自己这两发，⛔ 抄编排者的数）

- 根 A：逐字 `--- PASS: TestGateJudgesThePrivateSetByResolvedSID (0.01s)`／`--- PASS: TestSealNarrowsAndNamesThePrincipalItRemovedBySID (0.08s)`／`ok github.com/CarlosShao/wisp/internal/winsec 8.145s`／`gotest_rc=0`；计数尺同行＝`RUN=101 PASS=58 FAIL=0 SKIP=0`、`grep -c 'TempDir RemoveAll cleanup' = 0`。
- 根 B：逐字 `--- FAIL: TestGateJudgesThePrivateSetByResolvedSID (5.45s)`（两枚子例 `1.55s`／`1.93s`）＋`--- FAIL: TestSealNarrowsAndNamesThePrincipalItRemovedBySID (1.57s)`＋`FAIL github.com/CarlosShao/wisp/internal/winsec 14.710s`＋`gotest_rc=1`；计数＝`RUN=101 PASS=56 FAIL=2 SKIP=0`、`cleanup_lines=4`。
- ⇒ 与票面现量**逐字同数**（`58/0` ⇄ `56/2`、4 行 ⇄ 0 行），且**同一枚 `internal/winsec` 代码、本腿⛔ 改过一字**。

## 5. 尺与代用尺（本机缺的工具具名写出，⛔ 沉默跳过）

| 用途 | 用的尺 | 本机可得性 |
|---|---|---|
| 逐 ACE 打 DACL（SID＋掩码＋继承位） | `powershell -NoProfile -File ps1/dump-acl.ps1 <路径…>`（`Get-Acl` → 每枚 `FileSystemAccessRule` 的 `Translate(SecurityIdentifier)`＋`[int]FileSystemRights`＋`InheritanceFlags`＋`PropagationFlags`；并同行落 `sddl=`／`dacl_protected_flag=`／`owner=`） | 可用 |
| 同一读数的第二根棍子（OS 自己的渲染） | 同一件里 `& $env:WINDIR\System32\icacls.exe`（逐字 `ICACLS| …:(I)(M)` 那些行） | 可用 |
| 枚名字尺 0x40 是哪一枚权利 | `dump-acl.ps1` 头部 `[int][System.Security.AccessControl.FileSystemRights]::…`，逐字 `maskbit DeleteSubdirectoriesAndFiles=64`／`Modify_has_0x40=False`／`FullControl_has_0x40=True`／`0x1301BF_has_0x40=False`（`logs/11-…`） | 可用 |
| token 身份与组（"这枚 SID 我到底持⛔ 持"） | `ps1/gate-roster.ps1` 的 `[Security.Principal.WindowsIdentity]::GetCurrent()` ＋ `/c/Windows/System32/whoami.exe //groups //fo csv` | 可用；⚠ PowerShell 那把**⛔ 报 deny-only 组成员**（见 `30-cleanup-failure.md` 第 3 节），故两把并排 |
| 名册（⛔ 拿空输出当 0 枚） | 起手／交件后两发 `logs/05-gate-start.txt`／`logs/30-gate-end.txt`：`Get-ChildItem $env:TEMP -Directory` 里 `Test*` 枚数＋家族尺 `*.test.exe`＋正控 `svchost.exe`／`*.exe` 同行 | 可用 |
| 根内条目名册 | `logs/17-…`（bash `find -type d`）⇄ `logs/14-…`（`ps1/deep-list.ps1` 逐层 `Get-ChildItem -Directory -Force`，⛔ `-Recurse`）——⚠ `-Recurse` 那把**进不去锁死叶**，会把枚数读少（`logs/15-…` 里 `enumerate_denied_count=7`、`ps_errors=8` 就是它报错的地方） | 可用 |
| "封之后"的 ACL 在**根 A** 上取一份 | `ps1/plant-probe.ps1`＝**代用尺**：`Set-Acl`＋`DirectorySecurity.SetSecurityDescriptorSddlForm($sddl, AccessControlSections::Access)`（走同一枚 Win32 `SetNamedSecurityInfo`，⛔ 是 wisp 的码路） | 部分可用：`D:PAI…` 那一枚（带 `P`）在两枚根上落出与真码残留**逐字同形**的 6-ACE protected DACL；`D:(A;;FA;;;LA)`（不带 `P`）那一枚 .NET **⛔ 复现 `PROTECTED_DACL_SECURITY_INFORMATION`**，落出 4/9 ACE 的"仍带继承"形 ⇒ 那一枚只算半根棍子，真凭据＝根 B 上那 6 枚**真码留下的残留**逐字读数 |
| `tasklist`／`Get-CimInstance` 家族尺 | `gate-roster-304.md` §1 那四条精确名＋宽松家族尺 | 可用（本程只用家族尺＋正控，四条精确名⛔ 起进程，仍落枚数） |
| 长度门 | `sh scripts/check-path-length-budget.sh --with-self-test`（⛔ 三枚在飞的入口件，那三枚 `scripts/*.sh`＋`tools/d22scan/runtests.sh` 此刻归落地腿 `310-r2`） | 见 `90-final.md` §4 的 rc |

## 6. 抄给后续程的一条已知坑（本腿⛔ 用那枚脚本，但票面 §3 末要这条）

⚠ `bash scripts/winsec-tests.sh --scope=winsec` 会被它自己的 **GUARD 1 打回 `rc=2`**（正确形＝`bash scripts/winsec-tests.sh ./internal/winsec/`）。本腿两发**完全绕开**那三枚入口件，理由是 `310-r2` 正在改它们＝借在飞态当台面；这条拼写坑是**转抄票面 `:12`／正典 §1**，本腿⛔ 重跑它（⛔ 用那枚脚本 ⇒ ⛔ 有新的 rc=2 读数）。
