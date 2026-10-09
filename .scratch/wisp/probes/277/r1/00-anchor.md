# 00 起手锚 — 腿 277-r1（票 277 AC#1：真机跑 A 道，看常驻进程到底举没举出一张卡）

- 本腿性质：**真机执行腿**。写面＝本目录 `.scratch/wisp/probes/277/r1/*.md` ＋本腿的 commit message。
  ⛔ 产码零改动、⛔ 台件（`cmd/wisp/**`）零改动、⛔ 零 push、⛔ 不翻任何票面框、⛔ 不改台账。
- 落点说明：编排者转述写的是「`probes/277/r1/*.md`」。仓里没有顶层 `probes/` 目录，
  既有各腿一律落在 `.scratch/wisp/probes/<腿名>/…`（例：`probes/246/r2`、`probes/card-proof-prep-1`），
  本件按仓内既有约定落 `.scratch/wisp/probes/277/r1/`，具名登记这一处转述与盘上形状的差别。

## 起手读数（本腿现跑）

```
$ date
Fri Oct  9 09:49:11 CST 2026

$ git rev-parse HEAD
bf9974c229974b4dbd288d0c2aa49a41673d8243

$ git log --oneline -3
bf9974c2 A751/A752 落账：收 281-r1（翻框审计 10 对 10 逐对同文＋Status 两枚原句逐字留档，07/115 一字未动）＋收 283-r1（新尺 227 行 +0/-0 提交形状即硬）⇒ 票 283 三格翻勾并 -done
c87ad40e 283-r1 收件：查重三处＋两处补尺落点＋三发正控四件套读数（种 A/B/C 皆红并逐字还原）
58a4b2fe 281-r1 收口件：八张归位逐笔凭据＋AC#4 复跑对拉（未勾 14→4、-done 98 不变、零撤名）

$ git status --porcelain -- cmd/wisp internal/agent/approval
（空——本腿地界内零在飞改动；工作树别处的 design/**／.gitignore／probes/** 改动不是我的地界，未触碰）
```

## 桌面独占（配方 §1-P6，姿势抄 `docs/evidence/s1/245-esc-not-a-standby-global-hotkey-r1.md`）

```
$ tasklist //FI "IMAGENAME eq balldebug.exe"
INFO: No tasks are running which match the specified criteria.

$ tasklist //FI "IMAGENAME eq wisp.exe"
INFO: No tasks are running which match the specified criteria.
```

两枚皆 0 枚 ⇒ 用例头部那句「Run it with no other Wisp or balldebug process on the desktop」
（`cmd/wisp/resident_task_source_live_246_windows_test.go:38`，本腿直读）满足，A 道可以起。

## 命令原文核对（AC#1 要跑的那一发）

直读 `cmd/wisp/resident_task_source_live_246_windows_test.go:38-41`（`git show` 之外另用 Read 取整行）：

```
// Run it with no other Wisp or balldebug process on the desktop:
//
//	PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -tags winlive ./cmd/wisp \
//	  -run 'TestLive246ResidentPipeline|TestLive246ExitCancels' -v
```

与配方 §2 甲、与编排者转述**三方同文**，照抄执行，未自创。

## 跑前前置自检（配方 §5「A 道跑得动的前提」四条，本腿现尺）

```
$ which go gcc
/d/work/base/go/bin/go
/e/work/base/msys64/mingw64/bin/gcc          # 前提①：go 工具链＋C 编译器在 PATH（台件自建 exe 用）

$ ls third_party/sherpa-onnx/*.dll
onnxruntime.dll  sherpa-onnx-c-api.dll  sherpa-onnx-cxx-api.dll   # 前提②：三枚 dll 在盘

# 前提③真桌面＋零别的 wisp/balldebug 进程：见上一节两发 tasklist
# 前提④Esc 键空闲：只能由跑出来的读数自证（被占则用例自报「判不了」红句，见配方 §4 第三栏）
```

## 一条硬约束的落点确认

`build/wisp.exe` 本腿**不重建、不覆盖**（机主未授权）。A 道吃的是台件在临时目录现建的 exe
（`cmd/wisp/secret_argv_windows_test.go:161/:174`，配方 §2 甲①），与这枚旧 exe 无关。
