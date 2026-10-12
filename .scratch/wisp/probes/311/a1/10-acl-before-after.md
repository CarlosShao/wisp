# 311-a1 `10-acl-before-after` — 两枚根各自封前／封后的逐 ACE 名册（`AC#0` 问①）

尺＝`powershell -NoProfile -File ps1/dump-acl.ps1`（`Get-Acl`＋逐枚 `FileSystemAccessRule` 解析成 SID＋数值掩码），
同行第二根棍子＝同一件里跑的 `icacls.exe` 逐字渲染。读数件＝`logs/10-acl-inherit-baseline.txt`、`logs/11-acl-before-fresh-dir-both-roots.txt`、
`logs/13-acl-after-plant2-both-roots.txt`、`logs/16-acl-after-residue-tree-B.txt`、`logs/18-acl-after-residue-leaves-B.txt`、`logs/19c-acl-after-guests-leaf-B.txt`。
**本程⛔ 改过任何一字节码／⛔ 动过任何一枚既有目录的 ACL**（下面所有"封后"读数要么来自本腿自己新建的根里的真码残留，要么来自本腿新建的探针目录）。

## 1. 封之前＝两枚根各自"新建一枚子目录"落到的继承集合

### 根 A ＝ `C:\Users\swq\AppData\Local\wisp-test-gotmp\311a1-A-20261012-133040-1594\aclprobe`（`logs/11` 逐字）

```
sddl=…D:(A;OICIID;FA;;;SY)(A;OICIID;FA;;;BA)(A;OICIID;FA;;;S-1-5-21-1228170099-895614386-1166154857-1001)
ACE[0] type=Allow sid=S-1-5-18               who=NT AUTHORITY\SYSTEM       rights=FullControl mask=0x1F01FF inh=ContainerInherit, ObjectInherit prop=None isinherited=True
ACE[1] type=Allow sid=S-1-5-32-544           who=BUILTIN\Administrators    rights=FullControl mask=0x1F01FF inh=ContainerInherit, ObjectInherit prop=None isinherited=True
ACE[2] type=Allow sid=S-1-5-21-1228170099-895614386-1166154857-1001 who=DESKTOP-LVS7839\swq rights=FullControl mask=0x1F01FF inh=ContainerInherit, ObjectInherit prop=None isinherited=True
ace_count=3
icacls: NT AUTHORITY\SYSTEM:(I)(OI)(CI)(F) / BUILTIN\Administrators:(I)(OI)(CI)(F) / DESKTOP-LVS7839\swq:(I)(OI)(CI)(F)
```
父链同形（`logs/10`）：`C:\Users\swq\AppData\Local` 与 `…\Local\wisp-test-gotmp` 各自 `ace_count=3`、trustee 完全同上三枚。

### 根 B ＝ `D:\tmp\wisp-test-gotmp\311a1-B-20261012-133040-1594\aclprobe`（`logs/10`＋`logs/11` 逐字）

```
sddl=…D:(A;ID;FA;;;BA)(A;OICIIOID;GA;;;BA)(A;ID;FA;;;SY)(A;OICIIOID;GA;;;SY)(A;ID;0x1301bf;;;AU)(A;OICIIOID;SDGXGWGR;;;AU)(A;ID;0x1200a9;;;BU)(A;OICIIOID;GXGR;;;BU)
ACE[0] Allow sid=S-1-5-32-544 who=BUILTIN\Administrators        rights=FullControl              mask=0x1F01FF    inh=None prop=None isinherited=True
ACE[1] Allow sid=S-1-5-32-544 who=BUILTIN\Administrators        rights=268435456                mask=0x10000000  inh=ContainerInherit, ObjectInherit prop=InheritOnly isinherited=True
ACE[2] Allow sid=S-1-5-18     who=NT AUTHORITY\SYSTEM           rights=FullControl              mask=0x1F01FF    inh=None prop=None isinherited=True
ACE[3] Allow sid=S-1-5-18     who=NT AUTHORITY\SYSTEM           rights=268435456                mask=0x10000000  inh=ContainerInherit, ObjectInherit prop=InheritOnly isinherited=True
ACE[4] Allow sid=S-1-5-11     who=NT AUTHORITY\Authenticated Users rights=Modify, Synchronize    mask=0x1301BF    inh=None prop=None isinherited=True
ACE[5] Allow sid=S-1-5-11     who=NT AUTHORITY\Authenticated Users rights=-536805376             mask=0xE0010000  inh=ContainerInherit, ObjectInherit prop=InheritOnly isinherited=True
ACE[6] Allow sid=S-1-5-32-545 who=BUILTIN\Users                 rights=ReadAndExecute, Synchronize mask=0x1200A9  inh=None prop=None isinherited=True
ACE[7] Allow sid=S-1-5-32-545 who=BUILTIN\Users                 rights=-1610612736               mask=0xA0000000  inh=ContainerInherit, ObjectInherit prop=InheritOnly isinherited=True
ace_count=8
icacls: BUILTIN\Administrators:(I)(F) /(I)(OI)(CI)(IO)(F) / NT AUTHORITY\SYSTEM:(I)(F) /(I)(OI)(CI)(IO)(F) /
        NT AUTHORITY\Authenticated Users:(I)(M) /(I)(OI)(CI)(IO)(M) / BUILTIN\Users:(I)(RX) /(I)(OI)(CI)(IO)(GR,GE)
```
父链同形（`logs/10`：`D:\`、`D:\tmp`、`D:\tmp\wisp-test-gotmp` 三处 `ace_count=8`、trustee 同上八枚；`D:\` 那发还逐字带着 `NT AUTHORITY\Authenticated Users:(M)` 这一枚来源）。

## 2. 问① 的答复：差的就是下面这几枚 trustee

| trustee（逐枚 SID＋掩码） | 根 A | 根 B | 对本 token 有效吗 |
|---|---|---|---|
| `S-1-5-21-1228170099-895614386-1166154857-1001`（`DESKTOP-LVS7839\swq`＝**本进程自己的 SID**） | **有**，`mask=0x1F01FF`(F)，`(OI)(CI)`，非 InheritOnly | **没有**（这枚 SID 在两枚根的集合里只出现在 A 侧 3/10/11 那几发） | A 侧唯一给本 token 有效权利的 ACE；`0x1F01FF` 含 0x40（见 §4） |
| `S-1-5-11`（`NT AUTHORITY\Authenticated Users`）`mask=0x1301BF`(Modify,Synchronize) ＋ 兄弟 `mask=0xE0010000`(inherit-only，`SDGXGWGR`) | **没有** | **有**，两枚 | 本 token 持 `S-1-5-11`（`whoami //groups` 逐字 `"NT AUTHORITY\Authenticated Users","Well-known group","S-1-5-11","Mandatory group, Enabled by default, Enabled group"`）⇒ 生效，但 0x1301BF **⛔ 含 0x40** |
| `S-1-5-32-545`（`BUILTIN\Users`）`mask=0x1200A9`(ReadAndExecute,Synchronize) ＋ 兄弟 `mask=0xA0000000`(inherit-only，`GXGR`) | **没有** | **有**，两枚 | 本 token 持 `S-1-5-32-545` ⇒ 生效，但 `0x1200A9` 只给读＋ traverse，⛔ 写⛔ 删 |
| `S-1-5-18`（SYSTEM）`mask=0x1F01FF`（＋B 侧 inherit-only 0x10000000） | 有 | 有 | **⛔ 生效**（本 token⛔ 是 SYSTEM）⇒ 两枚根上这枚都是纸面权利 |
| `S-1-5-32-544`（`BUILTIN\Administrators`）`mask=0x1F01FF`（＋B 侧 inherit-only 0x10000000） | 有 | 有 | **⛔ 生效且反向**：`whoami //groups` 逐字 `"BUILTIN\Administrators","Alias","S-1-5-32-544","Group used for deny only"` ⇒ **deny-only**，它在 ACL 里的 F 对本 token 是⛔ 不给、而是拒〔真机量的组属性；"deny-only 使 F 变拒绝"那半句是读码推的〕 |

**一句话**：两枚根差的不是掩码写法，是**一枚 trustee 的有无**——A 侧继承了"本 token 自己 `…-1001`：F(0x1F01FF)、`(OI)(CI)`"，B 侧没有它、只有 `S-1-5-11`：0x1301BF 与 `S-1-5-32-545`：0x1200A9（本 token 另外能命中的两枚 SYSTEM/Administrators 一枚⛔ 在 token 里、一枚在 token 里是 **deny-only**）。
⇒ **在 A 的整条父链上本 token 都带 0x40（`DeleteSubdirectoriesAndFiles`＝FILE_DELETE_CHILD），在 B 的整条父链上一枚都⛔ 带。** 这一枚差＝本票全部颜色差的来源（问③＝`30-cleanup-failure.md`）。

## 3. 封／种之后（同一份码、⛔ 换判据）

### 根 B＝真码留下的 6 枚锁死叶（`logs/16`／`logs/18`／`logs/19c` 逐字）

| 残留叶（相对根 B） | `dacl_protected_flag` | `sddl` 的 DACL 部分 | ace_count | 逐枚 trustee＋掩码 |
|---|---|---|---|---|
| `TestGateJudgesThePrivateSetByResolvedSID3972594964\001\la` | **True** | `D:PAI(A;;FA;;;LA)` | 1 | `S-1-5-21-…-500`(Administrator) F 0x1F01FF, isinherited=False |
| `TestGateJudgesThePrivateSetByResolvedSIDrunner_descriptor_is_ju1386323932\001\runner` | **True** | `D:PAI(A;;FA;;;SY)(A;OICIIO;GA;;;SY)(A;;FA;;;BA)(A;OICIIO;GA;;;BA)(A;;FA;;;LA)(A;OICIIO;GA;;;LA)` | 6 | `S-1-5-18` F 0x1F01FF／`S-1-5-18` 0x10000000(IO)／`S-1-5-32-544` F 0x1F01FF／`S-1-5-32-544` 0x10000000(IO)／`S-1-5-21-…-500` F 0x1F01FF／`S-1-5-21-…-500` 0x10000000(IO)，全部 isinherited=False |
| `…a_name_spelling_cannot_2711238538\001\named` | True | 同 `runner` 那六枚（逐字同串） | 6 | 同上 |
| `…a_name_spelling_cannot_2711238538\001\numeric` | True | 同 `runner` 那六枚（逐字同串） | 6 | 同上 |
| `TestSealNarrowsAndNamesThePrincipalItRemovedBySID1250499198\002\la` | True | `D:PAI(A;;FA;;;LA)` | 1 | `S-1-5-21-…-500` F 0x1F01FF |
| `TestSealNarrowsAndNames…1250499198\003\la` | True | `D:PAI(A;;FA;;;BG)` | 1 | `S-1-5-32-546`(BUILTIN\Guests) F 0x1F01FF |

**这六枚里没有一枚点名本 token 实际持有的 SID**（`…-500`⛔ 在 token、`S-1-5-32-546`⛔ 在 token、`S-1-5-18`⛔ 在 token、`S-1-5-32-544` 在 token 里是 deny-only）⇒ 封／种之后本 token 对这六枚**既⛔ 能列举⛔ 能删**。
旁证（同一把尺的独立形状）：`deep-list.ps1` 逐层走时这六枚**枚枚**报 `CANNOT_ENUMERATE=… err=UnauthorizedAccessException: Access to the path '…' is denied.`（`logs/14`，根 B `enumerate_denied_count=7` ＝这六枚＋本腿自己的 `plant2-runner`）。
而它们的父目录（`\001`、`\002`、`\003`）逐字仍是 §1 那八枚继承 ACE（`logs/16`/`logs/18` 里 `dacl_protected_flag=False`＋`sddl=…(A;ID;0x1301bf;;;AU)…`）⇒ 父目录能列举、却⛔ 能给"按名字删掉一枚自己无权的孩子"供上 FILE_DELETE_CHILD。

### 根 A＝同一条 SDDL 在 A 上落出什么（⚠ 代用尺，见 `00-anchor.md` §5 末行）

`logs/13` 逐字（`ps1/plant-probe.ps1` 新建 `…\311a1-A-…\plant2-runner` 后 `Get-Acl`）：

```
dacl_protected_flag=True
sddl=…D:PAI(A;;FA;;;SY)(A;OICIIO;GA;;;SY)(A;;FA;;;BA)(A;OICIIO;GA;;;BA)(A;;FA;;;LA)(A;OICIIO;GA;;;LA)
ACE[0]…ACE[5] 与根 B 的 \001\runner 逐枚同 sid／同 mask=0x1F01FF|0x10000000／同 isinherited=False
ace_count=6
```
根 B 同一枚探针（`…\311a1-B-…\plant2-runner`）＝**逐字同串**（同 `dacl_protected_flag=True`、同 6 枚、同掩码）。
⇒ **封之后的 DACL⛔ 随卷变；随卷变的只有"父链给本 token 哪几枚权利"。** 本腿⛔ 试图删任何一枚（⛔ chmod、⛔ takeown），所以"同一枚锁死叶在 A 上删得掉、在 B 上删不掉"这一条**⛔ 由本腿真机演示**，它的真机证据＝§4 那两发 `go test` 的收尾结果（同一枚码在 A 上把同形的叶删干净了＝残留枚数 0，在 B 上留了 6 枚）。

## 4. 支持问①那句结论的尺件行（逐字，`logs/11-acl-before-fresh-dir-both-roots.txt`）

```
FileSystemRightsEnumNames=ListDirectory,ReadData,WriteData,CreateFiles,CreateDirectories,AppendData,ReadExtendedAttributes,WriteExtendedAttributes,Traverse,ExecuteFile,DeleteSubdirectoriesAndFiles,ReadAttributes,WriteAttributes,Write,Delete,ReadPermissions,Read,ReadAndExecute,Modify,ChangePermissions,TakeOwnership,Synchronize,FullControl
maskbit DeleteSubdirectoriesAndFiles=64
maskbit Modify=197055
maskbit FullControl=2032127
Modify_has_0x40=False
FullControl_has_0x40=True
0x1301BF_has_0x40=False
0x1F01FF_has_0x40=True
```
⇒ 0x40 的具名名字＝`DeleteSubdirectoriesAndFiles`（＝FILE_DELETE_CHILD／icacls 里的 `Dc`）；`0x1301BF`（B 侧唯一落到本 token 的写权利）**⛔ 含它**，`0x1F01FF`（A 侧 `…-1001` 那枚）**含它**。〔掩码＝真机量的；"这一枚就是 RemoveAll 删锁死叶所缺的权利"＝读码推的，凭据见 `30-cleanup-failure.md` §4〕
