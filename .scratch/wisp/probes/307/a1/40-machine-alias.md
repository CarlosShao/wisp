# 307-a1 · 40-machine-alias —— ★必答："本机绿"能不能当凭据

只读取数，⛔ 改任何全局/用户配置（无 `fsutil 8dot3name set`、无注册表写、⛔ 删任何目录）。
脚本＝`logs/40-alias-probe.ps1`／`logs/41-alias-world-probe.ps1`（`Scripting.FileSystemObject` 的 `GetFolder(...).ShortPath`，写成 `.ps1` 用 `powershell -NoProfile -ExecutionPolicy Bypass -File` 调）。
输出＝`logs/40-alias-probe.txt`／`logs/41-alias-world-probe.txt`，各带 `rc_ps1*=`／`fsutil_*_rc=` 自落行。

## 1. 现量：各段 8.3 别名取得到取不到（〔现跑〕2026-10-10）

台面 env（`logs/40-alias-probe.txt` 头部）：`COMPUTERNAME=DESKTOP-LVS7839`、`USERNAME=swq`、
`USERPROFILE=C:\Users\swq`、`TEMP=TMP=C:\Users\swq\AppData\Local\Temp`、**`TMPDIR=`（空）**、`LOCALAPPDATA=C:\Users\swq\AppData\Local`。

| 路径 | ShortPath 读数 | 判定 |
|---|---|---|
| `C:\Users\swq`（＝`USERPROFILE`） | `C:\Users\swq` | **SHORT==LONG ⇒ 取不到别名** |
| `C:\Users` | `C:\Users` | SHORT==LONG |
| `C:\Users\swq\AppData\Local` | 同长形 | SHORT==LONG |
| `C:\Users\swq\AppData\Local\Temp`（＝`TEMP`/`TMP`） | 同长形 | **SHORT==LONG ⇒ `TEMP` 链上任一段都⛔ 别名**（逐段 walk：`C:\Users`/`swq`/`AppData`/`Local`/`Temp` 五段全 ==LONG） |
| `C:\Users\swq\tmp` | `C:\Users\swq\tmp` | SHORT==LONG（票面点名要量的那一枚） |
| `C:\Users\swq\tmp\riskattrib1` | `C:\Users\swq\tmp\RISKAT~1` | **SHORT!=LONG ⇒ 别名取到** |
| `C:\Users\swq\tmp\riskattrib1\RunnerTempAreaForAttribution1` | `C:\Users\swq\tmp\RISKAT~1\RUNNER~1` | **SHORT!=LONG ⇒ 取到，且与票面 §1 那两枚短名逐字相同**（票面那两发用的目录**仍在盘上**，本腿没删、也没重建） |

`C:\Users` 下逐账户（同一把尺，`logs/40-alias-probe.txt` 的 USERS 段）：
`All Users`→`ALLUSE~1`、`CarlosShao`→`CARLOS~1`、`CodexSandboxOffline`→`CODEXS~1`、`Default User`→`DEFAUL~1`、`WsiAccount`→`WSIACC~1` 有别名；
`Default`、`Public`、`shaowq`、`swq`、`zaiuser1` **SHORT==LONG**。
⇒〔读码推到（成因层）〕这台机器"绿"的原因是**账户名恰好 ≤8 字符且全 8.3 合法字符**，⛔ 是"这台机器不生成短名"。

## 2. 现量：这台机器／这个卷到底能不能造出短名世界（〔现跑〕）

`logs/41-alias-world-probe.txt`：在 `%TEMP%` 下新建两枚长名目录（`wisp307a1AliasWorldProbe\LongSegmentNameBeyondEight`，**只建不删**，落仓库外）——

- `C:\Users\swq\AppData\Local\Temp\wisp307a1AliasWorldProbe` → `...\Temp\WI0253~1` ⇒ **SHORT!=LONG**
- `...\WI0253~1\LongSegmentNameBeyondEight` → `...\WI0253~1\LONGSE~1` ⇒ **SHORT!=LONG**

⇒ **别名生成对这枚卷是开着的**（只要名字超过 8 字符或含 8.3 非法字符就会分配短名）。
`fsutil 8dot3name query` ⇒ `The registry state is: 2 (Per volume setting - the default)`，`fsutil_machine_rc=0`；
`fsutil 8dot3name query C:` ⇒ `Failed to open the volume. Error 5: Access is denied.`，`fsutil_volume_rc=1` ⇒ **逐卷状态本腿取不到（需提权），具名欠这一枚读数**（本腿⛔ 提权）。

## 3. ★具名回答那一问

> "如果换一台 profile 没有短名别名的机器，这 5 枚会不会变成静悄悄的绿？"

**会。** 判据链条（前两环现跑、第三环读码）：

1. 〔现跑〕这 5 枚的 `home` 拼写来自 `t.TempDir()` ⇒ 继承 `%TEMP%`/`%TMP%`（本机＝`C:\Users\swq\AppData\Local\Temp`，`TMPDIR` 空），而这条链**每一段 SHORT==LONG**（第 1 节逐段 walk）⇒ 这台机器上 `t.TempDir()` 只会回长形 ⇒ `s.home`（`normPath` 短形折法）与 `cand`（`Resolve` 长形折法）**逐字相同** ⇒ 那枚不对称**在断言之前就不存在**。
2. 〔现跑〕票面 §1 那两发（别名台面 `FAIL=12` vs 本机台面 `PASS=12`）正是把第 1 环当成唯一变量：只差 `TEMP/TMP/TMPDIR` 一枚 env。
3. 〔读码推到〕⇒ 在一台 profile ⊕ TEMP 全链无别名的机器（＝**这台**、以及任何账户名 ≤8 字符的机器）上，`:162/:213/:233/:355/:402` 会**全绿且零信号**：⛔ `--- SKIP`、⛔ 失败消息、⛔ 任何"世界缺了这一枚形状"的自陈。绿读起来像"折叠对称性已被验证"，实际验证的是"这台机器⛔ 具备触发条件"。**⇒ 本机绿⛔ 能当凭据。**
4. 反向读数（⛔ 修法，只报这台机器给不给得出那个世界）：第 2 节证明**这台机器造得出**短名（新局长名目录有 `WI0253~1`/`LONGSE~1`），且包内**已经存在两枚 OS 级取短名的夹具函数**——`pathresolver_junction_windows_test.go:58 getShortPath`（`syscall.GetShortPathName`）与 `syncdirs_redteam_windows_test.go:29 tryShortPath`（`windows.GetShortPathName`），后者体内注释逐字"Every spelling below is produced by the OS (GetShortPathNameW), never typed"〔现跑 grep，输出在 `logs/g-f-83fixtures.txt`〕。⚠ 这两枚今天都带"无别名就 `t.Skip`"的出口（`pathresolver_junction_windows_test.go:128`、`syncdirs_redteam_windows_test.go:54`）——票面"边界"节已把那枚 `t.Skip` 逃逸口登记在 `A831`，**本腿⛔ 把它当本票的修法建议、⛔ 碰 `AC#3` 框**，只报到"这台机器给得出短名、包内已有取短名的现成手"这一读数层。

## 4. 与本腿其它件的咬合

- `AC#3` 那道判据（"⛔ 依赖机器形状"）在本机上的**唯一诚实读数**＝第 3 节第 1 环：判定所需的"短形一侧"来自 env 继承，⛔ 来自夹具 ⇒ **修复后两把台面各自的颜色**必须由能同时给出两把台面的人跑（本腿⛔ 跑编译面）。
- 跨面调用者读数（〔现跑〕，`logs/g-k-crossface.txt`）：`grep -rn --include=*.go "SyncDetectionComplete(\|IsSyncPath(" internal/panel internal/perm cmd` ⇒ **0 行、rc=1**；`IsSyncPath(` 非测试面全仓＝2 行＝定义 `syncdirs.go:187` ＋调用 `internal/risk/provenance.go:928`（与票面"生产调用者＝1 枚"同号同文）。⚠ `syncdirs.go:181-182` 那句注释另称面板 Security page 消费本判定 ⇒ 本腿**按注释读⛔ 采信**，也⛔ 复核 `SyncRoots()` 的消费面＝具名欠读数。
