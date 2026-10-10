# 票 293 · 写腿 `293-r1` · `30` `AC#3` 越界自查

钟点＝本腿 `2026-10-10 10:1x +08` 现跑；尺全部逐字给出，读数原样抄。

## 1. 名册（逐笔 `git diff-tree -r --name-only --no-renames`）

| commit | 内容 | 名册 |
|---|---|---|
| `2eed98eb` | 起手锚（`.scratch/wisp/probes/293/r1/00-anchor.md`） | 该件路径（`.md`） |
| `26289b9e` | 产码＋测试 | `cmd/wisp/resident_audio_windows.go`<br>`cmd/wisp/resident_ball_windows.go`<br>`cmd/wisp/resident_tray_mute_293_windows_test.go`<br>`cmd/wisp/resident_windows.go` |
| 本笔 | 件（`.scratch/wisp/probes/293/r1/*.md` ＋ `logs/*`） | 见 `git show --name-only` |

十枚禁列逐枚比对（尺＝上面那串名册，逐路径看）：`frontend/**`／`design/**`／`PLAN.md`／`docs/specs/**`／
`internal/observe/thresholds.go`／golden（`testdata/golden`）／`tools/d22scan/allowlist.txt`／
`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／
`internal/perm/ticket90_persist_test.go` ⇒ **0 命中**。另：⛔ 票 299 的射程（`cmd/wisp/panel_host_windows.go`、
`internal/panel/**`）**一枚未碰**；⛔ 三枚冻结件未碰；⛔ 默认值三枚（`schema.go`）未碰；⛔ D43 表／球侧 `Muted` 态未碰。

合法名册口径也守住＝`cmd/wisp/**` ＋ `.scratch/wisp/probes/293/r1/**`，两笔 commit 都带**显式 pathspec**、零 push。

## 2. `internal/ball` 零改动（甲形的硬尺）

```
git diff 6ef14788..HEAD -- internal/ball | wc -c   ⇒  0        （ball_diff_bytes=0）
git diff 6ef14788..HEAD --stat -- internal/ball    ⇒  空
```

## 3. 新包级依赖边 0／新协程 0

```
GOFLAGS= go list -deps ./internal/ball   ⇒  rc=0，145 枚依赖，internal/audio 命中 0
                                            （与 293-a1／编排者 09:4x 那两发同号同值 ⇒ ball 侧名册一字未动）
git diff 6ef14788..HEAD -- cmd/wisp | grep -c 'go func('   ⇒  0
sh scripts/d22scan.sh   ⇒  rc=0（裸 `go func(`＝ban #1 亦由该闸扫，见 20 件的逐字尾部读数）
```
本发没有给 `cmd/wisp` 加任何一条新 import（三处改动用的都是文件里已有的包），也没有扩 `internal/ball` 的导出面
（只是**调用**它本来就导出的 `SetTrayChecks`）。

## 4. 三条禁区的越界面读数

- **禁区①（⛔ 死字段当真相源）**：`grep -rn 'mutedAtBoot' cmd/wisp/` ⇒ 仍只有
  `resident_audio_windows.go:101` 声明 ＋ `:274`（＝`6ef14788` 的 `:256`，本腿 `trayMuteState` 插在其前 18 行 ⇒ 号挪、文本不变）
  一处**写**、**零读** ⇒ 写读比仍是 1:0，本腿没接它（唯一新增的读侧是 `gate.Muted()`）。
  第三命中是新测试文件里**提到**这个名字的一句注释，非代码引用。
- **禁区②（第二实参）**：生产里对托盘显示面只有**一处**写，逐字
  `cmd/wisp/resident_windows.go:366:  rb.attachTrayMuteProjection(raudio.trayMuteState, rb.b.SetTrayChecks)`，
  真正推送那一行逐字 `push(muted, false)`（`resident_ball_windows.go` 的 `mirrorTrayMute` 体内）⇒ 暂停唤醒半**恒显式 false**、
  理由在注释里；本腿**没有**为它新造任何状态位。
- **禁区③（nil 护栏）**：`if rb.b != nil { … }` 包住取值表达式与投影那一跳（`resident_windows.go:365-368`），
  护栏早于 `rb.b.SetTrayChecks`；票 290 那条位于 `attachMuteGate` 之前的早退分支（`startResidentBall` 的
  `if err != nil { … return rb }`）**位置一字未动**。

## 5. ★与票面文字的一枚冲突（具名报回，⛔ 我未据此扩权）

票面「现量」第二条那把尺后面挂着的话：
> `git grep -n "SetTrayChecks" HEAD -- internal cmd` … ⇒ **调用者 0 枚，改前 0、改后仍 0**

「改后仍 0」与 `AC#2` 的要求**自相矛盾**：勾要跟着门动，`internal/ball` 又零改动，那么改后**必然**多出至少一枚
生产写侧（本腿那一枚在 `cmd/wisp/resident_windows.go:366`）。我按 `AC#2` 那一格裁形，并把这条按票面缺陷申报，
不自行改写票面（票面只追加 Progress log 一条，⛔ 未翻任何框）。
改后同尺的实际读数（`grep -rn 'SetTrayChecks' cmd/wisp/ internal/ball/`）＝
定义 `internal/ball/ball_windows.go:953` ＋其注释 `:952` ＋**生产写侧 1 枚**（`resident_windows.go:366`）＋
本腿新测试里的 1 枚编译期断言 `var _ func(*ball.Ball, bool, bool) = (*ball.Ball).SetTrayChecks` ＋若干提及它的注释行。
⛔ 这枚计数**没有被当作 `AC#2` 的凭据**（凭据是 `10` 件的两枚突变体正控）。
