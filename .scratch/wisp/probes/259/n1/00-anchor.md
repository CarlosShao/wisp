# 259-n1 — 起手五把尺（原样，现跑现量）

腿代号 `259-n1`；写面**只一枚**＝`.scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md`（末尾追加一节）。
本目录只准 `.md`（派单约束 4）；本目录里的尺一律**带 `[[:space:]]*`**（不带会整枚漏掉缩进项）。
所有命令都在仓库根 `D:\work\workspace\projects plans\Wisp` 下跑（CWD＝仓库根，具名声明；按 A715 定式①，复认名册／行内容一律同场再跑一次 `git show HEAD:` 快照对照）。

## 尺 1 — 现量时刻

```
$ date '+%Y-%m-%d %H:%M %z'
2026-10-08 13:26 +0800
```

（同一把尺在写本节时复跑＝`2026-10-08 13:29 +0800`；追加节的节头时刻取写盘那一刻的现量，⛔ 不估。）

## 尺 2 — HEAD 与分支

```
$ git log --oneline -1
53d73db4 A715 落账：收我自己那枚票 242 射程注（cade79b3，⛔ 零翻框、原判据一字未改，……）   [长单行，此处只截开头]
$ git rev-parse --abbrev-ref HEAD
dev
$ git rev-parse HEAD
53d73db4f66cdfc24abb2a209d3e00769fbc06da
```

## 尺 3 — 票面工作树状态（必须干净）

```
$ git status --porcelain -- .scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md
(回空＝干净，rc=0)
```

⇒ 干净，**不停手**（派单：若不干净立刻停手报回）。

## 尺 4 — 改前 `wc -l`

```
$ wc -l .scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md
45 .scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md
```

改前真尾部（`tail -5` 见 `evidence.md` 第 7 节，末字节 `od -c`＝`\n`）⇒ 追加只会落在**第 46 行起**，⛔ 不动 1–45 任何一行。

## 尺 5 — 改前两把框尺

```
$ grep -cE '^[[:space:]]*- \[ \]' .scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md
5
$ grep -cE '^[[:space:]]*- \[x\]' .scratch/wisp/issues/259-grant-binding-identity-ruler-holes.md
1
```

⇒ 改前 **未勾 5 枚／已勾 1 枚**（已勾那枚＝AC#0；未勾 5 枚＝AC#1～AC#5）。
本腿⛔ 一枚框都不翻、⛔ 不改 AC#0 那个已勾的框、⛔ 不改任何原判据文字；改后这两把尺**必须仍是 5 与 1**。
