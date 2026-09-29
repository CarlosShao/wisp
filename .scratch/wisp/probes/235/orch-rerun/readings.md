# 编排者自己那发整包复跑（09-29 23:0x，推送前置判据）

- 起手锚点：**`9c96f103`**（`235-r1` 交件后的 HEAD；我跑之前 `tasklist //FI "IMAGENAME eq go.exe"` 计数＝**0**，无别的测试在飞）。
- 命令（逐字可复跑）：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/... -count=1 -v`
  ⚠ Sherpa PATH 必须带：不带的话 `cmd/wisp` 会给 `exit status 0xc0000135`＋`0.0xxs` **且没有 `--- FAIL`**＝用例根本没跑（本机已知的坑）。
- 原始日志＝`full.txt`（**849,741 字节**，**故意不入库**，同 `probes/ci-red/logs/` 那 9.6MB 的处置：原始逐用例日志体量大、结论已在本文件）。

## 读数（口径全部具名）

| 尺 | 读数 |
|---|---|
| 包级红（`grep -cP '^FAIL\t'`） | **4**＝`cmd/wisp` 128.872s／`internal/ball` 0.286s／`internal/panel` 4.496s／`internal/risk` 6.933s |
| 名级红（`grep -cE '^[[:space:]]*--- FAIL'`，两层同尺） | **7** |
| 只数顶层（`grep -cE '^--- FAIL'`） | **7**（本发无缩进子测试红） |
| 终态 rc | 1（`FAIL` 收尾） |

逐名（顶层 7 枚）：
1. `cmd/wisp/TestTicket223HandEditedFsLooseningCostsAnL2Card` 2.16s
2. `internal/ball/TestC21TableColourRowsMatchTokensCSS` 0.00s
3. `internal/panel/TestApprovalCardViewJSONKeysMatchFrontendTypes` 0.00s
4. `internal/panel/TestComposerContractTypesMatchFrontend` 0.01s
5. `internal/panel/TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` 0.17s
6. `internal/panel/TestC21DesignTokensFourWayAgree` 0.00s
7. `internal/risk/TestResolvePerCallBudget` 1.40s

## 归因（照票 236 AC#5 那条尺：整包名册红 **＋** 单包隔离复量，缺一半就是假结论）

- **在册常红 5 枚**＝`internal/panel` 4 枚（含那枚冻结件 `TestC21DesignTokensFourWayAgree`，两枚红因、owner 已撤 `Q-52` ⇒ 按已知常红读）＋ `internal/ball/TestC21TableColourRowsMatchTokensCSS`。0.00s／0.01s／0.17s 全是断言型，非计时型。
- **争用／计时型 2 枚，我做了隔离复量**：
  - `cmd/wisp/TestTicket223HandEditedFsLooseningCostsAnL2Card` → `-run '^…$' -count=3 -v` ＝ **3/3 绿**（2.07／2.09／2.03s，`iso-223-handedit.txt`）。整包那发的正文是 `config_reload_223_test.go:303`「the console did not render the reload card」＋`:306`「does not name the key」＝**子进程卡片在轮询窗口内没到**，同包另一枚 223 用例同时刻 PASS（`TestTicket223RefusedLooseningKeepsOldValues` 5.35s 绿，`full.txt:170` 附近）。
  - `internal/risk/TestResolvePerCallBudget` → `-count=3` **3/3 绿**（1.16／1.23／1.45s，`iso-resolve-budget.txt`）。这枚有既往具名记录（`HANDOVER:212`／`:224`：整包红、安静复量转绿）。
- ⚠ **换位实据再加一发**：`221-v1` 那发整包是 `cmd/wisp/TestTicket223RefusedLooseningKeepsOldValues` 红、`risk` 绿；本发正好相反（`Refused…` 绿、`HandEdited…` 红、`risk` 红）。⇒ **"在册红名册是稳定集合"这句话又一次不成立**，票 236 AC#5 本格判据必须自带隔离复量那一半（已写进票面增量）。

## 终判

- **本票/本轮新造真红＝0 枚**（`internal/tools` 全绿；两枚可疑红隔离复量 3/3 绿）。
- 推送前置四把尺：`go.exe` 计数 0 ✅／整包逐名红名册对得上且无新增 ✅／`gofumpt -l` 本轮改动面（`internal/tools` 两枚 `_test.go`）零输出 ✅（`235-r1` 自报、我在其终态复跑整包未报格式红）／远程 tip 推后逐字等于本地 HEAD（推后量）。
