# 裁决表 · 票 308（拼法不变量 `vol[:1]` 零守卫让整包静默少跑）— 验收腿 `308-v1`

- 裁决者＝`308-v1`（非实现者；⛔ 写过一个字节交付码）；被裁＝`308-r1` 五笔 `510c95ca`→`1ad48333`→`ee0a93e4`→`79242df7`→`1270e36e`
- 起手 HEAD＝`0d2260eb`（`dev`，共享工作树）；读数全部本腿现跑，⛔ 转抄 `probes/308/{a1,r1}`
- 射程三层（每格都标）＝工作树 ⇄ HEAD blob ⇄ 仓外导出树/容器；原始读数＝`.scratch/wisp/probes/308/v1/**`
- `AC#0` 已由编排者 13:1x 现跑闭合（`git log -S` 归因＝`0d2260eb`），本腿⛔ 重开、⛔ 动任何框

| 格 | 判语 | 本格尺（可逐字复跑）与那一行最强证据 |
|---|---|---|
| **AC#1** 守卫落点＋三条硬约束 | **成立**（两处具名偏移见下） | 产码 blob 七枚 ref 恒等：`pathresolver.go=15349e28…`、`pathresolver_windows.go=7b00b022…`；`git diff --stat 510c95ca^ -- <两枚>`＝空。守卫行逐字仍在前：`:262 if len(drive)==2 && drive[1]=='$' …` → `:263 ToUpper(drive[:1])`。执行位 `t.Skip(` 计数 pre⇄post＝anchor `0⇄0`／junction `1⇄1`（五笔新增的 3 行 `t.Skip` 字样全是注释）。新用例打印分母＝`7 volume-free spellings … 2 UNC spellings dropped by name [unc-admin unc-extended]`（工作树 windows 台面与 docker linux 台面同一句）。拼法级：可达台面 **9⇄9**（改前＝blob map 字面量数出来的；改后＝仓外探针跑出来的 `keep=9 dropped=[]`），喂无卷标字面量才 9→7 |
| **AC#2** 反形敏感 (a)(b)(c) | **成立** | (a) 拔守卫：`rc=1`＋逐字 `panic: runtime error: slice bounds out of range [:1] with length 0 [recovered, repanicked]`，**整包同尺 `RUN=201 PASS=132` → `RUN=35 PASS=27`**（＝票的坏形状被当场重造）。(b) `cmp_after_restore_rc=0`、`sha256=f2c0130702384bfd…070e` 三发恒等＝r1 自报值；复绿 `rc=0`。(c) 别名台面（`RISKAT~1\RUNNER~1`，`os.TempDir()` 当场证用上）名册 `132→133`，作差＝**只 ＋1 枚** `PASS TestVolumeInvariantSpellingsVolumeLess308`；改后⇄工作树名册 `diff` 空。★反空转两枚：M2（helper ⛔ 再去掉）→`dropped=[], want exactly [unc-admin unc-extended]` ⛔panic 红；M3b（抽掉 `"lower"`，落地尺 1→0）→`volume-free spelling "lower" went missing` 红 |
| **AC#3** 分母与档位 | **附条件成立**（条件＝`Q-308-v1-A`） | 档位具名：`ci.yml:419/487` core(ubuntu)→`portable-tests.sh:240` 含 `./internal/risk/...`；`:533/797` windows→`:251` 含 `./internal/risk/`；`:813` 那步实测 `go test -list 'TestPathResolverJunctionWindows'` **只匹配 1 个名字**⇒⛔ 覆盖两枚拼法用例；`:672` `go vet -tags winlive` 只打 `cmd/wisp`+`internal/ball`，`internal/risk` 里 `winlive` **0 命中**，`go vet -tags winlive ./internal/risk/`＝0。★新件**⛔ 任何 build tag**（`grep -n go:build` 只命中前两枚）⇒ ubuntu 档真求值：`GOOS=linux go vet`=0、`GOOS=linux go test -c`=0、`docker golang:1.27`（容器自报 `linux amd64 go1.27.1`）整包 **`rc=0 RUN=169 PASS=108 FAIL=0 SKIP=1 PANIC=0`**，那枚 SKIP＝`TestC26RewrittenSyncRootDoesNotDisarmSuspectNet`＝ledger 已在册；改前树同尺 `RUN=168 PASS=107`⇒名册作差＋1。**条件**：票 `AC#3` 原句"⛔ 用本机跑过顶替 CI 色"——docker/本机皆属"本机"，跑机那两发色⛔ 取（⛔ push） |
| **AC#4** 门禁与越界 | **成立** | `GOFLAGS= go build ./...`＝0、`go vet ./internal/risk/`＝0、`sh scripts/d22scan.sh`＝0（"clean - no D22 ban violations"），三枚件各**自落 `rc=N` 行**（build/vet 的 0 字节 stdout 是绿的本形，已把 rc 追进同一枚件）。格式名册两把并排：窄 `gofmt -l internal/risk`＝**1**（`provenance.go`，CR=**1115 B**、`i/lf w/crlf attr/text eol=lf`、`core.autocrlf=true` ⇒ CRLF 假枚，⛔ 本票引入）；宽 `git ls-files '*.go' \| xargs gofmt -l`＝**26**（22 `.scratch` 变异拷贝＋4 非：`cmd/wisp/models.go`/`internal/agent/approval/pending_read.go`/`internal/agent/tools.go`/`internal/risk/provenance.go`）；三枚交付件 CR＝**0**、`w/lf`。⚠ 与 r1 的 27/23 差一枚在 `.scratch` 桶＝别腿在两发之间动了 tracked `.go`（⛔ 本腿造 `.go`）。逐笔越界＝**0**。长度门 `--with-self-test`＝**GREEN**（`over-hat=57 ⇄ roster=57`、worst=228、⛔ 动别人的名册） |
| **AC#5** 对账 | **成立** | 五笔对票 72/300/301/302/307 的 `git show --numstat -- <五枚票>` **逐笔空**；票 308 总共只被 `ffcb21e9`（立票）与 `0d2260eb`（编排者盖章）两笔动过 ⇒ r1"⛔ 翻框、⛔ commit 票文件"由 **commit 归因**证实（⛔ 只是自报）。"同文件不同件事"有 hunk 级边界（`1ad48333` 在 junction 只 1 枚 hunk `@@ -149,7 +149,18 @@`，⛔ 碰 `f278c57f` 一行）；300/301/302 的搬运活（`portable-tests.sh`/`ci.yml`/`runtests.sh`）⛔ 被碰 |

## 两处具名偏移（⛔ 致命，⛔ 抹）

- ⓐ 票面只许可去掉 `unc-admin` **一枚**，交付把 `unc-admin`＋`unc-extended` **两枚**一起去掉。判：⛔ 越界（`unc-extended` 同由 `vol` 构造，无卷标时同样造不出；留它＝留一枚必然红的假拼法），且两枚都**按名打印** ⇒ 落在"具名去掉构造⛔ 出的拼法"同一件事里。
- ⓑ `TestClassifySpellingInvariance` 的 `if len(dropped) > 0` 一支与 `TestPathResolverUNCAListDenied` 的 `t.Fatalf(…cannot be constructed)` 一支，**在可构造台面上到达数＝0**（四发哨兵命中 0）。⇒ AC#1 那句"无卷标那一形下这枚**用例**以可解释的颜色结束"＝形状保证＋输入面载具，**⛔ 是"那枚用例真在无卷标台面上跑过"**。诚实口径已写进 `Q-308-v1-B`。

## 欠账（具名，⛔ 本腿代答）

`Q-308-v1-A`＝三枚用例在 `ci.yml:487`/`:797` 的**跑机**实跑色（归编排者，与 `A813` 同批；并入 r1 的 `Q-308-r1-A/B`）。
`Q-308-v1-B`＝两枚 windows-tagged 件的守卫分支⛔ 可达夹具（永久欠，除非换机器/换工具）。
`Q-308-v1-C`＝历史那发的**真成因未定**（`TMPDIR` 归因已由 stdlib 源＋env 网格两路反证；记录形状＝`filepath.Join("/", name)` 的产物，本腿网格量 `GetTempPath` ⛔ 交回该值）。
`Q-308-v1-D`＝两把**分辨⛔ 出结论**的尺：a1 的 `probe-sweep-*`（8 行右侧全空／只 2 行有正文／0 字节）⛔ 再当"逐枚坐实"；r1 的 `ls …/307*` 当"`-done` 已标"（该 slug 本来就⛔ 叫 `…not-sync-done.md`，本仓 `-done.md` 结尾 111 枚）。
`Q-308-v1-E`＝新件 `:19` "This is the single place the drive letter is taken." 缺射程（补"in this test surface"；全仓求值位普查＝2：产码 `:263` 有守卫＋helper `:27`）。
`Q-308-r1-C` **本腿已裁**＝`anchor:100-101` "every other spelling is unconditional, so the test never quietly stops measuring" 现在**被同函数代码正面否定**（两枚 UNC 拼法已变条件）⇒ **需要它自己那一刀改写**，建议替换句逐字在 `probes/308/v1/70-overclaim-audit.md` 末节；刀⛔ 归本腿（票面"登记、⛔ 顺手改对"）。

## 越claim 扫描结果

⛔ 一句把这一族写成"已防住"（扫描词 `prevented|never again|杜绝|根除|再也不会|cannot panic|fully fixed|guarantee` 命中＝0）。三句要具名：`anchor:100-101`（⛔ 立，见上）、`anchor:108-110` "the remaining seven volume-free spellings **still run** -- the case ends with a verdict"（外推，helper 那半已量、用例那半⛔ 可达）、`junction:157-159` "…and on the **windows CI gear**"（⛔ 读数，`fsutil` 本机 Access denied）⇒ 都归 `Q-308-v1-A/B/E` 与 `Q-308-r1-C`，⛔ 够把绿名册推翻，⛔ 够让本腿判 AC#2 不成立。
