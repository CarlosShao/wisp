# 票 103 在 HEAD 的复核读数 —— AC#1（R-c 密封缝）为第一件可做 AC，本腿判"已由历史提交实现，无需新增生产码"

**腿：** 实现腿（票 103 派单；票面 `Status: rejected-needs-fix`，退回单 = 票 108）
**交件时的树：** branch `dev` @ `d26c04988`（派单写的是 `0d809934e`，已用 `git merge-base --is-ancestor 0d809934e HEAD` 核过：`YES_ANCESTOR`，两枚之间的漂移是别人的收 wave 提交，与本票无关）
**写面核对（开工前实测）：** `git status --porcelain internal/winsec/ internal/memory/` 空 ⇒ 没有人的未提交改动挡路，本腿也没有覆盖任何东西。
**环境：** `export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP`；`go version go1.27.1 windows/amd64`；全程未用 `CGO_ENABLED=0`。

---

## 0. 一句话结论

票 103 的 AC#1 在 HEAD 上**已经是绿的、且绿在守卫上**：修前红 + 修后绿 + "把守卫去掉就红"三组读数都能量到（下面 §3/§4/§5）。
实现它的生产码是历史提交 `0717bf296`（票 103 自己那枚），而票面被退回的三枚形状（P1b 注册顺序、P2 分隔符、P3 底线同错法）
已分别由票 **108**（`c02c609f7`，`-done`）与票 **113**（`accepted-done`）闭合；票 108 的票头写着 `Blocks: 票 103 结案`。
因此本腿**没有新增一行生产码，也没有放松任何断言**——再写一遍守卫只会与 `0717bf296`/`c02c609f7` 重复，属于本仓明确禁止的"改松/重写"方向。
本腿的交付物是这份可复跑的复核读数，供裁决者把 AC#1/AC#2 重新判一次。

**AC 框本腿不自勾**（裁决表归验收方），票面文件名与 `Status` 一字未动，未 push。

---

## 1. 票面读法与准入（逐条，可核）

| 项 | 读数 | 结论 |
|---|---|---|
| 票面状态 | `rejected-needs-fix`（非 `-done`、非 withdrawn） | 可干 |
| `Blocked by` | `nothing` | 无前置 |
| 退回单 | 票 108，其票头 `Blocks: 票 103 结案`，`Status: done`（A1537） | 阻塞物已清 |
| 写面 | `internal/winsec/**` + `internal/memory/**` + focused `_test.go` + `docs/evidence/s1/103-*.md` | 本腿只写最后一项 |
| 实机/前端项 | 无（本票纯 `internal/`；junction 只在临时目录造） | 不触发停机条件 |

---

## 2. 开工时发现的形状（决定"第一个可做 AC 是什么"）

- `internal/winsec/resolve.go` 的 `SetPathResolver` 已是**单向闩锁**：`nil` 在已安装/曾安装后被拒并写 ERROR 审计（`resolve.go` 的 nil 分支），
  第二个**不同类型**的解析器被拒，**同类型**重装保留第一份（幂等、只窄不宽）；安装前跑 `resolverConformanceFailure`（敌意形状探针 + 树归属探针）。
- `ResolvePath` 对已装 resolver 交回的**答案**再过一遍 `builtinVerifier`，并拒收"账上记为改写"的答案。
- `RemoveUnlinked` 先 `filepath.IsAbs` 拒相对拼写，再 `firstLinkAncestor`（`pathPieces` 按平台分隔符切分，Windows 上 `/` 也算）拒"祖先链里有链接"。
- `internal/memory/removeStray` 仍是 `filepath.WalkDir`（Lstat 基、不下降），tripwire 用例在位。

⇒ AC#1 与 AC#2 都已有实现。按"第一个可做 AC = AC#1"，本腿要回答的是：**它今天还成不成立**（不是"再实现一遍"）。

---

## 3. RED-before：修前树真的红（仓外快照，`git archive 184af2222`）

快照 `C:/Users/swq/AppData/Local/Temp/wisp103-impl-<pid>`（带会话后缀，仓内未建 worktree、未 checkout）。
该提交的 `SetPathResolver` 原文（`grep -n "func SetPathResolver" -A6`）：

```go
func SetPathResolver(r C26Resolver) {
	resolverMu.Lock()
	defer resolverMu.Unlock()
	resolver = r
}
```

`go test ./internal/winsec/ -run "TestAC1|TestAC2" -count=1 -v` ⇒ **rc=1**：

```
--- FAIL: TestAC1SeamRejectsARubberStampAndLeavesTheForeignDaclAlone
--- FAIL: TestAC1SeamIsSingleUse
--- FAIL: TestAC2RemoveUnlinkedRefusesAPathThroughAJunction
FAIL	github.com/CarlosShao/wisp/internal/winsec	6.506s
```

同一次运行里绿着的对照组（点名，避免被读成"全红"）：
`TestAC1RefusedInstallLeavesTheSealWorking` PASS、`TestAC2RemoveUnlinkedStillUnlinksAStandaloneLink` PASS、
`TestAC1BaselineProductionPaths` / `TestAC1BaselineForeignACEPropagatesIntoModeOnlyWrites` /
`TestAC2SealedDirCoversFilesItNeverTouched` / `TestAC2SealDirNarrowsChildrenThatCarryTheirOwnExplicitACEs` /
`TestAC2MigrationBackupIsPrivate` / `TestAC2PreExistingMigrationBackupIsRepaired` 全 PASS。

---

## 4. GREEN-after：同一批用例在 HEAD 上绿

```
go test ./internal/winsec/ -run "TestAC1Seam|TestAC1Refused|TestAC2RemoveUnlinked" -count=2 -timeout 420s -v
--- PASS: TestAC1SeamVouchesForTwoRealObjectsThatDifferOnlyInAbsoluteness (0.01s)
--- PASS: TestAC1SeamCannotBeFreedThenGivenATreeMovingFake (2.96s)
--- PASS: TestAC1SeamRejectsARubberStampAndLeavesTheForeignDaclAlone (0.39s)
--- PASS: TestAC1SeamIsOneWayUseIsTheOnlyDirection (0.00s)      # 票 108 把 TestAC1SeamIsSingleUse 改成这个名字
--- PASS: TestAC1RefusedInstallLeavesTheSealWorking (0.31s)
--- PASS: TestAC2RemoveUnlinkedRefusesAPathThroughAJunction (0.16s)
--- PASS: TestAC2RemoveUnlinkedStillUnlinksAStandaloneLink (0.21s)
（以上 7 条 -count=2 第二样本同色，逐条 PASS；）
ok  github.com/CarlosShao/wisp/internal/winsec  5.611s
```

注：票面用例名 `TestAC1SeamIsSingleUse` 在 HEAD 已不存在（`grep -rn "func TestAC1SeamIsSingleUse"` 空），
现名是 `TestAC1SeamIsOneWayUseIsTheOnlyDirection`（`seam_guard_windows_test.go:242`），所以 §3 的修前红用的是旧名、§4 的修后绿用的是现名。

SID 级判据不是"什么都没发生"：用例自己打 `icacls` 前后两遍，读数里 `S-1-1-0`（Everyone）在密封后**仍在**外来文件上
（`victim DACL before: sids=[S-1-1-0 S-1-5-32-544 S-1-5-18 S-1-5-21-...-1001]`）。

---

## 5. MUT-1（HEAD 快照，只改守卫两行）⇒ 红在守卫、witness 仍绿

```bash
SNAP=.../wisp103-mut-<pid>                 # git archive HEAD
sed -i 's/if reason != "" {/if reason != "" && false {/; s/if seamLatched {/if seamLatched && false {/' internal/winsec/resolve.go
grep -n 'if reason != ""\|if seamLatched' internal/winsec/resolve.go
157:	if reason != "" && false {
164:	if seamLatched && false {
go build ./internal/winsec/   # rc=0（编译失败不算变异，先量到）
```

`go test ./internal/winsec/ -run "TestAC1Seam|TestAC1Refused" -count=1 -v` ⇒ **4 FAIL / 1 PASS，rc=1**：

- `--- FAIL: TestAC1SeamVouchesForTwoRealObjectsThatDifferOnlyInAbsoluteness`（票 129 的缝案例，随守卫一起红）
- `--- FAIL: TestAC1SeamCannotBeFreedThenGivenATreeMovingFake`
  `AC#1 RED: after the nil reset the guard ACCEPTED a fake that rewrites every answer to ...victim\sub\keep-me.txt`
  `AC#1 RED: S-1-1-0 (Everyone) was stripped from the foreign file, i.e. its DACL was rewritten`
  `AC#1 RED: foreign file DACL changed: before=[S-1-1-0 ...] after=[S-1-5-18 S-1-5-32-544 S-1-5-21-...-1001]`
- `--- FAIL: TestAC1SeamRejectsARubberStampAndLeavesTheForeignDaclAlone`
  `AC#1 leg 1 RED: a pass-through rubber stamp was installed into the sealing seam`
- `--- FAIL: TestAC1SeamIsOneWayUseIsTheOnlyDirection`
  `AC#1 RED: the seam is not single-use - winsec_test.narrowOnlyResolverB replaced *winsec_test.countingResolver`
- `--- PASS: TestAC1RefusedInstallLeavesTheSealWorking` ⇒ **AC#3② 要的"注册了但被拒仍有一条绿"成立**，红在守卫不在噪声。

还原：`git show HEAD:internal/winsec/resolve.go > $SNAP/internal/winsec/resolve.go`，
`diff -q <(git show HEAD:internal/winsec/resolve.go) $SNAP/internal/winsec/resolve.go` 无输出 ⇒ **CLEAN**。

---

## 6. 验收当时把票 103 退回的三枚形状，今天在 HEAD 的读数

| 形状 | 用例 | 读数 |
|---|---|---|
| P1b 先 `nil` 释放再装"恒改写"伪造 | `TestAC1SeamCannotBeFreedThenGivenATreeMovingFake` | PASS（3.49s） |
| P2 全 `/` 或混合分隔符绕过祖先检查 | `TestAC2AncestorGuardHoldsForEverySeparatorSpelling` | PASS，8 条子形状逐条：native-backslash / all-forward-slash / mixed-separators / trailing-separator / doubled-separator / dot-segment / extended-length-prefix / volume-only-relative-tail |
| P2 阳性对照（普通文件照删） | `TestAC2AncestorGuardStillUnlinksAPlainFile` | PASS |
| P3 `platformVerifyPlacement` 同错法 | `TestAC3PlacementFloorHoldsForEverySeparatorSpelling` | PASS，同一组 8 条子形状 |
| R-103-1 后半"答案的树归属无人查" | `TestAC4TreeOwnershipIsPartOfTheConformanceContract`、`TestAC4ResolverThatCannotAccountForItsTreeCannotSealAnything` | 全 PASS |
| AC#2 memory 侧 tripwire | `TestAC2ReclaimWalkNeverDescendsIntoAJunction` | PASS（0.40s），`ok memory 0.510s` |

`go test ... -run "<上表七名>" -count=1 -v` ⇒ `ok winsec 9.219s / ok memory 0.510s`，0 FAIL。

---

## 7. AC#4 门（真实树，全部在本腿复跑）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | **rc=0**（用裸 `; echo $?` 量的，不是管道尾值） |
| 回归 | `go test -count=2 ./internal/winsec/ ./internal/memory/` | `ok winsec 80.414s` / `ok memory 95.551s`，**FAIL=0** |
| 逐条点名 | `go test -count=1 -v`（同两包） | **320** 条 `=== RUN`；`--- SKIP` 恰 **2** 条：`TestSubprocessCrashWriter`、`TestAppendSessionTurnCrashChild`（均为既有环境条件用例）；`--- FAIL` **0** 条 |
| 格式 | `/d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/ internal/memory/` | **空** |
| vet 原生 | `go vet ./internal/winsec/ ./internal/memory/` | rc=0 |
| vet linux | `GOOS=linux go vet ./internal/winsec/ ./internal/memory/` | rc=0 |
| vet darwin | `GOOS=darwin go vet ./internal/winsec/` | rc=0（**只有编译期**，见 §9） |
| D22 | `sh scripts/d22scan.sh` | **rc=0**，`clean - no D22 ban violations`；各 scope 计数 bans #1-5 internal/=363、cmd/=58、#6 frontend/=116、#7 internal/tools/=48、#8 design/=39、frontend/=116、internal/=904、cmd/=211（不降） |

---

## 8. 本腿做了什么、没做什么（逐条，免得被读成"全做"）

**做了（可复跑）**：§3 的修前红、§4 的 HEAD 绿、§5 的 HEAD 变异四红一绿与还原、§6 的七名反探针、§7 的八道门。

**没做 / 没碰**：
- 未新增或修改任何生产码、任何既有用例、任何 SKIP 条件、任何断言（`git status --porcelain internal/winsec/ internal/memory/` 仍空）。
- 未动 `internal/risk/**`、`tools/d22scan/**`、`allowlist.txt`、`thresholds.go`、golden、`docs/PLAN.md`、`docs/specs/**`、`frontend/**`、`design/**`。
- **AC#3③（tripwire 的"不下降"改成"下降"）本腿未重测**：该读数由票 103 实现方与 `acceptor-ticket103` 各自测过（`artifacts.go` 换成跟随链接的递归 ⇒ `TestAC2ReclaimWalkNeverDescendsIntoAJunction` FAIL），
  本腿只复跑了它的绿侧（§6），不冒充"我重测过"。
- 未勾 AC 框、未改票面、未把票改名 `-done`、未 push、未在仓内建 worktree / checkout。
- 环境观察（与票无关、未动）：`internal/winsec/ৠ旵翺` 是一个**空目录**（git 不跟踪空目录，`git status` 里不出现），像是某轮用例的残留。

---

## 9. 仍然开着的账（不是本腿能结的，如实登记）

- **R-103-5（计数口径）**：票面交件写"SKIP 恰 1 条"，本腿实测是 **2** 条（多出 `TestAppendSessionTurnCrashChild`，票 108 之后才出现的环境条件用例）；`=== RUN` 计数按"每轮"还是"两轮合计"也仍无统一口径。
- **R-103-6（darwin）**：`GOOS=darwin go vet` 只是编译期；本机无 macOS runner，**运行期读数仍然没有**。
- **R-103-7（macOS `/tmp`、`/var` 是 symlink）**：票 113 已在 POSIX 侧接上同一条链接腿，代价是"数据根真挂在别人链接底下时改拒"；本腿未复核 macOS 行为。

---

## 10. 给裁决者

1. 若按本表判 AC#1：建议依据 §3（修前红）/§4（修后绿）/§5（去守卫即红、witness 绿）/§6（P1b 形状被拒）；**AC#1 的两条腿与 SID 级判据都齐**。
2. 若按本表判 AC#2：建议依据 §4（junction 拒 + 不删目标）/§6（8 条分隔符子形状 + tripwire 绿）；AC#3③ 的变异读数请引 `docs/evidence/s1/103-adversarial-acceptance.md` 与票 108 的交件，不要记在本腿名下。
3. 票 108 已 `done`（A1537）、票 113 已 `accepted-done` ⇒ 票 103 的"退回"不再有未清的前置；`R-103-5/6/7` 是口径与平台账，建议按票 108 的做法单开或登记，不要压在本票结案上。

**复跑最短路径（复制即用）**

```bash
cd "/d/work/workspace/projects plans/Wisp"
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP
go build ./cmd/... ./internal/... ./tools/... ; echo "BUILD_RC=$?"
go test -count=2 ./internal/winsec/ ./internal/memory/
go test ./internal/winsec/ -run "TestAC1Seam|TestAC1Refused|TestAC2AncestorGuard|TestAC3PlacementFloor|TestAC4" -count=1 -v
go test ./internal/memory/ -run "TestAC2ReclaimWalkNeverDescendsIntoAJunction" -count=1 -v
/d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/ internal/memory/
go vet ./internal/winsec/ ./internal/memory/ ; GOOS=linux go vet ./internal/winsec/ ./internal/memory/ ; GOOS=darwin go vet ./internal/winsec/
sh scripts/d22scan.sh ; echo "D22_RC=$?"
```
