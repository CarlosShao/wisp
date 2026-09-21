# 票 79 — 对抗验收裁决表（编排者亲自复跑，2026-09-21）

**被验对象**：票 79 `artifactName` 把不同 tool-call id 折成同一个磁盘名 + `writeFileExclusive` 并不排他。
**被验 commit**：`6901aa9`（实现）+ `878d9f8`（变异与门禁记录）+ `de3781b`（clean-HEAD 复测记录）。
**验收树**：`git archive de3781b | tar -x -C /tmp/wisp79`（**仓外纯净树**，不建 worktree，A38④）
⇒ 我读的不是任何人在飞的未提交改动。

三档证据标签的含义：**〔独立复现〕**＝我亲自敲命令、读真实 exit code；
**〔日志＋归档，我抽验〕**＝数字来自代理，但我核对文件存在且字段对得上；
**〔仅自述，不背书〕**＝只有代理的话，我没验。

## 逐 AC 裁决

| AC | 判据 | 我的裁决 | 证据档 |
|---|---|---|---|
| **AC#1** | 先前折叠的不同 id 现在得到不同磁盘名，且报告里那对（`p/q` vs `pq`）是**命名子测试** | **PASS** | 〔独立复现〕 |
| **AC#2** | `writeFileExclusive` 真的排他，且**同 id 重试**的行为被定死并测死 | **PASS** | 〔独立复现〕＋我**追加**一次代理没做的变异 |
| **AC#3** | 变异：把注入式编码退回"剥掉" ⇒ AC#1 必须红；先 grep 证明落地，还原后证干净 | **PASS**（我自己重做了一遍，不采信它的日志） | 〔独立复现〕 |
| **AC#4** | Defect 2 取 (a)/(b) 之一并有测试；把 `if e.IsDir() { continue }` 装回 ⇒ 那条测试必须红 | **PASS**（取 (a) 递归计量；红→数字逐字复现） | 〔独立复现〕 |
| **AC#5** | 三门禁 + RUN 计数不变式 + **点名已知 SKIP** | **PASS** | 〔独立复现〕 |

**结论：票 79 转 `-done`。** 五框全绿，且我这一侧**没有一条判据停留在"仅自述"档**。

## 我自己跑出来的数字（全部真实 exit code）

**基线（未变异）**
```
go vet ./internal/agent/ ./internal/memory/           rc=0
gofmt -l internal/agent/ internal/memory/             0 行
go test -count=2 -v ./internal/agent/ ./internal/memory/
  rc=0 · ok agent 3.531s · ok memory 26.015s
  === RUN 共 274 行  ·  --- FAIL 共 0 行  ·  --- SKIP 共 2 行
  两条 SKIP 都是同一个名字 TestSubprocessCrashWriter
  （internal/memory/concurrent_test.go:186；-count=2 所以一个名字出现两行）
```
⇒ 代理报的 `274`、`0 FAIL`、"2 行 SKIP / 1 个名字"、`ok agent ~3s / ok memory ~25s` 与我的独立复跑**逐字对上**。
274 是偶数且 = 2 × 137，与 `-count=2` 的结构自洽 ⇒ 没有"某一步被静默跳过"。

**M-1（= 它的 AC#3，我自己下的刀）** 锚点是 `encodeArtifactID` 循环体里真正承载注入性的三行
（`b.WriteByte('%')` + 两次 `upperHex`），整段换成 `continue`：
```
grep 证明落地：MUTATION-AC3 命中 1 行，b.WriteByte(upperHex 命中 0 行
go test -count=1 -run 'TestArtifactName|TestSpilled' ./internal/agent/   EXIT=1
  --- FAIL: TestArtifactNameDoesNotFoldDistinctIDs     （:74 共 8 条折叠对）
      ids "p/q" and "pq" both name "tool-output-pq.txt"            ← 报告里那对，原样复现
      ids "p\q" and "pq" / "p/q" and "p\q" / "a:b" and "\\a\b" / "." and ".." / ".." and "" …
  --- FAIL: TestSpilledBytesSurviveANameThatUsedToCollide
  --- FAIL: TestArtifactNameRoundTripsToTheExactID     （7 条 does not decode back）
还原：grep -c MUTATION-AC3 = 0，diff 与快照字节相同，重跑 ok 0.048s
```
⚠ **一处代理计数与我不符，方向是"它少报了"**：它写"6 个命名子测试全 FAIL"，我量到 `:74` 打出的折叠对是 **8 条**。
不是缺陷（覆盖比声称的更宽），但记下来：本仓的"多少条红"要靠**全量输出仪器**，别靠印象。

**M-2（= 它的 AC#4）** 把 `if e.IsDir() { continue }` 原样装回 `listArtifactsDir`：
```
go test -count=1 -run 'TestStray|TestPurgeArtifactsReclaims' ./internal/memory/   EXIT=1
  --- FAIL: TestStraySubdirectoryCannotHideBytesFromTheQuota
      artifacts tree still holds 4106 bytes against a 20-byte quota (freed=0)
  --- FAIL: TestStraySubdirectorySharesTheQuotasLRUQueue
  --- FAIL: TestPurgeArtifactsReclaimsStraySubdirectory
      after a purge the artifacts tree still holds 1024 bytes in 3 entries
  --- FAIL: TestStrayRemovalDoesNotFollowLinks
还原：MUTATION-AC4 命中 0，artifacts.go 与快照字节相同
```
数字 `4106`/`20-byte`/`1024 bytes in 3 entries` 与它的日志**逐字一致** ⇒ 那条判据不是为绿现写的。
"500MB 上限变 2GB" 的缩小版实拍成立。

**M-3（我追加的，代理没做）** AC#2 只被它的日志覆盖，所以我把 `O_EXCL` 换成 `O_TRUNC`：
```
sed 后的真实行：os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
go test -count=1 -run 'TestSpill|TestWrite' ./internal/agent/   EXIT=1
  --- FAIL: TestWriteFileExclusiveIsExclusive
      writeFileExclusive overwrote an existing file: O_EXCL is missing again
还原：ok 1.541s
```
⚠ 过程如实记：我**第一次**下刀把 `// 注释` 写进了函数实参列表 ⇒ `FAIL [build failed]`。
**编译失败不算行为变异**（本仓规矩），所以我换成合法表达式重跑才拿到上面那条红。

## 代码侧我读过的三处（不是只跑测试）
- `internal/agent/spill.go:245`：`os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)` —— 名字不再是谎话。
- `internal/agent/spill.go:112`：`Prepare` 对 `fs.ErrExist` 走"同一个逻辑 id 的自己人"分支（temp+rename 原子替换），
  所以真排他**没有**把同 id 重试变成错误路径。
- `internal/agent/spill_name_injectivity_test.go:242-270`：超长 id 的截断路径被钉住三件事
  —— 名字有界、**必须带 `-<16 hex>` 摘要尾巴**、两个只在第 96 字节之后才不同的 240 字节 id 不得同名；
  另外"转义之外不得出现大写"（NTFS 折叠大小写）与"名字不得停在半个 `%` 上"各有一条断言。

## 它对票 76 判据的影响（我最担心的一条：借修 bug 之名松判据）
`git diff d8a62d2..de3781b` 逐行读过两份票 76 测试：
- 改的是 **4 个 `wantName` 字符串**（那些字符串**就是**本票要消掉的折叠结果，例如
  `tool-output-pq.txt` → `tool-output-p%2Fq.txt`）+ 一个 seq 退例子测试**改名并加强**：
  旧断言在庆祝"4 个 id 折进同一个名字"，新断言要求 **5 个 id 得 5 个名字**，
  且 bare/不含 `/ \ : ..` 的 containment 断言**一条没少**。
- 票 76 的四形状守卫、AST sink 审计、containment 差分、阳性对照：`artifacts_path_invariant_test.go` 只动了
  (e) 段"nested 目录必须活过 purge"（按本票判据反过来钉成"必须不活过"，并把 removed 集合逐个名字钉死）。
⇒ **判定：没有借道放宽。** 两处反向变化都是把弱断言变强。

## 它如实上报、本票不修、需要有人接手的四件事 → 登记为 A51
1. `0o600` 在 Windows 上是装饰（实测落盘 `-rw-rw-rw-`，权限由目录 ACL 继承）⇒ artifact"只有我可读"从来没成立过；
   要真做到得走 ACL，**另一票**。
2. `os.Remove` 对"指向目录的符号链接"失败 ⇒ 一个裸名 artifact 可能清不掉（它的游离子树回收路径）。
3. **文档面已经过期**：`docs/specs/SPEC-05-agent-core.md:115`、`docs/specs/SPEC-02-data-storage.md:180`
   仍写 `artifacts\tool-output-<id>.txt`；票 20 面 `:107` 那句"只留 `[A-Za-z0-9_-]`，整串剥光退回 `seq<N>`"被本票作废。
   前两个是**冻结契约文件**、票 20 面是**编排者的面** ⇒ 代理没动是对的，改法归编排者/owner（D22）。
4. 磁盘名里现在可能出现 `%` ⇒ 任何"把 artifact 路径拼进 `cmd.exe` 命令"的新路由会先撞上 `%XX%` 展开；
   以及 `tool-output-*` 的 8.3 短名在前 6 字符就分叉（短名开着才是问题）。

**next（编排者）**：A51 入 `docs/reports/pending-and-issues.md`；票 20 面 `:107` 那句由我改写（不重开票 20）；
SPEC-05:115 / SPEC-02:180 的过期措辞要不要现在动，是 D22 判定 ⇒ 挂进 Q 清单。
