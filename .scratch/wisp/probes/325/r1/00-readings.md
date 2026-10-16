325-r1 读数件落盘：混形四枚真红（未修码）＋物理退回反控＋还原哈希

台面＝仓外导出树（git archive 46007856 | tar -x 到 C:/tmp/wisp325r1/pre，六件套逐枚验在、
find -name '*.go' | wc -l ＝ 1090 枚非空之后才信它的数）；改后树＝同法导出 428ef21e。
共享工作树这一刻有 322-r1b 在写 internal/winsec，所以整树/整档读数一律没在仓里跑。

logs/ 现量（逐枚）：
- a1-pre-d22scan.txt：改前（混形）sh scripts/d22scan.sh rc=0 clean，发射行逐字
  GOTMPDIR=C:\\Users\\swq\\AppData\\Local\\wisp-test-gotmp/d22scan-20261016-093651-329
  ＋正控制那步 PASS=35 FAIL=0 SKIP=0 === RUN=77 ⇒ 门绿着、坑没被看见（缺世界实证）
- a1-post-d22scan.txt：改后 rc=0 clean，同一行只差一个分隔符（\\d22scan-），四数同形
- s1-mixed-emitted.txt：未修码上由 runtests.sh 自己生成的混形根打那四枚 ⇒ rc=1、
  顶层 PASS=0 FAIL=4 SKIP=0 === RUN=4（:59／:116 回 Parameter format not correct，
  :182／:286 在前置条件文案里回显分隔符不一致）
- a2-post-green-four.txt：改后同构造 ⇒ rc=0、PASS=4 FAIL=0
- a2-counter-control-stdout.txt＋a2-revert-mixed-red.txt＋a2-restore-green-four.txt＋
  a2-90/91/92-sha-*.txt：在导出树上把那一行**物理**退回混形（⛔ -overlay）⇒
  mixed-shape-hits=1／native-shape-hits=0（突变真落了）⇒ rc=1 四红 ⇒ 用
  git show 428ef21e 的字节还原 ⇒ sha256 回到 ad6f6a8d…（＝提交态那枚）⇒ rc=0 四绿
- a3-post-scope-windows.txt＋a3-per-package-four-numbers.txt：AC#3 那发 rc=0，
  GOTMPDIR 行逐字原生形，逐包四数 11 枚（派单写"八包"＝过期计数，现量 11），
  总计 === RUN=662 / --- PASS=462 / FAIL=0 / SKIP=0，红名册＝空
- a4-gates-stdout.txt＋a4-g2/g3/g4-*.txt：四枚门禁逐枚 rc=0（bash -n／d22scan.sh／
  长度门 --with-self-test VERDICT GREEN／go build ./...）；gofumpt 这台机器 which rc=1、
  四处候选目录零命中 ⇒ 本腿写面 0 枚 .go，⛔ 用 gofmt 的 0 冒充它的 0
- 仪器＝ac2-counter-control.sh／ac4-gates.sh（只建不删）；commit-ac1.txt＝AC#1 那笔正文

只 commit，push 归编排者。
