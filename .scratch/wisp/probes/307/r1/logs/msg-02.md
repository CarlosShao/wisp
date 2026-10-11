票 307 AC#1 落地：syncdirs 同步判定的两支都换成包内既有的形比较族

- internal/risk/syncdirs.go match()：把 isUnder(cand, e.canon) 与
  isUnder(cand, s.home) 两支都换成 aUnder(candForms, anchorForms(...))，
  candForms = formsOf(cand) 在根循环之前折一次（编排者裁语＝形 C 换比较函数、
  两处出口同换）。候选侧多出的那枚 pathForms 是本落点唯一的新写法。
- 存值一字节未动：syncdirs.go:88 / provenance.go:345 的 home 赋值点、
  syncdirs.go:152/156/160 的 c.canon 赋值点、finalize()、
  resolveTarget() 里 errTargetUnverified / ErrReparseDenied 两句 fail-closed。
- 未新造折法（只复用 formsOf / anchorForms / aUnder / aRelates /
  uncertainAnchorMiss），risk 包非测试面新增 filepath.Clean|Abs = 0 枚。
- 新增 internal/risk/syncdirs_alias_307_windows_test.go（//go:build windows）：
  自建 8.3 别名世界的两枚夹具，取不到短形即 t.Fatal（票 115 定式，零 t.Skip）。
  syncdirs_test.go 那 5 枚未加 tag、未改一字；syncdirs_other_test.go:99 未动。
- AC#1 凭据 = probes/307/r1/10-ac1-colours.md：真别名台面 5 枚全绿
  （FAIL 12 -> 7，余 7 枚 = 甲族，归票 72），本机默认台面 FAIL 0 不变，
  两台面各两发取交集 = 新增红 0 枚。
- ⛔ 本格自勾 AC#1/AC#2/AC#3（裁决者 != 实现者）。
