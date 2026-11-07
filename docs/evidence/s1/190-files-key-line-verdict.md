# 190 — 顶层 `files` 键 reader 腿裁决（2026-11-07 14:1x +08）

## 裁决

- 修法选型：**甲是正确修法但非本腿写面；乙否；本腿执行丙**。
- 现跑取证（`git grep` 全仓＋逐行读）——`files` 键**有读者**，不是死键：
  - 生产链全通：`cmd/wisp/run.go:977`（`Files: rt.fileTreeView`）→ `cmd/wisp/panel_pump.go:294-296`（`ReadFileTree190("")`／`ReadFileTreeForWorkspace190`）→ `internal/panel/pump.go:288-290`（reader 存在才填 `composer.Files`）→ `internal/panel/composer.go:285`（`Files *TreeView json:"files,omitempty"` 上线）。
  - 页侧读者未落：frontend 的 FilesPane 今天吃静态 fixture（`right-rail.tsx:211 fileTree = TREE`、`App.tsx:349 fileTree={[]}`），不读 `composer.files`。
- **甲**（给 TS 接口补 `files` 声明）：方向正确——与本仓已判先例同刻度（`frontend/src/lib/panel.ts` 的 `git`/`currentModel`/`credentialState` 就是页侧补声明，票 400 AC#1；新 TS 键声明＝线契约变更，先例 A1378＝机主预批）。**但**台账 `A1487` 已明裁"按 181/197/248 先例＝Go 先出键、页面声明归 owner UI 腿"、"本队 ⛔ 零写面 frontend"，且该步在待机主解闸清单上（"③前端 panel.ts 声明 files 键"）。⇒ 非本腿执行，解闸路径不变。
- **乙**（Go 侧删 `files` 键）：前提"没人读"不成立（上文生产链）；删＝拆毁已裁决、已落地、tree_190 一族 6/6 绿的宿主功能，与 `A1487` 方向相反。⇒ 否。
- **丙**（本腿执行）：线级具名用例把 Go 侧那一半合同钉死（键名、层级、上线/不上线语义、第三态随键上线），页侧 UI 腿落地时对着的是不漂移的线；合同尺红照实记，不放宽断言、不 skip。

## 交付

- `internal/panel/composer_files_190_test.go`（新件，三枚具名用例）：
  1. `Test190ComposerFilesKeyAbsentWithoutReader` — 无 reader 时键整体不上线（omitempty 承诺："a section with no reader off the wire entirely"）；
  2. `Test190ComposerFilesKeyWireShapeWithReader` — 有 reader 时键名恰为 `files`、居 composer 顶层（与 `git`/`currentModel` 同层、不嵌套）、载荷恰为 TreeView 的七个 json 名＋entries 元素形（path/isDir/size）；
  3. `Test190ComposerFilesKeyCarriesRefusalNotDropped` — not_chosen 拒绝形随键上线（kind/reason/空 entries 数组）：拒绝是键上的一句话，不是键的缺席。

## 门禁读数（2026-11-07 14:0x-14:1x 现跑，Git Bash；测试前 unset GOTMPDIR/TMPDIR/TEMP/TMP、export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH）

| 尺 | 读数 |
|---|---|
| 现跑红句 | `composer_test.go:78: Go ComposerState emits [files] that interface ComposerState does not declare` |
| `gofumpt -l internal/panel` | 新件干净；`composer.go` 列脏＝**145 腿未提交注释插入拆了对齐组**（该件 mtime 14:01:36、composer_test.go 14:03:08，wave54"145 筛查"在飞；HEAD 版 composer.go gofumpt 干净 rc=0 已证）⇒ 非本腿引入、非本腿可收 |
| `go vet ./internal/panel` | rc=0 |
| `go build ./internal/panel ./cmd/wisp` | rc=0 |
| `go test -run 'Test190' -count=1 -v` | 3/3 PASS |
| `go test -run 'TestComposerContractTypesMatchFrontend|Test190' -count=1` | Test190 三枚绿；合同尺红＝页侧待 UI 腿（A1487 裁决的中间态） |
| `go test ./internal/panel -count=1 -timeout 20m` | 红=2：`TestComposerContractTypesMatchFrontend`（本键，页侧待声明）＋`TestC21DesignTokensFourWayAgree`（既有 tokens 域） |
| `sh scripts/d22scan.sh` | clean（无 D22 ban 违规） |

## 边界（没做的，如实）

- 未动 `frontend/**`（本队零写面；页侧声明归 owner 外部 UI 腿，台账待机主项③）。
- 未动 `internal/panel/composer.go`／`composer_test.go`——本腿修法不需要改 Go 侧，且**此刻这两件有 145 腿未提交在飞件**（mtime 14:01/14:03），本腿对其零写入、零提交。
- 未翻票面 AC 勾/Status；仅追加 Progress log 一行。
