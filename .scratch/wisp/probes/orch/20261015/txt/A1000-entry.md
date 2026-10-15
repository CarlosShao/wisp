## A1000（2026-10-15 23:5x +08，编排者）放 317-v1＋★自报一枚我漏接了一整天的格＋结案票残余普查（现跑）＋票池分母现量＋销一笔过期任务（294）

**① ★自报（账在我）**：票 317 的实现腿 `317-r1` 于 **10-14 13:48:22** 交件、改动**当天 13:50 就已入库**（`git log -1 78805129`＝`ticket 317 (317-r1): compose the test temp root with the native separator…`，现跑命中），它交件正文末尾自己写了 `next＝非实现者腿判 AC#1/AC#2`。**那句我没接住** ⇒ 票面两格（现量在票面 14／15 行）挂了 26 小时仍 `- [ ]`。
- 定式（写下来防再犯）：**交件里带 `next＝` 的，收件当场转成排程或当场派**；⛔ 让"腿等我"变成"我忘了"。这一枚尤其贵——它判的是**测量完整性**（`GOTMPDIR` 分隔符写法会让 windows 腿的 junction／路径字面量族用例**假红**），拖一天＝本编队每天在这种坏尺上读数的风险多一天。
- 处置＝放 `317-v1`（只读验收腿，面＝`scripts/**` 读数＋`./internal/risk/` 两形对照；⛔ 改脚本），三把尺现跑：`78805129` 真在／`git status --porcelain -- scripts tools`＝**0 行**／CPU 48%＜70%。存档＝批次提交 `55e11b2f`（`317v1-dispatch.txt`＋`317v1-send.txt`）。

**② 派单里按现量写的三条更正（三条都交腿自己复核，⛔ 我盖章）**
- 票面"元凶行＝`scripts/portable-tests.sh:144`"已被这枚修改动过行位 ⇒ 今天现量在 `:154`–`:156`（`:155` 基座 `"$LOCALAPPDATA\\wisp-test-gotmp"`、`:156` 尾接 `"$gotmp_base\\portable-…"`＝两处都原生 ⇒ 修法确实在盘上）。
- 实现腿报"三处 `wisp_311_ref_path` 判不修"，我 23:5x 现量＝**4 处赋值**：`scripts/portable-tests.sh:818`／`scripts/winsec-tests.sh:174`／`scripts/wisp-cli-tests.sh:186`／`tools/d22scan/runtests.sh:155`（全同形 `${LOCALAPPDATA//\\//}`）⇒ 差的那一处是"写面外所以没算"还是"真漏了一处"，交腿裁。
- `tools/d22scan/runtests.sh:92` 今天现量是**混形**（基座用 `\\`、尾接用 `/`）而⛔ 实现腿措辞的"仍折 posix 形"⇒ "会不会爆"要按混形重述。

**③ 结案票残余普查（我自己现跑，一把尺具名）**
- 尺＝对 `.scratch/wisp/issues/*-done.md` 全量跑 `grep -cE '^[[:space:]]*- \[ \]'`（这把尺含缩进项，＝10-02 起的换形尺；⛔ 与旧的 `^- [ ]` 相减）。分母＝**133 枚 `-done`**（⛔ 计 `README.md`）。
- 读数＝**命中 3 枚文件／共 4 格**：`115-…-done.md` 的 AC#4（变异反证）与 AC#6（门禁四数点名）各 1 格、`92-panel-composer-mode-attachments-workspace-done.md` 的 AC#5（变异三向）1 格、`07-ball-state-machine-core-done.md` 的"Visual: 20 states…human screenshot review"1 格。
- 逐枚定性（⛔ 批量撤名、⛔ 替机主判签收）：`07` 那格要求的是**人看图**——机主 09-30 已把样式/外观后置（"形状算你过关，好不好看以后再说"），我不催签，只登记它仍缺；`115`／`92` 那三格缺的都是**非实现者读数**⇒ 归队列，等编队面空再派，不自己盖章。
- ★结论说实话：结案票里的残余比我以为的小得多（3/133 枚），所以这一族今天不值得占一枚腿。

**④ 销一笔过期任务（先跑尺再销，⛔ 按旧清单复读）**：任务里挂着"收 `294-r1`（若识别为死腿则补派 r1b）"。盘上三把尺＝① `.scratch/wisp/probes/294/` 下 `a1 r1 r1b v1` 四节都在；② 证据件 `docs/evidence/s1/ticket294-294-v1-accept-v1.md` 在库；③ 票名已带 `-done` 且台账 `A880`（10-13 21:23）就是那次收件记录 ⇒ **这条是过期账，已发生**，就地销掉，⛔ 再排。

**⑤ 编队与排程（枚数＋写面类型一起报）**
- 在飞 5 枚＝**写码 2**（`322-r1`/`internal/winsec` 现量 6 行未提交、`166-r2`/`internal/session` 现量 1 行未提交）＋**只读验收 3**（`287-v1`/`internal/tools`、`324-v1b`/`cmd/wisp` 定向跑、`317-v1`/`scripts`＋`internal/risk` 读数）。
- 队列按住不变：`120-v2` 按 `322-r1` 交完 → `165-v1` 按 `324-v1b` 交完 → `135-r1`（面含 `cmd/wisp`，按 `324-v1b` 交完＋CPU＜70%）。
- 未推＝本地领先 `origin/dev` **13 笔**（尺＝`git rev-list --count origin/dev..HEAD`）；⛔ 现在推的理由具名在 §4.0ds（推会自启本机 `slo-full` 抢 CPU＋CI 色要一次对拉才好读）。
- 票池分母现量：`.scratch/wisp/issues/` 共 **311 枚**（⛔ README）＝`-done` 133 ＋开放 178。⚠ 开放数比 10-02 的 154 涨 24，那是"立新票"涨的、⛔ 读成进度倒退；引用这两个数时必须带今天这把尺的射程（同一换形尺）。
