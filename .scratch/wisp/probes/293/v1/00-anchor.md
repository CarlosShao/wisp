# 票 293 · 非实现者对抗验收腿 `293-v1` · 00 起手件

时刻＝`2026-10-10 11:05:06 +0800`（`date` 的 stdout，未手打）。
本件在任何长跑命令之前 commit。本腿只裁 `AC#2`／`AC#3`／`AC#4`；`AC#0`／`AC#1` 不归我；
`AC#5`＝机主肉眼那枚，**我判不了、也不算腿的欠账**。
**我⛔ 不翻任何 `- [ ]` 框**（翻框＝编排者权限）。判语只进本目录的件。

---

## 1. 锚（现量，不采信派单）

- 派单给的产码足迹末笔＝`0d884bb5`（写腿 `293-r1` 第 6 笔）。
- **起手 HEAD 现量＝`b3593507`**，分支 `dev`。⇒ **锚已漂 6 枚**（`f4f1463a`→`2af4c221`→`a96e64ee`→
  `5ad1ab94`→`ed07e53c`→`b3593507`＝票 297 普查腿 `297-a1` 五笔＋编排者收 293 一笔）。
  这是本仓惯例、不是异常，但**本腿全部读数都取自 `b3593507`，不是 `0d884bb5`**。
- 尺＝`git rev-parse --short HEAD` ⇒ `b3593507`；`git rev-parse --abbrev-ref HEAD` ⇒ `dev`。

## 2. 我自己复跑的足迹（⛔ 不引用实现者的结论当凭据）

| 尺（逐字） | 我的读数 |
|---|---|
| `git diff --no-renames --stat 6ef14788..HEAD -- cmd internal` | 5 枚文件全在 `cmd/wisp/`：`config_readers_255.go` 4／`resident_audio_windows.go` 20／`resident_ball_windows.go` 83／`resident_tray_mute_293_windows_test.go` 299／`resident_windows.go` 38 ＝ **439 insertions / 5 deletions** |
| `git diff 6ef14788..HEAD -- internal/ball \| wc -c` | **0 字节** |
| `grep -c '^func TestAC293' cmd/wisp/resident_tray_mute_293_windows_test.go` | **7**（名册逐名见下） |
| `git diff-tree -r --name-only --no-renames <六笔>` | 六笔逐枚名册见 §3；**十枚禁列路径 0 命中** |
| `grep -c '^- \[x\]' .scratch/wisp/issues/293-…md` | **2**（＝`AC#0`／`AC#1`） |
| `git diff --stat 6ef14788..HEAD -- .scratch/wisp/issues/293-…md` | **1 insertion / 0 deletion** ⇒ 写腿一枚框没翻、票面除追加那行一字未动 |

7 枚用例逐名（`grep -o '^func TestAC293[A-Za-z]*'`）：
`TestAC293TrayItemCheckmarkFollowsTheGate`／`TestAC293MuteHotkeyCheckmarkFollowsTheGate`／
`TestAC293BootProjectionShowsTheFactoryMutePosture`／
`TestAC293RefusedDeviceStillProjectsTheGateNotTheHopedOutcome`／`TestAC293NoGateWritesNoCheckmark`／
`TestAC293ProjectionIsNilSafeWithoutABallWindow`／`TestAC293UndrivenTrayFlagIsNotTheSourceOfTruth`。

## 3. 六笔 commit 逐枚存在性（`git log -1` 现跑，⛔ 不认派单里的号）

`2eed98eb`（起手锚，只 .md）→ `26289b9e`（落地：4 枚 `cmd/wisp` 产码/测试）→ `d1f17fc7`（2 枚 .md）→
`9dade451`（**复位**：`config_readers_255.go` + `resident_audio_windows.go` + `resident_ball_windows.go`）→
`e96da9f4`（AC#4 门禁件 + 票面追加）→ `0d884bb5`（票面 Progress log 归位）。
⇒ 产码只出现在 `26289b9e` 与 `9dade451` 两笔；`internal/` 名下**一枚都没有**。

## 4. 环境（起跑前）

- `tasklist | grep -icE "wisp|balldebug"` ⇒ **0 枚**（尺回 `rc=1`＝零命中）⇒ 允许开整包。
- ⚠ **工作树本来就脏 805 枚**（`git status --porcelain | wc -l` ⇒ 805；含 `design/**`、`frontend/**`、
  `.gitignore`、多枚 `probes/**/logs/*.txt` 等，**全部不是我改的**）。本腿⛔ 碰、⛔ add、⛔ 清理。
  回报口径＝**"我动过的路径恢复原状"**，⛔ 说"工作树干净"。
- go 测试一律带 harness：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test …`。
- ⛔ push／⛔ `add -A`／⛔ `amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／⛔ `--no-verify`／仓内⛔ worktree。
- 突变体还原用**我自己的备份文件写回**＋`sha256sum` 对值，⛔ `git checkout`。

## 5. 本腿要交的四处（欠账在此具名，跑完在 40 件逐格判）

1. `AC#2` ⓐ：7 枚用例逐枚答"断的是**门态**还是断'被调过'"，并按 `A710` 那把尺问"判据换成反形它会不会绿"。
2. `AC#2` ⓑ：**我自己造**两枚反控突变体（恒假／摘掉手势那一跳），具名跑到哪几枚红＋红句逐字。
3. `AC#4`：`GOFLAGS= go build ./...`／`sh scripts/d22scan.sh`／`gofmt -l` ＋ `$(go env GOPATH)/bin/gofumpt.exe -l` 并排／
   整包（`./cmd/wisp/ ./internal/panel/`，⚠ 票面写的是 `./cmd/wisp/ ./internal/audio/ ./internal/ball/`，
   两把我都跑，见 20 件）。
4. 跨票动仪器（`config_readers_255.go` `:316`→`:322`）——**独立裁**，见 30 件。
