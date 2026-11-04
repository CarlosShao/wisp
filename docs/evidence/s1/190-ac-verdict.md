# 190 AC 裁决（宿主侧工作区文件树；编排者 ＋ 独立只读核验腿 190-verify-min，2026-11-04 16:0x +08）

## 落地（8 commit，动件仅 8 件）
- `0d5aa9b4c` types 骨架｜`53b0d592d` reader｜`025cf41f9` 修 json tag 与 `-i`＋屯据绿｜`4d3daacea` 接线1（composer 顶层 `files *TreeView omitempty`＋pump reader/装配＋`ReadFileTree190` 导出）｜`2ed557db1` 接线2（`cmd/wisp` `fileTreeView` 用 `workspaceView().Canonical` 同坐标＋run 装配）｜`36ac694c0` AC#4 三形枚举｜`5229a5e41` reader 补 AC#4 三形＋AC#3 只读＋AC#2 not_chosen 判据。
- 动件：`internal/panel/{tree_types_190,tree_reader_190,tree_reader_190_test,tree_190_ac234_test,composer,pump}.go` ＋ `cmd/wisp/{panel_pump,run}.go`。零触 PLAN/specs/thresholds/golden/D34/allowlist。

## 逐 AC 判（独立只读核验腿 190-verify-min 现跑：build panel+tools+cmd/wisp rc=0；Tree190 6/6 PASS；vet rc=0）
- **AC#1 设计核**：成立（`docs/evidence/s1/190-file-tree-design-a1.md`；票面已 [x]）。
- **AC#2 只读落地＋判据**：成立。`TestReadGitTree190RealRepo`（临时目录合成树断 `a.txt` 在 entries）＋`TestReadGitTree190EmptyRoot`（空根→`not_chosen`，不落整盘）。
- **AC#3 只读边界＋判据**：成立。`TestTree190ReadSurfaceHasNoWriteCalls`（扫 reader 源码无 `os.Remove/WriteFile/...`）；tree 件零 `filepath.Clean/Abs`。
- **AC#4 越权三形**：成立。reader：Lstat 不存在→`unreadable`、权限→`permission_denied`、symlink→`link_outside`、非目录→`unreadable`；`TestTree190RootDoesNotExistUnreadable`＋`TestTree190LinkOutsideRefused`（真 symlink 跑通）钉住两形，permission 走代码路径。
- **AC#5 契约轴**：成立。八 commit 零触 PLAN/specs/thresholds/golden/D34/allowlist；tree 件零 `filepath.Clean/Abs`；`ReadFileTree190` 是包函数、composer `Files` 是结构字段，均非 C17 方法。
- **AC#6 门禁**：成立（host 侧）。build/vet rc=0；Tree190 6/6；d22scan clean(rc=0)。panel 整包红=2 照实记不修：`TestC21DesignTokensFourWayAgree`（既有 tokens 域，clone --shared 父提交 `66fdb2275` 正控证非本票引入）；`TestComposerContractTypesMatchFrontend`（Go 出 `files` 键、`frontend/src/lib/panel.ts` `ComposerState` 未声明——按 **181/197/248 先例** Go 先出键、页面声明归 owner UI 腿，属**契约移动正常刻度**，非缺陷）。

## 边界与遗留（诚实）
1. **文件名 scrub**（census §落点列的 `display_scrub_216.go` 硬化项）：AC 未点名，本波未接；登记为**硬化遗留**（不影响只读读面正确性——路径取自 git 索引，非用户手输）。
2. **页侧**：`frontend/src/lib/panel.ts` `ComposerState` 需声明 `files` 键（owner 外部 UI agent，本队 ⛔ 零写面）。兑现前 `TestComposerContractTypesMatchFrontend` 保持红＝已知契约移动。

## 判语
190 **宿主侧功能成立，可 -done**（按面板簇"宿主侧供数据/前端只渲染"分工与 181/197/248 先例）；页侧 `files` 声明移交 owner UI agent；文件名 scrub 登记为硬化遗留。
