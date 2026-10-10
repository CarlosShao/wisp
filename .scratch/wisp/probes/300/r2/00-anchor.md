# 300-r2 起手锚（票 300 `AC#6`，落点＝编排者裁的 C2 形）

跑时刻（`date '+%F %T %z'` 的 stdout）＝`2026-10-10 18:38:03 +0800`

## 尺与读数（逐字）

尺①＝`date '+%F %T %z'`
```
2026-10-10 18:38:03 +0800
```

尺②＝`git log -1 --format='%h %ad %s' --date=iso`
```
05db4bc6 2026-10-10 18:36:33 +0800 303-r1 收＋裁完（A823）：★回执那一跳接上（我母仓复跑拿到页面原话 REPLIED×3），而票面那句"修后三枚全绿"是我造的——nail2 在回归之前的同一枚台面上就是逐字同一句红 ⇒ 两枚残余症状单独立成票 305；⇒ 只翻 AC#2（AC#3 事实层全核到⛔ 由我按改过的目标自盖章，交 303-v1 三道必答）
```

尺③＝`git status --porcelain -- .scratch/wisp/issues internal cmd docs | head -40`
```
（0 行）
```

尺④＝`git rev-parse --abbrev-ref HEAD`
```
dev
```

尺⑤＝进程前置尺 `tasklist //FI "IMAGENAME eq <名>" | grep -c "<名>"`（三枚各一发）
```
wisp.exe      = 0
balldebug.exe = 0
mockllm.exe   = 0   （tasklist 回显 "No matches found"，grep -c 计数 0）
```

## 起手锚 SHA

**`05db4bc6`**（= 尺② 那一行的 `%h`）。本腿全部读数落在这枚锚之后；交回时再报一枚终锚。

## 本腿射程（派单原文的允许写面，逐字照抄）

- 新增一枚 `internal/audio/` 下带 `//go:build windows` 的 `_test.go`。
- 自己的证据件目录 `.scratch/wisp/probes/300/r2/**`（只新建 `.md`／`.txt`）。
- 票面 `.scratch/wisp/issues/300-*.md` 追加一节（append-only）。

⛔ 动产码、⛔ 碰 `internal/audio/parse_wave_format_300_windows_test.go`、⛔ 碰 `wavinjector_test.go`、
⛔ 动 `scripts/**`／`.github/**`／`frontend/**`／`design/**`／`tools/**`／冻结面、⛔ push。

rc-anchor=0
