307-r1 AC#2/AC#3/AC#4/AC#5 凭据（仓外导出树做的四形变异）

- 20-ac2-shapes.md：(a) 只把 :428 退回 isUnder ⇒ 5 枚指名全红（红句逐字；
  具名报回一处：TestSyncFixtureFallbackAndMatch 在这一形下红在 :175 而非 :162，
  因 :422 已被修好、:162 改由根成员那一支拿到 Sync:true）；
  (b) sha256 + cmp -s rc=0 逐字节还原 ⇒ 复绿（round-trip 自检 = pristine 变异同 sha）；
  (c) 中间形 = 只换 :428 ⇒ 具名回答：Sync:true 来自 :428（Root.Provider=sync-suspect /
  Source=suspect-fallback / Path=折过的 s.home），编排者读码推到成立，导出树里的打印用例
  把整枚 st.Root 落进件里（zz307probe 只活在仓外树，仓内零新 .go）；
  (d) 恒真形 || true ⇒ 5 枚里响 3 枚，具名交出不响的 2 枚（体内只有正向断言、天生无信号），
  包级另有 8 枚既有钉一起响。
- 30-ac3-fixture.md：新夹具 t.Skip = 0、t.Fatal = 15；在无别名默认台面上 fullrevert ⇒ FAIL=2、
  带修复 ⇒ PASS=2 ⇒ 夹具不依赖机器形状；★必答那一问答"会"，尺 = 改前默认台面那 5 枚
  --- PASS 且零 SKIP 零失败消息。
- 40-ac4-gates.md：rc_build=0 / rc_vet=0 / rc_d22scan=0；格式两把并排（工作树 1 枚 =
  i/lf w/crlf 的 provenance.go 假枚，HEAD blob 导出树 0 枚；gofumpt 同）；
  具名报回一处派单内部两条相撞 ⇒ AC#4/AC#5 的框本腿未翻（票面文件不在授权写面内）。
- 50-ac5-ledgers.md：三本账逐枚"⛔ 翻谁的格"；票 72 = 4/2、票 252 = 4/2、票 307 = 1/5 现量，
  git status --porcelain -- .scratch/wisp/issues docs = 0 行 ⇒ 本腿⛔ 动过任何一本账。
- 具名踩过：区间尺 git show cb7bb3c..HEAD 把编排者自己那两笔（249753a4 / 858ee7f4）数了进来
  ⇒ 印证派单"⛔ 区间尺"那一条；越界账只按逐笔名册算。
- 本笔零 push；写面 = probes/307/r1/**。
