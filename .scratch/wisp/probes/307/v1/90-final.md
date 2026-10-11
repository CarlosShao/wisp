# 307-v1 · 90-final —— 收口：与派单不符处、★两枚必答、具名欠读数

## 1. 与派单不符处（一律以原文／盘上为准，本腿具名报回）

1. ★**笔数与 HEAD 都腐烂**：派单 §0 说落地腿"已交完**三笔**"、"当时 HEAD `183032c1`"。
   盘上现量＝307-r1 名下**五笔**（`c6691756` 起手锚、`51a0a0f2` 产码、`183032c1` 凭据、`fb6a171e` 收口、**`d8743be5` 10:28「收口补件：共享工作树脏面存档＋节号更正」**），
   本腿起手 HEAD＝`d8743be5`。`d8743be5` 写面⛔ 涉产码（尺＝该笔名册＝`41-boundaries.md` ＋ `logs/msg-05.md`），
   且本腿复现产码 blob 在 `d8743be5` 与现在（`HEAD:internal/risk/syncdirs.go`）**同哈希** `7125b08a…` ⇒ 本票"改后"读数⛔ 因多这一笔而失真。
2. 派单 §3 `AC#2(a)` 那句"5 枚红（红在 `:162`）"⇄ 本腿量到**红在 `:175`** ⇒ 腿报的行号换形＝**事实**，派单那半句⛔ 准（见 `20-ac2.md` §1）。
3. **分母差异（⛔ 冲突、只是两把尺）**：腿那四发变异一律用 `-run TestSync`（尺＝`probes/307/r1/logs/run-mutation.ps1` 的默认 `$Run='TestSync'` ⇄ `20-ac2-shapes.md` §0），
   本腿把 **(a)(c) 两形换成更小的指名分母**（乙 5＋新 2＝7 枚）做归因，另把 **(d) 与 fullrevert 两形沿用腿的 `-run TestSync`** 同分母对拉
   ⇒ `v1-mut-d-aliased.txt` **11 枚红** ⇄ 腿 11 枚、`v1-mut-fullrevert-aliased.txt` **7 枚** ⇄ 腿 7 枚，**逐枚同名**；
   本腿那把 7 枚分母的 (a) 给 6 枚红、(c) 给 1 枚红＝**归因比腿那一发更尖**，⛔ 有一格靠"分母不同"糊过去。
4. 派单 §2 那枚"⚠ 0 字节尺必带 `-type f`"本腿照做（`find … -type f -size 0` ⇒ r1＝0 枚、v1＝0 枚）；派单 §5 "⛔ 在仓内建 worktree／checkout"本腿全程⛔ 触发（变异只在 `C:\Users\swq\tmp\wisp-307-v1-*` 四棵仓外树里做）。

## 2. ★必答①：这次修复在**这台机器的生产路径**上现在观测得到颜色差别吗？

**⛔ 观测不到。** 三把尺都是本腿自己跑的（件＝`logs/v1-production-rulers.txt`）：

1. **生产唯一消费者**＝`grep -rn --include=*.go '\.IsSyncPath(' internal cmd \| grep -v _test.go` ⇒ **1 行**＝`internal/risk/provenance.go:928`（`if p.IsSyncPath(t).Sync { return true, ChSyncWrite }`）⇒ ⛔ 是"建了没接"。
2. **生产唯一阳性出口＝`:428`**：非测试 `NewProvenance(` 调用者＝**2 枚**（`cmd/wisp/run.go:743`、`cmd/wisp/panel_assets.go:239`），
   本腿把两行正文逐字取回＝**都只传 `risk.ProvOptions{NoProbe: true}`**；
   `grep -rn 'SyncRoots:'` 与 `'SyncFixturePath:'` 非测试面各 **0** 命中 ⇒ `roots` 空 ⇒ `match()` 里 `for _, e := range s.roots` 那个循环（内容锚＝`if aUnder(candForms, anchorForms(e.canon))`）**生产零次迭代**；
   而 `home` 的两枚赋值点（`syncdirs.go:88`、`provenance.go:345`，本腿复看 blob＝两边正文逐字同形）都只过 `normPath`。
3. **这台机器的 profile 链⛔ 带别名**＝本腿自己用 COM 逐段量（⛔ 引用编排者的数）：
   `C:\Users\swq` → `ShortPath`＝`C:\Users\swq`（**同字**）；`APPDATA`＝`C:\Users\swq\AppData\Roaming` → `ShortPath` **同字** ⇒ `rc_com_short=0`。

⇒ **机制层结论**：生产面上那次写的候选（`resolveTarget` 折长）与锚（`s.home` 只 `normPath`）在**这台机器**上本来就是同一串拼法 ⇒
`isUnder` 与 `aUnder` 给出**同一个答案** ⇒ 改动在生产路径上**⛔ 改变颜色**（既⛔ 多拦、也⛔ 少拦）。
⇒ ⛔ 此，**`AC#1` 那一格的凭据性质本腿裁＝测试台面（test-bench）凭据，⛔ 生产台面凭据**。这句话必须跟着那枚勾走：
修的是**真缺陷**（生产消费者真在读这枚判定），但**"5 枚指名用例由红转绿"这一发颜色只存在于带 8.3 别名的测试台面上**。
⛔ 拿"CI 托管 runner 的 `TEMP` 是 `C:\Users\RUNNER~1`"替生产背书——那是**另一台机器上的另一个消费者**，而且 runner 的别名在 `TEMP` 那一支、⛔ 在 `USERPROFILE` 那一支（本腿 P7 那把尺量的正是后者）。

⇒ **两枚"生产面哪天会观测到差别"的形状，具名登记（⛔ 本票能修，属下一程／其它票射程）**：
- 形状 A＝**用户 profile 链自带 8.3 别名**的机器（用户名 >8 字符且卷开着 8.3 生成，`C:\Users\stephenw ⇒ STEPHENW~1` 那一族）⇒ 生产 `:428` 从此拦得住；
  这一枚的凭据本票交⛔ 了（本⛔ 第二台机器），本腿裁它⛔ 阻 `AC#1`（票面判据写的是测试台面那 5 枚＋default 台面⛔ 变红，两枚都到了）。
- 形状 B＝**`uncertainAnchorMiss` 那支 fail-closed 在生产上第一次真的可能生效**：换形之后候选侧与锚侧各带 `root`，
  若某一侧展开停在"存在但⛔ 肯开"的组件上（deny-filtered 目录／悬空 reparse 点），一次原本⛔ 命中的比较会变成命中 ⇒ **只朝"多疑"动、⛔ 朝"少拦"动**
  （本腿读了 `blacklist.go:314`–`:322` 的边界：需要两侧 root 非空、且⛔ 都 certain，且互为祖先）。
  这台机器⛔ 造得出那台机器吗？造得出（票 72 那一族 junction 夹具就是那条路），但那是**新仪器**，⛔ 本票授权 ⇒ **建议＝记一枚欠读数，落点归编排者定**（本腿倾向票 72 的下一程：它的射程正是"别名/junction 台面上的期望侧与仪器形状"）。

## 3. ★必答②：⛔ 翻任何勾

- 本腿全程⛔ 动任何一枚框、⛔ 改任何一句原句、⛔ 动票 72／252／302 的框（尺＝本腿自己每笔 `git show --name-only --format=` 名册，五笔全在 `probes/307/v1/**`＋`docs/evidence/s1/**`＋票 307 末尾那一枚 append-only bullet 之内）。
- 翻勾动作留给编排者，`50-ac5.md` §4 那张六行表就是给翻勾用的对照（每行带尺名＋射程＋出处三态）。

## 4. 具名欠读数（本腿⛔ 交的、本腿交⛔ 了的都写在这）

1. ★托管 CI runner（`C:\Users\RUNNER~1` 那一台）上这 5 枚＋2 枚夹具的颜色＝**本腿⛔ 取**（⛔ 本腿⛔ push，`@{u}`＝`cnb/dev` 落后 1032 枚）。
2. ★生产面形状 A（profile 链带别名的机器）＝⛔ 第二台机器 ⇒ 本腿只交⛔ 了"机制＋本机无色"这一半（§2）。
3. ★`uncertainAnchorMiss` 在生产上第一次生效（§2 形状 B）＝⛔ 仪器造过；本腿只读码到边界条件。
4. ★`alias307Pair` 里那枚 `t.Fatalf`（无别名即 Fatal）＝本机从未执行（需关卷的 8.3 生成＝改机器配置，⛔ 本票授权）；形状＝读码推到。
5. ★代价值那枚（`formsOf(cand)` 多一次祖先走步）＝本腿那把尺给出**反向**数（改后总耗时更短 1.19 s）⇒ 证⛔ 了"这把仪器量不出那枚代价"，⛔ 能用来主张代价为零⛔ 能主张代价为负；结构性上界（每查询一次、`anchorCache` 缓存锚侧）＝读码核过。
6. ★票 307 `AC#5` 那"合销"动作**尚未落进台账**（`A817` 欠账⑥那句仍在 `:15451`/`:15487`）＝编排者车道，⛔ 本腿写面。
7. 门禁三把之外本腿⛔ 跑的：`gofumpt` 两把（腿跑了、本腿⛔ 装那把工具⇒ 只交 `gofmt` 两把并排）；
   ⛔ 跑过 `cmd/wisp`（本票⛔ 需要）；⛔ 跑过 `git add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`／`push`／`fetch`。

## 5. 本腿五笔名册（commit-first 一路交到底）

| 笔 | 内容 |
|---|---|
| `c04ced26` | `00-anchor.md`（起手锚，长跑之前） |
| `ed1d2bad` | `10-ac1.md` ＋ `20-ac2.md` ＋ 全部 `logs/` 原始件 |
| （本笔） | `30-ac3.md` ＋ `40-ac4.md` ＋ `50-ac5.md` ＋ `90-final.md` ＋ `v1-production-rulers.txt` |
| （裁决表笔） | `docs/evidence/s1/ticket307-307-r1-accept-v1.md` |
| （Progress 笔） | 票 307 末尾 append 一条 bullet（numstat 必 0 删除） |
