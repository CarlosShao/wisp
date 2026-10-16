收件四笔：317-v1 判语落地⇒我翻两格并加 AC#3 闸门｜324-v1b 断流但六格判语已写完⇒代提补齐＋翻六格｜★立票 325（第四处混形活坑）｜我的两发现量复跑

一、票 317（测量完整性·GOTMPDIR 分隔符）
- 317-v1（非实现者）判语表已在盘并自提交 ceaa50ac：docs/evidence/s1/317-gotmp-separator.md（16489 字节，占位符 0）＝AC#1 成立／AC#2 成立／第四处同形建议立新票。
- 我自己复跑成对读数（仓外导出树、同一个目录只差分隔符）：
  GOTMPDIR='C:\tmp\wisp317orch\pair'  => RUN=4 PASS=4 FAIL=0 SKIP=0 rc=0
  GOTMPDIR='C:/tmp/wisp317orch/pair'   => RUN=4 PASS=0 FAIL=4 rc=1
  红句同行同号（internal/risk/syncdirs_ancestor_actable_leg_116_test.go:182／:286 ＋ reparse_ticket105:59／:116 的 Invalid switch - "tmp".）。
  读数件归档在 .scratch/wisp/probes/orch/20261016/logs/（00-tree 证导出树 1090 枚 .go 非空）。
- 三处赋值字节现读＝原生形（scripts/portable-tests.sh:155-156／winsec-tests.sh:69／wisp-cli-tests.sh:82）。
- 因此翻 AC#1／AC#2 两格。★不改 -done：新加 AC#3＝整条 --scope=windows 八包读数一发（317-v1 具名说这一发它没跑），并把 Status 那句过期"未派落地腿"就地更正。

二、票 324（名册尺按文本匹配）
- 324-v1b 死于断流（宿主回 Model response interrupted），但六格判语已全部写完、只差提交：docs/evidence/s1/324-v1b-acceptance.md 15795 字节＋三枚突变读数件（logs/01 18214／02 35472／03 7663）。代提补齐三处写面（含票尾一条 bullet，我写、标〔编排者代提〕）。
- 判语＝AC#0/1/2/3/4/5 六格全成立 ⇒ 我翻六格；★不改 -done（腿具名交出五笔残余：8 枚 stdout／开窗族用例缺真窗、甲形与 [voice] 镜像读者重裁移交我、运行时回执缺窗台件、CI 色不在其射程、作差那把我未复）。
- 我只独立复了 AC#0/AC#3 那把尺：第一发＝台面故障（exit status 0xc0000135＋RUN=0，三枚 sherpa DLL 不在 PATH，已登记 .scratch/wisp/probes/orch/20261016/logs/orch-324-roster-attempt1-ruler-fault.txt），把 third_party/sherpa-onnx 加进 PATH 后＝RUN=1 PASS=1 FAIL=0 rc=0（同目录 attempt2 件）。

三、立票 325（tools/d22scan/runtests.sh:92 混形活坑）
- 现量：那一行基座用 \\（原生）、尾接用 /（正斜杠）＝混形；已修三处对照都尾接 \\。今天不爆只因该档射程里没有那四枚 junction 用例＝"坑活着、暂时没人踩"。
- 票面 AC#1 改形（⛔ 顺手把 wisp_311_ref_path 四处对照位一起改一致）／AC#2 成对反控（物理改字节、⛔ -overlay；若射程真不含那四枚就具名写"缺世界"）／AC#3 顺带交 317 AC#3 那发八包读数／AC#4 门禁（具名两枚台面雷）／AC#5 账（317 那句"同形、今天未爆"由我重述）。

四、我自己的两笔错（记在这笔里，⛔ amend 旧笔）
- c2b1ade5 正文那句"盘上唯一一枚中间态"不实（166-r2 的证据件与 23 枚读数件当时就在盘上）＝已在 464f5549 追加更正并补齐三处写面；同一枚坑今天第二次抓到（A999 那次＝165-r2）。
- 我那份复跑脚本的汇总行把 rc 写成了 grep -c '^rc='（数行数）而不是取值 ⇒ 件照原样留档（⛔ 删坏件），真 rc 我从日志逐行另读（nat rc=0／psx rc=1），本笔正文里的两行就是更正后的读数。

编队：此刻零枚在飞（287-v1／322-r1／166-r2／324-v1b 四枚终态＝死腿或已代提，317-v1＝已交件并自提交 ceaa50ac）。未推 22 笔（起笔前 `git rev-list --count origin/dev..HEAD` 现量，origin/dev＝`adec1db1`、本笔落完＝23）；⛔ 现在推的具名理由沿用 §4.0ds（推会自启本机 slo-full 抢 CPU），释放条件＝下一批写腿交完之后一次推＋推前核 fast-forward、⛔ force、只 origin/dev。
本笔零产码改动（只动 issues/evidence/probes）；⛔ push。
