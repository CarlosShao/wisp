# 票 73：孤儿 `.wisp-tmp-*` 清扫器落地，A18 判据②选"自愈"并测了四件事

**时间：** 2026-09-21（本地）· **执行者：** 票 73 实现代理
**判据来源：** `.scratch/wisp/issues/73-orphan-staging-file-sweep.md`（registry **A18** / **Q-16**）
**前提不重做：** `761447f` + `docs/evidence/s1/20-a18-taskkill-residue-characterization.md` 已经用
**真子进程 + 真 `taskkill /F`** 量过：每次打断恰好留 1 个 `.wisp-tmp-*`，且**没有任何东西扫它**。
本文件只记"加了清扫器之后什么是真的"，不重推那三条。

## 0. 落地的东西

| 文件 | 作用 |
|---|---|
| `internal/tools/fs_staging.go` | 命名/归属方案 + 扫描器 + D31 账本措辞 |
| `internal/tools/staging_live_windows.go` | `OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION)` + `GetExitCodeProcess` 判创建者是否还活着；`FILE_ATTRIBUTE_REPARSE_POINT` 判 reparse |
| `internal/tools/staging_live_other.go` | 非 Windows：**判不了 ⇒ 一律不删**（清扫退化成 no-op，不假装成功） |
| `internal/tools/fs_write.go:282/:505` | 挂点：`fs.write` 的 `stageAndRename` 与 `fs.move` 跨卷，都在 `os.CreateTemp` **之前** |
| `internal/tools/fs_staging_windows_test.go` | AC#2/#3/#4 三条判据用例 |
| `internal/tools/bridge_a18_kill_windows_test.go` | A18 第三条断言按票面翻向 |

**时机 = 下一次写盘**（不是桥启动）：`internal/tools` 没有可挂的启动单点（`New()` 只是构造，
`cmd/wisp` 不属于本票），而"下一次写盘进同一个目录"正是 A18 特征化第三条已经在做的动作——修和证落在同一个缝上。

**归属 = 文件名自带身份** `.wisp-tmp-<owner8>-<pid>-<rand>`：
`owner8 = sha256(可执行文件路径, 用户名, ".wisp-tmp-")` 前 8 位十六进制。
选它而不是"进程内记一本账"的理由是硬的：**死进程的残留必须由活进程认领**，进程内账本跨不过进程死亡；
选它而不是"写个 sidecar 清单"的理由是：清单本身会变成第二种垃圾，而且它崩溃时和暂存文件一样不一致。
代价写明：**Wisp 换安装路径或换用户运行后，旧 owner token 的孤儿永远不再可归因 ⇒ 变成永久垃圾**。
这是故意的保守侧——"扫了别人的文件"是删除原语，"多留一个 4 字节垃圾"不是。

## 1. 四条 AC 的真实结果（本机 2026-09-21，全部 `//go:build windows`、零 `t.Skip`）

| # | 断言 | 结果 |
|---|---|---|
| AC#1 | 真 kill 后下一次写桥 ⇒ 目录里 `.wisp-tmp-*` **归零**，目标 `existing.txt` 逐字未动 | `TestA18…TheNextWriteReclaimsTheStagingFile` **PASS**（3/3 子项） |
| AC#2 | 撞前缀的外人文件逐个点名存活；死 pid 的可归因孤儿被扫掉；本进程活 pid 的孤儿存活 | `TestSweepReclaimsOnlyItsOwnStagingFiles` **PASS** |
| AC#2b | 名字层的近似匹配全拒（owner 不对 / pid 非数字 / pid=0 / 段数不对 / 随机段带 `\`） | `TestSweepNamingSchemeRejectsNearMisses` **PASS**（跑出一处修复，见 §3） |
| AC#3 | 名字**完全可归因**的真 `mklink /J` junction 不被删，对面文件字节未动；同目录普通孤儿确实被扫掉（证明清扫真跑了） | `TestSweepNeverDeletesThroughARealJunction` **PASS** |
| AC#4 | 活进程持有句柄的暂存文件不删**且写盘成功**；释放后下一次写盘必须扫掉 | `TestSweepSparesATempFileHeldByAnotherProcess` **PASS** |

AC#4 的形状刻意做成两道独立防线各咬一次：
**A** 名字里是活 pid（ nobody 持有句柄）⇒ 靠判活跳过；
**B** 名字里是**已死** pid、但另一活进程以 `FILE_SHARE_READ|FILE_SHARE_WRITE`（**不给 `FILE_SHARE_DELETE`**）
打开并阻塞 ⇒ 只剩"重试后仍失败就跳过"这一道能救它；然后真 `taskkill /F` 释放，
下一次写盘**必须**扫掉——否则 B 的幸存就可能是"名字没匹配上"而不是"句柄挡住了"。

## 2. 变异检验（票面只要求 AC#1，AC#2/#3/#4 是危险面所以补做）

- **AC#1**（票面要求）：`fs_staging.go` 的 `sweepStagingOrphans` 首行后插 `return got // MUTATION-AC1`。
  先 grep 证明改动真落：`fs_staging.go:210: return got // MUTATION-AC1` ⇒
  `go test ./internal/tools/ -run TestA18 -count=1` **FAIL**
  （`真 kill 之后的 .wisp-tmp-* 残留 = […47166889-26212…, …47166889-51572…], want 1 个`）。
  删掉该行 ⇒ `grep -rn MUTATION-AC1 internal/tools/` **rc=1（零命中）** ⇒ 复跑 `ok 1.046s`。
- **AC#2/#3/#4**（把清扫退回票面禁止的形状："前缀即我的 + 不判盘上是什么"，
  插 `pid = "4294967295" // MUTATION-73`）：先 grep 命中 `fs_staging.go:224` ⇒ 三条用例**全红**，
  AC#3 报的是 `清扫器连 junction 本身都删了 … GetFileAttributesEx …\.wisp-tmp-69a1ec61-33652-777: The system cannot find the file specified`
  ——这正是"清扫变成越界删除原语"的形状。恢复后 grep 零命中、`TestSweep|TestA18` 全绿。

## 3. 测试揪出来的真洞（不是测试的锅，是实现的）

`TestSweepNamingSchemeRejectsNearMisses` 第一次跑就红：`stagingAttribution(".wisp-tmp-<owner>-0-1")`
返回了 true —— 我写匹配器时只判了"pid 段全是数字"，没排除 **pid 0**（System Idle Process，
不可能是文件创建者；`OpenProcess(0)` 在 Windows 上会打开 Idle 进程 ⇒ 判活为"活着" ⇒ 那条名字
永远扫不掉，属于"静默失真"而不是危险）。修法：匹配器显式拒绝 `pid == "0"`
（`fs_staging.go:159`）。**这条断言先于修复存在，所以它是测出来的不是我编出来的。**

## 4. 门禁数字（AC#5）

```
gofmt -l internal/tools                       → 空
go vet ./internal/tools/...                   → rc=0
go test ./internal/tools/... -count=2         → ok  27.617s
go test ./internal/tools/... -race -count=2    → ok  36.283s
sh scripts/d22scan.sh                          → 两步都过：
    positive control  go test ./... (tools/d22scan)  → ok   github.com/CarlosShao/wisp/tools/d22scan
    go run . -root "D:/work/workspace/projects plans/Wisp"
                                                       → d22scan: clean - no D22 ban violations
                                                       （ban #7 scope internal/tools/ examined 16）
```

**d22scan 能变红这件事是在我这棵树上证的，不是引用它的自测**：临时写
`internal/tools/zz_seed_violation_73.go`（`filepath.Clean(filepath.Abs(p))`）⇒
扫描 **rc=1**、两条 `[pathresolver-bypass]` 命中并打 `2 finding(s); D22 bans are not negotiable`；
删掉文件 ⇒ 复扫 rc=0 clean，`ls` 确认种子文件不在树里（没有留垃圾）。

⚠ 一条**不属于本票**的观察：第一次跑 `sh scripts/d22scan.sh` 时该模块**编译不过**
（`tools/d22scan/main.go:818:26: undefined: io`），而脚本的 `set -eu` 之外没有拦到这一步——
`script-rc=0` 是从 `tail` 读的。**脚本自身的退出码是否会把 go test 的失败带出去，属 d22scan 在飞代理的地盘**，
我没碰 `tools/d22scan/`（票面禁我进），只登记。

## 5. 没证明但看起来成立的事

- **换用户/换安装路径后的旧孤儿**：按 §0 的设计它们永久不可归因。我没有测"重装后残留堆积"的量级上限，
  也没测真机上 Wisp 的安装路径是否稳定（`cmd/wisp` 不属于本票）。
- **多个 Wisp 进程并发写同一目录**：判活规则保护了"创建者还活着"的那一份，但**pid 复用**
  （死进程的 pid 被别的进程拿走）会让该孤儿被判成"活的"而永久留下——方向是安全的（不删），
  我没测它发生的概率，也没加超时兜底（加了就会重新引入误删风险）。
- **孤儿里可能含明文敏感内容**：Q-16 的备选判据里提过。清扫器只删不读，
  没有对残留内容做任何"是否敏感"的判断，也没有把残留次数/字节数上报到 observe 层。
- **rename 之后、Remove 之前那一段被真 kill**：`Hooks.Kill` 覆盖不到，真 kill 时机我固定不了
  （沿用上张票的 §5 边界，没扩大)。
