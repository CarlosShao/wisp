# 票 20 第 5 框：桥层（不是 risk 层）的 junction / 8.3 短名拒绝——真产物 + 变异检验

**时间：** 2026-09-21 08:50–09:15（本地）· **执行者：** 票 20 接续代理（第 5 框）
**分支/提交：** `dev`，commit `test(20): 桥层真 junction/8.3 拒绝（第 5 框）`
**判据文件：** `.scratch/wisp/issues/20-host-bridge-fs-tools.md` 第 5 框（票面框无编号，按位置数第 5 条）

## 0. 这一框原来欠的是什么（不是"再证一次下层会拒"）

票 18 已经在 `internal/risk/pathresolver_junction_windows_test.go` 里用真产物证明
`Resolve()` 会拒绝 junction 与 8.3 穿越。开工前我先复现了编排者 23:3x 的那条 grep：

```
grep -rn "junction\|Junction\|mklink\|GetShortPathName" internal/tools/*_test.go   → 0 命中
```

⇒ 缺的是**上层**的证据：桥有没有把下层的"不"吞掉、降级成放行，或者根本没走到下层
（= registry **A33②** 那一族"实现存在但装配根没人调"）。本文件证明的是后者不成立。

生产可达性不在本框范围内重做：`cmd/wisp/run.go` 已 `tools.New` + `approval.New` + 注册全部
`BuiltinFSEntries`（票面头部更正块 ①，端到端在 `cmd/wisp/run_test.go:245`/`:330`）。

## 1. 真产物，没有一条用例用字符串假装

新文件 `internal/tools/bridge_junction_windows_test.go`（`//go:build windows`，与票 18 同惯例）：

| 产物 | 怎么造的 | 前置条件（缺了就**显式红**，本文件零 `t.Skip`） |
|---|---|---|
| 真 NTFS junction | `exec.Command("cmd","/c","mklink","/J",link,target)`（`/J` 不需要 SeCreateSymbolicLinkPrivilege） | 临时卷是 NTFS；`t.Fatalf` 会打印缺哪一条 |
| 真 8.3 短名 | `syscall.GetShortPathNameW` 两遍取回，并要求回显含 `~` 且与长名不同 | 该卷启用短名生成（`fsutil 8dot3name query <vol>` = 0）；不满足则 `t.Fatalf` 写明要管理员改什么 |
| 阳性对照 | 每个 fixture 建完 junction 后**先用 `os.ReadFile` 从链接对面读出那 94 字节**，读不到就直接判失败 | 没有这条，后面的"拒绝"可能只是产物坏了 |

目标目录改动前后都用 `dirSnapshot()`（逐文件 sha256）比对，"没有任何字节落在目标里"
断言的是整个目录 before == after，不是"顺手 stat 了那一个文件"。

**本机实跑结果：两种真产物都建成了**（同一台机、同一个 `internal/risk` 的既有产物惯例）。

## 2. 用例与断言（每行都是 `-count=2` 实测）

| 用例 | 断言的桥层行为 |
|---|---|
| `TestBridgeRefusesARealJunctionOnTheReadRoute`（+2 子项） | `fs.read`/`fs.list` 打 `<allowed>\jn\victim.txt`：**门是 AnswerAllow 的**，仍 `IsError`、目标内容/条目一个字节都没出现；判定 L2 via R2；红队拒绝的 reason 是"无法规范化"，**不与**"目标路径在授权目录之外"共用措辞（对照项：真越界读确实是后者措辞） |
| `TestBridgeWritesNothingThroughARealJunction`（+5 子项） | `fs.write`（覆写 / 新建）、`fs.trash`、`fs.move`、`fs.delete` 逐个穿过真 junction：全部被拒，目标目录 sha256 快照逐字节不变、`brand-new.txt` 没出现、`moved.txt` 没出现、被"回收/删除"的文件仍在原地、两侧目录都没有 `.wisp-tmp-*` 残留 |
| `TestJunctionInsideAnAllowedRootCannotReachAnAListFile` | 授权根里种一个指向**真 `.ssh` 树**的 junction：读不到密钥内容；见 §4 那条发现 |
| `TestBridgeRefusesTheRealShortNameOfAnAListFile` | A 档文件（`<home>\.ssh\id_ed25519_short_name_probe`，就住在 allowed 根里，所以 R2 挑不出来、只有 R3 能抓）用 `GetShortPathNameW` 的**真短名**调用：判定 `Deny` via **R3**、`error_class=permission_denied`、门**一次都没被调用**（A 档不给人点批准的机会）、卡面路径是解析后的长路径 |
| `TestShortNameSpellingGetsTheSameVerdictAsTheLongOne`（+1 子项） | 越界文件的短名拼法与长名拼法**判定逐字相同**（L2 via R2、卡面显示长路径）；子项钉住反向：8.3 不是红队门，人批准后就该读到那个长路径指向的文件（卡面与字节必须同一个文件） |
| `TestWavedThroughJunctionWriteBooksAnApprovedRowThatWroteNothing` | 真 SQLite 行：被批准的红队调用**没落一个字节**，但 `tool_call` 记的是 `risk=L2 decision=allow outcome=error error_class=tool`——见 §4 |

`-count=2 -v` 实跑（09:4x，本机，本票两个新文件都已落地）：**`=== RUN` 36、`--- PASS` 36（顶层 + 子项）、
`--- SKIP` 0、`--- FAIL` 0，连字符串 `skip` 都是 0 次**；`go test -count=2 ./internal/tools/` → `ok 25.192s`；
`go test -race -count=1 ./internal/tools/` → `ok 17.709s`。

## 3. 变异检验（否则这 6 条就是自证）

变异 = A17 说的那个形状：**把下层的"不"吞掉**。两行，各一处：

```go
// internal/tools/fs.go:95-98   (FSDeps.open)
-	if err != nil { return "", err }
+	if err != nil { return raw, nil } // MUTATION-20: swallow the resolver's denial
// internal/tools/fs_write.go:169-174 (FSDeps.canonical)
-	return d.Paths.Canonicalize(raw)
+	return raw, nil                   // MUTATION-20
```

先证变异真落盘（`grep -n MUTATION-20` → `fs.go:97`、`fs_write.go:173` 两行），再跑：
**6 条顶层用例 + 8 个子项全红**（`FAIL github.com/CarlosShao/wisp/internal/tools 1.918s`）。真实红字（截）：

```
TestBridgeRefusesARealJunctionOnTheReadRoute/fs.read_through_a_real_junction:
  桥穿过了一个真 junction（94 字节返回）: text="只有穿过 junction 才读得到的内容：s3cr3t..."
.../fs.list_through_a_real_junction:
  ... text="C:/Users/.../001/jn (1 条目)\nf victim.txt 94 ..."
TestBridgeWritesNothingThroughARealJunction/fs.write_over_the_junctioned_target:
  ... text="已写入 C:/Users/.../001/jn/victim.txt（28 字节，临时文件+原子重命名）"
.../fs.trash_through_the_junction:
  ... text="已放入回收站：.../001/jn/victim.txt（还原记录 $I41EUHE.txt，可在回收站还原）"
.../fs.delete_through_the_junction:
  ... text="已永久删除：.../001/jn/victim.txt（未进回收站；找回请用卷快照/备份）"
TestJunctionInsideAnAllowedRootCannotReachAnAListFile:
  穿过 junction 读到了 A 档文件: {Text:OPENSSH-KEY-MATERIAL-DO-NOT-LEAK IsError:false ...}
TestWavedThroughJunctionWriteBooksAnApprovedRowThatWroteNothing:
  桥在人批准之后把 junction 对面的文件写了 ...
```

⇒ 这些用例**真的在看着 junction**，不是"因为什么都被拒所以通过"。
8.3 的两条在这次变异下**没红**（诚实记着）：这条变异只吞工具层的第二次解析，
而 8.3 短名的拒绝发生在**判定层**（R2/R3 自己会 Canonicalize），两者不是同一条路径。

还原核对：`grep -rc MUTATION-20 internal/tools/` → 0；`git status internal/tools/` 只剩我的新测试文件；
复跑全绿（上面那串数字就是还原后跑的）。变异只活了几十秒，动手前 `tasklist` 确认无 `go.exe`/`wisp.exe`/`balldebug.exe`。

## 4. 同形排查：桥层"接受调用方给的路径"的入口清单

逐个查过的调用点（文件:行）与是否过 C26 `PathCanonicalizer`→`risk.Resolve`：

| # | 入口 | 位置 | 过不过解析器 |
|---|---|---|---|
| 1 | 桥的统一路径抽取 | `internal/tools/bridge.go:245`（`pathArgs(params, entry.Decl.PathParams)`）→ `bridge.go:247` `Assess` | ✅ R2 内部 `ctx.canon.Canonicalize`（`internal/risk/rules_gateway.go:33`） |
| 2 | 卡面/审计行的路径 | `internal/tools/bridge.go:622-636` `displayPaths` | ✅ `b.paths.Canonicalize`，失败时打"无法规范化"标记 |
| 3 | `fs.read` | `internal/tools/fs.go:140` → `fs.go:91-103` `open` | ✅（且 `os.Open(canon)` 用的是解析后的路径，`fs.go:144`） |
| 4 | `fs.list` | `internal/tools/fs.go:204` → `open`；列条目用 `joinForListing`（`fs.go:278`，父目录已 canonical、条目名来自 `Readdirnames` 不含分隔符） | ✅ |
| 5 | `fs.write` | `internal/tools/fs_write.go:252` → `canonical`；暂存目录 `dirOf(target)`（`fs_write.go:268/277`） | ✅ |
| 6 | `fs.trash` | `internal/tools/fs_write.go:375` → `shellTrash(target)`（`recycle_windows.go`） | ✅ |
| 7 | `fs.move` | `internal/tools/fs_write.go:446`（from）与 `:450`（to）→ `canonical`，跨卷分支 `:484` 用的仍是这两个值 | ✅ |
| 8 | `fs.delete` | `internal/tools/fs_write.go:594` → `canonical` 后才 `os.Remove(target)` | ✅ |
| 9 | R8 覆写探针 | `internal/tools/fs_write.go:624-638` `writeFacts`：先 `canonical` 再 `existsViaLstat`；解析失败按"是覆写"处理（贵的一侧） | ✅ |
| 10 | R8 跨卷探针 | `internal/tools/fs_write.go:647-670` `moveFacts`：from/to 都先 `canonical`，任一失败 → `Irreversible:["delete"]` fail-closed | ✅ |
| 11 | 存在性探测 | `internal/tools/fs_write.go:200-210` `existsViaLstat`：入参声明只接受 canonical | ✅（上游是 9/10） |
| 12 | spill / artifacts（宿主内部） | `internal/agent/spill.go:97-106` `filepath.Join(s.dir, artifactName(callID, seq))` | ❌ **不过 C26**——但调用方能给的只有 `callID`，`artifactName`（`spill.go:132-146`）把它裁成 `[A-Za-z0-9_-]`，**分隔符/盘符/`..` 进不来**，`s.dir` 是宿主给的（`loop.go:238` ← `Config.ArtifactsDir`）。⇒ 结论：不是守卫缺口，但它是"接受调用方输入的路径入口里唯一不走 C26 的那个"，本条按硬规矩**只上报不修** |
| 13 | artifacts 面板侧删除 | `internal/memory/artifacts.go:104-146` `DeleteArtifact` + `validArtifactName` | ❌ 不过 C26，用**裸名校验**（拒 `/` `\` `:` `.` `..`）+ `filepath.Join`。同一形状的第二个入口，同样**只上报** |

13 个入口：**fs 全套 11 个都过解析器（判定层与执行层各一次）**；两条 artifacts 路不过，
但它们今天**不接受调用方给的路径**（一个靠字符白名单、一个靠裸名校验）。
我没有给任何一处加新守卫（本项目硬规矩：加守卫要先讲清放哪层，由 owner 判）。

### 上报给 owner 的三条（都不是我这一框能顺手改的）

1. **`<allowed>\jn\<A 档文件>` 今天走的是"可批准的 L2"，不是 Deny。** 解析器拒绝规范化 ⇒
   R3 拿不到长路径，只能 fail-closed L2，卡面上目标显示成 `<原样> (无法规范化: …)`——
   **批准它的人看不见自己批准的是什么**。字节仍然落不下去（第二次解析挡住了），
   但"红队拒绝该不该在判定阶段就是不可批准的"是安全判定。用例已把今日形状原样钉住
   （`TestJunctionInsideAnAllowedRootCannotReachAnAListFile`，注释里写明"登记不修"）。
2. **被批准的红队调用在 `tool_call` 里记成 `decision=allow / outcome=error / error_class=tool`**，
   也就是"工具说不"而不是"策略说不"（`observe.ClassPermissionDenied` 今天只用于真 Deny/能力缺失）。
   要不要给它自己的 disposition，与 A21 那条同族；已钉成
   `TestWavedThroughJunctionWriteBooksAnApprovedRowThatWroteNothing`。
3. **第 7 框 (a) 需要的"机器可辨原因"今天不存在**：两类拒绝只在**措辞**上不同
   （"路径无法规范化，按越界处理" vs "目标路径在授权目录之外"），`Decision` 上没有原因枚举，
   且 `risk.ErrReparseDenied` 在桥里被折成字符串（`bridge.go:628-631`、`fs.go:96`）。
   做第 7 框的人必须先决定原因值放哪一层。本框**没勾第 7 框**。

## 5. 没证明但看起来成立的事（明说）

- **符号链接（symlink）没测**：只测了 junction（`mklink /J`）。二者是同一个
  `FILE_ATTRIBUTE_REPARSE_POINT` 分支（`internal/risk/pathresolver_windows.go:87`），
  但创建 symlink 需要特权或开发者模式，我没有造，**没有证据**。
- **`\?\` 与 UNC 在桥层**：票 18 的 Case 3/4 覆盖了 resolver，桥层本文件**没重做**——
  它们不是 reparse 拒绝而是规范化成功的路径，形状与 §2 第 5 行的 8.3 越界同族。
- **硬链接（hard link）不在任何一层的产物清单里**：它不是 reparse point，C26 看不见它。
  本框没测，也不知道有没有测过；若它算缺口应新开条目，不该塞进这一框。
- 跨卷 junction（junction 指向另一卷/网络路径）未测。
