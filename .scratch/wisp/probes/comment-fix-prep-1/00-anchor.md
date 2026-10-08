# 00-anchor — comment-fix-prep-1 起手锚

- 腿名：`comment-fix-prep-1`（准备料腿，只交料不改产码）
- 分支：`dev`
- 现量时间：`2026-10-08 18:27:18 +0800`
- 现量 HEAD：`9e8477edc690c46f3540f97069a1009c05ea2fc4`（`Thu Oct 8 18:26:30 2026 +0800`，题＝`167-a2b 两格读数：格一草稿=〔没字段〕（读者 0／写者=插座有接线断 panel_inbound.go:277 Message:nil）；格二崩溃自救=没造（出向读面 0 枚崩溃字段；G2 主尺 47 行=注释 32+产码 15，唯一产码族=配置 restart 档审计）`）
- 料源：`.scratch/wisp/probes/comment-truth-2/01-triage.md`（入库 commit `762b694e`，38 行）
- 本腿目标：为 01-triage 里 13 枚〔只欠文案〕逐处备可直接粘贴的替换文本 ⇒ `01-ready-to-apply.md`；⛔ 不碰 P14/P15（〔要动产码〕2 枚）。

## 起手时在飞（`git status --porcelain | head -20`，只登记）

```
 M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt
 M .scratch/wisp/probes/161/r6/logs/flip-2.txt
 M .scratch/wisp/probes/161/r6/logs/flip-3.txt
 M .scratch/wisp/probes/161/r6/logs/flip-4.txt
 M .scratch/wisp/probes/161/r6/logs/flip-5.txt
 M .scratch/wisp/probes/161/r6/logs/flip-6.txt
 M .scratch/wisp/probes/161/r6/logs/flip-baseline.txt
 M .scratch/wisp/probes/161/r6/logs/flip-restored.txt
 M .scratch/wisp/probes/242/r3/logs/probe-routed.txt
 M .scratch/wisp/probes/268/v1/evidence.md
 M cmd/wisp/panel_host_windows.go
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
```

（其余在飞未列，仅登记前 20 行；`frontend/**`、`design/**` 本腿不读不引，此处仅为 `git status` 原样登记。）
