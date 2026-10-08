# 起手锚 — 腿 `next-instruments-brief-1`（只读备料腿）

- 取锚命令现跑：`date` / `git log -1 --format='%H %ad %s'` / `git status --porcelain | head -20` / `git branch --show-current`
- 取锚时刻：`Thu Oct  8 19:01:32 CST 2026`
- HEAD：`5f52f310f1940a9cdb4b043f2024d8d71eafee7f` · `Thu Oct 8 19:01:13 2026 +0800` · `A736 落账：收 gate-snapshot-1（50b34971/8ace2272，只读零修）＋我复跑同数——d22scan rc=0 clean（分母@50b34971：bans1-5 internal/=228 cmd/=38、#6 f/=85、#7 tools/=23、#8 design/=39 f/=85 internal/=516 cmd/=113）；gofmt -l cmd internal 7 命中我逐枚问来历＝6 枚既有（5e8748b3 10-03×2／7f9d6e40 10-07／d03d166f 10-04／5413f46d 10-04×2）＋1 枚今天 3a343bc7 10-08 09:37（已提交树净）⇒ 无'今天新写入未提交'形，⛔ 不擅自修（改=顺手改；且 CI 那步量 gofumpt 非 gofmt）；go func( 词面总 1（唯一位点 cmd/wisp/testdata/esclistener/main.go）＝不是判定要人看；四路径 porcelain 18:58 全空（我 19:0x 复量见 259-r4 正写 subagent_selfapproval_197_test.go＝非异常）；未推 427→431；零翻框零 push`
- 分支：`dev`
- `git status --porcelain | head -20`（只登记，一字未动；含他腿在写的 `cmd/wisp/subagent_selfapproval_197_test.go`＝腿 259-r4）：
  - ` M .gitignore`
  - ` M .scratch/wisp/probes/152/my152.py`
  - ` M .scratch/wisp/probes/161/r6/logs/flip-1.txt` …（flip-2~6、flip-baseline、flip-restored 同形）
  - ` M .scratch/wisp/probes/242/r3/logs/probe-routed.txt`
  - ` M .scratch/wisp/probes/268/v1/evidence.md`
  - ` M cmd/wisp/subagent_selfapproval_197_test.go`
  - ` D design/assets/base.css` 等 design/** 删除项
  - ` M design/doubao/README.md` 等 design/** 修改项

## 任务与写面

- 任务：为两处已知盲区写"下一枚仪器该怎么装"的**派单级简案**——案一＝票 181 的 M4 盲区（`internal/panel/workspace.go:141` 成功路径填充点，出处 `.scratch/wisp/probes/181/v3/evidence.md` §5）；案二＝票 259 `AC#2` 两发"残余风险"突变（非实现者腿 `259-v1` 自报没种的两发，出处 `.scratch/wisp/probes/259/v1/01-evidence.md`）。
- 写面：只新建本目录 `.scratch/wisp/probes/next-instruments-brief-1/*.md`；零产码、零突变、零 Go 命令、零翻框、零 push；他腿/票面/台账一字未动。
- 禁区：`frontend/**`／`design/**` 不读不引；三枚冻结件一字不动；不新建票/判据；`git add -A`／`.`／`-a`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean` 全禁。
