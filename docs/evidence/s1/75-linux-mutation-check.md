# 票 75 — Linux 变异对照（独立重跑，证明 POSIX 判据真的咬得住）

**测量者**：本代理（2026-09-21 重发第 3 次派发；前两次死于平台错误 `Get API key for model failed`，零工具调用）。
**快照**：`git archive 17efc2c`（仓外解包，**未建 worktree**，未碰工作树，未 commit）。
**平台**：docker `golang:1.27` = `go version go1.27.1 linux/amd64`，容器 Debian GNU/Linux 13 (trixie)，
内核 `Linux 6.6.114.1-microsoft-standard-WSL2 ... x86_64`。
**日志目录**：`/tmp/wisp75logs/`（`baseline.log` / `mutant.log` / `restored.log`，均在仓外）。

---

## 1. 快照与容器（命令 + 真实读数）

```
mkdir -p /tmp/wisp75m && git archive 17efc2c | tar -x -C /tmp/wisp75m     # ARCHIVE_OK exit=0，763 个文件
docker volume create wisp75m-gomod ; docker volume create wisp75m-gobuild
docker run -d --name wisp75m \
  -v "C:/Users/swq/AppData/Local/Temp/wisp75m:/src" \
  -v wisp75m-gomod:/go/pkg/mod -v wisp75m-gobuild:/root/.cache/go-build \
  -w /src golang:1.27 sleep infinity                                      # EXIT=0，容器 id 63efa7bb088a
docker exec wisp75m go version                                           # go1.27.1 linux/amd64
docker exec wisp75m go build ./internal/risk/ ./internal/tools/           # BUILD_EXIT=0（首次拉依赖，1m54.3s）
```

宿主机 `/tmp` = `C:\Users\swq\AppData\Local\Temp`，故 bind 源写成 Windows 形式。
模块缓存与 build cache 用命名卷，不进快照树。
`pathresolver.go` 是 LF（`grep -c $'\r'` = 0），与仓库检出无差异。

**改动前留底（用于步骤 5 对拍）**：
`cp /tmp/wisp75m/internal/risk/pathresolver.go /tmp/wisp75m.orig.pathresolver.go`
→ 两份 md5 同为 `87496d505e058435f76ccbdd1d5431be`。

---

## 2. 基线（改动**之前**，Linux 真实读数）

```
docker exec wisp75m go test ./internal/risk/ ./internal/tools/ -count=1 -v
BASELINE_EXIT=1        real 0m40.216s
```

判据命令就是票面的 `go test ./internal/risk/ ./internal/tools/ -count=1`，本代理只多加了 `-v`
（只为逐条取名，不改判定行为；退出码语义不变）。

计数（`grep` 实测，非估算）：

| 口径 | 全包 | internal/risk | internal/tools |
| --- | --- | --- | --- |
| `=== RUN` | **169** | 113 | 56 |
| 顶层 `--- PASS` | **101** | 59 | 42 |
| 顶层 `--- FAIL` | **8** | 8 | **0** |
| 顶层 `--- SKIP` | 3 | 1 | 2 |
| 缩进子项 `--- PASS` | 57 | 45 | 12 |
| 缩进子项 `--- FAIL` | 0 | 0 | 0 |

包尾两行（逐字）：

```
FAIL	github.com/CarlosShao/wisp/internal/risk	2.935s
ok  	github.com/CarlosShao/wisp/internal/tools	9.588s
```

`--- FAIL` 全量 8 条，逐字名字：

1. `TestWriteGateNotSelectedByPayloadKey`
2. `TestWriteGatePlainLocalWriteNotFlagged`
3. `TestWriteGateEveryPathTargetJudged`
4. `TestWriteGateAllPathsNonSyncStaysExempt`
5. `TestSyncFallbackNotDisarmableByWeakRoot`
6. `TestSyncEnvConfiguredRoots`
7. `TestSyncNormalNewFileWriteNotFlagged`
8. `TestSyncDotDotTailFailsClosed`

⇒ **与票面预期一致**：`internal/tools` 全绿（0 FAIL、无超时），`internal/risk` 只剩 8 条 sync-root /
"confirmed root / registry-grade root" 家族（它们的失败信息全是
`provenance_test.go:469 / syncdirs_test.go:121 ...: precondition: the injected registry-grade root must confirm detection...`
那一族 ⇒ 记在票 55/票 82 的账上，**不是票 75 的**）。

---

## 3. 变异（把守卫退回"无条件折叠"）

定位（改动前）：

```
$ grep -n -B2 -A4 'not a UNC spelling' internal/risk/pathresolver.go   # 容器内同路径同读数
167-  func normalizeLocalUNC(p string) string {
168-  	if !strings.HasPrefix(p, `\\`) && !strings.HasPrefix(p, `//`) {
169: 		return p // not a UNC spelling: this function has no business here
170- 	}
171- 	u := strings.ReplaceAll(p, "/", `\`)
```

动作 = **删掉 168–170 这三行早返回**，让 `strings.ReplaceAll(p, "/", `\`)` 对所有输入生效
（与 `27c6fe5^` 的史前版本逐字同形，已用 `git show 27c6fe5^:internal/risk/pathresolver.go` 对过：
史前函数体就是这个循环 + `return u`，没有守卫）。

证明改动真落地 + **证明它编译得过**（同一条 `&&`/`set -e` 链，避免"编译失败当变异"）：

```
docker exec wisp75m sh -c 'set -e; cd /src;
  grep -c "not a UNC spelling" internal/risk/pathresolver.go || echo 0;   # → 0（守卫已消失）
  grep -n -A2 "^func normalizeLocalUNC" internal/risk/pathresolver.go;     # → 167 func / 168 u := ReplaceAll / 169 for
  md5sum internal/risk/pathresolver.go;                                    # → b4a6db06c713e7be88056ab0cb7328b2
  go build ./internal/risk/ ./internal/tools/ && echo BUILD_RC=$?'          # → BUILD_RC=0
```

变异体 md5 `b4a6db06c713e7be88056ab0cb7328b2` ≠ 底稿 `87496d505e058435f76ccbdd1d5431be`，且编译 rc=0 ⇒ 是**行为变异**。

⚠ 变异面比"整票回退"更**窄**：只退回 AC#4 交接段落里的 (1)（守卫），
(2) `sepStr`/`unifySeparators` 与 (3) `tailExistsBelow` 的 `sepStr` 重拼接**原地保留**。
⇒ 下面测出来的红数是**下界**。

---

## 4. 变异后重跑同一条命令（Linux）

```
docker exec wisp75m go test ./internal/risk/ ./internal/tools/ -count=1 -v
MUTANT_EXIT=1        real 10m10.991s
```

| 口径 | 全包 | internal/risk | internal/tools |
| --- | --- | --- | --- |
| `=== RUN` | **166** | 113 | 53（**3 条从未跑到**，被 panic 截断） |
| 顶层 `--- PASS` | 66 | 45 | 21 |
| 顶层 `--- FAIL` | **40** | **22** | **18** |
| 顶层 `--- SKIP` | 2 | 1 | 1 |
| 缩进子项 `--- PASS` | 41 | 32 | 9 |
| 缩进子项 `--- FAIL` | **16** | 13 | 3 |

包尾逐字：

```
FAIL	github.com/CarlosShao/wisp/internal/risk	2.707s
panic: test timed out after 10m0s
	running tests:
		TestVetoInsideTheWindowWritesNothing (4m48s)
FAIL	github.com/CarlosShao/wisp/internal/tools	600.135s
```

### 4.1 新增的红（`comm -13 baseline.mutant` 实测，逐字）

`internal/risk`：**+14 条顶层 + 13 条子项**（基线 8 条 → 变异 22 条）

顶层：`TestResolveCanonicalIsPlatformShaped`、`TestResolveKeepsSeparatorsOutOfPosixNames`、
`TestDoubleSlashSpellingIsNotUNC`、`TestTierAAnchorsMatchNativeSpelling`、
`TestOverrideKeyKeepsPosixBackslashesDistinct`、`TestExfilSyncWritePosixSpellingInvariant`、
`TestFourChannelExfilSuite`、`TestResolvePerCallBudget`、`TestSyncAncestorWalkVerifiesNativeChain`、
`TestSyncFixtureFallbackAndMatch`、`TestSyncUnverifiedRootKeepsFallback`、
`TestSyncSuspectFallbackIsComponentBounded`、`TestSyncSuspectFallbackWhenUndetectable`、
`TestSyncDetectionOnThisMachine`
子项：`TestTierAAnchorsMatchNativeSpelling/{~/.ssh/**, ~/.ssh/_dir_itself, ~/.git-credentials,
~/.aws/credentials, ~/.kube/config, .git/config, %APPDATA%\wisp\config.toml,
%APPDATA%\Microsoft\Protect\**, %LOCALAPPDATA%\Microsoft\Credentials\**, browser_credential_store}`（10 条）
+ `TestExfilSyncWritePosixSpellingInvariant/{write_into_sync_root, backslash_inside_a_POSIX_file_name}`
+ `TestFourChannelExfilSuite/fs.write_into_sync_dir`

`internal/tools`：**+18 条顶层 + 3 条子项**（基线 0 条 → 变异 18 条，另有 3 条根本没轮到跑）

顶层：`TestDeclaredRiskIsOnlyAFloor`、`TestToolCallRowsAreComplete`、`TestFSReadReturnsTheFileAndTaintsIt`、
`TestFSReadTaintFeedsR4`、`TestFSListSummarizesADirectory`、`TestFSListHonoursItsCap`、
`TestEmptyAllowlistAuthorizesNothing`、`TestSensitiveFileIsDeniedNotEscalated`、`TestD34WriteMatrix`、
`TestOverwriteDetectionFollowsTheCanonicalPath`、`TestAtomicWriteKillsMidWrite`、
`TestFSTrashGoesToTheRecycleBin`、`TestDeleteEnabledRegistersItAsL2`、`TestFSMoveSameVolume`、
`TestCanonicalizeReturnsAPathTheOSCanOpen`、`TestCanonicalizeAgreesWithTheOSName`、
`TestLoopPassesDeclaredL1WriteThroughTheGate (10.09s)`、`TestL1WriteGoesThroughTheRealBlockWindow (300.05s)`
子项：`TestAtomicWriteKillsMidWrite/{new_file_target_does_not_appear_at_all, a_clean_write_lands_and_round_trips}`、
`TestDeleteEnabledRegistersItAsL2/approved_it_deletes_and_says_it_is_permanent`

被 600s panic 截断、**一次都没执行到**的 3 条（`comm -23` 对两包的 `=== RUN` 名单实测，基线 56 条 → 变异 53 条）：
`TestLateVetoRendersTheApprovalLayersAppliedStepsReport`、`TestFSReadOnlyNeverOpensACard`、
`TestApprovedL2ThatStopsMidWriteIsNotRenderedAsNeverStarted`。
另有 1 条 `TestVetoInsideTheWindowWritesNothing` 起跑但**卡在审批窗口里被闹钟打断**（既无 PASS 也无 FAIL），
见 4.3。⇒ 这 4 条**不能被计入"变异新增的红"**，只能记作"未定（被截断）"，故本表 tools 的 18 是去掉这 4 条后的硬读数。

### 4.2 反向检查：有没有"基线红在变异后反而绿了"

```
comm -23 baseline.failnames mutant.failnames   # risk 与 tools 各一次
```

两包输出**均为空** ⇒ 变异没有把任何一条既有的红洗绿（不存在"用新红换旧红"的假象）。

### 4.3 600s 超时按名字/按等待归因（不是"它变快了"那种含糊话）

- `TestL1WriteGoesThroughTheRealBlockWindow (300.05s)` —— `wiring_test.go:124: the window ran out unopposed,
  which MEANS execute: {Text:审批超时（300 秒未确认），C18 一律判拒绝，已自动拒绝 ...}`：
  审批卡没等到投递，真等满 C18 的 300 秒墙钟。
- 随后 `TestVetoInsideTheWindowWritesNothing` 卡 **4m48s** 被 10m 闹钟打断，panic 的 goroutine 栈点名了等待处：
  `approval.(*Gate).PendingApproval` @ `internal/agent/approval/gate.go:473`
  ← `internal/tools/bridge.go:340` ← `wiring_test.go:162`；
  它前一行日志就是 `wiring_test.go:158: veto: approval: correlation_id 无对应待审批项`
  ⇒ 待审批项的 key（canonical 路径）在 POSIX 上被折成反斜杠串，veto 找不到它，于是永远等在 `PendingApproval`。
- 300.05s + 288s ≈ 588s ⇒ 正好撞上 `test timed out after 10m0s`（tools 包 600.135s）。
  这条与票 70 取证的 CI 症状同形（原话是 19 FAIL + 600s panic；本代理在只回退守卫的下界变异下测到 18 FAIL + 600s panic + 3 条未跑到，
  数目差就差在那 3 条根本没轮到，方向一致）。
- 顺带一条纯形状证据（`fs_test.go:19`）：
  `open /src/internal/tools/\tmp\TestFSReadReturnsTheFileAndTaintsIt2762579392\001\note.txt: no such file or directory`
  —— 一个没有任何 POSIX 调用能打开的串，正是票 75 描述的病灶。

### 4.4 判据咬人的直接证据（fail-open 类，值得 owner 看一眼）

- `pathshape_portable_test.go:246`（`TestOverrideKeyKeepsPosixBackslashesDistinct`）：
  `CROSS-FILE BLEED (fail-open): the override recorded for "...id_k\x.pem" unlocked the different file ...`
  ⇒ 折叠一旦无条件，B 档单文件豁免会跨文件解锁，这是**安全方向**的红，不是装饰。
- `pathshape_portable_test.go:151`（`TestTierAAnchorsMatchNativeSpelling`，票 72 的不变式）：
  `classify("...\\tmp\\...\.ssh\\id_platformkey") = none, want A (rule ~/.ssh/**)` 共 10 条
  ⇒ A 表锚点全部漏判，敏感文件在 Linux 上直接掉档。
- `provenance_syncdirs_other_test.go:71`（`TestExfilSyncWritePosixSpellingInvariant`）：
  `ESCAPIABLE: "/tmp/.../OneDrive/Notes/out.md" is inside the injected sync root and was not an exfil channel`。
- `pathshape_portable_test.go:107/111`（`TestDoubleSlashSpellingIsNotUNC`）：
  `canonical = "/src/internal/risk/\\localhost\\c$\\x": a POSIX double-slash spelling was rewritten into UNC shape`
  ⇒ 守卫的**另一半**（`//` 前缀）也在被钉：POSIX 上 `//localhost/c$/x` 不是 UNC。
- 唯一一条**非形状**的连带红：`TestResolvePerCallBudget`
  （`C26 Resolve: 5416534 ns/op = 5.417 ms/op (budget 1.000 ms, 206 samples)`）——
  形状坏了之后 Resolve 走更多回退，性能门顺带爆了；归为连带，不计入形状判据本体。

---

## 5. 还原与对拍

```
cp /tmp/wisp75m.orig.pathresolver.go /tmp/wisp75m/internal/risk/pathresolver.go
diff -q /tmp/wisp75m.orig.pathresolver.go /tmp/wisp75m/internal/risk/pathresolver.go   # 无输出，rc=0 ⇒ 字节一致
md5sum 两份                                                                            # 均 87496d505e058435f76ccbdd1d5431be
docker exec wisp75m grep -n -A2 "^func normalizeLocalUNC" internal/risk/pathresolver.go # 守卫 168-170 回来了
```

还原后**再跑一遍第 2 步的同一条命令**（证明还原是真的，不是只对了 md5）：

```
docker exec wisp75m go test ./internal/risk/ ./internal/tools/ -count=1 -v
RESTORED_EXIT=1        real 0m15.566s（build cache 已热）
=== RUN 169 / 顶层 PASS 101 / 顶层 FAIL 8 / 顶层 SKIP 3 / 子项 PASS 57 / 子项 FAIL 0
FAIL	github.com/CarlosShao/wisp/internal/risk	2.703s
ok  	github.com/CarlosShao/wisp/internal/tools	9.697s
```

8 条红与基线**逐字同名**（第 2 节那 8 条），`internal/tools` 回到 `ok 9.697s`。

**工作树**：全程只在 `/tmp/wisp75m` 快照里改；未 commit、未 push。
本代理写盘时 `git status --porcelain` 里另有 4 条 `M`（票 77 / 票 83 / `docs/reports/pending-and-issues.md` /
`internal/config/unwired_test.go`）与一条别人的 `?? docs/evidence/s1/75-independent-verification.md`，
**全部一字未动**（它们随后被各自的代理提交了，`dev` HEAD 因此从 `522efec` 走到 `e4082a0` —— 与本代理无关）。
本代理对仓库的唯一贡献 = 本文件这一条，且它至今是 `??`（未跟踪、未提交）。

顺带一句交叉印证（**不是**本代理的读数来源）：HEAD 上 `e4082a0` 那笔票 70 的更正写着
"Linux 上 `test-core` 只剩 `internal/risk` 1 包 8 条" —— 与本代理第 2 节独立测得的
`risk FAIL 8 / tools ok` 数目与包分布一致。本文所有数字仍全部来自上面这三条本代理亲自跑的命令。

---

## 6. 结论

1. **票 75 的 POSIX 判据在 Linux 上真的咬得住**：只把 `normalizeLocalUNC` 的守卫退回"无条件折叠"
   （AC#4 交接段落三处改动里最小的那一处，且编译 rc=0 已自证），Linux 上立刻
   **+32 条顶层红 / +16 条子项红**（`internal/risk` 8→22，`internal/tools` 0→18），
   外加 `internal/tools` 复现票面描述的 `panic: test timed out after 10m0s`（600.135s）。
   这不是"一条都没多"，也不是转述别人的日志 —— 是本代理在同一条命令上量出来的。
2. **基线与变异之间没有"洗绿"**：`comm -23` 两包皆空，红只增不减，判据不是互相抵消的噪声。
3. **归因干净**：基线残留的 8 条全是 sync-root / registry-grade 家族（票 55/票 82 的账），
   这 8 条在变异下**一条不少地仍是红**（`comm -23` 已证无洗绿）；而同一家族里基线本就绿的成员
   （`TestSyncUnverifiedRootKeepsFallback`、`TestSyncSuspectFallbackIsComponentBounded`、
   `TestSyncSuspectFallbackWhenUndetectable`、`TestSyncAncestorWalkVerifiesNativeChain`、
   `TestSyncFixtureFallbackAndMatch`、`TestSyncDetectionOnThisMachine`）因形状坏掉而**新增**转红
   —— 两笔账在表里是分开可数的（8 旧 + 14 新 = 22）。
4. **`internal/tools` 那 19 条在 Linux 上是真红，且红得有名有姓**：`fs_test.go` / `bridge_test.go` /
   `fs_write_test.go` 直接报"打开 `\tmp\...` 这种串打不开"，`pathshape_portable_test.go` 直接报
   "C26 交回的 canonical 带外来分隔符"。⇒ 修前 Ubuntu 上 19 FAIL 的取证与本代理独立测得的结果同向。
5. **下界声明**：本变异只回退守卫一处；若连 `sepStr`/`unifySeparators`（(2)(3)）一起回退，红数只增不减。
   故票 75 判据的强度**不低于**本表。
6. 本代理未 push，AC#6（runner 可见证据）不在本单范围内；本文只交付"Linux 上判据咬得住"这一件事。

---

## 7. 现场保留（便于 owner 复核，不再是本代理的读数）

- 容器 `wisp75m`（golang:1.27）与快照 `/tmp/wisp75m` 已**原样留在原地**（容器本代理收尾时 `docker stop` 了，
  数据在命名卷 + bind 上不会丢）：`docker start wisp75m && docker exec -it wisp75m bash` 即可进现场；
  快照树已还原到 `17efc2c` 的字节（守卫在位，md5 `87496d505e058435f76ccbdd1d5431be`）。
- 三份全量日志：`/tmp/wisp75logs/baseline.log`（383 行）、`mutant.log`（468 行，含 panic 的 goroutine 栈）、
  `restored.log`。底稿 `/tmp/wisp75m.orig.pathresolver.go` 未删。
- 本代理对仓库的唯一写入 = 本文件。

