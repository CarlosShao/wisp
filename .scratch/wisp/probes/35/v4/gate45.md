35-v4 闸门件（第 45 次调用线到达时落笔；编排者的 45 闸门要求＝"判不动那节＋起手锚"写满并 commit）

## 起手锚
- 锚笔＝commit `1d719434`（只含 `logs/anchor.txt` 一枚件），HEAD 起跑 `168a91f0747a7f71c36344dd48f75edd98fe7729`，分支 dev。
- 禁区现量空：`git diff --stat 2fc5f5c9..HEAD -- cmd/wisp/panel_host_windows.go internal/panel/ frontend/ docs/specs/ docs/PLAN.md`＝零行。
- 夹具与产码 worktree blob == HEAD blob（`31c96f37d32adcd1625a2c01748d5b30016227b2` / `26b5de83b93a9a141f1546dc19a9b03ffa45ede0`），两枚 CR 计数各 0。
- 全程零 push、零真窗（⛔ 没有任何一条命令带 `-tags winlive`；winlive 那枚文件由 build tag 排除，见 `logs/roster-coverage.txt`）。

## 判不动那节（今天交件时也只能这么写，缺的是别人的件不是我少跑一把）
1. **"WebView2 真绑 receiver / 真不许页面换 postMessage"这两格＝本腿判不动。**
   缺的东西具名：`-tags winlive` 的真窗读数（票 35 `:52`＋台账 `A684`／`A686` 那一格），而派单硬禁今天开窗。
   解释器里"丢了 receiver 会 panic"是我这台夹具自己写的规则，不是 WebView2 的行为记录。
   ⇒ 这一格只能判"两格不许互相作证"，不能判"这条规矩在真浏览器里成立"。
2. **"整包里没有第二枚用例会被这批坏法染色"＝我只做了名册级论证，没做全量复跑。**
   缺的东西具名：一次 full-package 逐突变复跑（本机实测 150/346 枚用了 ~14 分钟，×17 发＝不可行）。
   替代尺＝`logs/roster-coverage.txt`（全仓只有 3 枚 test 文件＋1 枚 winlive 文件读 `installPanelTransport`/
   `panelPostMessageForwardInit`），⛔ 那是"引用面"的读数不是"跑过"的读数。
3. **`writable`/`configurable` 在 WebView2 的真实形态＝判不动**，因为页面上没有 `Object.defineProperty`
   （`markNonWritable` 是 Go 侧台件）。这一样在 r4 的 `impl.md §4` 里已具名，我没有把它变成可判的东西。

## 本腿自抓的两把尺的缺陷（先记账，免得读数被当凭据）
- pass 1 的全部十三发读数＝no-op overlay 造出来的假绿（`logs/self-catch-noop-overlay.txt`＋`logs/pass1-noopoverlay-*`）。
  正控（重跑腿自己的 ma-hook，本该红）才把它抓出来。pass 2 的读数以修好的台件重跑。
- 第一版台件想跑整包（不可行，见上第 2 条），已换成 12 枚名册＝与腿 r4 同一组名册，可比。

## 目前进度
- 15 枚一行突变全部现量 `changed-line-count=2`（`logs/mut-landing.txt`，`all-one-line=YES`）。
- pass 2 十六发＋baseline 在跑（`logs/mut-summary.txt`）。四判正文在 `verdict.md`。

rc=0
