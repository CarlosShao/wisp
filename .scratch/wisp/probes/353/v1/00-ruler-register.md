# 票 353 终裁腿 `353-v1` 尺册（★非实现者·只读·裁决者≠实现者；本件＝`.scratch/wisp/probes/353/v1/00-ruler-register.md`）

判语主件＝`docs/evidence/s1/353-v1-verdict.md`；逐字读数＝`01-independent-readings.txt`；立票判据草案＝`02-ticket-draft-secC.md`。
尺脚本本体＝本目录 `sh/01|02|03|10|11|12|13|14|15|16|17|18|19|20|21|21b|22|23|24|25|26|27|28|30|31|32|40|41|42|50|51|52`（★全部可逐字重跑；`sh/21` 那枚因引号破损作废、改正尺＝`sh/21b`；`sh/42`＝入库 blob 件身尺（★正控先行＝自造 2 枚 CR 用同尺读回 2；本腿第一发把正控写成 `tr -cd '\r'|wc -l`＝数行数⇒读出 0＝一枚钝尺，已具名改正并留原件）；`sh/50|51|52`＝逐枚字面 pathspec 的提交尺形；另附读数原件 `90-truncated-nonzero-rows.txt`＝尺 R32 那 20 枚非零行的逐字底片；本腿交付件＝`00/01/02` 三枚＋`90` 一枚＋这批尺脚本＋主件 `docs/evidence/s1/353-v1-verdict.md`，`.go` 字节 **0 枚**）。
射程根：入库 blob＝`git cat-file blob <rev>:<path>` 导到 `D:\tmp\wisp\20261020\353v1\blob\`；底片＝`D:\tmp\wisp\20261020\353r1\logs\` 与 `…\353r1\{finalB,finalL,mut-final-1..3}\`（＝被裁腿自己的写面根，本腿**只读**）；名册＝`git log --all --format`。
本腿起手锚＝`51aa922e`（22:31）；⛔ 跑任何 `go test`／`go build`／`d22scan`（内含 `go run`）＝现势色那味具名归 `gw-1`。

| 尺 | 逐字命令形（简写） | 射程 | 量的/算的 | 答哪一问 |
|---|---|---|---|---|
| R01 `sh/01-anchors.sh` | `git log -1`／`git rev-parse <锚>:<path>`／`git diff --name-only a..b -- '*.go'`／逐枚 `git show --name-only` | 六笔提交元数据＋`internal/winsec` 三枚 blob | 量的 | 名册⇄锚⇄写面⇄`a59f5630` 基线串（`db077e30`/`f3c1b68e`/`63714a83`） |
| R02 `sh/02-blob-dump-diff.sh` | `git cat-file blob HEAD:…` ⇄ `a59f5630:…`＋`diff -u3` | `pinnedleaf_342_windows.go` | 量的 | 第 1 件＝乙本体逐字形状 |
| R03 `sh/03-diff-rest.sh` | 同上三枚件（`winsec.go`/`_unix.go`/`_fallback.go`） | 三枚 | 量的 | ⛔ 扩射程⇄unix/fallback 只同步签名 |
| R10 `sh/10-ac2-pair-recheck.sh` | `grep -c '^=== RUN'`／`'^--- PASS'`／`'^--- FAIL'`／`'^--- SKIP'`＋`grep -n` 分子⇄liveness⇄红句⇄GREEN SCOPE⇄CONTROL | `logs/final-dir-{B,L}.raw.txt`（旧批 `base353-s1/s2` 同尺并量） | 量的 | 第 3 件成对⇄第 2 件"绿行自带⛔ 覆盖声明" |
| R11 `sh/11-mutation-recheck.sh` | `git cat-file blob \| sha256sum`⇄两棵树在盘字节＋三发 `final-mut-*.raw.txt` 四数/分子/红名/正控 | 六枚件哈希＋三棵突变树 | 量的 | 第 3 件台面对齐⇄去牙三发 |
| R12 `sh/12-ac3-color-recheck.sh` | 整包四发四数＋本腿自取名集 `sort`/`comm -13`⇄`comm -23`＋container 两发＋`cat final-fourcounts.tsv` | `logs/final-full-*`／`final-col-*` | 量的 | 第 4 件逐名新增红⇄第二平台 |
| R13 `sh/13-teeth-control.sh` | `grep -E 'AC#2 CONTROL'`（mut-1）／`did not ENTER`（mut-2）／分子（mut-3）＋`grep -rh RESTORE_SHA256`＋四枚在盘哈希 | 三发底片＋`logs/*` 全量 | 量的 | 第 3 件牙的落点⇄还原凭据在不在库 |
| R14 `sh/14-ticket-face-check.sh` | `git cat-file blob HEAD:<票353>`＋`grep -cE '^- \[ \]'`/`'^- \[x\]'`＋六枚 token 打在**第 15 行**⇄第 14 行对照 | 票 353 入库 blob | 量的 | 第 2 件"票面写没写"＋框终量 |
| R15 `sh/15-scope-rulers.sh` | 逐笔 blob 串＋段取 sha256（自陈段）＋`_ = p.root.Remove(base)` 命中＋`func` 名集作差＋导出名并集作差＋旗标行⇄`fs.ErrExist`＋`privateFile` return 集合＋`:11-20` 段 diff＋包级 `var` 集合 | `pinnedleaf_342_windows.go`/`winsec.go`/`_unix`/`_fallback` 的 BASE⇄HEAD | 量的 | 第 1 件六味（含★自陈段那把尺与编排者那把不同结论） |
| R16 `sh/16-secC-recheck.sh` | `cat alt-summary.txt`＋逐 `alt-*.raw.txt` 抽 `provedReparse/flipErrs/truncated/deleted`＋`swing-*`＋**正控先行**（四枚串命中文件数⇄总命中） | `logs/alt-*`／`swing-*` | 量的 | 第 5 件①⇄④撤回对不对 |
| R17 `sh/17-familyA-numerators.sh` | `grep -oE 'form=(symlink-dir\|leaf-symlink\|junction) attempts=400 …'` 逐发＋`grep -rlE 'truncated=[1-9]'` | `alt-*`/`swing-*`/`final-full-*` | 量的 | 第 5 件④（先只用定向批，R19 扩到全批） |
| R18 `sh/18-landing-positions.sh` | 内容锚 `grep -n` 取行号＋`git log 0b75e40b..HEAD -- internal/winsec`＋门禁底片 `wc -c` 逐枚＋`cat gates.txt` | 仪器件/产码件/既有钉 HEAD blob＋`logs/gate-*` | 量的 | 第 2 件落点⇄第 4 件门禁凭据在位⇄winsec 未漂 |
| R19 `sh/19-full-table.sh` | 逐文件×逐 form 全表（`alt-*`/`swing-*`/`recheck*`/`mx-*`/`final-full-*`/`fullpair-*`/`pername-*`）＋非零行总表 | `logs/*` 全量分子行 | 量的 | 第 5 件④（★本把把腿自己的 `recheck-summary.txt` 也打了进去⇒非零数偏高，改正尺＝**R32**） |
| R20 `sh/20-misc-rulers.sh` | `grep -cF 'go func'/time.Sleep/atomic.'` 打在仪器件＋正控打在 `leafrace_342r2`；`fs.ErrExist` 两侧；剥注释可执行行作差 | 仪器件⇄`leafrace` 对照；三枚件 | 量的 | 第 1 件⛔ racer⇄笔④纯注释（★本把第 39 行引号破损⇒第 6/7/8 节作废，改正尺＝R21b） |
| R21b `sh/21b-misc.sh` | `diff \| grep -Ec`（正控 22 行）＋`deleted=[1-9]` 命中＋`cmd/wisp` 串命中＋**名册三把**（锚定/非锚定/单次导出 `index==1`）＋`git status --porcelain`＋族件名册 | `logs`＋`internal/winsec` | 量的 | 第 3/4 件⛔ 借⛔ 相减⇄名册⇄第 5 件③件数 |
| R22 `sh/22-ancestorwindow-and-swing.sh` | 那枚钉 BASE..HEAD 全 diff＋可执行行数⇄用例名集作差＋摆动名逐批发片色＋container 尾行＋三把锚命中 | `ancestorwindow_342r2` blob＋八批底片 | 量的 | 第 1 件重述②＝判据**变严**⇄摆动定因 |
| R23 `sh/23-mutation-trees.sh` | 三棵突变树 `sha256sum`⇄每棵相对 `finalL` 的 diff⇄三把新锚在四棵树命中数 | `mut-final-1..3`／`finalL`／`finalB` | 量的 | 第 3 件"落点⛔ 同"的独立证明 |
| R24 `sh/24-secC-semantics.sh` | `ancestorunlink_348` HEAD blob 里 `flipErrs/provedReparse` 声明⇄自增处⇄红句现读＋`git ls-tree` 数 118/120/342/348/353 结案态＋`grep -c '^## <号>'` 逐枚验台账号＋票 118⇄120 已批文字 | 两枚票面＋一枚仪器件＋台账 | 量的 | 第 5 件②（⛔ 属本票射程）＋⛔ 引空号 |
| R25 `sh/25-crossbatch-newred.sh` | 七批基线红名并集 ⇄ 七批落地红名并集（`comm -13`/`-23`/`-12`）＋两味精确行段打印＋四枚计数器使用点 | `logs/*` 全量＋仪器件 | 量的 | 第 4 件"跨批新增红＝0"的硬版＋第 3 件⛔ 合并 |
| R26/R27 | `git cat-file blob HEAD:<sh/35>`＋`cat alt-1-finalL.raw.txt`＋`342-v3-verdict.md:38-44` | 腿的尺本体⇄空底片一枚⇄母票签名节 | 量的 | 交替形定义＋⛔ 藏读数的复现 |
| R28 `sh/28-precedent-and-misc.sh` | 票 322 `AC#1` 那格逐字＋族件逐 pattern 计数＋`0b75e40b..HEAD` 逐枚列名＋读数件"合并/相减"字样＋`docs/evidence/s1/353*` 空闲尺＋本腿派单入库尺 | 三枚票面＋`internal/winsec` 名册＋`git log` | 量的 | 第 2 件先例形⇄第 5 件③件数⇄写面空闲 |
| R29 `sh/29-grep-322.sh` | `git grep -n '报窗口'`⇄`openleaf_322_other_test.go` 里弃权句现读 | 全仓 HEAD | 量的 | 第 2 件(a) 凭据 |
| R30 `sh/30-body-rulers.sh` | 件身尺＝`wc -c`/`wc -l`/`tr -cd '\r' \| wc -c`/`grep -o <双字节序列> \| wc -l`/`grep -cE` 打那三枚未完成标记词（本件⛔ 写出它们本身＝写了就自占一枚命中）/`tail -c 1 \| od -An -tx1`/`iconv` | 本腿四枚交付件（blob ⇄ 工作树两形） | 量的 | 纪律格（⛔ 对工作树量 CR＝钝尺那条） |
| R31 `sh/31-find-sect.sh` | 分节码位**定位**尺（逐行报那枚双字节序列所在行） | 本腿工作树件 | 量的 | 同上＝★真正控＝本把尺在本腿首版主件上读出 **1** 枚命中并定位到行（＝尺灵）、改正后复量 **0**；"打票 353 面＝0"只是另一发对照、⛔ 是正控 |
| R32 `sh/32-trunc-side-count.sh` | ★**改正尺**＝`for f in logs/*.raw.txt` 逐文件 `grep -oE 'form=… truncated=[0-9]+ deleted=[0-9]+'`＋`[ $tr -gt 0 ]`＋按文件名定侧 | **只打 `*.raw.txt`**（⛔ 混入腿自己的 `recheck-summary.txt`） | 量的 | 第 5 件④＝`truncated` 非零 **20/基线 9/落地 11**（作废 R19 那把报 24/10 的钝尺）；原件 `trunc-rows.txt` |

## 三条算的（⛔ 冒充量过）

1. "`flipErrs` 换档的**机制**＝乙之前 `Root.Remove(base)` 把攻击者那枚暂挂 link 删掉、racer 每发必重建"＝**算的**〔本腿现读 `pinnedleaf_342_windows.go:201-205` 那句自陈"Deleteat … opens the entry with `FILE_OPEN_REPARSE_POINT`"在场＋`winsec_windows.go` 未动；⛔ 跑计时、⛔ 读内核〕。量的只有两味＝两侧分子区间＋两侧无交叠。
2. "票面 `AC#2` bullet 那句'成对＝落地 `damage>0` ⇄ 基线 `entered=0`'起草于 `AC#1` 裁形之前"＝**算的**（＝同一段文字里 `AC#1` 的裁形时刻 19:2x 晚于该 bullet 的立票文字，且〔乙〕⇒落地侧分子必为 0）；量的只有两处分子行与两侧红名。
3. 第 2 件那条"先例本身就是写进判码、⛔ 写进票面 bullet"＝**算的判语**，其凭据（票 322 票面 `:30` 逐字、`openleaf_322_other_test.go:8` 逐字、"报窗口"三字在票面命中 0）＝量的。

## 枚数随 HEAD 漂

本件⛔ 写死任何"共几笔/几枚"的通则；名册尺只认逐枚 `%s` 起头（现量＝锚定 6／非锚定 6／`index==1` 6，正控 `^票 353 普查腿 353-a1`=6、`^票 353 终裁腿 353-v1`=0）。`0b75e40b..HEAD`＝本腿起手时读到 10 枚别腿笔（含 `gw-1` 笔①），随 HEAD 继续漂。
