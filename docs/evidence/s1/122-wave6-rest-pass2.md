# 证据件：票 122 — pass2 清账波 6（tools 之外的全部余包）

- 波次：wave 6（122-wave6-rest-pass2）
- 派单：`.scratch/wisp/probes/orch/20261113/txt/122-wave6-rest-pass2-dispatch.txt`
- 开工 HEAD：`283a24dfa`；staticcheck 2026.2.2 (0.8.2)；Go 1.27.1 windows/amd64
- 口径（自跑，与派单逐字同）：
  `GOOS=windows CGO_ENABLED=0 staticcheck -tags=winlive -checks=all,-SA9009 ./internal/winsec/ ./internal/audio/ ./internal/proc/ ./internal/ball/ ./internal/risk/ ./internal/memory/ ./internal/observe/ ./internal/secret/ ./internal/agent/`
- 起手复跑＝**23 条**，与派单清单逐条一致（派单列 24 条；1 条差异见 §4——`aecspike/rig_windows.go:49 iidIUnknown` 不在此目录口径射程内，见 §4 判定）。
- 终态：**同口径复跑 0 条（EXIT=0）**；无新增。（commit `45c3b6511` 后终验复跑一次，同读数。）

## 1. 逐条表（23 条，派单口径）

| # | 检查 | 位置 | 符号 | 档 | 处置与理由 |
|---|------|------|------|----|-----------|
| 1 | ST1000 | winsec/pinnedleaf_342_windows.go:3 | 包注释 | 真改 | 文件级注释改 `Package winsec - file-level note:` 形；与同族 unix 版 `pinnedleaf_342_unix.go:3` 逐字同型。注释内容零改动，仅加头 |
| 2 | U1000 | winsec/acl_windows_test.go:32 | `wellKnown` | ③删 | 全仓 grep 唯一出现＝声明行；两成分 `systemSID`/`adminsSID` 各有独立读者 |
| 3 | U1000 | winsec/leftstanding_354_windows_test.go:68 | `leftStandingInjected354r1` | ③删 | 全仓唯一出现＝声明行；同包他文件的注入串是**各文件自有的独立常量**（340/353/354 各持己份），共享判据不成立。文件内另一引用（:197 控制支）是独立 `errors.New` 字面量，不涉此常量 |
| 4 | U1000 | winsec/seam_guard_windows_test.go:47/49/53 | `narrowOnlyResolver` 族（type+2 法） | ③删 | 全仓唯一出现＝本文件自定义处；测试全程用孪生 `narrowOnlyResolverB`（:275）。文件头注释自述其意：与孪生 B 配对、给单次闸「两个**不同**的合法 resolver」——删 A 留 B，剩余组合（rubber-stamp 拒＋B 拒＋keep-first 断言）完整保住该意图 |
| 5 | U1000 | winsec/seam_guard_windows_test.go:70 | `countingResolver.installs` 字段 | ①钉 | `//lint:ignore U1000 <点名>`：类型语义＝身份见证（文件注释逐字「which an empty struct cannot」）——AC#1 钉的是**同一实例被保留**（指针同一），字段值从不被读；删字段会破坏"counting 形与空 struct 的差别"，空 struct 撑不住"两实例等价即被误认"的判别。API 按偏移读，字面不读值 |
| 6 | ST1003 | audio/mmdevice_windows.go:284/291/293/295 | `deviceId`→`deviceID` | 真改 | 4 处参数改名，零本地引用需动（全部仅接收不使用）；COM 回调槽位布局不动 |
| 7 | ST1003 | audio/wasapi_windows.go:36 | `pGetCurrentThreadId`→`pGetCurrentThreadID` | 真改 | 声明＋:42 调用点，两处齐改；DLL 名（`"GetCurrentThreadId"`）不动 |
| 8 | U1000 | audio/aecspike/rig_windows.go:49 | `iidIUnknown` | ③删 | 见 §4 专项论证（派单★提醒 #1 判据链逐条核） |
| 9 | ST1021 | proc/boot_windows.go:22 | `Runtime` 注释 | 真改 | `// Boot assembles...`→`// Runtime assembles...`，仅换头词，正文零动 |
| 10 | ST1003 | proc/jobscope_windows.go:31 | `jobObjectBasicProcessIdList`→`...ProcessIDList` | 真改 | 声明＋注释＋:200 引用点，三处齐改 |
| 11 | ST1003 | proc/treemetrics_windows.go:47 | `procGetExtendedTcpTable`→`...TCPTable` | 真改 | 声明＋:339 调用点，两处齐改；DLL 名 `"GetExtendedTcpTable"` 不动 |
| 12 | ST1012 | ball/hotkey_borrow_refused_260r5_test.go:101 | `t260r5ErrNoWindowDoor`→`errT260r5NoWindowDoor` | 真改 | 声明＋:216 return＋:268 `errors.Is`，3 处齐改 |
| 13 | ST1012 | ball/hotkey_refused_borrow_release_245_test.go:58 | `t245ErrNoWindowDoor`→`errT245NoWindowDoor` | 真改 | 声明＋:172 return＋:197 `errors.Is`，3 处齐改 |
| 14 | U1000 | risk/anchor_world_334_test.go:55/60 | `hostilePOSIX347`/`danglingTarget347 | ④钉 | 双口径（winlive＋默认）恒死：声明在无 tag 文件，全部读者在 `//go:build !windows` 齿文件。按**本文件对镜像问题 `hostile334` 的既有型制**钉：`var _ =`＋注释。不删：风险测试的"世界"文件是平台间共享词汇表，票 347 恰是为消灭平台副本才把它们并入此文件——挪进 `!windows` 文件＝反转票 347；删了则 linux 口径会报"声明被删而读者还在"。不挪：同上。镜像问题在同文件 ：49 已有先例钉法，本波照办并扩注。**linux 口径佐证：GOOS=linux 下这两条 finding 自然消失（LINUX_EXIT=0），win 口径 `var _` 钉死** |
| 14b | — | risk/anchor_world_334_test.go | — | — | （与 #14 同一条 finding 的两枚常量，表内合并记 1 条） |
| 15 | U1000 | memory/l1_warm_29_test.go:151 | `l1WaitFor` | ③删 | 全仓唯一出现＝声明；同文件测试已各自内联等价等待或同步观察 |
| 16 | U1000 | observe/logprivate_132_windows_test.go:99 | `dirIsSealed132` | ③删 | 全仓引用全在 `logprivate_132_other_test.go`（`//go:build !windows`）；windows 口径恒死（④形态但删法同③：**按平台文件分家**——helper 本就该住平台文件里，非"挪 tag"）。同文件 `judge132` 等 helpers 有本地读者，留 |
| 16b | — | observe/logprivate_132_windows_test.go | — | — | （16/16b 为同一 finding 的记账拆行，实改 1 处） |
| 17 | SA4010 | secret/delete_test.go:131 | `got = append(got, b.ID)` | 夹具疏忽→删 | 见 §3 专项论证（派单★提醒 #2） |
| 18 | U1000 | agent/reminder_timer_test.go:109 | `(*fakeReminderClock).releasedCount` | ③删 | 全仓唯一出现＝声明。**删前核了删后不制造新 U1000**：字段 `released`（:42 写 ：81 读）与 `due` 均有其他读者，`timers` map 同 |
| 19 | — | （14b/16b 合并行，非独立 finding） | — | — | — |

## 2. 门禁（逐项）

| 门禁 | 结果 |
|------|------|
| `go build`（九包） | rc=0 |
| `go vet`（九包） | rc=0 |
| `go test -count=1` 八包 | 全 ok：winsec 73.391s / proc 12.253s / risk 29.328s / memory 65.894s / observe 15.952s / secret 8.225s / agent 28.512s / ball 3.063s |
| gofumpt -l（九包） | 0 行 |
| d22scan.sh（Git bash） | **clean rc=0**（runtests.sh top-level PASS=45 FAIL=0 SKIP=0；scan clean，examined 480 production Go files） |
| staticcheck 口径 | 前=23 条 EXIT=1 → 后=**0 条 EXIT=0** |
| SKIP 项 | `internal/audio/aecspike` 测试 **SKIP 如实报**（sherpa-onnx cgo DLL 依赖本机起不来，staticcheck 探针输出 `build constraints exclude all Go files in .../sherpa-onnx-go-windows@v1.13.8 (compile)`）；`go test ./cmd/wisp/` 按票 98 豁免未跑 |
| 佐证 | GOOS=linux 静态检查 risk+memory：**EXIT=0**（④钉的两枚常量在 linux 口径自然活着，钉法双向成立） |

## 3. secret/delete_test.go SA4010 读码判定（★提醒 #2，非机械改）

`TestBlobsListsMetadataOnly` 中：

```go
var got []string
for i, b := range blobs {
    got = append(got, b.ID)   // ← SA4010：这个 append 的结果无人消费
    if b.ID != want[i] { ... }
    ...
}
```

判定＝**测试夹具疏忽，不是被测行为，更不是真 bug**：
1. 循环内的所有意图断言（顺序、ID、Ref、两个时间戳非零、Size 非零）**逐元素直接落在 `blobs` 本体上**；
2. `got` 收集完从未被读——`want` 与 `got` 从未 compare；若夹具想比 `reflect.DeepEqual(got, want)`，逐元素断言已逐字覆盖且更严（错误信息带下标与具体值）；
3. 丢弃 `append` 结果不改变任何被测行为，是遗留的半截夹具。
处置＝**删 debris**（`var got`＋append 行，-2 行），覆盖无损失；补建 DeepEqual 才是造新意图。

## 4. aecspike `iidIUnknown` 判定（★提醒 #1 判据链）

1. **全仓 grep**：唯一出现＝声明行 `rig_windows.go:49`；`unsafe.Pointer`/`IID` 路径**无人按偏移读它**——COM 调用传的是 `&clsidMMDeviceEnumerator`（:148）与 `&iidIMMDeviceEnumerator`（:149）的地址，**别的** GUID；
2. **双口径回归实验**（改后临时加回该行实测，弃前先记数）：
   - 带 `-tags=winlive`：`rig_windows.go:49:2: var iidIUnknown is unused (U1000)`；
   - 默认档（无 tag）：同报 U1000。
   ⇒ **不是④类**（④＝只在非默认 tag 有主，带 tag 后自然消失），两种口径都死 ⇒ 判据链指向③；
3. **先例判例**：wave 4 对 `wasapi_windows.go:49` 的同名孪生判③删，判据逐字同：「它是纯 Go 值变量（GUID 常量），从不作为内存块被 API 按偏移读——COM 调用传的是**别的** GUID 的地址」；
4. 派单口径为什么没列进 23 条：起手口径用 `./internal/audio/` **目录模式不递归子包**，aecspike 是 audio 的子包；本波处置后以显式子包口径（`./internal/audio/aecspike/`）验证：**finding 行＝0**（commit 后终验）；staticcheck 进程退出码＝1，**源自 sherpa-onnx cgo 依赖本机不可编译**（`build constraints exclude all Go files in .../sherpa-onnx-go-windows@v1.13.8 (compile)`），是环境噪声不是 finding——与门禁 SKIP 项同因。（更正：本件初稿在此误写"EXIT=0"，commit `45c3b6511` 后终验抓出，按追加更正。）

## 5. 写面审计

- 触点 17 个 .go（winsec 4 ／ audio 3 ／ proc 3 ／ ball 2 ／ risk 1 ／ memory 1 ／ observe 1 ／ secret 1 ／ agent 1），**全在派单 9 包内**；
- 16/16 触点 build tag 与改前一致（audit 脚本读取触点文件首 3 行逐一比对）；
- 派单禁面零触碰：`internal/tools/**`、`cmd/**`、`internal/panel/**`、`internal/plugin/**`、`internal/config/**`、`frontend/**`、`.github/**`、`scripts/**`、`docs/PLAN.md`、`docs/specs/**`、thresholds/golden/allowlist——零 diff；
- 零 `lint:ignore` 新增于 prod 码（仅 winsec 测试夹具 1 处①钉，含点名注释）；
- 临时件落 `D:\tmp\wisp\20261113\122-wave6\`（run-sc.ps1/run-gates.ps1/run-corroborate.ps1＋regression 脚本 `D:\tmp\winreg.ps1`），只建不删。

## 6. 波 6 四档统计（按 finding 粒度，合计必须＝24）

| 档 | 条数 | 明细 |
|----|------|------|
| 真改 | 11 | ST1003×7（deviceId×4、pGetCurrentThreadID、jobObjectBasicProcessIDList、procGetExtendedTCPTable）＋ ST1012×2 ＋ ST1021×1 ＋ ST1000×1 |
| ③删 | 10 | winsec 5（wellKnown、leftStandingInjected354r1、narrowOnlyResolver type+两法×3）＋ audio aecspike 1（iidIUnknown）＋ memory 1 ＋ observe 1 ＋ agent 1 ＋ secret SA4010 debris 1 |
| ①钉 | 1 | winsec countingResolver.installs（lint:ignore＋点名） |
| ④钉 | 2 | risk 两枚常量（一处 var _ 块＋扩注；linux 口径佐证 EXIT=0） |
| 合计 | **24** | = 派单口径 23 ＋ aecspike 1（§4） |

（更正记录：本表初稿按"处置动作"粒度误记 ③6/真改10，commit `45c3b6511` 后自检按 finding 粒度改正；§1 表的 14/14b、16/16b 合并行系记账拆行，此处已拆回。）

## 7. Commit

- 显式 pathspec：17 个触点 .go ＋ 本证据件 ＋ 票 122 Progress log，共 19 path；
- ⛔ push；临时件与共享树他腿脏件零触碰。
