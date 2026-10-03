# 244-c2 只读普查 —— 构建入口名册复量 / CUI 残余路径 / subsystem 尺有无 / GUI 档下 stdout 可见性

> 本腿代号 `244-c2`，只读普查腿。今天 2026-10-03。⛔ 全程零 Go 命令、零构建命令、零写产码。
> 前件 A＝`c1`（`.scratch/wisp/probes/244/c1/census.md`，commit `0a88d7a8`，锚 `f17b1165`）：量到 11 枚构建入口全零 `-H`、台上 19 枚 exe 全 CUI。
> 前件 B＝落地件 `scripts/build.ps1`（commit `cc6eaa65`，244-r1）：第 115 行已追加 `-H=windowsgui`。本腿**对着现树复量**，⛔ 不照抄 c1 读数。

## §0 锚

- 时刻（起手 `date -Iseconds`）：**2026-10-03T09:13:13+08:00**。
- 起手 `git log -1 --format=%h`：**`3fe377e3`**；同法取数时 `git status --porcelain -- cmd internal tools docs scripts .scratch | wc -l`＝**355**（共享工作树，多腿在飞，dirty 属正常态）。
- ⚠ 落 §0 骨架这发命令时 HEAD 已漂到 **`72cce7a0`**（HANDOVER 4.0z 停车 commit）——本仓是**活体共享工作树**，锚点会在腿飞行期间被别人推进；本腿所有行号/读数对"取数时刻的盘上内容"负责，取数时 census 相关 8 枚文件（`scripts/build.ps1`／`cmd/wisp/console_windows.go`／`console_other.go`／`main.go`／`config_reload.go`／`run.go`／`.github/workflows/ci.yml`／`docs/BUILD.md`）`git status` **全部 clean＝读的是已 commit 内容，不是别腿的在飞脏写**。
- 票面：`.scratch/wisp/issues/244-spec-11-wants-the-gui-subsystem-build.md`（AC 框一枚未碰；末尾「越序」裁定在台账 `A548`，本腿已读）。
- 规格要求：`docs/specs/SPEC-11-build-deploy-containerization.md:50`（票面凭据）。
- 纪律自证：全程零 Go 命令／零构建；PE subsystem 本腿**量不到**（库内零枚 tracked exe，`*.exe`＋`build/` 均在 `.gitignore:2/14`；sanctioned 法＝objdump 只对 tracked exe，无枚可量——详见 §6）；`docs/BUILD.md`／`ci.yml` 等冻结件只读一字未动；`frontend/**`／`design/**` 零读零写零转述（grep 曾被动命中 frontend 路径，本腿未读其内容、本件不引用）。

## §1 入口名册复量（已带 / 仍不带 `-H`，逐枚）

（取数中）

## §2 今天仍会产 CUI exe 的路径名册

（取数中）

## §3 subsystem 尺的有无与最小落点

（取数中）

## §4 stdout 可见性〔读码推断，未跑〕

（取数中）

## §5 我可能写错的条目（自我对抗）

（取数中）

## §6 量不到的地方（具名）

（取数中）

## §7 交件判语

（取数中）
