# 308-v1 · 90-final（五格判语＋欠账＋越格自报＋对 A846 的三处不同意见）

判语者＝`308-v1`（⛔ 实现者；本腿⛔ 写过一个字节交付码，突变只⛔ 留现场）。
射程三层（工作树⇄HEAD blob⇄仓外导出树/容器）与单位＝`00-anchor.md`。件图：`10`=AC#1（含分母与"⛔ skip 换绿"）、`20`=AC#2 三件＋三族、`30`=AC#3 档位与 CI、`40`=AC#4 门禁与越界、`50`=AC#5 对账、`60`=归因单裁、`70`=越claim 扫描。

## 1. 五格判语（尺在盘上，逐格一行最强证据）

| 格 | 判语 | 最强那一行（可逐字复跑） |
|---|---|---|
| **AC#1** | **成立** | `git rev-parse` 七枚 ref 上 `pathresolver.go`＝`15349e28…`、`pathresolver_windows.go`＝`7b00b022…` **一字⛔ 变**；且执行位 `t.Skip` 计数 pre⇄post＝anchor `0⇄0`／junction `1⇄1`，新用例打印 `7 volume-free … 2 UNC spellings dropped by name [unc-admin unc-extended]` |
| **AC#2** | **成立** | 拔守卫那发：单用例 `rc=1` ＋ 逐字 `panic: runtime error: slice bounds out of range [:1] with length 0 [recovered, repanicked]`，**整包同尺从 `RUN=201 PASS=132` 掉到 `RUN=35 PASS=27`**（＝票的存在理由被当场重造）；还原 `sha256=f2c01307…` 与 r1 自报值同一；别名台面名册作差＝**＋1 枚** |
| **AC#3** | **附条件成立** | 档位全具名（core `:487`／windows `:797`／`:813` 那枚 `-run` 实测只匹配 1 个名字⇒⛔ 覆盖两枚拼法用例／winlive `:672` ⛔ 碰本包且本包 0 命中）；★新件**无 tag⇒ubuntu 真求值**，本腿补到"同版本真跑"＝`docker golang:1.27 (linux/amd64 go1.27.1)` 整包 `rc=0 RUN=169 PASS=108 FAIL=0 SKIP=1(ledgered) PANIC=0`、新用例 native `"/tmp/…"` 绿。**条件＝CI 跑机色⛔ 本机色**（`Q-308-v1-A`） |
| **AC#4** | **成立** | `GOFLAGS= go build ./...`=0／`go vet ./internal/risk/`=0／`sh scripts/d22scan.sh`=0（"clean - no D22 ban violations"），各件自落 `rc=N`；长度门 `VERDICT GREEN`（over-hat 57 ⇄ roster 57，worst=228）；格式名册窄＝1（`provenance.go`，CR=1115 B、`i/lf w/crlf`＝CRLF 假枚）、宽＝26（22 `.scratch`＋4 非），逐笔越界＝0 |
| **AC#5** | **成立** | 五笔对票 72/300/301/302/307 的 numstat **逐笔为空**；`git log -S'- [x] **\`AC#0\`' -- <票308>` 只命中 `0d2260eb`（编排者那笔）⇒ "⛔ 本腿翻框"由 commit 归因正面证实 |

## 2. 欠读数（具名＋归谁＋⛔ 拿什么顶）

- **`Q-308-v1-A`**（归编排者，与 `A813` 那批同批取）：`TestClassifySpellingInvariance`／`TestPathResolverUNCAListDenied`／`TestVolumeInvariantSpellingsVolumeLess308` 在 `ci.yml:487`（ubuntu core）与 `ci.yml:797`（windows scope）两步的**跑机实跑色**。⛔ 用本腿 docker 或本机名册顶（本腿已把 r1 的 `Q-308-r1-A/B` 并进这一枚，B 的"本机 linux 色"部分已交，**只剩跑机**）。附带：`junction:157-159` 那句"…and on the windows CI gear"此刻＝⛔ 读数（`logs/v46-fsutil.txt` Access denied）。
- **`Q-308-v1-B`**（归票面，永久欠，除非换机器/换工具）：两枚 windows-tagged 件里**守卫那两支**（`len(dropped)>0` 的 `t.Logf`／`!haveShare` 的 `t.Fatalf`）**在本机能造的每一把台面上到达数＝0**（四发哨兵命中全 0）。⇒ "无卷标那一形下这枚**用例**以可解释的颜色结束"目前只有 (i) helper 输入面一枚钉住的载具＋(ii) 代码形状，⛔ 是"那枚用例真在无卷标台面上跑过"。⛔ 因此判 AC#1 不成立（AC#1 的硬约束与分母全满足），但**必须这样描述才诚实**。
- **`Q-308-v1-C`**（归编排者 A846 §3）：历史那发的**真成因未定**。`TMPDIR` 归因已被两路反证（stdlib 源＋env 网格），而记录里的形状（前导单反斜杠、无盘符）＝`filepath.Join("/", name)` 的产物，本腿网格量不出 `GetTempPath` 会交回那枚值 ⇒ ⛔ "已解释"，⛔ "⛔ 发生过"。
- **`Q-308-v1-D`**（尺子缺陷，归 A846/a1 续账，⛔ 翻框）：(i) `probes/308/a1/logs/probe-sweep-shapes.txt` 8 行右侧**全空**、`probe-sweep-gotmpdir.txt` 8 行只 2 行有正文、`probe-sweep-unset.txt` **0 字节** ⇒ ⛔ 再当"逐枚坐实"引；(ii) `r1/40-tags-and-gates.md:52` 用 `ls .scratch/wisp/issues/307*` 当"307 已 `-done`"的尺——那枚 slug 本来就⛔ 叫 `…not-sync-done.md`，这把尺**分辨⛔ 出**"名字里含 done"与"被标 done"（本仓 `-done.md` 结尾的票有 111 枚），结论按票面/台账采信，尺⛔ 立。
- **`Q-308-v1-E`**（措辞债，⛔ 判级影响）：`pathresolver_volumeless_spelling_308_test.go:19` "This is the single place the drive letter is taken." 缺射程 ⇒ 补"in this test surface"（普查读数见 70 件 #3）。
- **`Q-308-r1-C` 本腿已裁**（⛔ 新开欠账）：**需要它自己那一刀改写**，建议句逐字在 `70-overclaim-audit.md` 末节；归编排者派刀（票面"登记、⛔ 顺手改对"＝⛔ 本腿改，本腿⛔ 已改）。

## 3. 本腿**没**做的事（禁线逐条负尺；⛔ 这句读成"做过的事"）

⛔ 翻任何框、⛔ 改任何 Status 行、⛔ 改票面原句；⛔ 动产码、⛔ 动产码测试外的任何交付件（突变只在同一脚本内活、三次 `cmp`＋`sha256` 当场还原）；⛔ push、⛔ `git add -A`/`.`、⛔ `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；⛔ 在 `probes/**` 造 `.go`（分母探针那枚 `.go` 只在**仓外导出树**里）；⛔ 跑 `cmd/wisp` 测试包；⛔ 在仓根丢东西；⛔ 删任何文件（7 MB 交叉编译二进制＝**迁出**仓，⛔ 删，见 `logs/v23-binary-relocated.txt`）。

## 4. 越格自报（⛔ 等下一位核，⛔ 自我盖章）

1. **写了 7 枚仓外 `.go`/可执行物**：`mktemp -d` 里的探针程序（`main.go`＋`go.mod`＋exe）与导出树里的 `internal/risk/zz_v1_denominator_probe_test.go`。**都在仓外** ⇒ 工作树/HEAD/`probes/**` 的 tracked-`.go` 分母**纹⛔ 动**（宽尺 26 枚里⛔ 出现它们，且 `git status --porcelain` 干净可复跑）。判语用了一把"在仓外拷贝里加探针"的尺＝超出"只读＋自家件"的窄形，**为的是把"名册⛔ 变"升级成"拼法枚数⛔ 变"**；若编排者认为那一把尺越格，撤回它只影响 10 件 §4 的"量的那一半"（算的那一半⛔ 依赖它）。
2. **跑了一次 `docker run golang:1.27`**（2 发，只打 `./internal/risk/`）：这是**新仪器**，票面与 A846 都⛔ 提过它。它把 AC#3 从"⛔ 本机顶替"推进到"本机以外的一枚真 linux 读数"，⛔ 等于 CI 跑机色（条件仍挂着）。⚠ 它需要 `MSYS_NO_PATHCONV=1` 才 mount 得动，且那一发在容器里下过 `go-toml` 模块（stderr 有 `go: downloading` 一行，⛔ 进了枚数——`portable-tests.sh` 里 ticket 250 那族坑的正身，本腿按"stdout 才算分母"处理）。
3. **`sed -i` 直接改过跟踪文件三次**（同一脚本内、trap 还原）：票面允"transient mutation"，但"⛔ 留脏"靠的是脚本没死在半路——本腿**每次都当场复跑** `cmp`＋`sha256sum`＋`git status --porcelain`，三发结果均在 `logs/teeth*-cmp.txt`（`cmp_final=0`、hash 恒等 `f2c01307…`）。⛔ 把"应该没事"当读数。
4. **⛔ 越界⛔ 顺手⛔ 判**了 a1/`r1` 的账目（§2 的 D/E 两枚＋70 件三句）——都只**具名登记**，⛔ 动它们一行字。
5. 一处本腿自己的**无效读数主动留在盘上**：`teeth2-M3`（突变没落地、`rc=0` 却看着像"断言宽松"）——`logs/teeth2-M3-landed.txt` 与 `logs/teeth3-M3b-landed.txt` 并排，就是给下一位看"⛔ 先证落地就会读到假结论"。

## 5. 对 A846（票 §编排者裁定）的三处不同意见

1. **§1 "两把台面"的更正方向对、口径仍宽**：`TMP=/tmp` 被补的是**当前盘符**（本腿 cmd 那发＝`C:\`，腿与编排者那发＝`D:\tmp`，差别在 cwd 在哪个盘）。⇒ "⛔ 造得出无卷标 temp"的**结论**成立，**理由**要换成 60 件三.2 那两条（补盘符＋UNC ⛔ 可创建）；否则换台 cwd 在别的盘的机器就会以为尺变了。
2. **§3 "腿的在场尺逐枚坐实"⛔ 完全坐实**：引的那一枚（`probe-T1-fwd-slash.txt`）确有正文；同批的三枚 sweep 件是空壳（`Q-308-v1-D`）。
3. **§2 ⓑ 那句"如果守卫只对 env 形状起作用，它就是一枚永远⛔ 被求值的分支"——本腿把它量成了实测而不止是理由**：`docker` 那发证明输入面载具在 linux 真求值；而 windows 那一侧的分支**到达数＝0** 也是实测（`Q-308-v1-B`）。⇒ 裁语⛔ 改，只是它现在**两边都有读数**了，A846 该把"理由"升级成"两枚读数＋一枚永久欠"。
