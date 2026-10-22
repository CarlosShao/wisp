# 票 360 普查腿 360-a1 交付件（只读，AC#0）

- 时刻＝2026-10-22 17:2x +08｜起手锚＝本腿 `git rev-parse HEAD` 现认 `a9b4513d`（`a9b4513dd825adf151b7781b06d1072ce061d546`）；交件时父提交＝`40339d01`，本腿两笔提交＝`48a37be7`＋`a97d2ab0`。
- 本腿零 `.go` 落仓、零产码字节、零票面/台账改动；临时件全在仓外 `D:\tmp\wisp\20261022\360a1\**`。
- 承重文件在起手锚与交件时**同串**（`git hash-object` 工作树 vs `git rev-parse HEAD:path` 逐枚同；起手锚..交件父提交之间 `git diff --name-only -- '*.go'` ＝ **0 枚**，6 件漂移全在 `.scratch`/`docs`）：writeguard.go=`0ebc415b0565`／allowdirs.go=`ef54b37c094f`／manager.go=`f53d8a0749cd`／permmode.go=`6f385353ec36`／loader.go=`9100ba92e220`／settings.go=`c4e3b09b2131`／config_reload.go=`990d148b1138`／writeguard_226_test.go=`0e5895ccb2e0`／writeguard_handedit_227_test.go=`f60b308b6b2e`／always_write_no_clobber_226_test.go=`38b32a55750e`／always_write_no_clobber_227_test.go=`058328d6b21e`。

## (1) 票面前提新鲜度（逐条对 HEAD blob 现跑复认）

1. 常驻用例存在＝真。`cmd/wisp/always_write_no_clobber_227_test.go:125 func TestAC6AHandEditIsStillSeenUnderTheRealPoll`。
2. 红点 `:198`＝真且本腿复现（读数见 (6a)）：仓外导出树插无条件 `statOwnWrite` 后红句逐字落 `always_write_no_clobber_227_test.go:198`。
3. "直接未钉"的字面依据＝真。票 227 `AC#6` 判据正文（HEAD blob `:31`）逐字只点名「真轮询之下，手改仍看得见」；"Adopt"在票 227 只出现在编排者追加裁决行 `:50`/`:87`，不在任何判据正文里。
4. 来路具名更正一处（不伤判语，只伤引号）：票面写来路＝`226-v2 sec 6 第 8 条`。本腿现跑：`226-config-write-no-clobber-v2.md` HEAD blob `:170` 的第 8 条讲的是"排程先于票 223／生产侧尚无轮询者"，该文**全部六枚历史 blob** 里 "Adopt 内存" 与 "确定性用例" 命中均 **0**。"「226 写回不 Adopt 内存」确定性用例"这句逐字首见于 `227-r1-impl-part2.md:58` 与 `227-v1-verdict.md:90`。=> 承重判语实际在 `227-v1` 第 (10) 项，票面引用链应改指它。
5. 两枚来源提交在历史上＝真：`6e838c16`（227-v1）与 `ecebc7a9`（227-r1）都是 HEAD 祖先（`git merge-base --is-ancestor` 各 rc=0）。
6. 生产轮询在位＝真：`cmd/wisp/config_reload.go:123` `Spawn("watchdog","config",...,rt.configReloadTick)`，`:92 = time.Second`，`:157 rt.mgr.CheckAndReload()`。=> 票 226 那句"生产侧尚无轮询者"已过期，票 360 量这一跳时 Adopt 路**确实有生产驱动**。

## (2) `statOwnWrite` 分层名单（逐名，NOT 只给计数；raw grep 会把注释算进来）

尺＝`git grep -n statOwnWrite HEAD -- internal cmd`：原样命中 **5 行**，分层如下（5 = 1 定义 + 1 代码调用 + 3 注释）。

- 代码定义（1）＝`internal/config/allowdirs.go:157 func (m *Manager) statOwnWrite()`，体 `:157-163`：`os.Stat` 失败即 `return`（保守方向），成功则只写 `m.seenMtime/seenSize/seenValid`。=> 它认领的是 **stat**，从不认领内容。
- 代码调用（1）＝`internal/config/writeguard.go:179 m.statOwnWrite()`，位于 `mergeWrite` 第 (3) 支 `if len(diverged) == 0`（`:171`）之内；分叉支 `:182-191` 只 `slog.Warn`、NOT 认领。
- 注释提及（3）＝`allowdirs.go:150`（doc 注释，逐字 "It is called from exactly one place - mergeWrite"）／`writeguard.go:105`（doc 注释）／`writeguard_226_test.go:186`（用例前注释）。
- 历史调用点已不存在（防误读）：旧普查（`probes/223/c1` 等）写的 `allowdirs.go:116`／`permmode.go:78` 调用点是票 226 之前的形状；今 HEAD `permmode.go` 零命中，认领全收进 `mergeWrite` 一处。

## (3) 写回路径 vs Adopt 路径（逐名，行号对 HEAD blob 现跑）

写回（单键 + 落盘），入口 3 枚、汇聚 1 处：

- `internal/config/allowdirs.go:62 AddAllowedDir` / `:109 SetAllowedDirs` -> `:139 writeAllowedDirs`（`:141` 先改内存，失败 `:143` 回滚）-> `:142 m.mergeWrite(keyFSAllowedDirs, ...)`。
- `internal/config/permmode.go:70 SetPermissionMode`（`:78` 先改内存，失败 `:82` 回滚）-> `:79 m.mergeWrite(keyRiskPermissionMode, ...)`。
- `internal/config/settings.go:259 writeOneKey`（`:302 apply(m.cur)`，失败 `:304` 回滚）-> `:303 m.mergeWrite(ownedKey, ...)`。
- 汇聚＝`internal/config/writeguard.go:111 mergeWrite`：`:118 readConfigFile`（定义 `loader.go:64`）把**盘读回来当 base**，`:143 setOn(base)` 只盖自己那一枚键，`:163 SaveFile(m.path, base)`（定义 `loader.go:252`；体＝`deepCopyConfig` + `MarshalCanonical` `parse.go:164` + `atomicWrite` `parse.go:196`，**从不回读**）。
- 生产触发点＝`cmd/wisp/approval_always.go:134`（AddAllowedDir）与 `internal/perm/store.go:208`（SetPermissionMode）。

Adopt（把盘上字节读回**接管内存**），全仓只有一条：

- `internal/config/manager.go:142 CheckAndReload` -> `:160 LoadFile` 读盘 -> `:169 plan` -> `:193 plan.commit()` 写 `cur`（`:195` 同时认领 stat）。生产驱动＝(1) 第 6 条那条 tick。
- 易混形状（具名钉住，NOT Adopt）：`writeguard.go:179` 的 `statOwnWrite` 与 `manager.go:162-165` 失败支的认领都**只改 seen-stat、不改 `m.cur`**；`mergeWrite:118` 的重读是"写前读当底"，读完即用于序列化，不回灌内存。=> 写回与 Adopt 是两条互不相连的路；票 360 量的那一跳（写回之后没有 Adopt）在结构上成立。

## (4) 复现读数：写回之后内存仍持旧值（今天可复现）

台面＝`-overlay` 挂载仓探针件（物理件 `D:\tmp\wisp\20261022\360a1\overlay\zz360a1_probe_test.go`，仓内 `internal/config/` 零 `.go`；形＝票 223-v1 的 overlay 定式）。命令＝`GOFLAGS= go test -overlay ... -count=1 -v -run Test360A1 ./internal/config/`。三处独立复跑读数逐字相同：主检出（`logs\01-*`）／首棵导出树（`logs\05-*`）／完整导出树（`logs\10-*`）。

形状＝开机 `matrixBase`（size=56）-> 手改文件 `size = 60` -> 跑一次受守卫写 `AddAllowedDir`：

- `READ A2 after writeback: file.ball.size=60 mem.ball.size=56`。
- `READ A4b diffKeyPaths(memory, file-after-write) n=1 [ball.size]`。
- `READ A5 REPRO: writeback landed, memory kept the old value (file=60 mem=56)`，探针 PASS（红闸＝探针 `:44` 的 `REPRO GONE`：若内存变 60 即红，本腿 NOT 自行改口径）。
- 该发正走分叉支：`mergeWrite` 落 `writeguard.go:184` 的 WARN，逐字带 `kept_in_file_not_in_memory=[ball.size]`。

=> 票 360 的前提**今天复现**，无须停手上报。

## (5) 正控读数：Adopt 之后内存是新值（同一台面能分辨两态）

同一台面同形状，写回后多调一次 `CheckAndReload`（= 让 (3) 那条 Adopt 路走一遍）：

- `READ B0 before Adopt: file.ball.size=60 mem.ball.size=56` -> `READ B1 after Adopt: report.Hot=[ball] mem.ball.size=60 file.ball.size=60`，探针 PASS（红闸＝探针 `:83`）。
- 第三臂（干净写）＝`READ C0 clean write: file=56 mem=56`，poll 后 `report=<nil>`（`statOwnWrite` 认领 stat；文件与内存本就同值）。=> 台面对 {旧值 56 / 新值 60} 给出**不同**读数，既非恒真也非恒红。

## (6) 间接已钉的复现 + "直接"缺口的决定性一发（突变全在仓外导出树）

(6a) 复现票面"间接已钉"（＝227-v1 logs/56 同形：`writeguard.go:169` 后插无条件 `m.statOwnWrite()`）：

- 控制发（未突变，显式 `-run '^TestAC6AHandEditIsStillSeenUnderTheRealPoll$' ./cmd/wisp/`）：PASS 3.46s，rc=0。
- 突变发：`go build ./...` rc=0；同一条 `-run` **红在 `always_write_no_clobber_227_test.go:198`**（红句 "the production reload tick never booked ... within 40s"），FAIL 42.55s，rc=1。
- 落刀前锚点命中 1、落刀后命中 1 + 标记 1；还原 sha256＝`efe12661e71b940f1fcf99d722f1c25e4403431997e90041491d5eae337d47bd`，与 `git cat-file blob HEAD:internal/config/writeguard.go` 逐字同串（14930 字节）。件＝`logs\09-227-full-arm.txt`（首棵树上同形读数见 `logs\02-*`）。

(6b) 直接缺口的决定性一发（本腿自造 M-ADOPT＝写回路径**真的**接管内存：`writeguard.go:169` 后插 `m.cur = deepCopyConfig(base)`）：

- `go build ./...` rc=0；`go test -count=1 ./internal/config/` 控制发 rc=0（全绿）-> 突变发 rc=1，**红 8 枚逐名**：`TestAC11ControlSnapshotWriteIsWhatRevertsTheHandEdit`（`settings_248_test.go:206`）／`TestAC1AllowedDirsWriteKeepsAHandEditedKey`（`writeguard_226_test.go:89`）／`TestAC3AdoptionClaimsOnlyWhatThisWriteProduced`（`:224`）／`TestAC3HandAddedEntryIsPreservedButNotAppliedToMemory`（`:248`+`:259`）／`TestAC4PermissionModeWriteKeepsAHandEditedKey`（`:285`）／`TestHandDeletedScalarLineComesBackAtSchemaDefault`（`writeguard_handedit_227_test.go:50`）／`TestHandDeletedWholeSectionComesBackFromDefaultsNotMemory`（`writeguard_unnailed_341_test.go:183`）／`TestHandDeletedProviderBaseUrlComesBackAsThePresetValue`（`:327`）。
- 宿主车道同突变、显式 `-run '^TestAC1AlwaysBranchDoesNotRevertAHandEditedKey$'`：控制 PASS 2.42s rc=0 -> 突变 **红在 `always_write_no_clobber_226_test.go:119`** rc=1（件＝`logs\08-host-full-lane.txt`）。两次还原 sha256 均逐字回原。件＝`logs\07-adopt-full-arm.txt`。
- 归因（本件最承重一句）：上述 8+1 枚红的都是**别的票、别格**的用例，"写回不许把文件内容灌进内存"今天只活在它们的**顺带断言**里；**没有任何一枚用例以这一跳为题目、且自带 (5) 那样的两态正控**。=> "间接已钉、直接未钉"成立；本腿 NOT 据此说"完全没做"，也 NOT 据此说"已有直接钉"。
- 交编排者裁 `AC#1` 的代价事实：〔乙〕"产码侧让写回路径显式 Adopt"＝M-ADOPT 同形，今天**当场红 9 枚**（包内 8 + 宿主 1），且分叉支的 Adopt 与 D36 rule 1（`ConfirmLocked` 先行）正面相撞 —— 任何乙/丙形都得先交代这 9 枚。〔甲〕"只加用例"钉的形状＝(4)+(5) 这组读数（断 file=60/mem=56 + 正控 Adopt 后 =60），零产码。
- 未跑具名：M-ADOPT 下 `TestAC6A...` 未跑（40s 一发，裁甲/乙不依赖它）；(6a) 突变下"226 宿主枚仍 PASS"沿 227-v1 记录，本腿未复跑。

## (7) 门禁各 rc（主检出；`GOFLAGS=`／unset GOTMPDIR/TMPDIR/TEMP/TMP／PATH 抬 sherpa-onnx+build）

- `go build ./...` rc=0。
- `sh scripts/d22scan.sh`（Git Bash）rc=0，末行逐字 `d22scan: clean - no D22 ban violations`（分母 `examined 289 production Go files` 是文件数，NOT 违规数）。
- `cmd/wisp` 全部读数只用显式 `-run` 单枚（票面 `A1126`/`A1127` 口径），无整包名册、无死点需登记。

## (8) 提交号 + 写面尺 + 正控尺

- 本腿两笔提交：`48a37be7`（件本体，父 `40339d01`）＋`a97d2ab0`（补记提交号，零读数改动）。写面尺 `git show --name-only --format=` 逐笔：各件数 1、`.go` **0**、越界 **0**（唯一路径 `docs/evidence/s1/360-a1-census.md`）。
- 正控尺（同尺打含 `.go` 的真提交 `ecebc7a9`）＝文件 4、`.go` **2**（尺 NOT 空转）；负控（纯文档笔 `9fc2be02`）＝`.go` **0**。
- 腿间漂移具名：起手 `a9b4513d` -> 交件父提交 `40339d01`，中间 6 件全在 `.scratch`/`docs`（`*.go` **0** 枚），故 (1)-(6) 的读数对交件 HEAD 仍成立。
- 只 commit、未 push；未翻任何框、未碰票 226/227/360 票面与台账一字。

## (9) 台面件清单（只建不删，全在仓外）

`D:\tmp\wisp\20261022\360a1\`：`overlay\zz360a1_probe_test.go` + 三份 overlay.json（主检出形／首棵导出树形／完整导出树形）／`tree-227\**`（首棵，见下）／`tree-check2\**`（完整导出树，17564 枚 = `git ls-tree -r` 同数）／`logs\01..10-*.txt`／`run-*-arm*.ps1` + `run-host-lane*.ps1`。

- 具名更正（本腿自己踩的台损）：**首棵导出树 `tree-227` 少 32 枚文件**——`git archive | tar -x` 经 PowerShell 管道时，32 枚中文名台件被转码（另落 17 枚乱码名），故 `tree-227` 只有 17549 枚。改用 Git Bash 直出 tar 后 `tree-check2` = **17564/17564**。少掉的 32 枚**全是 `.scratch` 台件、`.go` 0 枚**，但 (6a)/(6b)/(4)(5) 的承重读数已在 `tree-check2` 上**逐发重跑复现**（`logs\07..10-*`），本件引用的都是完整树读数。
