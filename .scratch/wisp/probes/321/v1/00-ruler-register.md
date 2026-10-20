# 票 321 终裁腿 321-v1 尺册（每把＝尺名⇄命令形⇄射程⇄量的还是算的；★⛔ 抄任何一枚腿的数）

腿＝票 321 终裁腿 321-v1（非实现者·只读）｜时刻＝2026-10-20 21:5x +08｜仓＝`D:\work\workspace\projects plans\Wisp`
硬闸＝全程⛔ 动 `.go` 字节、⛔ `go test`／`go build`；尺全＝`git cat-file`／`git grep HEAD`／`git log`／`git ls-tree`／仓外副本比对／原件 grep／`stat`／`tasklist`
临时根＝`D:\tmp\wisp\20261020\321v1\`（⛔ 盘符根、只建不删；探针副本一律 `.go.txt` 后缀＝⛔ 叫 `.go`）
脚本底片＝本目录同层 `D:\tmp\wisp\20261020\321v1\sh\R*.sh`＋读数件 `R*.txt`（未入库的那批⛔ 当凭据引、本件只登尺名与读数）

| 尺名 | 命令形（可逐字重跑） | 射程 | 量/算 |
|---|---|---|---|
| R1 判据件在位 | `git cat-file -e HEAD:cmd/wisp/exec_holder_toggle_321_test.go`；`git rev-parse 58de465a:P`／`HEAD:P`／`git hash-object P` | 单枚件×多修订 | 量的 |
| R1b 谁动过该路径 | `git log --all --format=… -- <path>`（逐笔列名，⛔ 区间归因）＋正控＝同尺对不存在路径 | 全史 --all | 量的 |
| R1c blob 维度 | `git cat-file blob HEAD:P`→仓外文件→`wc -c`／`wc -l`／`tr -cd '\015'|wc -c`／`od -An -tx1|tr|grep -c '^0d$'`／`sha256sum` | 入库 blob | 量的 |
| R2 谁在读盘 | `git grep -n -F -e PAT HEAD -- 'cmd/wisp/*'`（PAT＝`parser.ParseFile`／`os.ReadDir`／`filepath.WalkDir`／`os.ReadFile`／`os.Stat`／`overlay`）＋正控 `package main` | tracked@HEAD | 量的 |
| R2f 提取器本体 | `git cat-file blob HEAD:cmd/wisp/resident_hotkey_258_test.go`⇄`…resident_ball_228_test.go`＋`sed` 行区间；被测量件锚 `grep -F -c` | 两枚宿主件＋产码件 | 量的＋读正文 |
| R3 名册三把 | `git log --all --format='%H%x09%s'` 单次导出→awk `$2 ~ /^锚/` ⇄ `index($2,"串")>0` ⇄ 腿码任意位；正控＝`^票 115 出表腿 115-v1` | 全史 5689 笔 | 量的 |
| R4 证据件维度普查 | 28 枚交件/票面/派单逐枚 `git cat-file blob HEAD:p > 仓外` 再量字节/行/CR/blob id | 两批全部交件面 | 量的 |
| R5 票 321 框终量 | `git cat-file blob HEAD:<票>`→`grep -c '^- \[ \]'`⇄`'^- [x\]'`⇄正控 `'^- \['`＋`git diff --numstat HEAD` | 票面 blob⇄工作树 | 量的 |
| R6 常量⇄能力面 | `sed` 行区间取常量段；`grep -n 'AC#1 RED'` 对 blob；`awk '/^func \(h \*panelHolder\) toggle/,/^}/'` 取体 | 判据件＋产码件 blob | 量的 |
| R7 仓外原件重量 | `D:\tmp\wisp\20261020\321r1\logs\*.txt` 逐枚 `grep -c '^=== RUN'`／`'^--- PASS'`／`'^--- FAIL'`／`'^--- SKIP'`＋`comm`/`diff` 逐名作差＋`cmp`⇄`sha256sum` 对 HEAD blob | 重放腿仓外底片 41 枚 | 量的 |
| R8 门禁窗口 | `git rev-parse 4e0ac985:P`⇄`HEAD:P` 逐枚；`git diff --name-only 4e0ac985 HEAD -- cmd/wisp/`；`git show --name-only` 逐笔；`stat -c %y` | 两枚 HEAD 的编译面 | 量的 |
| R9 旧批底片 | `git cat-file blob HEAD:<logs/…>` 逐枚全文（三态⇄反面尺⇄255 名册） | 旧批 15 枚件 | 量的＋读正文 |
| R10 占位符纠正 | 对 blob `grep -F -c` 五枚红句子串（源串形）＋两棵树 `cmp` 全树逐枚 | 判据件＋仓外两树 | 量的 |
| R11 255 名册件 | 两枚 `-names.txt` `cmp` ⇄去时间戳后 `diff` | 旧批 2 枚件 | 量的 |
| R12 编排者放腿底片 | `git show --name-only 1104fea8`／`cat-file blob HEAD:…/20261020/txt/321r1-dispatch.txt`／`41673704`／`f1415309` | 派单存档面 | 量的＋读正文 |
| R13b 台账⇄防双派尺原文 | `grep -n 'A1068'`（★v1 用 `cut -c` 截中文＝非法 UTF-8，v2 改 `LC_ALL=C.UTF-8`＋`awk substr`）＋`cat-file blob HEAD:.scratch/wisp/probes/orch/predispatch-rulers.sh` | 台账＋尺件 | 量的 |
| R14/R15 A1068 定位 | `git grep -l -F A1068 HEAD`⇄`A1069`⇄`A1070`（正控）＋`grep -n '^## A10'` 序号列 | tracked@HEAD | 量的 |
| R16 防双派复量 | `git log 1104fea8 --grep='^票 321 落地腿 321-r1' --oneline|wc -l` ⇄ 严格形 `git log 1104fea8^ --format='%s'|grep -c '^…'`；`--since/--until` 限定 `cmd/wisp/`；票面最后一次入库笔 | 到 18:45 那刻的历史 | 量的 |
| R18/R19/R22 依赖面 | `git grep -l -F 'wisp/internal/winsec' HEAD -- '*.go'` 分包计数；逐件 `awk '/^import \(/,/^\)/'`＋`grep -c`；`git show --name-only` 定 winsec 漂移笔 | HEAD tracked | 量的（闭包那半＝算的，凭名册） |
| R21 变异复现 | 仓外副本：`git cat-file blob 58de465a:R`／`HEAD:R`→`sed "${LINE}s/holder.toggle/holder.show/"`（行号由 `grep -n` 现取＝内容锚）→`sha256sum`＋`diff` 变动行数＋`cmp` | 两批被测量件 | 量的（⛔ 跑测试） |
| R23 引用件在位 | `git cat-file -e HEAD:probes/319/v/…`⇄`.scratch/wisp/probes/319/v/…`；`git grep -n -F 'func TestBallPanelGesturesReachThePanelThread' HEAD`；`git log -1 --format=%B c ⇄ grep -c` | HEAD＋commit message | 量的 |
| R24 异腿旁证 | 外来 8 枚全包底片逐枚 `git cat-file blob`→`grep -E 'TestAC321'`＋PASS/FAIL 计数；`git ls-tree -r --name-only HEAD -- cmd|internal` 计数；逐件 `head -1` 抽 `//go:build` | 别的票的入库件 | 量的 |
| R26/R27 件身尺 | 被禁码位＝`od -An -tx1 -v|awk 相邻对`（★正控先行＝自造件读 2 枚）；CR 两把；`iconv` rc；末字节；占位符＝扫"未填记号族"四串英文/尖括号标记（★逐字命令只写在临时根 `sh/R27.sh`，⛔ 在本表写出＝本表会自撞那把尺报 1 命中） | 本腿交付件 | 量的 |
| C1 落笔自证 | `git add -- <字面 pathspec>`→`git commit -F D:\tmp\…\msg-1.txt -- <同 pathspec>`→`git log -1`＋`git show --name-only`＋入库 blob 复量 | 本腿写面 | 量的 |

## 本腿自抓（⛔ 藏，具名）

1. **R3b 第 h 段坏尺**＝`$(git cat-file blob …|wc -c)` 型命令替换吃掉行尾换行 ⇒ 编排者代提那四枚件各少报 1 字节／1 行；R4 改成先重定向到文件再量＝正数（3394/44、4255/66、2228/34、1352/22），与 `c7f813c7` 主题自报逐枚同号。旧尺件保留、⛔ 删。
2. **R13 v1 坏尺**＝`cut -c1-180` 按字节截中文 ⇒ `iconv` 报 `cannot convert`（＝公共块 台面四味④ 明令那一形，本腿犯了）；R13b 改 `LC_ALL=C.UTF-8`＋`awk substr` ⇒ `ICONV_RC=0`。
3. **两发 `git grep` 参数序 fatal**（R2 v1、R10b、R14 那三枚 `-e`／`-F` 落在 rev 之后）＝当场按"报 0 命中前先跑正控＋看 stderr"识破、⛔ 产生任何假读数。
4. **R1c 那枚 `drives NO executor on newPanelHolder`＝0 命中**＝尺射程所致（源串用 `%s` 占位符，运行时才插入 `newPanelHolder`），R10a 改锚 `drives NO executor on` 读出 1 ⇒ ⛔ 据此判腿报假。
