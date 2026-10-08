# 255-a2 起手锚（只读普查腿）

任务：票 255 `AC#4`「改面板宽度＝改完立刻见效、不用重启」的**落点表**（六问）。
本程零 go 命令、零产码写点；写点只有 `.scratch/wisp/probes/255/a2/` 与工单 255 追加一节。

## 现量锚（本腿自己量，不写「应为 <sha>」）

```
$ git log -1 --format="%H %ad %s" --date=iso
5738661192bad76d95a8c5abc659a1c95f8d453d 2026-10-08 10:35:14 +0800 189-a1(10-08) 起手锚：...
rc=0

$ git branch --show-current
dev
rc=0
```

- HEAD 现量 `5738661192bad76d95a8c5abc659a1c95f8d453d`，**不早于** 派单里给的地板 `914177e6`（差 1 笔以上，编排者话术已确认并行腿在推 HEAD）。
- 起手时刻：2026-10-08（本机 +0800）。

## 工作树现状（只记录，不动）

`git status --short` 前 40 行已见：`M .gitignore`、`D design/**`（一批）、`M design/doubao/**`、
`M .scratch/wisp/probes/152/my152.py`、`M .scratch/wisp/probes/161/r6/logs/*.txt`、
`M .scratch/wisp/probes/268/v1/evidence.md`、若干 `?? .scratch/…` 临时件。
⇒ **以上全部不是我该碰的**，本程只 `git add --` 我自己产的两三枚文件。

## 查重对象已就位（起手核存在性）

```
$ ls -d .scratch/wisp/probes/255 .scratch/wisp/probes/182/c1 .scratch/wisp/probes/180/c1 .scratch/wisp/probes/33/r10
.scratch/wisp/probes/180/c1
.scratch/wisp/probes/182/c1
.scratch/wisp/probes/255
.scratch/wisp/probes/33/r10
rc=0
```

## 相关工单（现量文件名）

```
255-config-says-these-sections-took-effect-immediately-panel-for-a-width-nobody-reads-while-the-panel-host-never-receives-config-and-tier-has-zero-readers.md
258-the-schema-calls-hotkey-hot-tier-but-the-resident-leg-builds-the-ball-from-default-hotkeys-and-the-only-reloader-caller-is-balldebug.md
rc=0
```
